package service

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"douyin/util"
	"log"
	"strconv"
	"strings"
)

type FavoriteServiceImpl struct {
}

func (fsi FavoriteServiceImpl) CountFavoritesByToVideoId(toVideoId int64) (int64, error) {
	// <toVideoId, fromUserId>
	toVideoIdStr := strconv.FormatInt(toVideoId, 10)
	keyStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + toVideoIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return 0, err
		}
		// 存在这个键，直接查
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count the number of value by key in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		// 根据toVideoId查找该用户点赞的视频
		userIds, err := dao.FindFavoriteUserIdsByToVideoId(toVideoId)
		if err != nil {
			log.Println("failed to find users who love this video")
			return 0, err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, userId := range userIds {
			userIdStr := strconv.FormatInt(userId, 10)
			vStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + userIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return 0, err
			}
		}
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count the number of value by key in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	}
}

func (fsi FavoriteServiceImpl) CheckFavoriteByBothId(fromUserId, toVideoId int64) (bool, error) {
	// <fromUserId, toVideoId>
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	toVideoIdStr := strconv.FormatInt(toVideoId, 10)
	keyStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + toVideoIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return false, err
		}
		// 存在这个键，直接查
		exist, err := redis.RedisCli.SIsMember(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to check value in redis", err)
			return false, err
		}
		return exist, nil
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return false, err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return false, err
		}
		// 根据fromUserId查找该用户点赞的视频
		videoIds, err := dao.FindFavoriteVideoIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find user's favorite videos")
			return false, err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, videoId := range videoIds {
			videoIdStr := strconv.FormatInt(videoId, 10)
			vStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + videoIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return false, err
			}
		}
		exist, err := redis.RedisCli.SIsMember(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to check value in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return false, err
		}
		return exist, nil
	}
}

// AddFavoriteByBothId 分成两部分，一个是<fromUserId, toVideoId>，一个是<toVideoId, fromUserId>
// 对于其中一个，先判断redis里有没有存在键
// 如果存在键，直接添加值，并把信息发送给mq
// 如果不存在，先添加键，设置过期时间，添加默认值（防止redis删除最后一个值后，key为空，又从数据库中加载数据，导致与真实情况不一致）
// 再查询数据库应有的信息，一个一个添加进去，最后再添加现在的信息，通知mq
func (fsi FavoriteServiceImpl) AddFavoriteByBothId(fromUserId, toVideoId int64) error {
	// <fromUserId, toVideoId>
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	toVideoIdStr := strconv.FormatInt(toVideoId, 10)
	keyStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + toVideoIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return err
		}
		// 存在这个键，直接添加值，通知mq
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			return err
		}
		rabbitmq.Favoritemq.Produce(util.MQLikeType, fromUserIdStr, toVideoIdStr)
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 根据fromUserId查找该用户点赞的视频
		videoIds, err := dao.FindFavoriteVideoIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find user's favorite videos")
			return err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, videoId := range videoIds {
			videoIdStr := strconv.FormatInt(videoId, 10)
			vStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + videoIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		rabbitmq.Favoritemq.Produce(util.MQLikeType, fromUserIdStr, toVideoIdStr)
	}

	// <toVideoId, fromUserId>
	keyStr, valueStr = valueStr, keyStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return err
		}
		// 存在这个键，直接添加值
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			return err
		}
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 根据fromUserId查找该用户点赞的视频
		userIds, err := dao.FindFavoriteUserIdsByToVideoId(toVideoId)
		if err != nil {
			log.Println("failed to find users who love this video")
			return err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, userId := range userIds {
			userIdStr := strconv.FormatInt(userId, 10)
			vStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + userIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
	}
	return nil
}

// DeleteFavoriteByBothId 逻辑和上个函数几乎一摸一样
func (fsi FavoriteServiceImpl) DeleteFavoriteByBothId(fromUserId, toVideoId int64) error {
	// <fromUserId, toVideoId>
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	toVideoIdStr := strconv.FormatInt(toVideoId, 10)
	keyStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + toVideoIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return err
		}
		// 存在这个键，直接添加值，通知mq
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove key in redis", err)
			return err
		}
		rabbitmq.Favoritemq.Produce(util.MQDisLikeType, fromUserIdStr, toVideoIdStr)
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 根据fromUserId查找该用户点赞的视频
		videoIds, err := dao.FindFavoriteVideoIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find user's favorite videos")
			return err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, videoId := range videoIds {
			videoIdStr := strconv.FormatInt(videoId, 10)
			vStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + videoIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		rabbitmq.Favoritemq.Produce(util.MQDisLikeType, fromUserIdStr, toVideoIdStr)
	}

	// <toVideoId, fromUserId>
	keyStr, valueStr = valueStr, keyStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return err
		}
		// 存在这个键，直接添加值
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove key in redis", err)
			return err
		}
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		// 根据fromUserId查找该用户点赞的视频
		userIds, err := dao.FindFavoriteUserIdsByToVideoId(toVideoId)
		if err != nil {
			log.Println("failed to find users who love this video")
			return err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, userId := range userIds {
			userIdStr := strconv.FormatInt(userId, 10)
			vStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + userIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
	}
	return nil
}

func (fsi FavoriteServiceImpl) FindFavoriteVideoIdsByFromUserId(fromUserId int64) ([]int64, error) {
	// <fromUserId, toVideoId>
	videoIds := make([]int64, 0)
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	keyStr := "favorite" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	// 存在这个键
	if n > 0 {
		if err != nil {
			log.Println("failed to query key in redis", err)
			return videoIds, err
		}
		// 存在这个键，直接取值
		vsStr, err := redis.RedisCli.SMembers(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to get members in redis", err)
			return videoIds, err
		}
		for _, vStr := range vsStr {
			//log.Println(vStr)
			if vStr == util.RedisDefaultValue {
				continue
			}
			videoId, _ := strconv.ParseInt(strings.Split(vStr, util.RedisSplit)[2], 10, 64)
			videoIds = append(videoIds, videoId)
		}
		return videoIds, nil
	} else { // 不存在这个键
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return videoIds, err
		}
		// 设置过期时间，1天到3天内随机，防止缓存雪崩
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiring time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return videoIds, err
		}
		// 根据fromUserId查找该用户点赞的视频
		videoIdsDao, err := dao.FindFavoriteVideoIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find user's favorite videos")
			return videoIds, err
		}
		// 一个一个添加，如果又一个添加失败，则删除键（为了保证redis内容和数据库同步，一旦添加失败就不同步了）
		for _, videoIdDao := range videoIdsDao {
			videoIdStr := strconv.FormatInt(videoIdDao, 10)
			vStr := "favorite" + util.RedisSplit + "toVideoId" + util.RedisSplit + videoIdStr
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add key in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return videoIds, err
			}
		}
		return videoIdsDao, nil
	}
}

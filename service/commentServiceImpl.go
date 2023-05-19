package service

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"douyin/util"
	"log"
	"strconv"
	"strings"
	"time"
)

type CommentServiceImpl struct {
}

// CountCommentsByToVideoId 先查询redis是否存在toVideId键
// 如果存在，则直接返回数量
// 如果不存在，则先添加默认值，再设置过期时间，然后查数据库，一个一个添加进来，最后返回数量
func (csi CommentServiceImpl) CountCommentsByToVideoId(toVideoId int64) (int64, error) {
	keyStr := "comment" + util.RedisSplit + "toVideoId" + util.RedisSplit + strconv.FormatInt(toVideoId, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about comment in redis", err)
		return 0, err
	}
	if n > 0 {
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count the number of key about comment in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	} else {
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add default value in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about comment in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		ids, err := dao.FindCommentIdsByToVideoId(toVideoId)
		if err != nil {
			log.Println("failed to find all comments of the video", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		for _, id := range ids {
			tmpValueStr := "comment" + util.RedisSplit + "id" + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, tmpValueStr).Result()
			if err != nil {
				log.Println("failed to add value about comment in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return 0, err
			}
		}
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count the number of comments of the video in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	}
}

// AddCommentByAll 直接对数据库操作，返回commentId，更新redis，整个函数返回commentId
// redis维护两个数据结构<toVideoId, commentIds> <commentId, videoId>
func (csi CommentServiceImpl) AddCommentByAll(fromUserId, toVideoId int64, content string, createDate time.Time) (int64, error) {
	// <toVideoId, commentId>
	id, err := dao.AddCommentByAll(fromUserId, toVideoId, content, createDate)
	if err != nil {
		log.Println("failed to add comment to database", err)
		return util.NotExistId, err
	}
	// 查看redis有没有toVideoId的键
	keyStr := "comment" + util.RedisSplit + "toVideoId" + util.RedisSplit + strconv.FormatInt(toVideoId, 10)
	valueStr := "comment" + util.RedisSplit + "id" + util.RedisSplit + strconv.FormatInt(id, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key in redis", err)
		return util.NotExistId, err
	}
	if n > 0 { //存在这个键，直接添加值
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			return util.NotExistId, err
		}
	} else { // 不存在这个键，先添加默认值，设置过期时间，再从数据库一个一个添加，不需要最后再添加现在的评论id
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return util.NotExistId, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return util.NotExistId, err
		}
		ids, err := dao.FindCommentIdsByToVideoId(toVideoId)
		if err != nil {
			log.Println("failed to find comments of video", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return util.NotExistId, err
		}
		for _, id := range ids {
			valueStr = "comment" + util.RedisSplit + "id" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
			if err != nil {
				log.Println("failed to add value in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return util.NotExistId, err
			}
		}
	}
	keyStr, valueStr = valueStr, keyStr
	// <id, toVideoId> 数据结构是1对1的，所以不添加默认值
	_, err = redis.RedisCli.Set(redis.Ctx, keyStr, valueStr, util.RandomDuration()).Result()
	if err != nil {
		redis.RedisCli.Del(redis.Ctx, keyStr)
		log.Println("failed to set key-value in redis", err)
		return util.NotExistId, err
	}
	return id, nil
}

// DeleteCommentById 先通过另一个redis数据结构查找toVideoId，如果找到直接获取，然后删除该key，如果找不到则去数据库查
// 然后通过toVideoId查找redis是否存在key，如果存在，则删除对应的id，然后mq通知数据库
// 如果不存在，则先查数据库，添加默认值，设置过期时间，ids加到redis里，然后删除对应的id，最后mq通知数据库
// 此处如果redis找不到key，不添加<id, toVideoId>这种数据结构
// 函数返回删除的comment
func (csi CommentServiceImpl) DeleteCommentById(id int64) error {
	keyStr := "comment" + util.RedisSplit + "id" + util.RedisSplit + strconv.FormatInt(id, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about comment in redis", err)
		return err
	}
	var toVideoId int64
	var toVideoIdStr string

	// 这段代码只是获得toVideoId
	if n > 0 { // 如果存在id这个key，直接获取videoId，然后把这个key删除
		valueStr, err := redis.RedisCli.Get(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to get value about comment in redis", err)
			return err
		}
		toVideoIdStr = strings.Split(valueStr, util.RedisSplit)[2]
		toVideoId, _ = strconv.ParseInt(toVideoIdStr, 10, 64)
		_, err = redis.RedisCli.Del(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to delete key about comment in redis", err)
			return err
		}
	} else { // 如果不存在id这个key，直接查数据库，不用加到redis来
		toVideoId, err = dao.FindCommentToVideoIdById(id)
		if err != nil {
			log.Println("failed to find video to which comment belongs", err)
			return err
		}
		toVideoIdStr = strconv.FormatInt(toVideoId, 10)
	}

	//至此已经获得toVideoId
	valueStr := keyStr
	keyStr = "comment" + util.RedisSplit + "toVideoId" + util.RedisSplit + toVideoIdStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about comment in redis", err)
		return err
	}
	if n > 0 { // 如果存在toVideoId这个key
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about comment in redis", err)
			return err
		}
		rabbitmq.Commentmq.Produce(util.MQCancelCommentType, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, strconv.FormatInt(id, 10))
	} else { // 不存在这个key
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add key about comment in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about comment in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		ids, err := dao.FindCommentIdsByToVideoId(toVideoId)
		if err != nil {
			log.Println("failed to find all comments of the video", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		for _, id := range ids {
			tmpValueStr := "comment" + util.RedisSplit + "id" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, tmpValueStr).Result()
			if err != nil {
				log.Println("failed to add value about comment in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about comment in redis", err)
			return err
		}
		rabbitmq.Commentmq.Produce(util.MQCancelCommentType, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, util.MQEmptyValue, strconv.FormatInt(id, 10))
	}
	return nil
}

package service

import (
	"douyin/dao"
	"douyin/middleware/rabbitmq"
	"douyin/middleware/redis"
	"douyin/util"
	"log"
	"strconv"
)

type RelationServiceImpl struct {
}

func (rsi RelationServiceImpl) CountRelationsByFromUserId(fromUserId int64) (int64, error) {
	// <from_user_id, to_user_id> 表示谁关注了谁
	keyStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + strconv.FormatInt(fromUserId, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return 0, err
	}
	if n > 0 { //存在
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count the number of values about relation in redis", err)
			return 0, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		toUserIds, err := dao.FindRelationToUserIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find idols of the user in database", err)
			return 0, err
		}
		for _, id := range toUserIds {
			vStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return 0, err
			}
		}
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count values about relation in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	}
}

func (rsi RelationServiceImpl) CountRelationsByToUserId(toUserId int64) (int64, error) {
	// <from_user_id, to_user_id> 表示谁关注了谁
	keyStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(toUserId, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return 0, err
	}
	if n > 0 { //存在
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count the number of values about relation in redis", err)
			return 0, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return 0, err
		}
		fromUserIds, err := dao.FindRelationFromUserIdsByToUserId(toUserId)
		if err != nil {
			log.Println("failed to find idols of the user in database", err)
			return 0, err
		}
		for _, id := range fromUserIds {
			vStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return 0, err
			}
		}
		cnt, err := redis.RedisCli.SCard(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println("failed to count values about relation in redis", err)
			return 0, err
		}
		return cnt - 1, nil
	}
}

func (rsi RelationServiceImpl) CheckRelationByBothId(fromUserId, toUserId int64) (bool, error) {
	// <from_user_id, to_user_id> 表示谁关注了谁
	keyStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + strconv.FormatInt(fromUserId, 10)
	valueStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(toUserId, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return false, err
	}
	if n > 0 { //存在
		res, err := redis.RedisCli.SIsMember(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to check value about relation in redis", err)
			return false, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return false, err
		}
		return res, nil
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return false, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return false, err
		}
		toUserIds, err := dao.FindRelationToUserIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find idols of the user in database", err)
			return false, err
		}
		for _, id := range toUserIds {
			vStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return false, err
			}
		}
		res, err := redis.RedisCli.SIsMember(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to count values about relation in redis", err)
			return false, err
		}
		return res, nil
	}
}

// AddRelationByBothId 获取用户id，先通知mq
// 然后更新redis （需要更新过期时间）（redis的更新不用goroutine，怕数据不一致）
// redis存储三种数据结构 <from_user_id, to_user_id> <to_user_id, from_user_id> <from_user_id, to_user_id> （最后一个表示朋友关系）

// AddRelationByBothId 对于第一个数据结构 <from_user_id, to_user_id>
// 首先看有没有键，如果有，直接添加，然后改过期时间
// 如果没有，先添加默认值，然后设置过期时间，从数据库中查信息，一个一个添加，然后添加最新消息

// AddRelationByBothId 对于第二个数据结构 <to_user_id, from_user_id>
// 首先看有没有键，如果有，直接添加，然后改过期时间
// 如果没有，先添加默认值，然后设置过期时间，从数据库中查信息，一个一个添加，然后添加最新信息

// AddRelationByBothId 对于第三个数据结构 <from_user_id, to_user_id> （表示朋友关系）
// 由于是双向的，所以又按key分为两个
// 假设A关注了B
// 先查B是否关注了A（这个用一个函数，内部也是先查redis，查不到再查数据库）
// 如果B没有关注A，直接退出
// 如果B关注了A，则确定朋友关系
// 先查key是否又A，如果有，则直接添加B的值
// 如果没有，则先添加默认值，设置过期时间，查数据库，然后一个一个添加，最后添加最新信息
// 对于另一个，先查key是否有B，如果有，则直接添加A的值
// 如果没有，则先添加默认值，设置过期时间，查数据库，然后一个一个添加，最后添加最新信息
func (rsi RelationServiceImpl) AddRelationByBothId(fromUserId, toUserId int64) error {
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	toUserIdStr := strconv.FormatInt(toUserId, 10)
	rabbitmq.Relationmq.Produce(util.MQFollowType, fromUserIdStr, toUserIdStr)

	// <from_user_id, to_user_id> 表示谁关注了谁
	keyStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + toUserIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		toUserIds, err := dao.FindRelationToUserIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find idols of the user in database", err)
			return err
		}
		for _, id := range toUserIds {
			vStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
	}

	// <to_user_id, from_user_id> 表示谁被谁关注
	keyStr, valueStr = valueStr, keyStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		fromUserIds, err := dao.FindRelationFromUserIdsByToUserId(toUserId)
		if err != nil {
			log.Println("failed to find fans of the user in database", err)
			return err
		}
		for _, id := range fromUserIds {
			vStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
	}

	// <from_user_id, to_user_id> 表示某人的朋友
	res, err := dao.CheckRelationByBothId(toUserId, fromUserId)
	if err != nil {
		log.Println("failed to check relation in database", err)
		return err
	}
	if !res {
		return nil
	}
	keyStr = "relationFriend" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr = "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + toUserIdStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		friendIds, err := dao.FindRelationFriendIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find friends of the user in database", err)
			return err
		}
		for _, id := range friendIds {
			vStr := "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
	}

	keyStr = "relationFriend" + util.RedisSplit + "fromUserId" + util.RedisSplit + toUserIdStr
	valueStr = "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + fromUserIdStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		friendIds, err := dao.FindRelationFriendIdsByFromUserId(toUserId)
		if err != nil {
			log.Println("failed to find friends of the user in database", err)
			return err
		}
		for _, id := range friendIds {
			vStr := "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			return err
		}
	}
	return nil
}

// 获取用户id，先通知mq
// 然后更新redis （需要更新过期时间）（redis的更新不用goroutine，怕数据不一致）
// redis存储三种数据结构 <from_user_id, to_user_id> <to_user_id, from_user_id> <from_user_id, to_user_id> （最后一个表示朋友关系）

// 对于第一个数据结构 <from_user_id, to_user_id>
// 首先看有没有键，如果有，直接删除，然后改过期时间
// 如果没有，先添加默认值，然后设置过期时间，从数据库中查信息，一个一个添加，然后删除最新消息

// 对于第二个数据结构 <to_user_id, from_user_id>
// 首先看有没有键，如果有，直接删除，然后改过期时间
// 如果没有，先添加默认值，然后设置过期时间，从数据库中查信息，一个一个添加，然后删除最新信息

// DeleteRelationByBothId 对于第三个数据结构 <from_user_id, to_user_id> （表示朋友关系）
// 由于是双向的，所以又按key分为两个
// 假设A取关了B
// 先查B是否关注了A，如果没有，直接退出
// 如果有
// 先查key是否有A，如果有，则直接删除B的值
// 如果没有，则先添加默认值，设置过期时间，查数据库，然后一个一个添加，最后删除B
// 对于另一个，先查key是否有B，如果有，则直接删除A的值
// 如果没有，则先添加默认值，设置过期时间，查数据库，然后一个一个添加，最后删除B
func (rsi RelationServiceImpl) DeleteRelationByBothId(fromUserId, toUserId int64) error {
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	toUserIdStr := strconv.FormatInt(toUserId, 10)
	rabbitmq.Relationmq.Produce(util.MQUnfollowType, fromUserIdStr, toUserIdStr)

	// <from_user_id, to_user_id> 表示谁关注了谁
	keyStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + toUserIdStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		toUserIds, err := dao.FindRelationToUserIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find idols of the user in database", err)
			return err
		}
		for _, id := range toUserIds {
			vStr := "relation" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
	}

	// <to_user_id, from_user_id> 表示谁被谁关注
	keyStr, valueStr = valueStr, keyStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		fromUserIds, err := dao.FindRelationFromUserIdsByToUserId(toUserId)
		if err != nil {
			log.Println("failed to find fans of the user in database", err)
			return err
		}
		for _, id := range fromUserIds {
			vStr := "relation" + util.RedisSplit + "fromUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
	}

	// <from_user_id, to_user_id> 表示某人的朋友
	res, err := dao.CheckRelationByBothId(toUserId, fromUserId)
	if err != nil {
		log.Println("failed to check relation in database", err)
		return err
	}
	if !res {
		return nil
	}
	keyStr = "relationFriend" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr
	valueStr = "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + toUserIdStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		friendIds, err := dao.FindRelationFriendIdsByFromUserId(fromUserId)
		if err != nil {
			log.Println("failed to find friends of the user in database", err)
			return err
		}
		for _, id := range friendIds {
			vStr := "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
	}

	keyStr = "relationFriend" + util.RedisSplit + "fromUserId" + util.RedisSplit + toUserIdStr
	valueStr = "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + fromUserIdStr
	n, err = redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println("failed to check key about relation in redis", err)
		return err
	}
	if n > 0 { //存在
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			return err
		}
	} else { //不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println("failed to add value about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println("failed to set expiration time about relation in redis", err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
			return err
		}
		friendIds, err := dao.FindRelationFriendIdsByFromUserId(toUserId)
		if err != nil {
			log.Println("failed to find friends of the user in database", err)
			return err
		}
		for _, id := range friendIds {
			vStr := "relationFriend" + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(id, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println("failed to add value about relation in redis", err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return err
			}
		}
		_, err = redis.RedisCli.SRem(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println("failed to remove value about relation in redis", err)
			return err
		}
	}
	return nil
}

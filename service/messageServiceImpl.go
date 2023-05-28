package service

import (
	"douyin/dao"
	"douyin/middleware/redis"
	"douyin/util"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MessageServiceImpl struct {
}

// AddMessageByAll 直接添加到数据库，返回message id，然后对redis进行操作
// redis维护数据结构 <fromUserId-toUserId, messageIds>
// 先查redis有没有该键，有则直接添加，重设过期时间
// 没有，查数据库，redis先设默认值，设置过期时间，一个一个添加，不需要添加新值
func (msi MessageServiceImpl) AddMessageByAll(fromUserId, toUserId int64, content string, createTime time.Time) (int64, error) {
	id, err := dao.AddMessageByAll(fromUserId, toUserId, content, createTime)
	if err != nil {
		log.Println(err)
		return util.NotExistId, err
	}
	fromUserIdStr := strconv.FormatInt(fromUserId, 10)
	toUserIdStr := strconv.FormatInt(toUserId, 10)
	idStr := strconv.FormatInt(id, 10)
	// "message-fromUserId-2-toUserId-3"
	// "message-id-4"
	keyStr := "message" + util.RedisSplit + "fromUserId" + util.RedisSplit + fromUserIdStr + util.RedisSplit + "toUserId" + util.RedisSplit + toUserIdStr
	valueStr := "message" + util.RedisSplit + "id" + util.RedisSplit + idStr
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println(err)
		return util.NotExistId, err
	}
	if n > 0 { // 存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, valueStr).Result()
		if err != nil {
			log.Println(err)
			return util.NotExistId, err
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println(err)
			return util.NotExistId, err
		}
	} else { // 不存在
		_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
		if err != nil {
			log.Println(err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
		}
		_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
		if err != nil {
			log.Println(err)
			redis.RedisCli.Del(redis.Ctx, keyStr)
		}
		ids, err := dao.FindMessageIdsByFromUserIdAndToUserId(fromUserId, toUserId)
		if err != nil {
			log.Println(err)
			return util.NotExistId, err
		}
		for _, idd := range ids {
			vStr := "message" + util.RedisSplit + "id" + util.RedisSplit + strconv.FormatInt(idd, 10)
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, vStr).Result()
			if err != nil {
				log.Println(err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return util.NotExistId, err
			}
		}
	}
	return id, nil
}

// FindMessageIdsByFromUserIdAndToUserId 先查redis，查不到查数据库，顺便更新redis
func (msi MessageServiceImpl) FindMessageIdsByFromUserIdAndToUserId(fromUserId, toUserId int64) ([]int64, error) {
	ids := make([]int64, 0)
	keyStr := "message" + util.RedisSplit + "fromUserId" + util.RedisSplit + strconv.FormatInt(fromUserId, 10) + util.RedisSplit + "toUserId" + util.RedisSplit + strconv.FormatInt(toUserId, 10)
	n, err := redis.RedisCli.Exists(redis.Ctx, keyStr).Result()
	if err != nil {
		log.Println(err)
		return ids, err
	}
	if n > 0 { // 存在
		vsStr, err := redis.RedisCli.SMembers(redis.Ctx, keyStr).Result()
		if err != nil {
			log.Println(err)
			return ids, err
		}
		for _, vStr := range vsStr {
			if vStr == util.RedisDefaultValue {
				continue
			}
			idStr := strings.Split(vStr, util.RedisSplit)[2]
			id, _ := strconv.ParseInt(idStr, 10, 64)
			ids = append(ids, id)
		}
	} else { // 不存在
		ids, err = dao.FindMessageIdsByFromUserIdAndToUserId(fromUserId, toUserId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		go func() {
			_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, util.RedisDefaultValue).Result()
			if err != nil {
				log.Println(err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return
			}
			_, err = redis.RedisCli.Expire(redis.Ctx, keyStr, util.RandomDuration()).Result()
			if err != nil {
				log.Println(err)
				redis.RedisCli.Del(redis.Ctx, keyStr)
				return
			}
			for _, id := range ids {
				tmpValueStr := "message" + util.RedisSplit + "id" + util.RedisSplit + strconv.FormatInt(id, 10)
				_, err = redis.RedisCli.SAdd(redis.Ctx, keyStr, tmpValueStr).Result()
				if err != nil {
					log.Println(err)
					redis.RedisCli.Del(redis.Ctx, keyStr)
					return
				}
			}
		}()
	}
	// 对ids排序
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] > ids[j]
	})
	return ids, nil
}

func (msi MessageServiceImpl) FindMessageById(id int64) (dao.Message, error) {
	message, err := dao.FindMessageById(id)
	if err != nil {
		log.Println(err)
		return message, err
	}
	return message, nil
}

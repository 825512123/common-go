package common_go

import (
	"context"
	"strconv"
	"time"

	//"github.com/garyburd/redigo/redis"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	REDIS *redis.Client
	ctx   = context.Background()
)

func InitRedis(add, pass string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     add,  //"localhost:6379",
		Password: pass, // no password set
		DB:       db,   // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		//Error级别的日志
		Logger().WithFields(logrus.Fields{
			"name": "hanyun",
		}).Error("redis connect ping failed, err:", "Error")
		panic(err)
	} else {
		Logger().WithFields(logrus.Fields{
			"name": "hanyun",
		}).Info("redis connect ping response:, err:"+pong, "Info")
		REDIS = client
	}
	return client
}

func RedisSet(key string, value interface{}) {
	err := REDIS.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func RedisSetOut(key string, value interface{}, t time.Duration) {
	err := REDIS.Set(ctx, key, value, time.Second*t).Err()
	if err != nil {
		panic(err)
	}
}

func RedisGet(key string) string {
	val, err := REDIS.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

func RedisGetInt(key string) int {
	val, err := strconv.Atoi(RedisGet(key))
	if err != nil {
		return 0
	}
	return val
}

// RedisIncrBy redis 给当前key自增val
func RedisIncrBy(key string, val int64) int64 {
	val, err := REDIS.IncrBy(ctx, key, val).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func RedisHset(key string, key1 string, value interface{}) {
	err := REDIS.HSet(ctx, key, key1, value).Err()
	if err != nil {
		panic(err)
	}
}
func RedisHDel(key string, key1 string) {
	err := REDIS.HDel(ctx, key, key1).Err()
	if err != nil {
		panic(err)
	}
}
func RedisHGetAll(key string) map[string]string {
	res, err := REDIS.HGetAll(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return res
}

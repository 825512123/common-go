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

// RedisSet 设置key缓存value
func RedisSet(key string, value interface{}) {
	err := REDIS.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

// RedisSetOut 设置key缓存value,过期时间t
func RedisSetOut(key string, value interface{}, t time.Duration) {
	err := REDIS.Set(ctx, key, value, time.Second*t).Err()
	if err != nil {
		panic(err)
	}
}

// RedisGet 获取key缓存value
func RedisGet(key string) string {
	val, err := REDIS.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

// RedisDel 删除key缓存
func RedisDel(key string) {
	REDIS.Del(ctx, key)
}

// RedisGetInt 获取key缓存值 int类型
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

// RedisHSet 设置hash表指定key的value
func RedisHSet(key string, key1 string, value interface{}) {
	err := REDIS.HSet(ctx, key, key1, value).Err()
	if err != nil {
		panic(err)
	}
}

// RedisHGet 获取hash表指定key1的value
func RedisHGet(key string, key1 string) string {
	val, err := REDIS.HGet(ctx, key, key1).Result()
	if err != nil {
		return ""
	}
	return val
}

// RedisHDel 删除hash表指定key1名
func RedisHDel(key string, key1 string) {
	REDIS.HDel(ctx, key, key1)
}

// RedisHMSet 设置hash表指定key-value键值对
func RedisHMSet(key string, value map[string]interface{}) bool {
	return REDIS.HMSet(ctx, key, value).Val()
}

// RedisHMGet 获取hash表指定keys的值
func RedisHMGet(key string, keys ...string) []interface{} {
	res, err := REDIS.HMGet(ctx, key, keys...).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisHGetAll 获取hash表所有值
func RedisHGetAll(key string) map[string]string {
	res, err := REDIS.HGetAll(ctx, key).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisHKeys 获取hash表所有key值
func RedisHKeys(key string) []string {
	res, err := REDIS.HKeys(ctx, key).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisHVals 获取hash表所有value值
func RedisHVals(key string) []string {
	res, err := REDIS.HVals(ctx, key).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisHLen 获取hash表长度
func RedisHLen(key string) int64 {
	res, err := REDIS.HLen(ctx, key).Result()
	if err != nil {
		return 0
	}
	return res
}

//RedisHExists 判断hash表中key是否存在
func RedisHExists(key string, key1 string) bool {
	return REDIS.HExists(ctx, key, key1).Val()
}

//RedisHSetNX 设置hash表指定key-value键值对，如果key存在则忽略
func RedisHSetNX(key string, key1 string, value interface{}) bool {
	return REDIS.HSetNX(ctx, key, key1, value).Val()
}

//RedisLPush 将一条数据添加到列表的头部（类似入栈）可直接添加切片[]string
func RedisLPush(key string, values ...interface{}) {
	REDIS.LPush(ctx, key, values)
}

//RedisRPush 将一条数据添加到列表的尾部 可直接添加切片[]string
func RedisRPush(key string, values ...interface{}) {
	REDIS.RPush(ctx, key, values)
}

//RedisLPop 弹出列表的头部 移除List的第一个元素（头元素）
func RedisLPop(key string) string {
	val, err := REDIS.LPop(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

//RedisRPop 弹出列表的尾部 移除List的最后一个元素（尾元素）
func RedisRPop(key string) string {
	val, err := REDIS.RPop(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

//RedisLLen 获取列表的长度
func RedisLLen(key string) int {
	val, err := REDIS.LLen(ctx, key).Result()
	if err != nil {
		return 0
	}
	return int(val)
}

//RedisLRange 获取列表指定范围内的元素
// 获取List中的元素：起始索引~结束索引，当结束索引 > llen(list)或=-1时，取出全部数据
// 遍历List，获取每一个元素
// 注意取出来的顺序！！！
func RedisLRange(key string, start, end int64) []string {
	val, err := REDIS.LRange(ctx, key, start, end).Result()
	if err != nil {
		return nil
	}
	return val
}

//RedisLRem 删除列表指定索引的元素
// 移除剩下的count个值：value(当移除的个数count大于该值的实际个数时，全部移除)
func RedisLRem(key string, count int64, value interface{}) {
	REDIS.LRem(ctx, key, count, value)
}

//RedisLTrim 保留列表指定索引范围内的数据
// 保留指定索引范围内的数据：0~-1，保留全部数据 截取的结束下标大于List长度或者-1时,一直截取到末尾
// 保留指定索引范围内的数据：0~10，保留索引0~10的数据
func RedisLTrim(key string, start, end int64) {
	REDIS.LTrim(ctx, key, start, end)
}

//RedisLIndex 获取列表指定索引的元素
func RedisLIndex(key string, index int64) string {
	val, err := REDIS.LIndex(ctx, key, index).Result()
	if err != nil {
		return ""
	}
	return val
}

//RedisLSet 设置指定索引的值
func RedisLSet(key string, index int64, value interface{}) {
	REDIS.LSet(ctx, key, index, value)
}

//RedisLInsertBefore 插入元素到列表中   当标志位不存在时,插入值失败
// 在元素前插入元素：value(当插入的个数count大于该值的实际个数时，全部插入)
func RedisLInsertBefore(key string, pivot, value interface{}) bool {
	err := REDIS.LInsertBefore(ctx, key, pivot, value).Err()
	if err != nil {
		return false
	}
	return true
}

//RedisLInsertAfter 插入元素到列表中  当标志位不存在时,插入值失败
// 在元素后插入元素：value(当插入的个数count大于该值的实际个数时，全部插入)
func RedisLInsertAfter(key string, pivot, value interface{}) bool {
	err := REDIS.LInsertAfter(ctx, key, pivot, value).Err()
	if err != nil {
		return false
	}
	return true
}

// RedisZAdd 增加有序集合
func RedisZAdd(key string, score float64, value string) bool {
	err := REDIS.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: value,
	}).Err()
	if err != nil {
		return false
	}
	return true
}

// RedisZRange 升序：查询zset中指定区间的成员,-1代表取到最后
func RedisZRange(key string, start, stop int64) []string {
	res, err := REDIS.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRevRange 降序：查询zset中指定区间的成员,-1代表取到最后
func RedisZRevRange(key string, start, stop int64) []string {
	res, err := REDIS.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRangeByScore 升序：查询 指定下标start, stop区间的成员 zset中指定分数min, max区间的成员,
func RedisZRangeByScore(key, min, max string, start, stop int64) []string {
	res, err := REDIS.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: start,
		Count:  stop,
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRevRangeByScore 降序：查询 指定下标start, stop区间的成员 zset中指定分数min, max区间的成员,
func RedisZRevRangeByScore(key, min, max string, start, stop int64) []string {
	res, err := REDIS.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: start,
		Count:  stop,
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRangeWithScores  升序：查询 指定下标start, stop区间的成员
func RedisZRangeWithScores(key string, start, stop int64) []redis.Z {
	res, err := REDIS.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil
	}
	return res
}

//  RedisZRevRangeWithScores  降序：查询 指定下标start, stop区间的成员
func RedisZRevRangeWithScores(key string, start, stop int64) []redis.Z {
	res, err := REDIS.ZRevRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRangeByScoreWithScores 升序：查询 指定下标start, stop区间的成员 zset中指定分数min, max区间的成员,
func RedisZRangeByScoreWithScores(key, min, max string, start, stop int64) []redis.Z {
	res, err := REDIS.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: start,
		Count:  stop,
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRevRangeByScoreWithScores 降序：查询 指定下标start, stop区间的成员 zset中指定分数min, max区间的成员,
func RedisZRevRangeByScoreWithScores(key, min, max string, start, stop int64) []redis.Z {
	res, err := REDIS.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: start,
		Count:  stop,
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRangeByLex
func RedisZRangeByLex(key, min, max string, start, stop int64) []string {
	res, err := REDIS.ZRangeByLex(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: start,
		Count:  stop,
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZRevRangeByLex 降序：查询 指定下标start, stop区间的成员 zset中指定分数min, max区间的成员,
func RedisZRevRangeByLex(key, min, max string, start, stop int64) []string {
	res, err := REDIS.ZRevRangeByLex(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: start,
		Count:  stop,
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

// RedisZScore 获取指定成员的score
func RedisZScore(key, member string) float64 {
	return REDIS.ZScore(ctx, key, member).Val()
}

// RedisZRank 获取指定成员的排名
func RedisZRank(key, member string) int64 {
	return REDIS.ZRank(ctx, key, member).Val()
}

// RedisZRevRank 获取指定成员的排名
func RedisZRevRank(key, member string) int64 {
	return REDIS.ZRevRank(ctx, key, member).Val()
}

// RedisZCount 返回指定区间中成员的个数
func RedisZCount(key, min, max string) int64 {
	return REDIS.ZCount(ctx, key, min, max).Val()
}

// RedisZCard 返回集合中成员的个数
func RedisZCard(key string) int64 {
	return REDIS.ZCard(ctx, key).Val()
}

// RedisZRem 删除指定成员
func RedisZRem(key, member string) int64 {
	return REDIS.ZRem(ctx, key, member).Val()
}

// RedisZRemRangeByRank 删除指定排名区间的成员
func RedisZRemRangeByRank(key string, start, stop int64) int64 {
	return REDIS.ZRemRangeByRank(ctx, key, start, stop).Val()
}

// RedisZRemRangeByScore 删除指定分数区间的成员
func RedisZRemRangeByScore(key, min, max string) int64 {
	return REDIS.ZRemRangeByScore(ctx, key, min, max).Val()
}

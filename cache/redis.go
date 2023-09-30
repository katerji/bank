package cache

import (
	"context"
	"fmt"
	"github.com/katerji/bank/envs"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisInstance *RedisClient

func GetRedisInstance() *RedisClient {
	if redisInstance == nil {
		redisInstance = &RedisClient{
			instance: newInstance(),
		}
	}
	return redisInstance
}

func newInstance() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", envs.GetInstance().GetDbHost(), envs.GetInstance().GetRedisPort()),
	})
}

type RedisClient struct {
	instance *redis.Client
}

func (r *RedisClient) GetBool(key string) (bool, error) {
	return r.instance.Get(context.Background(), key).Bool()
}

func (r *RedisClient) SetWithDefaultExpiry(key string, val any) {
	r.instance.Set(context.Background(), key, val, time.Hour*24)
}

func (r *RedisClient) Delete(key string) {
	r.instance.Del(context.Background(), key)
}

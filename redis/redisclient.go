package redis

import (
	"context"
	"github.com/339-Labs/exchange-market/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(config config.RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        config.Address,
		Username:    config.Username,
		Password:    config.Password,
		DialTimeout: 5 * time.Second,
		DB:          0,
	})
	return &RedisClient{rdb: rdb}, nil
}

func (r *RedisClient) CachePrice(symbol string, price string, ttl time.Duration) error {
	key := "price:" + symbol
	return r.rdb.Set(context.Background(), key, price, ttl).Err()
}

func (r *RedisClient) GetPrice(symbol string) (string, error) {
	key := "price:" + symbol
	return r.rdb.Get(context.Background(), key).Result()
}

func (r *RedisClient) Set(key string, value string, ttl time.Duration) error {
	return r.rdb.Set(context.Background(), key, value, ttl).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.rdb.Get(context.Background(), key).Result()
}

func (r *RedisClient) TryLock(lockKey string, value string, ttl time.Duration) (bool, error) {
	ok, err := r.rdb.SetNX(context.Background(), lockKey, value, ttl).Result()
	return ok, err
}

func (r *RedisClient) Unlock(lockKey string, value string) error {
	// 脚本防止误删别人的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	return r.rdb.Eval(context.Background(), script, []string{lockKey}, value).Err()
}

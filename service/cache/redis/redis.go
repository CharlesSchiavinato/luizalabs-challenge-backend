package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/util"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client     *redis.Client
	Expiration time.Duration
	OrderInst  cache.Order
}

func NewRedis(config *util.Config) (cache.Cache, error) {
	opt, err := redis.ParseURL(config.CacheURL)

	if err != nil {
		return nil, err
	}

	expiration, err := time.ParseDuration(config.CacheExpiration)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	err = client.Ping(context.Background()).Err()

	if err != nil {
		return nil, err
	}

	cacheRedis := &Redis{
		Client:     client,
		Expiration: expiration,
	}

	cacheRedis.OrderInst = NewRedisOrder(cacheRedis)

	return cacheRedis, nil
}

func (redis *Redis) Close() error {
	return redis.Client.Close()
}

func (redis *Redis) Check() error {
	return redis.Client.Ping(context.Background()).Err()
}

func (redis *Redis) Order() cache.Order {
	return redis.OrderInst
}

func RedisKeyFormat(identifier, fieldName, fieldValue string) string {
	return fmt.Sprintf("%v:%v:%v", identifier, fieldName, fieldValue)
}

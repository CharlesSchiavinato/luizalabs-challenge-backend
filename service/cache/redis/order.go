package cache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache"
)

type RedisOrder struct {
	Cache *Redis
}

func NewRedisOrder(cache *Redis) cache.Order {
	return &RedisOrder{Cache: cache}
}

func (redisOrder *RedisOrder) SetDetailsByOrderID(modelOrderDetails *model.OrderDetails) error {
	key := RedisKeyFormat("order", "id", strconv.FormatInt((*modelOrderDetails).Orders[0].OrderID, 10))
	value, err := json.Marshal(modelOrderDetails)

	if err != nil {
		return err
	}

	return redisOrder.Cache.Client.Set(context.Background(), key, value, redisOrder.Cache.Expiration).Err()
}

func (redisOrder *RedisOrder) GetDetailsByOrderID(orderID int64) (*model.OrderDetails, error) {
	key := RedisKeyFormat("order", "id", strconv.FormatInt(orderID, 10))
	value, err := redisOrder.Cache.Client.Get(context.Background(), key).Result()

	if err != nil {
		return nil, err
	}

	modelOrderDetails := &model.OrderDetails{}
	err = json.Unmarshal([]byte(value), modelOrderDetails)

	return modelOrderDetails, err
}

func (redisOrder *RedisOrder) DelDetailsByOrderID(orderID int64) error {
	key := RedisKeyFormat("order", "id", strconv.FormatInt(orderID, 10))
	return redisOrder.Cache.Client.Del(context.Background(), key).Err()
}

func (redisOrder *RedisOrder) ClearAll() error {
	return redisOrder.Cache.Client.FlushDB(context.Background()).Err()
}

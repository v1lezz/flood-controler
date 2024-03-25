package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"task/internal/config"
)

type RedisConn struct {
	client *redis.Client
}

func NewRedisConn(cfg config.RedisConfig) *RedisConn {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return &RedisConn{
		client: client,
	}
}

func (rc *RedisConn) IncrementIfLessK(ctx context.Context, userID int64, k int) (bool, error) {
	sResult, err := rc.client.Get(ctx, strconv.Itoa(int(userID))).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("%d:%d\n", userID, 0)
			err = rc.client.Set(ctx, strconv.Itoa(int(userID)), "1", 0).Err()
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	res, err := strconv.Atoi(sResult)
	if err != nil {
		return false, err
	}
	log.Printf("%d:%d\n", userID, res)
	if res < k {
		err = rc.client.Set(ctx, strconv.Itoa(int(userID)), strconv.Itoa(res+1), 0).Err()
		if err != nil {
			return true, err
		}
	}
	return res < k, nil
}

func (rc *RedisConn) Decrement(ctx context.Context, userID int64) error {
	sResult, err := rc.client.Get(ctx, strconv.Itoa(int(userID))).Result()
	if err != nil {
		return err
	}
	res, err := strconv.Atoi(sResult)
	if err != nil {
		return err
	}
	return rc.client.Set(ctx, strconv.Itoa(int(userID)), strconv.Itoa(res-1), 0).Err()
}

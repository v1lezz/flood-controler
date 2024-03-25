package redis_checker

import "C"
import (
	"context"
	"golang.org/x/sync/errgroup"
	"task/internal/config"
	"task/internal/ports/redis"
	"task/internal/repository"
	"time"
)

type RedisChecker struct {
	repo repository.Repository
	cfg  config.CheckerConfig
	eg   *errgroup.Group
}

func NewRedisChecker(cfg config.RedisConfig, cgCFG config.CheckerConfig, eg *errgroup.Group) *RedisChecker {
	return &RedisChecker{
		repo: redis.NewRedisConn(cfg),
		cfg:  cgCFG,
		eg:   eg,
	}
}

func (rc *RedisChecker) Check(ctx context.Context, userID int64) (bool, error) {
	flag, err := rc.repo.IncrementIfLessK(ctx, userID, rc.cfg.K)
	if err != nil {
		return false, err
	}
	if flag {
		rc.eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(time.Second * time.Duration(rc.cfg.N)):
				err = rc.repo.Decrement(ctx, userID)
				if err != nil {
					return err
				}
				return nil
			}
		})
	}
	return flag, nil
}

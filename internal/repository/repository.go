package repository

import "context"

type Repository interface {
	IncrementIfLessK(ctx context.Context, userID int64, k int) (bool, error)
	Decrement(ctx context.Context, userID int64) error
}

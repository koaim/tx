package tx

import (
	"context"
)

type Manager[T any] interface {
	InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error
}
package tx

import (
	"context"
)

type Manager[T any] interface {
	InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error
}

type FakeManager[T any] struct {
	repo T
}

func NewFakeManager[T any](repo T) *FakeManager[T] {
	return &FakeManager[T]{repo: repo}
}

func (m *FakeManager[T]) InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error {
	return f(ctx, m.repo)
}

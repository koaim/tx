package tx

import (
	"context"
)

// Manager runs callbacks with a repository of type T within a transaction.
type Manager[T any] interface {
	// InTx runs f with a transaction-bound repository.
	InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error
}

// FakeManager runs callbacks with a fixed repository without starting transactions.
// It should only be used in tests.
type FakeManager[T any] struct {
	repo T
}

// NewFakeManager creates a manager that passes repo to each callback.
func NewFakeManager[T any](repo T) *FakeManager[T] {
	return &FakeManager[T]{repo: repo}
}

// InTx calls f with the stored repository and returns its error.
func (m *FakeManager[T]) InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error {
	return f(ctx, m.repo)
}

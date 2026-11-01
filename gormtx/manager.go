package gormtx

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// Manager runs callbacks in GORM transactions with a repository of type T.
type Manager[T any] struct {
	conn     *gorm.DB
	wrapToTx func(*gorm.DB) T
	opts     *sql.TxOptions
}

// NewManager creates a manager with default transaction options.
// wrapToTx builds a repository from each transaction.
func NewManager[T any](conn *gorm.DB, wrapToTx func(*gorm.DB) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx}
}

// NewManagerWithOpts creates a manager with the given transaction options.
// wrapToTx builds a repository from each transaction.
func NewManagerWithOpts[T any](conn *gorm.DB, opts *sql.TxOptions, wrapToTx func(*gorm.DB) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx, opts: opts}
}

// InTx runs f in a new transaction, committing on success and rolling back otherwise.
func (m *Manager[T]) InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error {
	return m.conn.Transaction(func(tx *gorm.DB) error {
		return f(ctx, m.wrapToTx(tx))
	}, m.opts)
}

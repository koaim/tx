package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Manager runs callbacks in database/sql transactions with a repository of type T.
type Manager[T any] struct {
	conn     *sql.DB
	wrapToTx func(*sql.Tx) T
	opts     *sql.TxOptions
}

// NewManager creates a manager with default transaction options.
// wrapToTx builds a repository from each transaction.
func NewManager[T any](conn *sql.DB, wrapToTx func(*sql.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx}
}

// NewManagerWithOpts creates a manager with the given transaction options.
// wrapToTx builds a repository from each transaction.
func NewManagerWithOpts[T any](conn *sql.DB, opts *sql.TxOptions, wrapToTx func(*sql.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx, opts: opts}
}

// InTx runs f in a new transaction, committing on success and rolling back otherwise.
func (m *Manager[T]) InTx(ctx context.Context, f func(ctx context.Context, repo T) error) (err error) {
	tx, err := m.conn.BeginTx(ctx, m.opts)
	if err != nil {
		return err
	}

	defer func() {
		rlbckErr := tx.Rollback()
		if rlbckErr != nil && !errors.Is(rlbckErr, sql.ErrTxDone) {
			err = errors.Join(err, rlbckErr)
		}
	}()

	err = f(ctx, m.wrapToTx(tx))
	if err != nil {
		return fmt.Errorf("transaction callback error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

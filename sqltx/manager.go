package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Manager[T any] struct {
	conn     *sql.DB
	wrapToTx func(*sql.Tx) T
	opts     *sql.TxOptions
}

func NewManager[T any](conn *sql.DB, wrapToTx func(*sql.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx}
}

func NewManagerWithOpts[T any](conn *sql.DB, opts *sql.TxOptions, wrapToTx func(*sql.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx, opts: opts}
}

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

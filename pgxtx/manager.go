package pgxtx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type txBeginner interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type Manager[T any] struct {
	conn     txBeginner
	wrapToTx func(pgx.Tx) T
	opts     pgx.TxOptions
}

func NewManager[T any](conn txBeginner, wrapToTx func(pgx.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx}
}

func NewManagerWithOpts[T any](conn txBeginner, opts pgx.TxOptions, wrapToTx func(pgx.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx, opts: opts}
}

func (m *Manager[T]) InTx(ctx context.Context, f func(ctx context.Context, repo T) error) (err error) {
	tx, err := m.conn.BeginTx(ctx, m.opts)
	if err != nil {
		return err
	}

	defer func() {
		rlbckErr := tx.Rollback(ctx)
		if rlbckErr != nil && !errors.Is(err, pgx.ErrTxClosed) {
			err = errors.Join(err, rlbckErr)
		}
	}()

	err = f(ctx, m.wrapToTx(tx))
	if err != nil {
		return fmt.Errorf("transaction callback error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

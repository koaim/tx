package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Manager[T any] struct {
	conn     *sql.DB
	wrapToTx func(*sql.Tx) T
}

func NewManager[T any](conn *sql.DB, wrapToTx func(*sql.Tx) T) *Manager[T] {
	return &Manager[T]{conn: conn, wrapToTx: wrapToTx}
}

func (m *Manager[T]) InTx(ctx context.Context, f func(ctx context.Context, repo T) error) (err error) {
	tx, err := m.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		rlbckErr := tx.Rollback()
		if rlbckErr != nil {
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

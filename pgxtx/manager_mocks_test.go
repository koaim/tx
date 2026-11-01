package pgxtx_test

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type mockTx struct {
	pgx.Tx
	errRollback error
}

func (m *mockTx) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return m, nil
}

func (m *mockTx) Rollback(ctx context.Context) error {
	return m.errRollback
}

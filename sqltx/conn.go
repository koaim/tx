package sqltx

import (
	"context"
	"database/sql"
)

// Conn provides operations shared by sql.DB, sql.Conn, and sql.Tx.
type Conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

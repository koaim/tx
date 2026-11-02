package testdb

import (
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func PgxConnect(t *testing.T, dsn string) *pgx.Conn {
	t.Helper()

	db, err := pgx.Connect(t.Context(), dsn)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	return db
}

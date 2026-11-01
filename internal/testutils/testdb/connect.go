package testdb

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SqlConnect(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	db, err := sql.Open("pgx", "postgres://tx:tx@127.0.0.1:5432/tx_test?sslmode=disable")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	return db
}

func PgxConnect(t *testing.T, dsn string) *pgx.Conn {
	t.Helper()

	db, err := pgx.Connect(t.Context(), dsn)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	return db
}

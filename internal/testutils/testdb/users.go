package testdb

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	dsn = "postgres://tx:tx@127.0.0.1:5432/tx_test?sslmode=disable"
)

type DB struct {
	Conn *sql.DB
}

func New(t *testing.T) *DB {
	t.Helper()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	return &DB{Conn: db}
}

func (d *DB) WithPgxConn(t *testing.T, f func(*pgx.Conn) error) error {
	t.Helper()

	conn, err := d.Conn.Conn(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	return conn.Raw(func(driverConn any) error {
		return f(driverConn.(*stdlib.Conn).Conn())
	})
}

func (d *DB) UsersExists(t *testing.T, name string) bool {
	t.Helper()

	var exists bool
	err := d.Conn.QueryRowContext(t.Context(), "SELECT EXISTS (SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}

	return exists
}

func (d *DB) DeleteUser(t *testing.T, name string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 5*time.Second)
	defer cancel()

	_, err := d.Conn.ExecContext(ctx, "DELETE FROM users WHERE name = $1", name)
	if err != nil {
		t.Fatal(err)
	}
}

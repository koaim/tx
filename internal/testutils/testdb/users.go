package testdb

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func UserExists(t *testing.T, conn *sql.DB, name string) bool {
	t.Helper()

	var exists bool
	err := conn.QueryRowContext(t.Context(), "SELECT EXISTS (SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}

	return exists
}

func DeleteUser(t *testing.T, conn *sql.DB, name string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 5*time.Second)
	defer cancel()

	_, err := conn.ExecContext(ctx, "DELETE FROM users WHERE name = $1", name)
	if err != nil {
		t.Fatal(err)
	}
}

func PgxUserExists(t *testing.T, conn *pgx.Conn, name string) bool {
	t.Helper()

	var exists bool
	err := conn.QueryRow(t.Context(), "SELECT EXISTS (SELECT 1 FROM users WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}

	return exists
}

func PgxDeleteUser(t *testing.T, conn *pgx.Conn, name string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, "DELETE FROM users WHERE name = $1", name)
	if err != nil {
		t.Fatal(err)
	}
}

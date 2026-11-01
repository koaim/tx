package sqlxtx_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/koaim/tx/internal/testutils/testdb"
	"github.com/koaim/tx/sqlxtx"
)

const (
	dsn = "postgres://tx:tx@127.0.0.1:5432/tx_test?sslmode=disable"
)

func TestInTx_OK(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })
	t.Cleanup(func() { testdb.DeleteUser(t, db, name) })

	conn := sqlx.NewDb(db, "pgx")

	mngr := sqlxtx.NewManager(conn, func(tx *sqlx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return repo.Save(ctx, name)
	})

	if err != nil {
		t.Fatal(err)
	}
	if !testdb.UserExists(t, db, name) {
		t.Fatalf("user %q not exists in db", name)
	}
}

func TestInTx_ErrDb(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()
	tooLongName := uuid.NewV7().String() + uuid.NewV7().String()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })

	conn := sqlx.NewDb(db, "pgx")

	mngr := sqlxtx.NewManager(conn, func(tx *sqlx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		repo.Save(ctx, name)
		return repo.Save(ctx, tooLongName)
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
	if testdb.UserExists(t, db, tooLongName) {
		t.Fatalf("user %q should not have been saved to the database", tooLongName)
	}
	if testdb.UserExists(t, db, name) {
		t.Fatalf("user %q should not have been saved to the database", name)
	}
}

func TestInTx_ErrCallback(t *testing.T) {
	t.Parallel()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })

	conn := sqlx.NewDb(db, "pgx")

	mngr := sqlxtx.NewManager(conn, func(tx *sqlx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return errors.New("error")
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
}

func TestInTx_ErrReadonlyTx(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })

	conn := sqlx.NewDb(db, "pgx")

	mngr := sqlxtx.NewManagerWithOpts(conn, &sql.TxOptions{ReadOnly: true}, func(tx *sqlx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return repo.Save(ctx, name)
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
	if testdb.UserExists(t, db, name) {
		t.Fatalf("user %q should not have been saved to the database", name)
	}
}

func TestInTx_ErrRollback(t *testing.T) {
	t.Parallel()

	var (
		errRollback = errors.New("rollback err")
		errCallback = errors.New("callback err")
	)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errRollback)

	conn := sqlx.NewDb(db, "sqlmock")

	mngr := sqlxtx.NewManager(conn, func(tx *sqlx.Tx) repo { return newRepo(tx) })

	err = mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return errCallback
	})

	if !errors.Is(err, errCallback) {
		t.Errorf("\nwant: %q\nhave: %q", errCallback, err)
	}
	if !errors.Is(err, errRollback) {
		t.Errorf("\nwant: %q\nhave: %q", errRollback, err)
	}
}

type repo struct {
	conn sqlxtx.Conn
}

func newRepo(conn sqlxtx.Conn) repo {
	return repo{conn: conn}
}

func (r *repo) Save(ctx context.Context, name string) error {
	_, err := r.conn.ExecContext(ctx, "INSERT INTO users (name) VALUES($1)", name)
	if err != nil {
		return err
	}

	return nil
}

package sqltx_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/koaim/tx/internal/testutils/testdb"
	"github.com/koaim/tx/sqltx"
)

const (
	dsn = "postgres://tx:tx@127.0.0.1:5432/tx_test?sslmode=disable"
)

func TestInTx_OK(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })
	t.Cleanup(func() { testdb.DeleteRow(t, db, name) })

	mngr := sqltx.NewManager(db, func(tx *sql.Tx) testdb.Repo { return testdb.NewRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo testdb.Repo) error {
		return repo.Save(ctx, name)
	})

	if err != nil {
		t.Fatal(err)
	}
	if !testdb.RowExists(t, db, name) {
		t.Fatalf("name %q not exists in db", name)
	}
}

func TestInTx_ErrDb(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()
	tooLongName := uuid.NewV7().String() + uuid.NewV7().String()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })

	mngr := sqltx.NewManager(db, func(tx *sql.Tx) testdb.Repo { return testdb.NewRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo testdb.Repo) error {
		repo.Save(ctx, name)
		return repo.Save(ctx, tooLongName)
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
	if testdb.RowExists(t, db, tooLongName) {
		t.Fatalf("name %q should not have been saved to the database", tooLongName)
	}
	if testdb.RowExists(t, db, name) {
		t.Fatalf("name %q should not have been saved to the database", name)
	}
}

func TestInTx_ErrCallback(t *testing.T) {
	t.Parallel()

	db := testdb.SqlConnect(t, dsn)
	t.Cleanup(func() { db.Close() })

	mngr := sqltx.NewManager(db, func(tx *sql.Tx) testdb.Repo { return testdb.NewRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo testdb.Repo) error {
		return errors.New("error")
	})

	if err == nil {
		t.Fatal("expected not nil err")
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

	mngr := sqltx.NewManager(db, func(tx *sql.Tx) testdb.Repo { return testdb.NewRepo(tx) })

	err = mngr.InTx(t.Context(), func(ctx context.Context, repo testdb.Repo) error {
		return errCallback
	})

	if !errors.Is(err, errCallback) {
		t.Errorf("want %q, got %q", errCallback, err)
	}
	if !errors.Is(err, errRollback) {
		t.Errorf("want %q, got %q", errRollback, err)
	}
}

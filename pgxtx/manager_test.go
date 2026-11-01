package pgxtx_test

import (
	"context"
	"errors"
	"testing"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/koaim/tx/internal/testutils/testdb"
	"github.com/koaim/tx/pgxtx"
)

const (
	dsn = "postgres://tx:tx@127.0.0.1:5432/tx_test?sslmode=disable"
)

func TestInTx_OK(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()

	db := testdb.PgxConnect(t, dsn)
	t.Cleanup(func() { db.Close(t.Context()) })
	t.Cleanup(func() { testdb.PgxDeleteUser(t, db, name) })

	mngr := pgxtx.NewManager(db, func(tx pgx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return repo.Save(ctx, name)
	})

	if err != nil {
		t.Fatal(err)
	}
	if !testdb.PgxUserExists(t, db, name) {
		t.Fatalf("user %q not exists in db", name)
	}
}

func TestInTx_ErrDb(t *testing.T) {
	t.Parallel()

	name := uuid.NewV7().String()
	tooLongName := uuid.NewV7().String() + uuid.NewV7().String()

	db := testdb.PgxConnect(t, dsn)
	t.Cleanup(func() { db.Close(t.Context()) })

	mngr := pgxtx.NewManager(db, func(tx pgx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		repo.Save(ctx, name)
		return repo.Save(ctx, tooLongName)
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
	if testdb.PgxUserExists(t, db, tooLongName) {
		t.Fatalf("user %q should not have been saved to the database", tooLongName)
	}
	if testdb.PgxUserExists(t, db, name) {
		t.Fatalf("user %q should not have been saved to the database", name)
	}
}

func TestInTx_ErrCallback(t *testing.T) {
	t.Parallel()

	db := testdb.PgxConnect(t, dsn)
	t.Cleanup(func() { db.Close(t.Context()) })

	mngr := pgxtx.NewManager(db, func(tx pgx.Tx) repo { return newRepo(tx) })

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

	db := testdb.PgxConnect(t, dsn)
	t.Cleanup(func() { db.Close(t.Context()) })

	mngr := pgxtx.NewManagerWithOpts(db, pgx.TxOptions{AccessMode: pgx.ReadOnly}, func(tx pgx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return repo.Save(ctx, name)
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
	if testdb.PgxUserExists(t, db, name) {
		t.Fatalf("user %q should not have been saved to the database", name)
	}
}

func TestInTx_ErrRollback(t *testing.T) {
	t.Parallel()

	var (
		errRollback = errors.New("rollback err")
		errCallback = errors.New("callback err")
	)

	db := &mockTx{errRollback: errRollback}

	mngr := pgxtx.NewManager(db, func(tx pgx.Tx) repo { return newRepo(tx) })

	err := mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
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
	conn pgxtx.Conn
}

func newRepo(conn pgxtx.Conn) repo {
	return repo{conn: conn}
}

func (r *repo) Save(ctx context.Context, name string) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO users (name) VALUES($1)", name)
	if err != nil {
		return err
	}

	return nil
}

package gormtx_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"uuid"

	"github.com/koaim/tx/gormtx"
	"github.com/koaim/tx/internal/testutils/testdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}

	mngr := gormtx.NewManager(conn.WithContext(t.Context()), func(tx *gorm.DB) repo { return newRepo(tx) })

	err = mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
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

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}

	mngr := gormtx.NewManager(conn.WithContext(t.Context()), func(tx *gorm.DB) repo { return newRepo(tx) })

	err = mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
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

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}

	mngr := gormtx.NewManager(conn.WithContext(t.Context()), func(tx *gorm.DB) repo { return newRepo(tx) })

	err = mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
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

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}

	mngr := gormtx.NewManagerWithOpts(conn.WithContext(t.Context()), &sql.TxOptions{ReadOnly: true}, func(tx *gorm.DB) repo { return newRepo(tx) })

	err = mngr.InTx(t.Context(), func(ctx context.Context, repo repo) error {
		return repo.Save(ctx, name)
	})

	if err == nil {
		t.Fatal("expected not nil err")
	}
	if testdb.UserExists(t, db, name) {
		t.Fatalf("user %q should not have been saved to the database", name)
	}
}

type repo struct {
	conn *gorm.DB
}

func newRepo(conn *gorm.DB) repo {
	return repo{conn: conn}
}

func (r *repo) Save(ctx context.Context, name string) error {
	err := r.conn.WithContext(ctx).Exec("INSERT INTO users (name) VALUES(?)", name).Error
	if err != nil {
		return err
	}

	return nil
}

package gormtx_test

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"testing"

	"github.com/koaim/tx/gormtx"
	"github.com/koaim/tx/internal/testutils/testdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestInTx_OK(t *testing.T) {
	t.Parallel()

	name := rand.Text()

	db := testdb.New(t)
	t.Cleanup(func() { db.DeleteUser(t, name) })

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db.Conn}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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
	if !db.UsersExists(t, name) {
		t.Fatalf("user %q not exists in db", name)
	}
}

func TestInTx_ErrDb(t *testing.T) {
	t.Parallel()

	name := rand.Text()
	tooLongName := rand.Text() + rand.Text()

	db := testdb.New(t)

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db.Conn}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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
	if db.UsersExists(t, tooLongName) {
		t.Fatalf("user %q should not have been saved to the database", tooLongName)
	}
	if db.UsersExists(t, name) {
		t.Fatalf("user %q should not have been saved to the database", name)
	}
}

func TestInTx_ErrCallback(t *testing.T) {
	t.Parallel()

	db := testdb.New(t)

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db.Conn}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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

	name := rand.Text()

	db := testdb.New(t)

	conn, err := gorm.Open(postgres.New(postgres.Config{Conn: db.Conn}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
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
	if db.UsersExists(t, name) {
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

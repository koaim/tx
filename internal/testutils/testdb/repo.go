package testdb

import (
	"context"

	"github.com/koaim/tx/sqltx"
)

type Repo struct {
	conn sqltx.Conn
}

func NewRepo(conn sqltx.Conn) Repo {
	return Repo{conn: conn}
}

func (r *Repo) Save(ctx context.Context, name string) error {
	_, err := r.conn.ExecContext(ctx, "INSERT INTO users (name) VALUES($1)", name)
	if err != nil {
		return err
	}

	return nil
}

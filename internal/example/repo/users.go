package repo

import (
	"context"

	"github.com/koaim/tx/sqltx"
)

type Users struct {
	conn sqltx.Conn
}

func NewUsers(c sqltx.Conn) *Users {
	return &Users{conn: c}
}

func (u *Users) SaveUser(ctx context.Context, name string) error {
	return nil
}

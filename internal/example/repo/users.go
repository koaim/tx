package repo

import (
	"context"

	"github.com/koaim/tx/internal/example/model"
	"github.com/koaim/tx/sqltx"
)

type Users struct {
	conn sqltx.Conn
}

func NewUsers(c sqltx.Conn) *Users {
	return &Users{conn: c}
}

func (u *Users) Save(ctx context.Context, user model.User) error {
	return nil
}

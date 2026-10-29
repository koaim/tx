package repo

import (
	"context"

	"github.com/koaim/tx/sql/driver"
)

type Users struct {
	conn driver.Conn
}

func NewUsers(c driver.Conn) *Users {
	return &Users{conn: c}
}

func (u *Users) SaveUser(ctx context.Context, name string) error {
	return nil
}

package repo

import (
	"context"

	"github.com/koaim/tx/sqltx"
)

type Orders struct {
	conn sqltx.Conn
}

func NewOrders(c sqltx.Conn) *Orders {
	return &Orders{conn: c}
}

func (o *Orders) SaveOrder(ctx context.Context, name string) error {
	return nil
}

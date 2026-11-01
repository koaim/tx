package repo

import (
	"context"

	"github.com/koaim/tx/internal/example/model"
	"github.com/koaim/tx/sqltx"
)

type Orders struct {
	conn sqltx.Conn
}

func NewOrders(c sqltx.Conn) *Orders {
	return &Orders{conn: c}
}

func (o *Orders) Save(ctx context.Context, order model.Order) error {
	return nil
}

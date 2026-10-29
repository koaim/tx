package repo

import (
	"context"

	"github.com/koaim/tx/sql/driver"
)

type Orders struct {
	conn driver.Conn
}

func NewOrders(c driver.Conn) *Orders {
	return &Orders{conn: c}
}

func (o *Orders) SaveOrder(ctx context.Context, name string) error {
	return nil
}

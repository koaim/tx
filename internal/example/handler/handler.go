package handler

import (
	"context"

	"github.com/koaim/tx"
)

type users interface {
	SaveUser(ctx context.Context, name string) error
}

type orders interface {
	SaveOrder(ctx context.Context, name string) error
}

type Tx struct {
	users  users
	orders orders
}

func NewTx(u users, o orders) Tx {
	return Tx{users: u, orders: o}
}

type Handler struct {
	txmngr tx.Manager[Tx]
}

func NewHandler(mngr tx.Manager[Tx]) *Handler {
	return &Handler{txmngr: mngr}
}

func (h *Handler) Handle(ctx context.Context) error {
	return h.txmngr.InTx(ctx, func(ctx context.Context, repo Tx) error {
		repo.orders.SaveOrder(ctx, "name")
		repo.users.SaveUser(ctx, "name")

		return nil
	})
}

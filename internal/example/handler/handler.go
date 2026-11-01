package handler

import (
	"context"

	"github.com/koaim/tx"
	"github.com/koaim/tx/internal/example/model"
)

type users interface {
	Save(ctx context.Context, user model.User) error
}

type orders interface {
	Save(ctx context.Context, order model.Order) error
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
		repo.orders.Save(ctx, model.Order{})
		repo.users.Save(ctx, model.User{})

		return nil
	})
}

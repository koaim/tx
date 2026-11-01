# tx

[![Go version](https://img.shields.io/github/go-mod/go-version/koaim/tx)](https://github.com/koaim/tx/blob/main/go.mod)
[![GoDoc](https://godoc.org/github.com/koaim/tx?status.svg)](https://godoc.org/github.com/koaim/tx) [![license](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/koaim/tx/main/LICENSE)
[![codecov](https://codecov.io/gh/koaim/tx/branch/main/graph/badge.svg)](https://codecov.io/gh/koaim/tx)

Tiny transaction manager

## Installation

```bash
go get github.com/koaim/tx
go get github.com/koaim/tx/{sqltx,pgxtx,gormtx,sqlxtx}
```

## Features

- Does not store transactions in `context.Context`
- Works with `database/sql`, `sqlx`, `pgx`, and `gorm`
- Compatible with `sqlc` and `jet`

## Example

[Source code](internal/example)

```go
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
        // this code is running in transaction scope
		repo.orders.Save(ctx, model.Order{})
		repo.users.Save(ctx, model.User{})

		return nil
	})
}
```

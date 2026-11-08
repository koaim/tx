# tx

[![Go version](https://img.shields.io/github/go-mod/go-version/koaim/tx)](https://github.com/koaim/tx/blob/main/go.mod)
[![GoDoc](https://godoc.org/github.com/koaim/tx?status.svg)](https://godoc.org/github.com/koaim/tx) [![license](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/koaim/tx/main/LICENSE)
[![codecov](https://codecov.io/gh/koaim/tx/branch/main/graph/badge.svg)](https://codecov.io/gh/koaim/tx)

A tiny transaction manager

## Installation

Install the core library:

```bash
go get github.com/koaim/tx
```

Then install the adapter for your database library:

| Library | Installation command |
| --- | --- |
| `database/sql` | `go get github.com/koaim/tx/sqltx` |
| `sqlx` | `go get github.com/koaim/tx/sqlxtx` |
| `pgx` | `go get github.com/koaim/tx/pgxtx` |
| `gorm` | `go get github.com/koaim/tx/gormtx` |

## Features

- Does not store transactions in `context.Context`
- Works with [database/sql](https://pkg.go.dev/database/sql), [sqlx](https://github.com/jmoiron/sqlx/), [pgx](https://github.com/jackc/pgx) and [gorm](https://github.com/go-gorm/gorm)
- Compatible with [sqlc](https://github.com/sqlc-dev/sqlc), [jet](https://github.com/go-jet/jet) and [squirrel](https://github.com/masterminds/squirrel)
- Configurable transaction isolation levels

## Example

[Source code](internal/example)

Business logic layer:

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
        // this code runs within a transaction
        repo.orders.Save(ctx, model.Order{})
        repo.users.Save(ctx, model.User{})

        return nil
    })
}
```
`main.go`:
```go
package main

import (
	"context"
	"database/sql"

	"github.com/koaim/tx/internal/example/handler"
	"github.com/koaim/tx/internal/example/repo"
	"github.com/koaim/tx/sqltx"
)

func main() {
	db, _ := sql.Open("pgx", "dsn")
	h := handler.NewHandler(
		sqltx.NewManager(db, func(tx *sql.Tx) handler.Tx { return handler.NewTx(repo.NewUsers(tx), repo.NewOrders(tx)) }),
	)

	h.Handle(context.Background())
}
```

## Philosophy

Working with transactions in Go can be awkward. Three common approaches are:

- Storing the transaction in `context.Context`
- Working directly with transaction type (`sql.Tx`, `pgx.Tx` etc.) in the business logic layer
- Using a single god object for all entities and database operations

Each of these approaches has drawbacks that I am not willing to accept, so I built an alternative.

This project defines an interface for running transactions:

```go
type Manager[T any] interface {
    InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error
}
```

Adapters implement this interface for [database/sql](https://pkg.go.dev/database/sql), [sqlx](https://github.com/jmoiron/sqlx/), [pgx](https://github.com/jackc/pgx) and [gorm](https://github.com/go-gorm/gorm).

Each adapter (expect [gorm](https://github.com/go-gorm/gorm)) provides a `Conn` interface that defines the methods shared by its transaction and connection types. This interface should be injected into all repositories. As a result, the repository never knows whether it is working with a regular connection or a transaction.

## Trade-off

- Nested transactions are not supported: each call to `InTx` creates a new independent transaction
- The library requires repositories to use a connection pool rather than a single connection
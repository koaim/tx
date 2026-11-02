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

### Compatibility

All modules require Go 1.25.0 or later, including their test dependencies.

Tests live alongside the adapter implementations. Run `make test` to test all modules against the test database, or `make coverage` to generate a coverage report.

## Features

- Does not store transactions in `context.Context`
- Works with `database/sql`, `sqlx`, `pgx`, and `gorm`
- Compatible with `sqlc` and `jet`
- Configurable transaction isolation levels

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
        // this code runs within a transaction
        repo.orders.Save(ctx, model.Order{})
        repo.users.Save(ctx, model.User{})

        return nil
    })
}
```

## Philosophy

Working with transactions in Go can be awkward. Three common approaches are:

- Storing the transaction in `context.Context`
- Working directly with `*sql.Tx` in the business logic layer
- Using a single god object for all entities and database operations

Each of these approaches has drawbacks that I am not willing to accept, so I built an alternative.

This project defines an interface for running transactions:

```go
type Manager[T any] interface {
    InTx(ctx context.Context, f func(ctx context.Context, repo T) error) error
}
```

Adapters implement this interface for `database/sql`, `sqlx`, `pgx`, and `gorm`.
For libraries such as `database/sql` and `pgx`, repositories should depend on an abstraction such as `sqltx.Conn` or `pgxtx.Conn` rather than a concrete type such as `*sql.DB` or `*pgxpool.Pool`. These interfaces expose operations shared by connections and transactions, allowing the same repository implementation to work with either.

As a result, the repository never knows whether it is working with a regular connection or a transaction.

### Trade-off

- The library requires a connection pool (`*sql.DB` or `*pgxpool.Pool`) rather than a single connection (`*sql.Conn` or `*pgx.Conn`)
- Nested transactions are not supported

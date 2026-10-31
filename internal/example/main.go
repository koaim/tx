package example

import (
	"context"
	"database/sql"

	"github.com/koaim/tx/internal/example/handler"
	"github.com/koaim/tx/internal/example/repo"
	"github.com/koaim/tx/sqltx"
)

func main() {
	h := handler.NewHandler(
		sqltx.NewManager(nil, func(tx *sql.Tx) handler.Tx { return handler.NewTx(repo.NewUsers(tx), repo.NewOrders(tx)) }),
	)

	h.Handle(context.Background())
}

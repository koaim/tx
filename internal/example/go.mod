module github.com/koaim/tx/internal/example

go 1.27.0

require (
	github.com/koaim/tx v0.0.0
	github.com/koaim/tx/sqltx v0.0.0
)

replace github.com/koaim/tx => ../..

replace github.com/koaim/tx/sqltx => ../../sqltx

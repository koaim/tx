.PHONY: test-db-reset test tests t coverage

test-db-reset:
	docker compose rm --stop --force postgres
	docker compose up -d --wait postgres
	@set -e; for migration in internal/migrations/*.sql; do \
		[ -f "$$migration" ] || continue; \
		echo "Applying $$migration"; \
		docker compose exec -T postgres \
			psql -X --username=tx --dbname=tx_test --set=ON_ERROR_STOP=1 \
			< "$$migration"; \
	done

test tests t:
	go test ./... ./gormtx/... ./pgxtx/... ./sqltx/... ./sqlxtx/... ./internal/example/...

coverage:
	go test ./... ./gormtx/... ./pgxtx/... ./sqltx/... ./sqlxtx/... -coverpkg=github.com/koaim/tx,github.com/koaim/tx/gormtx,github.com/koaim/tx/pgxtx,github.com/koaim/tx/sqltx,github.com/koaim/tx/sqlxtx -coverprofile=coverage.out -covermode=atomic

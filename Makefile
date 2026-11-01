.PHONY: test-db-reset

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

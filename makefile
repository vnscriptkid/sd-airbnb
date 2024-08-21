up:
	docker compose up -d

down:
	docker compose down --remove-orphans --volumes

cli_redis:
	docker compose exec redis redis-cli

cli_pg:
	docker compose exec pg psql -U postgres
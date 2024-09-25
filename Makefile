include .env
export
# Migrations (Go)
.PHONY: migrate.up migrate.up.all migrate.down migrate.down.all migration migrate.force
migrate.up:
	migrate -path=$(MIGRATIONS_ROOT) -database $(DATABASE_URL) up $(n)
migrate.up.all:
	migrate -path=$(MIGRATIONS_ROOT) -database $(DATABASE_URL) up
migrate.down:
	migrate -path=$(MIGRATIONS_ROOT) -database $(DATABASE_URL) down $(n)
migrate.down.all:
	migrate -path=$(MIGRATIONS_ROOT) -database $(DATABASE_URL) down -all
migration:
	migrate create -seq -ext=.sql -dir=$(MIGRATIONS_ROOT) $(n)
migrate.force:
	migrate -path=$(MIGRATIONS_ROOT) -database=$(DATABASE_URL) force $(n)
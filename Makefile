APP_NAME=simplerag-go

.PHONY: up down run tidy migrate-local

up:
	docker compose up -d

down:
	docker compose down

run:
	go run ./cmd/api

tidy:
	go mod tidy

migrate-local:
	psql $$POSTGRES_DSN -f migrations/0001_init_extensions.sql
	psql $$POSTGRES_DSN -f migrations/0002_init_schema.sql
	psql $$POSTGRES_DSN -f migrations/0003_init_indexes.sql
	psql $$POSTGRES_DSN -f migrations/0004_add_document_content.sql

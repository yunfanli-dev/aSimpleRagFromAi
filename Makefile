APP_NAME=simplerag-go

.PHONY: up down run tidy

up:
	docker compose up -d

down:
	docker compose down

run:
	go run ./cmd/api

tidy:
	go mod tidy

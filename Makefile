.PHONY: help build-server build-client run-server migrate-up migrate-down docker-up docker-down test clean

DB_DSN ?= postgres://postgres:123456@localhost:5432/gophkeeper?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make build-server    - Build server binary"
	@echo "  make build-client    - Build client binary"
	@echo "  make run-server      - Run server"
	@echo "  make migrate-up      - Apply migrations"
	@echo "  make migrate-down    - Rollback migrations"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean binaries"

build-server:
	go build -o bin/gophkeeper-server cmd/server/main.go

build-client:
	go build -o bin/gophkeeper-client cmd/client/main.go

run-server:
	go run cmd/server/main.go

migrate-up:
	migrate -path ./migrations/migrations -database "$(DB_DSN)" up

migrate-down:
	migrate -path ./migrations/migrations -database "$(DB_DSN)" down


test:
	go test -v -cover ./...

clean:
	rm -rf bin/
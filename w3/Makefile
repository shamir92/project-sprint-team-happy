BIN_DIR := ./bin
MIGRATIONS_DIR=./migrations

PHONY: run-migrate-create
run-migrate-create:
	@echo 'Creating migrations file for ${name}'
	migrate create -ext sql -dir ${MIGRATIONS_DIR} ${name}

PHONY: run-migrate-up
run-migrate-up: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" up

PHONY: run-migrate-down-by-1
run-migrate-down-by-1: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" down 1

PHONY: run-migrate-down-all
run-migrate-down-all: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" down -all

PHONY: run-migrate-drop-force
run-migrate-drop-force: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" drop -f

PHONY: server/start
server/start:
	go run ./protocol/api/cmd/main.go
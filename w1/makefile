run-migrate-up: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" up

run-migrate-down-by-1: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" down 1

run-migrate-down-all: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" down -all

run-migrate-drop-force: 
	migrate -database "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" -path "./migrations" drop -f


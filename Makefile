default: run
.PHONY: run

PATH_DB=postgres://evans:evans@localhost:8001/postgres?sslmode=disable
FILE_MIGRATIONS = ./migrations

run:
	go run cmd/main.go -config ./config.env

run-default:
	go run cmd/main.go -config default

migrate:
	go run cmd/migrator/migrateup.go -storage-path $(PATH_DB) -migrations-path $(FILE_MIGRATIONS)
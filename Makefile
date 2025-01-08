default: run
.PHONY: run

PATH_DB = 
FILE_MIGRATIONS = 

run:
	go run cmd/main.go -config ./config.env

run-default:
	go run cmd/main.go -config default

migrate:
	go run cmd/migrator/migrationup.go -storage-path $(PATH_DB) -migrations-path $(FILE_MIGRATIONS)
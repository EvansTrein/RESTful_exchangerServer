package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

// main is the entry point of the migration script.
// It reads command-line flags for the database connection path and the path to the migration files.
// If either path is not provided, the program panics.
// It then initializes the migration tool and applies all pending migrations.
// If no migrations are needed, it logs a message and exits.
// If an error occurs during migration, the program panics.
func main() {
	var pathDB string
	var fileMigrationPath string

	flag.StringVar(&pathDB, "storage-path", "", "table creation path")
	flag.StringVar(&fileMigrationPath, "migrations-path", "", "path to migration file")
	flag.Parse()

	if pathDB == "" || fileMigrationPath == "" {
		panic("the path of the file with migrations or the path for database creation is not specified")
	}

	migrateDb, err := migrate.New("file://"+fileMigrationPath, pathDB)
	if err != nil {
		panic(err)
	}

	if err := migrateDb.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	log.Println("migrations have been successfully applied")
}

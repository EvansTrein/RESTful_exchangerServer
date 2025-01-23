package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

// PostgresDB represents a connection to a PostgreSQL database.
// It includes a database connection and a logger.
type PostgresDB struct {
	db  *sql.DB
	log *slog.Logger
}

// New creates a new PostgresDB instance and establishes a connection to the PostgreSQL server.
// It takes the database connection string and a logger as parameters.
// If the connection fails, it returns an error.
func New(storagePath string, log *slog.Logger) (*PostgresDB, error) {
	log.Debug("database: connection to Postgres started")

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("database: connect to Postgres successfully")
	return &PostgresDB{db: db, log: log}, nil
}

// Close closes the connection to the PostgreSQL server.
// If the connection is already closed, it returns an error.
func (s *PostgresDB) Close() error {
	s.log.Debug("database: stop started")

	if s.db == nil {
		return fmt.Errorf("database connection is already closed")
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	s.db = nil

	s.log.Info("database: stop successful")
	return nil
}

package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db  *sql.DB
	log *slog.Logger
}

// database connection
func New(storagePath string, log *slog.Logger) (*PostgresDB, error) {
	log.Debug("Connection to Postgres started")

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("connect to Postgres successfully")

	return &PostgresDB{db: db, log: log}, nil
}

// database disconnection
func (s *PostgresDB) Close() error {
	s.log.Debug("begins disconnecting from the database")

	if s.db == nil {
		return fmt.Errorf("database connection is already closed")
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	s.db = nil

	s.log.Info("database connection successfully closed")

	return nil
}

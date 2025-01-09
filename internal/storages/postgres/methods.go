package postgres

import (
	"fmt"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

func (s *PostgresDB) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {

	return &models.RegisterResponse{}, nil
}

func (s *PostgresDB) Login(req models.LoginRequest) (*models.LoginResponse, error) {

	return &models.LoginResponse{}, nil
}

func (s *PostgresDB) TestConnect() (int, error) {
	s.log.Debug("TestConnect DB")

	var result int
	err := s.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if result != 1 {
		return 0, fmt.Errorf("unexpected query result: %d", result)
	}

	return result, nil
}

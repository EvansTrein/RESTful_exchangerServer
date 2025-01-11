package postgres

import (
	"log/slog"
	"strings"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

func (db *PostgresDB) Register(req models.RegisterRequest) (uint, error) {
	op := "Database: user registration"
	log := db.log.With(slog.String("operation", op), slog.Any("requets data", req))
	log.Debug("Register func call")

	query := `WITH new_user AS (
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id)
	
	INSERT INTO accounts (user_id, currency_code)
	SELECT new_user.id, code
	FROM currencies, new_user
	RETURNING user_id;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return 0, err
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(req.Name, req.Email, req.HashPassword).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return 0, storages.ErrEmailAlreadyExists
		}
		log.Error("fail to execute SQL query", "error", err)
		return 0, err
	}

	db.log.Info("user registered successfully")
	return id, nil
}

func (db *PostgresDB) Login(req models.LoginRequest) (*models.LoginResponse, error) {

	return &models.LoginResponse{}, nil
}

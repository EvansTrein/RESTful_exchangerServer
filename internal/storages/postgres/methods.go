package postgres

import (
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

func (db *PostgresDB) CreateUser(req models.RegisterRequest) (uint, error) {
	op := "Database: user registration"
	log := db.log.With(slog.String("operation", op))
	log.Debug("Register func call", slog.Any("requets data", req))

	query := `WITH new_user AS (
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id)
	
	INSERT INTO accounts (user_id, currency_code)
	SELECT new_user.id, code
	FROM currencies
	CROSS JOIN new_user
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

	log.Info("data has been successfully saved in the database")
	return id, nil
}

func (db *PostgresDB) SearchUser(req models.LoginRequest) (*models.User, error) {
	op := "Database: user login"
	log := db.log.With(slog.String("operation", op))
	log.Debug("Login func call", slog.Any("requets data", req))

	query := `SELECT id, name, email, password_hash
		FROM users
		WHERE email = $1;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return nil, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.HashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("user with this email is not in the database", "email", req.Email)
			return nil, storages.ErrUserNotFound
		}
		log.Error("fail to execute SQL query", "error", err)
		return nil, err
	}

	log.Info("database successfully found the user")
	return &user, nil
}

func (db *PostgresDB) AllAccountsBalance(userId uint) (map[string]float32, error){
	op := "Database: balancing all accounts "
	log := db.log.With(slog.String("operation", op))
	log.Debug("AllAccountsBalance func call", slog.Any("requets data", userId))

	query := `SELECT currency_code, balance
		FROM accounts
		WHERE user_id = $1;`

	stmt, err := db.db.Prepare(query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
    if err != nil {
        log.Error("failed to execute SQL query", "error", err)
        return nil, err
    }
    defer rows.Close()

	accounts := make(map[string]float32)
    for rows.Next() {
        var currencyCode string
        var balance float32
        if err := rows.Scan(&currencyCode, &balance); err != nil {
            log.Error("failed to scan row", "error", err)
            return nil, err
        }
        accounts[currencyCode] = balance
    }

	if err := rows.Err(); err != nil {
        log.Error("error during rows iteration", "error", err)
        return nil, err
    }

	log.Info("database successfully returned data on all accounts")
	return accounts, nil
}

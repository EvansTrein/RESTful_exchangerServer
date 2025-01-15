package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

func (db *PostgresDB) CreateUser(ctx context.Context, req models.RegisterRequest) (uint, error) {
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

	stmt, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return 0, err
	}
	defer stmt.Close()

	var id uint
	err = stmt.QueryRowContext(ctx, req.Name, req.Email, req.HashPassword).Scan(&id)
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

func (db *PostgresDB) SearchUser(ctx context.Context, req models.LoginRequest) (*models.User, error) {
	op := "Database: user login"
	log := db.log.With(slog.String("operation", op))
	log.Debug("Login func call", slog.Any("requets data", req))

	query := `SELECT id, name, email, password_hash
		FROM users
		WHERE email = $1;`

	time.Sleep(time.Second * 7)
	stmt, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return nil, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRowContext(ctx, req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.HashPassword)
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

func (db *PostgresDB) AllAccountsBalance(ctx context.Context, userId uint) (map[string]float32, error) {
	op := "Database: balancing all accounts "
	log := db.log.With(slog.String("operation", op))
	log.Debug("AllAccountsBalance func call", slog.Any("requets data", userId))

	query := `SELECT currency_code, balance
		FROM accounts
		WHERE user_id = $1;`

	stmt, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("failed to prepare SQL query", "error", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId)
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

func (db *PostgresDB) ReplenishAccount(ctx context.Context, req models.DepositRequest) (map[string]float32, error) {
	op := "Database: account top-up "
	log := db.log.With(slog.String("operation", op))
	log.Debug("ReplenishAccount func call", slog.Any("requets data", req))

	currencyCheckQuery := `SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)`

	lockAccountsQuery := `
        SELECT EXISTS(
            SELECT 1 FROM accounts
            WHERE user_id = $1 AND currency_code = $2
            FOR UPDATE
        )`

	updateQuery := `
        UPDATE accounts
        SET balance = balance + $1
        WHERE user_id = $2 AND currency_code = $3`

	selectNewBalanceQuery := `
        SELECT currency_code, balance
        FROM accounts
        WHERE user_id = $1`

	currencyCheckStmt, err := db.db.PrepareContext(ctx, currencyCheckQuery)
	if err != nil {
		log.Error("failed to prepare currency check SQL query", "error", err)
		return nil, err
	}
	defer currencyCheckStmt.Close()

	lockAccountsStmt, err := db.db.PrepareContext(ctx, lockAccountsQuery)
	if err != nil {
		log.Error("failed to prepare lock accounts SQL query", "error", err)
		return nil, err
	}
	defer lockAccountsStmt.Close()

	updateStmt, err := db.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		log.Error("failed to prepare update SQL query", "error", err)
		return nil, err
	}
	defer updateStmt.Close()

	selectNewBalanceStmt, err := db.db.PrepareContext(ctx, selectNewBalanceQuery)
	if err != nil {
		log.Error("failed to prepare select balance SQL query", "error", err)
		return nil, err
	}
	defer selectNewBalanceStmt.Close()

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", "error", err)
		return nil, err
	}

	var currencyExists bool
	if err = tx.StmtContext(ctx, currencyCheckStmt).QueryRow(req.Currency).Scan(&currencyExists); err != nil {
		tx.Rollback()
		log.Error("failed to check currency", "error", err)
		return nil, err
	}

	if !currencyExists {
		tx.Rollback()
		log.Warn("currency not found", "currency", req.Currency)
		return nil, storages.ErrCurrencyNotFound
	}

	var accountExists bool
	if err = tx.StmtContext(ctx, lockAccountsStmt).QueryRow(req.UserID, req.Currency).Scan(&accountExists); err != nil {
		tx.Rollback()
		log.Error("failed to lock accounts", "error", err)
		return nil, storages.ErrAccountNotFound
	}

	if !accountExists {
		tx.Rollback()
		log.Error("account not found", "user_id", req.UserID, "currency", req.Currency)
		return nil, err
	}

	if _, err = tx.StmtContext(ctx, updateStmt).Exec(req.Amount, req.UserID, req.Currency); err != nil {
		tx.Rollback()
		log.Error("failed to update account balance", "error", err)
		return nil, err
	}

	rows, err := tx.StmtContext(ctx, selectNewBalanceStmt).Query(req.UserID)
	if err != nil {
		tx.Rollback()
		log.Error("failed to query account balances", "error", err)
		return nil, err
	}
	defer rows.Close()

	accounts := make(map[string]float32)
	for rows.Next() {
		var currencyCode string
		var balance float32
		if err := rows.Scan(&currencyCode, &balance); err != nil {
			tx.Rollback()
			log.Error("failed to scan row", "error", err)
			return nil, err
		}
		accounts[currencyCode] = balance
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		log.Error("error occurred during row iteration", "error", err)
		return nil, err
	}

	if len(accounts) == 0 {
		tx.Rollback()
		log.Error("no rows were returned by the query")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Error("failed to commit transaction", "error", err)
		return nil, err
	}

	return accounts, nil
}

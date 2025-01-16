package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
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

func (db *PostgresDB) AccountOperation(ctx context.Context, req *models.AccountOperationRequest) (map[string]float32, error) {
	op := "Database: account change"
	log := db.log.With(slog.String("operation", op))
	log.Debug("AccountOperation func call", slog.Any("requets data", req))

	if req.Operation == "" {
		log.Error("no database operation specified", "error", storages.ErrUnspecifiedOperation)
		return nil, storages.ErrUnspecifiedOperation
	}

	log.Debug("еxecuting an account transaction", "account operation", req.Operation)

	currencyCheckQuery := `SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)`

	getBalanceAndLockQuery := `
        SELECT balance
        FROM accounts
        WHERE user_id = $1 AND currency_code = $2
        FOR UPDATE`

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

	getBalanceAndLockStmt, err := db.db.PrepareContext(ctx, getBalanceAndLockQuery)
	if err != nil {
		log.Error("failed to prepare lock accounts SQL query", "error", err)
		return nil, err
	}
	defer getBalanceAndLockStmt.Close()

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

	log.Debug("all SQL queries for the transaction have been prepared successfully")

	// Start transaction
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", "error", err)
		return nil, err
	}

	var currencyExists bool
	if err = tx.StmtContext(ctx, currencyCheckStmt).QueryRow(req.Currency).Scan(&currencyExists); err != nil {
		tx.Rollback()
		log.Error("failed to check currency", "error", err, "transaction", "rollback")
		return nil, err
	}

	if !currencyExists {
		tx.Rollback()
		log.Warn("currency not found", "currency", req.Currency, "transaction", "rollback")
		return nil, storages.ErrCurrencyNotFound
	}

	var currentBalance float32
	if err = tx.StmtContext(ctx, getBalanceAndLockStmt).QueryRow(req.UserID, req.Currency).Scan(&currentBalance); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("account not found", "user id", req.UserID, "currency", req.Currency, "transaction", "rollback")
			return nil, storages.ErrAccountNotFound
		}
		log.Error("failed to get current balance", "error", err)
		return nil, err
	}

	if req.Operation == services.OperationWithdraw && currentBalance < req.Amount {
		tx.Rollback()
		log.Warn("insufficient funds", "current balance", currentBalance, "requested amount", req.Amount, "transaction", "rollback")
		return nil, storages.ErrInsufficientFunds
	}

	var amount float32
	switch req.Operation {
	case services.OperationDeposit:
		amount = req.Amount
	case services.OperationWithdraw:
		amount = -req.Amount
	default:
		tx.Rollback()
		log.Error("invalid operation type", "operation type", req.Operation)
		return nil, storages.ErrInvalidOperationType
	}

	log.Debug("all business logic checks have been completed successfully")

	if _, err = tx.StmtContext(ctx, updateStmt).Exec(amount, req.UserID, req.Currency); err != nil {
		tx.Rollback()
		log.Error("failed to update account balance", "error", err, "transaction", "rollback")
		return nil, err
	}

	rows, err := tx.StmtContext(ctx, selectNewBalanceStmt).Query(req.UserID)
	if err != nil {
		tx.Rollback()
		log.Error("failed to query account balances", "error", err, "transaction", "rollback")
		return nil, err
	}
	defer rows.Close()

	accounts := make(map[string]float32)
	for rows.Next() {
		var currencyCode string
		var balance float32
		if err := rows.Scan(&currencyCode, &balance); err != nil {
			tx.Rollback()
			log.Error("failed to scan row", "error", err, "transaction", "rollback")
			return nil, err
		}
		accounts[currencyCode] = balance
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		log.Error("error occurred during row iteration", "error", err, "transaction", "rollback")
		return nil, err
	}

	if len(accounts) == 0 {
		tx.Rollback()
		log.Error("no rows were returned by the query", "transaction", "rollback")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Error("!!!ATTENTION!!! failed to commit transaction", "error", err)
		return nil, err
	}

	log.Info("transaction successfully completed")
	return accounts, nil
}

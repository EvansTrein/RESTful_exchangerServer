package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

const transactionTimeout = time.Second * 15

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

func (db *PostgresDB) AllAccountsBalance(userId uint) (map[string]float32, error) {
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

func (db *PostgresDB) ReplenishAccount(req models.DepositRequest) (map[string]float32, error) {
	op := "Database: account top-up "
	log := db.log.With(slog.String("operation", op))
	log.Debug("ReplenishAccount func call", slog.Any("requets data", req))

	ctx, cancelTimeout := context.WithTimeout(context.Background(), transactionTimeout)
	defer cancelTimeout()

	resultChan := make(chan map[string]float32, 1)
	errorChan := make(chan error, 1)

	currencyCheckQuery := `SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)`

	// lockAccountsQuery := `
    //     SELECT EXISTS(
    //         SELECT 1 FROM accounts
    //         WHERE user_id = $1 AND currency_code = $2
    //         FOR UPDATE
    //     )`

	updateQuery := `
        UPDATE accounts
        SET balance = balance + $1
        WHERE user_id = $2 AND currency_code = $3`

	selectNewBalanceQuery := `
        SELECT currency_code, balance
        FROM accounts
        WHERE user_id = $1`

	currencyCheckStmt, err := db.db.Prepare(currencyCheckQuery)
	if err != nil {
		log.Error("failed to prepare currency check SQL query", "error", err)
		return nil, err
	}
	defer currencyCheckStmt.Close()

	updateStmt, err := db.db.Prepare(updateQuery)
	if err != nil {
		log.Error("failed to prepare update SQL query", "error", err)
		return nil, err
	}
	defer updateStmt.Close()

	selectNewBalanceStmt, err := db.db.Prepare(selectNewBalanceQuery)
	if err != nil {
		log.Error("failed to prepare select balance SQL query", "error", err)
		return nil, err
	}
	defer selectNewBalanceStmt.Close()

	// run the transaction in a separate goroutine
	go func() {
		var goroutineName = "deposit transaction"

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			log.Error("failed to begin transaction", "error", err, "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("failed to begin transaction: %w", err)
			return
		}

		var exists bool
		if err = tx.StmtContext(ctx, currencyCheckStmt).QueryRow(req.Currency).Scan(&exists); err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				log.Warn("currency not found", "currency", req.Currency, "GOROUTINE", goroutineName)
				errorChan <- storages.ErrCurrencyNotFound
				return
			}
			log.Error("failed to check currency", "error", err, "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("failed to check currency: %w", err)
			return
		}

		if _, err = tx.StmtContext(ctx, updateStmt).Exec(req.Amount, req.UserID, req.Currency); err != nil {
			tx.Rollback()
			log.Error("failed to update account balance", "error", err, "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("failed to update account balance: %w", err)
			return
		}

		// was specifically added to simulate the long runtime between database queries
		time.Sleep(time.Second * 10)

		rows, err := tx.StmtContext(ctx, selectNewBalanceStmt).Query(req.UserID)
		if err != nil {
			tx.Rollback()
			log.Error("failed to query account balances", "error", err, "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("failed to query account balances: %w", err)
			return
		}
		defer rows.Close()

		accounts := make(map[string]float32)
		for rows.Next() {
			var currencyCode string
			var balance float32
			if err := rows.Scan(&currencyCode, &balance); err != nil {
				tx.Rollback()
				log.Error("failed to scan row", "error", err, "GOROUTINE", goroutineName)
				errorChan <- fmt.Errorf("failed to scan row: %w", err)
				return
			}
			accounts[currencyCode] = balance
		}

		if err := rows.Err(); err != nil {
			tx.Rollback()
			log.Error("error occurred during row iteration", "error", err, "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("error occurred during row iteration: %w", err)
			return
		}

		if len(accounts) == 0 {
			tx.Rollback()
			log.Error("no rows were returned by the query", "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("no rows were returned by the query")
			return
		}

		if err = tx.Commit(); err != nil {
			log.Error("failed to commit transaction", "error", err, "GOROUTINE", goroutineName)
			errorChan <- fmt.Errorf("failed to commit transaction: %w", err)
			return
		}

		resultChan <- accounts
	}()

	// timeout tracking in a separate goroutine, added to check simulation of a long operation in a transaction
	go func() {
		select {
		case <-ctx.Done():
			// context expired, send error
			log.Error("context timeout expired")
			// equivalently errors.Is(err, context.DeadlineExceeded)
			errorChan <- fmt.Errorf("operation timed out: %w", ctx.Err())
		case <-time.After(transactionTimeout + (time.Second * 2)):
			errorChan <- errors.New("transaction took too long")
		}
	}()

	// purely for training purposes, above we have run the transaction in a separate
	// and in another goroutine, timeout tracking is running, this was done to check for long latency between database queries
	// i.e. we start 2 goroutines and wait for what will happen first: timeout expires, transaction error or result comes
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

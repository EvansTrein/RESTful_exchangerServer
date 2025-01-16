package storages

import (
	"context"
	"errors"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

var (
	// errors StoreAuth
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidLoginData   = errors.New("invalid email or password")

	// errors StoreWallet
	ErrCurrencyNotFound = errors.New("currency not found")
	ErrAccountNotFound  = errors.New("account not found")
)

type StoreAuth interface {
	CreateUser(ctx context.Context, req models.RegisterRequest) (uint, error)
	SearchUser(ctx context.Context, req models.LoginRequest) (*models.User, error)
}

type StoreWallet interface {
	AllAccountsBalance(ctx context.Context, userId uint) (map[string]float32, error)
	Deposit(ctx context.Context, req models.DepositRequest) (map[string]float32, error)
	Withdraw(ctx context.Context, req models.WithdrawRequest) (map[string]float32, error)
}

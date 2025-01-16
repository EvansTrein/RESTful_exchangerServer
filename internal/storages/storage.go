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
	ErrCurrencyNotFound     = errors.New("currency not found")
	ErrAccountNotFound      = errors.New("account not found")
	ErrUnspecifiedOperation = errors.New("unspecified operation")
	ErrInsufficientFunds    = errors.New("insufficient account balance")
	ErrInvalidOperationType = errors.New("invalid operation type")
)

type StoreAuth interface {
	CreateUser(ctx context.Context, req models.RegisterRequest) (uint, error)
	SearchUser(ctx context.Context, req models.LoginRequest) (*models.User, error)
}

type StoreWallet interface {
	AllAccountsBalance(ctx context.Context, userId uint) (map[string]float32, error)
	AccountOperation(ctx context.Context, req *models.AccountOperationRequest) (map[string]float32, error)
}

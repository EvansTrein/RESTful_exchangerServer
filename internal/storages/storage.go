package storages

import (
	"errors"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

var (
	// errors StoreAuth
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidLoginData = errors.New("invalid email or password")

	// errors StoreWallet
	ErrCurrencyNotFound = errors.New("currency not found")
)

type StoreAuth interface {
	CreateUser(req models.RegisterRequest) (uint, error)
	SearchUser(req models.LoginRequest) (*models.User, error)
}

type StoreWallet interface {
	AllAccountsBalance(userId uint) (map[string]float32, error)
	ReplenishAccount(req models.DepositRequest) (map[string]float32, error)
}
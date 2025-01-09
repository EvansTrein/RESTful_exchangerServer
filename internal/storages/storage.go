package storages

import (
	"errors"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

var (
	// errors StoreAuth
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidLoginData = errors.New("invalid username or password")

	// errors StoreWallet

)

type StoreAuth interface {
	Register(req models.RegisterRequest) (*models.RegisterResponse, error)
	Login(req models.LoginRequest) (*models.LoginResponse, error)
}

type StoreWallet interface {
	TestConnect() (int, error)
}
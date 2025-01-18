package storages

import (
	"context"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

type StoreAuth interface {
	CreateUser(ctx context.Context, req models.RegisterRequest) (uint, error)
	SearchUser(ctx context.Context, req models.LoginRequest) (*models.User, error)
	DeleteUser(ctx context.Context, userId uint) error
}

type StoreWallet interface {
	AllAccountsBalance(ctx context.Context, userId uint) (map[string]float32, error)
	AccountOperation(ctx context.Context, req *models.AccountOperationRequest) (map[string]float32, error)
	SaveExchangeRateChanges(ctx context.Context, newData *models.CurrencyExchangeResult) error
}

type CacheDB interface {
	SetExchange(fromCurrency, toCurrency string, value float32) error
	GetExchange(fromCurrency, toCurrency string) (float32, error)
}
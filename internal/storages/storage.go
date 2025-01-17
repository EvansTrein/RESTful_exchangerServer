package storages

import (
	"context"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

type StoreAuth interface {
	CreateUser(ctx context.Context, req models.RegisterRequest) (uint, error)
	SearchUser(ctx context.Context, req models.LoginRequest) (*models.User, error)
}

type StoreWallet interface {
	AllAccountsBalance(ctx context.Context, userId uint) (map[string]float32, error)
	AccountOperation(ctx context.Context, req *models.AccountOperationRequest) (map[string]float32, error)
	SaveExchangeRateChanges(ctx context.Context, newData *models.CurrencyExchangeResult) error
}

type CacheDB interface {
	TestMethodSet(key, value string) error
	TestMethodGet(key string) (string, error)
}
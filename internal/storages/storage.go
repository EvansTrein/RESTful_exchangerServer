package storages

import (
	"context"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
)

// StoreAuth defines the interface for authentication-related database operations.
// It includes methods for creating, searching, and deleting users.
type StoreAuth interface {
	CreateUser(ctx context.Context, req models.RegisterRequest) (uint, error)
	SearchUser(ctx context.Context, req models.LoginRequest) (*models.User, error)
	DeleteUser(ctx context.Context, userId uint) error
}

// StoreWallet defines the interface for wallet-related database operations.
// It includes methods for retrieving account balances, performing account operations, and saving exchange rate changes.
type StoreWallet interface {
	AllAccountsBalance(ctx context.Context, userId uint) (map[string]float32, error)
	AccountOperation(ctx context.Context, req *models.AccountOperationRequest) (map[string]float32, error)
	SaveExchangeRateChanges(ctx context.Context, newData *models.CurrencyExchangeResult) error
}

// CacheDB defines the interface for cache-related operations.
// It includes methods for setting and retrieving exchange rates.
type CacheDB interface {
	SetExchange(fromCurrency, toCurrency string, value float32) error
	GetExchange(fromCurrency, toCurrency string) (float32, error)
}
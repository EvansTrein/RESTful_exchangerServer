package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
)

const (
	OperationDeposit  = "deposit"
	OperationWithdraw = "withdraw"
)

var (
	ErrCurrencyNotFound     = errors.New("currency not found")
	ErrAccountNotFound      = errors.New("account not found")
	ErrUnspecifiedOperation = errors.New("unspecified operation")
	ErrInsufficientFunds    = errors.New("insufficient account balance")
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrNegativeBalance      = errors.New("negative balance")
)

type Wallet struct {
	log        *slog.Logger
	clientGRPC grpcclient.ClientGRPC
	db         storages.StoreWallet
	cacheDB    storages.CacheDB
}

func New(log *slog.Logger, gRPC grpcclient.ClientGRPC, db storages.StoreWallet, cacheDB storages.CacheDB) *Wallet {
	log.Debug("service Wallet: started creating")

	log.Info("service Wallet: successfully created")
	return &Wallet{
		log:        log,
		clientGRPC: gRPC,
		db:         db,
		cacheDB:    cacheDB,
	}
}

func (w *Wallet) Stop() error {
	w.log.Debug("service Wallet: stop started")

	w.clientGRPC = nil
	w.db = nil
	w.cacheDB = nil

	w.log.Info("service Wallet: stop successful")
	return nil
}

func (w *Wallet) Balance(ctx context.Context, req models.BalanceRequest) (*models.BalanceResponse, error) {
	op := "service Wallet: getting the balance of all accounts"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Balance func call", slog.Any("requets data", req))

	accounts, err := w.db.AllAccountsBalance(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get the balance of all accounts from the database", "error", err)
		return nil, err
	}

	log.Debug("balance data for all accounts successfully obtained from the database", "accounts", accounts)

	var resp models.BalanceResponse
	resp.Balance = accounts

	log.Info("balance data for all accounts successfully sent")
	return &resp, nil
}

func (w *Wallet) Deposit(ctx context.Context, req *models.AccountOperationRequest) (*models.AccountOperationResponse, error) {
	op := "service Wallet: deposit request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Deposit func call", slog.Any("requets data", req))

	req.Operation = OperationDeposit

	newBalance, err := w.db.AccountOperation(ctx, req)
	if err != nil {
		log.Error("failed to deposit in the database", "error", err)
		return nil, err
	}

	log.Debug("account operation successful", "new balance", newBalance)

	var resp models.AccountOperationResponse
	resp.Message = "successfully deposit"
	resp.NewBalance = newBalance

	log.Info("successfully deposit")
	return &resp, nil
}

func (w *Wallet) Withdraw(ctx context.Context, req *models.AccountOperationRequest) (*models.AccountOperationResponse, error) {
	op := "service Wallet: withdraw request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Withdraw func call", slog.Any("requets data", req))

	req.Operation = OperationWithdraw

	newBalance, err := w.db.AccountOperation(ctx, req)
	if err != nil {
		log.Error("failed to withdraw in the database", "error", err)
		return nil, err
	}

	log.Debug("account operation successful", "new balance", newBalance)

	var resp models.AccountOperationResponse
	resp.Message = "successfully withdrawn"
	resp.NewBalance = newBalance

	log.Info("successfully withdrawn")
	return &resp, nil
}

func (w *Wallet) Exchange(ctx context.Context, req models.ExchangeRequest) (*models.ExchangeResponse, error) {
	op := "service Wallet: currency exchange request"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Exchange func call", slog.Any("requets data", req))

	// TODO: запустить в горутине работу с курсом валют, получением его из Redis или БД и сохранением в Redis

	balanceUser, err := w.db.AllAccountsBalance(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get current user balance", "error", err)
		return nil, err
	}

	currentBaseAccountBalance, ok := balanceUser[req.FromCurrency]
	if !ok {
		log.Warn("no base currency to exchange", "currency", req.FromCurrency)
		return nil, ErrAccountNotFound
	}

	currentToAccountBalance, ok := balanceUser[req.ToCurrency]
	if !ok {
		log.Warn("not to currency exchange", "currency", req.ToCurrency)
		return nil, ErrCurrencyNotFound
	}

	if req.Amount > currentBaseAccountBalance {
		log.Warn("insufficient funds for exchange", "current balance", currentBaseAccountBalance, "requested amount", req.Amount)
		return nil, ErrInsufficientFunds
	}

	log.Debug("business logic check successfully completed")

	// TODO: проверить есть ли курс в Redis

	var rate models.ExchangeGRPC
	rate.FromCurrency = req.FromCurrency
	rate.ToCurrency = req.ToCurrency

	if err := w.clientGRPC.ExchangeRate(ctx, &rate); err != nil {
		log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	log.Debug("exchange rate successfully received from gRPC server", "rate from gRPC", rate)
	// TODO: тут получен курс, сохранить его в Redis

	// collect data to calculate new account balances
	exchangeData := models.CurrencyExchangeData{
		BaseBalance:  currentBaseAccountBalance,
		ToBalance:    currentToAccountBalance,
		ExchangeRate: rate.Rate,
		Amount:       req.Amount,
	}
	exchangeResult, err := w.CurrencyExchangeLogic(&exchangeData)
	if err != nil {
		log.Error("currency exchange failed", "error", err)
		return nil, err
	}
	// TODO: обмен произведен, сохранить изменения в БД
	exchangeResult.UserID = req.UserID
	exchangeResult.BaseCurrency = req.FromCurrency
	exchangeResult.ToCurrency = req.ToCurrency
	if err := w.db.SaveExchangeRateChanges(ctx, exchangeResult); err != nil {
		log.Error("failed to save currency exchange changes in the database", "error", err)
		return nil, err
	}

	// updating the current balance for the response
	balanceUser[req.FromCurrency] = exchangeResult.NewBaseBalance
	balanceUser[req.ToCurrency] = exchangeResult.NewToBalance

	// preparing response
	var resp models.ExchangeResponse
	resp.Message = "currency exchange successfully"
	resp.ExchangeRate = rate.Rate
	resp.SpentAccoutn = models.SpentAccoutn{Currency: req.FromCurrency, Amount: req.Amount}
	resp.ReceivedAccount = models.ReceivedAccount{Currency: req.ToCurrency, Amount: exchangeResult.Received}
	resp.NewBalance = balanceUser

	log.Info("successfully exchange")
	return &resp, nil
}

func (w *Wallet) ExchangeRates(ctx context.Context) (*models.ExchangeRatesResponse, error) {
	op := "service Wallet: obtaining all exchange rates"
	log := w.log.With(slog.String("operation", op))
	log.Debug("ExchangeRates func call")

	var resp models.ExchangeRatesResponse

	if err := w.clientGRPC.GetAllRates(ctx, &resp); err != nil {
		log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	resp.Message = "data successfully received"

	log.Info("all exchange rates have been successfully received")
	return &resp, nil
}

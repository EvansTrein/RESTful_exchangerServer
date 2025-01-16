package services

import (
	"context"
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/config"
	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
)

type Wallet struct {
	log        *slog.Logger
	clientGRPC grpcclient.ClientGRPC
	db         storages.StoreWallet
}

func New(log *slog.Logger, db storages.StoreWallet, conf *config.Services) *Wallet {
	log.Debug("service Wallet: started creating")

	client, err := grpcclient.New(log, conf.AddressGRPC, conf.PortGRPC)
	if err != nil {
		panic(err)
	}

	log.Info("service Wallet: successfully created")
	return &Wallet{
		log:        log,
		clientGRPC: client,
		db:         db,
	}
}

func (w *Wallet) Stop() error {
	w.log.Debug("service Wallet: stop started")

	if err := w.clientGRPC.Close(); err != nil {
		w.log.Error("failed to stop the Wallet service", "error", err)
		return err
	}

	w.clientGRPC = nil
	w.db = nil

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

func (w *Wallet) Deposit(ctx context.Context, req models.DepositRequest) (*models.DepositResponse, error) {
	op := "service Wallet: account replenishment"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Deposit func call", slog.Any("requets data", req))

	newBalance, err := w.db.Deposit(ctx, req)
	if err != nil {
		log.Error("failed to deposit in the database", "error", err)
		return nil, err
	}

	log.Debug("account was successfully funded", "new balance", newBalance)

	var resp models.DepositResponse
	resp.Message = "account topped up successfully"
	resp.NewBalance = newBalance

	log.Info("account topped up successfully")
	return &resp, nil
}

func (w *Wallet) Withdraw(ctx context.Context, req models.WithdrawRequest) (*models.WithdrawResponse, error) {
	op := "service Wallet: withdraw request received"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Withdraw func call", slog.Any("requets data", req))

	newBalance, err := w.db.Withdraw(ctx, req)
	if err != nil {
		log.Error("failed to withdraw in the database", "error", err)
		return nil, err
	}

	var resp models.WithdrawResponse
	resp.Message = "successfully withdrawn"
	resp.NewBalance = newBalance

	return &resp, nil
}

func (w *Wallet) Exchange(ctx context.Context, req models.ExchangeRequest) (*models.ExchangeResponse, error) {
	op := "service Wallet: currency exchange request"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Exchange func call", slog.Any("requets data", req))

	// TODO: проверить есть ли курс в Redis

	var resp models.ExchangeResponse
	var rate models.ExchangeGRPC

	rate.FromCurrency = req.FromCurrency
	rate.ToCurrency = req.ToCurrency

	if err := w.clientGRPC.ExchangeRate(ctx, &rate); err != nil {
		log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	log.Debug("exchange rate successfully received from gRPC server", "rate from gRPC", rate)
	// TODO: тут получен курс, сохранить его в Redis
	// TODO: тут получен курс, далее обращение к БД

	resp.Message = "data successfully received"

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

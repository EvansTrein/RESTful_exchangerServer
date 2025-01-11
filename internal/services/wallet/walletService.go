package services

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
)

type Wallet struct {
	log        *slog.Logger
	clientGRPC grpcclient.ClientGRPC
	db         storages.StoreWallet
}

func New(log *slog.Logger, db storages.StoreWallet, pathGRPC string) *Wallet {
	log.Debug("service Wallet: started creating")
	
	client, err := grpcclient.New(log, pathGRPC)
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

func (w *Wallet) Balance(req models.BalanceRequest) (*models.BalanceResponse, error) {

	return &models.BalanceResponse{}, nil
}

func (w *Wallet) Deposit(req models.DepositRequest) (*models.DepositResponse, error) {

	return &models.DepositResponse{}, nil
}

func (w *Wallet) Withdraw(req models.WithdrawRequest) (*models.WithdrawResponse, error) {
	
	return &models.WithdrawResponse{}, nil
}

func (w *Wallet) Exchange(req models.ExchangeRequest) (*models.ExchangeResponse, error) {
	op := "service Wallet: currency exchange request"
	log := w.log.With(slog.String("operation", op))
	log.Debug("Exchange func call", slog.Any("requets data", req))

	var resp models.ExchangeResponse
	var rate models.ExchangeGRPC

	rate.FromCurrency = req.FromCurrency
	rate.ToCurrency = req.ToCurrency

	if err := w.clientGRPC.ExchangeRate(&rate); err != nil {
		log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	log.Debug("exchange rate successfully received from gRPC server", "rate from gRPC", rate)
	// TODO: тут получен курс, сохранить его в Redis
	// TODO: тут получен курс, далее обращение к БД

	resp.Message = "data successfully received"
	
	return &resp, nil
}

func (w *Wallet) ExchangeRates() (*models.ExchangeRatesResponse, error) {
	op := "service Wallet: obtaining all exchange rates"
	log := w.log.With(slog.String("operation", op))
	log.Debug("ExchangeRates func call")

	var resp models.ExchangeRatesResponse

	if err := w.clientGRPC.GetAllRates(&resp); err != nil {
		log.Error("failed to get data from GRPC server", "error", err)
		return nil, err
	}

	resp.Message = "data successfully received"

	log.Info("all exchange rates have been successfully received")
	return &resp, nil
}

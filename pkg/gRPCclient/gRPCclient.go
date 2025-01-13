package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	pb "github.com/EvansTrein/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	gRPCTimeoutMethodCall = time.Second * 5

	ErrServerUnavailable = errors.New("gRPC server is unavailable")
	ErrServerTimeOut     = errors.New("gRPC method call execution timeout expired")
	ErrServerNotCurrency = errors.New("gRPC currency is not supported")
)

type ClientGRPC interface {
	GetAllRates(req *models.ExchangeRatesResponse) error
	ExchangeRate(req *models.ExchangeGRPC) error
	Close() error
}

type ServerGRPC struct {
	log     *slog.Logger
	conn    *grpc.ClientConn
}

func New(log *slog.Logger, address, port string) (*ServerGRPC, error) {
	grpcAddr := fmt.Sprintf("%s:%s", address, port)
	log.Debug("gRPC server: started creating", "address", grpcAddr)

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create a client for gRPC server", "error", err)
		return nil, err
	}

	log.Info("gRPC server: successfully created")
	return &ServerGRPC{log: log, conn: conn}, nil
}

func (s *ServerGRPC) Close() error {
	s.log.Debug("gRPC server: stop started")

	if err := s.conn.Close(); err != nil {
		s.log.Error("failed to close connection with gRPC server")
		return err
	}

	s.conn = nil

	s.log.Info("gRPC server: stop successful")
	return nil
}

func (s *ServerGRPC) GetAllRates(req *models.ExchangeRatesResponse) error {
	op := "gRPC server: obtaining all exchange rates"
	log := s.log.With(slog.String("operation", op))
	log.Debug("GetAllRates func call")

	client := pb.NewExchangeServiceClient(s.conn)

	ctx, cancel := context.WithTimeout(context.Background(), gRPCTimeoutMethodCall)
	defer cancel()

	callGRPC, err := client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Warn("timeout time for response from gRPC server has expired")
			return ErrServerTimeOut
		} else {
			log.Warn("failed to get data from gRPC server", "error", err)
			return fmt.Errorf("failed %s, error: %s", op, err.Error())
		}
	}

	if callGRPC.GetRates() == nil {
		log.Warn(ErrServerUnavailable.Error())
		return ErrServerUnavailable
	}

	req.Rates = callGRPC.GetRates()
	log.Info("data from gRPC server successfully received")
	return nil
}

func (s *ServerGRPC) ExchangeRate(req *models.ExchangeGRPC) error {
	op := "gRPC server: currency exchange rate request"
	log := s.log.With(slog.String("operation", op))
	log.Debug("ExchangeRate func call")

	client := pb.NewExchangeServiceClient(s.conn)

	ctx, cancel := context.WithTimeout(context.Background(), gRPCTimeoutMethodCall)
	defer cancel()

	var reqForGRPC pb.CurrencyRequest

	reqForGRPC.FromCurrency = req.FromCurrency
	reqForGRPC.ToCurrency = req.ToCurrency

	callGRPC, err := client.GetExchangeRateForCurrency(ctx, &reqForGRPC)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Warn("timeout time for response from gRPC server has expired")
			return ErrServerTimeOut
		} else if status.Code(err) == codes.NotFound {
			log.Warn("currency that is not supported has been requested")
			return ErrServerNotCurrency
		} else {
			log.Warn("failed to get data from gRPC server", "error", err)
			return fmt.Errorf("failed %s, error: %s", op, err.Error())
		}
	}

	if callGRPC.GetRate() == 0 {
		log.Warn(ErrServerUnavailable.Error())
		return ErrServerUnavailable
	}

	req.Rate = callGRPC.GetRate()
	log.Info("data from gRPC server successfully received")
	return nil
}

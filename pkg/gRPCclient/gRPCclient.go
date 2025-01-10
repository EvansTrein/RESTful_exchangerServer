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
	"google.golang.org/grpc/credentials/insecure"
)

var (
	gRPCTimeoutMethodCall = time.Second * 5

	ErrServerUnavailable = errors.New("gRPC server is unavailable")
	ErrServerTimeOut     = errors.New("gRPC method call execution timeout expired")
)

type ClientGRPC interface {
}

type ServerGRPC struct {
	log  *slog.Logger
	conn *grpc.ClientConn
}

func New(log *slog.Logger, grpcAddr string) (*ServerGRPC, error) {

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to create a client for gRPC server", "error", err)
		return nil, err
	}

	return &ServerGRPC{log: log, conn: conn}, nil
}

func (s *ServerGRPC) Close() error {
	s.log.Info("Start closing the connection to the gRPC server")

	if err := s.conn.Close(); err != nil {
		s.log.Error("failed to close connection with gRPC server")
		return err
	}

	s.conn = nil
	
	s.log.Info("connection with gRPC server successfully closed")
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

	s.log.Info("data from gRPC server successfully received")
	req.Rates = callGRPC.GetRates()

	return nil
}

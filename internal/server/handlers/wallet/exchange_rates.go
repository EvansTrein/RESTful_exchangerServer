package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
	"github.com/gin-gonic/gin"
)

type exchangeRatesServ interface {
	ExchangeRates(ctx context.Context) (*models.ExchangeRatesResponse, error)
}

// ExchangeRates is a Gin handler function that retrieves the current exchange all rates for supported currencies.
// It calls the service to fetch the exchange rates and returns the result.
// If the gRPC server is unavailable or the request times out, it returns a 503 Service Unavailable or 504 Gateway Timeout.
// On success, it returns a 200 OK response with the exchange rates.
//
// @Summary Get all exchange rates
// @Description Get the current exchange rates for supported currencies
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ExchangeRatesResponse
// @Failure 401 {object} models.HandlerResponse
// @Failure 500 {object} models.HandlerResponse
// @Failure 503 {object} models.HandlerResponse
// @Failure 504 {object} models.HandlerResponse
// @Router /exchange/rates [get]
func ExchangeRates(log *slog.Logger, serv exchangeRatesServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler ExchangeRates: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request received")

		result, err := serv.ExchangeRates(ctx.Request.Context())
		if err != nil {
			switch err {
			case grpcclient.ErrServerUnavailable:
				log.Warn("failed to send data", "error", err)
				ctx.JSON(503, models.HandlerResponse{
					Status: http.StatusServiceUnavailable, 
					Error: err.Error(), 
					Message: "failed to retrieve data",
				})
				return
			case grpcclient.ErrServerTimeOut:
				log.Error("failed to send data", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status: http.StatusGatewayTimeout, 
					Error: err.Error(), 
					Message: "response timeout expired on the GRPC server side",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to send data", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to send data", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status: http.StatusInternalServerError, 
					Error: err.Error(), 
					Message: "failed to retrieve data",
				})
				return
			}
		}

		log.Info("data successfully sent")
		ctx.JSON(200, result)
	}
}

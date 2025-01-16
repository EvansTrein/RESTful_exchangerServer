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

// Метод: **GET**
// URL: **/api/v1/exchange/rates**
// Заголовки:
// _Authorization: Bearer JWT_TOKEN_

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//     "rates":
//     {
//       "USD": "float",
//       "RUB": "float",
//       "EUR": "float"
//     }
// }
// ```

// • Ошибка: ```500 Internal Server Error```
// ```json
// {
//   "error": "Failed to retrieve exchange rates"
// }
// ```

// ▎Описание

// Получение актуальных курсов валют из внешнего gRPC-сервиса.
// Возвращает курсы всех поддерживаемых валют.

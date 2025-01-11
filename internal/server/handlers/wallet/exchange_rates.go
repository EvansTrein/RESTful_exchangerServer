package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
	"github.com/gin-gonic/gin"
)

type exchangeRatesServ interface {
	ExchangeRates() (*models.ExchangeRatesResponse, error)
}

func ExchangeRates(log *slog.Logger, serv exchangeRatesServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler ExchangeRates: call"
		castLog := log.With(
			slog.String("operation", op), 
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),	
		)
		castLog.Debug("request received")

		result, err := serv.ExchangeRates()
		if err != nil {
			switch err {
			case grpcclient.ErrServerUnavailable:
				castLog.Warn("failed to send data", "error", err)
				ctx.JSON(503, models.HandlerResponse{Status: http.StatusServiceUnavailable, Error: err.Error()})
				return
			case grpcclient.ErrServerTimeOut:
				castLog.Warn("failed to send data", "error", err)
				ctx.JSON(504, models.HandlerResponse{Status: http.StatusGatewayTimeout, Error: err.Error()})
				return
			default:
				fmt.Println(err)
				castLog.Error("failed to send data", "error", err)
				ctx.JSON(500, models.HandlerResponse{Status: http.StatusInternalServerError, Error: err.Error()})
				return
			}
		}

		log.Info("data successfully sent")
		ctx.JSON(200, models.HandlerResponse{Status: http.StatusOK, Message: "data successfully sent", Data: result})
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

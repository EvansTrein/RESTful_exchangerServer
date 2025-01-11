package handlers

import (
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type exchangeServ interface {
	Exchange(req models.ExchangeRequest) (*models.ExchangeResponse, error)
}

func Exchange(log *slog.Logger, serv exchangeServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Exchange: call"
		castLog := log.With(
			slog.String("operation", op), 
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),	
		)
		castLog.Debug("request received")

		var req models.ExchangeRequest

		if err := ctx.BindJSON(&req); err != nil {
			castLog.Warn("fail BindJSON", "error", err)
			return
		} 

		castLog.Debug("request data has been successfully validated", "data", req)
		

		result, err := serv.Exchange(req)

		if err != nil {
			// TODO: вернуть 404 если запрошенной валюты нет
			castLog.Error("failed to send data", "error", err)
			ctx.JSON(500, models.HandlerResponse{Status: http.StatusInternalServerError, Error: err.Error()})
			return
		}

		ctx.JSON(200, models.HandlerResponse{Status: http.StatusOK, Message: "data successfully sent", Data: result})
	}
}

// Метод: **POST**
// URL: **/api/v1/exchange**
// Заголовки:
// _Authorization: Bearer JWT_TOKEN_

// Тело запроса:
// ```json
// {
//   "from_currency": "USD",
//   "to_currency": "EUR",
//   "amount": 100.00
// }
// ```

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//   "message": "Exchange successful",
//   "exchanged_amount": 85.00,
//   "new_balance":
//   {
//   "USD": 0.00,
//   "EUR": 85.00
//   }
// }
// ```

// • Ошибка: 400 Bad Request
// ```json
// {
//   "error": "Insufficient funds or invalid currencies"
// }
// ```

// ▎Описание

// Курс валют осуществляется по данным сервиса exchange (если в течении небольшого времени был запрос от клиента курса валют (**/api/v1/exchange**) до обмена, то
// брать курс из кэша, если же запроса курса валют не было или он запрашивался слишком давно, то нужно осуществить gRPC-вызов к внешнему сервису, который предоставляет актуальные курсы валют)
// Проверяется наличие средств для обмена, и обновляется баланс пользователя.

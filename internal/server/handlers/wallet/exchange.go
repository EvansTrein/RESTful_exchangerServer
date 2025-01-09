package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

func ExchangeHandler(log *slog.Logger, wallet services.WalletService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("ExchangeHandler")
		res, _ := wallet.Exchange(models.ExchangeRequest{})

		ctx.JSON(200, gin.H{"ExchangeHandler": res})
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

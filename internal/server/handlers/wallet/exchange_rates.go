package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/gin-gonic/gin"
)

func ExchangeRatesHandler(log *slog.Logger, wallet services.Walletervice) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("ExchangeRatesHandler")
		res, _ := wallet.ExchangeRates()

		ctx.JSON(200, gin.H{"ExchangeRatesHandler": res})
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

package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type depositServ interface {
	Deposit(req models.DepositRequest) (*models.DepositResponse, error)
}

func DepositHandler(log *slog.Logger, serv depositServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("DepositHandler")
		res, _ := serv.Deposit(models.DepositRequest{})

		ctx.JSON(200, gin.H{"DepositHandler": res})
	}
}

// Метод: **POST**
// URL: **/api/v1/wallet/deposit**
// Заголовки:
// _Authorization: Bearer JWT_TOKEN_

// Тело запроса:
// ```
// {
//   "amount": 100.00,
//   "currency": "USD" // (USD, RUB, EUR)
// }
// ```

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//   "message": "Account topped up successfully",
//   "new_balance": {
//     "USD": "float",
//     "RUB": "float",
//     "EUR": "float"
//   }
// }
// ```

// • Ошибка: ```400 Bad Request```
// ```json
// {
// "error": "Invalid amount or currency"
// }
// ```

// ▎Описание

// Позволяет пользователю пополнить свой счет. Проверяется корректность суммы и валюты.
// Обновляется баланс пользователя в базе данных.

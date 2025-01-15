package handlers

import (
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type depositServ interface {
	Deposit(req models.DepositRequest) (*models.DepositResponse, error)
}

func Deposit(log *slog.Logger, serv depositServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Deposit: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("account top-up request received")

		var req models.DepositRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Warn("fail BindJSON", "error", err)
			ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error(), Message: "invalid data"})
			return
		}

		log.Debug("request data has been successfully validated", "data", req)
		
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(500, models.HandlerResponse{
                Status:  http.StatusInternalServerError,
                Error:   "userID not found in context",
                Message: "failed to retrieve user id from context",
            })
            return
		}

		userIdUint, ok := userID.(uint)
        if !ok {
            ctx.JSON(500, models.HandlerResponse{
                Status:  http.StatusInternalServerError,
                Error:   "invalid userID type in context",
                Message: "failed to convert user id to the required data type",
            })
            return
        }

		req.UserID = userIdUint
		log.Debug("user id was successfully obtained from the context and added to the request")

		result, err := serv.Deposit(req)
		if err != nil {
			log.Error("failed to send data", "error", err)
			ctx.JSON(500, models.HandlerResponse{
				Status: http.StatusInternalServerError, 
				Error: err.Error(),
				Message: "failed to deposit",
			})
			return
		}

		log.Info("deposit successful")
		ctx.JSON(200, result)
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

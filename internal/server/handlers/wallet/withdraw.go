package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/storages"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type withdrawServ interface {
	Withdraw(ctx context.Context, req *models.AccountOperationRequest) (*models.AccountOperationResponse, error)
}

func Withdraw(log *slog.Logger, serv withdrawServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Withdraw: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("debit withdrawal")

		var req models.AccountOperationRequest
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

		result, err := serv.Withdraw(ctx.Request.Context(), &req)
		if err != nil {
			switch err {
			case storages.ErrInsufficientFunds:
				log.Warn("failed to withdraw", "error", err)
				ctx.JSON(402, models.HandlerResponse{
					Status:  http.StatusPaymentRequired,
					Error:   err.Error(),
					Message: "insufficient funds",
				})
				return
			case storages.ErrCurrencyNotFound:
				log.Warn("failed to withdraw", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "currency is not supported",
				})
				return
			case storages.ErrAccountNotFound:
				log.Error("failed to withdraw", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "no account in the specified currency",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to withdraw", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to withdraw", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "failed to withdraw",
				})
				return
			}
		}

		log.Info("withdraw successful")
		ctx.JSON(200, result)
	}
}

// Метод: **POST**
// URL: **/api/v1/wallet/withdraw**
// Заголовки:
// _Authorization: Bearer JWT_TOKEN_

// Тело запроса:
// ```
// {
//     "amount": 50.00,
//     "currency": "USD" // USD, RUB, EUR)
// }
// ```

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//   "message": "Withdrawal successful",
//   "new_balance": {
//     "USD": "float",
//     "RUB": "float",
//     "EUR": "float"
//   }
// }
// ```

// • Ошибка: 400 Bad Request
// ```json
// {
//   "error": "Insufficient funds or invalid amount"
// }
// ```

// ▎Описание

// Позволяет пользователю вывести средства со своего счета.
// Проверяется наличие достаточного количества средств и корректность суммы.

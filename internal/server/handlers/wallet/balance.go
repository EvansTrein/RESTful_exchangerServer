package handlers

import (
	"log/slog"
	"net/http"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type balanceServ interface {
	Balance(req models.BalanceRequest) (*models.BalanceResponse, error)
}

func Balance(log *slog.Logger, serv balanceServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Balance: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request for balance of all accounts")

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

		var req models.BalanceRequest
		req.UserID = userIdUint
		log.Debug("user id was successfully obtained from the context and added to the request", "userID", userIdUint)

		result, err := serv.Balance(req)
		if err != nil {
			log.Error("failed to send data", "error", err)
			ctx.JSON(500, models.HandlerResponse{
				Status: http.StatusInternalServerError, 
				Error: err.Error(),
				Message: "failed to send data",
			})
			return
		}

		log.Info("data successfully sent")
		ctx.JSON(200, result)
	}
}

// Метод: **GET**
// URL: **/api/v1/balance**
// Заголовки:
// _Authorization: Bearer JWT_TOKEN_

// Ответ:

// • Успех: ```200 OK```

// ```json
// {
//   "balance":
//   {
//   "USD": "float",
//   "RUB": "float",
//   "EUR": "float"
//   }
// }
// ``

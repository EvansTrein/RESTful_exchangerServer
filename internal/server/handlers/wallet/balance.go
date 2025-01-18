package handlers

import (
	"log/slog"
	"net/http"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type balanceServ interface {
	Balance(ctx context.Context, req models.BalanceRequest) (*models.BalanceResponse, error)
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

		result, err := serv.Balance(ctx.Request.Context(), req)
		if err != nil {
			switch err {
			case services.ErrUserNotFound:
				log.Error("user not found", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status: http.StatusNotFound, 
					Error: err.Error(), 
					Message: "the balance of a non-existent user was requested",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to send data", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status: http.StatusGatewayTimeout, 
					Error: err.Error(), 
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to send data", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status: http.StatusInternalServerError, 
					Error: err.Error(),
					Message: "failed to send data",
				})
				return
			}
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

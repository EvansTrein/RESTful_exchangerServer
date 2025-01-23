package handlers

import (
	"context"
	"log/slog"
	"net/http"

	services "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type depositServ interface {
	Deposit(ctx context.Context, req *models.AccountOperationRequest) (*models.AccountOperationResponse, error)
}

// Deposit is a Gin handler function that handles depositing funds into a user's account.
// It binds the incoming JSON request to a struct, validates the data, and calls the service to deposit funds.
// If the data is invalid, it returns a 400 Bad Request.
// If the user ID is missing or invalid, it returns a 500 Internal Server Error.
// If the currency or account is not found, it returns a 404 Not Found.
// If the request times out, it returns a 504 Gateway Timeout.
// On success, it returns a 200 OK response with the deposit result.
//
// @Summary Deposit funds into an account
// @Description Deposit funds into a user's account for a specific currency
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.AccountOperationRequest true "Deposit request"
// @Success 200 {object} models.AccountOperationResponse
// @Failure 400 {object} models.HandlerResponse
// @Failure 401 {object} models.HandlerResponse
// @Failure 404 {object} models.HandlerResponse
// @Failure 500 {object} models.HandlerResponse
// @Failure 504 {object} models.HandlerResponse
// @Router /wallet/deposit [post]
func Deposit(log *slog.Logger, serv depositServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Deposit: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("account top-up request received")

		var req models.AccountOperationRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Warn("fail BindJSON", "error", err)
			ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error(), Message: "invalid data"})
			return
		}

		log.Debug("request data has been successfully validated", "data", req)

		// get user id from context
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

		result, err := serv.Deposit(ctx.Request.Context(), &req)
		if err != nil {
			switch err {
			case services.ErrCurrencyNotFound:
				log.Warn("failed to deposit", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "currency not found",
				})
				return
			case services.ErrAccountNotFound:
				log.Error("failed to deposit", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "account not found",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to deposit", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to deposit", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "failed to deposit",
				})
				return
			}
		}

		log.Info("deposit successful")
		ctx.JSON(200, result)
	}
}

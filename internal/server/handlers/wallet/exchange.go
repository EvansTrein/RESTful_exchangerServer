package handlers

import (
	"log/slog"
	"net/http"

	servWallet "github.com/EvansTrein/RESTful_exchangerServer/internal/services/wallet"
	servAuth "github.com/EvansTrein/RESTful_exchangerServer/internal/services/auth"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	grpcclient "github.com/EvansTrein/RESTful_exchangerServer/pkg/gRPCclient"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type exchangeServ interface {
	Exchange(ctx context.Context, req models.ExchangeRequest) (*models.ExchangeResponse, error)
}

// Exchange is a Gin handler function that handles currency exchange for the authenticated user.
// It binds the incoming JSON request to a struct, validates the data, and calls the service to perform the exchange.
// If the data is invalid or the currencies are the same, it returns a 400 Bad Request.
// If the user ID is missing or invalid, it returns a 500 Internal Server Error.
// If the user, account, or currency is not found, it returns a 404 Not Found.
// If there are insufficient funds, it returns a 402 Payment Required.
// If the request times out or the gRPC server is unavailable, it returns a 504 Gateway Timeout or 503 Service Unavailable.
// On success, it returns a 200 OK response with the exchange result.
//
// @Summary Exchange currency
// @Description Exchange one currency to another for the authenticated user
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.ExchangeRequest true "Exchange request"
// @Success 200 {object} models.ExchangeResponse
// @Failure 400 {object} models.HandlerResponse
// @Failure 401 {object} models.HandlerResponse
// @Failure 402 {object} models.HandlerResponse
// @Failure 404 {object} models.HandlerResponse
// @Failure 500 {object} models.HandlerResponse
// @Failure 503 {object} models.HandlerResponse
// @Failure 504 {object} models.HandlerResponse
// @Router /exchange [post]
func Exchange(log *slog.Logger, serv exchangeServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "Handler Exchange: call"
		log = log.With(
			slog.String("operation", op),
			slog.String("apiPath", ctx.FullPath()),
			slog.String("HTTP Method", ctx.Request.Method),
		)
		log.Debug("request received")

		var req models.ExchangeRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Warn("fail BindJSON", "error", err)
			ctx.JSON(400, models.HandlerResponse{Status: http.StatusBadRequest, Error: err.Error(), Message: "invalid data"})
			return
		}

		if req.FromCurrency == req.ToCurrency {
			log.Warn("same currencies")
			ctx.JSON(400, models.HandlerResponse{
				Status: http.StatusBadRequest,
				Error: "same currency is specified for buying and selling",
				Message: "invalid data",
			})
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
		log.Debug("user id was successfully obtained from the context and added to the request", "userID", userIdUint)

		result, err := serv.Exchange(ctx.Request.Context(), req)
		if err != nil {
			switch err {
			case servWallet.ErrInsufficientFunds:
				log.Warn("failed to exchanged", "error", err)
				ctx.JSON(402, models.HandlerResponse{
					Status:  http.StatusPaymentRequired,
					Error:   err.Error(),
					Message: "insufficient funds",
				})
				return
			case servAuth.ErrUserNotFound:
				log.Error("request was received for a non-existent user", "error", err, "user id", userIdUint)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "user does not exist",
				})
				return
			case servWallet.ErrAccountNotFound:
				log.Warn("failed to exchanged", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "no account in this currency",
				})
				return
			case servWallet.ErrCurrencyNotFound:
				log.Warn("failed to exchanged", "error", err)
				ctx.JSON(404, models.HandlerResponse{
					Status:  http.StatusNotFound,
					Error:   err.Error(),
					Message: "currency is not supported",
				})
				return
			case grpcclient.ErrServerUnavailable:
				log.Warn("failed to exchanged", "error", err)
				ctx.JSON(503, models.HandlerResponse{
					Status:  http.StatusServiceUnavailable,
					Error:   err.Error(),
					Message: "failed to exchanged",
				})
				return
			case grpcclient.ErrServerTimeOut:
				log.Error("failed to exchanged", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "response timeout expired on the GRPC server side",
				})
				return
			case context.DeadlineExceeded:
				log.Error("failed to exchanged", "error", err)
				ctx.JSON(504, models.HandlerResponse{
					Status:  http.StatusGatewayTimeout,
					Error:   err.Error(),
					Message: "the waiting time for a response from the internal service has expired",
				})
				return
			default:
				log.Error("failed to exchanged", "error", err)
				ctx.JSON(500, models.HandlerResponse{
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "failed to exchanged",
				})
				return
			}
		}

		log.Info("currency exchange successfully")
		ctx.JSON(200, result)
	}
}

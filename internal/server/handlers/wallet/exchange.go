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

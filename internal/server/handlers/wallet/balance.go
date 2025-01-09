package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

func BalanceHandler(log *slog.Logger, wallet services.WalletService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("BalanceHandler")
		res, _ := wallet.Balance(models.BalanceRequest{})

		ctx.JSON(200, gin.H{"BalanceHandler": res})
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

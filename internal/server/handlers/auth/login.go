package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

type loginServ interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
}

func Login(log *slog.Logger, serv loginServ) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("LoginHandler")
		res, _ := serv.Login(models.LoginRequest{})

		ctx.JSON(200, gin.H{"Login": res})
	}
}

// Метод: **POST**
// URL: **/api/v1/login**
// Тело запроса:
// ```json
// {
// "username": "string",
// "password": "string"
// }
// ```

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//   "token": "JWT_TOKEN"
// }
// ```

// • Ошибка: ```401 Unauthorized```
// ```json
// {
//   "error": "Invalid username or password"
// }
// ```

// ▎Описание

// Авторизация пользователя.
// При успешной авторизации возвращается JWT-токен, который будет использоваться для аутентификации последующих запросов.

package handlers

import (
	"log/slog"

	"github.com/EvansTrein/RESTful_exchangerServer/internal/services"
	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gin-gonic/gin"
)

func LoginHandler(log *slog.Logger, auth services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debug("LoginHandler")
		res, _ := auth.Login(models.LoginRequest{})

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

package handlers

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
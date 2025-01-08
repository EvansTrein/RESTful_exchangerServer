package handlers

// Метод: **POST**  
// URL: **/api/v1/register**  
// Тело запроса:
// ```json
// {
//   "username": "string",
//   "password": "string",
//   "email": "string"
// }
// ```

// Ответ:
// • Успех: ```201 Created```
// ```json
// { 
//   "message": "User registered successfully"
// }
// ```

// • Ошибка: ```400 Bad Request```
// ```json
// {
//   "error": "Username or email already exists"
// }
// ```

// ▎Описание

// Регистрация нового пользователя. 
// Проверяется уникальность имени пользователя и адреса электронной почты.
// Пароль должен быть зашифрован перед сохранением в базе данных.
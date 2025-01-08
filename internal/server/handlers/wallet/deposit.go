package handlers

// Метод: **POST**  
// URL: **/api/v1/wallet/deposit**  
// Заголовки:  
// _Authorization: Bearer JWT_TOKEN_

// Тело запроса:
// ```
// {
//   "amount": 100.00,
//   "currency": "USD" // (USD, RUB, EUR)
// }
// ```

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//   "message": "Account topped up successfully",
//   "new_balance": {
//     "USD": "float",
//     "RUB": "float",
//     "EUR": "float"
//   }
// }
// ```

// • Ошибка: ```400 Bad Request```
// ```json
// {
// "error": "Invalid amount or currency"
// }
// ```

// ▎Описание

// Позволяет пользователю пополнить свой счет. Проверяется корректность суммы и валюты.
// Обновляется баланс пользователя в базе данных.
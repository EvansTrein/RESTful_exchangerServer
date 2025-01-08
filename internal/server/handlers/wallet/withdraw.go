package handlers

// Метод: **POST**  
// URL: **/api/v1/wallet/withdraw**  
// Заголовки:  
// _Authorization: Bearer JWT_TOKEN_

// Тело запроса:
// ```
// {
//     "amount": 50.00,
//     "currency": "USD" // USD, RUB, EUR)
// }
// ```

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//   "message": "Withdrawal successful",
//   "new_balance": {
//     "USD": "float",
//     "RUB": "float",
//     "EUR": "float"
//   }
// }
// ```

// • Ошибка: 400 Bad Request
// ```json
// {
//   "error": "Insufficient funds or invalid amount"
// }
// ```

// ▎Описание

// Позволяет пользователю вывести средства со своего счета.
// Проверяется наличие достаточного количества средств и корректность суммы.
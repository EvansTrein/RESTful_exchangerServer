package handlers

// Метод: **GET**  
// URL: **/api/v1/exchange/rates**  
// Заголовки:  
// _Authorization: Bearer JWT_TOKEN_

// Ответ:

// • Успех: ```200 OK```
// ```json
// {
//     "rates": 
//     {
//       "USD": "float",
//       "RUB": "float",
//       "EUR": "float"
//     }
// }
// ```

// • Ошибка: ```500 Internal Server Error```
// ```json
// {
//   "error": "Failed to retrieve exchange rates"
// }
// ```

// ▎Описание

// Получение актуальных курсов валют из внешнего gRPC-сервиса.
// Возвращает курсы всех поддерживаемых валют.
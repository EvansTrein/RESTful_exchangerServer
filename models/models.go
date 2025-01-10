package models

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type BalanceRequest struct {
	UserID int `json:"id"`
}

type BalanceResponse struct {
	Balance map[string]float32 `json:"balance"`
}

type DepositRequest struct {
	UserID   int     `json:"id"`
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type DepositResponse struct {
	Message    string          `json:"message"`
	NewBalance BalanceResponse `json:"new_balance"`
}

type ExchangeRatesResponse struct {
	Rates map[string]float32 `json:"rates"`
}

type ExchangeRequest struct {
	UserID       int     `json:"id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float32 `json:"amount"`
}

type ExchangeResponse struct {
	Message         string          `json:"message"`
	ExchangedAmount float32         `json:"exchanged_amount"`
	NewBalance      BalanceResponse `json:"new_balance"`
}

type WithdrawRequest struct {
	UserID   int     `json:"id"`
	Currency string  `json:"currency"`
	Amount   float32 `json:"amount"`
}

type WithdrawResponse struct {
	Message    string          `json:"message"`
	NewBalance BalanceResponse `json:"new_balance"`
}

type HandlerResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

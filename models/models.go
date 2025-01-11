package models

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Name         string `json:"username" binding:"required,min=3"`
	HashPassword string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	UserID  uint   `json:"id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Name     string `json:"username"`
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
	Message string             `json:"message"`
	Rates   map[string]float32 `json:"rates"`
}

type ExchangeRequest struct {
	UserID       int     `json:"id" binding:"required"`
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Amount       float32 `json:"amount" binding:"required"`
}

type ExchangeGRPC struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Rate         float32 `json:"rate" binding:"required"`
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
}

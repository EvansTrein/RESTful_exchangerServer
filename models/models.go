package models

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	HashPassword string `json:"password"`
}

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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PayloadToken struct {
	UserID uint `json:"id"`
}

type BalanceRequest struct {
	UserID uint `json:"id"`
}

type BalanceResponse struct {
	Balance map[string]float32 `json:"balance"`
}

type DepositRequest struct {
	UserID   uint    `json:"id"`
	Amount   float32 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,min=3"`
}

type DepositResponse struct {
	Message    string             `json:"message"`
	NewBalance map[string]float32 `json:"new_balance"`
}

type ExchangeRatesResponse struct {
	Message string             `json:"message"`
	Rates   map[string]float32 `json:"rates"`
}

type ExchangeRequest struct {
	UserID       uint    `json:"id"`
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
	Message         string             `json:"message"`
	ExchangedAmount float32            `json:"exchanged_amount"`
	NewBalance      map[string]float32 `json:"new_balance"`
}

type WithdrawRequest struct {
	UserID   uint    `json:"id"`
	Amount   float32 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,min=3"`
}

type WithdrawResponse struct {
	Message    string             `json:"message"`
	NewBalance map[string]float32 `json:"new_balance"`
}

type HandlerResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

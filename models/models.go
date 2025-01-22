package models

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	HashPassword string `json:"password"`
}

type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Name         string `json:"username" binding:"required,min=3" example:"john"`
	HashPassword string `json:"password" binding:"required,min=6" example:"123456"`
}

type RegisterResponse struct {
	UserID  uint   `json:"id" example:"1"`
	Message string `json:"message" example:"user successfully created"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginResponse struct {
	Token string `json:"token" example:"JWT-token"`
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

type AccountOperationRequest struct {
	UserID    uint    `json:"-"`
	Amount    float32 `json:"amount" binding:"required,gt=0" example:"2000"`
	Currency  string  `json:"currency" binding:"required,min=3,max=6" example:"USD"`
	Operation string  `json:"-"`
}

type AccountOperationResponse struct {
	Message    string             `json:"message" example:"text message"`
	NewBalance map[string]float32 `json:"new_balance"`
}

type ExchangeRatesResponse struct {
	Message string             `json:"message" example:"text message"`
	Rates   map[string]float32 `json:"rates"`
}

type ExchangeRequest struct {
	UserID       uint    `json:"-"`
	FromCurrency string  `json:"from_currency" binding:"required,min=3,max=6" example:"USD"`
	ToCurrency   string  `json:"to_currency" binding:"required,min=3,max=6" example:"CNY"`
	Amount       float32 `json:"amount" binding:"required,gt=0" example:"500"`
}

type ExchangeRate struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Rate         float32 `json:"rate" binding:"required"`
}

type ExchangeResponse struct {
	Message         string             `json:"message" example:"text message"`
	ExchangeRate    float32            `json:"exchange_rate" example:"7.424683"`
	SpentAccoutn    SpentAccoutn       `json:"spent_accoutn"`
	ReceivedAccount ReceivedAccount    `json:"received_account"`
	NewBalance      map[string]float32 `json:"new_balance"`
}

type HandlerResponse struct {
	Status  int    `json:"status"`
	Error   string `json:"error" example:"text error"`
	Message string `json:"message" example:"text message"`
}

type CurrencyExchangeData struct {
	BaseBalance  float32
	ToBalance    float32
	ExchangeRate float32
	Amount       float32
}

type CurrencyExchangeResult struct {
	UserID         uint
	BaseCurrency   string
	NewBaseBalance float32
	ToCurrency     string
	NewToBalance   float32
	Received       float32
}

type SpentAccoutn struct {
	Currency string  `json:"currency" example:"USD"`
	Amount   float32 `json:"amount" example:"500"`
}

type ReceivedAccount struct {
	Currency string  `json:"currency" example:"CNY"`
	Amount   float32 `json:"amount" example:"3636.30"`
}

package models

// Currency - модель пользователя
type CurrencyRequest struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
}

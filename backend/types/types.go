package types


// Struct for Metal API Response
type MetalResponse struct {
	Status string `json:"status"`
	Data   struct {
		BaseCurrency string             `json:"base_currency"`
		Rates        map[string]float64 `json:"rates"` // Maps "XAU" -> 2350.50
	} `json:"data"`
}

// Struct for Currency API response
type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}
// What User Sends
type CalculationRequest struct {
	Currency	   string  `json:"currency"`
	GoldGrams      float64 `json:"gold_grams"`
	SilverGrams    float64 `json:"silver_grams"`
	Cash           float64 `json:"cash"`
	BusinessAssets float64 `json:"business_assets"`
	Liabilities    float64 `json:"liabilities"`
}

// What Backend Sends after Calculating
type CalculationResponse struct {
	TotalAssets    float64 `json:"total_assets"`
	NisabThreshold float64 `json:"nisab_threshold"`
	ZakahPayable   float64 `json:"zakah_payable"`
	Currency       string  `json:"currency"`
	LocalCurrency  string `json:"local_currency"`
	Message        string  `json:"message"`
}
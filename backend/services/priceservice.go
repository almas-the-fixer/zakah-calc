package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/almas-the-fixer/zakah-calc/types"
)

/*
TODO:
- Add Logging where errors can occur before returning them
*/


func GetGoldSilverPrices() (float64, float64, error) {
	// move to main.go later on as it will be the composition root
	apiKey := os.Getenv("APISED_SECRET_KEY")
	baseURL := "https://gold.g.apised.com/v1/latest"


	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Making Request
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return 0, 0, err
	}

	// Adding Query Parameters
	q := req.URL.Query()
	q.Add("metals", "XAU,XAG")
	q.Add("base_currency", "USD")
	q.Add("weight_unit", "gram")
	q.Add("currencies", "USD")
	req.URL.RawQuery = q.Encode() // Attaches ?metals=XAU... to the URL

	// Add Auth Headers - CRITICAL STEP
	req.Header.Add("x-api-key", apiKey)

	// Send the Request
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// G. Read and Parse Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var result types.MetalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, err
	}

	// Extract Prices (XAU is Gold, XAG is Silver)
	goldPrice := result.Data.MetalPrices["XAU"].Price
	silverPrice := result.Data.MetalPrices["XAG"].Price

	return goldPrice, silverPrice, nil
}


func GetExchangeRate(targetCurrency string) (float64, error) {
	// 1. If user wants USD, the rate is 1:1
	if targetCurrency == "USD" || targetCurrency == "" {
		return 1.0, nil
	}

	// 2. Call Free API (Frankfurter)
	// Note: Frankfurter supports major currencies (EUR, GBP, INR, CAD, etc.)
	// If we need every currency in the world, we might need a different API

	url := fmt.Sprintf("https://api.frankfurter.app/latest?from=USD&to=%s", targetCurrency)

	
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 3. Parse Response
	var result types.ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// 4. Return the specific rate
	rate, exists := result.Rates[targetCurrency]
	if !exists {
		// Fallback logic could go here, or return error
		return 0, fmt.Errorf("currency not supported")
	}

	return rate, nil
}
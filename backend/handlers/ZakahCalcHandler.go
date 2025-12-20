package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/almas-the-fixer/zakah-calc/types"
	"github.com/gofiber/fiber/v2"
)

// 1. Define the shape of the data you expect from the API
// (This struct matches the JSON response from APIsed)
// See types.go

// 2. The Function
func GetGoldSilverPrices() (float64, float64, error) {
	// A. Setup Configuration
	apiKey := os.Getenv("APISED_SECRET_KEY")
	postgres_user := os.Getenv("DB_USER")
	baseURL := "https://gold.g.apised.com/v1/latest"
	fmt.Println(postgres_user)
	// B. Create the Client (The Mailman) - Good practice to set a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// C. Create the Request (The Letter)
	// We use NewRequest so we can add headers later
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return 0, 0, err
	}

	// D. Add Query Parameters (The Address Details)
	q := req.URL.Query()
	q.Add("metals", "XAU,XAG")
	q.Add("base_currency", "USD")
	q.Add("weight_unit", "gram")
	q.Add("currencies", "USD")
	req.URL.RawQuery = q.Encode() // Attaches ?metals=XAU... to the URL

	// E. Add Headers (The Stamp/Auth) - CRITICAL STEP
	req.Header.Add("x-api-key", apiKey)

	// F. Send the Request
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close() // Always close the body when done!

	// G. Read and Parse Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var result types.MetalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, err
	}

	// H. Extract Prices (XAU is Gold, XAG is Silver)
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
	// If you need every currency in the world, we might need a different API


	url := fmt.Sprintf("https://api.frankfurter.app/latest?from=USD&to=%s", targetCurrency)

	resp, err := http.Get(url) // Simple Get is fine here
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



// 1. Define the Input (What the user sends)
// See types.go request struct
// 2. Define the Output (What we send back)
// See types.go response struct


// CalculateZakah godoc
// @Summary      Calculate Zakah
// @Description  Takes user assets and liabilities, converts currency, and calculates Zakah due.
// @Tags         Calculator
// @Accept       json
// @Produce      json
// @Param        request body types.CalculationRequest true "Calculation Request"
// @Success      200 {object} types.CalculationResponse
// @Failure      400 {object} map[string]interface{}
// @Router       /calculate-zakah [post]
func CalculateZakah(c *fiber.Ctx) error {
	// STEP 1: Parse Input
	input := new(types.CalculationRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input data"})
	}

	// STEP 2: Get Base Data
	goldPriceUSD, silverPriceUSD, err := GetGoldSilverPrices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch metal prices"})
	}

	exchangeRate, err := GetExchangeRate(input.Currency)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to get exchange rate"})
	}

	// STEP 3: The Correct Math
	
	// A. Convert UNIT PRICES to Local Currency
	localGoldPricePerGram := goldPriceUSD * exchangeRate
	localSilverPricePerGram := silverPriceUSD * exchangeRate

	// B. Calculate Total Asset Values
	goldWealth := input.GoldGrams * localGoldPricePerGram
	silverWealth := input.SilverGrams * localSilverPricePerGram
	
	// C. Total Net Wealth
	totalAssets := goldWealth + silverWealth + input.Cash + input.BusinessAssets
	netWealth := totalAssets - input.Liabilities

	// D. Calculate Nisab (Threshold) using Local Price
	// Nisab = 595g of Silver
	nisabThreshold := 595.0 * localSilverPricePerGram 

	// E. Calculate Zakah
	zakahDue := 0.0
	message := "No Zakah Due"

	if netWealth >= nisabThreshold {
		zakahDue = netWealth * 0.025 // 2.5%
		message = "Zakah is applicable"
	}

	// STEP 5: Response
	response := types.CalculationResponse{
		TotalAssets:    netWealth,
		NisabThreshold: nisabThreshold,
		ZakahPayable:   zakahDue,
		Currency:       input.Currency, // The values are now in this currency
		LocalCurrency:  input.Currency, 
		Message:        message,
	}

	return c.JSON(response)
}

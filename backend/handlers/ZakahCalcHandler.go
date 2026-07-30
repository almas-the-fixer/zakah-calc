package handlers

import (
	"github.com/almas-the-fixer/zakah-calc/services"
	"github.com/almas-the-fixer/zakah-calc/types"
	"github.com/almas-the-fixer/zakah-calc/zakah"
	"github.com/gofiber/fiber/v2"
)

// 1. Define the shape of the data you expect from the API
// (This struct matches the JSON response from APIsed)
// See types.go

// 2. The Function

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
func ZakahHandler(c *fiber.Ctx) error {
	// STEP 1: Parse Input
	input := new(types.CalculationRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input data"})
	}

	// STEP 2: Get Base Data
	goldPriceUSD, silverPriceUSD, err := services.GetGoldSilverPrices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch metal prices"})
	}
	targetCurrency := input.Currency
	exchangeRate, err := services.GetExchangeRate(targetCurrency)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to get exchange rate"})
	}

	netWealth, nisabThreshold := zakah.ComputeNetWealth(*input, goldPriceUSD, silverPriceUSD, exchangeRate)

	// Calculate Zakah
	zakahDue := zakah.CalculateZakah(netWealth, nisabThreshold, targetCurrency)
	return c.Status(200).JSON(zakahDue)
}

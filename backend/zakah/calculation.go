package zakah

import (
	"github.com/almas-the-fixer/zakah-calc/types"
)

const (
	silverNisab = 595.0
	//goldNisab   = 85.0 // TODO: Maybe add calc from gold nisab in future
	zakahRate   = 0.025 // 2.5%
)

// in zakah package, new function alongside CalculateZakah
func ComputeNetWealth(input types.CalculationRequest, goldPriceUSD, silverPriceUSD, exchangeRate float64) (netWealth float64, nisabThreshold float64) {
	localGoldPricePerGram := goldPriceUSD * exchangeRate
	localSilverPricePerGram := silverPriceUSD * exchangeRate

	goldWealth := input.GoldGrams * localGoldPricePerGram
	silverWealth := input.SilverGrams * localSilverPricePerGram

	totalAssets := goldWealth + silverWealth + input.Cash + input.BusinessAssets
	netWealth = totalAssets - input.Liabilities

	nisabThreshold = silverNisab * localSilverPricePerGram

	return netWealth, nisabThreshold
}

func CalculateZakah(netWealth float64, nisabThreshold float64, currencyCode string) types.CalculationResponse {
	if netWealth >= nisabThreshold {
		zakahDue := netWealth * zakahRate
		message := "Zakah is Applicable"

		// Make Response if Zakah Applicable
		response := types.CalculationResponse{
			TotalAssets:    netWealth,
			NisabThreshold: nisabThreshold,
			ZakahPayable:   zakahDue,
			Currency:       currencyCode, // The values are now in this currency
			Message:        message,
		}
		return response
	} else {
		zakahDue := 0.0
		message := "Zakah is Not Applicable"

		response := types.CalculationResponse{
			TotalAssets:    netWealth,
			NisabThreshold: nisabThreshold,
			ZakahPayable:   zakahDue,
			Currency:       currencyCode,
			Message:        message,
		}

		return response
	}
}

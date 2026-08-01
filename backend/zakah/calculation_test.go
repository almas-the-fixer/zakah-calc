package zakah

import (
	"testing"

	"github.com/almas-the-fixer/zakah-calc/types"
)

func TestCalculateZakah(t *testing.T) {
	tests := []struct {
		name           string
		netWealth      float64
		nisabThreshold float64
		currencyCode   string
		wantZakahDue   float64
		wantMessage    string
	}{
		// When Wealth is ABOVE Nisab
		{
			name:           "Zakah Applicable",
			netWealth:      150000.0,
			nisabThreshold: 120000.0,
			currencyCode:   "INR",
			wantZakahDue:   150000.0 * 0.025,
			wantMessage:    "Zakah is Applicable",
		},
		// When Wealth is EQUAL to Nisab
		{
			name:           "Zakah Applicable",
			netWealth:      150000.0,
			nisabThreshold: 150000.0,
			currencyCode:   "INR",
			wantZakahDue:   150000.0 * 0.025,
			wantMessage:    "Zakah is Applicable",
		},
		// When Wealth is BELOW Nisab
		{
			name:           "No Zakah Applicable",
			netWealth:      100000.0,
			nisabThreshold: 120000.0,
			currencyCode:   "INR",
			wantZakahDue:   0,
			wantMessage:    "Zakah is Not Applicable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateZakah(tt.netWealth, tt.nisabThreshold, tt.currencyCode)
			if got.ZakahPayable != tt.wantZakahDue {
				t.Errorf("Got: %v, Want: %v", got.ZakahPayable, tt.wantZakahDue)
			}
			if got.Message != tt.wantMessage {
				t.Errorf("Got: %v, Want: %v", got.Message, tt.wantMessage)
			}
			if got.Currency != tt.currencyCode {
				t.Errorf("Got: %v, Want: %v", got.Currency, tt.currencyCode)
			}
		})
	}
}

func TestComputeNetWealth(t *testing.T) {
	tests := []struct {
		name           string
		input          types.CalculationRequest
		goldPriceUSD   float64
		silverPriceUSD float64
		exchangeRate   float64
		wantNetWealth  float64
		wantNisab      float64
	}{
		// Case 1: Assets > Liabilities => netWealth must be +ve
		{
			name: "Case 1: Assets > Liabilities",
			input: types.CalculationRequest{
				Currency: "USD",
				GoldGrams: 85.0,
				SilverGrams: 595.0,
				Cash: 100000.0,
				BusinessAssets: 150000.0,
				Liabilities: 50000.0,
			},
			goldPriceUSD: 400.0,
			silverPriceUSD: 420.0,
			exchangeRate: 1.0,

			wantNetWealth: 100000.0 + 150000.0 + (85.0 * 400.0)+(420.0 * 595.0) - 50000.0, // 483,900 ??
			wantNisab: 249900.0, //
		},
		// Case 2: Liabilities > Assets
		{
			name: "Case 2: Liabilities > Assets",
			input: types.CalculationRequest{
				Currency: "USD",
				GoldGrams: 50.0,
				SilverGrams: 100.0,
				Cash: 50000.0,
				BusinessAssets: 25000.0,
				Liabilities: 500000.0,
			},
			goldPriceUSD: 400.0,
			silverPriceUSD: 420.0,
			exchangeRate: 1.0,

			wantNetWealth: (50.0 * 400.0)+(100.0 * 420.0) + 50000.0 + 25000.0 - 500000.0, // -363,000 ??
			wantNisab: 249900.0, // 249900
		},
	}

	for _, tt :=  range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			gotNetWealth, gotNisabThreshold := ComputeNetWealth(tt.input, tt.goldPriceUSD, tt.silverPriceUSD, tt.exchangeRate)

			if gotNetWealth != tt.wantNetWealth {
				t.Errorf("Got: %v Want: %v", gotNetWealth, tt.wantNetWealth)
			}
			if gotNisabThreshold != tt.wantNisab {
				t.Errorf("Got: %v Want: %v", gotNisabThreshold, tt.wantNisab)
			}
		})
	}
}
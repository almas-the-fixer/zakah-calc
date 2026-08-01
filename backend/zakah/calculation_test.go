	package zakah

	import (
		"testing"
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

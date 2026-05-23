// Package tax provides common tax rates and helpers.
package tax

// Rate is a VAT or sales tax rate as a fraction (e.g. 0.081 for 8.1%).
type Rate struct {
	Code        string
	Label       string
	Percent     float64
	Fraction    float64
	CountryCode string
}

// CalculateTax returns tax on amount in minor units.
func (r Rate) CalculateTax(amountMinor int64) int64 {
	if r.Fraction <= 0 {
		return 0
	}
	return int64(float64(amountMinor) * r.Fraction)
}

// Predefined rates (approximate standard rates; verify for production compliance).
var (
	CHStandard = Rate{Code: "CH_STANDARD", Label: "Swiss VAT 8.1%", Percent: 8.1, Fraction: 0.081, CountryCode: "CH"}
	DEReduced  = Rate{Code: "DE_REDUCED", Label: "German VAT 7%", Percent: 7, Fraction: 0.07, CountryCode: "DE"}
	DEStandard = Rate{Code: "DE_STANDARD", Label: "German VAT 19%", Percent: 19, Fraction: 0.19, CountryCode: "DE"}
	UKStandard = Rate{Code: "UK_STANDARD", Label: "UK VAT 20%", Percent: 20, Fraction: 0.20, CountryCode: "GB"}
	USNone     = Rate{Code: "US_NONE", Label: "No tax", Percent: 0, Fraction: 0, CountryCode: "US"}
)

// ByCode returns a predefined rate or false.
func ByCode(code string) (Rate, bool) {
	for _, r := range []Rate{CHStandard, DEReduced, DEStandard, UKStandard, USNone} {
		if r.Code == code {
			return r, true
		}
	}
	return Rate{}, false
}

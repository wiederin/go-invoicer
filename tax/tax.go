// Package tax provides tax rate calculations and common VAT rates
// for invoice generation across different jurisdictions.
package tax

import "github.com/shopspring/decimal"

// Rate represents a tax rate with a name, percentage, and category.
type Rate struct {
	Name       string          `json:"name"`
	Percentage decimal.Decimal `json:"percentage"`
	Category   Category        `json:"category"`
}

// Category represents the type of tax rate (standard, reduced, zero, exempt).
type Category string

// Tax rate categories.
const (
	CategoryStandard Category = "standard"
	CategoryReduced  Category = "reduced"
	CategoryZero     Category = "zero"
	CategoryExempt   Category = "exempt"
)

// NewRate creates a new tax rate with the given name, percentage, and category.
func NewRate(name string, percentage float64, category Category) Rate {
	return Rate{
		Name:       name,
		Percentage: decimal.NewFromFloat(percentage),
		Category:   category,
	}
}

// Calculate computes the tax amount for a given net amount.
func (r Rate) Calculate(amount decimal.Decimal) decimal.Decimal {
	factor := r.Percentage.Div(decimal.NewFromInt(100))
	return amount.Mul(factor).Round(2)
}

// AddToAmount adds the tax to the given amount and returns the gross amount.
func (r Rate) AddToAmount(amount decimal.Decimal) decimal.Decimal {
	tax := r.Calculate(amount)
	return amount.Add(tax)
}

// ExtractFromGross extracts the tax amount from a gross amount.
func (r Rate) ExtractFromGross(grossAmount decimal.Decimal) decimal.Decimal {
	divisor := decimal.NewFromInt(100).Add(r.Percentage)
	netAmount := grossAmount.Mul(decimal.NewFromInt(100)).Div(divisor)
	return grossAmount.Sub(netAmount).Round(2)
}

// CommonRates contains predefined tax rates for various countries.
var CommonRates = map[string]Rate{
	"CH_STANDARD": NewRate("Swiss Standard VAT", 8.1, CategoryStandard),
	"CH_REDUCED":  NewRate("Swiss Reduced VAT", 2.6, CategoryReduced),
	"CH_HOTEL":    NewRate("Swiss Hotel VAT", 3.8, CategoryReduced),
	"DE_STANDARD": NewRate("German Standard VAT", 19.0, CategoryStandard),
	"DE_REDUCED":  NewRate("German Reduced VAT", 7.0, CategoryReduced),
	"FR_STANDARD": NewRate("French Standard VAT", 20.0, CategoryStandard),
	"FR_REDUCED":  NewRate("French Reduced VAT", 5.5, CategoryReduced),
	"UK_STANDARD": NewRate("UK Standard VAT", 20.0, CategoryStandard),
	"UK_REDUCED":  NewRate("UK Reduced VAT", 5.0, CategoryReduced),
	"US_ZERO":     NewRate("US No Federal VAT", 0.0, CategoryZero),
	"EU_EXEMPT":   NewRate("EU VAT Exempt", 0.0, CategoryExempt),
}

// GetRate returns the tax rate for a given country/rate code.
func GetRate(code string) (Rate, bool) {
	rate, ok := CommonRates[code]
	return rate, ok
}

// Calculator accumulates multiple tax rates and calculates combined taxes.
type Calculator struct {
	rates []Rate
}

// NewCalculator creates a new tax calculator.
func NewCalculator() *Calculator {
	return &Calculator{
		rates: []Rate{},
	}
}

// AddRate adds a tax rate to the calculator.
func (c *Calculator) AddRate(rate Rate) *Calculator {
	c.rates = append(c.rates, rate)
	return c
}

// CalculateTax computes the total tax for all rates on a net amount.
func (c *Calculator) CalculateTax(netAmount decimal.Decimal) decimal.Decimal {
	total := decimal.Zero
	for _, rate := range c.rates {
		total = total.Add(rate.Calculate(netAmount))
	}
	return total.Round(2)
}

// CalculateGross computes the gross amount after adding all taxes.
func (c *Calculator) CalculateGross(netAmount decimal.Decimal) decimal.Decimal {
	tax := c.CalculateTax(netAmount)
	return netAmount.Add(tax)
}

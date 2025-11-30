package tax

import "github.com/shopspring/decimal"

type Rate struct {
	Name       string          `json:"name"`
	Percentage decimal.Decimal `json:"percentage"`
	Category   Category        `json:"category"`
}

type Category string

const (
	CategoryStandard Category = "standard"
	CategoryReduced  Category = "reduced"
	CategoryZero     Category = "zero"
	CategoryExempt   Category = "exempt"
)

func NewRate(name string, percentage float64, category Category) Rate {
	return Rate{
		Name:       name,
		Percentage: decimal.NewFromFloat(percentage),
		Category:   category,
	}
}

func (r Rate) Calculate(amount decimal.Decimal) decimal.Decimal {
	factor := r.Percentage.Div(decimal.NewFromInt(100))
	return amount.Mul(factor).Round(2)
}

func (r Rate) AddToAmount(amount decimal.Decimal) decimal.Decimal {
	tax := r.Calculate(amount)
	return amount.Add(tax)
}

func (r Rate) ExtractFromGross(grossAmount decimal.Decimal) decimal.Decimal {
	divisor := decimal.NewFromInt(100).Add(r.Percentage)
	netAmount := grossAmount.Mul(decimal.NewFromInt(100)).Div(divisor)
	return grossAmount.Sub(netAmount).Round(2)
}

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

func GetRate(code string) (Rate, bool) {
	rate, ok := CommonRates[code]
	return rate, ok
}

type Calculator struct {
	rates []Rate
}

func NewCalculator() *Calculator {
	return &Calculator{
		rates: []Rate{},
	}
}

func (c *Calculator) AddRate(rate Rate) *Calculator {
	c.rates = append(c.rates, rate)
	return c
}

func (c *Calculator) CalculateTax(netAmount decimal.Decimal) decimal.Decimal {
	total := decimal.Zero
	for _, rate := range c.rates {
		total = total.Add(rate.Calculate(netAmount))
	}
	return total.Round(2)
}

func (c *Calculator) CalculateGross(netAmount decimal.Decimal) decimal.Decimal {
	tax := c.CalculateTax(netAmount)
	return netAmount.Add(tax)
}

package invoice

import "github.com/shopspring/decimal"

// LineItem represents a single item or service on an invoice.
type LineItem struct {
	Description string          `json:"description"`
	Quantity    decimal.Decimal `json:"quantity"`
	UnitPrice   Money           `json:"unit_price"`
	TaxRate     decimal.Decimal `json:"tax_rate"`
	Discount    decimal.Decimal `json:"discount,omitempty"`
}

// NewLineItem creates a new line item with the given values.
func NewLineItem(description string, quantity float64, unitPrice Money, taxRate float64) LineItem {
	return LineItem{
		Description: description,
		Quantity:    decimal.NewFromFloat(quantity),
		UnitPrice:   unitPrice,
		TaxRate:     decimal.NewFromFloat(taxRate),
		Discount:    decimal.Zero,
	}
}

// WithDiscount returns a copy of the line item with a discount percentage applied.
func (li LineItem) WithDiscount(discountPercent float64) LineItem {
	li.Discount = decimal.NewFromFloat(discountPercent)
	return li
}

// SubTotal returns the quantity times unit price before any discounts.
func (li LineItem) SubTotal() Money {
	amount := li.UnitPrice.Amount.Mul(li.Quantity)
	return Money{Amount: amount, Currency: li.UnitPrice.Currency}
}

// DiscountAmount returns the discount amount for this line item.
func (li LineItem) DiscountAmount() Money {
	if li.Discount.IsZero() {
		return Money{Amount: decimal.Zero, Currency: li.UnitPrice.Currency}
	}
	subtotal := li.SubTotal()
	discountFactor := li.Discount.Div(decimal.NewFromInt(100))
	return subtotal.Mul(discountFactor)
}

// NetAmount returns the amount after discount but before tax.
func (li LineItem) NetAmount() Money {
	subtotal := li.SubTotal()
	discount := li.DiscountAmount()
	net, _ := subtotal.Sub(discount)
	return net
}

// TaxAmount returns the tax amount for this line item.
func (li LineItem) TaxAmount() Money {
	net := li.NetAmount()
	taxFactor := li.TaxRate.Div(decimal.NewFromInt(100))
	return net.Mul(taxFactor)
}

// GrossAmount returns the total amount including tax.
func (li LineItem) GrossAmount() Money {
	net := li.NetAmount()
	tax := li.TaxAmount()
	gross, _ := net.Add(tax)
	return gross
}

// Validate checks that the line item has all required fields.
func (li LineItem) Validate() error {
	if li.Description == "" {
		return ErrMissingDescription
	}
	if li.Quantity.IsNegative() || li.Quantity.IsZero() {
		return ErrInvalidQuantity
	}
	if li.UnitPrice.IsNegative() {
		return ErrInvalidUnitPrice
	}
	if li.TaxRate.IsNegative() {
		return ErrInvalidTaxRate
	}
	return nil
}

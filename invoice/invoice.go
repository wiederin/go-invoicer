// Package invoice defines the core invoice domain model.
package invoice

import (
	"fmt"
	"sort"
	"time"
)

// Money represents a monetary amount in minor units (e.g. cents).
type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// NewMoney creates money in minor units with ISO 4217 currency code.
func NewMoney(amountMinor int64, currency string) Money {
	return Money{Amount: amountMinor, Currency: currency}
}

func (m Money) String() string {
	if m.Currency == "" {
		return fmt.Sprintf("%d", m.Amount)
	}
	return fmt.Sprintf("%s %d", m.Currency, m.Amount)
}

// Party is a seller or buyer on an invoice.
type Party struct {
	Name       string  `json:"name"`
	Email      string  `json:"email,omitempty"`
	Phone      string  `json:"phone,omitempty"`
	Department string  `json:"department,omitempty"`
	Address    Address `json:"address,omitempty"`
	VATID      string  `json:"vat_id,omitempty"`
}

// Address is a postal address.
type Address struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"` // ISO 3166-1 alpha-2
}

// FormatSingleLine returns a one-line postal address.
func (a Address) FormatSingleLine() string {
	parts := []string{a.Line1, a.Line2, a.PostalCode, a.City, a.Region, a.Country}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return joinNonEmpty(out, ", ")
}

func joinNonEmpty(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	s := parts[0]
	for i := 1; i < len(parts); i++ {
		s += sep + parts[i]
	}
	return s
}

// LineItem is a single billable row.
type LineItem struct {
	Description      string            `json:"description"`
	Quantity         float64           `json:"quantity"`
	UnitPrice        Money             `json:"unit_price"`
	TaxRate          float64           `json:"tax_rate"` // fraction, e.g. 0.081 for 8.1% VAT
	DiscountPct      float64           `json:"discount_pct,omitempty"` // percentage 0-100
	CustomFields     map[string]string `json:"custom_fields,omitempty"`
	SupplyDateStart  string            `json:"supply_date_start,omitempty"` // YYYY-MM-DD
	SupplyDateEnd    string            `json:"supply_date_end,omitempty"`   // YYYY-MM-DD
}

// NewLineItem creates a line item with quantity, unit price in minor units, and tax fraction.
func NewLineItem(description string, quantity float64, unitPrice Money, taxRate float64) LineItem {
	return LineItem{
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		TaxRate:     taxRate,
	}
}

// LineTotal returns quantity * unit price in minor units, after any discount.
func (l LineItem) LineTotal() int64 {
	gross := float64(l.UnitPrice.Amount) * l.Quantity
	if l.DiscountPct > 0 && l.DiscountPct <= 100 {
		gross *= 1 - l.DiscountPct/100
	}
	return int64(gross)
}

// TaxAmount returns tax for this line in minor units.
func (l LineItem) TaxAmount() int64 {
	return int64(float64(l.LineTotal()) * l.TaxRate)
}

// GrossTotal returns line total including tax.
func (l LineItem) GrossTotal() int64 {
	return l.LineTotal() + l.TaxAmount()
}

// TaxLine summarizes tax at one rate.
type TaxLine struct {
	Rate       float64
	NetAmount  int64
	TaxAmount  int64
	Currency   string
}

// CustomField is an invoice-level key/value metadata entry displayed on the PDF.
type CustomField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Invoice is the root aggregate for generation and rendering.
type Invoice struct {
	Number         string            `json:"number"`
	Kind           DocumentKind      `json:"kind,omitempty"`
	RelatedNumber  string            `json:"related_number,omitempty"` // source invoice for credit notes
	IssuedAt       time.Time         `json:"issued_at"`
	DueAt          time.Time         `json:"due_at,omitempty"`
	Seller         Party             `json:"seller"`
	Buyer          Party             `json:"buyer"`
	LineItems      []LineItem        `json:"line_items"`
	Notes          string            `json:"notes,omitempty"`
	Footer         string            `json:"footer,omitempty"`
	TaxID          string            `json:"tax_id,omitempty"`
	CustomFields   []CustomField     `json:"custom_fields,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// New creates an invoice with the given number and issue date set to now (UTC).
func New(number string) *Invoice {
	return &Invoice{
		Number:   number,
		IssuedAt: time.Now().UTC(),
		Metadata: make(map[string]string),
	}
}

func (i *Invoice) SetSeller(p Party) { i.Seller = p }
func (i *Invoice) SetBuyer(p Party)   { i.Buyer = p }

func (i *Invoice) AddLine(item LineItem) {
	i.LineItems = append(i.LineItems, item)
}

// Currency returns the invoice currency from line items.
func (i *Invoice) Currency() string {
	for _, l := range i.LineItems {
		if l.UnitPrice.Currency != "" {
			return l.UnitPrice.Currency
		}
	}
	return ""
}

// Subtotal returns the sum of line totals before tax.
func (i *Invoice) Subtotal() Money {
	var sum int64
	for _, l := range i.LineItems {
		sum += l.LineTotal()
	}
	return Money{Amount: sum, Currency: i.Currency()}
}

// TotalTax returns total tax across all lines.
func (i *Invoice) TotalTax() Money {
	var sum int64
	for _, l := range i.LineItems {
		sum += l.TaxAmount()
	}
	return Money{Amount: sum, Currency: i.Currency()}
}

// Total returns subtotal plus tax.
func (i *Invoice) Total() Money {
	sub := i.Subtotal()
	tax := i.TotalTax()
	return Money{Amount: sub.Amount + tax.Amount, Currency: sub.Currency}
}

// TaxBreakdown groups tax amounts by rate.
func (i *Invoice) TaxBreakdown() []TaxLine {
	byRate := make(map[float64]*TaxLine)
	for _, l := range i.LineItems {
		if l.TaxRate <= 0 {
			continue
		}
		line, ok := byRate[l.TaxRate]
		if !ok {
			line = &TaxLine{Rate: l.TaxRate, Currency: l.UnitPrice.Currency}
			byRate[l.TaxRate] = line
		}
		line.NetAmount += l.LineTotal()
		line.TaxAmount += l.TaxAmount()
	}
	rates := make([]float64, 0, len(byRate))
	for r := range byRate {
		rates = append(rates, r)
	}
	sort.Float64s(rates)
	out := make([]TaxLine, 0, len(rates))
	for _, r := range rates {
		out = append(out, *byRate[r])
	}
	return out
}

// Validate checks required fields for rendering.
func (i *Invoice) Validate() error {
	if i.Number == "" {
		return ErrMissingNumber
	}
	if i.Seller.Name == "" {
		return ErrMissingSeller
	}
	if i.Buyer.Name == "" {
		return ErrMissingBuyer
	}
	if len(i.LineItems) == 0 {
		return ErrMissingLineItems
	}
	if !i.DueAt.IsZero() && i.DueAt.Before(i.IssuedAt) {
		return ErrDueBeforeIssue
	}
	if i.NormalizedKind() == DocumentCreditNote && i.RelatedNumber == "" {
		return ErrMissingRelated
	}
	cur := ""
	for _, l := range i.LineItems {
		if l.UnitPrice.Currency == "" {
			return ErrMissingCurrency
		}
		if cur == "" {
			cur = l.UnitPrice.Currency
		} else if l.UnitPrice.Currency != cur {
			return ErrCurrencyMismatch
		}
	}
	return nil
}

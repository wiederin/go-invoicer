package invoice

import (
	"time"

	"github.com/shopspring/decimal"
)

// Status represents the current state of an invoice.
type Status string

// Invoice status constants.
const (
	StatusDraft     Status = "draft"
	StatusIssued    Status = "issued"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusOverdue   Status = "overdue"
)

// Invoice represents a complete invoice document.
type Invoice struct {
	Number      string     `json:"number"`
	IssueDate   time.Time  `json:"issue_date"`
	DueDate     time.Time  `json:"due_date"`
	Currency    string     `json:"currency"`
	CountryCode string     `json:"country_code"`
	Supplier    Party      `json:"supplier"`
	Customer    Party      `json:"customer"`
	LineItems   []LineItem `json:"line_items"`
	Notes       string     `json:"notes,omitempty"`
	Terms       string     `json:"terms,omitempty"`
	Status      Status     `json:"status"`
	Metadata    Metadata   `json:"metadata,omitempty"`
}

// Metadata stores arbitrary key-value pairs for an invoice.
type Metadata map[string]any

// Builder provides a fluent interface for constructing invoices.
type Builder struct {
	inv Invoice
}

// New creates a new invoice builder with default values.
func New() *Builder {
	return &Builder{
		inv: Invoice{
			Status:   StatusDraft,
			Currency: "USD",
			Metadata: make(Metadata),
		},
	}
}

// Number sets the invoice number.
func (b *Builder) Number(number string) *Builder {
	b.inv.Number = number
	return b
}

// IssueDate sets the invoice issue date.
func (b *Builder) IssueDate(date time.Time) *Builder {
	b.inv.IssueDate = date
	return b
}

// DueDate sets the invoice due date.
func (b *Builder) DueDate(date time.Time) *Builder {
	b.inv.DueDate = date
	return b
}

// Currency sets the invoice currency code.
func (b *Builder) Currency(currency string) *Builder {
	b.inv.Currency = currency
	return b
}

// CountryCode sets the invoice country code for tax purposes.
func (b *Builder) CountryCode(code string) *Builder {
	b.inv.CountryCode = code
	return b
}

// Supplier sets the invoice supplier (seller).
func (b *Builder) Supplier(supplier Party) *Builder {
	b.inv.Supplier = supplier
	return b
}

// Customer sets the invoice customer (buyer).
func (b *Builder) Customer(customer Party) *Builder {
	b.inv.Customer = customer
	return b
}

// AddItem adds a line item to the invoice.
func (b *Builder) AddItem(item LineItem) *Builder {
	b.inv.LineItems = append(b.inv.LineItems, item)
	return b
}

// Notes sets additional notes on the invoice.
func (b *Builder) Notes(notes string) *Builder {
	b.inv.Notes = notes
	return b
}

// Terms sets the payment terms on the invoice.
func (b *Builder) Terms(terms string) *Builder {
	b.inv.Terms = terms
	return b
}

// Status sets the invoice status.
func (b *Builder) Status(status Status) *Builder {
	b.inv.Status = status
	return b
}

// SetMetadata sets a metadata key-value pair.
func (b *Builder) SetMetadata(key string, value any) *Builder {
	b.inv.Metadata[key] = value
	return b
}

// Build validates and returns the constructed invoice.
func (b *Builder) Build() (*Invoice, error) {
	inv := b.inv
	if err := inv.Validate(); err != nil {
		return nil, err
	}
	inv.RecalculateTotals()
	return &inv, nil
}

// Validate checks that all required fields are present and valid.
func (inv *Invoice) Validate() error {
	if inv.Number == "" {
		return ErrMissingInvoiceNumber
	}
	if inv.IssueDate.IsZero() {
		return ErrMissingIssueDate
	}
	if inv.DueDate.IsZero() {
		return ErrMissingDueDate
	}
	if inv.DueDate.Before(inv.IssueDate) {
		return ErrDueDateBeforeIssue
	}
	if err := inv.Supplier.Validate(); err != nil {
		return ErrMissingSupplier
	}
	if err := inv.Customer.Validate(); err != nil {
		return ErrMissingCustomer
	}
	if len(inv.LineItems) == 0 {
		return ErrNoLineItems
	}
	for _, item := range inv.LineItems {
		if err := item.Validate(); err != nil {
			return err
		}
		if item.UnitPrice.Currency != inv.Currency {
			return ErrCurrencyMismatch
		}
	}
	return nil
}

// RecalculateTotals recalculates all invoice totals from line items.
func (inv *Invoice) RecalculateTotals() {
}

// SubTotal returns the sum of all line item subtotals before discounts.
func (inv *Invoice) SubTotal() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		subtotal := item.SubTotal()
		total, _ = total.Add(subtotal)
	}
	return total
}

// TotalDiscount returns the sum of all line item discounts.
func (inv *Invoice) TotalDiscount() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		discount := item.DiscountAmount()
		total, _ = total.Add(discount)
	}
	return total
}

// TotalNet returns the total after discounts but before taxes.
func (inv *Invoice) TotalNet() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		net := item.NetAmount()
		total, _ = total.Add(net)
	}
	return total
}

// TotalTax returns the sum of all line item taxes.
func (inv *Invoice) TotalTax() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		tax := item.TaxAmount()
		total, _ = total.Add(tax)
	}
	return total
}

// TotalGross returns the final invoice total including taxes.
func (inv *Invoice) TotalGross() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		gross := item.GrossAmount()
		total, _ = total.Add(gross)
	}
	return total
}

// TaxBreakdown returns a map of tax rates to their total amounts.
func (inv *Invoice) TaxBreakdown() map[string]Money {
	breakdown := make(map[string]Money)
	for _, item := range inv.LineItems {
		rate := item.TaxRate.String()
		tax := item.TaxAmount()
		if existing, ok := breakdown[rate]; ok {
			breakdown[rate], _ = existing.Add(tax)
		} else {
			breakdown[rate] = tax
		}
	}
	return breakdown
}

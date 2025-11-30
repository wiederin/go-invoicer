package invoice

import (
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusIssued    Status = "issued"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusOverdue   Status = "overdue"
)

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

type Metadata map[string]any

type Builder struct {
	inv Invoice
}

func New() *Builder {
	return &Builder{
		inv: Invoice{
			Status:   StatusDraft,
			Currency: "USD",
			Metadata: make(Metadata),
		},
	}
}

func (b *Builder) Number(number string) *Builder {
	b.inv.Number = number
	return b
}

func (b *Builder) IssueDate(date time.Time) *Builder {
	b.inv.IssueDate = date
	return b
}

func (b *Builder) DueDate(date time.Time) *Builder {
	b.inv.DueDate = date
	return b
}

func (b *Builder) Currency(currency string) *Builder {
	b.inv.Currency = currency
	return b
}

func (b *Builder) CountryCode(code string) *Builder {
	b.inv.CountryCode = code
	return b
}

func (b *Builder) Supplier(supplier Party) *Builder {
	b.inv.Supplier = supplier
	return b
}

func (b *Builder) Customer(customer Party) *Builder {
	b.inv.Customer = customer
	return b
}

func (b *Builder) AddItem(item LineItem) *Builder {
	b.inv.LineItems = append(b.inv.LineItems, item)
	return b
}

func (b *Builder) Notes(notes string) *Builder {
	b.inv.Notes = notes
	return b
}

func (b *Builder) Terms(terms string) *Builder {
	b.inv.Terms = terms
	return b
}

func (b *Builder) Status(status Status) *Builder {
	b.inv.Status = status
	return b
}

func (b *Builder) SetMetadata(key string, value any) *Builder {
	b.inv.Metadata[key] = value
	return b
}

func (b *Builder) Build() (*Invoice, error) {
	inv := b.inv
	if err := inv.Validate(); err != nil {
		return nil, err
	}
	inv.RecalculateTotals()
	return &inv, nil
}

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

func (inv *Invoice) RecalculateTotals() {
}

func (inv *Invoice) SubTotal() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		subtotal := item.SubTotal()
		total, _ = total.Add(subtotal)
	}
	return total
}

func (inv *Invoice) TotalDiscount() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		discount := item.DiscountAmount()
		total, _ = total.Add(discount)
	}
	return total
}

func (inv *Invoice) TotalNet() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		net := item.NetAmount()
		total, _ = total.Add(net)
	}
	return total
}

func (inv *Invoice) TotalTax() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		tax := item.TaxAmount()
		total, _ = total.Add(tax)
	}
	return total
}

func (inv *Invoice) TotalGross() Money {
	total := Money{Amount: decimal.Zero, Currency: inv.Currency}
	for _, item := range inv.LineItems {
		gross := item.GrossAmount()
		total, _ = total.Add(gross)
	}
	return total
}

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

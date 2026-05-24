// Package stripe maps Stripe invoice data to the go-invoicer domain model.
//
// Use the API client, webhook HTTP handler, or unmarshalled Stripe JSON.
package stripe

import (
	"strings"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
)

// StripeInvoiceIDMeta is the metadata key linking a domain invoice to Stripe.
const StripeInvoiceIDMeta = "stripe_invoice_id"

// Invoice is a subset of Stripe invoice fields needed for PDF rendering.
type Invoice struct {
	ID           string
	Number       string
	Currency     string
	Created      time.Time
	DueDate      *time.Time
	CustomerName string
	CustomerEmail string
	Lines        []Line
	Metadata     map[string]string
}

// Line is a Stripe invoice line item.
type Line struct {
	Description string
	Quantity    float64
	AmountMinor int64 // amount in minor units (Stripe uses cents)
	TaxRate     float64 // fraction, e.g. 0.081
}

// Supplier is the seller party on mirrored invoices.
type Supplier struct {
	Name    string
	Email   string
	VATID   string
	Address invoice.Address
}

// ImportOptions configures Stripe → domain mapping.
type ImportOptions struct {
	Supplier Supplier
}

// ToDomain maps a Stripe invoice to invoice.Invoice (not yet validated).
func (s Invoice) ToDomain(opts ImportOptions) *invoice.Invoice {
	num := s.Number
	if num == "" {
		num = s.ID
	}
	inv := invoice.New(num)
	inv.IssuedAt = s.Created.UTC()
	if s.DueDate != nil {
		inv.DueAt = s.DueDate.UTC()
	}
	inv.Seller = invoice.Party{
		Name:    opts.Supplier.Name,
		Email:   opts.Supplier.Email,
		VATID:   opts.Supplier.VATID,
		Address: opts.Supplier.Address,
	}
	inv.Buyer = invoice.Party{
		Name:  s.CustomerName,
		Email: s.CustomerEmail,
	}
	for _, line := range s.Lines {
		inv.AddLine(invoice.NewLineItem(
			line.Description,
			line.Quantity,
			invoice.NewMoney(line.AmountMinor, strings.ToUpper(s.Currency)),
			line.TaxRate,
		))
	}
	for k, v := range s.Metadata {
		inv.Metadata[k] = v
	}
	return inv
}

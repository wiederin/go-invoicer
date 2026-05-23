package stripe

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v81"
	stripeapi "github.com/stripe/stripe-go/v81/invoice"
)

// Client fetches invoices from the Stripe API.
type Client struct {
	apiKey string
}

// NewClient creates a Stripe API client.
func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// FetchInvoice retrieves an invoice by ID (e.g. in_xxx) and maps it to the import type.
func (c *Client) FetchInvoice(ctx context.Context, invoiceID string) (Invoice, error) {
	if c.apiKey == "" {
		return Invoice{}, fmt.Errorf("stripe: API key is required")
	}
	stripe.Key = c.apiKey

	params := &stripe.InvoiceParams{}
	params.Context = ctx
	inv, err := stripeapi.Get(invoiceID, params)
	if err != nil {
		return Invoice{}, fmt.Errorf("stripe: fetch invoice: %w", err)
	}
	return mapAPIInvoice(inv), nil
}

func mapAPIInvoice(inv *stripe.Invoice) Invoice {
	var due *time.Time
	if inv.DueDate > 0 {
		t := time.Unix(inv.DueDate, 0).UTC()
		due = &t
	}
	lines := make([]Line, 0)
	if inv.Lines != nil {
		for _, l := range inv.Lines.Data {
			amount := l.Amount
			if amount == 0 && l.UnitAmountExcludingTax != 0 {
				amount = int64(l.UnitAmountExcludingTax)
			}
			qty := float64(1)
			if l.Quantity > 0 {
				qty = float64(l.Quantity)
			}
			desc := l.Description
			lines = append(lines, Line{
				Description: desc,
				Quantity:    qty,
				AmountMinor: amount,
				TaxRate:     0, // expand tax rates via Stripe API for precise mapping
			})
		}
	}
	meta := make(map[string]string)
	for k, v := range inv.Metadata {
		meta[k] = v
	}
	return Invoice{
		ID:            inv.ID,
		Number:        inv.Number,
		Currency:      strings.ToUpper(string(inv.Currency)),
		Created:       time.Unix(inv.Created, 0).UTC(),
		DueDate:       due,
		CustomerName:  inv.CustomerName,
		CustomerEmail: inv.CustomerEmail,
		Lines:         lines,
		Metadata:      meta,
	}
}

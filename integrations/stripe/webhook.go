package stripe

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Webhook event types used by go-invoicer.
const (
	EventInvoiceFinalized = "invoice.finalized"
	EventInvoicePaid      = "invoice.paid"
)

// WebhookEvent is a minimal Stripe event envelope.
type WebhookEvent struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Created int64           `json:"created"`
	Data    WebhookEventData `json:"data"`
}

// WebhookEventData holds the event object.
type WebhookEventData struct {
	Object json.RawMessage `json:"object"`
}

// APIInvoice mirrors key Stripe invoice JSON fields (API version agnostic subset).
type APIInvoice struct {
	ID               string `json:"id"`
	Number           string `json:"number"`
	Currency         string `json:"currency"`
	Created          int64  `json:"created"`
	DueDate          int64  `json:"due_date"`
	CustomerName     string `json:"customer_name"`
	CustomerEmail    string `json:"customer_email"`
	AmountDue        int64  `json:"amount_due"`
	HostedInvoiceURL string `json:"hosted_invoice_url"`
	Lines            struct {
		Data []APIInvoiceLine `json:"data"`
	} `json:"lines"`
	Metadata map[string]string `json:"metadata"`
}

// APIInvoiceLine is a Stripe invoice line item.
type APIInvoiceLine struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Amount      int64   `json:"amount"`
}

// ParseWebhookEvent unmarshals a Stripe webhook JSON body.
func ParseWebhookEvent(body []byte) (*WebhookEvent, error) {
	var ev WebhookEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("stripe: parse webhook: %w", err)
	}
	return &ev, nil
}

// InvoiceFromEvent extracts an APIInvoice from supported invoice events.
func InvoiceFromEvent(ev *WebhookEvent) (*APIInvoice, error) {
	switch ev.Type {
	case EventInvoiceFinalized, EventInvoicePaid:
	default:
		return nil, fmt.Errorf("stripe: unsupported event type %q", ev.Type)
	}
	var inv APIInvoice
	if err := json.Unmarshal(ev.Data.Object, &inv); err != nil {
		return nil, fmt.Errorf("stripe: parse invoice object: %w", err)
	}
	return &inv, nil
}

// ToDomainInvoice maps a Stripe API invoice to the domain import type.
func (a *APIInvoice) ToDomainInvoice() Invoice {
	var due *time.Time
	if a.DueDate > 0 {
		t := time.Unix(a.DueDate, 0).UTC()
		due = &t
	}
	lines := make([]Line, 0, len(a.Lines.Data))
	for _, l := range a.Lines.Data {
		lines = append(lines, Line{
			Description: l.Description,
			Quantity:    l.Quantity,
			AmountMinor: l.Amount,
		})
	}
	return Invoice{
		ID:            a.ID,
		Number:        a.Number,
		Currency:      strings.ToUpper(a.Currency),
		Created:       time.Unix(a.Created, 0).UTC(),
		DueDate:       due,
		CustomerName:  a.CustomerName,
		CustomerEmail: a.CustomerEmail,
		Lines:         lines,
		Metadata:      a.Metadata,
	}
}

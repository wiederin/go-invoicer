package stripe

import (
	"context"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v81/webhook"

	domain "github.com/wiederin/go-invoicer/invoice"
)

// HandlerOptions configures the Stripe webhook HTTP handler.
type HandlerOptions struct {
	WebhookSecret string
	Supplier      Supplier
	OnInvoice     func(*domain.Invoice) error
}

// WebhookHandler returns an http.Handler for Stripe invoice events.
func WebhookHandler(opts HandlerOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "read body", http.StatusBadRequest)
			return
		}

		var ev WebhookEvent
		if opts.WebhookSecret != "" {
			headers := r.Header.Get("Stripe-Signature")
			e, err := webhook.ConstructEvent(body, headers, opts.WebhookSecret)
			if err != nil {
				http.Error(w, "invalid signature", http.StatusBadRequest)
				return
			}
			ev = WebhookEvent{
				ID:   e.ID,
				Type: string(e.Type),
				Data: WebhookEventData{Object: e.Data.Raw},
			}
		} else {
			evPtr, err := ParseWebhookEvent(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ev = *evPtr
		}

		apiInv, err := InvoiceFromEvent(&ev)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		inv := apiInv.ToDomainInvoice().ToDomain(ImportOptions{Supplier: opts.Supplier})
		if inv.Metadata == nil {
			inv.Metadata = map[string]string{}
		}
		inv.Metadata["stripe_invoice_id"] = apiInv.ID
		if err := inv.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if opts.OnInvoice != nil {
			if err := opts.OnInvoice(inv); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"received":true}`))
	}
}

// SyncAndRender fetches a Stripe invoice and returns the domain model.
func (c *Client) SyncInvoice(ctx context.Context, invoiceID string, opts ImportOptions) (*domain.Invoice, error) {
	si, err := c.FetchInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	inv := si.ToDomain(opts)
	if err := inv.Validate(); err != nil {
		return nil, err
	}
	return inv, nil
}

package stripe_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wiederin/go-invoicer/integrations/stripe"
	domain "github.com/wiederin/go-invoicer/invoice"
)

func TestWebhookHandler_invoicePaid(t *testing.T) {
	var got *domain.Invoice
	h := stripe.WebhookHandler(stripe.HandlerOptions{
		Supplier: stripe.Supplier{Name: "Acme"},
		OnInvoice: func(inv *domain.Invoice) error {
			got = inv
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", bytes.NewBufferString(sampleWebhook))
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
	if got == nil || got.Number != "INV-9" {
		t.Fatalf("invoice %+v", got)
	}
}

func TestWebhookHandler_unsupportedEvent(t *testing.T) {
	h := stripe.WebhookHandler(stripe.HandlerOptions{
		Supplier: stripe.Supplier{Name: "Acme"},
	})
	body := `{"id":"evt","type":"customer.created","data":{"object":{}}}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestWebhookHandler_methodNotAllowed(t *testing.T) {
	h := stripe.WebhookHandler(stripe.HandlerOptions{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status %d", rec.Code)
	}
}

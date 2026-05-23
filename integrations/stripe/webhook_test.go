package stripe_test

import (
	"testing"

	"github.com/wiederin/go-invoicer/integrations/stripe"
)

const sampleWebhook = `{
  "id": "evt_1",
  "type": "invoice.paid",
  "created": 1714500000,
  "data": {
    "object": {
      "id": "in_1",
      "number": "INV-9",
      "currency": "chf",
      "created": 1714400000,
      "customer_name": "Test Co",
      "customer_email": "a@b.com",
      "lines": {
        "data": [
          {"description": "Plan", "quantity": 1, "amount": 5000}
        ]
      }
    }
  }
}`

func TestParseWebhookInvoicePaid(t *testing.T) {
	ev, err := stripe.ParseWebhookEvent([]byte(sampleWebhook))
	if err != nil {
		t.Fatal(err)
	}
	apiInv, err := stripe.InvoiceFromEvent(ev)
	if err != nil {
		t.Fatal(err)
	}
	si := apiInv.ToDomainInvoice()
	domain := si.ToDomain(stripe.ImportOptions{Supplier: stripe.Supplier{Name: "Seller"}})
	if domain.Number != "INV-9" {
		t.Fatalf("number %s", domain.Number)
	}
}

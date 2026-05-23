package stripe_test

import (
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/integrations/stripe"
	"github.com/wiederin/go-invoicer/tax"
)

func TestToDomain(t *testing.T) {
	due := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	si := stripe.Invoice{
		ID:            "in_123",
		Number:        "INV-STRIPE-1",
		Currency:      "CHF",
		Created:       time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		DueDate:       &due,
		CustomerName:  "Stripe Customer",
		CustomerEmail: "c@example.com",
		Lines: []stripe.Line{{
			Description: "SaaS plan",
			Quantity:    1,
			AmountMinor: 9900,
			TaxRate:     tax.CHStandard.Fraction,
		}},
	}

	domain := si.ToDomain(stripe.ImportOptions{Supplier: stripe.Supplier{Name: "My SaaS Co"}})
	if err := domain.Validate(); err != nil {
		t.Fatal(err)
	}
	if domain.Number != "INV-STRIPE-1" {
		t.Fatalf("number: %s", domain.Number)
	}
	if domain.Buyer.Name != "Stripe Customer" {
		t.Fatalf("buyer: %s", domain.Buyer.Name)
	}
}

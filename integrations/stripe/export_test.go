package stripe_test

import (
	"context"
	"testing"

	"github.com/wiederin/go-invoicer/integrations/stripe"
	"github.com/wiederin/go-invoicer/invoice"
)

func TestPushInvoice_requiresAPIKey(t *testing.T) {
	c := stripe.NewClient("")
	inv := invoice.New("INV-1")
	inv.Seller = invoice.Party{Name: "Seller"}
	inv.Buyer = invoice.Party{Name: "Buyer", Email: "buyer@example.com"}
	inv.AddLine(invoice.NewLineItem("Item", 1, invoice.NewMoney(1000, "CHF"), 0))

	_, err := c.PushInvoice(context.Background(), inv, stripe.ExportOptions{})
	if err == nil {
		t.Fatal("expected error without API key")
	}
}

func TestPushInvoice_requiresCustomer(t *testing.T) {
	c := stripe.NewClient("sk_test_dummy")
	inv := invoice.New("INV-2")
	inv.Seller = invoice.Party{Name: "Seller"}
	inv.Buyer = invoice.Party{Name: "Buyer"}
	inv.AddLine(invoice.NewLineItem("Item", 1, invoice.NewMoney(1000, "CHF"), 0))

	_, err := c.PushInvoice(context.Background(), inv, stripe.ExportOptions{})
	if err == nil {
		t.Fatal("expected error without buyer email or customer id")
	}
}

func TestPushInvoice_nilInvoice(t *testing.T) {
	c := stripe.NewClient("sk_test_dummy")
	_, err := c.PushInvoice(context.Background(), nil, stripe.ExportOptions{})
	if err == nil {
		t.Fatal("expected error for nil invoice")
	}
}

func TestLinkBrandedPDF_requiresAPIKey(t *testing.T) {
	c := stripe.NewClient("")
	err := c.LinkBrandedPDF(context.Background(), "in_123", "https://example.com/pdf", "job-1", "file-1")
	if err == nil {
		t.Fatal("expected error without API key")
	}
}

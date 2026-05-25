package quickbooks_test

import (
	"context"
	"testing"

	"github.com/wiederin/go-invoicer/integrations/quickbooks"
	"github.com/wiederin/go-invoicer/invoice"
)

func TestPushInvoice_requiresTokens(t *testing.T) {
	c := &quickbooks.Client{}
	inv := sampleInvoice(t)
	_, err := c.PushInvoice(context.Background(), inv, quickbooks.ExportOptions{})
	if err == nil {
		t.Fatal("expected error without tokens")
	}
}

func TestPushInvoice_nilInvoice(t *testing.T) {
	c := &quickbooks.Client{AccessToken: "tok", RealmID: "realm"}
	_, err := c.PushInvoice(context.Background(), nil, quickbooks.ExportOptions{})
	if err == nil {
		t.Fatal("expected error for nil invoice")
	}
}

func TestPushInvoice_requiresBuyer(t *testing.T) {
	c := &quickbooks.Client{AccessToken: "tok", RealmID: "realm"}
	inv := invoice.New("INV-1")
	inv.Seller = invoice.Party{Name: "Seller"}
	inv.AddLine(invoice.NewLineItem("Item", 1, invoice.NewMoney(1000, "USD"), 0))
	_, err := c.PushInvoice(context.Background(), inv, quickbooks.ExportOptions{})
	if err == nil {
		t.Fatal("expected error without buyer")
	}
}

func sampleInvoice(t *testing.T) *invoice.Invoice {
	t.Helper()
	inv := invoice.New("INV-1")
	inv.Seller = invoice.Party{Name: "Seller"}
	inv.Buyer = invoice.Party{Name: "Buyer", Email: "buyer@example.com"}
	inv.AddLine(invoice.NewLineItem("Item", 1, invoice.NewMoney(1000, "USD"), 0))
	return inv
}

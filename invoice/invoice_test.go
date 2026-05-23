package invoice_test

import (
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/tax"
)

func TestInvoiceTotal(t *testing.T) {
	inv := invoice.New("INV-001")
	inv.SetSeller(invoice.Party{Name: "Seller"})
	inv.SetBuyer(invoice.Party{Name: "Buyer"})
	inv.AddLine(invoice.LineItem{
		Description: "Item",
		Quantity:    2,
		UnitPrice:   invoice.Money{Amount: 1000, Currency: "CHF"},
		TaxRate:     0.1,
	})

	if err := inv.Validate(); err != nil {
		t.Fatal(err)
	}

	sub := inv.Subtotal()
	if sub.Amount != 2000 {
		t.Fatalf("subtotal: got %d want 2000", sub.Amount)
	}

	total := inv.Total()
	if total.Amount != 2200 {
		t.Fatalf("total: got %d want 2200", total.Amount)
	}
}

func TestBuilder(t *testing.T) {
	issued := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	due := issued.AddDate(0, 0, 30)

	inv, err := invoice.NewBuilder().
		Number("INV-2026-001").
		IssuedAt(issued).
		DueAt(due).
		Seller(invoice.Party{Name: "Acme GmbH", Address: invoice.Address{City: "Zürich", Country: "CH"}}).
		Buyer(invoice.Party{Name: "Client AG"}).
		AddLine(invoice.NewLineItem(
			"Consulting",
			10,
			invoice.NewMoney(15000, "CHF"),
			tax.CHStandard.Fraction,
		)).
		Notes("Thank you for your business.").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	if inv.Total().Amount <= inv.Subtotal().Amount {
		t.Fatal("expected tax to increase total")
	}
	if len(inv.TaxBreakdown()) != 1 {
		t.Fatalf("tax breakdown: got %d lines", len(inv.TaxBreakdown()))
	}
}

func TestValidateErrors(t *testing.T) {
	_, err := invoice.NewBuilder().Build()
	if err != invoice.ErrMissingNumber {
		t.Fatalf("got %v want ErrMissingNumber", err)
	}

	_, err = invoice.NewBuilder().
		Number("X").
		IssuedAt(time.Now()).
		DueAt(time.Now().AddDate(0, 0, -1)).
		Seller(invoice.Party{Name: "S"}).
		Buyer(invoice.Party{Name: "B"}).
		AddLine(invoice.NewLineItem("x", 1, invoice.NewMoney(100, "CHF"), 0)).
		Build()
	if err != invoice.ErrDueBeforeIssue {
		t.Fatalf("got %v want ErrDueBeforeIssue", err)
	}
}

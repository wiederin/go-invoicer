package settlement_test

import (
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/settlement"
)

func TestBuildSettlement(t *testing.T) {
	from := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	inv1, _ := invoice.NewBuilder().
		Number("INV-1").
		IssuedAt(from.AddDate(0, 0, 2)).
		Seller(invoice.Party{Name: "Seller"}).
		Buyer(invoice.Party{Name: "A"}).
		AddLine(invoice.NewLineItem("X", 1, invoice.NewMoney(1000, "CHF"), 0)).
		Build()
	inv2, _ := invoice.NewBuilder().
		Number("INV-2").
		IssuedAt(from.AddDate(0, 0, 10)).
		Seller(invoice.Party{Name: "Seller"}).
		Buyer(invoice.Party{Name: "B"}).
		AddLine(invoice.NewLineItem("Y", 1, invoice.NewMoney(2500, "CHF"), 0)).
		Build()

	r, err := settlement.Build("Acme", from, to, []settlement.InvoiceInput{
		{Number: "INV-1", Status: "paid", Invoice: *inv1},
		{Number: "INV-2", Status: "draft", Invoice: *inv2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if r.InvoiceCount != 2 || r.TotalMinor != 3500 || r.Currency != "CHF" {
		t.Fatalf("got count=%d total=%d cur=%q", r.InvoiceCount, r.TotalMinor, r.Currency)
	}
}

package render_test

import (
	"strings"
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/settlement"
)

func TestRenderSettlement(t *testing.T) {
	from := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	inv, _ := invoice.NewBuilder().
		Number("INV-1").
		IssuedAt(from.AddDate(0, 0, 5)).
		Seller(invoice.Party{Name: "Seller"}).
		Buyer(invoice.Party{Name: "Buyer"}).
		AddLine(invoice.NewLineItem("S", 1, invoice.NewMoney(5000, "CHF"), 0)).
		Build()
	rep, err := settlement.Build("Acme", from, to, []settlement.InvoiceInput{
		{Number: "INV-1", Status: "paid", Invoice: *inv},
	})
	if err != nil {
		t.Fatal(err)
	}
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	html, err := engine.RenderSettlement(rep)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Settlement report") || !strings.Contains(html, "INV-1") {
		t.Fatal("missing settlement content")
	}
}

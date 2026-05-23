package render_test

import (
	"strings"
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/tax"
)

func TestRenderDefault(t *testing.T) {
	inv, err := invoice.NewBuilder().
		Number("INV-TEST-001").
		IssuedAt(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)).
		DueAt(time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)).
		Seller(invoice.Party{Name: "Acme GmbH", Email: "billing@acme.example"}).
		Buyer(invoice.Party{Name: "Customer AG"}).
		AddLine(invoice.NewLineItem(
			"Platform subscription",
			1,
			invoice.NewMoney(9900, "CHF"),
			tax.CHStandard.Fraction,
		)).
		Notes("Payment due within 30 days.").
		Build()
	if err != nil {
		t.Fatal(err)
	}

	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}

	html, err := engine.RenderDefault(inv)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"INV-TEST-001",
		"Acme GmbH",
		"Customer AG",
		"Platform subscription",
		"CHF",
		"Payment due within 30 days",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("HTML missing %q", want)
		}
	}
}

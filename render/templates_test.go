package render_test

import (
	"strings"
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/render/i18n"
	"github.com/wiederin/go-invoicer/tax"
)

func sampleInvoice(t *testing.T) *invoice.Invoice {
	t.Helper()
	inv, err := invoice.NewBuilder().
		Number("INV-TPL-1").
		IssuedAt(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)).
		Seller(invoice.Party{Name: "Seller"}).
		Buyer(invoice.Party{Name: "Buyer"}).
		AddLine(invoice.NewLineItem("Item", 1, invoice.NewMoney(1000, "CHF"), tax.CHStandard.Fraction)).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	return inv
}

func TestRenderMinimal(t *testing.T) {
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	html, err := engine.RenderMinimal(sampleInvoice(t))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "INV-TPL-1") {
		t.Fatal("missing invoice number")
	}
}

func TestRenderMultilingualDE(t *testing.T) {
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	html, err := engine.RenderMultilingual(sampleInvoice(t), i18n.LocaleDE)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Rechnung") || !strings.Contains(html, "Zwischensumme") {
		t.Fatal("expected German labels")
	}
}

func TestRenderSwiss(t *testing.T) {
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	inv := sampleInvoice(t)
	inv.Seller.Address = invoice.Address{Line1: "S", PostalCode: "1", City: "Z", Country: "CH"}
	inv.Buyer.Address = invoice.Address{Line1: "B", PostalCode: "2", City: "C", Country: "CH"}

	html, err := engine.RenderSwiss(inv, "CH9300762011623852957")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"INV-TPL-1", "data:image/png;base64,", "IBAN"} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q", want)
		}
	}
}

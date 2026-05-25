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

func TestRenderModern(t *testing.T) {
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	html, err := engine.RenderModern(sampleInvoice(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"INV-TPL-1", "invoice-badge", "Total due"} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q", want)
		}
	}
}

func TestRenderStudio(t *testing.T) {
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	html, err := engine.RenderStudio(sampleInvoice(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"INV-TPL-1", "class=\"hero\"", "Total due"} {
		if !strings.Contains(html, want) {
			t.Fatalf("missing %q", want)
		}
	}
}

func TestRenderBrandingPackTemplates(t *testing.T) {
	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}
	inv := sampleInvoice(t)
	cases := []struct {
		name string
		fn   func(*invoice.Invoice) (string, error)
		want []string
	}{
		{"stratosphere", engine.RenderStratosphere, []string{"INV-TPL-1", "sky-bar", "Total due"}},
		{"ocean", engine.RenderOcean, []string{"INV-TPL-1", "class=\"hero\"", "Total due"}},
		{"ledger", engine.RenderLedger, []string{"INV-TPL-1", "Tax invoice", "Space Mono"}},
		{"mist", engine.RenderMist, []string{"INV-TPL-1", "class=\"sheet\"", "Total due"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			html, err := tc.fn(inv)
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range tc.want {
				if !strings.Contains(html, w) {
					t.Fatalf("missing %q", w)
				}
			}
		})
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

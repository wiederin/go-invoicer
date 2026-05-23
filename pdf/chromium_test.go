//go:build integration

package pdf_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/pdf"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/tax"
)

func TestChromiumRenderInvoice(t *testing.T) {
	inv, err := invoice.NewBuilder().
		Number("INV-PDF-001").
		Seller(invoice.Party{Name: "Seller Co"}).
		Buyer(invoice.Party{Name: "Buyer Co"}).
		AddLine(invoice.NewLineItem("Service", 1, invoice.NewMoney(5000, "CHF"), tax.CHStandard.Fraction)).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	engine, err := render.DefaultEngine()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	renderer := pdf.NewChromiumRenderer()
	out, err := pdf.RenderInvoiceDefault(ctx, engine, renderer, inv)
	if err != nil {
		t.Skipf("chromium not available: %v", err)
	}
	if len(out) < 100 || out[0] != '%' {
		t.Fatalf("expected PDF header, got %d bytes", len(out))
	}
	_ = os.WriteFile("testdata/invoice.pdf", out, 0o644)
}

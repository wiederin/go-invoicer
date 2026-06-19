// Quote example: render a non-binding quote (same pipeline as invoices).
//
// Usage:
//
//	go run ./examples/quote
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/tax"
)

func main() {
	issued := time.Now().UTC()
	inv, err := invoice.NewBuilder().
		Number("QUO-2026-0042").
		Kind(invoice.DocumentQuote).
		IssuedAt(issued).
		DueAt(issued.AddDate(0, 0, 14)).
		Seller(invoice.Party{Name: "Survih GmbH"}).
		Buyer(invoice.Party{Name: "Prospect AG"}).
		AddLine(invoice.NewLineItem(
			"Implementation package",
			1,
			invoice.NewMoney(1200000, "CHF"),
			tax.CHStandard.Fraction,
		)).
		Notes("Valid for 14 days. Prices exclude any third-party fees.").
		Build()
	if err != nil {
		log.Fatal(err)
	}

	engine, err := render.DefaultEngine()
	if err != nil {
		log.Fatal(err)
	}
	html, err := engine.RenderDefault(inv)
	if err != nil {
		log.Fatal(err)
	}

	out := filepath.Join(".", "quote.html")
	if err := os.WriteFile(out, []byte(html), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote", out)
}

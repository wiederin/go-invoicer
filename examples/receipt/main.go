// Receipt example: payment confirmation document (same render pipeline as invoices).
//
// Usage:
//
//	go run ./examples/receipt
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
)

func main() {
	issued := time.Now().UTC()
	inv, err := invoice.NewBuilder().
		Number("REC-2026-0099").
		Kind(invoice.DocumentReceipt).
		IssuedAt(issued).
		Seller(invoice.Party{Name: "Acme GmbH"}).
		Buyer(invoice.Party{Name: "Customer AG"}).
		AddLine(invoice.NewLineItem(
			"Payment received — invoice INV-2026-0042",
			1,
			invoice.NewMoney(5900, "CHF"),
			0,
		)).
		Notes("Thank you for your payment.").
		Build()
	if err != nil {
		log.Fatal(err)
	}

	engine, err := render.DefaultEngine()
	if err != nil {
		log.Fatal(err)
	}
	html, err := engine.RenderMinimal(inv)
	if err != nil {
		log.Fatal(err)
	}

	out := filepath.Join(".", "receipt.html")
	if err := os.WriteFile(out, []byte(html), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote", out)
}

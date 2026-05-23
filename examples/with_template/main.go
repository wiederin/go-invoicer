// Example: render with a custom embedded HTML template.
//
//	go run ./examples/with_template
package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/tax"
)

//go:embed invoice.html
var templateFS embed.FS

func main() {
	inv, err := invoice.NewBuilder().
		Number("INV-CUSTOM-01").
		IssuedAt(time.Now().UTC()).
		Seller(invoice.Party{Name: "Studio Example", Email: "hello@example.com"}).
		Buyer(invoice.Party{Name: "Design Client"}).
		AddLine(invoice.NewLineItem(
			"Brand identity package",
			1,
			invoice.NewMoney(450000, "EUR"),
			tax.DEStandard.Fraction,
		)).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	engine, err := render.NewEngine(templateFS, "invoice.html")
	if err != nil {
		log.Fatal(err)
	}

	html, err := engine.RenderInvoice(inv, "invoice.html")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("invoice-custom.html", []byte(html), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote invoice-custom.html")
}

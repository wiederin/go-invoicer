// Example: Swiss QR-bill invoice (HTML + optional PDF).
//
//	go run ./examples/swiss_qr
//	go run ./examples/swiss_qr -pdf
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/pdf"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/tax"
)

// Demo IBAN (Swiss test account format — replace in production).
const demoIBAN = "CH93 0076 2011 6238 5295 7"

func main() {
	withPDF := flag.Bool("pdf", false, "render PDF")
	outDir := flag.String("out", ".", "output directory")
	flag.Parse()

	issued := time.Now().UTC()
	inv, err := invoice.NewBuilder().
		Number("CH-INV-2026-0042").
		IssuedAt(issued).
		DueAt(issued.AddDate(0, 0, 30)).
		Seller(invoice.Party{
			Name: "Survih GmbH",
			Address: invoice.Address{
				Line1: "Bahnhofstrasse 1", PostalCode: "8001", City: "Zürich", Country: "CH",
			},
		}).
		Buyer(invoice.Party{
			Name: "Musterkunde AG",
			Address: invoice.Address{
				Line1: "Marktgasse 3", PostalCode: "3011", City: "Bern", Country: "CH",
			},
		}).
		AddLine(invoice.NewLineItem(
			"Software-Lizenz Q2",
			1,
			invoice.NewMoney(89000, "CHF"),
			tax.CHStandard.Fraction,
		)).
		Notes("Zahlbar via QR-Rechnung.").
		Build()
	if err != nil {
		log.Fatal(err)
	}

	engine, err := render.DefaultEngine()
	if err != nil {
		log.Fatal(err)
	}

	html, err := engine.RenderSwiss(inv, demoIBAN)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	htmlPath := filepath.Join(*outDir, "invoice-swiss.html")
	if err := os.WriteFile(htmlPath, []byte(html), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote", htmlPath)

	if !*withPDF {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	renderer := pdf.NewChromiumRenderer()
	pdfBytes, err := renderer.RenderHTML(ctx, html)
	if err != nil {
		log.Fatal(err)
	}
	pdfPath := filepath.Join(*outDir, "invoice-swiss.pdf")
	if err := os.WriteFile(pdfPath, pdfBytes, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote", pdfPath)
}

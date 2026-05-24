// Generate example invoice PDFs for all built-in templates.
//
// Requires Chrome/Chromium on PATH.
//
//	go run ./examples/samples
//	go run ./examples/samples -out ./examples/samples
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
	"github.com/wiederin/go-invoicer/render/i18n"
	"github.com/wiederin/go-invoicer/tax"
)

const demoIBAN = "CH93 0076 2011 6238 5295 7"

func main() {
	outDir := flag.String("out", ".", "directory for sample PDF files")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	inv := sampleInvoice()
	engine, err := render.DefaultEngine()
	if err != nil {
		log.Fatal(err)
	}
	renderer := pdf.NewChromiumRenderer()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	jobs := []struct {
		file string
		fn   func() (string, error)
	}{
		{"invoice-default.pdf", func() (string, error) { return engine.RenderDefault(inv) }},
		{"invoice-minimal.pdf", func() (string, error) { return engine.RenderMinimal(inv) }},
		{"invoice-multilingual-en.pdf", func() (string, error) {
			return engine.RenderMultilingual(inv, i18n.LocaleEN)
		}},
		{"invoice-multilingual-de.pdf", func() (string, error) {
			return engine.RenderMultilingual(inv, i18n.LocaleDE)
		}},
		{"invoice-multilingual-fr.pdf", func() (string, error) {
			return engine.RenderMultilingual(inv, i18n.LocaleFR)
		}},
		{"invoice-swiss-qr.pdf", func() (string, error) {
			return engine.RenderSwiss(inv, demoIBAN)
		}},
	}

	for _, job := range jobs {
		html, err := job.fn()
		if err != nil {
			log.Fatalf("%s: render: %v", job.file, err)
		}
		pdfBytes, err := renderer.RenderHTML(ctx, html)
		if err != nil {
			log.Fatalf("%s: pdf: %v", job.file, err)
		}
		path := filepath.Join(*outDir, job.file)
		if err := os.WriteFile(path, pdfBytes, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Wrote", path)
	}
}

func sampleInvoice() *invoice.Invoice {
	issued := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	inv, err := invoice.NewBuilder().
		Number("INV-SAMPLE-2026-001").
		IssuedAt(issued).
		DueAt(issued.AddDate(0, 0, 30)).
		Seller(invoice.Party{
			Name:  "Survih GmbH",
			Email: "billing@survih.ch",
			VATID: "CHE-123.456.789",
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
			"Platform subscription — Pro",
			1,
			invoice.NewMoney(19900, "CHF"),
			tax.CHStandard.Fraction,
		)).
		AddLine(invoice.NewLineItem(
			"Consulting (8.1h)",
			1,
			invoice.NewMoney(121500, "CHF"),
			tax.CHStandard.Fraction,
		)).
		Notes("Thank you for your business. Payment due within 30 days.").
		Build()
	if err != nil {
		log.Fatal(err)
	}
	return inv
}

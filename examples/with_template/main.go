package main

import (
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
	"github.com/wiederin/go-invoicer/template"
)

//go:embed templates/*
var templates embed.FS

func main() {
	inv, err := invoice.New().
		Number("INV-2025-002").
		IssueDate(time.Now()).
		DueDate(time.Now().AddDate(0, 0, 30)).
		Currency("CHF").
		CountryCode("CH").
		Supplier(invoice.Party{
			Name: "Swiss Tech GmbH",
			Address: invoice.Address{
				Street:     "Bahnhofstrasse 1",
				City:       "Zürich",
				PostalCode: "8001",
				Country:    "Switzerland",
			},
			VATID: "CHE-123.456.789",
			IBAN:  "CH93 0076 2011 6238 5295 7",
		}).
		Customer(invoice.Party{
			Name: "German Corp AG",
			Address: invoice.Address{
				Street:     "Hauptstraße 10",
				City:       "Berlin",
				PostalCode: "10115",
				Country:    "Germany",
			},
			VATID: "DE123456789",
		}).
		AddItem(invoice.NewLineItem(
			"Software Development (40 hours)",
			40,
			invoice.NewMoney(150.00, "CHF"),
			8.1,
		)).
		AddItem(invoice.NewLineItem(
			"Server Infrastructure",
			1,
			invoice.NewMoney(500.00, "CHF"),
			8.1,
		)).
		Notes("Payment via bank transfer to the IBAN above.").
		Terms("Net 30 days. Late payments subject to 5% interest.").
		Build()

	if err != nil {
		fmt.Printf("Error building invoice: %v\n", err)
		os.Exit(1)
	}

	tmplManager := template.NewManager(
		template.NewEmbedSource(templates),
	)

	engine := render.NewEngine(tmplManager)
	pdf, err := engine.RenderInvoice(inv, "templates/invoice.html")
	if err != nil {
		fmt.Printf("Error rendering PDF: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("invoice_swiss.pdf", pdf, 0644); err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Invoice saved to invoice_swiss.pdf")
}

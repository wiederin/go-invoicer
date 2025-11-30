package main

import (
	"fmt"
	"os"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
)

func main() {
	inv, err := invoice.New().
		Number("INV-2025-001").
		IssueDate(time.Now()).
		DueDate(time.Now().AddDate(0, 0, 30)).
		Currency("USD").
		CountryCode("US").
		Supplier(invoice.Party{
			Name: "Acme Corporation",
			Address: invoice.Address{
				Street:     "123 Business Street",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "United States",
			},
			Email: "billing@acme.com",
			VATID: "US-123456789",
		}).
		Customer(invoice.Party{
			Name: "John Smith",
			Address: invoice.Address{
				Street:     "456 Customer Ave",
				City:       "Los Angeles",
				State:      "CA",
				PostalCode: "90001",
				Country:    "United States",
			},
			Email: "john@example.com",
		}).
		AddItem(invoice.NewLineItem(
			"Web Design Services",
			1,
			invoice.NewMoney(1500.00, "USD"),
			8.0,
		)).
		AddItem(invoice.NewLineItem(
			"Hosting (12 months)",
			12,
			invoice.NewMoney(29.99, "USD"),
			8.0,
		)).
		AddItem(invoice.NewLineItem(
			"Domain Registration",
			1,
			invoice.NewMoney(15.00, "USD"),
			0.0,
		)).
		Notes("Thank you for your business!").
		Terms("Payment is due within 30 days of invoice date.").
		Build()

	if err != nil {
		fmt.Printf("Error building invoice: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Invoice: %s\n", inv.Number)
	fmt.Printf("Subtotal: %s\n", inv.SubTotal())
	fmt.Printf("Tax: %s\n", inv.TotalTax())
	fmt.Printf("Total: %s\n", inv.TotalGross())

	renderer := render.NewSimpleRenderer()
	pdf, err := renderer.RenderInvoice(inv)
	if err != nil {
		fmt.Printf("Error rendering PDF: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("invoice.pdf", pdf, 0644); err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Invoice saved to invoice.pdf")
}

# go-invoicer

A professional invoice PDF generator library for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/wiederin/go-invoicer.svg)](https://pkg.go.dev/github.com/wiederin/go-invoicer)
[![Go Report Card](https://goreportcard.com/badge/github.com/wiederin/go-invoicer)](https://goreportcard.com/report/github.com/wiederin/go-invoicer)

## Features

- **Type-safe invoice creation** with fluent builder API
- **Automatic calculations** for subtotals, taxes, discounts, and totals
- **Multi-currency support** with proper formatting (USD, EUR, CHF, GBP, etc.)
- **Tax helpers** with common rates for Switzerland, EU, UK, and more
- **Template engine** with Go templates and embedded template support
- **PDF rendering** with customizable layouts
- **Validation** with clear error messages

## Installation

```bash
go get github.com/wiederin/go-invoicer
```

## Quick Start

```go
package main

import (
    "os"
    "time"

    "github.com/wiederin/go-invoicer/invoice"
    "github.com/wiederin/go-invoicer/render"
)

func main() {
    // Build an invoice
    inv, err := invoice.New().
        Number("INV-2025-001").
        IssueDate(time.Now()).
        DueDate(time.Now().AddDate(0, 0, 30)).
        Currency("USD").
        Supplier(invoice.Party{
            Name: "Your Company",
            Address: invoice.Address{
                Street: "123 Business St",
                City: "New York",
                PostalCode: "10001",
                Country: "USA",
            },
        }).
        Customer(invoice.Party{
            Name: "Customer Inc",
            Address: invoice.Address{
                Street: "456 Client Ave",
                City: "Los Angeles",
                PostalCode: "90001",
                Country: "USA",
            },
        }).
        AddItem(invoice.NewLineItem(
            "Consulting Services",
            10, // hours
            invoice.NewMoney(150, "USD"),
            8.0, // tax rate %
        )).
        Notes("Thank you for your business!").
        Build()

    if err != nil {
        panic(err)
    }

    // Render to PDF
    renderer := render.NewSimpleRenderer()
    pdf, _ := renderer.RenderInvoice(inv)
    os.WriteFile("invoice.pdf", pdf, 0644)
}
```

## Packages

### `invoice` - Core Domain Models

Create and validate invoices with automatic totals calculation:

```go
// Create invoice with builder
inv, err := invoice.New().
    Number("INV-001").
    IssueDate(time.Now()).
    DueDate(time.Now().AddDate(0, 0, 30)).
    Currency("CHF").
    Supplier(supplier).
    Customer(customer).
    AddItem(item1).
    AddItem(item2).
    Build()

// Access calculated totals
fmt.Println(inv.SubTotal())      // Net amount before tax
fmt.Println(inv.TotalTax())      // Total tax amount
fmt.Println(inv.TotalGross())    // Final amount due
fmt.Println(inv.TaxBreakdown())  // Tax by rate
```

### `tax` - Tax Calculations

Common tax rates and calculators:

```go
import "github.com/wiederin/go-invoicer/tax"

// Use predefined rates
swissVAT, _ := tax.GetRate("CH_STANDARD")  // 8.1%
germanVAT, _ := tax.GetRate("DE_STANDARD") // 19%

// Calculate tax
amount := decimal.NewFromFloat(1000)
taxAmount := swissVAT.Calculate(amount)  // 81.00
gross := swissVAT.AddToAmount(amount)    // 1081.00

// Custom calculator
calc := tax.NewCalculator().AddRate(swissVAT)
totalTax := calc.CalculateTax(amount)
```

### `currency` - Currency Formatting

Format amounts in different currencies:

```go
import "github.com/wiederin/go-invoicer/currency"

formatter, _ := currency.NewFormatter("CHF")
formatted := formatter.Format(decimal.NewFromFloat(1234.56))
// Output: CHF 1'234.56

formatter, _ = currency.NewFormatter("EUR")
formatted = formatter.Format(decimal.NewFromFloat(1234.56))
// Output: 1.234,56 €
```

### `template` - Template Engine

Use Go templates with embedded or file-based templates:

```go
import (
    "embed"
    "github.com/wiederin/go-invoicer/template"
)

//go:embed templates/*
var templates embed.FS

mgr := template.NewManager(
    template.NewEmbedSource(templates),
)

html, err := mgr.RenderHTML("templates/invoice.html", data)
```

### `render` - PDF Rendering

Generate PDFs with customizable options:

```go
import "github.com/wiederin/go-invoicer/render"

// Simple renderer (no template needed)
renderer := render.NewSimpleRenderer()
pdf, err := renderer.RenderInvoice(inv)

// Template-based renderer
engine := render.NewEngine(templateManager)
pdf, err := engine.RenderInvoice(inv, "invoice.html")

// Custom options
renderer := render.NewSimpleRenderer()
renderer.Options = render.Options{
    PageSize:    "A4",
    Orientation: "P",
    MarginTop:   15,
    FontFamily:  "Helvetica",
}
```

## Line Items

Create line items with quantities, prices, and optional discounts:

```go
// Basic line item
item := invoice.NewLineItem(
    "Web Development",  // description
    10,                 // quantity
    invoice.NewMoney(100, "USD"),  // unit price
    8.0,                // tax rate %
)

// With discount
item = item.WithDiscount(10)  // 10% discount

// Access calculations
item.SubTotal()      // qty * unit price
item.DiscountAmount() // discount amount
item.NetAmount()     // after discount
item.TaxAmount()     // tax on net
item.GrossAmount()   // final amount
```

## Validation

Invoices are validated on build with clear errors:

```go
inv, err := invoice.New().Build()
// err = ErrMissingInvoiceNumber

inv, err := invoice.New().
    Number("INV-001").
    IssueDate(time.Now()).
    DueDate(time.Now().AddDate(0, 0, -1)). // past date
    Build()
// err = ErrDueDateBeforeIssue
```

## Examples

See the [examples](examples/) directory for complete working examples:

- `basic/` - Simple invoice generation
- `with_template/` - Custom HTML templates
- `swiss_qr/` - Swiss QR bill (coming soon)

## Roadmap

- [x] Core invoice models with validation
- [x] Tax and currency helpers
- [x] Template engine with embedding
- [x] PDF rendering
- [ ] Swiss QR bill support
- [ ] EU e-invoicing formats
- [ ] Digital signatures
- [ ] Hosted API service

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

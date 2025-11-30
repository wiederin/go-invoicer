# go-invoicer

A professional invoice PDF generator library written in Go.

## Overview

`go-invoicer` is a lightweight, easy-to-use Go library for generating professional PDF invoices. It uses [gofpdf](https://github.com/go-pdf/fpdf) under the hood to create beautifully formatted invoices with customizable headers, footers, company details, line items, and payment terms.

## Features

- 📄 Generate professional PDF invoices
- 🎨 Customizable headers and footers
- 💼 Company and customer details
- 📊 Line items with quantities, unit costs, and descriptions
- 💰 Multi-currency support
- 📝 Custom notes and payment terms
- 🔢 Automatic total calculations

## Installation

```bash
go get github.com/wiederin/go-invoicer
```

## Quick Start

```go
package main

import (
    "github.com/wiederin/go-invoicer"
    "github.com/wiederin/go-invoicer/components"
    "github.com/wiederin/go-invoicer/constants"
)

func main() {
    // Initialize invoice document
    doc, err := invoicer.Init(&components.Config{
        TextInvoiceType:   constants.Invoice,
        TextRefTitle:      "Invoice No.",
        CurrencySymbol:    "$",
        CurrencyPrecision: 2,
        CurrencyDecimal:   ".",
        CurrencyThousand:  ",",
    })
    if err != nil {
        panic(err)
    }

    // Set invoice number
    doc.SetVersion("INV-2025-001")

    // Set company details
    doc.SetCompany(&components.Contact{
        Name: "Acme Corporation",
        Address: &components.Address{
            Address:    "123 Business Street",
            PostalCode: "10001",
            City:       "New York",
        },
    })

    // Set customer details
    doc.SetCustomer(&components.Contact{
        Name: "John Smith",
        Address: &components.Address{
            Address:    "456 Customer Ave",
            PostalCode: "90001",
            City:       "Los Angeles",
            Country:    "United States",
        },
    })

    // Add description
    doc.SetDescription("Thank you for your business!")

    // Add line items
    doc.AppendItem(&components.Item{
        Name:        "Web Design Services",
        Description: "Custom website design and development",
        UnitCost:    "1500.00",
        Quantity:    "1",
    })

    // Add payment notes
    doc.SetNotes("Payment is due within 30 days.")

    // Build the PDF
    pdf, err := invoicer.Build(doc)
    if err != nil {
        panic(err)
    }

    // Save to file
    err = pdf.OutputFileAndClose("invoice.pdf")
    if err != nil {
        panic(err)
    }
}
```

## API Reference

### Initialization

```go
doc, err := invoicer.Init(&components.Config{
    TextInvoiceType:   constants.Invoice,
    TextRefTitle:      "Invoice No.",
    CurrencySymbol:    "$",
    CurrencyPrecision: 2,
    CurrencyDecimal:   ".",
    CurrencyThousand:  ",",
})
```

### Setting Invoice Details

```go
// Invoice number
doc.SetVersion("INV-001")

// Header
doc.SetHeader(&components.HeaderFooter{
    Text:       "<center>INVOICE</center>",
    Pagination: true,
})

// Footer
doc.SetFooter(&components.HeaderFooter{
    Text:       "<center>Thank you for your business</center>",
    Pagination: true,
})

// Description
doc.SetDescription("Description text...")

// Notes
doc.SetNotes("Payment terms and additional notes...")
```

### Company and Customer Information

```go
// Company
doc.SetCompany(&components.Contact{
    Name: "Company Name",
    Logo: logoBytes, // optional []byte
    Address: &components.Address{
        Address:    "Street Address",
        PostalCode: "12345",
        City:       "City",
    },
})

// Customer
doc.SetCustomer(&components.Contact{
    Name: "Customer Name",
    Address: &components.Address{
        Address:    "Street Address",
        PostalCode: "12345",
        City:       "City",
        Country:    "Country",
    },
})
```

### Adding Line Items

```go
doc.AppendItem(&components.Item{
    Name:        "Product/Service Name",
    Description: "Description",
    UnitCost:    "100.00",
    Quantity:    "2",
})
```

### Building the PDF

```go
pdf, err := invoicer.Build(doc)

// Save to file
err = pdf.OutputFileAndClose("invoice.pdf")

// Or write to io.Writer
var buf bytes.Buffer
err = pdf.Output(&buf)
pdfBytes := buf.Bytes()
```

## Project Structure

```
go-invoicer/
├── components/       # Component types (Contact, Address, Item, etc.)
├── constants/        # Constants used throughout the library
├── build.go         # PDF building logic
├── generator.go     # Invoice initialization
└── generator_test.go # Tests
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and feature requests, please use the GitHub issue tracker.

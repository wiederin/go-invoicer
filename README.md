# go-invoicer

[![Go Reference](https://pkg.go.dev/badge/github.com/wiederin/go-invoicer.svg)](https://pkg.go.dev/github.com/wiederin/go-invoicer)

Developer-first invoicing library for Go — strongly typed, HTML templates, Chromium PDF export.

> Mirrored from the private [GitLab monorepo](https://gitlab.com/survih/go-invoicer) (`oss/`). Issues and PRs: [GitHub](https://github.com/wiederin/go-invoicer).

## Features

- **Invoice domain** — builder API, validation, tax breakdown
- **Currency formatting** — CHF, EUR, USD, GBP
- **Tax helpers** — common VAT rates (CH, DE, UK)
- **HTML rendering** — embedded default template + custom templates via `embed.FS`
- **PDF export** — headless Chromium (`pdf.ChromiumRenderer`)
- **Stripe import (MVP)** — map Stripe invoice payloads to domain types

## Install

```bash
go get github.com/wiederin/go-invoicer@latest
```

Requires **Go 1.26+** (for PDF/Chromium dependencies).

## Quick start

```go
inv, err := invoice.NewBuilder().
    Number("INV-001").
    Seller(invoice.Party{Name: "Acme GmbH"}).
    Buyer(invoice.Party{Name: "Customer AG"}).
    AddLine(invoice.NewLineItem(
        "Consulting", 10,
        invoice.NewMoney(15000, "CHF"),
        tax.CHStandard.Fraction,
    )).
    Build()

engine, _ := render.DefaultEngine()
html, _ := engine.RenderDefault(inv)
```

### PDF (requires Chrome/Chromium)

```go
renderer := pdf.NewChromiumRenderer()
pdfBytes, err := pdf.RenderInvoiceDefault(ctx, engine, renderer, inv)
```

## Example

```bash
go run ./examples/basic
go run ./examples/basic -pdf -out ./out
```

## Packages

| Package | Description |
|---------|-------------|
| `invoice` | Domain model, builder, validation |
| `currency` | Display formatting |
| `tax` | VAT rate helpers |
| `render` | HTML template engine |
| `templates` | Embedded default templates |
| `pdf` | HTML → PDF via Chromium |
| `integrations/stripe` | Stripe → domain mapper (MVP) |

## License

Apache License 2.0 — see [LICENSE](LICENSE).

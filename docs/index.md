# go-invoicer

Developer-first invoicing for Go: typed invoices, HTML templates, Chromium PDF export, Swiss QR-bills, and Stripe import helpers.

```bash
go get github.com/wiederin/go-invoicer@latest
```

## Highlights

- **Invoice builder** with validation and tax breakdown
- **HTML templates** — default, minimal, Swiss QR-bill
- **PDF export** via headless Chrome (`pdf.ChromiumRenderer`)
- **Stripe webhooks** — parse `invoice.paid` and map to domain types

See [Getting started](getting-started.md).

## Build docs locally

```bash
pip install mkdocs-material
cd oss && mkdocs serve
```

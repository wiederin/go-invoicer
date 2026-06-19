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

Using the **hosted dashboard and API**?

**On the website:** [invoicer.survih.ch/docs](https://invoicer.survih.ch/docs) — full guides (Stripe, Xero, QuickBooks, paid features) without leaving the dashboard.

**In this repo (MkDocs / GitHub):**

- [Hosted platform — plans & paid features](hosted-platform.md) — tiers, feature matrix, gates
- [Paid features guide](paid-features.md) — how to use each integration from the API
- [Stripe integration](stripe-integration.md) · [Replace Stripe Invoicing](replace-stripe-invoicing.md)
- [Xero](xero-integration.md) · [QuickBooks](quickbooks-integration.md) · [Accounting overview](accounting-integrations.md)

## Build docs locally

```bash
pip install mkdocs-material
cd oss && mkdocs serve
```

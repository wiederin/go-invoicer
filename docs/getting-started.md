# Getting started

## Install

```bash
go get github.com/wiederin/go-invoicer@latest
```

Requires Go 1.26+ and Chrome/Chromium for PDF export.

## Create an invoice

```go
inv, err := invoice.NewBuilder().
    Number("INV-001").
    Seller(invoice.Party{Name: "Your Company"}).
    Buyer(invoice.Party{Name: "Customer"}).
    AddLine(invoice.NewLineItem(
        "Services", 1,
        invoice.NewMoney(10000, "CHF"),
        tax.CHStandard.Fraction,
    )).
    Build()
```

## Render HTML

```go
engine, _ := render.DefaultEngine()
html, _ := engine.RenderDefault(inv)
```

## Export PDF

```go
ctx := context.Background()
pdfBytes, err := pdf.RenderInvoiceDefault(ctx, engine, pdf.NewChromiumRenderer(), inv)
```

## Templates

```go
engine, _ := render.DefaultEngine()
html, _ := engine.RenderMinimal(inv)
html, _ := engine.RenderSwiss(inv, "CH93...") // Swiss QR-bill
```

## Stripe

Map Stripe invoice payloads with `integrations/stripe`:

```go
domain := stripeInvoice.ToDomain(stripe.ImportOptions{Supplier: supplier})
```

Parse webhooks:

```go
ev, _ := stripe.ParseWebhookEvent(body)
apiInv, _ := stripe.InvoiceFromEvent(ev)
domain := apiInv.ToDomainInvoice().ToDomain(opts)
```

## Hosted platform

The OSS library runs locally without a subscription. If you use the **hosted API and dashboard**, plans gate integrations, compliance exports, e-invoices, and logo/accent branding — while line-item columns, numbering, and cron schedules stay on every tier. See [Hosted platform — plans & paid features](hosted-platform.md).

## Release v0.1.0

After mirroring to GitHub, tag the public repo:

```bash
git tag v0.1.0 && git push origin v0.1.0
```

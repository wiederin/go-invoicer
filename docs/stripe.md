# Stripe integration

Import Stripe invoices into the go-invoicer domain model and handle webhooks.

## Domain mapper

```go
si := stripe.Invoice{ /* from API or webhook JSON */ }
inv := si.ToDomain(stripe.ImportOptions{Supplier: supplier})
if err := inv.Validate(); err != nil { ... }
```

## API client

```go
client := stripe.NewClient(os.Getenv("STRIPE_SECRET_KEY"))
inv, err := client.SyncInvoice(ctx, "in_123", stripe.ImportOptions{Supplier: supplier})
```

## Webhooks

The OSS package provides an `http.Handler` with optional signature verification:

```go
http.Handle("/webhooks/stripe", stripe.WebhookHandler(stripe.HandlerOptions{
    WebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
    Supplier:      supplier,
    OnInvoice: func(inv *invoice.Invoice) error {
        // persist, enqueue PDF render, etc.
        return nil
    },
}))
```

Supported events: `invoice.finalized`, `invoice.paid`.

### Platform API

Hosted endpoint (no platform API key; Stripe signs requests):

`POST /v1/webhooks/stripe`

## Example

```bash
export STRIPE_SECRET_KEY=sk_test_...
go run ./examples/stripe_sync -invoice in_xxx
```

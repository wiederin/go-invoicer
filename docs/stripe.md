# Stripe (OSS library)

Apache 2.0 package: `github.com/wiederin/go-invoicer/integrations/stripe`

Use this in your own Go services. The **hosted platform** adds multi-tenant credentials, PDF upload, and plan gates — see [Stripe integration](stripe-integration.md).

---

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

// Export
stripeID, err := client.PushInvoice(ctx, domainInv, stripe.ExportOptions{Supplier: supplier})
err = client.UploadPDF(ctx, "invoice.pdf", pdfBytes)
err = client.LinkBrandedPDF(ctx, stripeID, pdfURL, jobID, fileID)
err = client.FinalizeInvoice(ctx, stripeID)
err = client.SendInvoice(ctx, stripeID)
```

## Webhooks

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

**Supported events:** `invoice.finalized`, `invoice.paid`

### Hosted platform endpoint

Per-organization URL (Stripe signs with the org's `whsec_…`):

```text
POST /v1/webhooks/stripe/:org_id
```

Platform sets local invoice `status` to `paid` on `invoice.paid`. See [Replace Stripe Invoicing](replace-stripe-invoicing.md).

## Example

```bash
export STRIPE_SECRET_KEY=sk_test_…
go run ./examples/stripe_sync -invoice in_xxx
```

## Hosted guides

- [Stripe integration](stripe-integration.md) — connect, import, export, automation
- [Replace Stripe Invoicing](replace-stripe-invoicing.md) — payments-only Stripe + tracking
- [Paid features](paid-features.md) — plan requirements

# Stripe integration (hosted platform)

Connect **your organization's Stripe account** to import invoices, export branded PDFs back to Stripe, and keep payment state in sync via webhooks.

!!! note "Two Stripe accounts"
    - **Platform billing** — pays for go-invoicer (Starter/Pro/…). Uses server env `STRIPE_SECRET_KEY` and `POST /v1/webhooks/platform-billing`.
    - **Per-org invoicing** — your customers' invoices. Uses keys saved in **Settings** and `POST /v1/webhooks/stripe/:org_id`.

    This page covers **per-org invoicing** only.

**Plan required:** `stripe` feature (Starter and above). See [Paid features](paid-features.md).

---

## 1. Connect Stripe

1. Upgrade to **Starter** or higher (`POST /v1/billing/checkout`).
2. Open **Settings → Integrations → Stripe**.
3. Paste:
   - **Secret key** — `sk_live_…` or `sk_test_…`
   - **Webhook signing secret** — `whsec_…` from the Stripe Dashboard endpoint you create in step 4.
4. Save (`PUT /v1/org/integrations/stripe`).

```http
PUT /v1/org/integrations/stripe
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "secret_key": "sk_test_…",
  "webhook_secret": "whsec_…"
}
```

```http
GET /v1/org/integrations/stripe
```

Response includes `webhook_url` — register this in Stripe:

```text
https://<your-api-host>/v1/webhooks/stripe/<org_id>
```

`org_id` is your organization UUID (visible in Settings or from `GET /v1/billing/status`).

### Stripe Dashboard — webhook events

Create an endpoint pointing at the URL above. Enable:

| Event | Purpose |
|-------|---------|
| `invoice.finalized` | Import / refresh invoice; queue branded PDF |
| `invoice.paid` | Re-import invoice; set local status to **`paid`** |

The platform verifies the `Stripe-Signature` header with your `whsec_…` secret.

### Verify connection

```http
GET /v1/integrations
```

```json
{
  "integrations": [
    { "id": "stripe", "name": "Stripe", "status": "ready", "description": "…" }
  ]
}
```

---

## 2. Import (Stripe → go-invoicer)

### Manual import

```http
POST /v1/integrations/stripe/import
Authorization: Bearer <token>
Content-Type: application/json

{ "stripe_invoice_id": "in_1ABC…" }
```

Flow:

1. Fetch invoice from Stripe API (`sk_…` from Settings).
2. Map to go-invoicer domain model (`oss/integrations/stripe`).
3. Upsert locally (`external_source: stripe`, `external_id: in_…`).
4. Enqueue async PDF render with org template defaults.

### Webhook import

When Stripe sends `invoice.finalized` or `invoice.paid`, the platform:

1. Validates signature.
2. Maps payload → domain invoice.
3. Upserts with status `synced` (finalized) or **`paid`** (`invoice.paid`).
4. Queues PDF render.

---

## 3. Export (go-invoicer → Stripe)

### Create Stripe draft from local invoice

```http
POST /v1/invoices/:id/sync/stripe
Authorization: Bearer <token>
```

1. Creates or updates a **draft** Stripe invoice with line items and customer mapping.
2. Stores `stripe_invoice_id` in invoice metadata.
3. Queues branded PDF generation.

### After PDF renders

When the PDF job completes, the platform automatically:

1. Uploads the PDF to Stripe Files.
2. Links it on the Stripe invoice (metadata + custom field “Branded invoice”).
3. If **Stripe billing mode** is `finalize_send` → calls Stripe **Finalize** + **Send** (see [Replace Stripe Invoicing](replace-stripe-invoicing.md)).

### Manual finalize & send

```http
POST /v1/invoices/:id/stripe/bill
Authorization: Bearer <token>
```

Use when mode is `attach_only` or you want to bill on demand.

---

## 4. Stripe billing mode (Settings)

| Mode | Behavior |
|------|----------|
| **`attach_only`** | Branded PDF attached to Stripe draft; you finalize/send in Stripe Dashboard or via `/stripe/bill`. |
| **`finalize_send`** | After PDF upload, platform auto-finalizes and emails the invoice through Stripe. |

Set in **Settings → Integrations → Stripe** (`stripe_billing_mode` on rendering profile).

---

## 5. Automation API

### Generate with Stripe sync

```http
POST /v1/invoices/generate
Authorization: Bearer <token>

{
  "invoice": { … },
  "status": "draft",
  "render_pdf": true,
  "sync_stripe": true
}
```

Same checks as manual export: `stripe` feature + configured credentials.

### Inbound hook

```http
POST /v1/hooks/{token}/invoices

{
  "invoice": { … },
  "render_pdf": true,
  "sync_stripe": true
}
```

No Bearer token — use the secret URL from **Settings → Automation**.

---

## 6. OSS library (self-hosted / custom apps)

The Apache 2.0 package `github.com/wiederin/go-invoicer/integrations/stripe` provides:

- **Mapper** — `stripe.Invoice.ToDomain()`
- **Client** — `FetchInvoice`, `PushInvoice`, `UploadPDF`, `LinkBrandedPDF`, `FinalizeInvoice`, `SendInvoice`
- **Webhook handler** — `stripe.WebhookHandler()` for `invoice.finalized` / `invoice.paid`

Example:

```bash
export STRIPE_SECRET_KEY=sk_test_…
go run ./examples/stripe_sync -invoice in_xxx
```

Hosted platform wiring: `platform/api/internal/sync/stripe/service.go`.

See [Stripe (OSS)](stripe.md) for code-level API.

---

## 7. Metadata & linking

| Key | Set by | Purpose |
|-----|--------|---------|
| `stripe_invoice_id` | Platform | Links local record ↔ `in_…` |
| `go_invoicer_pdf_url` | Export | Stripe invoice metadata |
| `go_invoicer_render_job_id` | Export | Trace PDF job |
| `go_invoicer_pdf_file_id` | Export | Stripe File id for PDF |

List invoices filtered by external id via stored `external_source` / `external_id` on the invoice record.

---

## 8. Troubleshooting

| Symptom | Check |
|---------|--------|
| **402** on Stripe endpoints | Plan includes `stripe`; `GET /v1/billing/status` |
| Integration `not_configured` | Secret key saved; `GET /v1/integrations` |
| Webhook 400 invalid signature | `whsec_…` matches endpoint; raw body not modified by proxy |
| PDF not on Stripe invoice | Export ran first (`sync/stripe`); wait for render job; check job PDF |
| Duplicate S3 / deploy errors | Staging/prod use separate `TF_VAR_project_name` stacks |

---

## Next steps

- [Replace Stripe Invoicing](replace-stripe-invoicing.md) — architecture when Stripe is payments-only
- [Paid features](paid-features.md) — Xero, QBO, compliance, branding

# Xero integration (hosted platform)

Sync **accounts receivable (ACCREC)** invoices between Xero and go-invoicer. Use go-invoicer for **branded PDFs** and template control; use Xero for **books, tax, and payment recording**.

**Plan required:** `xero` feature (**Pro** and above). See [Hosted platform](hosted-platform.md).

**Compare:** [QuickBooks integration](quickbooks-integration.md) · [Accounting overview](accounting-integrations.md) · [Stripe](stripe-integration.md)

---

## What you get

| Capability | Supported |
|------------|-----------|
| OAuth connect (per org) | ✓ |
| Import invoice by Xero `InvoiceID` (UUID) | ✓ |
| Export local invoice → Xero **DRAFT** ACCREC | ✓ |
| Auto-create Xero **Contact** from buyer name/email | ✓ |
| Branded PDF after import/export | ✓ (async queue) |
| Xero → go-invoicer webhooks | ✗ (re-import manually or via your automation) |
| Attach PDF to Xero invoice via API | ✗ (download PDF from go-invoicer; attach in Xero if needed) |

---

## Prerequisites

### 1. Platform OAuth app (operator / self-hosted)

Register a [Xero app](https://developer.xero.com/app/manage) and set the redirect URI:

```text
https://<api-host>/v1/oauth/xero/callback
```

Configure on the API server:

| Variable | Purpose |
|----------|---------|
| `GO_INVOICER_XERO_CLIENT_ID` | Xero app client ID |
| `GO_INVOICER_XERO_CLIENT_SECRET` | Xero app client secret |
| `GO_INVOICER_PUBLIC_API_URL` | Used to build callback URL |

See `docs/gitlab-deploy.md` in the monorepo for CI/Terraform.

**Required Xero scopes** (requested by the platform):

- `openid`, `profile`, `email`
- `accounting.transactions`
- `accounting.contacts`
- `offline_access`

### 2. Organization subscription

```http
GET /v1/billing/status
→ subscription.features includes "xero"
```

Upgrade via **Billing** → Checkout (**Pro** or higher) if missing.

---

## Connect Xero (dashboard)

1. Sign in as **owner** or **admin**.
2. **Settings → Integrations → Xero** → **Connect**.
3. Authorize the Xero organization (tenant). The platform stores access + refresh tokens per org.
4. On success you are redirected to Settings with `xero=connected`.

### API connect flow

```http
GET /v1/oauth/xero/connect
Authorization: Bearer <admin-token>
```

→ Redirects to Xero (or `?format=json` → `{ "authorize_url": "…" }`).

Callback (browser, no auth):

```text
GET /v1/oauth/xero/callback?code=…&state=…
```

Check status:

```http
GET /v1/org/integrations/xero
Authorization: Bearer <token>
```

```json
{
  "xero": {
    "configured": true,
    "connectable": true,
    "tenant_id": "…",
    "tenant_name": "Acme Ltd"
  }
}
```

Disconnect:

```http
DELETE /v1/org/integrations/xero
Authorization: Bearer <admin-token>
```

---

## Import (Xero → go-invoicer)

### Find the Xero invoice ID

In Xero: open the invoice → URL contains the UUID, or use the API. This is **`InvoiceID`**, not the invoice number.

### Import via API

```http
POST /v1/integrations/xero/import
Authorization: Bearer <token>
Content-Type: application/json

{ "xero_invoice_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890" }
```

Response:

```json
{
  "invoice": { "id": "…", "number": "INV-0042", "status": "synced", … },
  "message": "imported; PDF render queued"
}
```

**Flow:**

1. Fetch invoice from Xero Accounting API.
2. Map lines, buyer (contact), dates, currency → domain model.
3. Upsert locally: `external_source: xero`, `external_id: <InvoiceID>`.
4. Queue PDF render (org template, locale, IBAN from **Settings**).

Counts toward monthly invoice quota (same as `POST /v1/invoices`).

### Dashboard

**Settings → Integrations → Xero** — paste UUID and **Import**.

---

## Export (go-invoicer → Xero)

```http
POST /v1/invoices/:id/sync/xero
Authorization: Bearer <token>
```

**Flow:**

1. Load local invoice.
2. If `metadata.xero_invoice_id` exists → **update** Xero invoice; else **create** ACCREC **DRAFT**.
3. Match or create **Contact** by buyer email (preferred) or name.
4. Upsert local record with Xero link; queue branded PDF.

You finalize, email, and record payments in **Xero** (or your bank feed). go-invoicer does not change Xero invoice status to AUTHORISED/PAID automatically.

### Buyer requirements

Export needs a resolvable buyer:

- **Email** — lookup existing contact, or
- **Name** — create new contact

Missing both → `400` with contact error.

---

## Branded PDFs

After import or export, the platform enqueues `POST /v1/invoices/:id/render-pdf` internally using:

- Template, locale, IBAN from `GET/PATCH /v1/org/rendering`
- Logo/accent if plan includes `branding`

Download the PDF from the invoice row when the job completes. Upload to Xero manually if your workflow requires the PDF inside Xero; the accounting record is already linked via `external_id`.

---

## Linking & metadata

| Field | Value |
|-------|--------|
| `external_source` | `xero` |
| `external_id` | Xero `InvoiceID` (UUID) |
| `invoice.metadata.xero_invoice_id` | Same UUID |
| Local `status` on import | `synced` |

Re-exporting the same local invoice updates the linked Xero draft (same UUID).

---

## Payment & status tracking

Xero is the **system of record for AR and bank reconciliation**. go-invoicer mirrors document data.

| Goal | Approach |
|------|----------|
| Mark paid in go-invoicer | `PUT /v1/invoices/:id` with `"status": "paid"` when you know Xero shows paid |
| Refresh from Xero | Re-run `POST /v1/integrations/xero/import` with the same `xero_invoice_id` |
| Automate | Cron job or ERP webhook → import API when Xero invoice status changes (build your worker; platform webhooks planned Phase 6) |
| Compliance archive | [Compliance export](paid-features.md#compliance-zip-compliance--pro) (Pro+) |

There is **no** inbound Xero webhook in the platform today.

---

## Recommended workflows

### A — Books in Xero, PDFs in go-invoicer

1. Create ACCREC invoice in Xero (or import from another system).
2. Import into go-invoicer → branded PDF.
3. Send PDF to customer from your process; record payment in Xero.

### B — Author in go-invoicer, post to Xero

1. Create invoice in dashboard or `POST /v1/invoices/generate`.
2. `POST /v1/invoices/:id/sync/xero` → draft in Xero.
3. Download PDF from go-invoicer; approve and send from Xero.

### C — With Stripe for payments

Use [Stripe integration](stripe-integration.md) for card/SEPA collection and Xero for GL — common on **Pro+** (both feature slugs).

---

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| **402** on import/export | Plan must include `xero`; check `GET /v1/billing/status` |
| `connect Xero in Settings first` | Complete OAuth; `GET /v1/org/integrations/xero` → `configured: true` |
| `OAuth not configured on server` | Set `GO_INVOICER_XERO_CLIENT_ID` / `SECRET` on API |
| `invalid_state` on callback | Retry connect; do not share callback URL |
| Scope / 403 from Xero | Add scopes in Xero app; **disconnect** and reconnect org |
| `buyer name required` | Set buyer name or email on invoice before export |
| Wrong tenant | Disconnect and reconnect; pick correct Xero org at authorize |
| Token expired | Platform refreshes automatically when within 2 minutes of expiry |

---

## OSS library

Package: `github.com/wiederin/go-invoicer/integrations/xero`

- `Client.FetchInvoice`, `Client.PushInvoice`
- `Invoice.ToDomain`, OAuth helpers

Hosted wiring: `platform/api/internal/sync/xero/service.go`.

---

## API reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/v1/oauth/xero/connect` | Admin | Start OAuth |
| GET | `/v1/oauth/xero/callback` | Public | OAuth callback |
| GET | `/v1/org/integrations/xero` | Member+ | Connection status |
| DELETE | `/v1/org/integrations/xero` | Admin | Disconnect |
| POST | `/v1/integrations/xero/import` | Write + `xero` | Import by UUID |
| POST | `/v1/invoices/:id/sync/xero` | Write + `xero` | Export to Xero draft |
| GET | `/v1/integrations` | Read | All integration statuses |

---

## Related

- [QuickBooks integration](quickbooks-integration.md)
- [Paid features](paid-features.md)
- [Replace Stripe Invoicing](replace-stripe-invoicing.md) — payments vs accounting

# QuickBooks Online integration (hosted platform)

Sync invoices between **QuickBooks Online (QBO)** and go-invoicer. Use go-invoicer for **branded PDFs**; use QuickBooks for **books, tax, and payment recording**.

**Plan required:** `quickbooks` feature (**Pro** and above).

**Compare:** [Xero integration](xero-integration.md) · [Accounting overview](accounting-integrations.md) · [Stripe](stripe-integration.md)

---

## What you get

| Capability | Supported |
|------------|-----------|
| OAuth connect (per org, per **realm**) | ✓ |
| Import by QBO invoice `Id` | ✓ |
| Export local invoice → QBO draft | ✓ |
| Auto-create QBO **Customer** from buyer | ✓ |
| Branded PDF after import/export | ✓ |
| Sandbox companies | ✓ (`GO_INVOICER_QBO_SANDBOX=true`) |
| QBO webhooks → go-invoicer | ✗ (re-import or automate externally) |
| Attach PDF to QBO via API | ✗ (download from go-invoicer) |

---

## Prerequisites

### 1. Intuit developer app (operator / self-hosted)

Create an app in the [Intuit Developer Portal](https://developer.intuit.com/) and add redirect URI:

```text
https://<api-host>/v1/oauth/quickbooks/callback
```

| Variable | Purpose |
|----------|---------|
| `GO_INVOICER_QBO_CLIENT_ID` | Intuit client ID |
| `GO_INVOICER_QBO_CLIENT_SECRET` | Intuit client secret |
| `GO_INVOICER_QBO_SANDBOX` | `true` → sandbox API (`sandbox-quickbooks.api.intuit.com`) |
| `GO_INVOICER_PUBLIC_API_URL` | Callback URL base |

**Scope:** `com.intuit.quickbooks.accounting`

### 2. Organization subscription

`GET /v1/billing/status` → `subscription.features` includes `"quickbooks"`.

---

## Connect QuickBooks (dashboard)

1. **Owner** or **admin** → **Settings → Integrations → QuickBooks** → **Connect**.
2. Sign in to Intuit and select the **company (realm)**.
3. Callback stores tokens; redirect `quickbooks=connected`.

### API

```http
GET /v1/oauth/quickbooks/connect
Authorization: Bearer <admin-token>
```

Status:

```http
GET /v1/org/integrations/quickbooks
```

```json
{
  "quickbooks": {
    "configured": true,
    "connectable": true,
    "tenant_id": "<realmId>",
    "tenant_name": ""
  }
}
```

`tenant_id` is the QuickBooks **realm ID** from the OAuth callback (`realmId` query param).

Disconnect:

```http
DELETE /v1/org/integrations/quickbooks
Authorization: Bearer <admin-token>
```

---

## Sandbox vs production

| Mode | API base | When |
|------|----------|------|
| Production | `quickbooks.api.intuit.com` | `GO_INVOICER_QBO_SANDBOX` unset or `false` |
| Sandbox | `sandbox-quickbooks.api.intuit.com` | `GO_INVOICER_QBO_SANDBOX=true` |

Use a **sandbox company** in Intuit when testing. Connect from Settings after the API is in sandbox mode.

!!! tip
    Sandbox and production use **separate** Intuit apps and credentials. Match `QBO_SANDBOX` to the app type you registered.

---

## Import (QuickBooks → go-invoicer)

### Find the invoice Id

In QuickBooks: invoice detail → URL or API. Use numeric/string **`Id`**, not only `DocNumber`.

### API

```http
POST /v1/integrations/quickbooks/import
Authorization: Bearer <token>
Content-Type: application/json

{ "quickbooks_invoice_id": "145" }
```

```json
{
  "invoice": { "id": "…", "number": "1042", "status": "synced", … },
  "message": "imported; PDF render queued"
}
```

Maps customer, lines, dates, currency; upserts `external_source: quickbooks`; queues PDF.

### Dashboard

**Settings → Integrations → QuickBooks** — enter Id and **Import**.

---

## Export (go-invoicer → QuickBooks)

```http
POST /v1/invoices/:id/sync/quickbooks
Authorization: Bearer <token>
```

**Flow:**

1. Load local invoice.
2. Create QBO invoice, or update existing (`metadata.quickbooks_invoice_id` + `SyncToken` from QBO).
3. Find/create **Customer** by email or display name.
4. Link record; queue PDF.

Record payments and send invoices from **QuickBooks** (or link to Stripe/other PSP outside go-invoicer).

### Buyer requirements

Same as Xero: buyer **email** (lookup) or **name** (create customer).

---

## Branded PDFs

Uses org rendering defaults (template, locale, IBAN, optional logo on Business+). Job runs asynchronously — download from invoice list when ready.

QBO does not receive automatic PDF upload from go-invoicer; the invoice **record** is linked for your books.

---

## Linking & metadata

| Field | Value |
|-------|--------|
| `external_source` | `quickbooks` |
| `external_id` | QBO invoice `Id` |
| `invoice.metadata.quickbooks_invoice_id` | Same Id |
| Local `status` on import | `synced` |

Updates require current QBO `SyncToken` (handled on export update path).

---

## Payment & status tracking

| Goal | Approach |
|------|----------|
| Paid in go-invoicer | `PUT /v1/invoices/:id` with `"status": "paid"` |
| Refresh from QBO | `POST /v1/integrations/quickbooks/import` again |
| Bank feed / QBO payments | Stay in QuickBooks; optional re-import for amounts |
| Archive | [Compliance ZIP](paid-features.md#compliance-zip-compliance--pro) |

No QBO webhook handler in the platform today.

---

## Recommended workflows

### A — Invoice exists in QuickBooks

Import → branded PDF → deliver to customer.

### B — Invoice authored in go-invoicer

Create locally → `sync/quickbooks` → finalize in QBO.

### C — US company + Stripe

Stripe for checkout; QBO for books — enable both on Pro+.

---

## Troubleshooting

| Symptom | Fix |
|---------|-----|
| **402** | Need `quickbooks` on plan |
| `connect QuickBooks in Settings first` | Finish OAuth |
| `OAuth not configured on server` | Set QBO client ID/secret |
| `AuthenticationFailed` (sandbox) | Set `GO_INVOICER_QBO_SANDBOX=true`; use sandbox company |
| `Stale Object Error` / SyncToken | Re-import from QBO then export again |
| Empty customer | Set buyer name/email |
| Wrong company | Disconnect; reconnect and pick correct realm |

---

## OSS library

Package: `github.com/wiederin/go-invoicer/integrations/quickbooks`

Hosted: `platform/api/internal/sync/quickbooks/service.go`.

---

## API reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/v1/oauth/quickbooks/connect` | Admin | Start OAuth |
| GET | `/v1/oauth/quickbooks/callback` | Public | OAuth callback (`realmId`) |
| GET | `/v1/org/integrations/quickbooks` | Member+ | Status |
| DELETE | `/v1/org/integrations/quickbooks` | Admin | Disconnect |
| POST | `/v1/integrations/quickbooks/import` | Write + feature | Import by Id |
| POST | `/v1/invoices/:id/sync/quickbooks` | Write + feature | Export draft |

---

## Related

- [Xero integration](xero-integration.md)
- [Paid features](paid-features.md)

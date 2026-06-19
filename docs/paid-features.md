# Paid features — integration guide

The **hosted platform** adds multi-tenant orgs, async PDF jobs, plan gates, and integrations on top of the free OSS library. This guide shows how to use each paid capability from the dashboard and API.

For plan prices and the feature matrix, see [Hosted platform — plans & paid features](hosted-platform.md).

## Quick reference

| Feature | Plan | Primary docs |
|---------|------|----------------|
| Hosted PDF queue + S3 | Free+ | [Getting started](getting-started.md), [Templates](templates.md) |
| Stripe import / export / webhooks | Starter+ | [Stripe integration](stripe-integration.md) |
| Replace Stripe Invoicing | Starter+ | [Replace Stripe Invoicing](replace-stripe-invoicing.md) |
| Xero / QuickBooks | Pro+ | [Xero](xero-integration.md) · [QuickBooks](quickbooks-integration.md) · [overview](accounting-integrations.md) |
| Compliance ZIP | Pro+ | Below |
| E-invoice XML / ZUGFeRD | Business+ | Below |
| Logo & accent on PDFs | Business+ | [Hosted platform](hosted-platform.md#custom-branding-branding--business) |
| Approval workflows | Enterprise | [Hosted platform](hosted-platform.md#approval-workflows-approvals--enterprise) |
| White-label portal | Enterprise | [Hosted platform](hosted-platform.md#white-label-portal-white_label--enterprise) |
| SSO (OIDC) | Enterprise | [Hosted platform](hosted-platform.md#enterprise-sso-sso--enterprise) |

All plans also include **line columns**, **invoice numbering**, **schedules**, **inbound hook**, and **generate API** — see [Automation](#automation-all-plans) below.

---

## Subscribing to go-invoicer (platform billing)

Your organization pays for go-invoicer through **Stripe Checkout** (platform Stripe account — separate from per-org invoice Stripe).

```http
GET /v1/billing/status
Authorization: Bearer <token>
```

Response includes `subscription.plan`, `subscription.features`, and monthly invoice usage.

```http
POST /v1/billing/checkout
Authorization: Bearer <token>
Content-Type: application/json

{ "plan_slug": "pro" }
```

→ `{ "checkout_url": "https://checkout.stripe.com/…" }`

Owners and admins manage checkout from the dashboard **Billing** page. When you exceed the included invoice count, `POST /v1/invoices` returns **402** until you upgrade or the billing period resets.

**Platform admin** (Survih operators): edit tier limits and feature slugs at `/admin/billing` or via `PUT /v1/admin/billing/plans/:id`.

---

## Stripe (`stripe`) — Starter+

Connect **your** Stripe account in **Settings → Integrations** (secret key + webhook signing secret).

| Action | API |
|--------|-----|
| Import Stripe invoice | `POST /v1/integrations/stripe/import` `{ "stripe_invoice_id": "in_…" }` |
| Export local → Stripe draft | `POST /v1/invoices/:id/sync/stripe` |
| Finalize & email on Stripe | `POST /v1/invoices/:id/stripe/bill` |
| Webhook (Stripe → you) | `POST /v1/webhooks/stripe/:org_id` |

**Deep dive:** [Stripe integration](stripe-integration.md) · [Replace Stripe Invoicing](replace-stripe-invoicing.md)

---

## Xero & QuickBooks (`xero`, `quickbooks`) — Pro+

OAuth connect from **Settings → Integrations**.

| System | Import | Export |
|--------|--------|--------|
| Xero | `POST /v1/integrations/xero/import` | `POST /v1/invoices/:id/sync/xero` |
| QuickBooks | `POST /v1/integrations/quickbooks/import` | `POST /v1/invoices/:id/sync/quickbooks` |

After import, the platform queues a branded PDF using your org template defaults (template, locale, IBAN).

**Guides:** [Xero integration](xero-integration.md) · [QuickBooks integration](quickbooks-integration.md) · [Overview](accounting-integrations.md)

---

## Compliance (`compliance`) — Pro+

Bulk archive for finance / audit (admin or owner):

```http
GET /v1/compliance/export
Authorization: Bearer <token>
```

ZIP includes invoices (HTML/PDF where available), e-invoice XML variants, and audit trail entries.

Requires the `compliance` feature slug and RBAC `PermComplianceExport`.

---

## E-invoice (`einvoice`) — Business+

Per-invoice downloads (dashboard or API):

| Format | Path |
|--------|------|
| XRechnung UBL | `GET /v1/invoices/:id/compliance/xrechnung` |
| PEPPOL BIS | `GET /v1/invoices/:id/compliance/peppol` |
| EN 16931 | `GET /v1/invoices/:id/compliance/en16931` |
| FatturaPA | `GET /v1/invoices/:id/compliance/fatturapa` |
| ZUGFeRD hybrid PDF | `GET /v1/invoices/:id/compliance/zugferd` |

Optional `?sign=1&key_id=…` for XML signing hooks (enterprise compliance module).

---

## Custom branding (`branding`) — Business+

```http
PATCH /v1/org/rendering
Authorization: Bearer <token>

{
  "logo_url": "https://cdn.example.com/logo.png",
  "accent_color": "#0ea5e9"
}
```

Without the `branding` feature, this returns **402**. Line item columns and custom headers work on **all** plans.

---

## Automation (all plans)

| Capability | Endpoint |
|------------|----------|
| Inbound hook (no API key) | `POST /v1/hooks/{token}/invoices` |
| Generate + optional PDF / Stripe | `POST /v1/invoices/generate` |
| Cron schedules | `POST /v1/invoice-schedules` |
| Next invoice number | `GET /v1/invoices/next-number` |

Example generate body:

```json
{
  "invoice": {
    "seller": { "name": "Acme AG" },
    "buyer": { "name": "Customer GmbH" },
    "line_items": [
      {
        "description": "SaaS subscription",
        "quantity": 1,
        "unit_price": { "amount": 5900, "currency": "CHF" },
        "tax_rate": 0.081
      }
    ]
  },
  "status": "draft",
  "render_pdf": true,
  "sync_stripe": true
}
```

Omit `invoice.number` to auto-assign from your numbering pattern in Settings.

---

## Entitlements and errors

| HTTP | Meaning |
|------|---------|
| **402** | Missing plan feature or invoice quota exceeded |
| **403** | RBAC — wrong role (e.g. member changing integrations) |
| **401** | Invalid or expired Bearer token |

Check features before building integrations:

```http
GET /v1/billing/status
→ subscription.features: ["hosted_pdf", "stripe", …]
```

Integration buttons on the invoice list require **both** the plan feature and `GET /v1/integrations` → `status: "ready"`.

---

## Open source vs hosted

| Capability | OSS only | Hosted + plan |
|------------|----------|---------------|
| Invoice model, templates, local PDF | ✓ | ✓ |
| Stripe mapper & webhook handler code | ✓ | ✓ (with credentials + `stripe` feature) |
| Async PDF, S3, multi-tenant | | ✓ |
| Xero / QBO / compliance / e-invoice | Enterprise code in monorepo | ✓ (plan-gated) |
| Subscribe to go-invoicer | | ✓ (Stripe Checkout) |

Library-only Stripe usage: [Stripe (OSS)](stripe.md).

---

## Related

- [Stripe integration](stripe-integration.md) — connect, webhooks, export pipeline
- [Replace Stripe Invoicing](replace-stripe-invoicing.md) — use go-invoicer for documents, Stripe for payments
- [Xero integration](xero-integration.md) · [QuickBooks integration](quickbooks-integration.md)
- [Hosted platform](hosted-platform.md) — plan matrix and RBAC

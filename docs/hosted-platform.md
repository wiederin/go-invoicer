# Hosted platform — plans & paid features

The **open-source library** (`github.com/wiederin/go-invoicer`) is free under Apache 2.0. The **hosted platform** (dashboard + API at `platform/`) adds organizations, async PDF jobs, integrations, compliance exports, and billing.

This page describes what each subscription tier unlocks, how gates work in the API, and what remains free on every plan.

## Plans (CHF / month)

| Plan | Price | Invoices / month | Overage |
|------|------:|------------------:|---------|
| **Free** | 0 | 25 | — |
| **Starter** | 29 | 100 | +5 invoices @ CHF 2 |
| **Pro** | 59 | 500 | +25 @ CHF 1 |
| **Business** | 129 | 2,000 | +100 @ CHF 1 |
| **Enterprise** | 199 | Unlimited | — |

**New organizations** register at `/signup` with **email and password** (recommended). The API still issues an owner API key in the database for automation, but it is only returned when signing up without email. Sign in at `/login` with email/password or an API key. Disable registration with `SIGNUP_ENABLED=false` on the API.

Subscribe from the dashboard **Billing** page (`POST /v1/billing/checkout`). Owners and admins manage checkout; members use features allowed by RBAC.

**Platform administrators** configure tier limits and feature slugs under **Plans** (`/admin/billing`): max invoices per month, overage allowance/price, and checkboxes for each feature gate. Changes persist in `billing_plans` and apply immediately to orgs on that plan. API: `GET /v1/admin/billing/features`, `PUT /v1/admin/billing/plans/:id`.

Check current plan and usage:

```http
GET /v1/billing/status
Authorization: Bearer <token>
```

When monthly invoice volume exceeds the included cap (and overage is configured), the API returns **402 Payment Required** on `POST /v1/invoices` until you upgrade or the next billing period starts.

## Feature matrix

Feature slugs are stored on `billing_plans.features` and returned in `subscription.features` from billing status.

| Feature slug | Label | Free | Starter | Pro | Business | Enterprise |
|--------------|-------|:----:|:-------:|:---:|:--------:|:----------:|
| `hosted_pdf` | Hosted PDF rendering | ✓ | ✓ | ✓ | ✓ | ✓ |
| `stripe` | Stripe import / export / webhooks | | ✓ | ✓ | ✓ | ✓ |
| `xero` | Xero OAuth import & export | | | ✓ | ✓ | ✓ |
| `quickbooks` | QuickBooks OAuth import & export | | | ✓ | ✓ | ✓ |
| `compliance` | Compliance ZIP export | | | ✓ | ✓ | ✓ |
| `audit` | Audit log API (listed on plan) | | | ✓ | ✓ | ✓ |
| `einvoice` | Per-invoice e-invoice XML / ZUGFeRD | | | | ✓ | ✓ |
| `branding` | Logo & accent on PDFs | | | | ✓ | ✓ |
| `approvals` | Invoice approval workflows | | | | | ✓ |
| `white_label` | White-label portal branding | | | | | ✓ |
| `sso` | Enterprise SSO (OIDC) | | | | | ✓ |

!!! note "Included on all plans"
    These capabilities are **not** paywalled:

    - **Line item columns** — show/hide quantity, unit price, tax, line total; optional custom column headers (`PATCH /v1/org/rendering` → `line_columns`).
    - **Invoice numbering** — prefix, pattern, auto-assign on create (`GET /v1/invoices/next-number`).
    - **Scheduled invoices** — cron schedules (`/v1/invoice-schedules`, Settings → Scheduled invoices). Disable the worker with `INVOICE_SCHEDULER_ENABLED=false` on the API.
    - **Inbound hook** — per-org secret URL (`POST /v1/hooks/{token}/invoices`) for ERP/Zapier without an API key; rotate in Settings → Automation.
    - **Invoice generate API** — `POST /v1/invoices/generate` with optional `render_pdf` and `sync_stripe` (same body as inbound hook; requires Bearer token).
    - **Template & locale defaults**, block layout, usage history charts (within invoice quota).

## Paid features in detail

### Stripe (`stripe`) — Starter+

- Connect per-org secret key and webhook signing secret in **Settings**.
- `POST /v1/integrations/stripe/import` — import by Stripe invoice id.
- `POST /v1/invoices/:id/sync/stripe` — export local invoice to Stripe.
- `POST /v1/invoices/:id/stripe/bill` — finalize & send (when **Replace Stripe Invoicing** mode is `finalize_send`).
- Webhook: `POST /v1/webhooks/stripe/:org_id`.

Integration buttons on the invoice list appear when the integration is **configured** (`GET /v1/integrations` → `status: ready`) **and** your plan includes the matching feature slug. The API still returns **402** if either check fails.

### Xero & QuickBooks (`xero`, `quickbooks`) — Pro+

- OAuth connect from Settings (`GET /v1/oauth/xero/connect`, QuickBooks equivalent).
- Import by external id; export from invoice row actions.
- PDF render is queued automatically using org template defaults.

### Compliance ZIP (`compliance`) — Pro+

Bulk export for finance / archive:

```http
GET /v1/compliance/export
```

ZIP includes invoices (HTML/PDF where available), e-invoice XML variants, and audit trail. Requires **admin/owner** RBAC.

### E-invoice (`einvoice`) — Business+

Per-invoice downloads (dashboard buttons or API):

| Format | Path |
|--------|------|
| XRechnung UBL | `GET /v1/invoices/:id/compliance/xrechnung` |
| PEPPOL BIS | `…/peppol` |
| EN 16931 | `…/en16931` |
| FatturaPA | `…/fatturapa` |
| ZUGFeRD hybrid PDF | `GET /v1/invoices/:id/compliance/zugferd` |

Optional query `sign=1` and `key_id` for XML signing placeholders.

### Custom branding (`branding`) — Business+

- **Logo URL** and **accent color** on rendered PDFs/HTML.
- Without this feature, saving logo/accent via `PATCH /v1/org/rendering` returns **402**.

Line item columns and header labels do **not** require `branding`.

### Approval workflows (`approvals`) — Enterprise

- `POST /v1/invoices/:id/submit-approval`
- `GET /v1/approvals`, approve/reject endpoints.
- Optional Temporal worker for durable workflows (`docs/on-prem.md`).

### White-label portal (`white_label`) — Enterprise

- `GET/PATCH /v1/org/portal` — dashboard name, logo, accent, hide “Powered by”.
- Public shell: `GET /v1/portal/branding`.

### Enterprise SSO (`sso`) — Enterprise

- OIDC login: `GET /v1/auth/oidc/login`, callback, config.
- Configure issuer, client id/secret, redirect URL on the API container.

## RBAC vs plan features

**Plan features** gate product capabilities (402 if missing).

**RBAC roles** gate who can call an endpoint:

| Role | Typical access |
|------|----------------|
| Owner / Admin | Integrations, billing, audit logs, compliance export, approvals |
| Member | Create/read invoices, render PDFs, schedules |

Example: audit logs (`GET /v1/audit-logs`) require `PermAuditRead` (admin/owner), not a separate plan check in the handler — but Pro+ plans list `audit` in marketing/features.

## Self-hosted / on-prem

Run the platform API yourself (`docs/on-prem.md` in the monorepo). Billing plans and feature slugs still apply if you seed `billing_plans`; set org plan via admin or database. Disable automatic invoice schedules with:

```bash
export INVOICE_SCHEDULER_ENABLED=false
```

## Open source vs hosted

| Capability | OSS library | Hosted platform |
|------------|-------------|-----------------|
| Invoice model & validation | ✓ | ✓ |
| HTML templates & local PDF | ✓ | ✓ |
| Stripe mapper (code) | ✓ | ✓ (with plan + credentials) |
| Async PDF queue, S3 storage | | ✓ |
| Multi-tenant orgs & API keys | | ✓ |
| Xero / QBO / compliance / e-invoice | Enterprise modules in monorepo | ✓ (plan-gated) |

For library-only usage, see [Getting started](getting-started.md) and [Stripe](stripe.md).

# Accounting integrations (Xero & QuickBooks)

**Plans:** **Pro**, **Business**, and **Enterprise** include accounting sync (`xero` and `quickbooks` feature slugs).

Use these integrations when your **general ledger** lives in Xero or QuickBooks Online and you want **go-invoicer** for professional invoice PDFs, templates, compliance exports, and API automation.

---

## Choose your system

| | [Xero](xero-integration.md) | [QuickBooks Online](quickbooks-integration.md) |
|--|----------------------------|-----------------------------------------------|
| **Regions** | Common in UK, AU, NZ, EU | Common in US, CA |
| **Connect** | OAuth → Xero tenant | OAuth → Intuit realm |
| **Import ID** | Invoice UUID (`InvoiceID`) | Invoice `Id` |
| **Export creates** | ACCREC **DRAFT** | QBO draft invoice |
| **Contacts** | Xero Contact | QBO Customer |
| **Sandbox** | Xero demo org | `GO_INVOICER_QBO_SANDBOX=true` |

You can connect **both** on Pro+ if you operate multiple entities.

---

## Shared behavior

Both integrations follow the same platform pattern:

1. **OAuth** — admin connects in Settings; tokens stored per organization.
2. **Import** — pull invoice by external ID → local record + **PDF render queued**.
3. **Export** — push domain invoice → external draft → link `external_id` + PDF queue.
4. **Entitlements** — `402` without feature slug; **503** if not connected.
5. **Quota** — imports count toward monthly invoice limit.
6. **Seller on import** — mapped from your organization name (platform supplier config).

Neither integration receives **inbound webhooks** from the accounting provider today. Refresh data with a repeat **import** call or update local `status` via API.

**Deep dives:**

- [Xero integration guide](xero-integration.md) — setup, scopes, workflows, troubleshooting
- [QuickBooks integration guide](quickbooks-integration.md) — sandbox, realm ID, SyncToken, troubleshooting

---

## Quick start

### Xero

```bash
# 1. Admin: Settings → Connect Xero
# 2. Import
curl -X POST https://invoicer-api.survih.ch/v1/integrations/xero/import \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"xero_invoice_id":"<uuid>"}'

# 3. Export local invoice
curl -X POST https://invoicer-api.survih.ch/v1/invoices/$ID/sync/xero \
  -H "Authorization: Bearer $TOKEN"
```

### QuickBooks

```bash
curl -X POST …/v1/integrations/quickbooks/import \
  -d '{"quickbooks_invoice_id":"145"}'

curl -X POST …/v1/invoices/$ID/sync/quickbooks \
  -H "Authorization: Bearer $TOKEN"
```

---

## Stripe vs accounting

| | Stripe | Xero / QuickBooks |
|--|--------|-------------------|
| **Best for** | Card/SEPA billing, subscriptions | GL, AR, tax, accountant handoff |
| **PDF on provider** | Attach to Stripe invoice | Manual attach in Xero/QBO if needed |
| **Payment webhooks** | `invoice.paid` → `status: paid` | Re-import or manual status |
| **Plan** | Starter+ (`stripe`) | Pro+ (`xero` / `quickbooks`) |

Many teams use **Stripe + Xero** or **Stripe + QuickBooks**: [Stripe integration](stripe-integration.md) · [Replace Stripe Invoicing](replace-stripe-invoicing.md).

---

## Operator setup (self-hosted / CI)

Register OAuth apps and set API environment variables — see `docs/gitlab-deploy.md`:

| Variable | Provider |
|----------|----------|
| `GO_INVOICER_XERO_CLIENT_ID` / `SECRET` | Xero |
| `GO_INVOICER_QBO_CLIENT_ID` / `SECRET` | Intuit |
| `GO_INVOICER_QBO_SANDBOX` | QuickBooks sandbox API |

Redirect URIs:

```text
https://<api>/v1/oauth/xero/callback
https://<api>/v1/oauth/quickbooks/callback
```

---

## UI and entitlements

Invoice list **Xero** / **QBO** buttons require:

1. `GET /v1/integrations` → `status: "ready"`, and  
2. `GET /v1/billing/status` → matching feature slug.

---

## Related

- [Paid features guide](paid-features.md)
- [Stripe integration](stripe-integration.md)
- [Hosted platform](hosted-platform.md)

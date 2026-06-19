# Hosted vs OSS — which should I use?

go-invoicer ships as **Apache 2.0 libraries** (`github.com/wiederin/go-invoicer`) and an optional **hosted platform** (dashboard + API). Both share the same invoice domain model and HTML → PDF rendering pipeline.

## Quick decision

| You want… | Choose |
|-----------|--------|
| Embed invoicing in your Go app with full control | **OSS libraries** — render PDFs locally, own your data |
| A ready-made SaaS with orgs, billing, integrations | **Hosted platform** — sign up at the dashboard |
| Your own VPC / air-gapped / single tenant | **Self-hosted platform** — deploy `platform/` from this monorepo |
| Only PDF/HTML from code, no multi-tenant API | **OSS only** — no `platform/` required |

## OSS libraries (free, Apache 2.0)

**Best for:** developers building invoicing into a product, CI pipelines, or custom backends.

- Invoice schema, validation, totals in `invoice/`
- HTML templates + Chromium PDF in `render/` and `pdf/`
- Examples and templates under `oss/examples/` and `oss/templates/`
- Stripe helpers in `integrations/stripe/` (library only)

**You provide:** storage, auth, job queue, hosting, compliance workflows if needed.

See [Getting started](getting-started.md) and [Examples](examples.md).

## Hosted platform (subscription)

**Best for:** teams that want invoices, PDF jobs, accounting sync, and a UI without operating the stack.

Includes everything in OSS plus:

- Multi-tenant orgs, RBAC, API keys, email/password auth
- Async PDF rendering, S3 storage, usage metering
- Stripe / Xero / QuickBooks integrations (plan-gated)
- Compliance ZIP, e-invoice exports, approvals (higher tiers)
- Automation: generate API, inbound hooks, outbound webhooks

Plans and limits: [Hosted platform & plans](hosted-platform.md). API surface: [API reference](api-reference.md) (`GET /v1/openapi.yaml`).

## Self-hosted platform

**Best for:** regulated environments, private cloud, or “hosted features on our infra.”

Same `platform/` code as SaaS; you run PostgreSQL, Redis, API, and dashboard. See [Self-hosted deployment](self-hosted.md).

| | Hosted SaaS | Self-hosted | OSS only |
|--|-------------|-------------|----------|
| Ops burden | Low | High (you) | Medium (app-specific) |
| Multi-tenant | Yes | Your choice | N/A |
| Integrations UI | Yes | Yes (configure env) | Build yourself |
| Source license | Proprietary runtime | Proprietary runtime | Apache 2.0 libs |

## Can I mix?

Yes. Common patterns:

1. **OSS in production, hosted for sales demos** — same templates; export sample JSON from Studio.
2. **OSS core + hosted PDF API later** — start with local `pdf/`; migrate to `POST /v1/invoices/generate` when you need scale.
3. **Self-hosted for prod, OSS libs in microservices** — platform stores invoices; workers call shared render packages.

## Related docs

- [Hosted platform & plans](hosted-platform.md)
- [Self-hosted deployment](self-hosted.md)
- [Paid features guide](paid-features.md)
- [Invoice customization](invoice-customization.md)

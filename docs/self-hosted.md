# Self-hosted deployment

Run the go-invoicer **platform** on your own infrastructure (single-tenant). The OSS Go libraries remain Apache 2.0; `platform/` and `enterprise/` are proprietary but included when you deploy from this monorepo.

## Architecture

| Component | Role |
|-----------|------|
| **PostgreSQL 16+** | Organizations, invoices, jobs, billing, webhooks |
| **Redis** (recommended) | Distributed PDF job queue (`REDIS_URL`) |
| **Platform API** | Gin server + Chromium PDF worker in one container |
| **Dashboard** | Static Vite build (nginx or S3 + CloudFront) |
| **S3** (optional) | PDF object storage (`S3_PDF_BUCKET`) |

Migrations run automatically when the API starts with `DATABASE_URL` set.

## Quick start — Docker Compose

```bash
docker compose --profile platform up -d --build
```

| URL | Service |
|-----|---------|
| http://localhost:8088 | Dashboard (proxies `/api` to API) |
| http://localhost:8080 | API |
| Sign-in | API key `gi_local_dev_key` (from compose env) |

Stop: `docker compose --profile platform down`

## Development — native processes

```bash
docker compose up -d postgres redis

export DATABASE_URL=postgres://invoicer:invoicer@localhost:5432/invoicer?sslmode=disable
export AUTH_API_KEYS=gi_dev_key
export AUTH_JWT_SECRET=dev-secret-change-me
export SIGNUP_ENABLED=true
export PDF_STORAGE_DIR=/tmp/invoicer-pdfs
export GO_INVOICER_PUBLIC_API_URL=http://localhost:8080
export GO_INVOICER_DASHBOARD_URL=http://localhost:5173

cd platform/api && go run ./cmd/server
```

Dashboard:

```bash
cd platform/dashboard
export VITE_API_URL=http://localhost:8080
npm install && npm run dev
```

## Essential environment variables

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | Postgres connection string |
| `AUTH_API_KEYS` | Comma-separated bootstrap API keys |
| `AUTH_JWT_SECRET` | JWT signing secret |
| `GO_INVOICER_PUBLIC_API_URL` | Public API base (webhooks, OAuth, inbound hooks) |
| `GO_INVOICER_DASHBOARD_URL` | Dashboard URL (Stripe Checkout return, OIDC) |
| `REDIS_URL` | Job queue broker |
| `S3_PDF_BUCKET` | Store completed PDFs in S3 |
| `SIGNUP_ENABLED` | Allow public registration (`true`/`false`) |

## Production checklist

- Build API image from `platform/api/Dockerfile` (includes Chromium).
- Terminate TLS at your load balancer; set CORS via `AUTH_ALLOWED_ORIGINS`.
- Configure **OIDC** for enterprise SSO: `OIDC_ISSUER_URL`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_REDIRECT_URL`.
- Disable built-in scheduler if you only want manual runs: `INVOICE_SCHEDULER_ENABLED=false`.
- Set deployment-wide portal branding with `GO_INVOICER_PORTAL_*` env vars (optional).

## What you get self-hosted

- Full invoice CRUD, PDF render queue, and org rendering profile (templates, logo, columns, numbering).
- Inbound and outbound webhooks, cron schedules, and `POST /v1/invoices/generate`.
- Stripe / Xero / QuickBooks integrations when you configure credentials.
- Enterprise modules (compliance exports, OIDC) when built with `enterprise/` in the image.

You operate billing to your own customers — platform Stripe billing for SaaS plans is optional.

## Staging pattern

This repository uses a separate AWS stack for the `staging` branch. See `docs/gitlab-staging.md` in the monorepo for CI deploy jobs.

## More detail

The monorepo `docs/on-prem.md` includes metered billing, usage exports, e-invoice endpoints, and Temporal approval workers.

# Outbound webhooks

Receive **invoice lifecycle events** at your own HTTPS endpoint. Configure per organization in **Settings → Automation** or via the API.

This is separate from the **inbound hook**, which creates invoices when *you* POST to go-invoicer.

## Event types

| Event | When fired |
|-------|------------|
| `invoice.created` | Invoice stored (API create, generate, inbound hook, schedule) |
| `invoice.updated` | Invoice updated via `PUT /v1/invoices/:id` |
| `invoice.pdf.ready` | Async PDF job completed successfully |

Subscribe to a subset with the `events` array on configuration.

## Configure

```http
GET /v1/org/outbound-webhook
PATCH /v1/org/outbound-webhook
POST /v1/org/outbound-webhook/test
```

Example:

```json
PATCH /v1/org/outbound-webhook
{
  "enabled": true,
  "url": "https://hooks.example.com/invoicer",
  "secret": "whsec_your_signing_secret",
  "events": ["invoice.created", "invoice.pdf.ready"]
}
```

- **url** — must be `https://`
- **secret** — optional HMAC key; only returned on write when you set a new value
- **test** — queues a sample `invoice.created` payload (requires enabled + URL)

## Delivery format

HTTP `POST` with `Content-Type: application/json`:

```json
{
  "id": "evt_…",
  "type": "invoice.created",
  "created_at": "2026-05-24T12:00:00Z",
  "organization_id": "org_…",
  "data": {
    "invoice": {
      "id": "inv_…",
      "number": "INV-2026-00042",
      "status": "draft"
    },
    "job": {
      "id": "job_…",
      "status": "completed"
    }
  }
}
```

The `job` object is present for `invoice.pdf.ready`.

## Verifying signatures

When a secret is configured, requests include:

| Header | Meaning |
|--------|---------|
| `X-Go-Invoicer-Event` | Event type |
| `X-Go-Invoicer-Delivery` | Unique delivery id |
| `X-Go-Invoicer-Timestamp` | Unix seconds |
| `X-Go-Invoicer-Signature` | `sha256=` + HMAC-SHA256 of `{timestamp}.{raw_body}` |

Example (pseudo-code):

```text
expected = HMAC_SHA256(secret, timestamp + "." + body)
compare(header "X-Go-Invoicer-Signature", "sha256=" + hex(expected))
```

Reject requests with timestamps too far in the past to limit replay.

## Reliability

Deliveries run **asynchronously** after the API transaction completes. Failed requests are retried up to **4 times** with backoff **0s → 1s → 4s → 16s** when:

- The connection fails
- The server returns **5xx**, **408**, or **429**

Most **4xx** responses are not retried. Design your endpoint to be **idempotent** using `X-Go-Invoicer-Delivery` (same id across retries for one event).

## Delivery log

```http
GET /v1/org/outbound-webhook/deliveries?limit=50
```

Returns the last deliveries for your organization (up to 100 stored). Each row includes event type, HTTP status, attempt count, success flag, and error message. The dashboard **Settings → Automation** section shows the same log.

## Permissions

Requires `invoice:read` for GET and `invoice:write` for PATCH/TEST (same as inbound hook management).

## Related

- [Invoice customization](invoice-customization.md) — defaults applied before events fire
- [Self-hosted deployment](self-hosted.md) — set `GO_INVOICER_PUBLIC_API_URL` for correct hook URLs on inbound; outbound uses your URL

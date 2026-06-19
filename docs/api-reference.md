# API reference

The hosted platform REST API is documented via an OpenAPI spec.

## OpenAPI spec download

Download the YAML spec (no authentication required):

```http
GET /v1/openapi.yaml
```

## Idempotency and webhooks

- `Idempotency-Key` is supported on:
  - `POST /v1/invoices`
  - `POST /v1/invoices/generate`
- Webhook events:
  - inbound: `POST /v1/hooks/:token/invoices`
  - outbound: configure `GET/PATCH /v1/org/outbound-webhook` and subscribe to
    `invoice.created`, `invoice.updated`, and `invoice.pdf.ready`


# Buyer invoice portal

Let customers view an invoice and download the PDF without a dashboard login.

## Create a share link

```http
POST /v1/invoices/{id}/share-link
Authorization: Bearer <token>
Content-Type: application/json

{ "expires_in_days": 30 }
```

Response includes `url` (dashboard route) and `token`.

## Public endpoints (no auth)

```http
GET /v1/buyer/invoices/{token}
GET /v1/buyer/invoices/{token}/pdf
```

The dashboard page is:

```
https://<dashboard-host>/buyer/{token}
```

Set `GO_INVOICER_DASHBOARD_URL` on the API so generated links use the correct host.

## Branding

Buyer pages use the organization **portal** profile (`GET/PATCH /v1/org/portal`). White-label branding requires the **white_label** plan feature.

## Dashboard

Click **Share** on an invoice row to copy a 30-day link to the clipboard.

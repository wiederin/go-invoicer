# Email delivery

Send completed invoice PDFs to customers from the hosted platform.

## API

```http
POST /v1/invoices/{id}/send-email
Authorization: Bearer <token>
Content-Type: application/json

{
  "to": "buyer@example.com",
  "subject": "Invoice INV-2026-0042",
  "body": "Thank you for your business.",
  "include_link": true,
  "expires_in_days": 30
}
```

Requirements:

- A **completed PDF render job** for the invoice (`POST /v1/invoices/:id/render-pdf` first).
- Recipient `to` or `invoice.buyer.email` on the invoice.

Delivery attempts are logged at `GET /v1/org/email-deliveries`.

## Configuration

| Variable | Description |
|----------|-------------|
| `EMAIL_PROVIDER` | `log` (default) or `smtp` |
| `SMTP_HOST`, `SMTP_PORT` | SMTP server |
| `SMTP_USER`, `SMTP_PASSWORD` | Optional auth |
| `SMTP_FROM` or `EMAIL_FROM` | From address |

With `EMAIL_PROVIDER=log`, messages are written to API logs (staging/dev). Use SMTP for Postmark, SES SMTP, Mailhog, etc.

## Dashboard

On **Invoices**, use **Email** after PDF generation. **Share** copies a buyer portal link (can be included automatically when sending email with `include_link: true`).

See [Buyer portal](buyer-portal.md).

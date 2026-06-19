# Automation example

Sample payloads for the hosted **generate API** and **inbound hook**, including layout and custom line columns.

Set `API_URL` and `API_KEY` (or use a JWT from `POST /v1/auth/login`).

## Generate invoice + PDF

```bash
export API_URL=https://staging-invoicer-api.survih.ch
export API_KEY=gi_your_key

curl -sS -X POST "$API_URL/v1/invoices/generate" \
  -H "Authorization: Bearer $(curl -sS -X POST "$API_URL/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"api_key\":\"$API_KEY\"}" | jq -r .access_token)" \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: demo-generate-$(date +%s)" \
  -d @- <<'EOF'
{
  "render_pdf": true,
  "template": "modern",
  "locale": "de",
  "layout_blocks": ["header", "parties", "line_items", "totals", "notes"],
  "line_columns": {
    "quantity": true,
    "unit_price": true,
    "tax": true,
    "line_total": true,
    "labels": {
      "description": "Leistung",
      "quantity": "Menge",
      "unit_price": "Einzelpreis",
      "tax": "MwSt.",
      "line_total": "Betrag"
    }
  },
  "invoice": {
    "kind": "invoice",
    "seller": { "name": "Acme GmbH" },
    "buyer": { "name": "Customer AG" },
    "line_items": [
      {
        "description": "Monthly platform fee",
        "quantity": 1,
        "unit_price": { "amount": 5900, "currency": "CHF" },
        "tax_rate": 0.081
      }
    ]
  }
}
EOF
```

## Inbound hook (token URL)

Configure the hook in **Settings → Automation**. Then:

```bash
curl -sS -X POST "$API_URL/v1/hooks/<your-token>/invoices" \
  -H "Content-Type: application/json" \
  -d '{"invoice":{"seller":{"name":"Hook Seller"},"buyer":{"name":"Hook Buyer"},"line_items":[{"description":"Webhook line","quantity":1,"unit_price":{"amount":1000,"currency":"CHF"},"tax_rate":0}]}}'
```

## Quote via generate

Use `"kind": "quote"` on the nested `invoice` object. Numbers preview with `GET /v1/invoices/next-number?kind=quote`.

See [Document types](../docs/document-types.md) and [API reference](../docs/api-reference.md).

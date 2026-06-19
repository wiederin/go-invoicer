# Document types (invoice, quote, credit note)

Commercial documents share the same `invoice` JSON shape. Set **`kind`** to change titles on PDFs and auto-numbering prefixes on the hosted platform.

## Kinds

| `kind` | Use | PDF title (EN) | Auto number prefix (hosted) |
|--------|-----|----------------|----------------------------|
| *(omit)* or `invoice` | Tax invoice | Invoice | Org prefix (default `INV-`) |
| `quote` | Non-binding offer | Quote | `QUO-` |
| `proforma` | Proforma invoice | Proforma invoice | `PRO-` |
| `receipt` | Payment receipt | Receipt | `REC-` |

## Credit notes

Credit notes must reference the source invoice:

```json
{
  "kind": "credit_note",
  "related_number": "INV-2026-0042",
  "number": "CN-2026-0007",
  ...
}
```

Validation returns an error if `related_number` is missing.

## API examples

```http
POST /v1/invoices
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "draft",
  "invoice": {
    "kind": "quote",
    "seller": { "name": "Acme GmbH" },
    "buyer": { "name": "Client AG" },
    "line_items": [
      {
        "description": "Discovery workshop",
        "quantity": 1,
        "unit_price": { "amount": 250000, "currency": "CHF" },
        "tax_rate": 0.081
      }
    ]
  }
}
```

Preview the next number for a kind:

```http
GET /v1/invoices/next-number?kind=quote
```

## Convert quote to invoice

When a quote is accepted, convert it in place (same record id) with a new invoice number:

```http
POST /v1/invoices/{id}/convert-to-invoice
Authorization: Bearer <token>
```

The API sets `kind` to `invoice`, allocates the next invoice number, stores `converted_from_quote` in metadata, and sets status to `draft`.

## OSS library

```go
inv, err := invoice.NewBuilder().
    Number("QUO-2026-001").
    Kind(invoice.DocumentQuote).
    Seller(seller).
    Buyer(buyer).
    AddLine(line).
    Build()
```

See `oss/examples/quote` and template rendering (titles use localized labels).

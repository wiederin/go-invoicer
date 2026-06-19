# Settlement reports

Summarize all invoices **issued** in a calendar month and export a PDF table.

## API

```http
GET /v1/reports/settlement?year=2026&month=5
Authorization: Bearer <token>
```

Or explicit range (end date exclusive):

```http
GET /v1/reports/settlement?from=2026-05-01&to=2026-06-01
```

```http
GET /v1/reports/settlement/pdf?year=2026&month=5
```

Returns `application/pdf` with one row per invoice (number, buyer, type, status, amount) and a period total. All amounts must share the same currency.

## OSS

- Aggregation: `github.com/wiederin/go-invoicer/settlement`
- HTML template: `oss/templates/settlement/settlement.html`
- Render: `render.Engine.RenderSettlement`

## Receipts

Payment receipts use document kind `receipt` on the invoice JSON (`kind: "receipt"`). They render with the same templates as invoices and auto-number with prefix `REC-` on the hosted platform.

See [Document types](document-types.md).

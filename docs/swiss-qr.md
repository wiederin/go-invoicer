# Swiss QR-bill

Generate Swiss QR-bill payloads and render invoices with an embedded payment QR code.

## Quick start

```go
html, err := engine.RenderSwiss(inv, "CH93 0076 2011 6238 5295 7")
```

Or build the payload directly:

```go
bill := render.BillFromInvoice(inv, iban)
bill.Reference = "21000000000313947143031766" // optional 27-digit QRR
payload, err := bill.Payload()
```

## QRR reference

When `Bill.Reference` is set, it must be a valid **27-digit QR reference** (mod10 recursive check digit). Use `swiss.ValidateQRR` before assigning.

## Example

```bash
go run ./examples/swiss_qr
go run ./examples/swiss_qr -pdf -out ./out
```

## Compliance note

This library implements the SPC QR payload format and rendering helpers. Full Swiss QR-bill compliance (layout, creditor identification, etc.) may require additional validation for production use — consult [SIX Swiss Payment Standards](https://www.six-group.com/en/products-services/banking-services/payment-standardization.html).

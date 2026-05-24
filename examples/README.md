# Examples

| Example | Command | Output |
|---------|---------|--------|
| `basic/` | `go run ./examples/basic` | `invoice.html` |
| `basic/` (PDF) | `go run ./examples/basic -pdf` | `invoice.html` + `invoice.pdf` |
| `with_template/` | `go run ./examples/with_template` | `invoice-custom.html` |
| `swiss_qr/` | `go run ./examples/swiss_qr` | `invoice-swiss.html` |
| `swiss_qr/` (PDF) | `go run ./examples/swiss_qr -pdf` | + `invoice-swiss.pdf` |
| `stripe_sync/` | `STRIPE_SECRET_KEY=sk_... go run ./examples/stripe_sync -id in_xxx` | `stripe-invoice.html` |
| `samples/` | `go run ./examples/samples -out ./examples/samples` | All template PDFs (see [samples/README.md](samples/README.md)) |

Pre-built PDFs: [samples/](samples/) (`invoice-default.pdf`, `invoice-swiss-qr.pdf`, …).

Run from the module root (monorepo: `go run ./oss/examples/basic`).

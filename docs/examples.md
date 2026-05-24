# Examples

| Example | Command |
|---------|---------|
| Basic HTML/PDF | `go run ./examples/basic` |
| Custom template | `go run ./examples/with_template` |
| Swiss QR-bill | `go run ./examples/swiss_qr` |
| All sample PDFs | `go run ./examples/samples -out ./examples/samples` |

## Sample PDFs (no Chromium required)

Download or preview the committed examples in [`examples/samples/`](../examples/samples/):

- [invoice-default.pdf](../examples/samples/invoice-default.pdf)
- [invoice-minimal.pdf](../examples/samples/invoice-minimal.pdf)
- [invoice-multilingual-de.pdf](../examples/samples/invoice-multilingual-de.pdf)
- [invoice-swiss-qr.pdf](../examples/samples/invoice-swiss-qr.pdf)

Run from the module root (`oss/` in the monorepo, or repo root via `go.work`).

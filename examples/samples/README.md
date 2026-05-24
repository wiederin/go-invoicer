# Sample invoice PDFs

Pre-generated examples of every built-in template (same invoice data, different layouts).

| File | Template |
|------|----------|
| [invoice-default.pdf](invoice-default.pdf) | Default |
| [invoice-minimal.pdf](invoice-minimal.pdf) | Minimal |
| [invoice-multilingual-en.pdf](invoice-multilingual-en.pdf) | Multilingual (English) |
| [invoice-multilingual-de.pdf](invoice-multilingual-de.pdf) | Multilingual (German) |
| [invoice-multilingual-fr.pdf](invoice-multilingual-fr.pdf) | Multilingual (French) |
| [invoice-swiss-qr.pdf](invoice-swiss-qr.pdf) | Swiss QR-bill |

## Regenerate (requires Chromium)

From the `oss/` module root:

```bash
go run ./examples/samples -out ./examples/samples
```

These files are committed so GitHub and docs can link to them without running Chrome in CI.

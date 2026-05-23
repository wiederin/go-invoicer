# Changelog

## [Unreleased]

### Added

- Localization (`render/i18n`) — English, German, French
- Multilingual invoice template
- Stripe API client (`FetchInvoice`, `SyncInvoice`)
- Stripe webhook HTTP handler with signature verification
- `examples/stripe_sync`
- Platform API MVP (`POST /v1/render/html`, `/v1/render/pdf`, `/v1/webhooks/stripe`)
- Renderer HTTP mode (`renderer -listen :8081`)

## [0.1.0] - 2026-05-23

### Added

- Invoice domain model with builder, validation, and tax breakdown
- `currency` and `tax` packages
- HTML rendering with default, minimal, and Swiss QR-bill templates
- PDF export via Chromium (`pdf.ChromiumRenderer`)
- Swiss QR-bill payload generation (`qr/swiss`)
- Stripe domain mapper and webhook parsing (`integrations/stripe`)
- Examples: `basic`, `with_template`, `swiss_qr`
- MkDocs documentation scaffold

[0.1.0]: https://github.com/wiederin/go-invoicer/releases/tag/v0.1.0

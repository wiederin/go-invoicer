# Templates & branding

go-invoicer renders invoices from HTML templates through Chromium (same pipeline in OSS and hosted platform).

## Built-in templates

List available templates:

```http
GET /v1/templates
```

Common IDs include `default`, `minimal`, `modern`, and Swiss QR-bill variants. Set your org default with `PATCH /v1/org/rendering`:

```json
{
  "template": "modern",
  "locale": "de-CH",
  "iban": "CH93…"
}
```

## Logo and accent (Business+)

On eligible plans, set branding on the rendering profile:

```json
{
  "logo_url": "https://cdn.example.com/logo.png",
  "accent_color": "#0ea5e9"
}
```

- **Logo** — shown in the invoice header (HTTPS URL; PNG/SVG recommended).
- **Accent** — used for headings, rules, and emphasis in supported templates.

The dashboard **Settings → Branding** and **Studio** preview apply these values live. Without the branding entitlement, logo/accent fields are ignored at PDF render time.

## Template Studio

Open **Studio** (`/preview` in the dashboard) to:

- Browse all templates with live HTML preview
- Reorder layout blocks (header, parties, line items, totals, notes)
- Try sample presets (consulting, SaaS, Swiss B2B)
- **Save org defaults** — persists template, locale, IBAN, and layout blocks to `PATCH /v1/org/rendering`

Line columns, numbering, and default notes are edited under **Settings**; Studio shows a summary and links there.

## Portal vs invoice branding

| Profile | Purpose |
|---------|---------|
| `PATCH /v1/org/rendering` | PDF/HTML invoice documents |
| `PATCH /v1/org/portal` | Login shell & dashboard white-label (Enterprise) |

You can use different logos for the customer portal and printed invoices.

## Studio preview

Open **Studio** from the sidebar to:

1. Pick template and locale
2. Toggle line columns and layout blocks
3. Render sample HTML or PDF before saving

Saving from Studio updates the org rendering profile used by async PDF jobs.

## Self-hosted / OSS

In the open-source library, pass branding via `render.RenderOptions` when calling `render.HTML` or the `pdf` package. Hosted settings map to those options through `renderprofile.Options`.

## Related guides

- [Invoice customization](invoice-customization.md) — columns, numbering, notes
- [Self-hosted deployment](self-hosted.md) — run the full stack on your infrastructure

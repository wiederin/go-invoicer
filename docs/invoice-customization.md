# Invoice customization

Every organization on the hosted platform can shape how invoices look and how numbers are assigned — without changing the shared Go domain model in the OSS library.

## Where to configure

| Area | Dashboard | API |
|------|-----------|-----|
| Template, locale, IBAN | **Settings → Appearance** | `GET/PATCH /v1/org/rendering` |
| Logo & accent color | **Settings → Branding** (Business+) | same `rendering` profile |
| Line item columns & labels | **Settings → Line columns** | `line_columns` on rendering profile |
| Layout block order | **Studio** (`/preview`) | `layout_blocks` on rendering profile |
| Invoice numbering | **Settings → Numbering** | `number_prefix`, `number_pattern`, `GET /v1/invoices/next-number` |
| Default footer notes | **Settings → Default notes** | `default_notes` on rendering profile |
| Stripe billing mode | **Settings → Integrations → Stripe** | `stripe_billing_mode` on rendering profile |

The **Customization overview** card on Settings summarizes the current profile and links to each section.

## Line item columns

Toggle which columns appear on the PDF table:

- Quantity
- Unit price
- Tax rate / amount
- Line total

You can also override column header text (for example “Qty” vs “Quantity”). At least **description** and a total column remain visible.

## Layout blocks

Templates support reorderable sections (header, parties, line items, totals, payment details, notes). Use **Studio** to drag blocks and preview HTML/PDF with your org defaults.

## Invoice numbering

Configure a prefix and pattern with tokens:

- `{YYYY}` — four-digit year
- `{YY}` — two-digit year
- `{SEQ}` — zero-padded sequence (width from pattern)

Leave `invoice.number` empty on create or generate — the API allocates the next value from your pattern.

## Default notes

`default_notes` is appended when an invoice has no `notes` field. Useful for payment terms, VAT disclaimers, or bank instructions that repeat on every document.

## Plan gates

| Feature | Typical plan |
|---------|----------------|
| Line columns, numbering, layout, default notes | All paid tiers |
| Logo & accent on PDFs | Business and above |
| White-label portal | Enterprise |

See [Paid features guide](paid-features.md) for the full matrix.

## Automation uses the same defaults

Inbound hooks, `POST /v1/invoices/generate`, and cron schedules all call `applyRenderingDefaults` — so API automation picks up your Settings profile automatically.

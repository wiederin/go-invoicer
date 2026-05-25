# Templates

| File | Template name | API ID | Render method |
|------|---------------|--------|---------------|
| `default/invoice.html` | `invoice.html` | `default` | `RenderDefault` |
| `modern/modern.html` | `modern.html` | `modern` | `RenderModern` |
| `studio/studio.html` | `studio.html` | `studio` | `RenderStudio` |
| `minimal/minimal.html` | `minimal.html` | `minimal` | `RenderMinimal` |
| `stratosphere/stratosphere.html` | `stratosphere.html` | `stratosphere` | `RenderStratosphere` |
| `ocean/ocean.html` | `ocean.html` | `ocean` | `RenderOcean` |
| `ledger/ledger.html` | `ledger.html` | `ledger` | `RenderLedger` |
| `mist/mist.html` | `mist.html` | `mist` | `RenderMist` |
| `multilingual/multilingual.html` | `multilingual.html` | `multilingual` | `RenderMultilingual(inv, locale)` |
| `swiss-qr/swiss.html` | `swiss.html` | `swiss` | `RenderSwiss(inv, iban)` |

Embedded via `templates.FS` and loaded by `render.DefaultEngine()`.

**Modern** — card layout, indigo accent, SaaS-style line items.  
**Studio** — dark header band, serif typography, premium editorial feel.  
**Stratosphere / Ocean / Ledger / Mist** — Stratosphere branding pack (sky `#38a8ff`, ocean `#1a7dd4`, mist `#daeeff`, cloud `#f0f7ff`, void `#07111f`); Outfit + Space Mono, matching the hosted dashboard palette.

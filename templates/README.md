# Templates

| File | Template name | API ID | Render method |
|------|---------------|--------|---------------|
| `default/invoice.html` | `invoice.html` | `default` | `RenderDefault` |
| `modern/modern.html` | `modern.html` | `modern` | `RenderModern` |
| `studio/studio.html` | `studio.html` | `studio` | `RenderStudio` |
| `minimal/minimal.html` | `minimal.html` | `minimal` | `RenderMinimal` |
| `multilingual/multilingual.html` | `multilingual.html` | `multilingual` | `RenderMultilingual(inv, locale)` |
| `swiss-qr/swiss.html` | `swiss.html` | `swiss` | `RenderSwiss(inv, iban)` |

Embedded via `templates.FS` and loaded by `render.DefaultEngine()`.

**Modern** — card layout, indigo accent, SaaS-style line items.  
**Studio** — dark header band, serif typography, premium editorial feel.

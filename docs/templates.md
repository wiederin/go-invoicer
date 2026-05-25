# Templates

| Template | Constant | Use case |
|----------|----------|----------|
| Default | `invoice.html` | Standard business invoice |
| Minimal | `minimal.html` | Clean one-page layout |
| Modern | `modern.html` | SaaS card layout |
| Studio | `studio.html` | Editorial premium layout |
| Stratosphere | `stratosphere.html` | Stratosphere pack — card + sky accent bar |
| Ocean | `ocean.html` | Stratosphere pack — dark hero header |
| Ledger | `ledger.html` | Stratosphere pack — monospace finance grid |
| Mist | `mist.html` | Stratosphere pack — soft mist/cloud background |
| Swiss QR | `swiss.html` | Swiss QR-bill with payment slip |

```go
engine, _ := render.DefaultEngine()
html, _ := engine.RenderDefault(inv)
html, _ := engine.RenderStratosphere(inv)
html, _ := engine.RenderOcean(inv)
html, _ := engine.RenderLedger(inv)
html, _ := engine.RenderMist(inv)
html, _ := engine.RenderSwiss(inv, "CH93...")
```

Custom templates: pass your own `embed.FS` to `render.NewEngine` — see `examples/with_template`.

Branding-pack templates use the same palette as the go-invoicer dashboard (Outfit, Space Mono, sky/ocean/mist/cloud/void tokens).

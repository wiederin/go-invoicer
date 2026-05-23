# Templates

| Template | Constant | Use case |
|----------|----------|----------|
| Default | `invoice.html` | Standard business invoice |
| Minimal | `minimal.html` | Clean one-page layout |
| Swiss QR | `swiss.html` | Swiss QR-bill with payment slip |

```go
engine, _ := render.DefaultEngine()
html, _ := engine.RenderDefault(inv)
html, _ := engine.RenderMinimal(inv)
html, _ := engine.RenderSwiss(inv, "CH93...")
```

Custom templates: pass your own `embed.FS` to `render.NewEngine` — see `examples/with_template`.

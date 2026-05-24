package pdf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ChromiumRenderer renders HTML to PDF using headless Chromium.
type ChromiumRenderer struct {
	opts   []chromedp.ExecAllocatorOption
	timeout time.Duration
}

// ChromiumOption configures ChromiumRenderer.
type ChromiumOption func(*ChromiumRenderer)

// WithChromiumTimeout sets the per-render timeout (default 30s).
func WithChromiumTimeout(d time.Duration) ChromiumOption {
	return func(r *ChromiumRenderer) { r.timeout = d }
}

// chromiumExecPath returns a Chrome/Chromium binary (chromedp defaults to google-chrome).
func chromiumExecPath() string {
	if p := os.Getenv("CHROMIUM_PATH"); p != "" {
		return p
	}
	if p := os.Getenv("CHROME_PATH"); p != "" {
		return p
	}
	for _, name := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

// NewChromiumRenderer creates a renderer using a local Chromium/Chrome binary.
func NewChromiumRenderer(opts ...ChromiumOption) *ChromiumRenderer {
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	if path := chromiumExecPath(); path != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(path))
	}
	r := &ChromiumRenderer{
		opts:    allocOpts,
		timeout: 30 * time.Second,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// RenderHTML implements Renderer.
func (r *ChromiumRenderer) RenderHTML(ctx context.Context, html string) ([]byte, error) {
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, r.opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var pdfBuf []byte
	err := chromedp.Run(browserCtx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			tree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(tree.Frame.ID, html).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("pdf: chromium render: %w", err)
	}
	return pdfBuf, nil
}

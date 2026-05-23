package pdf

import (
	"context"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render"
)

// RenderInvoice renders HTML via engine then converts to PDF.
func RenderInvoice(ctx context.Context, eng *render.Engine, ren Renderer, inv *invoice.Invoice, templateName string) ([]byte, error) {
	html, err := eng.RenderInvoice(inv, templateName)
	if err != nil {
		return nil, err
	}
	return ren.RenderHTML(ctx, html)
}

// RenderInvoiceDefault uses the default template.
func RenderInvoiceDefault(ctx context.Context, eng *render.Engine, ren Renderer, inv *invoice.Invoice) ([]byte, error) {
	return RenderInvoice(ctx, eng, ren, inv, render.DefaultTemplate)
}

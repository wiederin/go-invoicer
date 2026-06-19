package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render/i18n"
	"github.com/wiederin/go-invoicer/templates"
)

// Engine renders invoices using HTML templates.
type Engine struct {
	templates *template.Template
}

// NewEngine parses templates from an fs.FS (paths like "default/invoice.html").
func NewEngine(fsys fs.FS, patterns ...string) (*Engine, error) {
	if len(patterns) == 0 {
		patterns = []string{"default/*.html"}
	}
	funcs := template.FuncMap{
		"customFieldValue": func(fields []CustomFieldView, key string) string {
			for _, cf := range fields {
				if cf.Key == key {
					return cf.Value
				}
			}
			return ""
		},
	}
	tmpl, err := template.New("invoice").Funcs(funcs).ParseFS(fsys, patterns...)
	if err != nil {
		return nil, fmt.Errorf("render: parse templates: %w", err)
	}
	return &Engine{templates: tmpl}, nil
}

// DefaultEngine uses all embedded templates.
func DefaultEngine() (*Engine, error) {
	return NewEngine(templates.FS,
		"partials/*.html",
		"default/*.html", "minimal/*.html", "swiss-qr/*.html", "multilingual/*.html",
		"modern/*.html", "studio/*.html",
		"stratosphere/*.html", "ocean/*.html", "ledger/*.html", "mist/*.html",
		"settlement/*.html")
}

// RenderOptions configures a render call.
type RenderOptions struct {
	Locale          i18n.Locale
	Branding        *BrandingView
	LayoutBlocks    []string
	BilingualLabels bool // when true, labels show both the locale language and English (for multilingual template)
}

// RenderWithOptions renders using locale-aware view data.
func (e *Engine) RenderWithOptions(inv *invoice.Invoice, templateName string, opts RenderOptions) (string, error) {
	if err := inv.Validate(); err != nil {
		return "", err
	}
	locale := opts.Locale
	if locale == "" {
		locale = i18n.LocaleEN
	}
	var view InvoiceView
	if opts.BilingualLabels {
		view = invoiceViewFromLocaleBilingual(inv, locale)
	} else {
		view = InvoiceViewFromLocale(inv, locale)
	}
	data := WithLayout(WithBranding(view, opts.Branding), opts.LayoutBlocks)
	var buf bytes.Buffer
	if err := e.templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("render: execute %s: %w", templateName, err)
	}
	return buf.String(), nil
}

// RenderMultilingual renders the multilingual template.
func (e *Engine) RenderMultilingual(inv *invoice.Invoice, locale i18n.Locale) (string, error) {
	return e.RenderWithOptions(inv, templates.MultilingualInvoice, RenderOptions{Locale: locale})
}

// RenderMinimal renders the minimal template.
func (e *Engine) RenderMinimal(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.MinimalInvoice)
}

// RenderModern renders the modern SaaS-style template.
func (e *Engine) RenderModern(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.ModernInvoice)
}

// RenderStudio renders the editorial studio template.
func (e *Engine) RenderStudio(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.StudioInvoice)
}

// RenderStratosphere renders the Stratosphere branding-pack template.
func (e *Engine) RenderStratosphere(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.StratosphereInvoice)
}

// RenderOcean renders the Ocean branding-pack template.
func (e *Engine) RenderOcean(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.OceanInvoice)
}

// RenderLedger renders the Ledger branding-pack template.
func (e *Engine) RenderLedger(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.LedgerInvoice)
}

// RenderMist renders the Mist branding-pack template.
func (e *Engine) RenderMist(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, templates.MistInvoice)
}

// RenderInvoice renders an invoice with the named template (e.g. "default/invoice.html").
func (e *Engine) RenderInvoice(inv *invoice.Invoice, templateName string) (string, error) {
	if err := inv.Validate(); err != nil {
		return "", err
	}
	data := WithLayout(InvoiceViewFrom(inv), nil)
	var buf bytes.Buffer
	if err := e.templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("render: execute %s: %w", templateName, err)
	}
	return buf.String(), nil
}

// DefaultTemplate is the name of the embedded default invoice template.
const DefaultTemplate = "invoice.html"

// RenderDefault renders the default embedded invoice template.
func (e *Engine) RenderDefault(inv *invoice.Invoice) (string, error) {
	return e.RenderInvoice(inv, DefaultTemplate)
}

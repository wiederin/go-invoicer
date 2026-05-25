package render

// BrandingView controls optional logo, accent color, and line-item columns on PDF/HTML templates.
type BrandingView struct {
	LogoURL       string
	AccentColor   string
	ShowQuantity  bool
	ShowUnitPrice bool
	ShowTax       bool
	ShowLineTotal bool
	LabelDescription string
	LabelQuantity    string
	LabelUnitPrice   string
	LabelTax         string
	LabelLineTotal   string
}

// DefaultBrandingView shows all standard line columns and no logo override.
func DefaultBrandingView() BrandingView {
	return BrandingView{
		ShowQuantity:  true,
		ShowUnitPrice: true,
		ShowTax:       true,
		ShowLineTotal: true,
	}
}

// WithBranding attaches branding to a view; nil branding keeps default column visibility.
func WithBranding(view InvoiceView, branding *BrandingView) InvoiceView {
	if branding == nil {
		def := DefaultBrandingView()
		view.Branding = &def
		return view
	}
	b := *branding
	view.Branding = &b
	return view
}

// Package templates provides embedded invoice HTML templates.
package templates

import "embed"

// FS contains all embedded HTML templates.
//
//go:embed default/*.html minimal/*.html swiss-qr/*.html multilingual/*.html modern/*.html studio/*.html
var FS embed.FS

// Template file names (ParseFS uses the file basename as template name).
const (
	DefaultInvoice       = "invoice.html"
	MinimalInvoice       = "minimal.html"
	SwissInvoice         = "swiss.html"
	MultilingualInvoice = "multilingual.html"
	ModernInvoice       = "modern.html"
	StudioInvoice       = "studio.html"
)

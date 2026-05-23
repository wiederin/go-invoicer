package render

import (
	"fmt"
	"time"

	"github.com/wiederin/go-invoicer/currency"
	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render/i18n"
)

// PartyView is a party for templates.
type PartyView struct {
	Name    string
	Email   string
	VATID   string
	Address string
}

// LineView is a line item for templates.
type LineView struct {
	Description string
	Quantity    string
	UnitPrice   string
	TaxLabel    string
	LineTotal   string
}

// TaxLineView is a tax summary row for templates.
type TaxLineView struct {
	Label  string
	Amount string
}

// InvoiceView is template-ready invoice data.
type InvoiceView struct {
	Number   string
	IssuedAt string
	DueAt    string
	Seller   PartyView
	Buyer    PartyView
	Lines    []LineView
	Subtotal string
	TaxLines []TaxLineView
	Total    string
	Notes    string
	Labels   i18n.Labels
	Locale   string
}

// InvoiceViewFrom builds display data from a domain invoice (English labels).
func InvoiceViewFrom(inv *invoice.Invoice) InvoiceView {
	return InvoiceViewFromLocale(inv, i18n.LocaleEN)
}

// InvoiceViewFromLocale builds display data with localized labels.
func InvoiceViewFromLocale(inv *invoice.Invoice, locale i18n.Locale) InvoiceView {
	labels := i18n.LabelsFor(locale)
	cur := inv.Currency()
	view := InvoiceView{
		Number:   inv.Number,
		IssuedAt: formatDateLocale(inv.IssuedAt, locale),
		Seller:   partyView(inv.Seller),
		Buyer:    partyView(inv.Buyer),
		Subtotal: currency.FormatMinor(inv.Subtotal().Amount, cur),
		Total:    currency.FormatMinor(inv.Total().Amount, cur),
		Notes:    inv.Notes,
		Labels:   labels,
		Locale:   string(locale),
	}
	if !inv.DueAt.IsZero() {
		view.DueAt = formatDateLocale(inv.DueAt, locale)
	}
	for _, l := range inv.LineItems {
		view.Lines = append(view.Lines, LineView{
			Description: l.Description,
			Quantity:    formatQty(l.Quantity),
			UnitPrice:   currency.FormatMinor(l.UnitPrice.Amount, cur),
			TaxLabel:    formatTaxRate(l.TaxRate),
			LineTotal:   currency.FormatMinor(l.LineTotal(), cur),
		})
	}
	for _, tl := range inv.TaxBreakdown() {
		view.TaxLines = append(view.TaxLines, TaxLineView{
			Label:  fmt.Sprintf("%s %.1f%%", labels.TaxRate, tl.Rate*100),
			Amount: currency.FormatMinor(tl.TaxAmount, cur),
		})
	}
	return view
}

func partyView(p invoice.Party) PartyView {
	return PartyView{
		Name:    p.Name,
		Email:   p.Email,
		VATID:   p.VATID,
		Address: p.Address.FormatSingleLine(),
	}
}

func formatDateLocale(t time.Time, locale i18n.Locale) string {
	if t.IsZero() {
		return ""
	}
	if locale == i18n.LocaleDE {
		return t.Format("02.01.2006")
	}
	if locale == i18n.LocaleFR {
		return t.Format("02/01/2006")
	}
	return t.Format("2 Jan 2006")
}

func formatDate(t time.Time) string {
	return formatDateLocale(t, i18n.LocaleEN)
}

func formatQty(q float64) string {
	if q == float64(int64(q)) {
		return fmt.Sprintf("%d", int64(q))
	}
	return fmt.Sprintf("%.2f", q)
}

func formatTaxRate(r float64) string {
	if r <= 0 {
		return "—"
	}
	return fmt.Sprintf("%.1f%%", r*100)
}

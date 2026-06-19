package render

import (
	"fmt"
	"sort"
	"time"

	"github.com/wiederin/go-invoicer/currency"
	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/render/i18n"
)

// PartyView is a party for templates.
type PartyView struct {
	Name       string
	Email      string
	Phone      string
	Department string
	VATID      string
	Address    string
}

// CustomFieldView is a single custom metadata entry on a line item.
type CustomFieldView struct {
	Key   string
	Value string
}

// LineView is a line item for templates.
type LineView struct {
	Description       string
	Quantity          string
	UnitPrice         string
	TaxLabel          string
	LineTotal         string
	Discount          string // e.g. "10%" when DiscountPct > 0, otherwise ""
	OriginalLineTotal string // pre-discount total; non-empty only when Discount != ""
	CustomFields      []CustomFieldView
	SupplyDateStart   string // formatted date string, empty if not set
	SupplyDateEnd     string // formatted date string, empty if not set
}

// TaxLineView is a tax summary row for templates.
type TaxLineView struct {
	Label  string
	Amount string
}

// InvoiceCustomFieldView is a single invoice-level custom field for templates.
type InvoiceCustomFieldView struct {
	Key   string
	Value string
}

// InvoiceView is template-ready invoice data.
type InvoiceView struct {
	Number              string
	IssuedAt            string
	DueAt               string
	Seller              PartyView
	Buyer               PartyView
	Lines               []LineView
	Subtotal            string
	TaxLines            []TaxLineView
	Total               string
	Notes               string
	Footer              string
	FooterKeyInfo       bool
	TaxID               string
	CustomFields        []InvoiceCustomFieldView
	PaymentInstructions string
	Labels              i18n.Labels
	Locale              string
	Branding            *BrandingView
	Layout              LayoutView
}

// InvoiceViewFrom builds display data from a domain invoice (English labels).
func InvoiceViewFrom(inv *invoice.Invoice) InvoiceView {
	return InvoiceViewFromLocale(inv, i18n.LocaleEN)
}

// InvoiceViewFromLocale builds display data with monolingual labels (customer language only).
func InvoiceViewFromLocale(inv *invoice.Invoice, locale i18n.Locale) InvoiceView {
	return buildInvoiceView(inv, locale, false)
}

// invoiceViewFromLocaleBilingual builds display data with bilingual labels (for multilingual template).
func invoiceViewFromLocaleBilingual(inv *invoice.Invoice, locale i18n.Locale) InvoiceView {
	return buildInvoiceView(inv, locale, true)
}

func buildInvoiceView(inv *invoice.Invoice, locale i18n.Locale, bilingual bool) InvoiceView {
	var labels i18n.Labels
	if bilingual {
		labels = i18n.LabelsWithEnglishSuffix(locale)
	} else {
		labels = i18n.LabelsFor(locale)
	}
	labels.Invoice = documentTitle(inv, labels)
	cur := inv.Currency()
	view := InvoiceView{
		Number:   inv.Number,
		IssuedAt: formatDateLocale(inv.IssuedAt, locale),
		Seller:   partyView(inv.Seller),
		Buyer:    partyView(inv.Buyer),
		Subtotal: currency.FormatMinor(inv.Subtotal().Amount, cur),
		Total:    currency.FormatMinor(inv.Total().Amount, cur),
		Notes:               inv.Notes,
		Footer:              inv.Footer,
		FooterKeyInfo:       inv.Metadata["footer_key_info"] == "true",
		TaxID:               inv.TaxID,
		PaymentInstructions: inv.Metadata["payment_instructions"],
		Labels:   labels,
		Locale:   string(locale),
	}
	for _, cf := range inv.CustomFields {
		view.CustomFields = append(view.CustomFields, InvoiceCustomFieldView{Key: cf.Key, Value: cf.Value})
	}
	if !inv.DueAt.IsZero() {
		view.DueAt = formatDateLocale(inv.DueAt, locale)
	}
	for _, l := range inv.LineItems {
		lv := LineView{
			Description: l.Description,
			Quantity:    formatQty(l.Quantity),
			UnitPrice:   currency.FormatMinor(l.UnitPrice.Amount, cur),
			TaxLabel:    formatTaxRate(l.TaxRate),
			LineTotal:   currency.FormatMinor(l.LineTotal(), cur),
		}
		if l.DiscountPct > 0 {
			lv.Discount = fmt.Sprintf("%.4g%%", l.DiscountPct)
			gross := int64(float64(l.UnitPrice.Amount) * l.Quantity)
			lv.OriginalLineTotal = currency.FormatMinor(gross, cur)
		}
		if l.SupplyDateStart != "" {
			if t, err := time.Parse("2006-01-02", l.SupplyDateStart[:10]); err == nil {
				lv.SupplyDateStart = formatDateLocale(t, locale)
			} else {
				lv.SupplyDateStart = l.SupplyDateStart
			}
		}
		if l.SupplyDateEnd != "" {
			if t, err := time.Parse("2006-01-02", l.SupplyDateEnd[:10]); err == nil {
				lv.SupplyDateEnd = formatDateLocale(t, locale)
			} else {
				lv.SupplyDateEnd = l.SupplyDateEnd
			}
		}
		if len(l.CustomFields) > 0 {
			keys := make([]string, 0, len(l.CustomFields))
			for k := range l.CustomFields {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				lv.CustomFields = append(lv.CustomFields, CustomFieldView{Key: k, Value: l.CustomFields[k]})
			}
		}
		view.Lines = append(view.Lines, lv)
	}
	for _, tl := range inv.TaxBreakdown() {
		view.TaxLines = append(view.TaxLines, TaxLineView{
			Label:  fmt.Sprintf("%s %.1f%%", labels.TaxRate, tl.Rate*100),
			Amount: currency.FormatMinor(tl.TaxAmount, cur),
		})
	}
	return WithLayout(WithBranding(view, nil), nil)
}

func documentTitle(inv *invoice.Invoice, labels i18n.Labels) string {
	switch inv.NormalizedKind() {
	case invoice.DocumentQuote:
		return labels.Quote
	case invoice.DocumentProforma:
		return labels.Proforma
	case invoice.DocumentCreditNote:
		return labels.CreditNote
	case invoice.DocumentReceipt:
		return labels.Receipt
	default:
		return labels.Invoice
	}
}

func partyView(p invoice.Party) PartyView {
	return PartyView{
		Name:       p.Name,
		Email:      p.Email,
		Phone:      p.Phone,
		Department: p.Department,
		VATID:      p.VATID,
		Address:    p.Address.FormatSingleLine(),
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

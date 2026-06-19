// Package i18n provides localized invoice template labels.
package i18n

// Locale is a BCP 47 language tag (subset supported).
type Locale string

const (
	LocaleEN Locale = "en"
	LocaleDE Locale = "de"
	LocaleFR Locale = "fr"
)

// Labels are UI strings for invoice templates.
type Labels struct {
	Invoice     string
	Quote       string
	Proforma    string
	CreditNote  string
	Receipt     string
	Issued      string
	Due         string
	From        string
	BillTo      string
	Description string
	Quantity    string
	UnitPrice   string
	Tax         string
	Amount      string
	Subtotal    string
	TotalDue    string
	Notes               string
	PaymentInstructions string
	TaxRate             string // format with %.1f%% via helper in render
}

// LabelsFor returns labels for a locale (defaults to English).
func LabelsFor(locale Locale) Labels {
	switch locale {
	case LocaleDE:
		return Labels{
			Invoice: "Rechnung", Quote: "Angebot", Proforma: "Proforma-Rechnung", CreditNote: "Gutschrift", Receipt: "Quittung",
			Issued: "Datum", Due: "Zahlbar bis",
			From: "Von", BillTo: "An",
			Description: "Beschreibung", Quantity: "Menge", UnitPrice: "Einzelpreis",
			Tax: "MwSt.", Amount: "Betrag", Subtotal: "Zwischensumme",
			TotalDue: "Total", Notes: "Bemerkungen", PaymentInstructions: "Zahlungshinweise", TaxRate: "MwSt.",
		}
	case LocaleFR:
		return Labels{
			Invoice: "Facture", Quote: "Devis", Proforma: "Facture pro forma", CreditNote: "Avoir", Receipt: "Reçu",
			Issued: "Date", Due: "Échéance",
			From: "De", BillTo: "À",
			Description: "Description", Quantity: "Qté", UnitPrice: "Prix unitaire",
			Tax: "TVA", Amount: "Montant", Subtotal: "Sous-total",
			TotalDue: "Total dû", Notes: "Remarques", PaymentInstructions: "Instructions de paiement", TaxRate: "TVA",
		}
	default:
		return Labels{
			Invoice: "Invoice", Quote: "Quote", Proforma: "Proforma invoice", CreditNote: "Credit note", Receipt: "Receipt",
			Issued: "Invoice date", Due: "Due",
			From: "From", BillTo: "Bill to",
			Description: "Description", Quantity: "Qty", UnitPrice: "Unit price",
			Tax: "Tax", Amount: "Amount", Subtotal: "Subtotal",
			TotalDue: "Total due", Notes: "Notes", PaymentInstructions: "Payment instructions", TaxRate: "Tax",
		}
	}
}

// LabelsWithEnglishSuffix returns labels that append an English translation
// after a slash for non-English locales (e.g. "Datum / Invoice date").
// Document-type names (Invoice, Quote, etc.) and TaxRate are kept in the local
// language only because they appear in titles and calculation rows.
func LabelsWithEnglishSuffix(locale Locale) Labels {
	l := LabelsFor(locale)
	if locale == LocaleEN {
		return l
	}
	en := LabelsFor(LocaleEN)
	l.Issued = l.Issued + " / " + en.Issued
	l.Due = l.Due + " / " + en.Due
	l.From = l.From + " / " + en.From
	l.BillTo = l.BillTo + " / " + en.BillTo
	l.Description = l.Description + " / " + en.Description
	l.Quantity = l.Quantity + " / " + en.Quantity
	l.UnitPrice = l.UnitPrice + " / " + en.UnitPrice
	l.Tax = l.Tax + " / " + en.Tax
	l.Amount = l.Amount + " / " + en.Amount
	l.Subtotal = l.Subtotal + " / " + en.Subtotal
	l.TotalDue = l.TotalDue + " / " + en.TotalDue
	l.Notes = l.Notes + " / " + en.Notes
	l.PaymentInstructions = l.PaymentInstructions + " / " + en.PaymentInstructions
	return l
}

// ParseLocale normalizes a locale string.
func ParseLocale(s string) Locale {
	switch s {
	case "de", "de-CH", "de-DE":
		return LocaleDE
	case "fr", "fr-CH", "fr-FR":
		return LocaleFR
	default:
		return LocaleEN
	}
}

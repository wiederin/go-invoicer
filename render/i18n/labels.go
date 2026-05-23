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
	Notes       string
	TaxRate     string // format with %.1f%% via helper in render
}

// LabelsFor returns labels for a locale (defaults to English).
func LabelsFor(locale Locale) Labels {
	switch locale {
	case LocaleDE:
		return Labels{
			Invoice: "Rechnung", Issued: "Datum", Due: "Zahlbar bis",
			From: "Von", BillTo: "An",
			Description: "Beschreibung", Quantity: "Menge", UnitPrice: "Einzelpreis",
			Tax: "MwSt.", Amount: "Betrag", Subtotal: "Zwischensumme",
			TotalDue: "Total", Notes: "Bemerkungen", TaxRate: "MwSt.",
		}
	case LocaleFR:
		return Labels{
			Invoice: "Facture", Issued: "Date", Due: "Échéance",
			From: "De", BillTo: "À",
			Description: "Description", Quantity: "Qté", UnitPrice: "Prix unitaire",
			Tax: "TVA", Amount: "Montant", Subtotal: "Sous-total",
			TotalDue: "Total dû", Notes: "Remarques", TaxRate: "TVA",
		}
	default:
		return Labels{
			Invoice: "Invoice", Issued: "Issued", Due: "Due",
			From: "From", BillTo: "Bill to",
			Description: "Description", Quantity: "Qty", UnitPrice: "Unit price",
			Tax: "Tax", Amount: "Amount", Subtotal: "Subtotal",
			TotalDue: "Total due", Notes: "Notes", TaxRate: "Tax",
		}
	}
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

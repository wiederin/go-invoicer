package invoice

// DocumentKind classifies a commercial document for rendering and workflows.
type DocumentKind string

const (
	DocumentInvoice    DocumentKind = "invoice"
	DocumentQuote      DocumentKind = "quote"
	DocumentProforma   DocumentKind = "proforma"
	DocumentCreditNote DocumentKind = "credit_note"
	DocumentReceipt    DocumentKind = "receipt"
)

// NormalizedKind returns the document kind, defaulting to invoice when unset or unknown.
func (i *Invoice) NormalizedKind() DocumentKind {
	if i == nil {
		return DocumentInvoice
	}
	switch DocumentKind(i.Kind) {
	case DocumentQuote, DocumentProforma, DocumentCreditNote, DocumentReceipt:
		return DocumentKind(i.Kind)
	default:
		return DocumentInvoice
	}
}

// IsQuote reports whether the document is a quote or proforma (non-binding offer).
func (i *Invoice) IsQuote() bool {
	k := i.NormalizedKind()
	return k == DocumentQuote || k == DocumentProforma
}

package render

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/qr/swiss"
	"github.com/wiederin/go-invoicer/templates"
)

// SwissInvoiceView extends InvoiceView with Swiss QR-bill data.
type SwissInvoiceView struct {
	InvoiceView
	IBAN    string
	QRImage template.URL // data:image/png;base64,...
}

// SwissInvoiceViewFrom builds a Swiss template view from invoice + QR bill.
func SwissInvoiceViewFrom(inv *invoice.Invoice, bill swiss.Bill) (SwissInvoiceView, error) {
	payload, err := bill.Payload()
	if err != nil {
		return SwissInvoiceView{}, err
	}
	qr, err := swiss.QRDataURL(payload, 200)
	if err != nil {
		return SwissInvoiceView{}, err
	}
	return SwissInvoiceView{
		InvoiceView: InvoiceViewFrom(inv),
		IBAN:        swiss.CompactIBAN(bill.IBAN),
		QRImage:     template.URL(qr),
	}, nil
}

// BillFromInvoice builds a Swiss QR bill from an invoice and creditor IBAN.
func BillFromInvoice(inv *invoice.Invoice, iban string) swiss.Bill {
	return swiss.Bill{
		IBAN:                iban,
		Creditor:            inv.Seller,
		Debtor:              inv.Buyer,
		Amount:              inv.Total().Amount,
		Currency:            inv.Currency(),
		UnstructuredMessage: fmt.Sprintf("Invoice %s", inv.Number),
		BillingNumber:       inv.Number,
	}
}

// RenderSwiss renders the Swiss QR-bill template.
func (e *Engine) RenderSwiss(inv *invoice.Invoice, iban string) (string, error) {
	if err := inv.Validate(); err != nil {
		return "", err
	}
	bill := BillFromInvoice(inv, iban)
	view, err := SwissInvoiceViewFrom(inv, bill)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := e.templates.ExecuteTemplate(&buf, templates.SwissInvoice, view); err != nil {
		return "", fmt.Errorf("render: swiss template: %w", err)
	}
	return buf.String(), nil
}

package components

import (
  "fmt"
)

func (doc *Document) AddPaymentTerm() {
	if len(doc.PaymentTerm) > 0 {
		paymentTermString := fmt.Sprintf(
			"%s: %s",
			doc.encodeString(doc.Config.TextPaymentTermTitle),
			doc.encodeString(doc.PaymentTerm),
		)
		doc.Pdf.SetY(doc.Pdf.GetY() + 15)
		doc.Pdf.SetX(120)
		doc.Pdf.SetFont(doc.Config.BoldFont, "B", 10)
		doc.Pdf.CellFormat(80, 4, doc.encodeString(paymentTermString), "0", 0, "R", false, 0, "")
	}
}

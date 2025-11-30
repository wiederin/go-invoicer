package components

import (
    "github.com/wiederin/go-invoicer/constants"
)

func (doc *Document) AddNotes() {
	if len(doc.Notes) == 0 {
		return
	}

	currentY := doc.Pdf.GetY()

	doc.Pdf.SetFont(doc.Config.Font, "", 9)
	doc.Pdf.SetX(constants.BaseMargin)
	doc.Pdf.SetRightMargin(100)
	doc.Pdf.SetY(currentY + 10)

	_, lineHt := doc.Pdf.GetFontSize()
	html := doc.Pdf.HTMLBasicNew()
	html.Write(lineHt, doc.encodeString(doc.Notes))

	doc.Pdf.SetRightMargin(constants.BaseMargin)
	doc.Pdf.SetY(currentY)
}

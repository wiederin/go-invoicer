package invoicer

import (
    "github.com/go-pdf/fpdf"
		"github.com/wiederin/go-invoicer/constants"
		"github.com/wiederin/go-invoicer/components"
)

// Build pdf document from config
func Build(doc *components.Document) (*fpdf.Fpdf, error) {
    // todo: validate document data

    // Build base doc
    doc.Pdf.SetMargins(constants.BaseMargin, constants.BaseMarginTop, constants.BaseMargin)
    doc.Pdf.SetXY(10, 10)

    // todo: add invoice data

    if doc.Header != nil {
        if err := doc.Header.ApplyHeader(doc); err != nil {
            return nil, err
        }
    }

    if doc.Footer != nil {
        if err := doc.Footer.ApplyFooter(doc); err != nil {
            return nil, err
        }
    }

    // Add first page - required for non header/footer components
    doc.Pdf.AddPage()

    doc.Pdf.SetFont(doc.Config.Font, "", 12)
		doc.Title()
    doc.Meta()
  	companyContactBottom := doc.Company.AppendCompanyContactToDoc(doc)
  	customerContactBottom := doc.Customer.AppendCustomerContactToDoc(doc)
    doc.AddDescription()

    if customerContactBottom > companyContactBottom {
  		doc.Pdf.SetXY(10, customerContactBottom)
  	} else {
  		doc.Pdf.SetXY(10, companyContactBottom)
  	}

    doc.AddItems()

    // Check page height (total bloc height = 30, 45 when doc discount)
    offset := doc.Pdf.GetY() + 30
    if doc.Discount != nil {
      offset += 15
    }
    if offset > constants.MaxPageHeight {
      doc.Pdf.AddPage()
    }

    doc.AddNotes()
    doc.AddTotal()
    doc.AddPaymentTerm()

    return doc.Pdf, nil
}

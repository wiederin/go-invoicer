package invoicer

import (
    "github.com/wiederin/go-invoicer/components"
    "github.com/wiederin/go-invoicer/constants"

    "errors"
    "testing"
    "os"
)

func TestInitWithInvalidType(t *testing.T) {
    _, err := Init(&components.Config{
        TextInvoiceType: "INVALID",
    })

    if errors.Is(err, ErrInvalidDocumentType) {
        return
    }

    t.Fatalf("expected ErrInvalidDocumentType, got %v", err)
}

func TestInit(t *testing.T) {
    doc, err := Init(&components.Config{
        TextInvoiceType:   constants.Invoice,
        TextRefTitle:      "Invoice No.",
        CurrencyPrecision: 2,
    })
    logoBytes, _ := os.ReadFile("./malu_logo.jpeg")

    if err != nil {
        t.Fatalf("error generating pdf. got error %v", err)
    }

    doc.SetVersion("1")

    doc.SetHeader(&components.HeaderFooter{
        Text:       "<center>Test header content.</center>",
        Pagination: true,
    })


    doc.SetFooter(&components.HeaderFooter{
        Text:       "<center>Test footer content</center>",
        Pagination: true,
    })

    doc.SetCompany(&components.Contact{
        Name: "Test Company",
        Logo: logoBytes,
        Address: &components.Address{
            Address:    "Test Company Address, 12",
            PostalCode: "90000",
            City:       "Test City",
        },
    })


    doc.SetCustomer(&components.Contact{
        Name: "Test Customer",
        Address: &components.Address{
            Address:    "Test Customer Address, 12",
            PostalCode: "8003",
            City:       "Test City",
            Country:    "Test country",
        },
    })

    doc.SetDescription("Dear Sir or Madam, We hereby invoice you for the following services:")
    doc.SetNotes("Terms of payment: Full payment is due upon receipt of this invoice. Late payments may incur additional charges or interest as per applicable laws. Please transfer the invoice amount to the account given below, stating the invoice number. With kind regards, MALU Agency")

    doc.AppendItem(&components.Item{
      Name: "Item 1 - ",
      Description: "Test Item",
      UnitCost: "10.0",
      Commission: "3.00",
    })

    pdf, err := Build(doc)
    if err != nil {
        t.Errorf(err.Error())
    }

    err = pdf.OutputFileAndClose("test_generator_init.pdf")

    if err != nil {
        t.Errorf(err.Error())
    }
}

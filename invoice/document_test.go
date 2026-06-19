package invoice_test

import (
	"testing"

	"github.com/wiederin/go-invoicer/invoice"
)

func TestCreditNoteRequiresRelated(t *testing.T) {
	inv := invoice.New("CN-001")
	inv.Kind = invoice.DocumentCreditNote
	inv.SetSeller(invoice.Party{Name: "Seller"})
	inv.SetBuyer(invoice.Party{Name: "Buyer"})
	inv.AddLine(invoice.NewLineItem("Refund", 1, invoice.NewMoney(1000, "CHF"), 0))

	if err := inv.Validate(); err != invoice.ErrMissingRelated {
		t.Fatalf("Validate() = %v, want ErrMissingRelated", err)
	}

	inv.RelatedNumber = "INV-100"
	if err := inv.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestNormalizedKindDefault(t *testing.T) {
	inv := invoice.New("INV-1")
	if inv.NormalizedKind() != invoice.DocumentInvoice {
		t.Fatalf("got %q", inv.NormalizedKind())
	}
	inv.Kind = "quote"
	if !inv.IsQuote() {
		t.Fatal("expected quote")
	}
}

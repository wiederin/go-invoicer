package swiss_test

import (
	"strings"
	"testing"

	"github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/qr/swiss"
)

func TestPayload(t *testing.T) {
	b := swiss.Bill{
		IBAN:     "CH93 0076 2011 6238 5295 7",
		Currency: "CHF",
		Amount:   19900,
		Creditor: invoice.Party{
			Name: "Acme GmbH",
			Address: invoice.Address{
				Line1: "Bahnhofstrasse 1", PostalCode: "8001", City: "Zürich", Country: "CH",
			},
		},
		Debtor: invoice.Party{
			Name: "Client AG",
			Address: invoice.Address{
				Line1: "Dufourstrasse 2", PostalCode: "9000", City: "St. Gallen", Country: "CH",
			},
		},
		UnstructuredMessage: "Invoice INV-001",
	}
	payload, err := b.Payload()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(payload, "SPC\n0200\n") {
		t.Fatalf("unexpected payload start: %q", payload[:20])
	}
}

func TestValidate_withQRR(t *testing.T) {
	p26 := strings.Repeat("0", 26)
	cd, err := swiss.Mod10RecursiveCheckDigit(p26)
	if err != nil {
		t.Fatal(err)
	}
	b := swiss.Bill{
		IBAN: "CH9300762011623852957", Currency: "CHF", Amount: 100,
		Reference: p26 + string(cd),
		Creditor: invoice.Party{
			Name: "A",
			Address: invoice.Address{Line1: "x", PostalCode: "1", City: "z", Country: "CH"},
		},
	}
	if err := b.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestQRDataURL(t *testing.T) {
	b := swiss.Bill{
		IBAN: "CH9300762011623852957", Currency: "CHF", Amount: 1000,
		Creditor: invoice.Party{Name: "A", Address: invoice.Address{Line1: "x", PostalCode: "1", City: "z", Country: "CH"}},
		Debtor:   invoice.Party{Name: "B", Address: invoice.Address{Line1: "y", PostalCode: "2", City: "c", Country: "CH"}},
	}
	payload, err := b.Payload()
	if err != nil {
		t.Fatal(err)
	}
	url, err := swiss.QRDataURL(payload, 128)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(url, "data:image/png;base64,") {
		t.Fatal("expected data url")
	}
}

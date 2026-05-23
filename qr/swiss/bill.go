// Package swiss builds Swiss QR-bill payloads (SPC QR code, version 2.0 subset).
//
// See SIX Swiss Payment Standards for full compliance requirements.
package swiss

import (
	"fmt"
	"strings"

	"github.com/wiederin/go-invoicer/invoice"
)

// Bill holds creditor/debtor and payment data for a Swiss QR-bill.
type Bill struct {
	IBAN                  string
	Creditor              invoice.Party
	Amount                int64  // minor units
	Currency              string // CHF or EUR
	Reference             string // optional QR reference
	UnstructuredMessage   string
	Debtor                invoice.Party
	BillingNumber         string // shown on invoice, not in QR payload
}

// Validate checks required QR-bill fields.
func (b Bill) Validate() error {
	if b.IBAN == "" {
		return fmt.Errorf("swiss qr: IBAN is required")
	}
	if b.Creditor.Name == "" {
		return fmt.Errorf("swiss qr: creditor name is required")
	}
	if b.Currency == "" {
		return fmt.Errorf("swiss qr: currency is required")
	}
	if b.Amount < 0 {
		return fmt.Errorf("swiss qr: amount cannot be negative")
	}
	return nil
}

// Payload returns the SPC QR code payload string (non-reference, type 1).
func (b Bill) Payload() (string, error) {
	if err := b.Validate(); err != nil {
		return "", err
	}

	credAddr := b.Creditor.Address
	debtorAddr := b.Debtor.Address

	amountStr := ""
	if b.Amount > 0 {
		amountStr = formatAmount(b.Amount)
	}

	lines := []string{
		"SPC",                     // QR Type
		"0200",                    // Version
		"1",                       // Coding: UTF-8
		compactIBAN(b.IBAN),       // IBAN
		"S",                       // Creditor address type (structured)
		b.Creditor.Name,
		credAddr.Line1,
		joinZipCity(credAddr.PostalCode, credAddr.City),
		credAddr.Country,
		"",                        // Ultimate creditor (empty)
		"", "", "", "", "",
		amountStr,
		strings.ToUpper(b.Currency),
		"S", // Debtor address type
		b.Debtor.Name,
		debtorAddr.Line1,
		joinZipCity(debtorAddr.PostalCode, debtorAddr.City),
		debtorAddr.Country,
		refType(b.Reference),
		b.Reference,
		b.UnstructuredMessage,
		"EPD",
	}
	return strings.Join(lines, "\n"), nil
}

func refType(reference string) string {
	if reference == "" {
		return "NON"
	}
	return "QRR"
}

// CompactIBAN normalizes an IBAN for QR payloads.
func CompactIBAN(iban string) string {
	return strings.ReplaceAll(strings.ToUpper(iban), " ", "")
}

func compactIBAN(iban string) string {
	return CompactIBAN(iban)
}

func joinZipCity(zip, city string) string {
	if zip == "" {
		return city
	}
	if city == "" {
		return zip
	}
	return zip + " " + city
}

func formatAmount(minor int64) string {
	return fmt.Sprintf("%d.%02d", minor/100, minor%100)
}

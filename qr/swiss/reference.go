package swiss

import (
	"fmt"
	"strings"
	"unicode"
)

// mod10Table is the Swiss QR-bill / ISR recursive mod10 table.
var mod10Table = [10]int{0, 9, 4, 6, 8, 2, 7, 1, 3, 5}

// NormalizeReference strips spaces and returns digits-only reference text.
func NormalizeReference(ref string) string {
	var b strings.Builder
	for _, r := range ref {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Mod10RecursiveCheckDigit returns the check digit for the given digit string (without check digit).
func Mod10RecursiveCheckDigit(digits string) (byte, error) {
	if digits == "" {
		return 0, fmt.Errorf("swiss qr: empty reference")
	}
	carry := 0
	for _, r := range digits {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("swiss qr: reference must be numeric")
		}
		carry = mod10Table[(carry+int(r-'0'))%10]
	}
	return byte('0' + (10-carry)%10), nil
}

// ValidateQRR checks a 27-digit QR reference (26 payload digits + mod10 check digit).
func ValidateQRR(reference string) error {
	ref := NormalizeReference(reference)
	if len(ref) != 27 {
		return fmt.Errorf("swiss qr: QRR reference must be 27 digits, got %d", len(ref))
	}
	payload := ref[:26]
	check, err := Mod10RecursiveCheckDigit(payload)
	if err != nil {
		return err
	}
	if ref[26] != check {
		return fmt.Errorf("swiss qr: invalid QRR check digit")
	}
	return nil
}

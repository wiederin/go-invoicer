package swiss_test

import (
	"strings"
	"testing"

	"github.com/wiederin/go-invoicer/qr/swiss"
)

func TestValidateQRR_roundtrip(t *testing.T) {
	payload := strings.Repeat("0", 26)
	check, err := swiss.Mod10RecursiveCheckDigit(payload)
	if err != nil {
		t.Fatal(err)
	}
	ref := payload + string(check)
	if len(ref) != 27 {
		t.Fatalf("ref len %d", len(ref))
	}
	if err := swiss.ValidateQRR(ref); err != nil {
		t.Fatalf("expected valid QRR: %v", err)
	}
}

func TestValidateQRR_invalidCheckDigit(t *testing.T) {
	payload := strings.Repeat("1", 26)
	check, err := swiss.Mod10RecursiveCheckDigit(payload)
	if err != nil {
		t.Fatal(err)
	}
	bad := byte('0')
	if check != '0' {
		bad = check - 1
	} else {
		bad = '1'
	}
	ref := payload + string(bad)
	if err := swiss.ValidateQRR(ref); err == nil {
		t.Fatal("expected error for bad check digit")
	}
}

func TestValidateQRR_wrongLength(t *testing.T) {
	if err := swiss.ValidateQRR("123"); err == nil {
		t.Fatal("expected length error")
	}
}

func TestNormalizeReference(t *testing.T) {
	got := swiss.NormalizeReference("12 34 56 78 90")
	if got != "1234567890" {
		t.Fatalf("got %q", got)
	}
}

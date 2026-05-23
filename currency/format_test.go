package currency_test

import (
	"testing"

	"github.com/wiederin/go-invoicer/currency"
)

func TestFormatMinorCHF(t *testing.T) {
	got := currency.FormatMinor(123456, "CHF")
	want := "CHF 1'234.56"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

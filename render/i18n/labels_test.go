package i18n_test

import (
	"testing"

	"github.com/wiederin/go-invoicer/render/i18n"
)

func TestLabelsDE(t *testing.T) {
	l := i18n.LabelsFor(i18n.LocaleDE)
	if l.Invoice != "Rechnung" {
		t.Fatalf("got %q", l.Invoice)
	}
}

func TestParseLocale(t *testing.T) {
	if i18n.ParseLocale("de-CH") != i18n.LocaleDE {
		t.Fatal("expected de")
	}
}

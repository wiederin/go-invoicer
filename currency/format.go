// Package currency formats monetary amounts for display.
package currency

import (
	"fmt"
	"strings"
)

// FormatMinor formats an amount in minor units (cents) for the given ISO 4217 code.
func FormatMinor(amount int64, code string) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
	major := amount / 100
	minor := amount % 100

	sign := ""
	if negative {
		sign = "-"
	}

	switch strings.ToUpper(code) {
	case "CHF":
		return fmt.Sprintf("%sCHF %s.%02d", sign, groupThousands(major, "'"), minor)
	case "EUR":
		return fmt.Sprintf("%s%s,%02d €", sign, groupThousands(major, "."), minor)
	case "USD", "GBP":
		sym := "$"
		if code == "GBP" {
			sym = "£"
		}
		return fmt.Sprintf("%s%s%s.%02d", sign, sym, groupThousands(major, ","), minor)
	default:
		if code == "" {
			return fmt.Sprintf("%s%d.%02d", sign, major, minor)
		}
		return fmt.Sprintf("%s%s %s.%02d", sign, code, groupThousands(major, ","), minor)
	}
}

func groupThousands(n int64, sep string) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	return strings.Join(parts, sep)
}

package currency

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Currency struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	SymbolPosition  string `json:"symbol_position"`
	DecimalPlaces   int32  `json:"decimal_places"`
	DecimalSep      string `json:"decimal_sep"`
	ThousandsSep    string `json:"thousands_sep"`
}

var Currencies = map[string]Currency{
	"USD": {Code: "USD", Name: "US Dollar", Symbol: "$", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: ","},
	"EUR": {Code: "EUR", Name: "Euro", Symbol: "€", SymbolPosition: "after", DecimalPlaces: 2, DecimalSep: ",", ThousandsSep: "."},
	"CHF": {Code: "CHF", Name: "Swiss Franc", Symbol: "CHF", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: "'"},
	"GBP": {Code: "GBP", Name: "British Pound", Symbol: "£", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: ","},
	"JPY": {Code: "JPY", Name: "Japanese Yen", Symbol: "¥", SymbolPosition: "before", DecimalPlaces: 0, DecimalSep: "", ThousandsSep: ","},
	"CAD": {Code: "CAD", Name: "Canadian Dollar", Symbol: "CA$", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: ","},
	"AUD": {Code: "AUD", Name: "Australian Dollar", Symbol: "A$", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: ","},
	"CNY": {Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: ","},
	"INR": {Code: "INR", Name: "Indian Rupee", Symbol: "₹", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ".", ThousandsSep: ","},
	"BRL": {Code: "BRL", Name: "Brazilian Real", Symbol: "R$", SymbolPosition: "before", DecimalPlaces: 2, DecimalSep: ",", ThousandsSep: "."},
}

func Get(code string) (Currency, bool) {
	c, ok := Currencies[strings.ToUpper(code)]
	return c, ok
}

type Formatter struct {
	currency Currency
}

func NewFormatter(currencyCode string) (*Formatter, error) {
	c, ok := Get(currencyCode)
	if !ok {
		return nil, fmt.Errorf("unknown currency: %s", currencyCode)
	}
	return &Formatter{currency: c}, nil
}

func (f *Formatter) Format(amount decimal.Decimal) string {
	rounded := amount.Round(f.currency.DecimalPlaces)
	
	isNegative := rounded.IsNegative()
	if isNegative {
		rounded = rounded.Abs()
	}

	str := rounded.StringFixed(f.currency.DecimalPlaces)
	
	parts := strings.Split(str, ".")
	intPart := parts[0]
	
	var formattedInt strings.Builder
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			formattedInt.WriteString(f.currency.ThousandsSep)
		}
		formattedInt.WriteRune(digit)
	}
	
	var result string
	if f.currency.DecimalPlaces > 0 && len(parts) > 1 {
		result = formattedInt.String() + f.currency.DecimalSep + parts[1]
	} else {
		result = formattedInt.String()
	}
	
	if f.currency.SymbolPosition == "before" {
		result = f.currency.Symbol + " " + result
	} else {
		result = result + " " + f.currency.Symbol
	}
	
	if isNegative {
		result = "-" + result
	}
	
	return result
}

func (f *Formatter) Parse(s string) (decimal.Decimal, error) {
	cleaned := strings.ReplaceAll(s, f.currency.Symbol, "")
	cleaned = strings.ReplaceAll(cleaned, f.currency.ThousandsSep, "")
	cleaned = strings.ReplaceAll(cleaned, f.currency.DecimalSep, ".")
	cleaned = strings.TrimSpace(cleaned)
	
	return decimal.NewFromString(cleaned)
}

func FormatSimple(amount decimal.Decimal, currencyCode string) string {
	f, err := NewFormatter(currencyCode)
	if err != nil {
		return fmt.Sprintf("%s %s", amount.StringFixed(2), currencyCode)
	}
	return f.Format(amount)
}

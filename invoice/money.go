// Package invoice provides domain models for creating and managing invoices.
package invoice

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money represents a monetary value with currency.
type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

// NewMoney creates a new Money value from a float64 amount and currency code.
func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   decimal.NewFromFloat(amount),
		Currency: currency,
	}
}

// NewMoneyFromString creates a new Money value from a string amount and currency code.
func NewMoneyFromString(amount string, currency string) (Money, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount: %w", err)
	}
	return Money{Amount: d, Currency: currency}, nil
}

// Add adds two Money values with the same currency.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount.Add(other.Amount), Currency: m.Currency}, nil
}

// Sub subtracts another Money value from this one.
func (m Money) Sub(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount.Sub(other.Amount), Currency: m.Currency}, nil
}

// Mul multiplies the Money value by a decimal factor.
func (m Money) Mul(factor decimal.Decimal) Money {
	return Money{Amount: m.Amount.Mul(factor), Currency: m.Currency}
}

// MulFloat multiplies the Money value by a float64 factor.
func (m Money) MulFloat(factor float64) Money {
	return m.Mul(decimal.NewFromFloat(factor))
}

// Round rounds the Money amount to the given number of decimal places.
func (m Money) Round(places int32) Money {
	return Money{Amount: m.Amount.Round(places), Currency: m.Currency}
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool {
	return m.Amount.IsZero()
}

// IsNegative returns true if the amount is negative.
func (m Money) IsNegative() bool {
	return m.Amount.IsNegative()
}

// Float64 returns the amount as a float64.
func (m Money) Float64() float64 {
	f, _ := m.Amount.Float64()
	return f
}

// String returns a string representation of the Money value.
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount.StringFixed(2), m.Currency)
}

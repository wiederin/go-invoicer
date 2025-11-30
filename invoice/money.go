package invoice

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Money struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   decimal.NewFromFloat(amount),
		Currency: currency,
	}
}

func NewMoneyFromString(amount string, currency string) (Money, error) {
	d, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount: %w", err)
	}
	return Money{Amount: d, Currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount.Add(other.Amount), Currency: m.Currency}, nil
}

func (m Money) Sub(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.Currency, other.Currency)
	}
	return Money{Amount: m.Amount.Sub(other.Amount), Currency: m.Currency}, nil
}

func (m Money) Mul(factor decimal.Decimal) Money {
	return Money{Amount: m.Amount.Mul(factor), Currency: m.Currency}
}

func (m Money) MulFloat(factor float64) Money {
	return m.Mul(decimal.NewFromFloat(factor))
}

func (m Money) Round(places int32) Money {
	return Money{Amount: m.Amount.Round(places), Currency: m.Currency}
}

func (m Money) IsZero() bool {
	return m.Amount.IsZero()
}

func (m Money) IsNegative() bool {
	return m.Amount.IsNegative()
}

func (m Money) Float64() float64 {
	f, _ := m.Amount.Float64()
	return f
}

func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount.StringFixed(2), m.Currency)
}

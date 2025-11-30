package invoice

import (
	"testing"
	"time"
)

func TestInvoiceBuilder(t *testing.T) {
	inv, err := New().
		Number("INV-001").
		IssueDate(time.Now()).
		DueDate(time.Now().AddDate(0, 0, 30)).
		Currency("USD").
		Supplier(Party{Name: "Supplier Co", Address: Address{Street: "123 St", City: "NYC", PostalCode: "10001", Country: "US"}}).
		Customer(Party{Name: "Customer Inc", Address: Address{Street: "456 Ave", City: "LA", PostalCode: "90001", Country: "US"}}).
		AddItem(NewLineItem("Service", 1, NewMoney(100, "USD"), 10)).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inv.Number != "INV-001" {
		t.Errorf("expected number INV-001, got %s", inv.Number)
	}

	if len(inv.LineItems) != 1 {
		t.Errorf("expected 1 line item, got %d", len(inv.LineItems))
	}
}

func TestInvoiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		builder func() *Builder
		wantErr error
	}{
		{
			name: "missing number",
			builder: func() *Builder {
				return New().
					IssueDate(time.Now()).
					DueDate(time.Now().AddDate(0, 0, 30)).
					Currency("USD").
					Supplier(Party{Name: "S"}).
					Customer(Party{Name: "C"}).
					AddItem(NewLineItem("X", 1, NewMoney(100, "USD"), 0))
			},
			wantErr: ErrMissingInvoiceNumber,
		},
		{
			name: "missing issue date",
			builder: func() *Builder {
				return New().
					Number("INV-001").
					DueDate(time.Now().AddDate(0, 0, 30)).
					Currency("USD").
					Supplier(Party{Name: "S"}).
					Customer(Party{Name: "C"}).
					AddItem(NewLineItem("X", 1, NewMoney(100, "USD"), 0))
			},
			wantErr: ErrMissingIssueDate,
		},
		{
			name: "due date before issue",
			builder: func() *Builder {
				return New().
					Number("INV-001").
					IssueDate(time.Now()).
					DueDate(time.Now().AddDate(0, 0, -1)).
					Currency("USD").
					Supplier(Party{Name: "S"}).
					Customer(Party{Name: "C"}).
					AddItem(NewLineItem("X", 1, NewMoney(100, "USD"), 0))
			},
			wantErr: ErrDueDateBeforeIssue,
		},
		{
			name: "no line items",
			builder: func() *Builder {
				return New().
					Number("INV-001").
					IssueDate(time.Now()).
					DueDate(time.Now().AddDate(0, 0, 30)).
					Currency("USD").
					Supplier(Party{Name: "S"}).
					Customer(Party{Name: "C"})
			},
			wantErr: ErrNoLineItems,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.builder().Build()
			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLineItemCalculations(t *testing.T) {
	item := NewLineItem("Product", 2, NewMoney(100, "USD"), 10)

	subtotal := item.SubTotal()
	if subtotal.Float64() != 200 {
		t.Errorf("expected subtotal 200, got %v", subtotal.Float64())
	}

	tax := item.TaxAmount()
	if tax.Float64() != 20 {
		t.Errorf("expected tax 20, got %v", tax.Float64())
	}

	gross := item.GrossAmount()
	if gross.Float64() != 220 {
		t.Errorf("expected gross 220, got %v", gross.Float64())
	}
}

func TestLineItemWithDiscount(t *testing.T) {
	item := NewLineItem("Product", 1, NewMoney(100, "USD"), 10).WithDiscount(10)

	discount := item.DiscountAmount()
	if discount.Float64() != 10 {
		t.Errorf("expected discount 10, got %v", discount.Float64())
	}

	net := item.NetAmount()
	if net.Float64() != 90 {
		t.Errorf("expected net 90, got %v", net.Float64())
	}

	tax := item.TaxAmount()
	if tax.Float64() != 9 {
		t.Errorf("expected tax 9, got %v", tax.Float64())
	}

	gross := item.GrossAmount()
	if gross.Float64() != 99 {
		t.Errorf("expected gross 99, got %v", gross.Float64())
	}
}

func TestInvoiceTotals(t *testing.T) {
	inv, _ := New().
		Number("INV-001").
		IssueDate(time.Now()).
		DueDate(time.Now().AddDate(0, 0, 30)).
		Currency("USD").
		Supplier(Party{Name: "S"}).
		Customer(Party{Name: "C"}).
		AddItem(NewLineItem("A", 1, NewMoney(100, "USD"), 10)).
		AddItem(NewLineItem("B", 2, NewMoney(50, "USD"), 10)).
		Build()

	subtotal := inv.SubTotal()
	if subtotal.Float64() != 200 {
		t.Errorf("expected subtotal 200, got %v", subtotal.Float64())
	}

	tax := inv.TotalTax()
	if tax.Float64() != 20 {
		t.Errorf("expected tax 20, got %v", tax.Float64())
	}

	gross := inv.TotalGross()
	if gross.Float64() != 220 {
		t.Errorf("expected gross 220, got %v", gross.Float64())
	}
}

func TestMoneyOperations(t *testing.T) {
	m1 := NewMoney(100, "USD")
	m2 := NewMoney(50, "USD")

	sum, err := m1.Add(m2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.Float64() != 150 {
		t.Errorf("expected 150, got %v", sum.Float64())
	}

	diff, err := m1.Sub(m2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff.Float64() != 50 {
		t.Errorf("expected 50, got %v", diff.Float64())
	}

	m3 := NewMoney(100, "EUR")
	_, err = m1.Add(m3)
	if err == nil {
		t.Error("expected currency mismatch error")
	}
}

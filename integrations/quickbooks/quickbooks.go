package quickbooks

import (
	"context"
	"fmt"

	domain "github.com/wiederin/go-invoicer/invoice"
)

// Service implements sync for QuickBooks via Client.
type Service struct {
	Client   *Client
	Supplier Supplier
}

// ImportExternal fetches an invoice by QuickBooks Id.
func (s Service) ImportExternal(ctx context.Context, _, externalID string) (*domain.Invoice, error) {
	if s.Client == nil {
		return nil, fmt.Errorf("quickbooks: client not configured")
	}
	inv, err := s.Client.FetchInvoice(ctx, externalID)
	if err != nil {
		return nil, err
	}
	return inv.ToDomain(ImportOptions{Supplier: s.Supplier}), nil
}

// PushInvoice is not implemented yet.
func (Service) PushInvoice(context.Context, string, *domain.Invoice) (string, error) {
	return "", fmt.Errorf("quickbooks: export not implemented yet")
}

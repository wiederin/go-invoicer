package xero

import (
	"context"
	"fmt"

	domain "github.com/wiederin/go-invoicer/invoice"
)

// Service implements sync for Xero via Client.
type Service struct {
	Client *Client
	Supplier Supplier
}

// ImportExternal fetches an invoice by Xero InvoiceID.
func (s Service) ImportExternal(ctx context.Context, _, externalID string) (*domain.Invoice, error) {
	if s.Client == nil {
		return nil, fmt.Errorf("xero: client not configured")
	}
	inv, err := s.Client.FetchInvoice(ctx, externalID)
	if err != nil {
		return nil, err
	}
	return inv.ToDomain(ImportOptions{Supplier: s.Supplier}), nil
}

// PushInvoice creates or updates a Xero ACCREC draft invoice.
func (s Service) PushInvoice(ctx context.Context, _ string, inv *domain.Invoice) (string, error) {
	if s.Client == nil {
		return "", fmt.Errorf("xero: client not configured")
	}
	return s.Client.PushInvoice(ctx, inv, ExportOptions{})
}

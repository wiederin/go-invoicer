// Package quickbooks provides QuickBooks invoice sync (Phase 5 — scaffold).
package quickbooks

import (
	"context"
	"fmt"

	domain "github.com/wiederin/go-invoicer/invoice"
	"github.com/wiederin/go-invoicer/integrations/sync"
)

const ProviderID = "quickbooks"

// Service will implement sync.Syncer for QuickBooks.
type Service struct{}

func (Service) Provider() string { return ProviderID }

func (Service) ImportExternal(context.Context, string, string) (*domain.Invoice, error) {
	return nil, fmt.Errorf("quickbooks: sync not implemented yet")
}

func (Service) PushInvoice(context.Context, string, *domain.Invoice) (string, error) {
	return "", fmt.Errorf("quickbooks: sync not implemented yet")
}

var _ sync.Syncer = Service{}

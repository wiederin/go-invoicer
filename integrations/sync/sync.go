// Package sync defines bidirectional invoice sync for external accounting systems.
//
// Stripe is implemented in the stripe subpackage; Xero and QuickBooks follow the same
// Importer/Pusher interfaces (Phase 5).
package sync

import (
	"context"

	domain "github.com/wiederin/go-invoicer/invoice"
)

// Importer pulls invoices from an external system into go-invoicer.
type Importer interface {
	ImportExternal(ctx context.Context, orgID, externalID string) (*domain.Invoice, error)
}

// Pusher creates or updates invoices in an external system from go-invoicer.
type Pusher interface {
	PushInvoice(ctx context.Context, orgID string, inv *domain.Invoice) (externalID string, err error)
}

// Syncer combines import and export for a provider.
type Syncer interface {
	Importer
	Pusher
	Provider() string
}

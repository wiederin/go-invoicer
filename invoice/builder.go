package invoice

import "time"

// Builder constructs a validated Invoice.
type Builder struct {
	inv *Invoice
}

// NewBuilder starts building an invoice.
func NewBuilder() *Builder {
	return &Builder{inv: &Invoice{
		IssuedAt: time.Now().UTC(),
		Metadata: make(map[string]string),
	}}
}

func (b *Builder) Number(n string) *Builder {
	b.inv.Number = n
	return b
}

func (b *Builder) IssuedAt(t time.Time) *Builder {
	b.inv.IssuedAt = t.UTC()
	return b
}

func (b *Builder) DueAt(t time.Time) *Builder {
	b.inv.DueAt = t.UTC()
	return b
}

func (b *Builder) Seller(p Party) *Builder {
	b.inv.Seller = p
	return b
}

func (b *Builder) Buyer(p Party) *Builder {
	b.inv.Buyer = p
	return b
}

func (b *Builder) Notes(n string) *Builder {
	b.inv.Notes = n
	return b
}

func (b *Builder) AddLine(item LineItem) *Builder {
	b.inv.LineItems = append(b.inv.LineItems, item)
	return b
}

func (b *Builder) Metadata(key, value string) *Builder {
	b.inv.Metadata[key] = value
	return b
}

// Build validates and returns the invoice.
func (b *Builder) Build() (*Invoice, error) {
	if err := b.inv.Validate(); err != nil {
		return nil, err
	}
	return b.inv, nil
}

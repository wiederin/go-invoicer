// Package settlement builds period summary reports from invoice rows.
package settlement

import (
	"fmt"
	"sort"
	"time"

	domain "github.com/wiederin/go-invoicer/invoice"
)

// Report summarizes invoices in a period.
type Report struct {
	OrgName      string    `json:"org_name"`
	PeriodFrom   time.Time `json:"period_from"`
	PeriodTo     time.Time `json:"period_to"`
	Currency     string    `json:"currency"`
	Lines        []Line    `json:"lines"`
	TotalMinor   int64     `json:"total_minor"`
	InvoiceCount int       `json:"invoice_count"`
}

// Line is one invoice row on a settlement report.
type Line struct {
	Number      string    `json:"number"`
	IssuedAt    time.Time `json:"issued_at"`
	BuyerName   string    `json:"buyer_name"`
	Status      string    `json:"status"`
	Kind        string    `json:"kind"`
	AmountMinor int64     `json:"amount_minor"`
	Currency    string    `json:"currency"`
}

// InvoiceInput is a minimal row for building a report.
type InvoiceInput struct {
	Number  string
	Status  string
	Invoice domain.Invoice
}

// Build aggregates invoice rows into a settlement report.
func Build(orgName string, from, to time.Time, rows []InvoiceInput) (*Report, error) {
	if to.Before(from) {
		return nil, fmt.Errorf("settlement: period end before start")
	}
	r := &Report{
		OrgName:    orgName,
		PeriodFrom: from.UTC(),
		PeriodTo:   to.UTC(),
	}
	cur := ""
	for _, row := range rows {
		inv := row.Invoice
		at := inv.IssuedAt
		if at.IsZero() {
			continue
		}
		if at.Before(from) || !at.Before(to) {
			continue
		}
		lineCur := inv.Currency()
		if lineCur == "" {
			continue
		}
		if cur == "" {
			cur = lineCur
		} else if lineCur != cur {
			return nil, fmt.Errorf("settlement: mixed currencies in period")
		}
		amt := inv.Total().Amount
		r.Lines = append(r.Lines, Line{
			Number:      row.Number,
			IssuedAt:    at.UTC(),
			BuyerName:   inv.Buyer.Name,
			Status:      row.Status,
			Kind:        string(inv.NormalizedKind()),
			AmountMinor: amt,
			Currency:    lineCur,
		})
		r.TotalMinor += amt
	}
	sort.Slice(r.Lines, func(i, j int) bool {
		return r.Lines[i].IssuedAt.Before(r.Lines[j].IssuedAt)
	})
	r.Currency = cur
	r.InvoiceCount = len(r.Lines)
	return r, nil
}

package render

import (
	"bytes"
	"fmt"
	"time"

	"github.com/wiederin/go-invoicer/currency"
	"github.com/wiederin/go-invoicer/render/i18n"
	"github.com/wiederin/go-invoicer/settlement"
	"github.com/wiederin/go-invoicer/templates"
)

// SettlementView is template data for settlement reports.
type SettlementView struct {
	OrgName      string
	PeriodFrom   string
	PeriodTo     string
	InvoiceCount int
	Total        string
	Lines        []SettlementLineView
}

// SettlementLineView is one table row.
type SettlementLineView struct {
	IssuedAt  string
	Number    string
	BuyerName string
	Kind      string
	Status    string
	Amount    string
}

// RenderSettlement renders a settlement report HTML page.
func (e *Engine) RenderSettlement(rep *settlement.Report) (string, error) {
	if rep == nil {
		return "", fmt.Errorf("render: nil settlement report")
	}
	view := settlementViewFrom(rep)
	var buf bytes.Buffer
	if err := e.templates.ExecuteTemplate(&buf, templates.SettlementReport, view); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func settlementViewFrom(rep *settlement.Report) SettlementView {
	view := SettlementView{
		OrgName:      rep.OrgName,
		PeriodFrom:   formatDateLocale(rep.PeriodFrom, i18n.LocaleEN),
		PeriodTo:     formatDateLocale(rep.PeriodTo.Add(-time.Nanosecond), i18n.LocaleEN),
		InvoiceCount: rep.InvoiceCount,
		Total:        currency.FormatMinor(rep.TotalMinor, rep.Currency),
	}
	for _, l := range rep.Lines {
		view.Lines = append(view.Lines, SettlementLineView{
			IssuedAt:  formatDateLocale(l.IssuedAt, i18n.LocaleEN),
			Number:    l.Number,
			BuyerName: l.BuyerName,
			Kind:      l.Kind,
			Status:    l.Status,
			Amount:    currency.FormatMinor(l.AmountMinor, l.Currency),
		})
	}
	return view
}

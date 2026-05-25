package quickbooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
)

// Client calls the QuickBooks Online API.
type Client struct {
	HTTPClient  *http.Client
	AccessToken string
	RealmID     string
	APIBase     string
}

// FetchInvoice retrieves an invoice by QuickBooks Id.
func (c *Client) FetchInvoice(ctx context.Context, invoiceID string) (Invoice, error) {
	if c.AccessToken == "" || c.RealmID == "" {
		return Invoice{}, fmt.Errorf("quickbooks: access token and realm required")
	}
	base := c.APIBase
	if base == "" {
		base = qboProdAPI
	}
	url := fmt.Sprintf("%s/v3/company/%s/invoice/%s?minorversion=65", base, c.RealmID, invoiceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Invoice{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Accept", "application/json")
	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	res, err := hc.Do(req)
	if err != nil {
		return Invoice{}, fmt.Errorf("quickbooks: fetch invoice: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return Invoice{}, fmt.Errorf("quickbooks: fetch invoice: %s", strings.TrimSpace(string(body)))
	}
	var payload struct {
		Invoice apiInvoice `json:"Invoice"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return Invoice{}, err
	}
	if payload.Invoice.Id == "" {
		return Invoice{}, fmt.Errorf("quickbooks: invoice not found")
	}
	return mapAPIInvoice(payload.Invoice), nil
}

type apiInvoice struct {
	Id        string `json:"Id"`
	SyncToken string `json:"SyncToken"`
	DocNumber string `json:"DocNumber"`
	CurrencyRef struct {
		Value string `json:"value"`
	} `json:"CurrencyRef"`
	TxnDate   string `json:"TxnDate"`
	DueDate   string `json:"DueDate"`
	CustomerRef struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"CustomerRef"`
	CustomerEmail struct {
		Address string `json:"Address"`
	} `json:"BillEmail"`
	Line []struct {
		Description string `json:"Description"`
		Amount      float64 `json:"Amount"`
		SalesItemLineDetail struct {
			Qty       float64 `json:"Qty"`
			UnitPrice float64 `json:"UnitPrice"`
			TaxCodeRef struct {
				Value string `json:"value"`
			} `json:"TaxCodeRef"`
		} `json:"SalesItemLineDetail"`
	} `json:"Line"`
}

// Invoice is the import view of a QuickBooks invoice.
type Invoice struct {
	ID            string
	SyncToken     string
	Number        string
	Currency      string
	IssuedAt      time.Time
	DueAt         *time.Time
	CustomerName  string
	CustomerEmail string
	Lines         []Line
}

// Line is a QuickBooks line item.
type Line struct {
	Description string
	Quantity    float64
	AmountMinor int64
	TaxRate     float64
}

// Supplier is the seller on mirrored invoices.
type Supplier struct {
	Name  string
	Email string
}

// ImportOptions configures QuickBooks → domain mapping.
type ImportOptions struct {
	Supplier Supplier
}

// ToDomain maps a QuickBooks invoice to invoice.Invoice.
func (s Invoice) ToDomain(opts ImportOptions) *invoice.Invoice {
	num := s.Number
	if num == "" {
		num = s.ID
	}
	inv := invoice.New(num)
	inv.IssuedAt = s.IssuedAt.UTC()
	if s.DueAt != nil {
		inv.DueAt = s.DueAt.UTC()
	}
	inv.Seller = invoice.Party{Name: opts.Supplier.Name, Email: opts.Supplier.Email}
	inv.Buyer = invoice.Party{Name: s.CustomerName, Email: s.CustomerEmail}
	cur := strings.ToUpper(s.Currency)
	if cur == "" {
		cur = "USD"
	}
	for _, line := range s.Lines {
		inv.AddLine(invoice.NewLineItem(
			line.Description,
			line.Quantity,
			invoice.NewMoney(line.AmountMinor, cur),
			line.TaxRate,
		))
	}
	if inv.Metadata == nil {
		inv.Metadata = map[string]string{}
	}
	inv.Metadata[InvoiceIDMeta] = s.ID
	return inv
}

func mapAPIInvoice(inv apiInvoice) Invoice {
	issued := parseQBDate(inv.TxnDate)
	var due *time.Time
	if d := parseQBDate(inv.DueDate); !d.IsZero() {
		due = &d
	}
	cur := "USD"
	if inv.CurrencyRef.Value != "" {
		cur = strings.ToUpper(inv.CurrencyRef.Value)
	}
	lines := make([]Line, 0, len(inv.Line))
	for _, l := range inv.Line {
		if l.Description == "" && l.Amount == 0 {
			continue
		}
		qty := l.SalesItemLineDetail.Qty
		if qty <= 0 {
			qty = 1
		}
		unitMinor := int64(l.SalesItemLineDetail.UnitPrice * 100)
		if unitMinor == 0 {
			unitMinor = int64(l.Amount * 100 / qty)
		}
		lines = append(lines, Line{
			Description: l.Description,
			Quantity:    qty,
			AmountMinor: unitMinor,
		})
	}
	name := inv.CustomerRef.Name
	return Invoice{
		ID:            inv.Id,
		SyncToken:     inv.SyncToken,
		Number:        inv.DocNumber,
		Currency:      cur,
		IssuedAt:      issued,
		DueAt:         due,
		CustomerName:  name,
		CustomerEmail: inv.CustomerEmail.Address,
		Lines:         lines,
	}
}

func parseQBDate(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return t.UTC()
	}
	return time.Time{}
}

package xero

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wiederin/go-invoicer/invoice"
)

var xeroDateRE = regexp.MustCompile(`\d+`)

// Client calls the Xero Accounting API.
type Client struct {
	HTTPClient   *http.Client
	AccessToken  string
	TenantID     string
}

// FetchInvoice retrieves an ACCREC invoice by Xero InvoiceID (UUID).
func (c *Client) FetchInvoice(ctx context.Context, invoiceID string) (Invoice, error) {
	if c.AccessToken == "" || c.TenantID == "" {
		return Invoice{}, fmt.Errorf("xero: access token and tenant required")
	}
	url := fmt.Sprintf("https://api.xero.com/api.xro/2.0/Invoices/%s", invoiceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Invoice{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Xero-Tenant-Id", c.TenantID)
	req.Header.Set("Accept", "application/json")
	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	res, err := hc.Do(req)
	if err != nil {
		return Invoice{}, fmt.Errorf("xero: fetch invoice: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return Invoice{}, fmt.Errorf("xero: fetch invoice: %s", strings.TrimSpace(string(body)))
	}
	var payload struct {
		Invoices []apiInvoice `json:"Invoices"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return Invoice{}, err
	}
	if len(payload.Invoices) == 0 {
		return Invoice{}, fmt.Errorf("xero: invoice not found")
	}
	return mapAPIInvoice(payload.Invoices[0]), nil
}

type apiInvoice struct {
	InvoiceID     string `json:"InvoiceID"`
	InvoiceNumber string `json:"InvoiceNumber"`
	CurrencyCode  string `json:"CurrencyCode"`
	Date          string `json:"Date"`
	DueDate       string `json:"DueDate"`
	Contact       struct {
		Name         string `json:"Name"`
		EmailAddress string `json:"EmailAddress"`
	} `json:"Contact"`
	LineItems []struct {
		Description string  `json:"Description"`
		Quantity    float64 `json:"Quantity"`
		UnitAmount  float64 `json:"UnitAmount"`
		TaxAmount   float64 `json:"TaxAmount"`
		LineAmount  float64 `json:"LineAmount"`
	} `json:"LineItems"`
}

// Invoice is the import view of a Xero invoice.
type Invoice struct {
	ID            string
	Number        string
	Currency      string
	IssuedAt      time.Time
	DueAt         *time.Time
	CustomerName  string
	CustomerEmail string
	Lines         []Line
}

// Line is a Xero line item.
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

// ImportOptions configures Xero → domain mapping.
type ImportOptions struct {
	Supplier Supplier
}

// ToDomain maps a Xero invoice to invoice.Invoice.
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
		cur = "CHF"
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
	issued := parseXeroDate(inv.Date)
	var due *time.Time
	if d := parseXeroDate(inv.DueDate); !d.IsZero() {
		due = &d
	}
	cur := strings.ToUpper(inv.CurrencyCode)
	if cur == "" {
		cur = "CHF"
	}
	lines := make([]Line, 0, len(inv.LineItems))
	for _, l := range inv.LineItems {
		qty := l.Quantity
		if qty <= 0 {
			qty = 1
		}
		unitMinor := int64(l.UnitAmount * 100)
		if unitMinor == 0 && l.LineAmount != 0 {
			unitMinor = int64(l.LineAmount * 100 / qty)
		}
		taxRate := 0.0
		if l.LineAmount > 0 && l.TaxAmount > 0 {
			taxRate = l.TaxAmount / l.LineAmount
		}
		lines = append(lines, Line{
			Description: l.Description,
			Quantity:    qty,
			AmountMinor: unitMinor,
			TaxRate:     taxRate,
		})
	}
	return Invoice{
		ID:            inv.InvoiceID,
		Number:        inv.InvoiceNumber,
		Currency:      cur,
		IssuedAt:      issued,
		DueAt:         due,
		CustomerName:  inv.Contact.Name,
		CustomerEmail: inv.Contact.EmailAddress,
		Lines:         lines,
	}
}

func parseXeroDate(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return t.UTC()
	}
	m := xeroDateRE.FindString(raw)
	if m == "" {
		return time.Time{}
	}
	ms, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.UnixMilli(ms).UTC()
}

package quickbooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	domain "github.com/wiederin/go-invoicer/invoice"
)

// ExportOptions configures domain → QuickBooks export.
type ExportOptions struct {
	CustomerID string
}

// PushInvoice creates or updates a QuickBooks draft invoice from a domain invoice.
func (c *Client) PushInvoice(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (string, error) {
	if c.AccessToken == "" || c.RealmID == "" {
		return "", fmt.Errorf("quickbooks: access token and realm required")
	}
	if inv == nil {
		return "", fmt.Errorf("quickbooks: invoice is nil")
	}
	if err := inv.Validate(); err != nil {
		return "", err
	}

	qboID := inv.Metadata[InvoiceIDMeta]
	if qboID == "" {
		var err error
		qboID, err = c.createInvoice(ctx, inv, opts)
		if err != nil {
			return "", err
		}
	} else if err := c.updateInvoice(ctx, qboID, inv, opts); err != nil {
		return "", err
	}

	if inv.Metadata == nil {
		inv.Metadata = map[string]string{}
	}
	inv.Metadata[InvoiceIDMeta] = qboID
	return qboID, nil
}

func (c *Client) createInvoice(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (string, error) {
	payload, err := c.exportPayload(ctx, inv, opts)
	if err != nil {
		return "", err
	}
	resBody, err := c.apiRequest(ctx, http.MethodPost, c.invoiceURL(""), payload)
	if err != nil {
		return "", err
	}
	var out struct {
		Invoice struct {
			Id string `json:"Id"`
		} `json:"Invoice"`
	}
	if err := json.Unmarshal(resBody, &out); err != nil {
		return "", err
	}
	if out.Invoice.Id == "" {
		return "", fmt.Errorf("quickbooks: create invoice: empty response")
	}
	return out.Invoice.Id, nil
}

func (c *Client) updateInvoice(ctx context.Context, qboID string, inv *domain.Invoice, opts ExportOptions) error {
	existing, err := c.FetchInvoice(ctx, qboID)
	if err != nil {
		return err
	}
	payload, err := c.exportPayload(ctx, inv, opts)
	if err != nil {
		return err
	}
	payload["Id"] = qboID
	payload["SyncToken"] = existing.SyncToken
	_, err = c.apiRequest(ctx, http.MethodPost, c.invoiceURL(qboID), payload)
	return err
}

func (c *Client) exportPayload(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (map[string]any, error) {
	customerID := opts.CustomerID
	if customerID == "" {
		id, err := c.findOrCreateCustomer(ctx, inv.Buyer)
		if err != nil {
			return nil, err
		}
		customerID = id
	}
	if customerID == "" {
		return nil, fmt.Errorf("quickbooks: buyer name or email required for customer")
	}

	lines := make([]map[string]any, 0, len(inv.LineItems))
	for _, line := range inv.LineItems {
		qty := line.Quantity
		if qty <= 0 {
			qty = 1
		}
		unit := float64(line.UnitPrice.Amount) / 100.0
		amount := unit * qty
		lines = append(lines, map[string]any{
			"DetailType":  "SalesItemLineDetail",
			"Amount":      amount,
			"Description": line.Description,
			"SalesItemLineDetail": map[string]any{
				"Qty":       qty,
				"UnitPrice": unit,
			},
		})
	}

	payload := map[string]any{
		"DocNumber": inv.Number,
		"TxnDate":   formatQBDate(inv.IssuedAt),
		"CustomerRef": map[string]string{
			"value": customerID,
		},
		"Line": lines,
	}
	if !inv.DueAt.IsZero() {
		payload["DueDate"] = formatQBDate(inv.DueAt)
	}
	cur := strings.ToUpper(inv.Currency())
	if cur != "" {
		payload["CurrencyRef"] = map[string]string{"value": cur}
	}
	return payload, nil
}

func (c *Client) findOrCreateCustomer(ctx context.Context, buyer domain.Party) (string, error) {
	if buyer.Email != "" {
		q := fmt.Sprintf("select Id from Customer where PrimaryEmailAddr = '%s'", escapeQBString(buyer.Email))
		body, err := c.query(ctx, q)
		if err == nil {
			var out struct {
				QueryResponse struct {
					Customer []struct {
						Id string `json:"Id"`
					} `json:"Customer"`
				} `json:"QueryResponse"`
			}
			if json.Unmarshal(body, &out) == nil && len(out.QueryResponse.Customer) > 0 {
				return out.QueryResponse.Customer[0].Id, nil
			}
		}
	}
	if buyer.Name == "" {
		return "", fmt.Errorf("quickbooks: buyer name required")
	}
	resBody, err := c.apiRequest(ctx, http.MethodPost, c.customerURL(), map[string]any{
		"DisplayName": buyer.Name,
		"PrimaryEmailAddr": map[string]string{
			"Address": buyer.Email,
		},
	})
	if err != nil {
		return "", err
	}
	var out struct {
		Customer struct {
			Id string `json:"Id"`
		} `json:"Customer"`
	}
	if err := json.Unmarshal(resBody, &out); err != nil {
		return "", err
	}
	if out.Customer.Id == "" {
		return "", fmt.Errorf("quickbooks: create customer: empty response")
	}
	return out.Customer.Id, nil
}

func (c *Client) query(ctx context.Context, q string) ([]byte, error) {
	base := c.APIBase
	if base == "" {
		base = qboProdAPI
	}
	reqURL := fmt.Sprintf("%s/v3/company/%s/query?query=%s&minorversion=65", base, c.RealmID, url.QueryEscape(q))
	return c.apiRequest(ctx, http.MethodGet, reqURL, nil)
}

func (c *Client) invoiceURL(id string) string {
	base := c.APIBase
	if base == "" {
		base = qboProdAPI
	}
	if id == "" {
		return fmt.Sprintf("%s/v3/company/%s/invoice?minorversion=65", base, c.RealmID)
	}
	return fmt.Sprintf("%s/v3/company/%s/invoice?minorversion=65", base, c.RealmID)
}

func (c *Client) customerURL() string {
	base := c.APIBase
	if base == "" {
		base = qboProdAPI
	}
	return fmt.Sprintf("%s/v3/company/%s/customer?minorversion=65", base, c.RealmID)
}

func (c *Client) apiRequest(ctx context.Context, method, reqURL string, payload map[string]any) ([]byte, error) {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}
	var r io.Reader
	if len(body) > 0 {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, reqURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	res, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("quickbooks: request: %w", err)
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("quickbooks: %s %s: %s", method, reqURL, strings.TrimSpace(string(resBody)))
	}
	return resBody, nil
}

func formatQBDate(t time.Time) string {
	if t.IsZero() {
		return time.Now().UTC().Format("2006-01-02")
	}
	return t.UTC().Format("2006-01-02")
}

func escapeQBString(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

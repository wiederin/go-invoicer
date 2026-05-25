package xero

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

// ExportOptions configures domain → Xero export.
type ExportOptions struct {
	ContactID string
}

// PushInvoice creates or updates a Xero ACCREC draft invoice from a domain invoice.
func (c *Client) PushInvoice(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (string, error) {
	if c.AccessToken == "" || c.TenantID == "" {
		return "", fmt.Errorf("xero: access token and tenant required")
	}
	if inv == nil {
		return "", fmt.Errorf("xero: invoice is nil")
	}
	if err := inv.Validate(); err != nil {
		return "", err
	}

	xeroID := inv.Metadata[InvoiceIDMeta]
	if xeroID == "" {
		var err error
		xeroID, err = c.createInvoice(ctx, inv, opts)
		if err != nil {
			return "", err
		}
	} else if err := c.updateInvoice(ctx, xeroID, inv, opts); err != nil {
		return "", err
	}

	if inv.Metadata == nil {
		inv.Metadata = map[string]string{}
	}
	inv.Metadata[InvoiceIDMeta] = xeroID
	return xeroID, nil
}

func (c *Client) createInvoice(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (string, error) {
	payload, err := c.exportPayload(ctx, inv, opts)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(map[string]any{"Invoices": []any{payload}})
	if err != nil {
		return "", err
	}
	resBody, err := c.apiRequest(ctx, http.MethodPost, "https://api.xero.com/api.xro/2.0/Invoices", body)
	if err != nil {
		return "", err
	}
	var out struct {
		Invoices []struct {
			InvoiceID string `json:"InvoiceID"`
		} `json:"Invoices"`
	}
	if err := json.Unmarshal(resBody, &out); err != nil {
		return "", err
	}
	if len(out.Invoices) == 0 || out.Invoices[0].InvoiceID == "" {
		return "", fmt.Errorf("xero: create invoice: empty response")
	}
	return out.Invoices[0].InvoiceID, nil
}

func (c *Client) updateInvoice(ctx context.Context, xeroID string, inv *domain.Invoice, opts ExportOptions) error {
	payload, err := c.exportPayload(ctx, inv, opts)
	if err != nil {
		return err
	}
	payload["InvoiceID"] = xeroID
	body, err := json.Marshal(map[string]any{"Invoices": []any{payload}})
	if err != nil {
		return err
	}
	_, err = c.apiRequest(ctx, http.MethodPost, "https://api.xero.com/api.xro/2.0/Invoices/"+url.PathEscape(xeroID), body)
	return err
}

func (c *Client) exportPayload(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (map[string]any, error) {
	contactID := opts.ContactID
	if contactID == "" {
		id, err := c.findOrCreateContact(ctx, inv.Buyer)
		if err != nil {
			return nil, err
		}
		contactID = id
	}
	if contactID == "" {
		return nil, fmt.Errorf("xero: buyer name or email required for contact")
	}

	lineItems := make([]map[string]any, 0, len(inv.LineItems))
	for _, line := range inv.LineItems {
		qty := line.Quantity
		if qty <= 0 {
			qty = 1
		}
		unit := float64(line.UnitPrice.Amount) / 100.0
		lineItems = append(lineItems, map[string]any{
			"Description": line.Description,
			"Quantity":    qty,
			"UnitAmount":  unit,
		})
	}

	payload := map[string]any{
		"Type":          "ACCREC",
		"Status":        "DRAFT",
		"Contact":       map[string]string{"ContactID": contactID},
		"Date":          formatXeroDate(inv.IssuedAt),
		"InvoiceNumber": inv.Number,
		"LineItems":     lineItems,
	}
	if !inv.DueAt.IsZero() {
		payload["DueDate"] = formatXeroDate(inv.DueAt)
	}
	cur := strings.ToUpper(inv.Currency())
	if cur != "" {
		payload["CurrencyCode"] = cur
	}
	return payload, nil
}

func (c *Client) findOrCreateContact(ctx context.Context, buyer domain.Party) (string, error) {
	if buyer.Email != "" {
		where := fmt.Sprintf(`EmailAddress=="%s"`, strings.ReplaceAll(buyer.Email, `"`, `\"`))
		reqURL := "https://api.xero.com/api.xro/2.0/Contacts?where=" + url.QueryEscape(where)
		body, err := c.apiRequest(ctx, http.MethodGet, reqURL, nil)
		if err == nil {
			var out struct {
				Contacts []struct {
					ContactID string `json:"ContactID"`
				} `json:"Contacts"`
			}
			if json.Unmarshal(body, &out) == nil && len(out.Contacts) > 0 {
				return out.Contacts[0].ContactID, nil
			}
		}
	}
	if buyer.Name == "" {
		return "", fmt.Errorf("xero: buyer name required")
	}
	createBody, _ := json.Marshal(map[string]any{
		"Contacts": []map[string]any{{
			"Name":         buyer.Name,
			"EmailAddress": buyer.Email,
		}},
	})
	resBody, err := c.apiRequest(ctx, http.MethodPost, "https://api.xero.com/api.xro/2.0/Contacts", createBody)
	if err != nil {
		return "", err
	}
	var out struct {
		Contacts []struct {
			ContactID string `json:"ContactID"`
		} `json:"Contacts"`
	}
	if err := json.Unmarshal(resBody, &out); err != nil {
		return "", err
	}
	if len(out.Contacts) == 0 || out.Contacts[0].ContactID == "" {
		return "", fmt.Errorf("xero: create contact: empty response")
	}
	return out.Contacts[0].ContactID, nil
}

func (c *Client) apiRequest(ctx context.Context, method, reqURL string, body []byte) ([]byte, error) {
	var r io.Reader
	if len(body) > 0 {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, reqURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Xero-Tenant-Id", c.TenantID)
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
		return nil, fmt.Errorf("xero: request: %w", err)
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("xero: %s %s: %s", method, reqURL, strings.TrimSpace(string(resBody)))
	}
	return resBody, nil
}

func formatXeroDate(t time.Time) string {
	if t.IsZero() {
		return time.Now().UTC().Format("2006-01-02")
	}
	return t.UTC().Format("2006-01-02")
}

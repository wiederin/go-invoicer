package stripe

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	stripefile "github.com/stripe/stripe-go/v81/file"
	stripeinvoice "github.com/stripe/stripe-go/v81/invoice"
	stripeitem "github.com/stripe/stripe-go/v81/invoiceitem"

	domain "github.com/wiederin/go-invoicer/invoice"
)

const (
	// MetaGoInvoicerPDFURL points to the branded PDF hosted by go-invoicer.
	MetaGoInvoicerPDFURL = "go_invoicer_pdf_url"
	// MetaGoInvoicerJobID links a Stripe invoice to a render job.
	MetaGoInvoicerJobID = "go_invoicer_render_job_id"
	// MetaGoInvoicerFileID is a Stripe File object containing the branded PDF.
	MetaGoInvoicerFileID = "go_invoicer_pdf_file_id"
)

func (c *Client) withKey() error {
	if c.apiKey == "" {
		return fmt.Errorf("stripe: API key is required")
	}
	stripe.Key = c.apiKey
	return nil
}

// LinkBrandedPDF updates a Stripe invoice with metadata and a custom field for the branded PDF.
func (c *Client) LinkBrandedPDF(ctx context.Context, stripeInvoiceID string, pdfURL, jobID, fileID string) error {
	if err := c.withKey(); err != nil {
		return err
	}
	meta := map[string]string{
		MetaGoInvoicerPDFURL: pdfURL,
		MetaGoInvoicerJobID:  jobID,
	}
	if fileID != "" {
		meta[MetaGoInvoicerFileID] = fileID
	}
	params := &stripe.InvoiceParams{
		Metadata: meta,
	}
	params.Context = ctx
	if pdfURL != "" {
		params.CustomFields = []*stripe.InvoiceCustomFieldParams{
			{Name: stripe.String("Branded invoice"), Value: stripe.String(pdfURL)},
		}
		params.Footer = stripe.String("Branded PDF available via go-invoicer.")
	}
	_, err := stripeinvoice.Update(stripeInvoiceID, params)
	if err != nil {
		return fmt.Errorf("stripe: link branded pdf: %w", err)
	}
	return nil
}

// UploadPDF uploads a branded invoice PDF to Stripe Files and returns the file ID and link URL.
func (c *Client) UploadPDF(ctx context.Context, filename string, pdf []byte) (fileID, fileURL string, err error) {
	if err := c.withKey(); err != nil {
		return "", "", err
	}
	if filename == "" {
		filename = "invoice.pdf"
	}
	params := &stripe.FileParams{
		Purpose:    stripe.String(string(stripe.FilePurposeDisputeEvidence)),
		Filename:   stripe.String(filename),
		FileReader: bytes.NewReader(pdf),
	}
	params.Context = ctx
	f, err := stripefile.New(params)
	if err != nil {
		return "", "", fmt.Errorf("stripe: upload pdf: %w", err)
	}
	url := ""
	if f.Links != nil && len(f.Links.Data) > 0 {
		url = f.Links.Data[0].URL
	}
	return f.ID, url, nil
}

// ExportOptions configures domain → Stripe export.
type ExportOptions struct {
	CustomerID string
	Supplier   Supplier
}

// PushInvoice creates or updates a Stripe draft invoice from a domain invoice.
// Returns the Stripe invoice ID.
func (c *Client) PushInvoice(ctx context.Context, inv *domain.Invoice, opts ExportOptions) (string, error) {
	if err := c.withKey(); err != nil {
		return "", err
	}
	if inv == nil {
		return "", fmt.Errorf("stripe: invoice is nil")
	}
	if err := inv.Validate(); err != nil {
		return "", err
	}

	stripeID := inv.Metadata[StripeInvoiceIDMeta]
	if stripeID != "" {
		return stripeID, c.LinkBrandedPDF(ctx, stripeID, inv.Metadata[MetaGoInvoicerPDFURL], inv.Metadata[MetaGoInvoicerJobID], inv.Metadata[MetaGoInvoicerFileID])
	}

	customerID := opts.CustomerID
	if customerID == "" && inv.Buyer.Email != "" {
		custParams := &stripe.CustomerParams{
			Email: stripe.String(inv.Buyer.Email),
			Name:  stripe.String(inv.Buyer.Name),
		}
		custParams.Context = ctx
		cust, err := customer.New(custParams)
		if err != nil {
			return "", fmt.Errorf("stripe: create customer: %w", err)
		}
		customerID = cust.ID
	}
	if customerID == "" {
		return "", fmt.Errorf("stripe: customer_id or buyer email required")
	}

	currency := strings.ToLower(inv.Currency())
	if currency == "" {
		currency = "chf"
	}

	invParams := &stripe.InvoiceParams{
		Customer:         stripe.String(customerID),
		CollectionMethod: stripe.String(string(stripe.InvoiceCollectionMethodSendInvoice)),
		Metadata: map[string]string{
			"go_invoicer_number": inv.Number,
		},
	}
	invParams.Context = ctx
	stInv, err := stripeinvoice.New(invParams)
	if err != nil {
		return "", fmt.Errorf("stripe: create invoice: %w", err)
	}

	for _, line := range inv.LineItems {
		amount := line.LineTotal()
		if amount == 0 {
			amount = int64(float64(line.UnitPrice.Amount) * line.Quantity)
		}
		itemParams := &stripe.InvoiceItemParams{
			Customer:    stripe.String(customerID),
			Invoice:     stripe.String(stInv.ID),
			Currency:    stripe.String(currency),
			Description: stripe.String(line.Description),
		}
		itemParams.Context = ctx
		if line.Quantity > 0 {
			qty := int64(line.Quantity)
			if qty < 1 {
				qty = 1
			}
			itemParams.Quantity = stripe.Int64(qty)
			itemParams.UnitAmount = stripe.Int64(int64(float64(amount) / line.Quantity))
		} else {
			itemParams.Amount = stripe.Int64(amount)
		}
		if _, err := stripeitem.New(itemParams); err != nil {
			return "", fmt.Errorf("stripe: add line item: %w", err)
		}
	}

	if inv.Metadata == nil {
		inv.Metadata = map[string]string{}
	}
	inv.Metadata[StripeInvoiceIDMeta] = stInv.ID
	return stInv.ID, nil
}

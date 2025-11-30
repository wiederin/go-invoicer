package render

import (
        "bytes"
        "fmt"
        "io"
        "strings"

        "github.com/go-pdf/fpdf"
        "github.com/wiederin/go-invoicer/invoice"
        "github.com/wiederin/go-invoicer/template"
)

type PDFRenderer interface {
        RenderHTML(html string) ([]byte, error)
}

type Options struct {
        PageSize    string
        Orientation string
        MarginTop   float64
        MarginRight float64
        MarginBottom float64
        MarginLeft  float64
        FontFamily  string
        FontSize    float64
}

func DefaultOptions() Options {
        return Options{
                PageSize:     "A4",
                Orientation:  "P",
                MarginTop:    10,
                MarginRight:  10,
                MarginBottom: 10,
                MarginLeft:   10,
                FontFamily:   "Arial",
                FontSize:     10,
        }
}

type Engine struct {
        Templates *template.Manager
        Options   Options
}

func NewEngine(templates *template.Manager) *Engine {
        return &Engine{
                Templates: templates,
                Options:   DefaultOptions(),
        }
}

func (e *Engine) WithOptions(opts Options) *Engine {
        e.Options = opts
        return e
}

func (e *Engine) RenderInvoice(inv *invoice.Invoice, templateName string) ([]byte, error) {
        if err := inv.Validate(); err != nil {
                return nil, fmt.Errorf("invalid invoice: %w", err)
        }

        data := e.prepareTemplateData(inv)

        html, err := e.Templates.RenderHTML(templateName, data)
        if err != nil {
                return nil, fmt.Errorf("failed to render template: %w", err)
        }

        return e.htmlToPDF(html)
}

func (e *Engine) prepareTemplateData(inv *invoice.Invoice) map[string]any {
        return map[string]any{
                "Invoice":       inv,
                "SubTotal":      inv.SubTotal(),
                "TotalDiscount": inv.TotalDiscount(),
                "TotalNet":      inv.TotalNet(),
                "TotalTax":      inv.TotalTax(),
                "TotalGross":    inv.TotalGross(),
                "TaxBreakdown":  inv.TaxBreakdown(),
        }
}

func (e *Engine) htmlToPDF(html string) ([]byte, error) {
        pdf := fpdf.New(e.Options.Orientation, "mm", e.Options.PageSize, "")
        pdf.SetMargins(e.Options.MarginLeft, e.Options.MarginTop, e.Options.MarginRight)
        pdf.SetAutoPageBreak(true, e.Options.MarginBottom)
        pdf.AddPage()
        pdf.SetFont(e.Options.FontFamily, "", e.Options.FontSize)

        htmlContent := parseHTMLContent(html)
        _, lineHt := pdf.GetFontSize()
        
        for _, block := range htmlContent {
                switch block.Type {
                case "h1":
                        pdf.SetFont(e.Options.FontFamily, "B", 24)
                        pdf.MultiCell(0, 12, block.Text, "", "", false)
                        pdf.Ln(5)
                case "h2":
                        pdf.SetFont(e.Options.FontFamily, "B", 18)
                        pdf.MultiCell(0, 10, block.Text, "", "", false)
                        pdf.Ln(4)
                case "h3":
                        pdf.SetFont(e.Options.FontFamily, "B", 14)
                        pdf.MultiCell(0, 8, block.Text, "", "", false)
                        pdf.Ln(3)
                case "p", "div":
                        pdf.SetFont(e.Options.FontFamily, "", e.Options.FontSize)
                        pdf.MultiCell(0, lineHt*1.5, block.Text, "", "", false)
                        pdf.Ln(2)
                case "strong", "b":
                        pdf.SetFont(e.Options.FontFamily, "B", e.Options.FontSize)
                        pdf.MultiCell(0, lineHt*1.5, block.Text, "", "", false)
                default:
                        pdf.SetFont(e.Options.FontFamily, "", e.Options.FontSize)
                        if block.Text != "" {
                                pdf.MultiCell(0, lineHt*1.5, block.Text, "", "", false)
                        }
                }
        }

        var buf bytes.Buffer
        if err := pdf.Output(&buf); err != nil {
                return nil, fmt.Errorf("failed to generate PDF: %w", err)
        }

        return buf.Bytes(), nil
}

type htmlBlock struct {
        Type string
        Text string
}

func parseHTMLContent(html string) []htmlBlock {
        var blocks []htmlBlock
        var currentTag string
        var currentText strings.Builder
        inTag := false
        isClosing := false

        for _, r := range html {
                if r == '<' {
                        if currentText.Len() > 0 {
                                text := strings.TrimSpace(currentText.String())
                                if text != "" {
                                        blocks = append(blocks, htmlBlock{Type: currentTag, Text: text})
                                }
                                currentText.Reset()
                        }
                        inTag = true
                        isClosing = false
                        continue
                }
                if r == '>' {
                        inTag = false
                        continue
                }
                if inTag {
                        if r == '/' {
                                isClosing = true
                        } else if !isClosing && r != ' ' {
                                currentTag += string(r)
                        }
                        if r == ' ' || isClosing {
                                currentTag = strings.ToLower(currentTag)
                        }
                        continue
                }
                currentText.WriteRune(r)
        }

        if currentText.Len() > 0 {
                text := strings.TrimSpace(currentText.String())
                if text != "" {
                        blocks = append(blocks, htmlBlock{Type: currentTag, Text: text})
                }
        }

        return blocks
}

type SimpleRenderer struct {
        Options Options
}

func NewSimpleRenderer() *SimpleRenderer {
        return &SimpleRenderer{
                Options: DefaultOptions(),
        }
}

func (r *SimpleRenderer) RenderInvoice(inv *invoice.Invoice) ([]byte, error) {
        if err := inv.Validate(); err != nil {
                return nil, fmt.Errorf("invalid invoice: %w", err)
        }

        pdf := fpdf.New(r.Options.Orientation, "mm", r.Options.PageSize, "")
        pdf.SetMargins(r.Options.MarginLeft, r.Options.MarginTop, r.Options.MarginRight)
        pdf.SetAutoPageBreak(true, r.Options.MarginBottom)
        pdf.AddPage()

        r.renderHeader(pdf, inv)
        r.renderParties(pdf, inv)
        r.renderLineItems(pdf, inv)
        r.renderTotals(pdf, inv)
        r.renderFooter(pdf, inv)

        var buf bytes.Buffer
        if err := pdf.Output(&buf); err != nil {
                return nil, fmt.Errorf("failed to generate PDF: %w", err)
        }

        return buf.Bytes(), nil
}

func (r *SimpleRenderer) renderHeader(pdf *fpdf.Fpdf, inv *invoice.Invoice) {
        pdf.SetFont("Arial", "B", 24)
        pdf.Cell(0, 15, "INVOICE")
        pdf.Ln(20)

        pdf.SetFont("Arial", "", 10)
        pdf.Cell(95, 6, fmt.Sprintf("Invoice Number: %s", inv.Number))
        pdf.Ln(6)
        pdf.Cell(95, 6, fmt.Sprintf("Issue Date: %s", inv.IssueDate.Format("2006-01-02")))
        pdf.Ln(6)
        pdf.Cell(95, 6, fmt.Sprintf("Due Date: %s", inv.DueDate.Format("2006-01-02")))
        pdf.Ln(15)
}

func (r *SimpleRenderer) renderParties(pdf *fpdf.Fpdf, inv *invoice.Invoice) {
        startY := pdf.GetY()

        pdf.SetFont("Arial", "B", 10)
        pdf.Cell(95, 6, "From:")
        pdf.Ln(6)
        pdf.SetFont("Arial", "", 10)
        pdf.Cell(95, 5, inv.Supplier.Name)
        pdf.Ln(5)
        for _, line := range inv.Supplier.Address.Lines() {
                pdf.Cell(95, 5, line)
                pdf.Ln(5)
        }
        if inv.Supplier.VATID != "" {
                pdf.Cell(95, 5, fmt.Sprintf("VAT ID: %s", inv.Supplier.VATID))
                pdf.Ln(5)
        }

        pdf.SetXY(105, startY)
        pdf.SetFont("Arial", "B", 10)
        pdf.Cell(95, 6, "To:")
        pdf.Ln(6)
        pdf.SetX(105)
        pdf.SetFont("Arial", "", 10)
        pdf.Cell(95, 5, inv.Customer.Name)
        pdf.Ln(5)
        pdf.SetX(105)
        for _, line := range inv.Customer.Address.Lines() {
                pdf.Cell(95, 5, line)
                pdf.Ln(5)
                pdf.SetX(105)
        }
        if inv.Customer.VATID != "" {
                pdf.Cell(95, 5, fmt.Sprintf("VAT ID: %s", inv.Customer.VATID))
                pdf.Ln(5)
        }

        pdf.Ln(15)
}

func (r *SimpleRenderer) renderLineItems(pdf *fpdf.Fpdf, inv *invoice.Invoice) {
        pdf.SetFont("Arial", "B", 9)
        pdf.SetFillColor(240, 240, 240)
        pdf.CellFormat(80, 8, "Description", "1", 0, "", true, 0, "")
        pdf.CellFormat(20, 8, "Qty", "1", 0, "C", true, 0, "")
        pdf.CellFormat(30, 8, "Unit Price", "1", 0, "R", true, 0, "")
        pdf.CellFormat(20, 8, "Tax %", "1", 0, "C", true, 0, "")
        pdf.CellFormat(40, 8, "Amount", "1", 0, "R", true, 0, "")
        pdf.Ln(8)

        pdf.SetFont("Arial", "", 9)
        for _, item := range inv.LineItems {
                pdf.CellFormat(80, 7, item.Description, "1", 0, "", false, 0, "")
                pdf.CellFormat(20, 7, item.Quantity.String(), "1", 0, "C", false, 0, "")
                pdf.CellFormat(30, 7, item.UnitPrice.String(), "1", 0, "R", false, 0, "")
                pdf.CellFormat(20, 7, item.TaxRate.String()+"%", "1", 0, "C", false, 0, "")
                pdf.CellFormat(40, 7, item.GrossAmount().String(), "1", 0, "R", false, 0, "")
                pdf.Ln(7)
        }
        pdf.Ln(5)
}

func (r *SimpleRenderer) renderTotals(pdf *fpdf.Fpdf, inv *invoice.Invoice) {
        pdf.SetX(120)
        pdf.SetFont("Arial", "", 10)
        pdf.Cell(40, 6, "Subtotal:")
        pdf.Cell(30, 6, inv.SubTotal().String())
        pdf.Ln(6)

        if !inv.TotalDiscount().IsZero() {
                pdf.SetX(120)
                pdf.Cell(40, 6, "Discount:")
                pdf.Cell(30, 6, "-"+inv.TotalDiscount().String())
                pdf.Ln(6)
        }

        pdf.SetX(120)
        pdf.Cell(40, 6, "Tax:")
        pdf.Cell(30, 6, inv.TotalTax().String())
        pdf.Ln(6)

        pdf.SetFont("Arial", "B", 11)
        pdf.SetX(120)
        pdf.Cell(40, 8, "Total:")
        pdf.Cell(30, 8, inv.TotalGross().String())
        pdf.Ln(15)
}

func (r *SimpleRenderer) renderFooter(pdf *fpdf.Fpdf, inv *invoice.Invoice) {
        pdf.SetFont("Arial", "", 9)
        if inv.Notes != "" {
                pdf.SetFont("Arial", "B", 9)
                pdf.Cell(0, 6, "Notes:")
                pdf.Ln(6)
                pdf.SetFont("Arial", "", 9)
                pdf.MultiCell(0, 5, inv.Notes, "", "", false)
                pdf.Ln(5)
        }

        if inv.Terms != "" {
                pdf.SetFont("Arial", "B", 9)
                pdf.Cell(0, 6, "Payment Terms:")
                pdf.Ln(6)
                pdf.SetFont("Arial", "", 9)
                pdf.MultiCell(0, 5, inv.Terms, "", "", false)
        }
}

func (r *SimpleRenderer) RenderToWriter(inv *invoice.Invoice, w io.Writer) error {
        data, err := r.RenderInvoice(inv)
        if err != nil {
                return err
        }
        _, err = w.Write(data)
        return err
}

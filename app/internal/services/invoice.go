package services

import (
        "bytes"
        "database/sql"
        "encoding/json"
        "fmt"

        "github.com/wiederin/go-invoicer"
        "github.com/wiederin/go-invoicer-app/internal/models"
        "github.com/wiederin/go-invoicer/components"
        "github.com/wiederin/go-invoicer/constants"
)

type InvoiceService struct {
        db *sql.DB
}

func NewInvoiceService(db *sql.DB) *InvoiceService {
        return &InvoiceService{db: db}
}

type InvoiceRequest struct {
        InvoiceNumber   string `json:"invoiceNumber"`
        CompanyName     string `json:"companyName"`
        CompanyAddress  string `json:"companyAddress"`
        CompanyCity     string `json:"companyCity"`
        CompanyPostal   string `json:"companyPostal"`
        CustomerName    string `json:"customerName"`
        CustomerAddress string `json:"customerAddress"`
        CustomerCity    string `json:"customerCity"`
        CustomerPostal  string `json:"customerPostal"`
        CustomerCountry string `json:"customerCountry"`
        Description     string `json:"description"`
        Items           []struct {
                Name        string `json:"name"`
                Description string `json:"description"`
                UnitCost    string `json:"unitCost"`
                Quantity    string `json:"quantity"`
        } `json:"items"`
        Notes string `json:"notes"`
}

func (s *InvoiceService) GenerateInvoice(userID int, req *InvoiceRequest) ([]byte, error) {
        doc, err := invoicer.Init(&components.Config{
                TextInvoiceType:   constants.Invoice,
                TextRefTitle:      "Invoice No.",
                CurrencySymbol:    "$",
                CurrencyPrecision: 2,
                CurrencyDecimal:   ".",
                CurrencyThousand:  ",",
        })
        if err != nil {
                return nil, fmt.Errorf("failed to initialize invoice: %w", err)
        }

        doc.SetVersion(req.InvoiceNumber)

        doc.SetHeader(&components.HeaderFooter{
                Text:       "<center>INVOICE</center>",
                Pagination: true,
        })

        doc.SetFooter(&components.HeaderFooter{
                Text:       "<center>Thank you for your business</center>",
                Pagination: true,
        })

        doc.SetCompany(&components.Contact{
                Name: req.CompanyName,
                Address: &components.Address{
                        Address:    req.CompanyAddress,
                        PostalCode: req.CompanyPostal,
                        City:       req.CompanyCity,
                },
        })

        doc.SetCustomer(&components.Contact{
                Name: req.CustomerName,
                Address: &components.Address{
                        Address:    req.CustomerAddress,
                        PostalCode: req.CustomerPostal,
                        City:       req.CustomerCity,
                        Country:    req.CustomerCountry,
                },
        })

        if req.Description != "" {
                doc.SetDescription(req.Description)
        }

        for _, item := range req.Items {
                doc.AppendItem(&components.Item{
                        Name:        item.Name,
                        Description: item.Description,
                        UnitCost:    item.UnitCost,
                        Quantity:    item.Quantity,
                })
        }

        if req.Notes != "" {
                doc.SetNotes(req.Notes)
        }

        pdf, err := invoicer.Build(doc)
        if err != nil {
                return nil, fmt.Errorf("failed to build invoice: %w", err)
        }

        var buf bytes.Buffer
        err = pdf.Output(&buf)
        if err != nil {
                return nil, fmt.Errorf("failed to generate PDF: %w", err)
        }

        metadata, _ := json.Marshal(req)
        _, err = s.db.Exec(`
                INSERT INTO invoices (user_id, invoice_number, company_name, customer_name, metadata)
                VALUES ($1, $2, $3, $4, $5)
        `, userID, req.InvoiceNumber, req.CompanyName, req.CustomerName, metadata)

        if err != nil {
                return nil, fmt.Errorf("failed to save invoice record: %w", err)
        }

        return buf.Bytes(), nil
}

func (s *InvoiceService) GetUserInvoices(userID int, limit int) ([]models.Invoice, error) {
        if limit <= 0 {
                limit = 50
        }

        rows, err := s.db.Query(`
                SELECT id, user_id, invoice_number, company_name, customer_name, created_at
                FROM invoices
                WHERE user_id = $1
                ORDER BY created_at DESC
                LIMIT $2
        `, userID, limit)
        if err != nil {
                return nil, fmt.Errorf("failed to get invoices: %w", err)
        }
        defer rows.Close()

        var invoices []models.Invoice
        for rows.Next() {
                var inv models.Invoice
                err := rows.Scan(&inv.ID, &inv.UserID, &inv.InvoiceNumber, &inv.CompanyName, &inv.CustomerName, &inv.CreatedAt)
                if err != nil {
                        return nil, fmt.Errorf("failed to scan invoice: %w", err)
                }
                invoices = append(invoices, inv)
        }

        return invoices, nil
}

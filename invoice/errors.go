package invoice

import "errors"

// Validation errors returned by invoice and line item validation.
var (
	ErrMissingInvoiceNumber = errors.New("invoice number is required")
	ErrMissingIssueDate     = errors.New("issue date is required")
	ErrMissingDueDate       = errors.New("due date is required")
	ErrDueDateBeforeIssue   = errors.New("due date cannot be before issue date")
	ErrMissingSupplier      = errors.New("supplier is required")
	ErrMissingCustomer      = errors.New("customer is required")
	ErrMissingPartyName     = errors.New("party name is required")
	ErrNoLineItems          = errors.New("at least one line item is required")
	ErrMissingDescription   = errors.New("line item description is required")
	ErrInvalidQuantity      = errors.New("quantity must be positive")
	ErrInvalidUnitPrice     = errors.New("unit price cannot be negative")
	ErrInvalidTaxRate       = errors.New("tax rate cannot be negative")
	ErrCurrencyMismatch     = errors.New("all amounts must use the same currency")
)

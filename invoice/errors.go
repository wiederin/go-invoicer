package invoice

import "errors"

var (
	ErrMissingNumber     = errors.New("invoice: number is required")
	ErrMissingSeller     = errors.New("invoice: seller name is required")
	ErrMissingBuyer      = errors.New("invoice: buyer name is required")
	ErrMissingLineItems  = errors.New("invoice: at least one line item is required")
	ErrDueBeforeIssue    = errors.New("invoice: due date must be on or after issue date")
	ErrMissingCurrency   = errors.New("invoice: currency is required on line items")
	ErrCurrencyMismatch  = errors.New("invoice: all line items must use the same currency")
	ErrMissingRelated    = errors.New("invoice: credit note requires related_number (source invoice)")
)

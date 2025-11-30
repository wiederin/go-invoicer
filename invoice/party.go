package invoice

// Address represents a postal address.
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// IsEmpty returns true if the address has no fields set.
func (a Address) IsEmpty() bool {
	return a.Street == "" && a.City == "" && a.PostalCode == "" && a.Country == ""
}

// Lines returns the address as formatted lines for display.
func (a Address) Lines() []string {
	var lines []string
	if a.Street != "" {
		lines = append(lines, a.Street)
	}
	cityLine := ""
	if a.PostalCode != "" {
		cityLine = a.PostalCode + " "
	}
	if a.City != "" {
		cityLine += a.City
	}
	if a.State != "" {
		cityLine += ", " + a.State
	}
	if cityLine != "" {
		lines = append(lines, cityLine)
	}
	if a.Country != "" {
		lines = append(lines, a.Country)
	}
	return lines
}

// Party represents a business entity (supplier or customer).
type Party struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
	Email   string  `json:"email,omitempty"`
	Phone   string  `json:"phone,omitempty"`
	VATID   string  `json:"vat_id,omitempty"`
	IBAN    string  `json:"iban,omitempty"`
}

// Validate checks that the party has all required fields.
func (p Party) Validate() error {
	if p.Name == "" {
		return ErrMissingPartyName
	}
	return nil
}

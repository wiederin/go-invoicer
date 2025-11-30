package invoice

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

func (a Address) IsEmpty() bool {
	return a.Street == "" && a.City == "" && a.PostalCode == "" && a.Country == ""
}

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

type Party struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
	Email   string  `json:"email,omitempty"`
	Phone   string  `json:"phone,omitempty"`
	VATID   string  `json:"vat_id,omitempty"`
	IBAN    string  `json:"iban,omitempty"`
}

func (p Party) Validate() error {
	if p.Name == "" {
		return ErrMissingPartyName
	}
	return nil
}

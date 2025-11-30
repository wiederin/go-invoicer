package components

// UnicodeTranslateFunc
type UnicodeTranslateFunc func(string) string

// Config for Document
type Config struct {
    // fonts
    Font     string `default:"Helvetica"`
    BoldFont string `default:"Helvetica"`

    // colors
    DarkBgColor   []int `default:"[212,212,212]" json:"dark_bg_color,omitempty"`
    GreyBgColor   []int `default:"[232,232,232]" json:"grey_bg_color,omitempty"`
    BaseTextColor []int `default:"[35,35,35]" json:"base_text_color,omitempty"`
    GreyTextColor []int `default:"[82,82,82]" json:"grey_text_color,omitempty"`

    // currency
    CurrencySymbol    string `default:"€ " json:"currency_symbol,omitempty"`
    CurrencyPrecision int    `default:"2" json:"currency_precision,omitempty"`
    CurrencyDecimal   string `default:"." json:"currency_decimal,omitempty"`
    CurrencyThousand  string `default:" " json:"currency_thousand,omitempty"`

    // text
    TextInvoiceType      string `default:"INVOICE" json:"text_type_invoice,omitempty"`
    TextRefTitle         string `default:"Ref." json:"text_ref_title,omitempty"`
	  TextDateTitle        string `default:"Date" json:"text_date_title,omitempty"`
    TextVersionTitle     string `default:"Version" json:"text_version_title,omitempty"`
    TextTotalTotal      string `default:"SUBTOTAL" json:"text_total_total,omitempty"`
    TextTotalDiscounted string `default:"TOTAL DISCOUNTED" json:"text_total_discounted,omitempty"`
    TextTotalTax        string `default:"VAT applied" json:"text_total_tax,omitempty"`
    TextTotalWithTax    string `default:"TOTAL WITH TAX" json:"text_total_with_tax,omitempty"`
    TextTotalWithCommission    string `default:"TOTAL" json:"text_total_commission,omitempty"`
    TextFixedTransferFees    string `default:"Fixed Transfer Fees" json:"text_fixed_transfer_fees,omitempty"`
    TextPaymentTermTitle string `default:"Payment term" json:"text_payment_term_title,omitempty"`
    TextItemsNameTitle     string `default:"Name" json:"text_items_name_title,omitempty"`
    TextItemsUnitCostTitle string `default:"Ad Spend" json:"text_items_unit_cost_title,omitempty"`
    TextItemsCommissionTitle string `default:"Commission" json:"text_items_commission_title,omitempty"`
    TextItemsTotalTTCTitle string `default:"Total" json:"text_items_total_ttc_title,omitempty"`

    // fees et
    FixedFee bool `json:"fixed_fee,omitempty"`

    UnicodeTranslateFunc UnicodeTranslateFunc
}

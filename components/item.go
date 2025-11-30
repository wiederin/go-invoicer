package components

import (
    "github.com/wiederin/go-invoicer/constants"

    "fmt"

    "github.com/shopspring/decimal"
)

type Item struct {
    Name        string    `json:"name,omitempty" validate:"required"`
    Description string    `json:"description,omitempty"`
    UnitCost     string    `json:"cost,omitempty"`
    Commission  string    `json:"commission,omitempty"`
    Quantity    string    `json:"quantity,omitempty"`
    Tax         *Tax      `json:"tax,omitempty"`
    Discount    *Discount `json:"discount,omitempty"`

    _unit_cost decimal.Decimal
    _commission decimal.Decimal
}

func (doc *Document) AddItems() {
  	doc.drawsTableTitles()

  	doc.Pdf.SetX(10)
  	doc.Pdf.SetY(doc.Pdf.GetY() + 8) // was 8
  	doc.Pdf.SetFont(doc.Config.Font, "", 8)

  	for i := 0; i < len(doc.Items); i++ {
  		item := doc.Items[i]

  		// Check item tax
  		if item.Tax == nil {
  			item.Tax = doc.DefaultTax
  		}

  		// Append to pdf
  		item.appendColTo(doc.Config, doc)

  		if doc.Pdf.GetY() > constants.MaxPageHeight {
  			// Add page
  			doc.Pdf.AddPage()
  			doc.drawsTableTitles()
  			doc.Pdf.SetFont(doc.Config.Font, "", 8)
  		}

  		doc.Pdf.SetX(10)
  		doc.Pdf.SetY(doc.Pdf.GetY() + 6)
  	}
  }

func (i *Item) Prepare() error {
    unit_cost, err := decimal.NewFromString(i.UnitCost)
    if err != nil {
        return err
    }
    i._unit_cost = unit_cost

    // commission
    commission, err := decimal.NewFromString(i.Commission)
    if err != nil {
        return err
    }
    i._commission = commission

    // Tax
    if i.Tax != nil {
        if err := i.Tax.Prepare(); err != nil {
            return err
        }
    }

    // Discount
    if i.Discount != nil {
        if err := i.Discount.Prepare(); err != nil {
            return err
        }
    }

    return nil
}

// TotalWithoutTaxAndWithoutDiscount returns the total without tax and without discount
func (i *Item) TotalWithoutTaxAndWithoutDiscount() decimal.Decimal {
    //quantity, _ := decimal.NewFromString(i.Quantity)
    //price, _ := decimal.NewFromString(i.UnitCost)
    //total := price.Mul(quantity)
    total, _  := decimal.NewFromString(i.UnitCost)

    return total
}

// TotalWithoutTaxAndWithDiscount returns the total without tax and with discount
func (i *Item) TotalWithoutTaxAndWithDiscount() decimal.Decimal {
    total := i.TotalWithoutTaxAndWithoutDiscount()

    // Check discount
    if i.Discount != nil {
        dType, dNum := i.Discount.getDiscount()

        if dType == DiscountTypeAmount {
            total = total.Sub(dNum)
        } else {
            // Percent
            toSub := total.Mul(dNum.Div(decimal.NewFromFloat(100)))
            total = total.Sub(toSub)
        }
    }

    return total
}

// TotalWithTaxAndDiscount returns the total with tax and discount
func (i *Item) TotalWithTaxAndDiscount() decimal.Decimal {
    return i.TotalWithoutTaxAndWithDiscount().Add(i.TaxWithTotalDiscounted())
}

// TaxWithTotalDiscounted returns the tax with total discounted
func (i *Item) TaxWithTotalDiscounted() decimal.Decimal {
    result := decimal.NewFromFloat(0)

    if i.Tax == nil {
        return result
    }

    totalHT := i.TotalWithoutTaxAndWithDiscount()
    taxType, taxAmount := i.Tax.getTax()

    if taxType == TaxTypeAmount {
        result = taxAmount
    } else {
        divider := decimal.NewFromFloat(100)
        result = totalHT.Mul(taxAmount.Div(divider))
    }

    return result
}

func (i *Item) TotalWithCommission() decimal.Decimal {
    dCost := i.TotalWithoutTaxAndWithoutDiscount()
    dAmount := dCost.Mul(i._commission.Div(decimal.NewFromFloat(100)))
    return i.TotalWithoutTaxAndWithoutDiscount().Add(dAmount)
}

// appendColTo document doc
func (i *Item) appendColTo(config *Config, doc *Document) {
    // Get base Y (top of line)
    baseY := doc.Pdf.GetY()

    // Name
    doc.Pdf.SetX(constants.ItemColNameOffset)
    doc.Pdf.MultiCell(
        constants.ItemColUnitPriceOffset-constants.ItemColNameOffset,
        3,
        doc.encodeString(i.Name),
        "",
        "",
        false,
    )

    // Description
    if len(i.Description) > 0 {
        doc.Pdf.SetX(constants.ItemColNameOffset)
        doc.Pdf.SetY(doc.Pdf.GetY() + 1)

        doc.Pdf.SetFont(doc.Config.Font, "", constants.SmallTextFontSize)
        doc.Pdf.SetTextColor(
            doc.Config.GreyTextColor[0],
            doc.Config.GreyTextColor[1],
            doc.Config.GreyTextColor[2],
        )

        doc.Pdf.MultiCell(
            constants.ItemColUnitPriceOffset-constants.ItemColNameOffset,
            3,
            doc.encodeString(i.Description),
            "",
            "",
            false,
        )

        // Reset font
        doc.Pdf.SetFont(doc.Config.Font, "", constants.BaseTextFontSize)
        doc.Pdf.SetTextColor(
            doc.Config.BaseTextColor[0],
            doc.Config.BaseTextColor[1],
            doc.Config.BaseTextColor[2],
        )
    }

    // Compute line height
    colHeight := doc.Pdf.GetY() - baseY

    // Unit price
    doc.Pdf.SetY(baseY)
    doc.Pdf.SetX(constants.ItemColUnitPriceOffset)
    doc.Pdf.CellFormat(
        constants.ItemColQuantityOffset-constants.ItemColUnitPriceOffset,
        colHeight,
        doc.encodeString(doc.Accounting.FormatMoneyDecimal(i._unit_cost)),
        "0",
        0,
        "",
        false,
        0,
        "",
    )

    // Commission
    doc.Pdf.SetX(constants.ItemColCommissionOffset)
    var commissionTitle string
    var commissionDesc string
    commissionTitle = fmt.Sprintf("%s %s", i._commission, "%")
    // get amount from percent
    dCost := i.TotalWithoutTaxAndWithoutDiscount()
    dAmount := dCost.Mul(i._commission.Div(decimal.NewFromFloat(100)))
    commissionDesc = doc.Accounting.FormatMoneyDecimal(dAmount)
    // tax title
    // lastY := doc.pdf.GetY()
    doc.Pdf.CellFormat(
        constants.ItemColDiscountOffset-constants.ItemColTaxOffset,
        colHeight/2,
        doc.encodeString(commissionTitle),
        "0",
        0,
        "LB",
        false,
        0,
        "",
    )

    // tax desc
    doc.Pdf.SetXY(constants.ItemColCommissionOffset, baseY+(colHeight/2))
    doc.Pdf.SetFont(doc.Config.Font, "", constants.SmallTextFontSize)
    doc.Pdf.SetTextColor(
        doc.Config.GreyTextColor[0],
        doc.Config.GreyTextColor[1],
        doc.Config.GreyTextColor[2],
    )

    doc.Pdf.CellFormat(
        constants.ItemColDiscountOffset-constants.ItemColTaxOffset,
        colHeight/2,
        doc.encodeString(commissionDesc),
        "0",
        0,
        "LT",
        false,
        0,
        "",
    )

    // reset font and y
    doc.Pdf.SetFont(doc.Config.Font, "", constants.BaseTextFontSize)
    doc.Pdf.SetTextColor(
        doc.Config.BaseTextColor[0],
        doc.Config.BaseTextColor[1],
        doc.Config.BaseTextColor[2],
    )
    doc.Pdf.SetY(baseY)

    // TOTAL TTC
    doc.Pdf.SetX(constants.ItemColTotalTTCOffset)
    doc.Pdf.CellFormat(
        190-constants.ItemColTotalTTCOffset,
        colHeight,
        doc.encodeString(doc.Accounting.FormatMoneyDecimal(i.TotalWithCommission())),
        "0",
        0,
        "",
        false,
        0,
        "",
    )

    // Set Y for next line
    doc.Pdf.SetY(baseY + colHeight)

    // Discount
    doc.Pdf.SetX(constants.ItemColDiscountOffset)
    if i.Discount == nil {
        /* remove to display discount
        doc.pdf.CellFormat(
            ItemColTotalTTCOffset-ItemColDiscountOffset,
            colHeight,
            doc.encodeString("--"),
            "0",
            0,
            "",
            false,
            0,
            "",
        )
        */
    } else {
        // If discount
        discountType, discountAmount := i.Discount.getDiscount()
        var discountTitle string
        var discountDesc string

        dCost := i.TotalWithoutTaxAndWithoutDiscount()
        if discountType == DiscountTypePercent {
            discountTitle = fmt.Sprintf("%s %s", discountAmount, doc.encodeString("%"))

            // get amount from percent
            dAmount := dCost.Mul(discountAmount.Div(decimal.NewFromFloat(100)))
            discountDesc = fmt.Sprintf("-%s", doc.Accounting.FormatMoneyDecimal(dAmount))
        } else {
            discountTitle = fmt.Sprintf("%s %s", discountAmount, doc.encodeString("€"))

            // get percent from amount
            dPerc := discountAmount.Mul(decimal.NewFromFloat(100))
            dPerc = dPerc.Div(dCost)
            discountDesc = fmt.Sprintf("-%s %%", dPerc.StringFixed(2))
        }

        // discount title
        // lastY := doc.pdf.GetY()
        doc.Pdf.CellFormat(
            constants.ItemColTotalTTCOffset-constants.ItemColDiscountOffset,
            colHeight/2,
            doc.encodeString(discountTitle),
            "0",
            0,
            "LB",
            false,
            0,
            "",
        )

        // discount desc
        doc.Pdf.SetXY(constants.ItemColDiscountOffset, baseY+(colHeight/2))
        doc.Pdf.SetFont(doc.Config.Font, "", constants.SmallTextFontSize)
        doc.Pdf.SetTextColor(
            doc.Config.GreyTextColor[0],
            doc.Config.GreyTextColor[1],
            doc.Config.GreyTextColor[2],
        )

        doc.Pdf.CellFormat(
            constants.ItemColTotalTTCOffset-constants.ItemColDiscountOffset,
            colHeight/2,
            doc.encodeString(discountDesc),
            "0",
            0,
            "LT",
            false,
            0,
            "",
        )

        // reset font and y
        doc.Pdf.SetFont(doc.Config.Font, "", constants.BaseTextFontSize)
        doc.Pdf.SetTextColor(
            doc.Config.BaseTextColor[0],
            doc.Config.BaseTextColor[1],
            doc.Config.BaseTextColor[2],
        )
        doc.Pdf.SetY(baseY)
    }

    // Tax
    doc.Pdf.SetX(constants.ItemColTaxOffset)
    if i.Tax == nil {
        /* remove if want to show tax column
        // If no tax
        doc.pdf.CellFormat(
            ItemColDiscountOffset-ItemColTaxOffset,
            colHeight,
            doc.encodeString("--"),
            "0",
            0,
            "",
            false,
            0,
            "",
        )
        */
    } else {
        // If tax
        taxType, taxAmount := i.Tax.getTax()
        var taxTitle string
        var taxDesc string

        if taxType == TaxTypePercent {
            taxTitle = fmt.Sprintf("%s %s", taxAmount, "%")
            // get amount from percent
            dCost := i.TotalWithoutTaxAndWithDiscount()
            dAmount := dCost.Mul(taxAmount.Div(decimal.NewFromFloat(100)))
            taxDesc = doc.Accounting.FormatMoneyDecimal(dAmount)
        } else {
            taxTitle = fmt.Sprintf("%s %s", doc.Accounting.Symbol, taxAmount)
            dCost := i.TotalWithoutTaxAndWithDiscount()
            dPerc := taxAmount.Mul(decimal.NewFromFloat(100))
            dPerc = dPerc.Div(dCost)
            // get percent from amount
            taxDesc = fmt.Sprintf("%s %%", dPerc.StringFixed(2))
        }

        // tax title
        // lastY := doc.pdf.GetY()
        doc.Pdf.CellFormat(
            constants.ItemColDiscountOffset-constants.ItemColTaxOffset,
            colHeight/2,
            doc.encodeString(taxTitle),
            "0",
            0,
            "LB",
            false,
            0,
            "",
        )

        // tax desc
        doc.Pdf.SetXY(constants.ItemColTaxOffset, baseY+(colHeight/2))
        doc.Pdf.SetFont(doc.Config.Font, "", constants.SmallTextFontSize)
        doc.Pdf.SetTextColor(
            doc.Config.GreyTextColor[0],
            doc.Config.GreyTextColor[1],
            doc.Config.GreyTextColor[2],
        )

        doc.Pdf.CellFormat(
            constants.ItemColDiscountOffset-constants.ItemColTaxOffset,
            colHeight/2,
            doc.encodeString(taxDesc),
            "0",
            0,
            "LT",
            false,
            0,
            "",
        )

        // reset font and y
        doc.Pdf.SetFont(doc.Config.Font, "", constants.BaseTextFontSize)
        doc.Pdf.SetTextColor(
            doc.Config.BaseTextColor[0],
            doc.Config.BaseTextColor[1],
            doc.Config.BaseTextColor[2],
        )
        doc.Pdf.SetY(baseY)
    }
    /*
    // TOTAL TTC
    doc.pdf.SetX(ItemColTotalTTCOffset)
    doc.pdf.CellFormat(
        190-ItemColTotalTTCOffset,
        colHeight,
        doc.encodeString(doc.ac.FormatMoneyDecimal(i.TotalWithCommission())),
        "0",
        0,
        "",
        false,
        0,
        "",
    ) */

    // Set Y for next line
    doc.Pdf.SetY(baseY + colHeight)
}

// drawsTableTitles in document
func (doc *Document) drawsTableTitles() {
	// Draw table titles
	doc.Pdf.SetX(10)
	doc.Pdf.SetY(doc.Pdf.GetY() + 20)// was 5
	doc.Pdf.SetFont(doc.Config.BoldFont, "B", 8)

	// Draw rec
	doc.Pdf.SetFillColor(doc.Config.GreyBgColor[0], doc.Config.GreyBgColor[1], doc.Config.GreyBgColor[2])
	doc.Pdf.Rect(10, doc.Pdf.GetY(), 190, 6, "F")

	// Name
	doc.Pdf.SetX(constants.ItemColNameOffset)
	doc.Pdf.CellFormat(
		constants.ItemColUnitPriceOffset-constants.ItemColNameOffset,
		6,
		doc.encodeString(doc.Config.TextItemsNameTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Unit price
	doc.Pdf.SetX(constants.ItemColUnitPriceOffset)
	doc.Pdf.CellFormat(
		constants.ItemColQuantityOffset-constants.ItemColUnitPriceOffset,
		6,
		doc.encodeString(doc.Config.TextItemsUnitCostTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)

	// Commission
	doc.Pdf.SetX(constants.ItemColCommissionOffset)
	doc.Pdf.CellFormat(
		constants.ItemColQuantityOffset-constants.ItemColUnitPriceOffset,
		6,
		doc.encodeString(doc.Config.TextItemsCommissionTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)
	/* remove to show tax column
 	// Tax
	doc.pdf.SetX(ItemColTaxOffset)
	doc.pdf.CellFormat(
		ItemColDiscountOffset-ItemColTaxOffset,
		6,
		doc.encodeString(doc.Options.TextItemsTaxTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)
	*/

	/* remove to show discount column
	// Discount
	doc.pdf.SetX(ItemColDiscountOffset)
	doc.pdf.CellFormat(
		ItemColTotalTTCOffset-ItemColDiscountOffset,
		6,
		doc.encodeString(doc.Options.TextItemsDiscountTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	) */

	// TOTAL TTC
	doc.Pdf.SetX(constants.ItemColTotalTTCOffset)
	doc.Pdf.CellFormat(
		190-constants.ItemColTotalTTCOffset,
		6,
		doc.encodeString(doc.Config.TextItemsTotalTTCTitle),
		"0",
		0,
		"",
		false,
		0,
		"",
	)
}

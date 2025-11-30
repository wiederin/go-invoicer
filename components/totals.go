package components


import (
	"github.com/wiederin/go-invoicer/constants"
	"github.com/shopspring/decimal"
	"bytes"
)

func (doc *Document) TotalWithoutTaxAndWithoutDocumentDiscount() decimal.Decimal {
	total := decimal.NewFromInt(0)

	for _, item := range doc.Items {
		total = total.Add(item.TotalWithoutTaxAndWithDiscount())
	}

	return total
}

func (doc *Document) TotalWithoutTax() decimal.Decimal {
	total := doc.TotalWithoutTaxAndWithoutDocumentDiscount()

	if doc.Discount != nil {
		discountType, discountNumber := doc.Discount.getDiscount()

		if discountType == DiscountTypeAmount {
			total = total.Sub(discountNumber)
		} else {
			toSub := total.Mul(discountNumber.Div(decimal.NewFromFloat(100)))
			total = total.Sub(toSub)
		}
	}

	return total
}

func (doc *Document) TotalWithCommission() decimal.Decimal {
	totalWithCommission := decimal.NewFromFloat(0)
	for _, item := range doc.Items {
		totalWithCommission = totalWithCommission.Add(item.TotalWithCommission())
	}

	return totalWithCommission
}

func (doc *Document) TotalWithCommissionAndFixedFee() decimal.Decimal {
	totalWithCommissionAndFixedFee := doc.TotalWithCommission()
	if doc.Config.FixedFee {
		totalWithCommissionAndFixedFee = totalWithCommissionAndFixedFee.Add(decimal.NewFromFloat(15))
	}

	return totalWithCommissionAndFixedFee
}

func (doc *Document) TotalWithTax() decimal.Decimal {
	totalWithoutTax := doc.TotalWithoutTax()
	tax := doc.Tax()

	return totalWithoutTax.Add(tax)
}

func (doc *Document) Tax() decimal.Decimal {
	totalWithoutTaxAndWithoutDocDiscount := doc.TotalWithoutTaxAndWithoutDocumentDiscount()
	totalTax := decimal.NewFromFloat(0)

	if doc.Discount == nil {
		for _, item := range doc.Items {
			totalTax = totalTax.Add(item.TaxWithTotalDiscounted())
		}
	} else {
		discountType, discountAmount := doc.Discount.getDiscount()
		discountPercent := discountAmount
		if discountType == DiscountTypeAmount {
			discountPercent = discountAmount.Mul(decimal.NewFromFloat(100)).Div(totalWithoutTaxAndWithoutDocDiscount)
		}

		for _, item := range doc.Items {
			if item.Tax != nil {
				taxType, taxAmount := item.Tax.getTax()
				if taxType == TaxTypeAmount {
					totalTax = totalTax.Add(taxAmount)
				} else {
					// Else, remove doc discount % from item total without tax and item discount
					itemTotal := item.TotalWithoutTaxAndWithDiscount()
					toSub := discountPercent.Mul(itemTotal).Div(decimal.NewFromFloat(100))
					itemTotalDiscounted := itemTotal.Sub(toSub)

					// Then recompute tax on itemTotalDiscounted
					itemTaxDiscounted := taxAmount.Mul(itemTotalDiscounted).Div(decimal.NewFromFloat(100))

					totalTax = totalTax.Add(itemTaxDiscounted)
				}
			}
		}
	}

	return totalTax
}

func (doc *Document) AddTotal() {
	doc.Pdf.SetY(doc.Pdf.GetY() + 10)
	doc.Pdf.SetFont(doc.Config.Font, "", constants.LargeTextFontSize)
	doc.Pdf.SetTextColor(
		doc.Config.BaseTextColor[0],
		doc.Config.BaseTextColor[1],
		doc.Config.BaseTextColor[2],
	)

	// Draw TOTAL HT title
	doc.Pdf.SetX(120)
	doc.Pdf.SetFillColor(doc.Config.DarkBgColor[0], doc.Config.DarkBgColor[1], doc.Config.DarkBgColor[2])
	doc.Pdf.Rect(120, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(38, 10, doc.encodeString(doc.Config.TextTotalTotal), "0", 0, "R", false, 0, "")

	// Draw TOTAL HT amount
	doc.Pdf.SetX(162)
	doc.Pdf.SetFillColor(doc.Config.GreyBgColor[0], doc.Config.GreyBgColor[1], doc.Config.GreyBgColor[2])
	doc.Pdf.Rect(160, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.Accounting.FormatMoneyDecimal(doc.TotalWithCommission())),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	if doc.Discount != nil {
		baseY := doc.Pdf.GetY() + 10

		// Draw discounted title
		doc.Pdf.SetXY(120, baseY)
		doc.Pdf.SetFillColor(doc.Config.DarkBgColor[0], doc.Config.DarkBgColor[1], doc.Config.DarkBgColor[2])
		doc.Pdf.Rect(120, doc.Pdf.GetY(), 40, 15, "F")

		// title
		doc.Pdf.CellFormat(38, 7.5, doc.encodeString(doc.Config.TextTotalDiscounted), "0", 0, "BR", false, 0, "")

		// description
		doc.Pdf.SetXY(120, baseY+7.5)
		doc.Pdf.SetFont(doc.Config.Font, "", constants.BaseTextFontSize)
		doc.Pdf.SetTextColor(
			doc.Config.GreyTextColor[0],
			doc.Config.GreyTextColor[1],
			doc.Config.GreyTextColor[2],
		)

		var descString bytes.Buffer
		discountType, discountAmount := doc.Discount.getDiscount()
		if discountType == DiscountTypePercent {
			descString.WriteString("-")
			descString.WriteString(discountAmount.String())
			descString.WriteString(" % / -")
			descString.WriteString(doc.Accounting.FormatMoneyDecimal(
				doc.TotalWithoutTaxAndWithoutDocumentDiscount().Sub(doc.TotalWithoutTax())),
			)
		} else {
			descString.WriteString("-")
			descString.WriteString(doc.Accounting.FormatMoneyDecimal(discountAmount))
			descString.WriteString(" / -")
			descString.WriteString(
				discountAmount.Mul(decimal.NewFromFloat(100)).Div(doc.TotalWithoutTaxAndWithoutDocumentDiscount()).StringFixed(2),
			)
			descString.WriteString(" %")
		}

		doc.Pdf.CellFormat(38, 7.5, doc.encodeString(descString.String()), "0", 0, "TR", false, 0, "")

		doc.Pdf.SetFont(doc.Config.Font, "", constants.LargeTextFontSize)
		doc.Pdf.SetTextColor(
			doc.Config.BaseTextColor[0],
			doc.Config.BaseTextColor[1],
			doc.Config.BaseTextColor[2],
		)

		// Draw discount amount
		doc.Pdf.SetY(baseY)
		doc.Pdf.SetX(162)
		doc.Pdf.SetFillColor(doc.Config.GreyBgColor[0], doc.Config.GreyBgColor[1], doc.Config.GreyBgColor[2])
		doc.Pdf.Rect(160, doc.Pdf.GetY(), 40, 15, "F")
		doc.Pdf.CellFormat(
			40,
			15,
			doc.encodeString(doc.Accounting.FormatMoneyDecimal(doc.TotalWithoutTax())),
			"0",
			0,
			"L",
			false,
			0,
			"",
		)
		doc.Pdf.SetY(doc.Pdf.GetY() + 15)
	} else {
		doc.Pdf.SetY(doc.Pdf.GetY() + 10)
	}

	// Draw tax title
	doc.Pdf.SetX(120)
	doc.Pdf.SetFillColor(doc.Config.DarkBgColor[0], doc.Config.DarkBgColor[1], doc.Config.DarkBgColor[2])
	doc.Pdf.Rect(120, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(38, 10, doc.encodeString(doc.Config.TextTotalTax), "0", 0, "R", false, 0, "")

	// Draw tax amount
	doc.Pdf.SetX(162)
	doc.Pdf.SetFillColor(doc.Config.GreyBgColor[0], doc.Config.GreyBgColor[1], doc.Config.GreyBgColor[2])
	doc.Pdf.Rect(160, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.Accounting.FormatMoneyDecimal(doc.Tax())),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Draw fixed transfer fees title
	doc.Pdf.SetY(doc.Pdf.GetY() + 10)
	doc.Pdf.SetX(120)
	doc.Pdf.SetFillColor(doc.Config.DarkBgColor[0], doc.Config.DarkBgColor[1], doc.Config.DarkBgColor[2])
	doc.Pdf.Rect(120, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(38, 10, doc.encodeString(doc.Config.TextFixedTransferFees), "0", 0, "R", false, 0, "")


	var fixedFeeDesc = "-"
	// Draw total with fixed transfer fees
	if doc.Config.FixedFee {
		fixedFeeDesc = doc.Config.CurrencySymbol + "15,00"
	}
	doc.Pdf.SetX(162)
	doc.Pdf.SetFillColor(doc.Config.GreyBgColor[0], doc.Config.GreyBgColor[1], doc.Config.GreyBgColor[2])
	doc.Pdf.Rect(160, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(
		40,
		10,
		doc.encodeString(fixedFeeDesc),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)

	// Draw total with commission title
	doc.Pdf.SetY(doc.Pdf.GetY() + 10)
	doc.Pdf.SetX(120)
	doc.Pdf.SetFont(doc.Config.BoldFont, "B", constants.LargeTextFontSize)
	doc.Pdf.SetFillColor(doc.Config.DarkBgColor[0], doc.Config.DarkBgColor[1], doc.Config.DarkBgColor[2])
	doc.Pdf.Rect(120, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(38, 10, doc.encodeString(doc.Config.TextTotalWithCommission), "0", 0, "R", false, 0, "")

	// Draw total with commission amount
	doc.Pdf.SetX(162)
	doc.Pdf.SetFillColor(doc.Config.GreyBgColor[0], doc.Config.GreyBgColor[1], doc.Config.GreyBgColor[2])
	doc.Pdf.Rect(160, doc.Pdf.GetY(), 40, 10, "F")
	doc.Pdf.CellFormat(
		40,
		10,
		doc.encodeString(doc.Accounting.FormatMoneyDecimal(doc.TotalWithCommissionAndFixedFee())),
		"0",
		0,
		"L",
		false,
		0,
		"",
	)
	doc.Pdf.SetFont(doc.Config.Font, "", constants.LargeTextFontSize)
}

package emailTemplates

import (
	"strings"
	"time"
	"unicode"
)

var storeCreditFinalValues [][]interface{}
var creditAppliedFinalValues [][]interface{}
var shippingTrackerFinalValues [][]interface{}

func StoreCreditFlushFinalValues() {
	storeCreditFinalValues = nil
}

func CreditAppliedFlushFinalValues() {
	creditAppliedFinalValues = nil
}

func ShippingTrackerFlushFinalValues() {
	shippingTrackerFinalValues = nil
}

func GetStoreCreditReport(creditEmail string, InternalDate string, EmailReceiver string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	StoreCreditAmount := ""
	StoreCreditLink := ""
	AmountStartIndex := strings.Index(creditEmail, "sent you a $")
	if AmountStartIndex != -1 {
		AmountEndIndex := strings.Index(creditEmail[AmountStartIndex:], "eStore Credit")
		if AmountEndIndex != -1 {
			StoreCreditAmount = creditEmail[AmountStartIndex+len("sent you a ") : AmountStartIndex+AmountEndIndex]
			StoreCreditAmount = strings.Replace(StoreCreditAmount, " USD", "", -1)
		}
	}
	StoreCreditLinkStartIndex := strings.Index(creditEmail, "Go to  ")
	if StoreCreditLinkStartIndex != -1 {
		StoreCreditLinkEndIndex := strings.Index(creditEmail[StoreCreditLinkStartIndex:], "to receive instructions")
		if StoreCreditLinkEndIndex != -1 {
			StoreCreditLink = creditEmail[StoreCreditLinkStartIndex+len("Go to  ") : StoreCreditLinkStartIndex+StoreCreditLinkEndIndex-2]
		}
	}
	row = append(row, StoreCreditLink)
	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, StoreCreditAmount)
	row = append(row, EmailReceiver)

	storeCreditFinalValues = append(storeCreditFinalValues, row)
}

func GetStoreCreditFinalValues() [][]interface{} {
	return storeCreditFinalValues
}

func GetCreditAppliedFinalValues() [][]interface{} {
	return creditAppliedFinalValues
}

func GetShippingTrackerFinalValues() [][]interface{} {
	return shippingTrackerFinalValues
}

func GetCreditAppliedReport(creditEmail string, InternalDate string, EmailReceiver string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	OrderNumber := ""
	OrderDate := ""
	InternetNumber := ""
	CreditAmount := ""
	StoreSKU := ""

	OrderNumberStartIndex := strings.Index(creditEmail, "Order Number")
	if OrderNumberStartIndex != -1 {
		OrderNumberEndIndex := strings.Index(creditEmail[OrderNumberStartIndex:], "</td>")
		if OrderNumberEndIndex != -1 {
			OrderNumber = creditEmail[OrderNumberStartIndex+len("Order Number") : OrderNumberStartIndex+OrderNumberEndIndex]
			OrderNumber = strings.Replace(OrderNumber, ":</b>", "", -1)
		}
	}
	OrderDateStartIndex := strings.Index(creditEmail, "Order Date")
	if OrderDateStartIndex != -1 {
		OrderDateEndIndex := strings.Index(creditEmail[OrderDateStartIndex:], "</td>")
		if OrderDateEndIndex != -1 {
			OrderDate = creditEmail[OrderDateStartIndex+len("Order Date") : OrderDateStartIndex+OrderDateEndIndex]
			OrderDate = strings.Replace(OrderDate, ":</b>", "", -1)
		}
	}

	InternetNumberStartIndex := strings.Index(creditEmail, "Internet #")
	if InternetNumberStartIndex != -1 {
		InternetNumberEndIndex := strings.Index(creditEmail[InternetNumberStartIndex:], "<br")
		if InternetNumberEndIndex != -1 {
			InternetNumber = creditEmail[InternetNumberStartIndex+len("Internet #") : InternetNumberStartIndex+InternetNumberEndIndex]
		}
	}

	CreditAmountStartIndex := strings.Index(creditEmail, "credit has been")
	if CreditAmountStartIndex != -1 {
		CreditAmountEndIndex := strings.Index(creditEmail[CreditAmountStartIndex-15:CreditAmountStartIndex], "A")
		if CreditAmountEndIndex != -1 {
			CreditAmount = creditEmail[CreditAmountStartIndex-15 : CreditAmountStartIndex+CreditAmountEndIndex+10]
			CreditAmount = strings.Replace(CreditAmount, "credit has been", "", -1)
			CreditAmount = strings.Replace(CreditAmount, "A", "", -1)
			CreditAmount = strings.Replace(CreditAmount, " USD", "", -1)
		}
	}

	StoreSKUStartIndex := strings.Index(creditEmail, "Store SKU #")
	if StoreSKUStartIndex != -1 {
		StoreSKUEndIndex := strings.Index(creditEmail[StoreSKUStartIndex:], "<br")
		if StoreSKUEndIndex != -1 {
			StoreSKU = creditEmail[StoreSKUStartIndex+len("Store SKU #") : StoreSKUStartIndex+StoreSKUEndIndex]
		}
	}
	row = append(row, InternetNumber)
	row = append(row, stripSpaces(OrderNumber))
	row = append(row, stripSpaces(CreditAmount))
	row = append(row, EmailReceiver)
	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, OrderDate)
	row = append(row, stripSpaces(StoreSKU))

	creditAppliedFinalValues = append(creditAppliedFinalValues, row)
}

func GetShippingTrackerReport(creditEmail string, InternalDate string, EmailReceiver string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	OrderNumber := ""
	OrderDate := ""
	InternetNumber := ""
	CreditAmount := ""
	StoreSKU := ""
	Carrier := ""
	TrackingNumber := ""

	OrderNumberStartIndex := strings.Index(creditEmail, "Order Number")
	if OrderNumberStartIndex != -1 {
		OrderNumberEndIndex := strings.Index(creditEmail[OrderNumberStartIndex:], "</td>")
		if OrderNumberEndIndex != -1 {
			Substring := creditEmail[OrderNumberStartIndex+len("Order Number") : OrderNumberStartIndex+OrderNumberEndIndex]
			OrderNumberStartIndex = strings.Index(Substring, "font-weight:bold;")
			if OrderNumberStartIndex != -1 {
				OrderNumberEndIndex = strings.Index(Substring[OrderNumberStartIndex:], "</span>")
				if OrderNumberEndIndex != -1 {
					OrderNumber = Substring[OrderNumberStartIndex+len("font-weight:bold;")+2 : OrderNumberStartIndex+OrderNumberEndIndex]
					OrderNumber = strings.Replace(OrderNumber, ":</b>", "", -1)
				}
			}

		}
	}
	OrderDateStartIndex := strings.Index(creditEmail, "Order Date")
	if OrderDateStartIndex != -1 {
		OrderDateEndIndex := strings.Index(creditEmail[OrderDateStartIndex:], "</td >")
		if OrderDateEndIndex != -1 {
			SubString := creditEmail[OrderDateStartIndex+len("Order Date") : OrderDateStartIndex+OrderDateEndIndex]
			OrderDateStartIndex := strings.Index(SubString, "font-weight:bold;")
			if OrderDateStartIndex != -1 {
				OrderDateEndIndex := strings.Index(SubString[OrderDateStartIndex:], "</span>")
				if OrderDateEndIndex != -1 {
					OrderDate = SubString[OrderDateStartIndex+len("font-weight:bold;")+2 : OrderDateStartIndex+OrderDateEndIndex]
				}
			}

		}
	}

	InternetNumberStartIndex := strings.Index(creditEmail, "Internet #")
	if InternetNumberStartIndex != -1 {
		InternetNumberEndIndex := strings.Index(creditEmail[InternetNumberStartIndex:], "<br")
		if InternetNumberEndIndex != -1 {
			InternetNumber = creditEmail[InternetNumberStartIndex+len("Internet #") : InternetNumberStartIndex+InternetNumberEndIndex]
		}
	}

	CreditAmountStartIndex := strings.Index(creditEmail, "Order Total")
	if OrderDateStartIndex != -1 {
		CreditAmountEndIndex := strings.Index(creditEmail[CreditAmountStartIndex:], "</tr>")
		if CreditAmountEndIndex != -1 {
			SubString := creditEmail[CreditAmountStartIndex+len("Order Total") : CreditAmountStartIndex+CreditAmountEndIndex]
			CreditAmountStartIndex := strings.Index(SubString, "display:block;text-align:right;")
			if CreditAmountStartIndex != -1 {
				CreditAmountEndIndex := strings.Index(SubString[CreditAmountStartIndex:], "</span>")
				if CreditAmountEndIndex != -1 {
					CreditAmount = SubString[CreditAmountStartIndex+len("display:block;text-align:right;")+2 : CreditAmountStartIndex+CreditAmountEndIndex]
					CreditAmount = strings.Replace(CreditAmount, " USD", "", -1)
				}
			}

		}
	}

	StoreSKUStartIndex := strings.Index(creditEmail, "Store SKU #")
	if StoreSKUStartIndex != -1 {
		StoreSKUEndIndex := strings.Index(creditEmail[StoreSKUStartIndex:], "<br")
		if StoreSKUEndIndex != -1 {
			StoreSKU = creditEmail[StoreSKUStartIndex+len("Store SKU #") : StoreSKUStartIndex+StoreSKUEndIndex]
		}
	}

	CarrierStartIndex := strings.Index(creditEmail, "Carrier")
	if CarrierStartIndex != -1 {
		CarrierEndIndex := strings.Index(creditEmail[CarrierStartIndex:], "<br>")
		if CarrierEndIndex != -1 {
			SubString := creditEmail[CarrierStartIndex+len("Carrier") : CarrierStartIndex+CarrierEndIndex]
			CarrierStartIndex := strings.Index(SubString, "uppercase")
			if CarrierStartIndex != -1 {
				CarrierEndIndex := strings.Index(SubString[CarrierStartIndex:], "</b>")
				if CarrierEndIndex != -1 {
					Carrier = SubString[CarrierStartIndex+len("uppercase")+2 : CarrierStartIndex+CarrierEndIndex]
				}
			}

		}
	}

	TrackingStartIndex := strings.Index(creditEmail, "Tracking Number:")
	if TrackingStartIndex != -1 {
		TrackingEndIndex := strings.Index(creditEmail[TrackingStartIndex:], "</a>")
		if TrackingEndIndex != -1 {
			SubString := creditEmail[TrackingStartIndex+len("Tracking Number:") : TrackingStartIndex+TrackingEndIndex]
			TrackingStartIndex := strings.Index(SubString, ";color:#333")
			if TrackingStartIndex != -1 {
				TrackingEndIndex := strings.Index(SubString[TrackingStartIndex:], "<br/>")
				if TrackingEndIndex != -1 {
					TrackingNumber = SubString[TrackingStartIndex+len(";color:#333")+2 : TrackingStartIndex+TrackingEndIndex]
				}
			}

		}
	}

	row = append(row, InternetNumber)
	row = append(row, stripSpaces(OrderNumber))
	row = append(row, stripSpaces(CreditAmount))
	row = append(row, EmailReceiver)
	row = append(row, stripSpaces(TrackingNumber))
	row = append(row, Carrier)

	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, OrderDate)
	row = append(row, stripSpaces(StoreSKU))

	shippingTrackerFinalValues = append(shippingTrackerFinalValues, row)
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

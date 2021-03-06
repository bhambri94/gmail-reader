package emailTemplates

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

var storeCreditFinalValues [][]interface{}
var creditAppliedFinalValues [][]interface{}

func StoreCreditFlushFinalValues() {
	storeCreditFinalValues = nil
}

func CreditAppliedFlushFinalValues() {
	creditAppliedFinalValues = nil
}

// func ShippingTrackerFlushFinalValues() {
// 	shippingTrackerFinalValues = nil
// }

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

// func GetShippingTrackerFinalValues() [][]interface{} {
// 	return shippingTrackerFinalValues
// }

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

func GetShippingTrackerReport(creditEmail string, InternalDate string, EmailReceiver string) []interface{} {
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
	Address := ""
	Quantity := ""
	ItemName := ""

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
	if CreditAmountStartIndex != -1 {
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

	AddressStartIndex := strings.Index(creditEmail, "inside-store-detail")
	if AddressStartIndex != -1 {
		AddressEndIndex := strings.Index(creditEmail[AddressStartIndex:], "</div>")
		if AddressEndIndex != -1 {
			SubString := creditEmail[AddressStartIndex+len("inside-store-detail") : AddressStartIndex+AddressEndIndex+len("</div>")]
			AddressStartIndex := strings.Index(SubString, "<br/>")
			if AddressStartIndex != -1 {
				AddressEndIndex := strings.Index(SubString[AddressStartIndex:], "</div>")
				if AddressEndIndex != -1 {
					SubString2 := SubString[AddressStartIndex+len("<br/>") : AddressStartIndex+AddressEndIndex+len("</div>")]
					AddressStartIndex := strings.Index(SubString2, "<br/>")
					if AddressStartIndex != -1 {
						AddressEndIndex := strings.Index(SubString2[AddressStartIndex+4:], "<br/>")
						if AddressEndIndex != -1 {
							Address = SubString2[AddressStartIndex+len("<br/>") : AddressStartIndex+AddressEndIndex]
							Address = strings.Replace(Address, "<span>", "", -1)
							Address = strings.Replace(Address, "span style=", "", -1)
							Address = strings.Replace(Address, "text-transform:capitalize;", "", -1)
							Address = strings.Replace(Address, "</span>", "", -1)
							Address = strings.Replace(Address, "</div>", "", -1)
							Address = strings.Replace(Address, "<", "", -1)
							Address = strings.Replace(Address, ">", "", -1)
							Address = stripSpaces(Address)[2:]
						}
					}

				}
			}
		}
	}

	QuantityStartIndex := strings.Index(creditEmail, "Qty")
	if QuantityStartIndex != -1 {
		SubString := creditEmail[QuantityStartIndex+len("Qty")+4:]
		QuantityStartIndex := strings.Index(SubString, "Qty")
		if QuantityStartIndex != -1 {
			QuantityEndIndex := strings.Index(SubString[QuantityStartIndex:], "</span>")
			if QuantityEndIndex != -1 {
				Quantity = strings.Replace(SubString[QuantityStartIndex+len("<br/>"):QuantityStartIndex+QuantityEndIndex], "<span style=", "", -1)
				Quantity = strings.Replace(Quantity, "display:inline-block;float:right;", "", -1)
				Quantity = strings.Replace(Quantity, "b>", "", -1)
				Quantity = strings.Replace(Quantity, ">", "", -1)
				Quantity = (stripSpaces(Quantity))[2:]
			}
		}
	}

	// Unit Price
	ItemNameStartIndex := strings.Index(creditEmail, "Unit Price")
	if ItemNameStartIndex != -1 {
		ItemNameEndIndex := strings.Index(creditEmail[ItemNameStartIndex:], "<br />")
		if ItemNameEndIndex != -1 {
			SubString3 := creditEmail[ItemNameStartIndex+len("Unit Price") : ItemNameStartIndex+ItemNameEndIndex+10]
			ItemNameStartIndex := strings.Index(SubString3, "<a href=http://link.order.homedepot.com")
			if ItemNameStartIndex != -1 {
				ItemNameEndIndex := strings.Index(SubString3[ItemNameStartIndex:], "<br />")
				if ItemNameEndIndex != -1 {
					SubString := SubString3[ItemNameStartIndex+len("<a href=http://link.order.homedepot.com") : ItemNameStartIndex+ItemNameEndIndex+10]
					ItemNameStartIndex := strings.Index(SubString, "<a href=http://link.order.homedepot.com")
					if ItemNameStartIndex != -1 {
						ItemNameEndIndex := strings.Index(SubString[ItemNameStartIndex:], "<br />")
						if ItemNameEndIndex != -1 {
							SubString2 := SubString[ItemNameStartIndex+len("<a href=http://link.order.homedepot.com") : ItemNameStartIndex+ItemNameEndIndex+6]
							ItemNameStartIndex := strings.Index(SubString2, "_blank")
							// fmt.Println(SubString2)
							if ItemNameStartIndex != -1 {
								ItemNameEndIndex := strings.Index(SubString2[ItemNameStartIndex:], "<br />")
								if ItemNameEndIndex != -1 {
									ItemName = SubString2[ItemNameStartIndex+len("_blank")+2 : ItemNameStartIndex+ItemNameEndIndex]
									ItemName = stripSpaces(ItemName)
									ItemName = strings.Replace(ItemName, "</a></span>", "", -1)

								}
							}
						}
					}
				}
			}
		}
	}

	DiffDays := 0.0
	if len(InternalDate) > 5 && len(OrderDate) > 5 {
		OrderDate = stripSpaces(OrderDate)
		OrderDate = strings.Replace(OrderDate, ".", "-", -1)
		OrderDate = strings.Replace(OrderDate, ",", "-", -1)
		Format := OrderDate[7:11] + "-" + OrderDate[0:6]
		InternalDateFormatted, err := time.Parse("2006-01-02T15:04:05.000Z", InternalDate[0:10]+"T11:45:26.371Z")
		if err != nil {
			fmt.Println(err)
		}
		OrderDateFormatted, err := time.Parse("2006-Jan-02", Format)
		if err != nil {
			fmt.Println(err)
		}
		DiffDays = InternalDateFormatted.Sub(OrderDateFormatted).Hours() / 24
	}

	row = append(row, InternetNumber)
	row = append(row, stripSpaces(OrderNumber))
	row = append(row, EmailReceiver)
	row = append(row, stripSpaces(CreditAmount))
	row = append(row, stripSpaces(TrackingNumber))
	row = append(row, Carrier)

	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, OrderDate)
	row = append(row, DiffDays)
	row = append(row, stripSpaces(StoreSKU))
	row = append(row, Address)
	row = append(row, Quantity)
	row = append(row, ItemName)
	return row
}

func GetReturnedProductsReport(creditEmail string, InternalDate string, EmailReceiver string) []interface{} {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	OrderNumber := ""
	OrderDate := ""
	InternetNumber := ""
	CreditAmount := ""
	StoreSKU := ""
	Quantity := ""
	ItemName := ""

	OrderNumberStartIndex := strings.Index(creditEmail, "Order Number")
	if OrderNumberStartIndex != -1 {
		OrderNumberEndIndex := strings.Index(creditEmail[OrderNumberStartIndex:], "</span>")
		if OrderNumberEndIndex != -1 {
			Substring := creditEmail[OrderNumberStartIndex+len("Order Number") : OrderNumberStartIndex+OrderNumberEndIndex+9]
			OrderNumberStartIndex = strings.Index(Substring, "<span")
			if OrderNumberStartIndex != -1 {
				OrderNumberEndIndex = strings.Index(Substring[OrderNumberStartIndex:], "</span>")
				if OrderNumberEndIndex != -1 {
					OrderNumber = Substring[OrderNumberStartIndex+len("<span>") : OrderNumberStartIndex+OrderNumberEndIndex]
					OrderNumber = strings.Replace(OrderNumber, ":</b>", "", -1)
					OrderNumber = strings.Replace(OrderNumber, ">", "", -1)
				}
			}

		}
	}
	OrderDateStartIndex := strings.Index(creditEmail, "Order Date")
	if OrderDateStartIndex != -1 {
		OrderDateEndIndex := strings.Index(creditEmail[OrderDateStartIndex:], "</span>")
		if OrderDateEndIndex != -1 {
			SubString := creditEmail[OrderDateStartIndex+len("Order Date") : OrderDateStartIndex+OrderDateEndIndex+9]
			OrderDateStartIndex := strings.Index(SubString, "<span>")
			if OrderDateStartIndex != -1 {
				OrderDateEndIndex := strings.Index(SubString[OrderDateStartIndex:], "</span>")
				if OrderDateEndIndex != -1 {
					OrderDate = SubString[OrderDateStartIndex+len("<span>") : OrderDateStartIndex+OrderDateEndIndex]
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

	CreditAmountStartIndex := strings.Index(creditEmail, "$")
	if CreditAmountStartIndex != -1 {
		CreditAmountEndIndex := strings.Index(creditEmail[CreditAmountStartIndex:], "</div>")
		if CreditAmountEndIndex != -1 {
			CreditAmount = creditEmail[CreditAmountStartIndex-1 : CreditAmountStartIndex+CreditAmountEndIndex]
			// CreditAmount = strings.Replace(CreditAmount, " USD", "", -1)

		}
	}

	StoreSKUStartIndex := strings.Index(creditEmail, "Store SKU #")
	if StoreSKUStartIndex != -1 {
		StoreSKUEndIndex := strings.Index(creditEmail[StoreSKUStartIndex:], "<br")
		if StoreSKUEndIndex != -1 {
			StoreSKU = creditEmail[StoreSKUStartIndex+len("Store SKU #") : StoreSKUStartIndex+StoreSKUEndIndex]
		}
	}

	QuantityStartIndex := strings.Index(creditEmail, "<b>Qty</b>")
	if QuantityStartIndex != -1 {
		QuantityEndIndex := strings.Index(creditEmail[QuantityStartIndex:], "</span>")
		if QuantityEndIndex != -1 {
			Quantity = creditEmail[QuantityStartIndex+len("<b>Qty</b>") : QuantityStartIndex+QuantityEndIndex]
			Quantity = strings.Replace(Quantity, "<span style=", "", -1)
			Quantity = strings.Replace(Quantity, "display:inline-block;float:right", "", -1)
			Quantity = strings.Replace(Quantity, ";", "", -1)
			Quantity = (stripSpaces(Quantity))[3:]

		}
	}

	ItemNameStartIndex := strings.Index(creditEmail, "</td> <td style")
	if ItemNameStartIndex != -1 {
		ItemNameEndIndex := strings.Index(creditEmail[ItemNameStartIndex:], "</a>")
		if ItemNameEndIndex != -1 {
			SubString3 := creditEmail[ItemNameStartIndex+len("</td> <td style") : ItemNameStartIndex+ItemNameEndIndex+10]
			ItemNameStartIndex := strings.Index(SubString3, "<a href=http://link.order.homedepot.com")
			if ItemNameStartIndex != -1 {
				ItemNameEndIndex := strings.Index(SubString3[ItemNameStartIndex:], "</a>")
				if ItemNameEndIndex != -1 {
					SubString := SubString3[ItemNameStartIndex+len("<a href=http://link.order.homedepot.com") : ItemNameStartIndex+ItemNameEndIndex+10]
					ItemNameStartIndex := strings.Index(SubString, "<a href=http://link.order.homedepot.com")
					if ItemNameStartIndex != -1 {
						ItemNameEndIndex := strings.Index(SubString[ItemNameStartIndex:], "</a>")
						if ItemNameEndIndex != -1 {
							SubString2 := SubString[ItemNameStartIndex+len("<a href=http://link.order.homedepot.com") : ItemNameStartIndex+ItemNameEndIndex+6]
							ItemNameStartIndex := strings.Index(SubString2, "_blank")
							// fmt.Println(SubString2)
							if ItemNameStartIndex != -1 {
								ItemNameEndIndex := strings.Index(SubString2[ItemNameStartIndex:], "</a>")
								if ItemNameEndIndex != -1 {
									ItemName = SubString2[ItemNameStartIndex+len("_blank")+2 : ItemNameStartIndex+ItemNameEndIndex]
									ItemName = stripSpaces(ItemName)
									ItemName = strings.Replace(ItemName, "</a></span>", "", -1)

								}
							}
						}
					}
				}
			}
		}
	}

	row = append(row, InternetNumber)
	row = append(row, stripSpaces(OrderNumber))
	row = append(row, EmailReceiver)
	row = append(row, stripSpaces(CreditAmount))

	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, OrderDate)
	row = append(row, stripSpaces(StoreSKU))
	row = append(row, Quantity)
	row = append(row, ItemName)
	return row
}

func GetNewOrderReport(creditEmail string, InternalDate string, EmailReceiver string) []interface{} {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	OrderNumber := ""
	InternetNumber := ""
	CreditAmount := ""
	StoreSKU := ""
	Address := ""
	EstimatedArrival := ""
	Quantity := ""
	ItemTotal := ""
	ItemName := ""

	OrderNumberStartIndex := strings.Index(creditEmail, "Order Number </span>")
	if OrderNumberStartIndex != -1 {
		OrderNumberEndIndex := strings.Index(creditEmail[OrderNumberStartIndex+len("Order Number </span>"):], "</span>")
		if OrderNumberEndIndex != -1 {
			Substring := creditEmail[OrderNumberStartIndex+len("Order Number </span>") : OrderNumberStartIndex+OrderNumberEndIndex+len("Order Number </span> </span>")]
			OrderNumberStartIndex = strings.Index(Substring, "font-weight:bold;")
			if OrderNumberStartIndex != -1 {
				OrderNumberEndIndex = strings.Index(Substring[OrderNumberStartIndex:], "</span>")
				if OrderNumberEndIndex != -1 {
					OrderNumber = Substring[OrderNumberStartIndex+len("font-weight:bold;")+2 : OrderNumberStartIndex+OrderNumberEndIndex]
				}
			}

		}
	}
	// OrderDateStartIndex := strings.Index(creditEmail, "Order Date")
	// if OrderDateStartIndex != -1 {
	// 	OrderDateEndIndex := strings.Index(creditEmail[OrderDateStartIndex+20:], "</span>")
	// 	if OrderDateEndIndex != -1 {
	// 		SubString := creditEmail[OrderDateStartIndex+len("Order Date") : OrderDateStartIndex+OrderDateEndIndex]
	// 		OrderDateStartIndex := strings.Index(SubString, "font-weight:bold;")
	// 		if OrderDateStartIndex != -1 {
	// 			OrderDateEndIndex := strings.Index(SubString[OrderDateStartIndex:], "</span>")
	// 			if OrderDateEndIndex != -1 {
	// 				OrderDate = SubString[OrderDateStartIndex+len("font-weight:bold;")+2 : OrderDateStartIndex+OrderDateEndIndex]
	// 			}
	// 		}

	// 	}
	// }

	InternetNumberStartIndex := strings.Index(creditEmail, "Internet #")
	if InternetNumberStartIndex != -1 {
		InternetNumberEndIndex := strings.Index(creditEmail[InternetNumberStartIndex:], "<br")
		if InternetNumberEndIndex != -1 {
			InternetNumber = creditEmail[InternetNumberStartIndex+len("Internet #") : InternetNumberStartIndex+InternetNumberEndIndex]
		}
	}

	CreditAmountStartIndex := strings.Index(creditEmail, "Order Total")
	if CreditAmountStartIndex != -1 {
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
	i := 0
	// var ProductRow []interface{}
	ParentStoreSKUStartIndex := 0
	ParentAddressStartIndex := 0
	ParentQuantityStartIndex := 0
	ParentItemTotalStartIndex := 0
	ParentItemNameStartIndex := 0
	ParentEstArrivalStartIndex := 0

	for i < 1 {

		StoreSKUCreditEmail := creditEmail[ParentStoreSKUStartIndex:]
		StoreSKUStartIndex := strings.Index(StoreSKUCreditEmail, "Store SKU #")
		if StoreSKUStartIndex != -1 {
			ParentStoreSKUStartIndex = StoreSKUStartIndex + 20
			StoreSKUEndIndex := strings.Index(StoreSKUCreditEmail[StoreSKUStartIndex:], "<br")
			if StoreSKUEndIndex != -1 {
				StoreSKU = StoreSKUCreditEmail[StoreSKUStartIndex+len("Store SKU #") : StoreSKUStartIndex+StoreSKUEndIndex]
			}
		}

		ParentAddressCreditEmail := creditEmail[ParentAddressStartIndex:]
		AddressStartIndex := strings.Index(ParentAddressCreditEmail, "inside-store-detail")
		if AddressStartIndex != -1 {
			ParentAddressStartIndex = AddressStartIndex + 20
			AddressEndIndex := strings.Index(ParentAddressCreditEmail[AddressStartIndex:], "</div>")
			if AddressEndIndex != -1 {
				SubString := ParentAddressCreditEmail[AddressStartIndex+len("inside-store-detail") : AddressStartIndex+AddressEndIndex+len("</div>")]
				AddressStartIndex = strings.Index(SubString, "<br />")
				if AddressStartIndex != -1 {
					AddressEndIndex = strings.Index(SubString[AddressStartIndex:], "</div>")
					if AddressEndIndex != -1 {
						SubString2 := SubString[AddressStartIndex+len("<br />") : AddressStartIndex+AddressEndIndex+len("</div>")]
						AddressStartIndex = strings.Index(SubString2, "<br />")
						if AddressStartIndex != -1 {
							AddressEndIndex = strings.Index(SubString2[AddressStartIndex+6:], "<br />")
							if AddressEndIndex != -1 {
								Address = SubString2[AddressStartIndex+len("<br />") : AddressStartIndex+AddressEndIndex]
								Address = strings.Replace(Address, "<span>", "", -1)
								Address = strings.Replace(Address, "span style=", "", -1)
								Address = strings.Replace(Address, "text-transform:capitalize;", "", -1)
								Address = strings.Replace(Address, "</span>", "", -1)
								Address = strings.Replace(Address, "</div>", "", -1)
								Address = strings.Replace(Address, "<", "", -1)
								Address = strings.Replace(Address, ">", "", -1)
								Address = stripSpaces(Address)
								if len(Address) > 2 {
									Address = Address[2:]
								}
							}
						}
					}
				}
			}
		}

		EstArrivalCreditEmail := creditEmail[ParentEstArrivalStartIndex:]
		EstArrivalStartIndex := strings.Index(EstArrivalCreditEmail, "Est Arrival:")
		if EstArrivalStartIndex != -1 {
			ParentEstArrivalStartIndex = EstArrivalStartIndex + 20
			EstArrivalEndIndex := strings.Index(EstArrivalCreditEmail[EstArrivalStartIndex:], "</span>")
			if EstArrivalEndIndex != -1 {
				EstimatedArrival = EstArrivalCreditEmail[EstArrivalStartIndex+len("Est Arrival:") : EstArrivalStartIndex+EstArrivalEndIndex]
				EstimatedArrival = strings.Replace(EstimatedArrival, "<b style", "", -1)
				EstimatedArrival = strings.Replace(EstimatedArrival, ">", "", -1)
				EstimatedArrival = strings.Replace(EstimatedArrival, "color:#00873C", "", -1)
				EstimatedArrival = strings.Replace(EstimatedArrival, "</b", "", -1)
				EstimatedArrival = EstimatedArrival[4:]
			}
		}

		QuantityCreditEmail := creditEmail[ParentQuantityStartIndex:]
		QuantityStartIndex := strings.Index(QuantityCreditEmail, "Qty")
		if QuantityStartIndex != -1 {
			SubString := QuantityCreditEmail[QuantityStartIndex+len("Qty")+4:]
			ParentQuantityStartIndex = QuantityStartIndex + 20
			QuantityStartIndex = strings.Index(SubString, "Qty")
			if QuantityStartIndex != -1 {
				QuantityEndIndex := strings.Index(SubString[QuantityStartIndex:], "</span>")
				if QuantityEndIndex != -1 {
					Quantity = strings.Replace(SubString[QuantityStartIndex+len("<br/>"):QuantityStartIndex+QuantityEndIndex], "<span style=", "", -1)
					Quantity = strings.Replace(Quantity, "display:inline-block;float:right;", "", -1)
					Quantity = strings.Replace(Quantity, "b>", "", -1)
					Quantity = strings.Replace(Quantity, ">", "", -1)
					Quantity = (stripSpaces(Quantity))[3:]
				}
			}
		}

		ItemTotalCreditEmail := creditEmail[ParentItemTotalStartIndex:]
		ItemTotalStartIndex := strings.Index(ItemTotalCreditEmail, "<b> Item Total </b>")
		if ItemTotalStartIndex != -1 {
			ParentItemTotalStartIndex = ItemTotalStartIndex + 20
			ItemTotalEndIndex := strings.Index(ItemTotalCreditEmail[ItemTotalStartIndex:], "</span>")
			if ItemTotalEndIndex != -1 {
				ItemTotal = ItemTotalCreditEmail[ItemTotalStartIndex+len("<b> Item Total </b>") : ItemTotalStartIndex+ItemTotalEndIndex]
				ItemTotal = strings.Replace(ItemTotal, "display:inline-block;float:right;", "", -1)
				ItemTotal = strings.Replace(ItemTotal, ">", "", -1)
				ItemTotal = strings.Replace(ItemTotal, "<span style", "", -1)
				ItemTotal = ItemTotal[4:]
			}
		}

		// Unit Price
		ItemNameCreditEmail := creditEmail[ParentItemNameStartIndex:]
		ItemNameStartIndex := strings.Index(ItemNameCreditEmail, "<table>")
		if ItemNameStartIndex != -1 {
			ParentItemNameStartIndex = ItemNameStartIndex + 20
			ItemNameEndIndex := strings.Index(ItemNameCreditEmail[ItemNameStartIndex:], "<br/>")
			if ItemNameEndIndex != -1 {
				SubString3 := ItemNameCreditEmail[ItemNameStartIndex+len("<table>") : ItemNameStartIndex+ItemNameEndIndex+10]
				ItemNameStartIndex = strings.Index(SubString3, "<a href=http://link.order.homedepot.com")
				if ItemNameStartIndex != -1 {
					ItemNameEndIndex = strings.Index(SubString3[ItemNameStartIndex:], "<br/>")
					if ItemNameEndIndex != -1 {
						SubString := SubString3[ItemNameStartIndex+len("<a href=http://link.order.homedepot.com") : ItemNameStartIndex+ItemNameEndIndex+10]
						ItemNameStartIndex = strings.Index(SubString, "<a href=http://link.order.homedepot.com")
						if ItemNameStartIndex != -1 {
							ItemNameEndIndex = strings.Index(SubString[ItemNameStartIndex:], "<br/>")
							if ItemNameEndIndex != -1 {
								SubString2 := SubString[ItemNameStartIndex+len("<a href=http://link.order.homedepot.com") : ItemNameStartIndex+ItemNameEndIndex+6]
								ItemNameStartIndex = strings.Index(SubString2, "_blank")
								if ItemNameStartIndex != -1 {
									ItemNameEndIndex = strings.Index(SubString2[ItemNameStartIndex:], "<br/>")
									if ItemNameEndIndex != -1 {
										ItemName = SubString2[ItemNameStartIndex+len("_blank")+2 : ItemNameStartIndex+ItemNameEndIndex]
										ItemName = strings.Replace(ItemName, "</a></span>", "", -1)
										ItemName = strings.Replace(ItemName, "color:#333;text-decoration:none;", "", -1)
										ItemName = strings.Replace(ItemName, "style=", "", -1)
										ItemName = strings.Replace(ItemName, "</span>", "", -1)
										ItemName = strings.Replace(ItemName, "</a>", "", -1)
										ItemName = strings.Replace(ItemName, "<b>", "", -1)
										ItemName = strings.Replace(ItemName, ">", "", -1)
										ItemName = strings.Replace(ItemName, "</b", "", -1)
										ItemName = strings.Replace(ItemName, "\"", "", -1)
										ItemName = strings.Trim(ItemName, " ")
										ItemName = strings.Trim(ItemName, "\t")
										ItemName = strings.Trim(ItemName, "\n")
										ItemName = strings.Trim(ItemName, "                                                                                                                               ")
										fmt.Println(ItemName)
									}
								}
							}
						}
					}
				}
			}
		}
		i++
	}

	row = append(row, InternetNumber)
	row = append(row, stripSpaces(OrderNumber))
	row = append(row, EmailReceiver)
	row = append(row, stripSpaces(CreditAmount))
	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, stripSpaces(StoreSKU))
	row = append(row, Address)
	row = append(row, EstimatedArrival)
	row = append(row, Quantity)
	row = append(row, ItemTotal)
	row = append(row, ItemName)
	return row
}

func GetCashBackReport(creditEmail string, InternalDate string, EmailReceiver string) []interface{} {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	OrderNumber := ""
	DateOfPurchase := ""
	Store := ""
	CashBack := ""
	// fmt.Println(creditEmail)

	StartIndex := strings.Index(creditEmail, "Your order details:")
	if StartIndex != -1 {
		EndIndex := strings.Index(creditEmail[StartIndex+len("Your order details:"):], "<p style")
		if EndIndex != -1 {
			Substring := creditEmail[StartIndex+len("Your order details:") : StartIndex+EndIndex]
			OrderNumberStartIndex := strings.Index(Substring, "Order ID:")
			if OrderNumberStartIndex != -1 {
				OrderNumberEndIndex := strings.Index(Substring[OrderNumberStartIndex:], "<br>")
				if OrderNumberEndIndex != -1 {
					OrderNumber = Substring[OrderNumberStartIndex+len("Order ID:") : OrderNumberStartIndex+OrderNumberEndIndex]
					OrderNumber = strings.Replace(OrderNumber, "</b> ", "", -1)
				}
			}
			DateOfPurchaseStartIndex := strings.Index(Substring, "Date of Purchase:")
			if DateOfPurchaseStartIndex != -1 {
				DateOfPurchaseEndIndex := strings.Index(Substring[DateOfPurchaseStartIndex:], "<br>")
				if DateOfPurchaseEndIndex != -1 {
					DateOfPurchase = Substring[DateOfPurchaseStartIndex+len("Date of Purchase:") : DateOfPurchaseStartIndex+DateOfPurchaseEndIndex]
					DateOfPurchase = strings.Replace(DateOfPurchase, "</b> ", "", -1)
				}
			}
			StoreStartIndex := strings.Index(Substring, "Store:")
			if StoreStartIndex != -1 {
				StoreEndIndex := strings.Index(Substring[StoreStartIndex:], "<br>")
				if StoreEndIndex != -1 {
					Store = Substring[StoreStartIndex+len("Store:") : StoreStartIndex+StoreEndIndex]
					Store = strings.Replace(Store, "</b> ", "", -1)
				}
			}
			CashBackStartIndex := strings.Index(Substring, "Cash Back Pending:")
			if CashBackStartIndex != -1 {
				CashBackEndIndex := strings.Index(Substring[CashBackStartIndex:], "</p>")
				if CashBackEndIndex != -1 {
					CashBack = Substring[CashBackStartIndex+len("Cash Back Pending:") : CashBackStartIndex+CashBackEndIndex]
					CashBack = strings.Replace(CashBack, "</b> ", "", -1)
				}
			}

		}
	}

	row = append(row, stripSpaces(OrderNumber))
	row = append(row, EmailReceiver)
	if len(InternalDate) > 10 {
		row = append(row, InternalDate[0:10])
	} else {
		row = append(row, InternalDate)
	}
	row = append(row, DateOfPurchase)
	row = append(row, Store)
	row = append(row, CashBack)
	return row
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

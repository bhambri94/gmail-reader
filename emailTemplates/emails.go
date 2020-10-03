package emailTemplates

import (
	"strings"
	"time"
)

var storeCreditFinalValues [][]interface{}

func StoreCreditFlushFinalValues() {
	storeCreditFinalValues = nil
}

func GetStoreCreditReport(creditEmail string, InternalDate string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	row = append(row, InternalDate)
	StoreCreditAmount := ""
	OrderNumber := ""
	AmountStartIndex := strings.Index(creditEmail, "The Home Depot sent you a")
	if AmountStartIndex != -1 {
		AmountEndIndex := strings.Index(creditEmail[AmountStartIndex:], "eStore Credit")
		if AmountEndIndex != -1 {
			StoreCreditAmount = creditEmail[AmountStartIndex+len("The Home Depot sent you a") : AmountStartIndex+AmountEndIndex]
		}
	}
	OrderNumberStartIndex := strings.Index(creditEmail, "Please refer to order number ")
	if OrderNumberStartIndex != -1 {
		OrderNumberEndIndex := strings.Index(creditEmail[OrderNumberStartIndex:], "Please keep this ")
		if OrderNumberEndIndex != -1 {
			OrderNumber = creditEmail[OrderNumberStartIndex+len("Please refer to order number")+2 : OrderNumberStartIndex+OrderNumberEndIndex-2]
		}
	}
	row = append(row, StoreCreditAmount)
	row = append(row, OrderNumber)

	storeCreditFinalValues = append(storeCreditFinalValues, row)
}

func GetStoreCreditFinalValues() [][]interface{} {
	return storeCreditFinalValues
}

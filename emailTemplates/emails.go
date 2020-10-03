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

func GetStoreCreditReport(creditEmail string, InternalDate string, EmailReceiver string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	row = append(row, InternalDate)
	StoreCreditAmount := ""
	StoreCreditLink := ""
	AmountStartIndex := strings.Index(creditEmail, "sent you a $")
	if AmountStartIndex != -1 {
		AmountEndIndex := strings.Index(creditEmail[AmountStartIndex:], "eStore Credit")
		if AmountEndIndex != -1 {
			StoreCreditAmount = creditEmail[AmountStartIndex+len("sent you a ") : AmountStartIndex+AmountEndIndex]
		}
	}
	StoreCreditLinkStartIndex := strings.Index(creditEmail, "Go to  ")
	if StoreCreditLinkStartIndex != -1 {
		StoreCreditLinkEndIndex := strings.Index(creditEmail[StoreCreditLinkStartIndex:], "to receive instructions")
		if StoreCreditLinkEndIndex != -1 {
			StoreCreditLink = creditEmail[StoreCreditLinkStartIndex+len("Go to  ") : StoreCreditLinkStartIndex+StoreCreditLinkEndIndex-2]
		}
	}
	row = append(row, StoreCreditAmount)
	row = append(row, StoreCreditLink)
	row = append(row, EmailReceiver)

	storeCreditFinalValues = append(storeCreditFinalValues, row)
}

func GetStoreCreditFinalValues() [][]interface{} {
	return storeCreditFinalValues
}

func GetCreditAppliedFinalValues() [][]interface{} {
	return creditAppliedFinalValues
}

func GetCreditAppliedReport(creditEmail string, InternalDate string, EmailReceiver string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	row = append(row, InternalDate)
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
		fmt.Println(CreditAmountStartIndex)
		CreditAmountEndIndex := strings.Index(creditEmail[CreditAmountStartIndex-15:CreditAmountStartIndex], "A")
		if CreditAmountEndIndex != -1 {
			fmt.Println(CreditAmountEndIndex)
			CreditAmount = creditEmail[CreditAmountStartIndex-15 : CreditAmountStartIndex+CreditAmountEndIndex+10]
			CreditAmount = strings.Replace(CreditAmount, "credit has been", "", -1)
			CreditAmount = strings.Replace(CreditAmount, "A", "", -1)
		}
	}

	StoreSKUStartIndex := strings.Index(creditEmail, "Store SKU #")
	if StoreSKUStartIndex != -1 {
		StoreSKUEndIndex := strings.Index(creditEmail[StoreSKUStartIndex:], "<br")
		if StoreSKUEndIndex != -1 {
			StoreSKU = creditEmail[StoreSKUStartIndex+len("Store SKU #") : StoreSKUStartIndex+StoreSKUEndIndex]
		}
	}

	row = append(row, OrderNumber)
	row = append(row, OrderDate)
	row = append(row, InternetNumber)
	row = append(row, CreditAmount)
	row = append(row, StoreSKU)
	row = append(row, EmailReceiver)
	fmt.Println(stripSpaces(OrderNumber))
	fmt.Println(stripSpaces(OrderDate))
	fmt.Println(stripSpaces(InternetNumber))
	fmt.Println(stripSpaces(CreditAmount))
	fmt.Println(stripSpaces(StoreSKU))

	creditAppliedFinalValues = append(creditAppliedFinalValues, row)
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

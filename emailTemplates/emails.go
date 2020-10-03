package emailTemplates

import (
	"strings"
	"time"
)

var storeCreditFinalValues [][]interface{}

func StoreCreditFlushFinalValues() {
	storeCreditFinalValues = nil
}

func GetStoreCreditReport(creditEmail string, InternalDate string, EmailReceiver string) {
	var row []interface{}
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	row = append(row, currentTime.Format("2006-01-02 15:04:05"))
	row = append(row, InternalDate)
	StoreCreditAmount := ""
	StoreCreditLink := ""
	AmountStartIndex := strings.Index(creditEmail, "The Home Depot sent you a")
	if AmountStartIndex != -1 {
		AmountEndIndex := strings.Index(creditEmail[AmountStartIndex:], "eStore Credit")
		if AmountEndIndex != -1 {
			StoreCreditAmount = creditEmail[AmountStartIndex+len("The Home Depot sent you a") : AmountStartIndex+AmountEndIndex]
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

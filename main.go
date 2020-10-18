package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bhambri94/gmail-reader/gmailApis"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func main() {
	sugar.Infof("starting gmail reader app server...")
	defer logger.Sync() // flushes buffer, if any

	router := fasthttprouter.New()
	router.GET("/v1/gmail-reader/query=:query/afterDate=:afterDate", handleGmailSearch)
	router.GET("/v1/gmail-reader/search=:query/fromDate=:fromDate/toDate=:toDate", handleDynamicGmailSearch)
	log.Fatal(fasthttp.ListenAndServe(":7004", router.Handler))
}

/*
http://localhost:7004/v1/gmail-reader/query='StoreCredit'/afterDate='2020-10-04'
http://localhost:7004/v1/gmail-reader/search='subject:Credit Applied to Order'/fromDate='2020-10-06'/toDate='2020-10-06'
http://localhost:7004/v1/gmail-reader/search='subject:Your order just shipped'/fromDate='2020-07-01'/toDate='2020-07-10'
http://localhost:7004/v1/gmail-reader/search='subject:We have received your returned products'/fromDate='2020-07-01'/toDate='2020-07-10'
http://localhost:7004/v1/gmail-reader/search='subject:"We received your order"'/fromDate='2020-07-01'/toDate='2020-07-02'

*/

func handleDynamicGmailSearch(ctx *fasthttp.RequestCtx) {
	sugar.Infof("calling gmail dynamic reader api!")
	SearchQuery := ctx.UserValue("query")
	if SearchQuery != nil {
		sugar.Infof("SearchQuery is := " + SearchQuery.(string))
		SearchQuery = SearchQuery.(string)[1 : len(SearchQuery.(string))-1]
	}
	EmailAfterDate := ctx.UserValue("fromDate")
	if EmailAfterDate != nil {
		sugar.Infof("EmailAfterDate is := " + EmailAfterDate.(string))
		EmailAfterDate = EmailAfterDate.(string)[1 : len(EmailAfterDate.(string))-1]
	}
	EmailBeforeDate := ctx.UserValue("toDate")
	if EmailBeforeDate != nil {
		sugar.Infof("EmailBeforeDate is := " + EmailBeforeDate.(string))
		EmailBeforeDate = EmailBeforeDate.(string)[1 : len(EmailBeforeDate.(string))-1]
	}
	loc, err := time.LoadLocation("America/Bogota")
	BeforeTillTime, err := time.ParseInLocation("2006-01-02 15:04:05", EmailBeforeDate.(string)+" 00:00:00", loc)
	if err != nil {
		fmt.Println(err)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(200)
		ctx.SetBody([]byte("Failed! Incorrect toDate Format in URl, Please fix. Example: "))
		sugar.Infof("calling gmail reader api failure!")
		return
	}
	BeforeTimeUnix := int(BeforeTillTime.Unix())

	AfterTillTime, err := time.ParseInLocation("2006-01-02 15:04:05", EmailAfterDate.(string)+" 00:00:00", loc)
	if err != nil {
		fmt.Println(err)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetStatusCode(200)
		ctx.SetBody([]byte("Failed! Incorrect fromDate Format in URl, Please fix."))
		sugar.Infof("calling gmail reader api failure!")
		return
	}
	AfterTimeUnix := int(AfterTillTime.Unix())

	var finalValues [][]interface{}
	FinalSearchQuery := SearchQuery.(string) + " before:" + strconv.Itoa(BeforeTimeUnix) + " after:" + strconv.Itoa(AfterTimeUnix)
	fmt.Println(FinalSearchQuery)
	finalValues = gmailApis.SearchForEmailDynamic(FinalSearchQuery, EmailAfterDate.(string)+" 00:00:00")
	loc, _ = time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	var header []string
	var CSVName string
	if strings.Contains(SearchQuery.(string), "Credit Applied") {
		header = []string{"EmailWorkflow_Refresh_time", "Internet Number", "Order Number", "Amount Credit", "To Email Address", "Email Received Timestamp", "Order Date", "Store SKU"}
		CSVName = "StoreCreditCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	} else if strings.Contains(SearchQuery.(string), "just shipped") {
		header = []string{"EmailWorkflow_Refresh_time", "Internet Number", "Order Number", "To Email Address", "Amount Credit", "Tracking Number", "Carrier", "Shipment Date", "Order Date", "Shipped-Order Date", "Store SKU", "Address", "Quantity", "Item Name"}
		CSVName = "ShippedOrdersCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	} else if strings.Contains(SearchQuery.(string), "returned products") {
		header = []string{"EmailWorkflow_Refresh_time", "Internet Number", "Order Number", "To Email Address", "Amount Total", "Order Date", "Email Received Date", "Store SKU", "Quantity"}
		CSVName = "ReturnedOrdersCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	} else if strings.Contains(SearchQuery.(string), "received your order") {
		header = []string{"EmailWorkflow_Refresh_time", "Internet Number", "Order Number", "To Email Address", "Amount Credit", "New Order Date", "Store SKU", "Address", "Estimated Arrival", "Quantity", "Item Total", "Item Name"}
		CSVName = "NewOrdersCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	} else if strings.Contains(SearchQuery.(string), "Cash Back") {
		header = []string{"EmailWorkflow_Refresh_time", "Order Number", "To Email Address", "Email Received Timestamp", "Date Of Purchase", "Store", "CashBack Amount"}
		CSVName = "CashBackCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	}

	f, err := os.Create(CSVName)
	if err != nil {
		log.Fatal(err)
	}
	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write(header)
	stringfinalValues := make([][]string, len(finalValues)+5)
	i := 0
	for i < len(finalValues) {
		for _, value := range finalValues[i] {
			a := fmt.Sprintf("%v", value)
			stringfinalValues[i] = append(stringfinalValues[i], a)
		}
		writer.Write(stringfinalValues[i])
		writer.Flush()
		i++
	}
	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Response.Header.Set("Content-Disposition", "attachment;filename="+CSVName)
	ctx.SendFile(CSVName)
	err = os.Remove(CSVName)
	if err != nil {
		fmt.Println("Unable to delete file")
	} else {
		fmt.Println("File Deleted")
	}
	err = os.Remove(CSVName + ".fasthttp.gz")
	if err != nil {
		fmt.Println("Unable to delete file")
	} else {
		fmt.Println("File Deleted")
	}
}

func handleGmailSearch(ctx *fasthttp.RequestCtx) {
	sugar.Infof("calling gmail reader api!")
	SearchQuery := ctx.UserValue("query")
	if SearchQuery != nil {
		sugar.Infof("SearchQuery is := " + SearchQuery.(string))
		SearchQuery = SearchQuery.(string)[1 : len(SearchQuery.(string))-1]
	}
	EmailAfterDate := ctx.UserValue("afterDate")
	if EmailAfterDate != nil {
		sugar.Infof("EmailAfterDate is := " + EmailAfterDate.(string))
		EmailAfterDate = EmailAfterDate.(string)[1 : len(EmailAfterDate.(string))-1]
	}
	var finalValues [][]interface{}
	if SearchQuery.(string) == "StoreCredit" {
		finalValues = gmailApis.SearchForEmail("subject:You've received an eStore Credit", EmailAfterDate.(string)+" 00:00:00")
	}

	// if SearchQuery.(string) == "StoreCredit" {
	// 	finalValues = gmailApis.SearchForEmail("subject:You've received an eStore Credit", EmailAfterDate.(string)+" 00:00:00")
	// }

	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	CSVName := "StoreCreditCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	f, err := os.Create(CSVName)
	if err != nil {
		log.Fatal(err)
	}
	writer := csv.NewWriter(f)
	defer writer.Flush()
	header := []string{"StoreCredit_Refresh_time", "Store Credit Link", "Email Received Timestamp", "Store Credit Amount", "To Email Address"}
	writer.Write(header)
	stringfinalValues := make([][]string, len(finalValues)+5)
	i := 0
	for i < len(finalValues) {
		for _, value := range finalValues[i] {
			a := fmt.Sprintf("%v", value)
			stringfinalValues[i] = append(stringfinalValues[i], a)
		}
		writer.Write(stringfinalValues[i])
		writer.Flush()
		i++
	}
	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Content-Type", "text/csv")
	ctx.Response.Header.Set("Content-Disposition", "attachment;filename="+CSVName)
	ctx.SendFile(CSVName)
	err = os.Remove(CSVName)
	if err != nil {
		fmt.Println("Unable to delete file")
	} else {
		fmt.Println("File Deleted")
	}
	err = os.Remove(CSVName + ".fasthttp.gz")
	if err != nil {
		fmt.Println("Unable to delete file")
	} else {
		fmt.Println("File Deleted")
	}
}

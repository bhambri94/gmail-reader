package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	sugar.Infof("starting ecommerce manager app server...")
	defer logger.Sync() // flushes buffer, if any

	router := fasthttprouter.New()
	router.GET("/v1/gmail-reader/query=:query/afterDate=:afterDate", handleGmailSearch)
	router.GET("/v1/gmail-reader/search=:query/afterDate=:afterDate", handleDynamicGmailSearch)
	log.Fatal(fasthttp.ListenAndServe(":7004", router.Handler))
}

/*
http://localhost:7004/v1/gmail-reader/query='StoreCredit'/afterDate='2020-10-04'
http://localhost:7004/v1/gmail-reader/search='subject:Credit Applied to Order'/afterDate='2020-08-20'
http://localhost:7004/v1/gmail-reader/search='subject:Your order just shipped'/afterDate='2020-07-20'
*/

func handleDynamicGmailSearch(ctx *fasthttp.RequestCtx) {
	sugar.Infof("calling gmail dynamic reader api!")
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
	finalValues = gmailApis.SearchForEmailDynamic(SearchQuery.(string), EmailAfterDate.(string)+" 00:00:00")
	loc, _ := time.LoadLocation("America/Bogota")
	currentTime := time.Now().In(loc)
	var header []string
	var CSVName string
	if strings.Contains(SearchQuery.(string), "Credit Applied") {
		header = []string{"EmailWorkflow_Refresh_time", "Internet Number", "Order Number", "Amount Credit", "To Email Address", "Email Received Timestamp", "Order Date", "Store SKU"}
		CSVName = "StoreCreditCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
	} else if strings.Contains(SearchQuery.(string), "just shipped") {
		header = []string{"EmailWorkflow_Refresh_time", "Internet Number", "Order Number", "Amount Credit", "To Email Address", "Tracking Number", "Carrier", "Email Received Timestamp", "Order Date", "Store SKU"}
		CSVName = "ShippedOrdersCSV_" + currentTime.Format("2006-01-02 15:04:05") + ".csv"
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

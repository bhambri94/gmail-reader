package gmailApis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bhambri94/gmail-reader/emailTemplates"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var gmailService *gmail.Service

// Retrieve a token, saves the token, then returns the generated client.
func getClient() *gmail.Service {

	b, err := ioutil.ReadFile("gmailApis/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	gmailClient := config.Client(context.Background(), tok)
	srv, err := gmail.New(gmailClient)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	gmailService = srv
	return srv
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetLabelsFromEmail() {

	if gmailService == nil {
		gmailService = getClient()
	}

	user := "me"
	r, err := gmailService.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}

func SearchForEmail(SearchQuery string, EmailsAfterTime string) [][]interface{} {
	loc, err := time.LoadLocation("America/Bogota")
	TillTime, err := time.ParseInLocation("2006-01-02 15:04:05", EmailsAfterTime, loc)
	if err != nil {
		fmt.Println(err)
	}
	TillTimeUnix := TillTime.Unix()
	emailTemplates.StoreCreditFlushFinalValues()

	if gmailService == nil {
		gmailService = getClient()
	}

	user := "me"
	var r *gmail.ListMessagesResponse
	NextToken := "First"
	BreakAllLoops := false
	for NextToken != "" {
		if NextToken == "First" {
			r, err = gmailService.Users.Messages.List(user).Q(SearchQuery).MaxResults(500).Do()
			if err != nil {
				fmt.Println("Unable to retrieve Message with Search Query")
				return nil
			}
			NextToken = r.NextPageToken
		} else {
			r, err = gmailService.Users.Messages.List(user).Q(SearchQuery).MaxResults(500).PageToken(NextToken).Do()
			if err != nil {
				fmt.Println("Unable to retrieve Message with Search Query")
				return nil
			}
			NextToken = r.NextPageToken
		}
		if len(r.Messages) == 0 {
			fmt.Println("No Message found.")
			break
		}
		fmt.Println("Iterating over Messages:")
		for _, l := range r.Messages {
			msg, err := gmailService.Users.Messages.Get(user, l.Id).Format("full").Do()
			if err != nil {
				fmt.Printf("Unable to retrieve Message with Search Query: %v", err)
			}
			EmailTimings := msg.InternalDate / 1000
			t := time.Unix(EmailTimings, 0)
			if err != nil {
				fmt.Println(err)
			}
			t = t.In(loc)
			CentralTime := t.Format("2006-01-02 15:04:05")
			if t.Unix() < TillTimeUnix {
				BreakAllLoops = true
				break
			}
			Header := msg.Payload.Headers
			EmailReceiver := ""
			for _, s := range Header {
				if s.Name == "To" {
					EmailReceiver = s.Value
				}
			}
			if EmailReceiver == "tuanshivam@gmail.com" {
				continue
			}
			Output := msg.Payload.Parts[0].Body.Data
			Output = strings.Replace(Output, "-", "+", -1)
			Output = strings.Replace(Output, "_", "/", -1)
			emailBody := DecodeB64(Output)
			emailTemplates.GetStoreCreditReport(emailBody, CentralTime, EmailReceiver)
		}
		if BreakAllLoops {
			break
		}
	}
	return emailTemplates.GetStoreCreditFinalValues()
}

func SearchForEmailDynamic(SearchQuery string, EmailsAfterTime string) [][]interface{} {
	var shippingTrackerFinalValues [][]interface{}
	loc, err := time.LoadLocation("America/Bogota")
	TillTime, err := time.ParseInLocation("2006-01-02 15:04:05", EmailsAfterTime, loc)
	if err != nil {
		fmt.Println(err)
	}
	TillTimeUnix := TillTime.Unix()
	emailTemplates.CreditAppliedFlushFinalValues()

	if gmailService == nil {
		gmailService = getClient()
	}

	user := "me"
	var r *gmail.ListMessagesResponse
	NextToken := "First"
	BreakAllLoops := false
	Output := ""
	for NextToken != "" {
		if NextToken == "First" {
			r, err = gmailService.Users.Messages.List(user).Q(SearchQuery).MaxResults(20).Do()
			if err != nil {
				fmt.Println("Unable to retrieve Message with Search Query")
				return nil
			}
			NextToken = r.NextPageToken
		} else {
			r, err = gmailService.Users.Messages.List(user).Q(SearchQuery).MaxResults(20).PageToken(NextToken).Do()
			if err != nil {
				fmt.Println("Unable to retrieve Message with Search Query")
				return nil
			}
			NextToken = r.NextPageToken
		}
		if len(r.Messages) == 0 {
			fmt.Println("No Message found.")
			break
		}
		fmt.Println("Iterating over Messages:")

		for _, l := range r.Messages {
			msg, err := gmailService.Users.Messages.Get(user, l.Id).Format("full").Do()
			if err != nil {
				log.Fatalf("Unable to retrieve Message with Search Query: %v", err)
			}
			EmailTimings := msg.InternalDate / 1000
			t := time.Unix(EmailTimings, 0)
			if err != nil {
				fmt.Println(err)
			}
			t = t.In(loc)
			CentralTime := t.Format("2006-01-02 15:04:05")
			if t.Unix() < TillTimeUnix {
				BreakAllLoops = true
				break
			}
			Header := msg.Payload.Headers
			EmailReceiver := ""
			for _, s := range Header {
				if s.Name == "To" {
					EmailReceiver = s.Value
				}
			}
			if EmailReceiver == "tuanshivam@gmail.com" {
				continue
			}
			if len(msg.Payload.Parts) > 0 {
				Output = msg.Payload.Parts[0].Body.Data
				Output = strings.Replace(Output, "-", "+", -1)
				Output = strings.Replace(Output, "_", "/", -1)
				emailBody := DecodeB64(Output)
				emailTemplates.GetCreditAppliedReport(emailBody, CentralTime, EmailReceiver)
			} else {
				Output = msg.Payload.Body.Data
				Output = strings.Replace(Output, "-", "+", -1)
				Output = strings.Replace(Output, "_", "/", -1)
				emailBody := DecodeB64(Output)
				if strings.Contains(SearchQuery, "Credit Applied") {
					emailTemplates.GetCreditAppliedReport(emailBody, CentralTime, EmailReceiver)
				} else if strings.Contains(SearchQuery, "just shipped") {
					shippingTrackerFinalValues = append(shippingTrackerFinalValues, emailTemplates.GetShippingTrackerReport(emailBody, CentralTime, EmailReceiver))
				}

			}
		}
		if BreakAllLoops {
			break
		}
	}

	if strings.Contains(SearchQuery, "Credit Applied") {
		return emailTemplates.GetCreditAppliedFinalValues()
	} else if strings.Contains(SearchQuery, "just shipped") {
		return shippingTrackerFinalValues
	}
	return nil
}

func DecodeB64(message string) string {
	emailBody, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		fmt.Println(err)
	}
	return string(emailBody)
}

//yjt-nuefu-wb

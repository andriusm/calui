package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

const (
	configDir = ".config/calui"
)

var (
	homeDir string
)

func init() {
	homeDir = os.Getenv("HOME")
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := filepath.Join(homeDir, configDir, "token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n\nEnter code:", authURL)

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

func getEvents() [][]string {
	lines := make([][]string, 5)

	credFile := filepath.Join(homeDir, configDir, "credentials.json")
	b, err := ioutil.ReadFile(credFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	n := time.Now()

	n = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	beg := n.AddDate(0, 0, -int(n.Weekday()-1))
	end := beg.AddDate(0, 0, 5)

	//t := time.Now().Format(time.RFC3339)
	listCall := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true)
	//events, err := listCall.TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	events, err := listCall.
		TimeMin(beg.Format(time.RFC3339)).
		TimeMax(end.Format(time.RFC3339)).
		MaxResults(25).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			//line := fmt.Sprintf("%v (%v)\n", item.Summary, date)

			start, _ := time.Parse(time.RFC3339, date)
			d := int(start.Weekday()) - 1

			// item.HangoutLink
			// item.Description
			line := fmt.Sprintf("@%v - %v", start.Format("15:04"), item.Summary)
			lines[d] = append(lines[d], line)
		}
	}

	return lines
}

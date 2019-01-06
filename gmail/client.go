package gmail

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

type Client struct {
	service *gmail.Service
	regexp  *regexp.Regexp
}

const User = "me"
const EmailFinder = "<(.*)>"

// Retrieve a token, saves the token, then returns the generated client.
func GetClient() *Client {
	b, err := ioutil.ReadFile("secrets/oauth_credentials.json")
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
	tokFile := "secrets/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		panic(err)
	}
	client := config.Client(context.Background(), tok)
	srv, err := gmail.New(client)
	if err != nil {
		panic(err)
	}

	return &Client{service: srv, regexp: regexp.MustCompile(EmailFinder)}
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

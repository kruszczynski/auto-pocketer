package gmail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

type Client struct {
	service       *gmail.Service
	regexp        *regexp.Regexp
	lastHistoryId uint64
	watch         func(string, *gmail.WatchRequest) *gmail.UsersWatchCall
}

const User = "me"
const EmailFinder = "<(.*)>"

// Retrieve a token, saves the token, then returns the generated client.
func GetClient() (*Client, error) {
	oauthSecretPath := os.Getenv("OAUTH_SECRET_PATH")
	b, err := ioutil.ReadFile(oauthSecretPath)
	if err != nil {
		log.Printf("Unable to read client secret file: %v\n", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.GmailModifyScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v\n", err)
		return nil, err
	}

	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := os.Getenv("GMAIL_TOKEN_PATH")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background(), tok)
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	return &Client{
		service: srv,
		regexp:  regexp.MustCompile(EmailFinder),
		watch:   srv.Users.Watch,
	}, nil
}

func (client *Client) ListMessageIds(historyId uint64) ([]string, error) {
	ret := []string{}
	if historyId < client.lastHistoryId {
		return ret, nil
	}
	c := client.service.Users.History.List(User)
	c.StartHistoryId(client.lastHistoryId)
	r, err := c.Do()
	if err != nil {
		log.Printf("Unable to retrieve labels: %v\n", err)
		return ret, err
	}
	lastHistoryId := historyId
	for _, h := range r.History {
		if h.Id > lastHistoryId {
			lastHistoryId = h.Id
		}
		if len(h.MessagesAdded) > 0 {
			for _, m := range h.MessagesAdded {
				ret = append(ret, m.Message.Id)
			}
		}
	}
	client.lastHistoryId = lastHistoryId
	return ret, nil
}

func (client *Client) ProcessMessages(messageIds []string) ([]*ProcessedMessage, error) {
	ret := []*ProcessedMessage{}
	for _, msgId := range messageIds {
		fetchedMessage, err := client.fetchMessage(msgId)
		if err != nil {
			return ret, err
		}
		ret = append(ret, fetchedMessage)
	}
	return ret, nil
}

func (client *Client) Archive(messageId string) error {
	fmt.Printf("Archiving message %s\n", messageId)
	request := &gmail.ModifyMessageRequest{RemoveLabelIds: []string{"INBOX"}}
	_, err := client.service.Users.Messages.Modify(User, messageId, request).Do()
	return err
}

func (client *Client) fetchMessage(messageId string) (*ProcessedMessage, error) {
	fmt.Printf("Fetching message %s\n", messageId)
	r, err := client.service.Users.Messages.Get(User, messageId).Do()
	if err != nil {
		log.Println("Unable to fetch messages from Gmail")
		return nil, err
	}

	from := ""
	to := ""
	body := ""
	for _, header := range r.Payload.Headers {
		if header.Name == "From" {
			from = client.regexp.FindString(header.Value)
		}
		if header.Name == "To" {
			to = client.regexp.FindString(header.Value)
		}
	}
	for _, mp := range r.Payload.Parts {
		if mp.MimeType == "text/plain" {
			dec, err := base64.URLEncoding.DecodeString(mp.Body.Data)
			if err != nil {
				log.Printf("Body decoding failed %v\n", err)
				return nil, err
			}

			// remove threads
			threadRegexp := regexp.MustCompile("^>+ ")
			filteredBody := []string{}
			for _, line := range strings.Split(string(dec), "\n") {
				if !threadRegexp.MatchString(line) {
					filteredBody = append(filteredBody, line)
				}
			}
			body = strings.Join(filteredBody, "\n")
		}
	}
	return &ProcessedMessage{Id: messageId, From: from, To: to, Body: body}, nil
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

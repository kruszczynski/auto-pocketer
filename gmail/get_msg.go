package gmail

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type ProcessedMessage struct {
	From string
	To   string
	Body string
}

func (client Client) ListMessageIds(startHistoryId uint64, historyId uint64) (ret []string) {
	c := client.service.Users.History.List(User)
	c.StartHistoryId(startHistoryId)
	r, err := c.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	for _, h := range r.History {
		if len(h.MessagesAdded) > 0 {
			for _, m := range h.MessagesAdded {
				ret = append(ret, m.Message.Id)
			}
		}
	}
	return
}

func (client Client) ProcessMessages(messageIds []string) (ret []*ProcessedMessage) {
	for _, msgId := range messageIds {
		ret = append(ret, client.fetchMessage(msgId))
	}
	return
}

func (client Client) fetchMessage(messageId string) *ProcessedMessage {
	fmt.Printf("Fetching message %s\n", messageId)
	r, err := client.service.Users.Messages.Get(User, messageId).Do()
	if err != nil {
		panic(err)
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
				log.Fatalf("Body decoding failed %v", err)
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
	return &ProcessedMessage{From: from, To: to, Body: body}
}

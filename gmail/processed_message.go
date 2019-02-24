package gmail

import (
	"fmt"
	"os"
	"strings"

	xurls "github.com/mvdan/xurls"
)

var allowedSenders []string

func init() {
	allowedSenders = strings.Split(os.Getenv("ALLOWED_SENDERS"), ",")
}

type ProcessedMessage struct {
	Id   string
	From string
	To   string
	Body string
}

func (pMsg ProcessedMessage) FindLink() string {
	link := xurls.Strict.FindString(pMsg.Body)
	if link != "" {
		fmt.Printf("Found a link: %s\n", link)
	} else {
		fmt.Printf("No link found in the message\n")
	}
	return link
}

func (pMsg ProcessedMessage) AllowedSender() bool {
	for _, sender := range allowedSenders {
		if sender == pMsg.From {
			return true
		}
	}
	return false
}

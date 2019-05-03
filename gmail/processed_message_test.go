package gmail

import (
	"os"
	"testing"
)

func TestFindLink_FindsLink(t *testing.T) {
	msg := &ProcessedMessage{
		Body: "hello https://google.com",
	}

	link := msg.FindLink()
	if link != "https://google.com" {
		t.Error("Should've found https://google.com")
	}
}

func TestFindLink_FindsFirst(t *testing.T) {
	msg := &ProcessedMessage{
		Body: `hello my dear friend. There are two urls to be found in this message
					 hello https://party.com
					 hello https://google.com`,
	}

	link := msg.FindLink()
	if link != "https://party.com" {
		t.Error("Should've found https://party.com")
	}
}

func TestFindLink_NothingFound(t *testing.T) {
	msg := &ProcessedMessage{
		Body: "hello my dear friend. There is no url to be found in this message",
	}

	link := msg.FindLink()
	if link != "" {
		t.Error("Empty string was not returned")
	}
}

func TestAllowedSender(t *testing.T) {
	os.Setenv("ALLOWED_SENDERS", "hello@pentagon.com,hello@newsletter.com")

	msg := &ProcessedMessage{From: "hello@newsletter.com"}
	if !msg.AllowedSender() {
		t.Error("hello@newsletter.com should have been allowed")
	}

	msg = &ProcessedMessage{From: "john@cleese.com"}
	if msg.AllowedSender() {
		t.Error("john@cleese.com should not have been allowed")
	}
}

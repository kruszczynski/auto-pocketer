package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/kruszczynski/auto-pocketer/gmail"
	"github.com/kruszczynski/auto-pocketer/googlepubsub"
	"mvdan.cc/xurls"
)

func main() {
	gmailClient := gmail.GetClient()
	startHistoryId := gmailClient.Watch()

	sub := googlepubsub.GetSubscription()
	var mu sync.Mutex
	errr := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		// locks because of startHistoryId
		// it's not an issue though
		mu.Lock()
		defer mu.Unlock()

		msg.Ack()

		var message googlepubsub.Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got message: %q\n", message)
		messageIds := gmailClient.ListMessageIds(startHistoryId, message.HistoryId)
		fmt.Printf("%d new messages received\n", len(messageIds))
		startHistoryId = message.HistoryId

		processedMessages := gmailClient.ProcessMessages(messageIds)
		for _, pm := range processedMessages {
			findLink(pm)
		}
	})

	if errr != nil {
		panic(errr)
	}
}

func findLink(msg *gmail.ProcessedMessage) string {
	link := xurls.Strict().FindString(msg.Body)
	fmt.Println(link)
	return link
}

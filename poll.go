package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/kruszczynski/auto-pocketer/gmail"
	"github.com/kruszczynski/auto-pocketer/googlepubsub"
	"github.com/kruszczynski/auto-pocketer/pocket"
)

func main() {
	fmt.Println("This is updated version")
	gmailClient := gmail.GetClient()
	pocketClient := pocket.GetClient()
	gmailClient.Watch()

	sub := googlepubsub.GetSubscription()
	var mu sync.Mutex
	errr := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		// locks because of startHistoryId is shared
		mu.Lock()
		defer mu.Unlock()

		msg.Ack()

		var message googlepubsub.Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil {
			panic(err)
		}
		messageIds := gmailClient.ListMessageIds(message.HistoryId)
		fmt.Printf("%d new messages received\n", len(messageIds))

		processedMessages := gmailClient.ProcessMessages(messageIds)
		filteredMessages := filterMessages(processedMessages)
		for _, pm := range filteredMessages {
			link := pm.FindLink()
			if link != "" {
				pocketClient.Add(link)
				gmailClient.Archive(pm.Id)
			}
		}
	})

	if errr != nil {
		panic(errr)
	}
}

func filterMessages(messages []*gmail.ProcessedMessage) (ret []*gmail.ProcessedMessage) {
	for _, msg := range messages {
		fmt.Printf("Message sender is: %s\n", msg.From)
		if msg.AllowedSender() {
			ret = append(ret, msg)
		}
	}
	return
}

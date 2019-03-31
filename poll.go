package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/getsentry/raven-go"
	"github.com/kruszczynski/auto-pocketer/gmail"
	"github.com/kruszczynski/auto-pocketer/googlepubsub"
	"github.com/kruszczynski/auto-pocketer/pocket"
)

func main() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	gmailClient := gmail.GetClient()
	pocketClient := pocket.GetClient()
	gmailClient.Watch()
	gmailWatchCalledOn := daysSinceBeginning()

	sub := googlepubsub.GetSubscription()
	var mu sync.Mutex
	errr := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		// locks because of startHistoryId is shared
		mu.Lock()
		defer mu.Unlock()

		// call watch if more than one day has elapsed
		if gmailWatchCalledOn < daysSinceBeginning() {
			gmailClient.Watch()
			gmailWatchCalledOn = daysSinceBeginning()
		}

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
				_ = pocketClient.Add(link)
				// TODO(kruszczynski) report that to Raven
				gmailClient.Archive(pm.Id)
			}
		}
	})

	if errr != nil {
		// TODO(kruszczynski) no more panic
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

func daysSinceBeginning() int64 {
	return time.Now().Unix() / 60 / 60 / 24
}

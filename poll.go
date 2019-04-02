package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	gmailClient, err := gmail.GetClient()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}
	pocketClient := pocket.GetClient()
	gmailClient.Watch()
	gmailWatchCalledOn := daysSinceBeginning()

	sub := googlepubsub.GetSubscription()
	var mu sync.Mutex
	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
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
			raven.CaptureError(err, nil)
			return
		}
		messageIds, err := gmailClient.ListMessageIds(message.HistoryId)
		if err != nil {
			raven.CaptureError(err, nil)
			return
		}
		fmt.Printf("%d new messages received\n", len(messageIds))

		processedMessages, err := gmailClient.ProcessMessages(messageIds)
		if err != nil {
			raven.CaptureError(err, nil)
			return
		}
		filteredMessages := filterMessages(processedMessages)
		for _, pm := range filteredMessages {
			link := pm.FindLink()
			if link != "" {
				err := pocketClient.Add(link)
				if err != nil {
					raven.CaptureError(err, nil)
					continue
				}
				err = gmailClient.Archive(pm.Id)
				if err != nil {
					raven.CaptureError(err, nil)
					continue
				}
			}
		}
	})

	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
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

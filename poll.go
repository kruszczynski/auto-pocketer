// Sample pubsub-quickstart creates a Google Cloud Pub/Sub topic.
package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/kruszczynski/auto-pocketer/gmail"
	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
)

func main() {
	gmail.Watch()
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "auto-pocketer"

	// Creates a client.
	credentialsPath := "secrets/pub_sub_credentials.json"
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new topic.
	// topicName := "incoming-emails"
	subName := "auto-pocketer-subscription"

	// Consume 10 messages.
	var mu sync.Mutex
	sub := client.Subscription(subName)
	cctx, _ := context.WithCancel(ctx)
	errr := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		fmt.Printf("Got message: %q\n", string(msg.Data))
		fmt.Printf("ID: %q\n", string(msg.ID))
		mu.Lock()
		defer mu.Unlock()
	})
	if errr != nil {
		panic(errr)
	}
}

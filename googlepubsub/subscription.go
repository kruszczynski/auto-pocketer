package googlepubsub

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type Message struct {
	EmailAddress string `json:"emailAddress"`
	HistoryId    uint64 `json:"historyId"`
}

func GetSubscription() *pubsub.Subscription {
	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	client, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	subName := os.Getenv("PUBSUB_SUBSCRIPTION_NAME")
	return client.Subscription(subName)
}

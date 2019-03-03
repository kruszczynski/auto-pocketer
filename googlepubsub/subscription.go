package googlepubsub

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Message struct {
	EmailAddress string `json:"emailAddress"`
	HistoryId    uint64 `json:"historyId"`
}

func GetSubscription() *pubsub.Subscription {
	credentialsPath := os.Getenv("PUB_SUB_CREDENTIALS_FILE_PATH")
	projectId := os.Getenv("GOOGLE_PROJECT_ID")
	client, err := pubsub.NewClient(context.Background(), projectId, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	subName := os.Getenv("PUBSUB_SUBSCRIPTION_NAME")
	return client.Subscription(subName)
}

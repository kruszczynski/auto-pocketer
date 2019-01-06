package googlepubsub

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const subName = "auto-pocketer-subscription"

// Sets your Google Cloud Platform project ID.
const projectID = "auto-pocketer"

type Message struct {
	EmailAddress string `json:"emailAddress"`
	HistoryId    uint64 `json:"historyId"`
}

func GetSubscription() *pubsub.Subscription {
	credentialsPath := "secrets/pub_sub_credentials.json"
	client, err := pubsub.NewClient(context.Background(), projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client.Subscription(subName)
}

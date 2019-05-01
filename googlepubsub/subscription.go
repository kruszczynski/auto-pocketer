package googlepubsub

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type Message struct {
	EmailAddress string `json:"emailAddress"`
	HistoryId    uint64 `json:"historyId"`
}

type Client struct {
	CredentialsPath  string
	ProjectID        string
	SubscriptionName string
}

func GetClient() *Client {
	return &Client{
		CredentialsPath:  os.Getenv("PUB_SUB_CREDENTIALS_FILE_PATH"),
		ProjectID:        os.Getenv("GOOGLE_PROJECT_ID"),
		SubscriptionName: os.Getenv("PUBSUB_SUBSCRIPTION_NAME"),
	}
}

func (c *Client) GetSubscription() (*pubsub.Subscription, error) {
	client, err := pubsub.NewClient(context.Background(), c.ProjectID, option.WithCredentialsFile(c.CredentialsPath))
	if err != nil {
		return nil, err
	}

	return client.Subscription(c.SubscriptionName), nil
}

package pocket

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const AddUrl = "https://getpocket.com/v3/add"

type PocketClient struct {
	ConsumerKey string
	AccessToken string
}

func GetClient() *PocketClient {
	return &PocketClient{
		AccessToken: os.Getenv("POCKET_ACCESS_TOKEN"),
		ConsumerKey: os.Getenv("POCKET_CONSUMER_KEY"),
	}
}

func (client *PocketClient) Add(link string) {
	resp, err := http.PostForm(AddUrl,
		url.Values{"url": {link}, "access_token": {client.AccessToken}, "consumer_key": {client.ConsumerKey}})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Pocket response status: %s\n", resp.Status)
}

package pocket

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const AddUrl = "https://getpocket.com/v3/add"

type httpClient interface {
	PostForm(url string, data url.Values) (*http.Response, error)
}

type PocketClient struct {
	ConsumerKey string
	AccessToken string
	HttpClient  httpClient
}

func GetClient() *PocketClient {
	return &PocketClient{
		AccessToken: os.Getenv("POCKET_ACCESS_TOKEN"),
		ConsumerKey: os.Getenv("POCKET_CONSUMER_KEY"),
		HttpClient:  &http.Client{},
	}
}

func (client *PocketClient) Add(link string) error {
	resp, err := client.HttpClient.PostForm(
		AddUrl,
		url.Values{
			"url":          {link},
			"access_token": {client.AccessToken},
			"consumer_key": {client.ConsumerKey},
		},
	)

	if err != nil {
		return err
	}
	fmt.Printf("Pocket response status: %s\n", resp.Status)
	return nil
}

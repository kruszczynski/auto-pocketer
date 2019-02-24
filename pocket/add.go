package pocket

import (
	"fmt"
	"net/http"
	"net/url"
)

type PocketClient struct {
	ConsumerKey string
	AccessToken string
}

func GetClient() *PocketClient {
	consumerKey := "83146-d20c312392e77ac6a0da235c"
	accessToken := "a1c4ed5d-b4fb-4fee-a558-f9d9e9"
	return &PocketClient{AccessToken: accessToken, ConsumerKey: consumerKey}
}

func (client *PocketClient) Add(link string) {
	resp, err := http.PostForm("https://getpocket.com/v3/add",
		url.Values{"url": {link}, "access_token": {client.AccessToken}, "consumer_key": {client.ConsumerKey}})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Pocket response status: %s\n", resp.Status)
}

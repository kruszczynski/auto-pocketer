package pocket

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestGetClient(t *testing.T) {
	os.Setenv("POCKET_ACCESS_TOKEN", "letestaccesstoken")
	os.Setenv("POCKET_CONSUMER_KEY", "leconsumerkey")

	client := GetClient()
	if client.AccessToken != "letestaccesstoken" {
		t.Error("Wrong Access Token")
	}
	if client.ConsumerKey != "leconsumerkey" {
		t.Error("Wrong Consumer Key")
	}
}

func TestAdd_Success(t *testing.T) {
	client := &PocketClient{
		AccessToken: "letestaccesstoken",
		ConsumerKey: "leconsumerkey",
		HttpClient:  &mockPocketClient{t: t},
	}

	err := client.Add("https://example.com")
	if err != nil {
		t.Error("This should not have returned an error")
	}
}

func TestAdd_Fail(t *testing.T) {
	client := &PocketClient{
		AccessToken: "letestaccesstoken",
		ConsumerKey: "leconsumerkey",
		HttpClient:  &mockFailedPocketClient{},
	}

	err := client.Add("https://example.com")

	if err.Error() != "Somehow request to Pocket Failed" {
		t.Error("Wrong error returned")
	}
}

type mockPocketClient struct {
	t *testing.T
}

func (c *mockPocketClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	if url != AddUrl {
		c.t.Error("Wrong Pocket API Url")
	}

	if data["url"][0] != "https://example.com" {
		c.t.Error("Incorrect URL argument")
	}

	if data["access_token"][0] != "letestaccesstoken" {
		c.t.Error("Incorrect Access Token")
	}

	if data["consumer_key"][0] != "leconsumerkey" {
		c.t.Error("Incorrect Consumer Key")
	}

	return &http.Response{Status: "200 OK"}, nil
}

type mockFailedPocketClient struct{}

func (c *mockFailedPocketClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return nil, errors.New("Somehow request to Pocket Failed")
}

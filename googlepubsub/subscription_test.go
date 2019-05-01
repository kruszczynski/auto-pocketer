package googlepubsub

import (
	"os"
	"testing"
)

func TestGetClient(t *testing.T) {
	os.Setenv("PUB_SUB_CREDENTIALS_FILE_PATH", "credentialsfile")
	os.Setenv("GOOGLE_PROJECT_ID", "googleprojectid")
	os.Setenv("PUBSUB_SUBSCRIPTION_NAME", "thequeuename")

	client := GetClient()
	if client.CredentialsPath != "credentialsfile" {
		t.Error("Wrong Credentials Path")
	}
	if client.ProjectID != "googleprojectid" {
		t.Error("Wrong Project ID")
	}
	if client.SubscriptionName != "thequeuename" {
		t.Error("Wrong Subscription Name")
	}
}

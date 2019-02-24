package gmail

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/gmail/v1"
)

func (client *Client) Watch() {
	topicName := os.Getenv("PUBSUB_TOPIC_NAME")
	watchRequest := &gmail.WatchRequest{TopicName: topicName, LabelIds: []string{"INBOX"}}
	r, err := client.service.Users.Watch(User, watchRequest).Do()
	if err != nil {
		log.Fatalf("Unsuccessful watch request: %v", err)
	}
	fmt.Printf("Watch successful, expiration: %d, historyId: %d\n", r.Expiration, r.HistoryId)
	if client.lastHistoryId == 0 {
		fmt.Println("Setting initial history id")
		client.lastHistoryId = r.HistoryId
	}
}

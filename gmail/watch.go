package gmail

import (
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

func (client Client) Watch() uint64 {
	watchRequest := &gmail.WatchRequest{TopicName: "projects/auto-pocketer/topics/incoming-emails", LabelIds: []string{"INBOX"}}
	r, err := client.service.Users.Watch(User, watchRequest).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	fmt.Printf("Watch successful, expiration: %d\n", r.Expiration)
	return r.HistoryId
}

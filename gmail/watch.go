package gmail

import (
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

func Watch() {
	client := getClient()

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	watchRequest := &gmail.WatchRequest{TopicName: "projects/auto-pocketer/topics/incoming-emails", LabelIds: []string{"INBOX"}}
	r, err := srv.Users.Watch(user, watchRequest).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	fmt.Println("Expiration:")
	fmt.Printf("%s\n", r.Expiration)
	fmt.Println("History:")
	fmt.Printf("%s\n", r.HistoryId)
}

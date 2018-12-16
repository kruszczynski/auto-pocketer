package gmail

import (
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

func GetMsg(startHistoryId uint64, historyId uint64) {
	client := getClient()

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	c := srv.Users.History.List(user)
	c.StartHistoryId(startHistoryId)
	r, err := c.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	fmt.Println("Histories:")
	for _, h := range r.History {
		fmt.Printf("HISTORY ID: %d", h.Id)
		fmt.Printf("Labels added: %d\n", len(h.LabelsAdded))
		fmt.Printf("Labels removed: %d\n", len(h.LabelsRemoved))
		fmt.Printf("MSG added: %d\n", len(h.MessagesAdded))
		fmt.Printf("MSG removed: %d\n\n", len(h.MessagesDeleted))
		// for _, m := range h.Messages {
		// 	r, _ := srv.Users.Messages.Get(user, m.Id).Do()
		// 	for _, mp := range r.Payload.Parts {
		// 		dec, _ := base64.StdEncoding.DecodeString(mp.Body.Data)
		// 		fmt.Println(string(dec))
		// 	}
		// }
	}
}

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the consumer_key: ")
	consumerKey, _ := reader.ReadString('\n')
	consumerKey = strings.TrimSpace(consumerKey)

	// obtain the request token
	res, err := http.PostForm("https://getpocket.com/v3/oauth/request",
		url.Values{"consumer_key": {consumerKey}, "redirect_uri": {"localhost"}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v+\n", res)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	requestCode := strings.Split(string(bodyBytes), "=")[1]
	url := "https://getpocket.com/auth/authorize?request_token=" + requestCode + "&redirect_uri=localhost"
	cmd := exec.Command("open", url)
	cmd.Run()
	fmt.Println(requestCode)
}

package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/context"
)

const endpoint = "https://slack.com/api/users.admin.invite?token=%s&email=%s"

// SlackInvite ...
func SlackInvite(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("SLACK_TOKEN")
	email := r.URL.Query().Get("email")

	u := fmt.Sprintf(endpoint, token, email)
	w.Header().Set("content-type", "application/json")

	client := context.Get(r, "client").(interface {
		Get(url string) (*http.Response, error)
	})
	resp, err := client.Get(u)
	if err != nil {
		fmt.Fprintln(w, `{"ok":false}`)
		return
	}

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintln(w, string(b))
}

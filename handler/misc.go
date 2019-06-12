package handler

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	cip "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

func githubEndpoint(code string) string {
	// default method
	method := "authorize"

	// env
	id := os.Getenv("GITHUB_CLIENT_ID")
	secret := os.Getenv("GITHUB_CLIENT_SECRET")
	callback := os.Getenv("GITHUB_CALLBACK")

	if code != "" {
		method = "access_token"
	}

	// format
	ep := fmt.Sprintf(
		"%s/%s?client_id=%s&callback_url=%s",
		oauthEndpoint,
		method,
		id,
		callback,
	)

	// apend if code exists
	if code != "" {
		ep += fmt.Sprintf("&client_secret=%s&code=%s", secret, code)
	}

	return ep
}

func sessionUser(r *http.Request) map[string]string {
	ss := context.Get(r, "session").(interface {
		GetUser(r *http.Request) map[string]string
	})

	su := ss.GetUser(r)
	if su == nil {
		return nil
	}

	u := context.Get(r, "userService").(interface {
		AccountDetails(token string) (*cip.GetUserOutput, error)
	})

	token := su["token"]
	o, err := u.AccountDetails(token)
	if err != nil {
		return nil
	}

	if len(o.UserAttributes) == 0 {
		return nil
	}

	ua := map[string]string{}
	for _, v := range o.UserAttributes {
		ua[*v.Name] = *v.Value
	}

	return ua
}

func render(w http.ResponseWriter, base, content string, ctx interface{}) {
	tmpl, err := template.New("").ParseFiles(
		fmt.Sprintf("assets/templates/layout/%s.tmpl", base),
		fmt.Sprintf("assets/templates/%s.tmpl", content),
		fmt.Sprintf("assets/templates/main-nav.tmpl"),
		// main css file
		fmt.Sprintf("assets/css/main.css"),
		// main javascript file
		fmt.Sprintf("assets/scripts/main.js"),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if content == "notfound" {
		w.WriteHeader(http.StatusNotFound)
	}
	tmpl.ExecuteTemplate(w, "layout", ctx)
}

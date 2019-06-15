package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

const (
	oauthEndpoint = "https://github.com/login/oauth"
)

// Render parses templates and writes them into the the passed in ResponseWriter
// interface
func Render(w http.ResponseWriter, base, content string, ctx interface{}) {
	tmpl, err := template.New("").ParseFiles(
		fmt.Sprintf("assets/templates/layout/%s.tmpl", base),
		fmt.Sprintf("assets/templates/%s.tmpl", content),
		fmt.Sprintf("assets/templates/main-nav.tmpl"),
		fmt.Sprintf("assets/templates/user-card.tmpl"),
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

// GithubEndpoint parses endpoint for github request
func GithubEndpoint(code string) string {
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

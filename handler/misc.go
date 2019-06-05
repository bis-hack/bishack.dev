package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

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

func render(w http.ResponseWriter, base, content string, ctx interface{}) {
	tmpl, err := template.New("").ParseFiles(
		fmt.Sprintf("templates/layout/%s.tmpl", base),
		fmt.Sprintf("templates/%s.tmpl", content),
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

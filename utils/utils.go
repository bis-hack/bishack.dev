package utils

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/golang-commonmark/markdown"
)

const (
	oauthEndpoint = "https://github.com/login/oauth"
)

func md(input string) template.HTML {
	md := markdown.New(markdown.Linkify(false))
	out := md.RenderToString([]byte(input))
	out = strings.Replace(out, "<pre>", "<pre class=\"prettyprint\">", -1)
	return template.HTML(out)
}

func date(fmt string, input int64) string {
	t := time.Unix(input, 0)
	return t.Format(fmt)
}

// Render parses templates and writes them into the the passed in ResponseWriter
// interface
func Render(w http.ResponseWriter, base, content string, ctx interface{}) {
	fns := template.FuncMap{
		"md":   md,
		"date": date,
	}
	tmpl, err := template.New("").Funcs(fns).ParseFiles(
		fmt.Sprintf("assets/templates/layout/%s.tmpl", base),
		fmt.Sprintf("assets/templates/%s.tmpl", content),
		// components
		"assets/templates/components/main-nav.tmpl",
		"assets/templates/components/user-card.tmpl",
		"assets/templates/components/posts.tmpl",
		// stylesheets
		"assets/css/main.css",
		// scripts
		"assets/scripts/turbolinks.js",
		"assets/scripts/axios.js",
		"assets/scripts/main.js",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if content == "notfound" {
		w.WriteHeader(http.StatusNotFound)
	}
	_ = tmpl.ExecuteTemplate(w, "layout", ctx)
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

package session

import "github.com/gorilla/sessions"

// Flash ...
type Flash struct {
	Type  string
	Value string
}

// Client ...
type Client struct {
	Store *sessions.CookieStore
}

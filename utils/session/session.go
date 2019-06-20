package session

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// New ...
func New() *Client {
	return &Client{store}
}

// SetUser ...
func (s *Client) SetUser(
	w http.ResponseWriter,
	r *http.Request,
	username,
	token string,
) {
	session, _ := s.Store.Get(r, "user")
	session.Values["username"] = username
	session.Values["token"] = token
	_ = session.Save(r, w)
}

// GetUser ...
func (s *Client) GetUser(r *http.Request) map[string]string {
	session, _ := s.Store.Get(r, "user")

	username := session.Values["username"]
	token := session.Values["token"]

	if username == nil || token == nil {
		return nil
	}

	return map[string]string{
		"username": username.(string),
		"token":    token.(string),
	}
}

// DeleteUser ...
func (s *Client) DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := s.Store.Get(r, "user")

	session.Values["username"] = nil
	session.Values["token"] = nil

	_ = session.Save(r, w)
}

// SetFlash sets the flash message with the given
// type and value
func (s *Client) SetFlash(
	w http.ResponseWriter,
	r *http.Request,
	t,
	v string,
) {
	session, _ := s.Store.Get(r, "notification")
	session.AddFlash(fmt.Sprintf("%s<>%s", t, v))
	_ = session.Save(r, w)
}

// GetFlash ...
func (s *Client) GetFlash(w http.ResponseWriter, r *http.Request) *Flash {
	session, _ := s.Store.Get(r, "notification")

	if flashes := session.Flashes(); len(flashes) > 0 {
		chunks := strings.Split(flashes[0].(string), "<>")
		f := &Flash{chunks[0], chunks[1]}
		_ = session.Save(r, w)
		return f
	}

	return nil
}

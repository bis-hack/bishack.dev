package session

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// Get wraps the session store Getter
func Get(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

// SetUser ...
func SetUser(w http.ResponseWriter, r *http.Request, email, token string) {
	session, _ := store.Get(r, "user")
	session.Values["email"] = email
	session.Values["token"] = token
	session.Save(r, w)
}

// GetUser ...
func GetUser(r *http.Request) map[string]string {
	session, err := store.Get(r, "user")
	if err != nil {
		return nil
	}

	email := session.Values["email"]
	token := session.Values["token"]

	if email == nil || token == nil {
		return nil
	}

	return map[string]string{
		"email": email.(string),
		"token": token.(string),
	}
}

// DeleteUser ...
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")

	session.Values["email"] = nil
	session.Values["token"] = nil

	session.Save(r, w)
}

// SetFlash sets the flash message with the given
// type and value
func SetFlash(w http.ResponseWriter, r *http.Request, t, v string) {
	session, _ := store.Get(r, "notification")
	session.AddFlash(fmt.Sprintf("%s<>%s", t, v))
	session.Save(r, w)
}

// GetFlash ...
func GetFlash(w http.ResponseWriter, r *http.Request) *Flash {
	session, _ := store.Get(r, "notification")

	if flashes := session.Flashes(); len(flashes) > 0 {
		chunks := strings.Split(flashes[0].(string), "<>")
		f := &Flash{chunks[0], chunks[1]}
		session.Save(r, w)
		return f
	}

	return nil
}

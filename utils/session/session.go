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

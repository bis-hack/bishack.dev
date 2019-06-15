package middleware

import (
	"net/http"
	"regexp"

	"github.com/gorilla/context"
)

// AuthRedirects middleware will redirect user to the root page
// if the user is trying to access auth based endpoint like: /login, /signup etc
func AuthRedirects(h http.Handler) http.Handler {
	rx := regexp.MustCompile(`(?i)^/(signup|login|verify)`)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.Get(r, "user")
		if r.URL.Path != "/" && rx.MatchString(r.URL.Path) && user != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

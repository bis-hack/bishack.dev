package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"bishack.dev/handler"
	"github.com/gorilla/csrf"
	"github.com/gorilla/pat"
)

// PUBLICFOLDER ...
const PUBLICFOLDER = "public"

func main() {
	// init route
	r := pat.New()

	// not found
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			handler.Home(w, r)
		} else {
			handler.NotFound(w, r)
		}
	})

	// handlers et al
	r.Get("/signup", handler.Signup)
	r.Get("/verify", handler.Verify)
	r.Post("/signup", handler.FinishSignup)

	// csrf
	csrfSecure := false
	if regexp.MustCompile("`(?i)stag|prod").MatchString(os.Getenv("UP_STAGE")) {
		csrfSecure = true
	}
	protect := csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.Secure(csrfSecure))

	// push
	// read from public dir
	o, _ := exec.Command("ls", "-a", PUBLICFOLDER).Output()
	list := strings.Split(string(o), "\n")
	// ginore space and dots
	assets := list[2 : len(list)-1]

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(
		port,
		protect(push(assets, r)),
	))
}

// HTTP2/Push middleware baby!
func push(assets []string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			for _, a := range assets {
				err := pusher.Push("/"+a, nil)
				if err != nil {
					log.Println("Could not push", a)
				} else {
					log.Println(a, "successfully pushed to client")
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

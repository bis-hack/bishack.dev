package main

import (
	"log"
	"net/http"
	"os"

	"bishack.dev/handler"
	"github.com/gorilla/pat"
)

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

	// launch
	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, r))
}

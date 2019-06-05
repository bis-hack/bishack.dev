package handler

import "net/http"

// Home ...
func Home(w http.ResponseWriter, r *http.Request) {
	render(w, "main", "index", struct {
		Title string
	}{"Home"})
}

// NotFound ...
func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, "main", "notfound", struct {
		Title string
	}{"Not Found"})
}

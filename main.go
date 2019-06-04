package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := ":" + os.Getenv("PORT")
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "soon...")
}

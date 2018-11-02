package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"summaryGenerater/app"
)

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	text := r.FormValue("text")
	delimiter := r.FormValue("delimiter")
	app.Generate(text, delimiter)
	fmt.Fprintf(w, "text: %v\n", text)
	fmt.Fprintf(w, "delimiter: %v\n", delimiter)
}

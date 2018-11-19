package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ramenjuniti/lexrank"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Cache-Control, Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	text := r.FormValue("text")
	delimiter := r.FormValue("delimiter")
	summary := lexrank.New()
	summary.Summarize(text, delimiter)
	data, err := json.Marshal(summary.LexRankScores)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
}

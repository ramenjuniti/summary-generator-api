package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	w.Header().Set("Content-Type", "application/json")

	text := r.FormValue("text")
	delimiter := r.FormValue("delimiter")

	maxLine, err := strconv.ParseInt(r.FormValue("maxLine"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	maxCharacter, err := strconv.ParseInt(r.FormValue("maxCharacter"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	threshold, err := strconv.ParseFloat(r.FormValue("threshold"), 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tolerance, err := strconv.ParseFloat(r.FormValue("tolerance"), 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	damping, err := strconv.ParseFloat(r.FormValue("damping"), 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	lambda, err := strconv.ParseFloat(r.FormValue("lambda"), 64)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	summary := lexrank.New(
		lexrank.MaxLines(int(maxLine)),
		lexrank.MaxCharacters(int(maxCharacter)),
		lexrank.Threshold(threshold),
		lexrank.Tolerance(tolerance),
		lexrank.Damping(damping),
		lexrank.Lambda(lambda),
	)
	summary.Summarize(text, delimiter)
	data, err := json.Marshal(summary.LexRankScores)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(data))
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ramenjuniti/lexrank-mmr"
)

const (
	defaultMaxLines      = 0
	defaultMaxCharacters = 0
	defaultThreshold     = 0.001
	defaultTolerance     = 0.0001
	defaultDamping       = 0.85
	defaultLambda        = 1.0
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

	maxLines := defaultMaxLines
	maxCharacters := defaultMaxCharacters
	threshold := defaultThreshold
	tolerance := defaultTolerance
	damping := defaultDamping
	lambda := defaultLambda
	var err error

	if r.FormValue("maxLines") != "" {
		maxLines, err = strconv.Atoi(r.FormValue("maxLines"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	if r.FormValue("maxCharacters") != "" {
		maxCharacters, err = strconv.Atoi(r.FormValue("maxCharacters"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	if r.FormValue("threshold") != "" {
		threshold, err = strconv.ParseFloat(r.FormValue("threshold"), 64)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	if r.FormValue("tolerance") != "" {
		tolerance, err = strconv.ParseFloat(r.FormValue("tolerance"), 64)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	if r.FormValue("damping") != "" {
		damping, err = strconv.ParseFloat(r.FormValue("damping"), 64)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	if r.FormValue("lambda") != "" {
		lambda, err = strconv.ParseFloat(r.FormValue("lambda"), 64)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

	summary := lexrank.New(
		lexrank.MaxLines(maxLines),
		lexrank.MaxCharacters(maxCharacters),
		lexrank.Threshold(threshold),
		lexrank.Tolerance(tolerance),
		lexrank.Damping(damping),
		lexrank.Lambda(lambda),
	)
	summary.Summarize(text)
	data, err := json.Marshal(summary)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	fmt.Fprint(w, string(data))
}

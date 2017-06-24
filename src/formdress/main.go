package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"encoding/json"

	"context"
	"flag"
	"net/url"
	"strings"
	"time"
)

func CheckURL(u string) error {
	v, err := url.Parse(u)
	if err != nil {
		return err
	}

	if v.Hostname() != "docs.google.com" {
		return errors.New("Invalid Google Forms address")
	}

	if !strings.HasSuffix(v.EscapedPath(), "/viewform") {
		return errors.New("Please, use the public URL")
	}

	return nil
}

func JSONError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(struct {
		Error string
	}{err.Error()})
}

func FormDressHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Check form url
	url := r.URL.Query().Get("url")
	if err := CheckURL(url); err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	// Fetch form
	resp, err := httpClient.Get(url)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Extract data
	form, err := FormExtract(resp.Body)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(form)
}

var addr = flag.String("l", "0.0.0.0:8000", "Bind Address")
var distDir = flag.String("d", "dist", "Static Assets Directory")
var fetch = flag.String("f", "", "Just fetch the Google Form data")

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

func FetchAndExit(url string) {
	res, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	form, err := FormExtract(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := json.MarshalIndent(form, "", "    ")
	log.Println(string(out))
}

func main() {
	flag.Parse()

	if *fetch != "" {
		FetchAndExit(*fetch)
		return
	}

	log.Println("Serving assets from: " + *distDir)
	log.Println("Bind Address: " + *addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	mux := http.NewServeMux()
	mux.HandleFunc("/formdress", FormDressHandler)
	mux.Handle("/", http.FileServer(http.Dir(*distDir)))

	server := &http.Server{Addr: *addr, Handler: mux}
	go func() {
		log.Println("Listening...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

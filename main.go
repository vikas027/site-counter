package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gopkg.in/redis.v4"
)

var (
	client *redis.Client
)

type jsonStruct struct {
	Status string `json:"status"`
}

func health(w http.ResponseWriter, r *http.Request) {
	// TO-DO: Check the status of Redis too
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	jsonResp := jsonStruct{Status: "ok"}
	json.NewEncoder(w).Encode(jsonResp)
}

func allowedURI(w http.ResponseWriter, r *http.Request) {
	allowedURI := [2]string{
		"health",
		"counter",
	}
	w.WriteHeader(501)
	io.WriteString(w, fmt.Sprintln("Use one of these URIs:", allowedURI))
}

func counter(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	addrs, _ := net.LookupIP(hostname)

	if r.URL.Path == "/favicon.ico" {
		io.WriteString(w, "favicon")
		return
	}
	count, err := client.Incr("counter").Result()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Redis is unhappy")
	} else {
		if os.Getenv("SHOW_IP") == "false" {
			w.WriteHeader(200)
			io.WriteString(w, fmt.Sprintln(hostname, " - ", "View Count: ", count))
		} else {
			w.WriteHeader(200)
			io.WriteString(w, fmt.Sprintln(hostname, " - ", addrs, " - ", "View Count: ", count))
		}
	}
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	http.HandleFunc("/", allowedURI)
	http.HandleFunc("/health", health)
	if os.Getenv("NEW_RELIC_LICENSE_KEY") == "" {
		http.HandleFunc("/counter", counter)
	} else {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("site-counter"),
			newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		)
		if err != nil {
			log.Panic("Error in using New Relic Agent")
			os.Exit(1)
		}
		http.HandleFunc(newrelic.WrapHandleFunc(app, "/counter", counter))
	}
	http.ListenAndServe(":80", nil)
}

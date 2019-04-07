package main

import (
	"fmt"
	"gopkg.in/redis.v4"
	"io"
	"net/http"
	"os"
	"net"
)

var (
	client *redis.Client
)

func hello(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	addrs, _ := net.LookupIP(hostname)

	if r.URL.Path == "/favicon.ico" {
		io.WriteString(w, "favicon")
		return
	}
	count, err := client.Incr("counter").Result()
	if err != nil {
		io.WriteString(w, "Redis is unhappy")
	} else {
		if os.Getenv("SHOW_IP") == "false" {
			io.WriteString(w, fmt.Sprintln(hostname, " - " , "View Count: ", count))
		} else {
			io.WriteString(w, fmt.Sprintln(hostname, " - ", addrs, " - " , "View Count: ", count))
		}

	}
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	http.HandleFunc("/", hello)
	http.ListenAndServe(":80", nil)
}

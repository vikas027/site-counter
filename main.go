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
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	addrs, _ := net.LookupIP(hostname)
	for _, addr := range addrs {
	    if ipv4 := addr.To4(); ipv4 != nil {
	        fmt.Println("IPv4: ", ipv4)
	    }
	}

	if r.URL.Path == "/favicon.ico" {
		io.WriteString(w, "favicon")
		return
	}
	count, err := client.Incr("counter").Result()
	if err != nil {
		io.WriteString(w, "Redis is unhappy")
	} else {
		io.WriteString(w, fmt.Sprintln(hostname, " - ", addrs, " - " , "View Count: ", count))
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

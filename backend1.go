package main

import (
	"fmt"
	"log"
	"net/http"
)
func main() {
	port := "9001"
	log.Printf("Starting backend server on :%s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend 1 received request from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Hello from Backend 1!\n")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend 1 health check received from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Backend 1 is healthy!\n")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}



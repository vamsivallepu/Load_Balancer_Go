package main

import (
	"fmt"
	"log"
	"net/http"
)
func main() {
	port := "9002"
	log.Printf("Starting backend server on :%s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend 2 received request from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Hello from Backend 2!\n")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Backend 2 health check received from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Backend 2 is healthy!\n")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}



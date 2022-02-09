package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

// Home page handler
func Home(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "Hello, it's home page!")
	if err != nil {
		log.Println("Got error when return request", err)
	}

	fmt.Println(fmt.Sprintf("Bytes written: %d", n))
}

// About page handler
func About(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "Hello, it's about page!")
	if err != nil {
		log.Println("Got error when return request", err)
	}

	fmt.Println(fmt.Sprintf("Bytes written: %d", n))
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	http.ListenAndServe(portNumber, nil)
}

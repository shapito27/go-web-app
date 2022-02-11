package main

import (
	"fmt"
	"net/http"
	"github.com/shapito27/go-web-app/pkg/handlers"
)

//app listen this port
const portNumber = ":8080"

func main() {
	//routes
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	
	http.ListenAndServe(portNumber, nil)
}

package main

import (
	"net/http"
	"orders/internal/app"

	"github.com/gorilla/mux"
)

func main() {
	var myApp app.App
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", myApp.HomeHandler)
	// r.HandleFunc("/api/hello", helloHandler).Methods("GET")
	// r.HandleFunc("/api/greet/{name}", greetHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}

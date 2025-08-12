package main

import (
	"fmt"
	"log"
	"net/http"
	"orders/internal/app"

	"github.com/gorilla/mux"
)

func main() {

	// config init
	dbuser := "postgres"
	dbpassword := "1234"
	dbname := "order"

	myApp, err := app.NewApp(fmt.Sprintf("postgres://%s:%s@db:5432/%s", dbuser, dbpassword, dbname))
	if err != nil {
		log.Fatal("Failed to init")
	}
	defer myApp.Close()

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", myApp.HomeHandler)
	r.HandleFunc("/api/add", myApp.Insert)
	r.HandleFunc("/order/{order_uid}", myApp.GetById)
	log.Println("Server is up")
	http.ListenAndServe(":8080", r)
}

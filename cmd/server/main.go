package main

import (
	"fmt"
	"log"
	"net/http"
	"orders/internal/app"
	"orders/internal/config"

	"github.com/gorilla/mux"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	myApp, err := app.NewApp(fmt.Sprintf("postgres://%s:%s@db:5432/%s",
		config.DBUser,
		config.DBPassword,
		config.DBName),
	)
	if err != nil {
		log.Fatal("Failed to init")
	}
	defer myApp.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", myApp.HomeHandler)
	r.HandleFunc("/api/add", myApp.Insert)
	r.HandleFunc("/order/{order_uid}", myApp.GetById)
	log.Println("Server is up")
	http.ListenAndServe(":8080", r)
}

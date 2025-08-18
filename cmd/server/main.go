package main

import (
	"fmt"
	"log"
	"net/http"
	_ "orders/docs"
	"orders/internal/app"
	"orders/internal/config"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Orders Service API
// @version 1.0
// @description API documentation for Orders microservice
// @host localhost:8081
// @BasePath /
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
	r.HandleFunc("/order/random/{count}", myApp.GetNOrders).Methods("GET")
	r.HandleFunc("/order/{order_uid}", myApp.GetById).Methods("GET")

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	log.Println("Server is up")
	http.ListenAndServe(":8081", r)
}

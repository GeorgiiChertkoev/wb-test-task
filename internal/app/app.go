package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"orders/internal/repository"
	"orders/pkg/models"

	"github.com/gorilla/mux"
)

type App struct {
	// conn *pgx.Conn
	repo repository.OrderRepo
}

func NewApp(dburl string) (*App, error) {
	app := &App{}
	err := app.repo.InitRepo(dburl)
	if err != nil {
		log.Printf("Failed to init repo: %v", err)
		return nil, err
	}

	return app, nil
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Opened home page")
	w.Write([]byte("Welcome Home!\n\n"))
	rows := a.repo.GetAllRows()
	for rows.Next() {
		var new_order models.Order

		err := rows.Scan(&new_order.OrderUID, &new_order.TrackNumber, &new_order.Entry, &new_order.Locale,
			&new_order.InternalSignature, &new_order.CustomerID, &new_order.DeliveryService,
			&new_order.Shardkey, &new_order.SmID, &new_order.DateCreated, &new_order.OofShard)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		json_data, err := json.MarshalIndent(new_order, "", "\t")
		if err != nil {
			log.Printf("Error making json: %v", err)
		}
		fmt.Fprintf(w, "%s\n", json_data)
		// fmt.Fprintf(w, "Order UID: %s\n", orderUID)
		// fmt.Fprintf(w, "Track Number: %s\n", trackNumber)
		// fmt.Fprintf(w, "Entry: %s\n", entry)
		// fmt.Fprintf(w, "Locale: %s\n", locale)
		// fmt.Fprintf(w, "Internal Signature: %s\n", internalSignature)
		// fmt.Fprintf(w, "Customer ID: %s\n", customerID)
		// fmt.Fprintf(w, "Delivery Service: %s\n", deliveryService)
		// fmt.Fprintf(w, "Shardkey: %s\n", shardkey)
		// fmt.Fprintf(w, "SM ID: %d\n", smID)
		// fmt.Fprintf(w, "Date Created: %s\n", dateCreated.Format(time.RFC3339))
		// fmt.Fprintf(w, "OOF Shard: %s\n", oofShard)
		fmt.Fprintf(w, "----------------------------------------\n")
	}
}

func (a *App) GetById(w http.ResponseWriter, r *http.Request) {
	order_uid := mux.Vars(r)["order_uid"]
	order, found, err := a.repo.Find(order_uid)
	if !found {
		fmt.Fprintf(w, "order %v not found\n", order_uid)
		return
	} else if err != nil {
		fmt.Fprintf(w, "internal error\n")
		log.Printf("Error while searching by id: %v", err)
		return
	}
	json_data, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		log.Printf("Error making json: %v", err)
	}
	fmt.Fprintf(w, "%s\n", json_data)

}

func (a *App) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Opened insert page")

	a.repo.Store(models.MakeRandomOrder())

}

func (a *App) Close() {
	// a.repo.conn.Close(context.Background()) //TODO
}

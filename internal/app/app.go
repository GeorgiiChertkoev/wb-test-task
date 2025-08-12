package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"orders/internal/models"
	"orders/internal/repository"

	"github.com/gorilla/mux"
)

type App struct {
	// conn *pgx.Conn
	repo repository.OrderRepo
}

/*
GET http://localhost:8081/order/<order_uid> должен вернуть JSON с информацией о заказе
-> получение по order_id

добавление из json
*/

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (a *App) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Opened insert page")
	uid := randomString(19)
	log.Printf("random uid is %v\n", uid)

	order := models.Order{
		OrderUID:          uid,
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:          "1",
		Delivery: models.Delivery{
			OrderUID: uid,
			Name:     "Test Testov",
			Phone:    "+9720000000",
			Zip:      "2639809",
			City:     "Kiryat Mozkin",
			Address:  "Ploshad Mira 15",
			Region:   "Kraiot",
			Email:    "test@gmail.com",
		},
		Payment: models.Payment{
			OrderUID:     uid,
			Transaction:  uid,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ID:          1,
				OrderUID:    uid,
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}

	a.repo.Store(&order)

}

func (a *App) Close() {
	// a.repo.conn.Close(context.Background()) //TODO
}

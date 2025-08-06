package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"orders/internal/models"

	"github.com/jackc/pgx/v5"
)

type App struct {
	conn *pgx.Conn
}

func NewApp(dburl string) (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dburl)
	if err != nil {
		log.Printf("Couldn't connect to db by url %s: %s\n", dburl, err)
		return nil, err
	}
	app := &App{
		conn: conn,
	}

	return app, nil
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Opened home page")
	w.Write([]byte("Welcome Home!\n\n"))
	rows, err := a.conn.Query(context.Background(), `SELECT order_uid, track_number, entry, locale,
       internal_signature, customer_id,delivery_service,shardkey, sm_id,
       date_created, oof_shard FROM orders`)
	if err != nil {
		fmt.Fprintf(w, "Shit happened: %s", err)
		return
	}
	for rows.Next() {
		var (
			orderUID          string
			trackNumber       string
			entry             string
			locale            string
			internalSignature string
			customerID        string
			deliveryService   string
			shardkey          string
			smID              int
			dateCreated       time.Time
			oofShard          string
		)

		err := rows.Scan(&orderUID, &trackNumber, &entry, &locale,
			&internalSignature, &customerID, &deliveryService,
			&shardkey, &smID, &dateCreated, &oofShard)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		fmt.Fprintf(w, "Order UID: %s\n", orderUID)
		fmt.Fprintf(w, "Track Number: %s\n", trackNumber)
		fmt.Fprintf(w, "Entry: %s\n", entry)
		fmt.Fprintf(w, "Locale: %s\n", locale)
		fmt.Fprintf(w, "Internal Signature: %s\n", internalSignature)
		fmt.Fprintf(w, "Customer ID: %s\n", customerID)
		fmt.Fprintf(w, "Delivery Service: %s\n", deliveryService)
		fmt.Fprintf(w, "Shardkey: %s\n", shardkey)
		fmt.Fprintf(w, "SM ID: %d\n", smID)
		fmt.Fprintf(w, "Date Created: %s\n", dateCreated.Format(time.RFC3339))
		fmt.Fprintf(w, "OOF Shard: %s\n", oofShard)
		fmt.Fprintf(w, "----------------------------------------\n")
	}
}

func (a *App) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Opened insert page")

	tx, err := a.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v\n", err)
	}
	defer tx.Rollback(context.Background()) // Безопасный откат если не будет Commit

	// Тестовые данные
	order := models.Order{
		OrderUID:          "pgx_test_123",
		TrackNumber:       "TRACK_PGX_001",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer",
		DeliveryService:   "pgx_delivery",
		Shardkey:          "1",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// 1. Вставка заказа
	_, err = tx.Exec(context.Background(), `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		log.Fatalf("Insert order failed: %v\n", err)
	}

	// 2. Вставка доставки
	_, err = tx.Exec(context.Background(), `
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID, "PGX Test User", "+1234567890", "123456", "PGX City", "PGX Address", "PGX Region", "pgx@test.com")
	if err != nil {
		log.Fatalf("Insert delivery failed: %v\n", err)
	}

	// 3. Вставка платежа
	_, err = tx.Exec(context.Background(), `
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, "pgx_txn_123", "", "USD", "pgx_pay", 1000, time.Now().Unix(), "pgx_bank", 500, 500, 0)
	if err != nil {
		log.Fatalf("Insert payment failed: %v\n", err)
	}

	// 4. Вставка товаров
	items := []struct {
		ChrtID      int64
		TrackNumber string
		Price       int
		Rid         string
		Name        string
		Sale        int
		Size        string
		TotalPrice  int
		NmID        int64
		Brand       string
		Status      int
	}{
		{9934930, "TRACK_PGX_001", 453, "rid_pgx_1", "PGX Product 1", 30, "0", 317, 2389212, "PGX Brand", 202},
		{11223344, "TRACK_PGX_001", 1000, "rid_pgx_2", "PGX Product 2", 10, "1", 900, 55667788, "PGX Brand", 200},
	}

	for _, item := range items {
		_, err = tx.Exec(context.Background(), `
			INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
				sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Fatalf("Insert item failed: %v\n", err)
		}
	}

	// Фиксация транзакции
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Commit failed: %v\n", err)
	}

}
func (a *App) Close() {
	a.conn.Close(context.Background())
}

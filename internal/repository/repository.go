package repository

// пакет отвечающий за получение и сохранение информации

import (
	"context"
	"log"
	"sync"
	"time"

	"orders/internal/models"

	"github.com/jackc/pgx/v5"
)

type OrderRepo struct {
	conn  *pgx.Conn
	cache sync.Map
}

func (repo *OrderRepo) InitRepo(dburl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	repo.conn, err = pgx.Connect(ctx, dburl)
	if err != nil {
		log.Printf("Couldn't connect to db by url %s: %s\n", dburl, err)
		return err
	}
	return nil

}

func (repo *OrderRepo) Store(ord models.Order) error {
	repo.saveToDB(ord)
	return nil
}
func (repo *OrderRepo) GetAllRows() pgx.Rows {
	rows, err := repo.conn.Query(context.Background(), `SELECT order_uid, track_number, entry, locale,
       internal_signature, customer_id,delivery_service,shardkey, sm_id,
       date_created, oof_shard FROM "order"`)
	if err != nil {
		log.Printf("Shit happened: %s", err)
	}
	return rows
}

func (repo *OrderRepo) saveToDB(order models.Order) error {
	tx, err := repo.conn.Begin(context.Background())
	if err != nil {
		log.Printf("Unable to begin transaction: %v\n", err)
		return err
	}
	defer tx.Rollback(context.Background()) // Автоматически откатит если если не будет коммита

	_, err = tx.Exec(context.Background(), insertIntoOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard)
	if err != nil {
		log.Printf("Insert order failed: %v\n", err)
		return err
	}
	delivery := &order.Delivery
	_, err = tx.Exec(context.Background(), insertIntoDelivery,
		delivery.OrderUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	if err != nil {
		log.Printf("Insert delivery failed: %v\n", err)
		return err
	}
	payment := &order.Payment
	_, err = tx.Exec(context.Background(), insertIntoPayment,
		payment.OrderUID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)

	if err != nil {
		log.Printf("Insert payment failed: %v\n", err)
		return err
	}

	for i := 0; i < len(order.Items); i++ {
		item := &order.Items[i]
		_, err = tx.Exec(context.Background(), insertIntoItem,
			item.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			log.Printf("Insert item failed: %v", err)
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Commit failed: %v\n", err)
		return err
	}
	return nil
}

// Get

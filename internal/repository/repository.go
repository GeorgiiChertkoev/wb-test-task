package repository

// пакет отвечающий за получение и сохранение информации

import (
	"context"
	"log"

	"orders/internal/cache"
	"orders/pkg/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const cache_capacity = 500
const max_db_connections = 50

type OrderRepo struct {
	pool  *pgxpool.Pool
	cache cache.Cache
}

func (repo *OrderRepo) InitRepo(dburl string) error {
	config, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		log.Printf("Couldn't parse config: %v", err)
		return err
	}
	config.MaxConns = max_db_connections

	repo.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("Couldn't make pool: %v", err)
		return err
	}

	repo.cache = *cache.MakeCache(cache_capacity)
	// fill cache
	orders, err := repo.GetOrders(cache_capacity)
	if err != nil {
		log.Printf("failed to init cache: %v", err)
		return err
	}
	for i := 0; i < len(orders); i++ {
		repo.cache.Add(&orders[i])
	}

	return nil

}

func (repo *OrderRepo) Store(ord *models.Order) error {
	err := repo.saveToDB(ord)
	if err != nil {
		log.Printf("Failed to save to db: %v", err)
		return err
	}
	repo.cache.Add(ord)
	return nil
}

func (repo *OrderRepo) Find(order_uid string) (order models.Order, found bool, err error) {
	// check cache
	cacher_order, found := repo.cache.Get(order_uid)
	if found {
		return *cacher_order, true, nil
	}

	return repo.getFromDB(order_uid)
}

func (repo *OrderRepo) GetOrders(quantity int) ([]models.Order, error) {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Failed to get conn from pool: %v", err)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		`SELECT order_uid FROM "order" LIMIT $1`, quantity)
	if err != nil {
		log.Printf("Failed to fetch ids to to get orders: %v", err)
		return nil, err
	}

	var uids []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			log.Printf("Error while scanning ids: %v", err)
			return nil, err
		}
		uids = append(uids, uid)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
		return nil, err
	}
	rows.Close()

	var orders []models.Order
	for i := 0; i < len(uids); i++ {
		order, found, err := repo.Find(uids[i])
		if !found {
			log.Printf("order %v not found\n", uids[i])
		} else if err != nil {
			log.Printf("Error while searching by id: %v", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (repo *OrderRepo) Close() {
	repo.pool.Close()
}

func (repo *OrderRepo) saveToDB(order *models.Order) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Failed to get conn from pool: %v", err)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
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
		order.OrderUID,
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
		order.OrderUID,
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
			order.OrderUID,
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

func (repo *OrderRepo) getFromDB(order_uid string) (order models.Order, found bool, err error) {
	ctx := context.Background()
	found = true

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Failed to get conn from pool: %v", err)
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `SELECT * FROM "order" WHERE order_uid = $1`, order_uid)
	scanOrder(&order, row)

	if err != nil {
		if err == pgx.ErrNoRows {
			found = false
			err = nil
			log.Printf("didn't find order with id = %v\n", order_uid)
			return
		}
		log.Printf("Querry by id failed: %v", err)
		return
	}
	row = tx.QueryRow(ctx, "SELECT * FROM delivery WHERE order_uid = $1", order_uid)
	scanDelivery(&order, row)

	if err != nil && err != pgx.ErrNoRows { // order without delivery is possible
		log.Printf("Querry by id failed at delivery: %v", err)
		return
	}
	row = tx.QueryRow(ctx, "SELECT * FROM payment WHERE order_uid = $1", order_uid)
	scanPayment(&order, row)

	if err != nil && err != pgx.ErrNoRows { // order without payment is possible
		log.Printf("Querry by id failed at payment: %v", err)
		return
	}

	rows, err := tx.Query(ctx, "SELECT * FROM item WHERE order_uid = $1", order_uid)
	if err != nil {
		log.Printf("Querry by id failed at items: %v", err)
		return
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Item])
	if err != nil {
		log.Printf("Error while collecting items: %v", err)
		return
	}
	order.Items = items

	return

}

func scanOrder(order *models.Order, row pgx.Row) error {
	return row.Scan(&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
}

func scanDelivery(order *models.Order, row pgx.Row) error {
	return row.Scan(
		&order.Delivery.OrderUID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
}

func scanPayment(order *models.Order, row pgx.Row) error {
	return row.Scan(
		&order.Payment.OrderUID,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
}

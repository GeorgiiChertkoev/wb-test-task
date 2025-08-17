package models

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Order struct {
	OrderUID          string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery" db:"-"`
	Payment           Payment   `json:"payment" db:"-"`
	Items             []Item    `json:"items" db:"-"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	Shardkey          string    `json:"shardkey" db:"shardkey"`
	SmID              int       `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

func MakeRandomOrder() *Order {
	orderUID := gofakeit.UUID()
	trackNumber := gofakeit.Regex(`[A-Z]{2}\d{9}[A-Z]{2}`)

	itemCount := gofakeit.Number(1, 3)
	items := make([]Item, itemCount)
	totalPrice := 0

	for i := 0; i < itemCount; i++ {
		price := gofakeit.Number(100, 10000)
		sale := gofakeit.Number(0, 50)
		totalItemPrice := price * (100 - sale) / 100

		items[i] = Item{
			ID:          i,
			OrderUID:    orderUID,
			ChrtID:      gofakeit.Int64(),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        sale,
			Size:        gofakeit.RandomString([]string{"S", "M", "L", "XL", "XXL"}),
			TotalPrice:  totalItemPrice,
			NmID:        gofakeit.Int64(),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.Number(100, 400),
		}
		totalPrice += totalItemPrice
	}

	deliveryCost := gofakeit.Number(100, 1000)
	goodsTotal := totalPrice

	return &Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       gofakeit.RandomString([]string{"WBIL", "WBIL2", "WBIL3"}),
		Delivery: Delivery{
			OrderUID: orderUID,
			Name:     gofakeit.Name(),
			Phone:    gofakeit.Phone(),
			Zip:      gofakeit.Zip(),
			City:     gofakeit.City(),
			Address:  gofakeit.Street(),
			Region:   gofakeit.State(),
			Email:    gofakeit.Email(),
		},
		Payment: Payment{
			Transaction:  gofakeit.UUID(),
			RequestID:    gofakeit.UUID(),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.RandomString([]string{"wbpay", "applepay", "googlepay"}),
			Amount:       goodsTotal + deliveryCost,
			PaymentDt:    gofakeit.Date().Unix(),
			Bank:         gofakeit.Company(),
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    gofakeit.Number(0, 100),
		},
		Items:             items,
		Locale:            gofakeit.Language(),
		InternalSignature: gofakeit.UUID(),
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   gofakeit.RandomString([]string{"meest", "nova_poshta", "ukrposhta"}),
		Shardkey:          gofakeit.DigitN(1),
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          gofakeit.DigitN(1),
	}
}

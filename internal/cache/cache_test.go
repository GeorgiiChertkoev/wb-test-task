package cache

import (
	"math/rand"
	"orders/internal/models"
	"testing"
	"time"
)

func TestCache_BaseFunctionality(t *testing.T) {
	cache := MakeCache(5)

	saved_order := makeRandomOrder()
	cache.Add(saved_order)

	from_cache, found := cache.Get(saved_order.OrderUID)
	if !found {
		t.Errorf("Cache didn't find added order")
	}
	if from_cache != saved_order {
		t.Errorf("Got different order. Got: %s, wanted: %s", from_cache.OrderUID, saved_order.Entry)
	}
}

func TestCache_SearchInEmptyCache(t *testing.T) {
	cache := MakeCache(2)
	not_existing_id := "13s"
	order, found := cache.Get(not_existing_id)
	if found {
		t.Errorf("Found order in empty cache. got %s, searched for %s", order.OrderUID, not_existing_id)
	}
}

func TestCache_CacheEviction(t *testing.T) {
	cache := MakeCache(2)

	order1 := makeRandomOrder()
	order2 := makeRandomOrder()
	order3 := makeRandomOrder()

	cache.Add(order1)
	cache.Add(order2)
	cache.Add(order3) // order1 вытеснется

	_, found := cache.Get(order1.OrderUID)
	if found {
		t.Error("Order1 should be evicted")
	}

	if _, found := cache.Get(order2.OrderUID); !found {
		t.Error("Order2 should still be in cache")
	}
	if _, found := cache.Get(order3.OrderUID); !found {
		t.Error("Order3 should still be in cache")
	}
}

func TestCache_LRUOrderCheck(t *testing.T) {
	cache := MakeCache(2)

	order1 := makeRandomOrder()
	order2 := makeRandomOrder()
	order3 := makeRandomOrder()

	cache.Add(order1)
	cache.Add(order2)

	cache.Get(order1.OrderUID)

	cache.Add(order3) // order2 вытесняется т.к. дольше всего не использовался

	if _, found := cache.Get(order2.OrderUID); found {
		t.Error("Order2 should be evicted")
	}

	if _, found := cache.Get(order1.OrderUID); !found {
		t.Error("Order1 should still be in cache")
	}
	if _, found := cache.Get(order3.OrderUID); !found {
		t.Error("Order3 should still be in cache")
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func makeRandomOrder() *models.Order {
	uid := randomString(19)

	order := &models.Order{
		OrderUID:          uid,
		TrackNumber:       randomString(10),
		Entry:             randomString(4),
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
	return order
}

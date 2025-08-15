package cache

import (
	"orders/pkg/models"
	"testing"
)

func TestCache_BaseFunctionality(t *testing.T) {
	cache := MakeCache(5)

	saved_order := models.MakeRandomOrder()
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

	order1 := models.MakeRandomOrder()
	order2 := models.MakeRandomOrder()
	order3 := models.MakeRandomOrder()

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

	order1 := models.MakeRandomOrder()
	order2 := models.MakeRandomOrder()
	order3 := models.MakeRandomOrder()

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

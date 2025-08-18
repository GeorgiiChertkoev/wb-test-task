package cache

import (
	"container/list"
	"sync"

	"orders/pkg/models"
)

type Cache struct {
	/*
		Кэш с LRU политикой вытеснения
		Элементы лежат в List
		Map для поиска за O(1)
		Когда элемент используется он перемещается в начало списка
		Когда нужно вытеснить элемент, вытесняется последний в списке
	*/

	capacity  int
	cacheMap  map[string]*list.Element
	cacheList *list.List
	mutex     sync.Mutex
}

func MakeCache(capacity int) *Cache {
	return &Cache{
		capacity:  capacity,
		cacheMap:  make(map[string]*list.Element),
		cacheList: list.New(),
	}
}

func (c *Cache) Get(order_uid string) (order *models.Order, found bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, found := c.cacheMap[order_uid]
	if !found {
		return nil, false
	}

	c.cacheList.MoveToFront(element)
	return element.Value.(*models.Order), true
}

func (c *Cache) Add(order *models.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, found := c.cacheMap[order.OrderUID]; found {
		element.Value = order
		c.cacheList.MoveToFront(element)
		return
	}

	c.cacheList.PushFront(order)
	c.cacheMap[order.OrderUID] = c.cacheList.Front()
	if c.cacheList.Len() > c.capacity {
		delete(c.cacheMap, c.cacheList.Back().Value.(*models.Order).OrderUID)
		c.cacheList.Remove(c.cacheList.Back())
	}

}

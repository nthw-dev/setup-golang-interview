// Package lru — เฉลยข้อ 11: LRU cache ด้วย map + doubly linked list + generics
package lru

import "container/list"

type entry[K comparable, V any] struct {
	key   K
	value V
}

// Cache — Get/Put ต้องเป็น O(1)
//   - map: หา element จาก key ได้ O(1)
//   - list: หน้าสุด = ใช้ล่าสุด, ท้ายสุด = ใช้นานที่สุด (ตัวที่จะถูกเตะออก)
type Cache[K comparable, V any] struct {
	capacity int
	ll       *list.List
	items    map[K]*list.Element
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		capacity: capacity,
		ll:       list.New(),
		items:    make(map[K]*list.Element, capacity),
	}
}

// Get คืนค่าและย้าย key นั้นไปเป็น "ใช้ล่าสุด"
func (c *Cache[K, V]) Get(key K) (V, bool) {
	el, ok := c.items[key]
	if !ok {
		var zero V // zero value ของ type V ใด ๆ
		return zero, false
	}
	c.ll.MoveToFront(el)
	return el.Value.(*entry[K, V]).value, true
}

// Put เพิ่ม/อัปเดตค่า ถ้าเกิน capacity ให้ลบตัวที่ไม่ได้ใช้นานที่สุดออก
func (c *Cache[K, V]) Put(key K, value V) {
	if c.capacity <= 0 {
		return
	}
	if el, ok := c.items[key]; ok {
		el.Value.(*entry[K, V]).value = value
		c.ll.MoveToFront(el)
		return
	}
	c.items[key] = c.ll.PushFront(&entry[K, V]{key, value})
	if c.ll.Len() > c.capacity {
		oldest := c.ll.Back()
		c.ll.Remove(oldest)
		delete(c.items, oldest.Value.(*entry[K, V]).key) // ต้องลบจาก map ด้วย ไม่งั้น memory leak
	}
}

func (c *Cache[K, V]) Len() int { return c.ll.Len() }

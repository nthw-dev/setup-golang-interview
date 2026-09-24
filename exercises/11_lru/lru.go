// Package lru — ข้อ 11: LRU cache (Get/Put ต้องเป็น O(1)) + generics
// hint: map + container/list
package lru

type Cache[K comparable, V any] struct {
	// TODO
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	// TODO
	return &Cache[K, V]{}
}

// Get คืนค่าและทำให้ key นั้นเป็น "ใช้ล่าสุด"
func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO
	var zero V
	return zero, false
}

// Put เพิ่ม/อัปเดตค่า (ถือเป็นการใช้ล่าสุด)
// ถ้าเกิน capacity ให้ลบตัวที่ไม่ได้ใช้นานที่สุดออก
func (c *Cache[K, V]) Put(key K, value V) {
	// TODO
}

func (c *Cache[K, V]) Len() int {
	// TODO
	return 0
}

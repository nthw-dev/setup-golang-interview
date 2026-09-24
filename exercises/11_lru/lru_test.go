package lru

import "testing"

func TestLRUBasic(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v; want 1, true", v, ok)
	}
	c.Put("c", 3) // b ไม่ได้ใช้นานสุด → ถูกเตะออก
	if _, ok := c.Get("b"); ok {
		t.Errorf("b should be evicted")
	}
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("a should still exist, got %v, %v", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Errorf("c should exist, got %v, %v", v, ok)
	}
	if c.Len() != 2 {
		t.Errorf("Len = %d; want 2", c.Len())
	}
}

func TestLRUUpdate(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 10) // อัปเดต a → a กลายเป็นใช้ล่าสุด
	c.Put("c", 3)  // b ต้องถูกเตะ ไม่ใช่ a
	if v, ok := c.Get("a"); !ok || v != 10 {
		t.Errorf("Get(a) = %v, %v; want 10, true", v, ok)
	}
	if _, ok := c.Get("b"); ok {
		t.Errorf("b should be evicted")
	}
	if c.Len() != 2 {
		t.Errorf("Len = %d; want 2", c.Len())
	}
}

func TestLRUMissAndGenerics(t *testing.T) {
	c := New[int, string](1)
	if v, ok := c.Get(1); ok || v != "" {
		t.Errorf("Get on empty = %q, %v; want \"\", false", v, ok)
	}
	c.Put(1, "one")
	c.Put(2, "two")
	if _, ok := c.Get(1); ok {
		t.Errorf("1 should be evicted with capacity 1")
	}
	if v, _ := c.Get(2); v != "two" {
		t.Errorf("Get(2) = %q; want two", v)
	}
}

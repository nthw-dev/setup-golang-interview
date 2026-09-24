// Package counter — เฉลยข้อ 8: race condition, sync.Mutex
package counter

import "sync"

// Counter นับจำนวนครั้งต่อ key และต้องปลอดภัยเมื่อถูกเรียกจากหลาย goroutine
// map ของ Go ไม่ thread-safe: เขียนพร้อมกัน → "fatal error: concurrent map writes"
type Counter struct {
	mu sync.Mutex // ไม่ต้อง init, zero value ใช้งานได้เลย
	m  map[string]int
}

func NewCounter() *Counter {
	return &Counter{m: make(map[string]int)}
}

func (c *Counter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key]++
}

func (c *Counter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[key]
}

// Snapshot คืน "สำเนา" ของข้อมูล — ถ้าคืน c.m ตรง ๆ คนนอกจะแก้ map ได้โดยไม่ผ่าน lock
func (c *Counter) Snapshot() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int, len(c.m))
	for k, v := range c.m {
		out[k] = v
	}
	return out
}

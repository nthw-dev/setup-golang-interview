// Package counter — ข้อ 8: race condition, sync.Mutex
// BUG: ถูกเรียกจากหลาย goroutine แล้วพัง (fatal error: concurrent map writes)
// ลองรัน: go test -race ./exercises/08_mutex
package counter

type Counter struct {
	m map[string]int
}

func NewCounter() *Counter {
	return &Counter{m: make(map[string]int)}
}

func (c *Counter) Inc(key string) {
	c.m[key]++
}

func (c *Counter) Value(key string) int {
	return c.m[key]
}

// Snapshot คืนข้อมูลทั้งหมด — คนเรียกแก้ผลลัพธ์แล้วต้องไม่กระทบ Counter
func (c *Counter) Snapshot() map[string]int {
	return c.m
}

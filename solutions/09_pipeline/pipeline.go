// Package pipeline — เฉลยข้อ 9: pipeline, fan-in, context cancellation, goroutine leak
package pipeline

import (
	"context"
	"sync"
)

// Generate ส่งตัวเลขออกทาง channel ทีละตัว แล้ว close เมื่อหมด
// ต้องหยุด (และ close) ทันทีเมื่อ ctx ถูก cancel — ไม่งั้น goroutine จะค้างที่ `out <- n` ตลอดไป (leak)
func Generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Square อ่านจาก in แล้วส่ง n*n ออกไป
func Square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Merge รวมหลาย channel เป็น channel เดียว (fan-in)
// close out เมื่อ input ทุกตัวถูก close แล้ว
func Merge(ctx context.Context, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range c {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	// goroutine แยกคอย close — ถ้า wg.Wait() ใน Merge ตรง ๆ จะ deadlock เพราะยังไม่มีใครอ่าน out
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

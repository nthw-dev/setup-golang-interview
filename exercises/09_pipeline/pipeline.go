// Package pipeline — ข้อ 9: pipeline, fan-in, context cancellation, goroutine leak
// BUG: Generate และ Square ใช้งานได้ แต่ถ้าผู้อ่านเลิกอ่านกลางทาง goroutine จะค้าง (leak)
package pipeline

import "context"

// Generate ส่งตัวเลขออกทาง channel ทีละตัว แล้ว close เมื่อหมด
// ต้องหยุดทันทีเมื่อ ctx ถูก cancel
func Generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Square อ่านจาก in แล้วส่ง n*n ออกไป, ต้องหยุดเมื่อ ctx ถูก cancel
func Square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Merge รวมหลาย channel เป็น channel เดียว (fan-in)
// close out เมื่อ input ทุกตัวถูก close แล้ว, ต้องหยุดเมื่อ ctx ถูก cancel
func Merge(ctx context.Context, cs ...<-chan int) <-chan int {
	// TODO
	out := make(chan int)
	close(out)
	return out
}

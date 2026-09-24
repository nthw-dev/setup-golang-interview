// Package workerpool — ข้อ 7: goroutine, channel, sync.WaitGroup
package workerpool

// ParallelMap เรียก fn กับทุกตัวใน nums โดยใช้ goroutine "ไม่เกิน" workers ตัวพร้อมกัน
// ผลลัพธ์ต้องเรียงตามลำดับเดิมของ nums
// ทดสอบ race ด้วย: go test -race ./exercises/07_workerpool
func ParallelMap(nums []int, workers int, fn func(int) int) []int {
	// TODO: ตอนนี้ทำงานทีละตัว (ช้า) — ทำให้เป็น worker pool
	out := make([]int, len(nums))
	for i, n := range nums {
		out[i] = fn(n)
	}
	return out
}

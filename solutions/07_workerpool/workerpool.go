// Package workerpool — เฉลยข้อ 7: goroutine, channel, sync.WaitGroup, worker pool
package workerpool

import "sync"

// ParallelMap เรียก fn กับทุกตัวใน nums โดยใช้ goroutine ไม่เกิน workers ตัวพร้อมกัน
// ผลลัพธ์ต้องเรียงตามลำดับเดิมของ nums
func ParallelMap(nums []int, workers int, fn func(int) int) []int {
	if workers < 1 {
		workers = 1
	}
	results := make([]int, len(nums)) // จอง slice ไว้ก่อน ให้แต่ละ worker เขียนคนละ index → ไม่มี data race
	jobs := make(chan int)            // ส่ง "index" ของงาน

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs { // วนจนกว่า jobs จะถูก close
				results[i] = fn(nums[i])
			}
		}()
	}

	for i := range nums {
		jobs <- i
	}
	close(jobs) // บอก worker ว่าไม่มีงานแล้ว → range ข้างบนจบ
	wg.Wait()   // รอทุก worker เสร็จก่อนคืนผล
	return results
}

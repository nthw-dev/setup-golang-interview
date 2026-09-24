package workerpool

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelMapResult(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := ParallelMap(nums, 3, func(n int) int { return n * n })
	want := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParallelMap = %v; want %v", got, want)
	}
	if got := ParallelMap(nil, 3, func(n int) int { return n }); len(got) != 0 {
		t.Errorf("ParallelMap(nil) = %v; want empty", got)
	}
}

func TestParallelMapConcurrency(t *testing.T) {
	var running, maxRunning int64
	fn := func(n int) int {
		cur := atomic.AddInt64(&running, 1)
		for {
			m := atomic.LoadInt64(&maxRunning)
			if cur <= m || atomic.CompareAndSwapInt64(&maxRunning, m, cur) {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt64(&running, -1)
		return n
	}

	start := time.Now()
	ParallelMap(make([]int, 8), 4, fn)
	elapsed := time.Since(start)

	if maxRunning > 4 {
		t.Errorf("max concurrent = %d; must not exceed 4 workers", maxRunning)
	}
	if maxRunning < 2 {
		t.Errorf("max concurrent = %d; work is not running in parallel", maxRunning)
	}
	if elapsed > 300*time.Millisecond {
		t.Errorf("took %v; 8 jobs x 50ms with 4 workers should take ~100ms", elapsed)
	}
}

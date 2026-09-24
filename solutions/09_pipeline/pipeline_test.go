package pipeline

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func sum(c <-chan int) int {
	total := 0
	for n := range c {
		total += n
	}
	return total
}

func TestPipeline(t *testing.T) {
	ctx := context.Background()
	if got := sum(Square(ctx, Generate(ctx, 1, 2, 3, 4, 5))); got != 55 {
		t.Fatalf("sum of squares = %d; want 55", got)
	}
}

func TestMergeFanOut(t *testing.T) {
	ctx := context.Background()
	src := Generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	// fan-out: 3 worker อ่านจาก source เดียวกัน, แล้ว fan-in ด้วย Merge
	got := sum(Merge(ctx, Square(ctx, src), Square(ctx, src), Square(ctx, src)))
	if got != 385 {
		t.Fatalf("fan-out/fan-in sum = %d; want 385", got)
	}
	if got := sum(Merge(ctx)); got != 0 {
		t.Errorf("Merge() with no input should close immediately")
	}
}

func TestNoGoroutineLeak(t *testing.T) {
	base := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	out := Merge(ctx, Square(ctx, Generate(ctx, nums...)))
	if first := <-out; first != 1 {
		t.Fatalf("first value = %d; want 1", first)
	}
	cancel() // ผู้บริโภคเลิกอ่านกลางทาง

	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > base {
		if time.Now().After(deadline) {
			t.Fatalf("goroutine leak: %d goroutines still running (base %d) — handle ctx.Done()",
				runtime.NumGoroutine(), base)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

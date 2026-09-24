package memory

import (
	"reflect"
	"testing"
)

func TestSquares(t *testing.T) {
	if got := Squares(5); !reflect.DeepEqual(got, []int{0, 1, 4, 9, 16}) {
		t.Fatalf("Squares(5) = %v", got)
	}
	allocs := testing.AllocsPerRun(100, func() { _ = Squares(1000) })
	if allocs > 1 {
		t.Errorf("Squares(1000) allocs = %v; want 1 (hint: make with capacity)", allocs)
	}
}

func TestJoinInts(t *testing.T) {
	tests := []struct {
		nums []int
		sep  string
		want string
	}{
		{nil, ",", ""},
		{[]int{7}, ",", "7"},
		{[]int{1, -2, 300}, ", ", "1, -2, 300"},
	}
	for _, tt := range tests {
		if got := JoinInts(tt.nums, tt.sep); got != tt.want {
			t.Errorf("JoinInts(%v, %q) = %q; want %q", tt.nums, tt.sep, got, tt.want)
		}
	}

	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i
	}
	allocs := testing.AllocsPerRun(100, func() { _ = JoinInts(nums, ",") })
	if allocs > 3 {
		t.Errorf("JoinInts(100 nums) allocs = %v; want <= 3 (hint: strings.Builder + Grow)", allocs)
	}
}

func TestHead(t *testing.T) {
	big := make([]byte, 1<<20)
	for i := range big {
		big[i] = byte(i)
	}
	h := Head(big, 10)
	if len(h) != 10 || h[9] != 9 {
		t.Fatalf("Head = %v", h)
	}
	big[0] = 255
	if h[0] != 0 {
		t.Errorf("Head shares memory with input (hint: copy)")
	}
	if cap(h) > 64 {
		t.Errorf("cap(Head) = %d; result still holds the big backing array", cap(h))
	}
	if got := Head([]byte("hi"), 10); string(got) != "hi" {
		t.Errorf("Head(hi, 10) = %q; want hi", got)
	}
}

func BenchmarkJoinInts(b *testing.B) {
	nums := make([]int, 1000)
	for b.Loop() {
		JoinInts(nums, ",")
	}
}

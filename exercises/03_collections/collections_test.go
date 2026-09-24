package collections

import (
	"reflect"
	"testing"
)

func TestWordFrequency(t *testing.T) {
	got := WordFrequency("Go is fun. go, GO! Is it?")
	want := map[string]int{"go": 3, "is": 2, "fun": 1, "it": 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("WordFrequency = %v; want %v", got, want)
	}
	if got := WordFrequency(""); len(got) != 0 {
		t.Errorf("WordFrequency(\"\") = %v; want empty", got)
	}
}

func TestTopK(t *testing.T) {
	freq := map[string]int{"go": 3, "is": 2, "fun": 2, "it": 1}
	tests := []struct {
		k    int
		want []string
	}{
		{1, []string{"go"}},
		{3, []string{"go", "fun", "is"}},
		{10, []string{"go", "fun", "is", "it"}},
	}
	for _, tt := range tests {
		if got := TopK(freq, tt.k); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("TopK(k=%d) = %v; want %v", tt.k, got, tt.want)
		}
	}
	if got := TopK(freq, 0); len(got) != 0 {
		t.Errorf("TopK(k=0) = %v; want empty", got)
	}
}

func TestUnique(t *testing.T) {
	got := Unique([]int{3, 1, 3, 2, 1, 5})
	want := []int{3, 1, 2, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unique = %v; want %v", got, want)
	}
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Chunk = %v; want %v", got, want)
	}
	if got := Chunk([]int{1, 2}, 0); got != nil {
		t.Errorf("Chunk(size=0) = %v; want nil", got)
	}

	// append ลง chunk แรก ต้องไม่ไปทับ chunk ที่สอง
	nums := []int{1, 2, 3, 4}
	chunks := Chunk(nums, 2)
	_ = append(chunks[0], 99)
	if chunks[1][0] != 3 {
		t.Errorf("append to chunk[0] overwrote chunk[1]: %v (hint: full slice expression)", chunks)
	}
}

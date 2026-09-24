package basics

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestSumStrings(t *testing.T) {
	got, err := SumStrings([]string{"1", " 2 ", "39"})
	if err != nil || got != 42 {
		t.Fatalf("SumStrings = %d, %v; want 42, nil", got, err)
	}

	got, err = SumStrings(nil)
	if err != nil || got != 0 {
		t.Fatalf("SumStrings(nil) = %d, %v; want 0, nil", got, err)
	}

	_, err = SumStrings([]string{"1", "2", "abc"})
	if err == nil {
		t.Fatal("expected error for \"abc\"")
	}
	if !strings.Contains(err.Error(), "index 2") {
		t.Errorf("error %q should contain \"index 2\"", err)
	}
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Errorf("error should wrap *strconv.NumError (use %%w)")
	}
}

func TestAverage(t *testing.T) {
	got, err := Average([]int{1, 2})
	if err != nil || got != 1.5 {
		t.Fatalf("Average([1 2]) = %v, %v; want 1.5, nil", got, err)
	}
	if _, err := Average(nil); !errors.Is(err, ErrEmpty) {
		t.Fatalf("Average(nil) error = %v; want ErrEmpty", err)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"hello", "olleh"},
		{"กขค", "คขก"},
		{"Go🚀", "🚀oG"},
	}
	for _, tt := range tests {
		if got := Reverse(tt.in); got != tt.want {
			t.Errorf("Reverse(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestCharCount(t *testing.T) {
	b, r := CharCount("สวัสดี")
	if b != 18 || r != 6 {
		t.Errorf("CharCount(สวัสดี) = %d, %d; want 18, 6", b, r)
	}
}

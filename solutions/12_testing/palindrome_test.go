package palindrome

import (
	"fmt"
	"testing"
)

// Table-driven test + subtest (t.Run) — สไตล์มาตรฐานของ Go
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"empty", "", true},
		{"single char", "a", true},
		{"simple", "racecar", true},
		{"not palindrome", "golang", false},
		{"mixed case", "Level", true},
		{"with spaces and punctuation", "A man, a plan, a canal: Panama", true},
		{"digits", "12321", true},
		{"multi-byte unicode", "été", true},
		{"thai", "กขก", true},
		{"only punctuation", "!!", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsPalindrome(tt.in); got != tt.want {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.in, got, tt.want)
			}
		})
	}
}

// Benchmark — รันด้วย: go test -bench=. -benchmem ./solutions/12_testing
func BenchmarkIsPalindrome(b *testing.B) {
	s := "A man, a plan, a canal: Panama"
	for b.Loop() { // Go 1.24+ (เวอร์ชันเก่าใช้ for i := 0; i < b.N; i++)
		IsPalindrome(s)
	}
}

// Example — เป็นทั้ง test (เทียบกับ // Output:) และเอกสารใน go doc
func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("Never odd or even"))
	// Output: true
}

// Fuzz test — รันด้วย: go test -fuzz=FuzzIsPalindrome ./solutions/12_testing
// property: s + reverse(s) ต้องเป็น palindrome เสมอ
func FuzzIsPalindrome(f *testing.F) {
	f.Add("hello")
	f.Add("สวัสดี")
	f.Fuzz(func(t *testing.T, s string) {
		r := []rune(s)
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		if !IsPalindrome(s + string(r)) {
			t.Errorf("IsPalindrome(%q) = false", s+string(r))
		}
	})
}

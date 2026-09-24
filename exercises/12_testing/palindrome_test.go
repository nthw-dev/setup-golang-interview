package palindrome

import "testing"

// TODO:
//  1. เปลี่ยนเป็น table-driven test ที่ใช้ t.Run และมีเคสครอบคลุม (อย่างน้อย 8 เคส)
//  2. เขียน BenchmarkIsPalindrome
//  3. เขียน ExampleIsPalindrome (มี // Output:)
//  4. (โบนัส) เขียน FuzzIsPalindrome
func TestIsPalindrome(t *testing.T) {
	if !IsPalindrome("racecar") {
		t.Error("racecar should be a palindrome")
	}
}

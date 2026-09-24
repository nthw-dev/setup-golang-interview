// Package palindrome — ข้อ 12: เขียน test หาบั๊ก แล้วแก้
package palindrome

import "strings"

// IsPalindrome คืน true ถ้าอ่านจากหน้าไปหลังและหลังไปหน้าได้เหมือนกัน
// โดยไม่สนตัวพิมพ์เล็ก/ใหญ่ และข้ามทุกอย่างที่ไม่ใช่ตัวอักษรหรือตัวเลข (รองรับ Unicode)
//
// ตัวอย่าง: "A man, a plan, a canal: Panama" → true, "été" → true
//
// ⚠️ ฟังก์ชันนี้มีบั๊กอย่างน้อย 2 จุด — เขียน test ให้เจอก่อน แล้วค่อยแก้
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

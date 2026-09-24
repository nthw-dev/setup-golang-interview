// Package palindrome — เฉลยข้อ 12: เขียน test หาบั๊ก แล้วแก้
package palindrome

import "unicode"

// IsPalindrome คืน true ถ้าอ่านจากหน้าไปหลังและหลังไปหน้าได้เหมือนกัน
// โดยไม่สนตัวพิมพ์เล็ก/ใหญ่ และข้ามทุกอย่างที่ไม่ใช่ตัวอักษรหรือตัวเลข (รองรับ Unicode)
//
// บั๊กของเวอร์ชันเดิม:
//  1. เทียบทีละ byte → ตัวอักษร UTF-8 หลาย byte (é, ก) พัง
//  2. ไม่ข้ามช่องว่าง/เครื่องหมายวรรคตอน
func IsPalindrome(s string) bool {
	rs := make([]rune, 0, len(s))
	for _, r := range s { // range บน string ได้ทีละ rune
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			rs = append(rs, unicode.ToLower(r))
		}
	}
	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		if rs[i] != rs[j] {
			return false
		}
	}
	return true
}

// Package basics — ข้อ 1: ตัวแปร, zero value, type conversion, string/rune
// ดูโจทย์เต็มใน PROBLEMS.md แล้วรัน: go test ./exercises/01_basics
package basics

import "errors"

// ErrEmpty ใช้เมื่อ input ว่าง
var ErrEmpty = errors.New("empty input")

// SumStrings แปลงทุกตัวเป็น int แล้วรวมกัน
// - ตัด space หัวท้ายก่อนแปลง
// - ถ้าแปลงไม่ได้ ให้คืน error ที่มีคำว่า "index <i>" และ wrap error เดิมไว้
func SumStrings(items []string) (int, error) {
	// TODO
	return 0, nil
}

// Average คืนค่าเฉลี่ยแบบ float64, ถ้า nums ว่างให้คืน ErrEmpty
func Average(nums []int) (float64, error) {
	// TODO
	return 0, nil
}

// Reverse กลับลำดับตัวอักษร (ต้องรองรับภาษาไทยและ emoji)
func Reverse(s string) string {
	// TODO
	return ""
}

// CharCount คืนจำนวน byte และจำนวนตัวอักษร (rune) ของ s
func CharCount(s string) (bytes int, runes int) {
	// TODO
	return 0, 0
}

// Package basics — เฉลยข้อ 1: ตัวแปร, zero value, type conversion, string/rune
package basics

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ErrEmpty ใช้เมื่อ input ว่าง (sentinel error)
var ErrEmpty = errors.New("empty input")

// SumStrings แปลงทุกตัวเป็น int แล้วรวมกัน
// - ตัด space หัวท้ายก่อนแปลง
// - ถ้าแปลงไม่ได้ ให้คืน error ที่มีคำว่า "index <i>" และ wrap error เดิมด้วย %w
func SumStrings(items []string) (int, error) {
	total := 0 // zero value ของ int คือ 0 อยู่แล้ว แต่เขียนให้ชัด
	for i, s := range items {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			// %w ทำให้ caller ใช้ errors.As หา *strconv.NumError ต่อได้
			return 0, fmt.Errorf("index %d: %w", i, err)
		}
		total += n
	}
	return total, nil
}

// Average คืนค่าเฉลี่ยแบบ float64 (ระวัง integer division!)
func Average(nums []int) (float64, error) {
	if len(nums) == 0 {
		return 0, ErrEmpty
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	// ต้องแปลงเป็น float64 ก่อนหาร ไม่งั้น 3/2 = 1
	return float64(sum) / float64(len(nums)), nil
}

// Reverse กลับลำดับตัวอักษร (ระดับ rune ไม่ใช่ byte)
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// CharCount คืนจำนวน byte และจำนวน rune ของ s
func CharCount(s string) (bytes int, runes int) {
	return len(s), utf8.RuneCountInString(s)
}

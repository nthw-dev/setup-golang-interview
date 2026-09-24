// Package collections — เฉลยข้อ 3: slice, map, sort
package collections

import (
	"sort"
	"strings"
	"unicode"
)

// WordFrequency นับความถี่ของคำ (ไม่สนตัวพิมพ์เล็ก/ใหญ่)
// "คำ" = กลุ่มตัวอักษรหรือตัวเลขที่ติดกัน อย่างอื่นถือเป็นตัวคั่น
func WordFrequency(text string) map[string]int {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	freq := make(map[string]int, len(words))
	for _, w := range words {
		freq[w]++ // key ที่ไม่มีอยู่ จะได้ zero value (0) มาก่อนบวก
	}
	return freq
}

// TopK คืน k คำที่พบบ่อยที่สุด เรียงจากมากไปน้อย
// ถ้าจำนวนเท่ากัน ให้เรียงตามตัวอักษร (a→z)
// หมายเหตุ: การวนลูป map ใน Go "ไม่มีลำดับ" จึงต้อง sort เอง
func TopK(freq map[string]int, k int) []string {
	if k <= 0 {
		return nil
	}
	words := make([]string, 0, len(freq))
	for w := range freq {
		words = append(words, w)
	}
	sort.Slice(words, func(i, j int) bool {
		if freq[words[i]] != freq[words[j]] {
			return freq[words[i]] > freq[words[j]]
		}
		return words[i] < words[j]
	})
	if k > len(words) {
		k = len(words)
	}
	return words[:k]
}

// Unique ลบตัวซ้ำโดยคงลำดับเดิม
func Unique(nums []int) []int {
	seen := make(map[int]struct{}, len(nums)) // struct{} ใช้ memory 0 byte
	out := make([]int, 0, len(nums))
	for _, n := range nums {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}

// Chunk แบ่ง slice เป็นก้อนละ size ตัว (ก้อนสุดท้ายอาจสั้นกว่า)
// size <= 0 ให้คืน nil
func Chunk(nums []int, size int) [][]int {
	if size <= 0 {
		return nil
	}
	chunks := make([][]int, 0, (len(nums)+size-1)/size)
	for i := 0; i < len(nums); i += size {
		end := min(i+size, len(nums))
		// full slice expression [i:end:end] จำกัด cap
		// ป้องกันไม่ให้ append ของ chunk หนึ่งไปทับข้อมูลของ chunk ถัดไป
		chunks = append(chunks, nums[i:end:end])
	}
	return chunks
}

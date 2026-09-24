// Package memory — ข้อ 13: preallocation, strings.Builder, slice ค้าง backing array
// โค้ดตอนนี้ "ถูก" แต่ไม่มีประสิทธิภาพ — ปรับให้ test เรื่อง allocation ผ่าน
package memory

import "strconv"

// Squares คืน [0, 1, 4, 9, ...] จำนวน n ตัว — ต้อง allocate แค่ 1 ครั้ง
func Squares(n int) []int {
	var out []int
	for i := 0; i < n; i++ {
		out = append(out, i*i)
	}
	return out
}

// JoinInts ต่อตัวเลขด้วย sep เช่น JoinInts([1 2 3], ",") = "1,2,3"
func JoinInts(nums []int, sep string) string {
	s := ""
	for i, n := range nums {
		if i > 0 {
			s += sep
		}
		s += strconv.Itoa(n)
	}
	return s
}

// Head คืน n byte แรกของ data (ถ้า n > len ให้คืนทั้งหมด)
// ผลลัพธ์ต้อง "ไม่แชร์" memory กับ data
func Head(data []byte, n int) []byte {
	if n > len(data) {
		n = len(data)
	}
	return data[:n]
}

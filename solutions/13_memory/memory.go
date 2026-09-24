// Package memory — เฉลยข้อ 13: preallocation, strings.Builder, slice ค้าง backing array
package memory

import (
	"strconv"
	"strings"
)

// Squares คืน [0, 1, 4, 9, ...] จำนวน n ตัว — ต้อง allocate แค่ 1 ครั้ง
// append ลง slice ที่ cap ไม่พอ → Go ต้องสร้าง array ใหม่ใหญ่ขึ้นแล้ว copy (หลายรอบ)
func Squares(n int) []int {
	out := make([]int, 0, n) // จอง cap ไว้ครั้งเดียว
	for i := range n {
		out = append(out, i*i)
	}
	return out
}

// JoinInts ต่อตัวเลขด้วย sep เช่น JoinInts([1 2 3], ",") = "1,2,3"
// ห้ามใช้ s += ... ในลูป เพราะ string immutable → สร้าง string ใหม่ทุกรอบ (O(n^2))
func JoinInts(nums []int, sep string) string {
	if len(nums) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(nums) * (3 + len(sep))) // ประมาณขนาดไว้ก่อน
	var buf [20]byte                   // พอสำหรับ int64 ทุกค่า, อยู่บน stack
	for i, n := range nums {
		if i > 0 {
			b.WriteString(sep)
		}
		b.Write(strconv.AppendInt(buf[:0], int64(n), 10))
	}
	return b.String()
}

// Head คืน n byte แรกของ data (ถ้า n > len ให้คืนทั้งหมด)
// ผลลัพธ์ต้อง "ไม่แชร์" backing array กับ data
//
// ถ้าคืน data[:n] ตรง ๆ: แม้จะใช้แค่ 10 byte แต่ GC เก็บ array ขนาด 100MB ไม่ได้
// เพราะยังมี slice ชี้อยู่ → memory leak แบบเงียบ ๆ
func Head(data []byte, n int) []byte {
	n = min(n, len(data))
	out := make([]byte, n)
	copy(out, data)
	return out
}

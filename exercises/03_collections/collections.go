// Package collections — ข้อ 3: slice, map, sort
package collections

// WordFrequency นับความถี่ของคำ (ไม่สนตัวพิมพ์เล็ก/ใหญ่)
// "คำ" = กลุ่มตัวอักษรหรือตัวเลขที่ติดกัน อย่างอื่นถือเป็นตัวคั่น
func WordFrequency(text string) map[string]int {
	// TODO
	return nil
}

// TopK คืน k คำที่พบบ่อยที่สุด เรียงจากมากไปน้อย
// ถ้าจำนวนเท่ากัน ให้เรียงตามตัวอักษร (a→z), k <= 0 คืน slice ว่าง
func TopK(freq map[string]int, k int) []string {
	// TODO
	return nil
}

// Unique ลบตัวซ้ำโดยคงลำดับเดิม
func Unique(nums []int) []int {
	// TODO
	return nil
}

// Chunk แบ่ง slice เป็นก้อนละ size ตัว (ก้อนสุดท้ายอาจสั้นกว่า), size <= 0 คืน nil
// เงื่อนไขพิเศษ: append ลง chunk หนึ่ง ต้องไม่ไปทับข้อมูลของ chunk ถัดไป
func Chunk(nums []int, size int) [][]int {
	// TODO
	return nil
}

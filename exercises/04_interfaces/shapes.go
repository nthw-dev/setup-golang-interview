// Package shapes — ข้อ 4: struct, method, interface, polymorphism, type switch
package shapes

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

type Circle struct {
	Radius float64
}

// TODO: เขียน method Area() และ Perimeter() ของ Rectangle และ Circle ให้ถูก (ใช้ math.Pi)
func (r Rectangle) Area() float64      { return 0 }
func (r Rectangle) Perimeter() float64 { return 0 }
func (c Circle) Area() float64         { return 0 }
func (c Circle) Perimeter() float64    { return 0 }

// TotalArea รับกี่ shape ก็ได้ (variadic) แล้วคืนผลรวมพื้นที่
func TotalArea(shapes ...Shape) float64 {
	// TODO
	return 0
}

// Largest คืน shape ที่พื้นที่มากที่สุด, ถ้า slice ว่างคืน (nil, false)
func Largest(shapes []Shape) (Shape, bool) {
	// TODO
	return nil, false
}

// Describe
//
//	Rectangle{2,3} -> "rectangle 2x3"
//	Circle{1.5}    -> "circle r=1.5"
//	อื่น ๆ          -> "unknown shape"
func Describe(s Shape) string {
	// TODO
	return ""
}

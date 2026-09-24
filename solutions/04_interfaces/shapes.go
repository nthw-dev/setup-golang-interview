// Package shapes — เฉลยข้อ 4: struct, method, interface, polymorphism, type switch
package shapes

import (
	"fmt"
	"math"
)

// Shape — type ไหนมี method ครบ 2 ตัวนี้ ถือว่า implement Shape อัตโนมัติ
// (implicit interface — ไม่ต้องเขียน implements)
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// compile-time check ว่า type implement interface จริง
var (
	_ Shape = Rectangle{}
	_ Shape = Circle{}
)

// TotalArea รับกี่ shape ก็ได้ (variadic) — ทดแทน method overloading ที่ Go ไม่มี
func TotalArea(shapes ...Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// Largest คืน shape ที่พื้นที่มากที่สุด, ถ้า slice ว่างคืน (nil, false)
func Largest(shapes []Shape) (Shape, bool) {
	if len(shapes) == 0 {
		return nil, false
	}
	best := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > best.Area() {
			best = s
		}
	}
	return best, true
}

// Describe ใช้ type switch แยกประเภท
//
//	Rectangle{2,3} -> "rectangle 2x3"
//	Circle{1.5}    -> "circle r=1.5"
//	อื่น ๆ          -> "unknown shape"
func Describe(s Shape) string {
	switch v := s.(type) {
	case Rectangle:
		return fmt.Sprintf("rectangle %gx%g", v.Width, v.Height)
	case Circle:
		return fmt.Sprintf("circle r=%g", v.Radius)
	default:
		return "unknown shape"
	}
}

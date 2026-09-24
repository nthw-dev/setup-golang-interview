package shapes

import (
	"math"
	"testing"
)

type square struct{ side float64 }

func (s square) Area() float64      { return s.side * s.side }
func (s square) Perimeter() float64 { return 4 * s.side }

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestAreaPerimeter(t *testing.T) {
	r := Rectangle{Width: 2, Height: 3}
	if r.Area() != 6 || r.Perimeter() != 10 {
		t.Errorf("Rectangle area=%v perimeter=%v; want 6, 10", r.Area(), r.Perimeter())
	}
	c := Circle{Radius: 1}
	if !almostEqual(c.Area(), math.Pi) || !almostEqual(c.Perimeter(), 2*math.Pi) {
		t.Errorf("Circle area=%v perimeter=%v", c.Area(), c.Perimeter())
	}
}

func TestTotalArea(t *testing.T) {
	got := TotalArea(Rectangle{2, 3}, Circle{1}, square{2})
	if !almostEqual(got, 6+math.Pi+4) {
		t.Errorf("TotalArea = %v; want %v", got, 10+math.Pi)
	}
	if TotalArea() != 0 {
		t.Errorf("TotalArea() should be 0")
	}
}

func TestLargest(t *testing.T) {
	got, ok := Largest([]Shape{Rectangle{1, 1}, Circle{2}, square{3}})
	if !ok || got != (Circle{2}) {
		t.Errorf("Largest = %v, %v; want Circle{2}, true", got, ok)
	}
	if got, ok := Largest(nil); ok || got != nil {
		t.Errorf("Largest(nil) = %v, %v; want nil, false", got, ok)
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		s    Shape
		want string
	}{
		{Rectangle{2, 3}, "rectangle 2x3"},
		{Circle{1.5}, "circle r=1.5"},
		{square{1}, "unknown shape"},
	}
	for _, tt := range tests {
		if got := Describe(tt.s); got != tt.want {
			t.Errorf("Describe(%v) = %q; want %q", tt.s, got, tt.want)
		}
	}
}

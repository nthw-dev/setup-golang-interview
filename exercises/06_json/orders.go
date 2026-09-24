// Package orders — ข้อ 6: encoding/json, struct tags, omitempty
package orders

import "errors"

var ErrMissingID = errors.New("order id is required")

// TODO: เพิ่ม struct tag ให้ตรงกับ JSON ในโจทย์
//
//	Item:  sku, qty, price
//	Order: id, customer_name, items, note (ไม่ต้องใส่ key ถ้า note ว่าง)
type Item struct {
	SKU   string
	Qty   int
	Price float64
}

type Order struct {
	ID       string
	Customer string
	Items    []Item
	Note     string
}

// Total = ผลรวมของ Qty * Price
func (o Order) Total() float64 {
	// TODO
	return 0
}

// ParseOrder แปลง JSON เป็น Order
// JSON ผิดรูปแบบ → error, ไม่มี id → ErrMissingID
func ParseOrder(data []byte) (Order, error) {
	// TODO
	return Order{}, nil
}

// Summary ต้อง marshal เป็น {"count":..., "revenue":..., "by_customer": {...}}
type Summary struct {
	Count      int
	Revenue    float64
	ByCustomer map[string]float64
}

// Summarize รับ JSON array ของ orders แล้วคืน JSON ของ Summary
func Summarize(data []byte) ([]byte, error) {
	// TODO
	return nil, nil
}

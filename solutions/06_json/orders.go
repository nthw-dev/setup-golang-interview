// Package orders — เฉลยข้อ 6: encoding/json, struct tags, omitempty, raw string
package orders

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrMissingID = errors.New("order id is required")

// struct tag บอก encoding/json ว่า field นี้ map กับ key อะไรใน JSON
// ต้องเป็น exported field (ตัวใหญ่) เท่านั้น json ถึงจะมองเห็น
type Item struct {
	SKU   string  `json:"sku"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

type Order struct {
	ID       string `json:"id"`
	Customer string `json:"customer_name"`
	Items    []Item `json:"items"`
	Note     string `json:"note,omitempty"` // omitempty: ไม่ใส่ key ถ้าเป็นค่าว่าง
}

// Total = ผลรวมของ Qty * Price
func (o Order) Total() float64 {
	total := 0.0
	for _, it := range o.Items {
		total += float64(it.Qty) * it.Price
	}
	return total
}

// ParseOrder แปลง JSON เป็น Order และตรวจว่ามี id
func ParseOrder(data []byte) (Order, error) {
	var o Order
	if err := json.Unmarshal(data, &o); err != nil {
		return Order{}, fmt.Errorf("parse order: %w", err)
	}
	if o.ID == "" {
		return Order{}, ErrMissingID
	}
	return o, nil
}

type Summary struct {
	Count      int                `json:"count"`
	Revenue    float64            `json:"revenue"`
	ByCustomer map[string]float64 `json:"by_customer"`
}

// Summarize รับ JSON array ของ orders แล้วคืน JSON ของ Summary
func Summarize(data []byte) ([]byte, error) {
	var orders []Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, fmt.Errorf("summarize: %w", err)
	}
	sum := Summary{ByCustomer: make(map[string]float64)}
	for _, o := range orders {
		t := o.Total()
		sum.Count++
		sum.Revenue += t
		sum.ByCustomer[o.Customer] += t
	}
	return json.Marshal(sum)
}

package orders

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// raw string literal (backtick) — ไม่ต้อง escape เครื่องหมาย " และขึ้นบรรทัดใหม่ได้
const orderJSON = `{
  "id": "A001",
  "customer_name": "Tah",
  "items": [
    {"sku": "KB-01", "qty": 2, "price": 10.5},
    {"sku": "MS-01", "qty": 1, "price": 4}
  ]
}`

func TestParseOrder(t *testing.T) {
	o, err := ParseOrder([]byte(orderJSON))
	if err != nil {
		t.Fatal(err)
	}
	if o.ID != "A001" || o.Customer != "Tah" || len(o.Items) != 2 {
		t.Fatalf("ParseOrder = %+v (hint: json tags)", o)
	}
	if o.Items[0].SKU != "KB-01" || o.Items[0].Qty != 2 {
		t.Errorf("first item = %+v", o.Items[0])
	}
	if o.Total() != 25 {
		t.Errorf("Total = %v; want 25", o.Total())
	}
}

func TestParseOrderErrors(t *testing.T) {
	if _, err := ParseOrder([]byte(`{"customer_name":"x"}`)); !errors.Is(err, ErrMissingID) {
		t.Errorf("missing id error = %v; want ErrMissingID", err)
	}
	if _, err := ParseOrder([]byte(`{not json`)); err == nil {
		t.Errorf("expected error for invalid JSON")
	}
}

func TestMarshalOmitEmpty(t *testing.T) {
	b, err := json.Marshal(Order{ID: "B1", Customer: "Bee"})
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `"customer_name":"Bee"`) {
		t.Errorf("marshal = %s; want key customer_name", s)
	}
	if strings.Contains(s, `"note"`) {
		t.Errorf("marshal = %s; empty note should be omitted", s)
	}
	b, _ = json.Marshal(Order{ID: "B2", Note: "fragile"})
	if !strings.Contains(string(b), `"note":"fragile"`) {
		t.Errorf("marshal = %s; want note", b)
	}
}

func TestSummarize(t *testing.T) {
	in := `[
	  {"id":"1","customer_name":"Tah","items":[{"sku":"a","qty":2,"price":10}]},
	  {"id":"2","customer_name":"Bee","items":[{"sku":"b","qty":1,"price":5.5}]},
	  {"id":"3","customer_name":"Tah","items":[{"sku":"c","qty":3,"price":1}]}
	]`
	out, err := Summarize([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("output is not valid JSON: %s", out)
	}
	want := map[string]any{
		"count":       3.0,
		"revenue":     28.5,
		"by_customer": map[string]any{"Tah": 23.0, "Bee": 5.5},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Summarize = %s\nwant %v", out, want)
	}
	if _, err := Summarize([]byte(`{}`)); err == nil {
		t.Errorf("expected error when input is not an array")
	}
}

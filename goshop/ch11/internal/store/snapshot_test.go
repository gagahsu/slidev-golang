package store

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"goshop/internal/shop"
)

func TestDecodeProducts(t *testing.T) {
	ps, err := DecodeProducts(strings.NewReader(
		`[{"sku":"A","name":"咖啡豆","price":450,"stock":3}]`))
	if err != nil || len(ps) != 1 || ps[0].Price != 450 {
		t.Fatalf("DecodeProducts() = %v, %v", ps, err)
	}
	_, err = DecodeProducts(strings.NewReader(`[{"sku":"A","colour":"red"}]`))
	if err == nil {
		t.Error("不認得的欄位應該回傳錯誤")
	}
}

func TestSaveLoad(t *testing.T) {
	m := NewMemory(shop.Product{SKU: "A", Name: "咖啡豆", Price: 450, Stock: 3})
	o := shop.Order{Lines: []shop.Line{{SKU: "A", Qty: 1}}, Total: 450}
	if err := m.PlaceOrder(&o); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := m.Save(&buf); err != nil {
		t.Fatal(err)
	}
	m2 := NewMemory()
	if err := m2.Load(&buf); err != nil {
		t.Fatal(err)
	}
	if p, _ := m2.Product("A"); p.Stock != 2 {
		t.Errorf("庫存 = %d，want 2", p.Stock)
	}
	if got, _ := m2.Order(1); got.Total != 450 {
		t.Errorf("訂單金額 = %v，want 450", got.Total)
	}
}

func TestStatusJSON(t *testing.T) {
	b, err := json.Marshal(shop.Order{ID: 7, Status: shop.Paid})
	if err != nil || !strings.Contains(string(b), `"status":"paid"`) {
		t.Fatalf("Marshal() = %s, %v", b, err)
	}
	var o shop.Order
	if err := json.Unmarshal(b, &o); err != nil || o.Status != shop.Paid {
		t.Errorf("Unmarshal() = %v, %v", o.Status, err)
	}
	err = json.Unmarshal([]byte(`{"status":"lost"}`), &o)
	if err == nil {
		t.Error("未知的狀態應該回傳錯誤")
	}
}

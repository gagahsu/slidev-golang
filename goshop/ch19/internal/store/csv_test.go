package store

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"goshop/internal/shop"
)

func TestReadProductsCSV(t *testing.T) {
	in := "sku,name,price,stock\nA,咖啡豆,450,3\nB,磨豆機,二千,1\nC,,-5,1\n"
	ps, err := ReadProductsCSV(strings.NewReader(in))
	if len(ps) != 1 || ps[0].SKU != "A" {
		t.Errorf("products = %v，want 只有 A", ps)
	}
	if _, ok := errors.AsType[*strconv.NumError](err); !ok {
		t.Errorf("err = %v，want *strconv.NumError", err)
	}
	for _, want := range []string{"第 3 行", "第 4 行：欄位 name", "欄位 price 不符合規則 min=1"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v，應該包含 %q", err, want)
		}
	}
}

func TestWriteOrdersCSV(t *testing.T) {
	var sb strings.Builder
	orders := []shop.Order{{ID: 1, Total: 1880, Status: shop.Paid, PaidBy: "貨到付款"}}
	if err := WriteOrdersCSV(&sb, orders); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 2 || !strings.HasSuffix(lines[1], ",1880,已付款,貨到付款") {
		t.Errorf("CSV = %q", sb.String())
	}
}

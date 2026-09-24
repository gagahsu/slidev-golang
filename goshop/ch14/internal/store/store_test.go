package store

import (
	"errors"
	"os"
	"testing"
	"time"

	"goshop/internal/shop"
)

// testStore 是「合約測試」：任何 Store 實作都必須通過同一組測試。
func testStore(t *testing.T, st Store) {
	ctx := t.Context()
	for _, p := range []shop.Product{
		{SKU: "A", Name: "咖啡豆", Price: 450, Stock: 10},
		{SKU: "B", Name: "手沖壺", Price: 1280, Stock: 1},
	} {
		if err := st.SaveProduct(ctx, p); err != nil {
			t.Fatal(err)
		}
	}

	o := shop.Order{
		Lines:    []shop.Line{{SKU: "A", Name: "咖啡豆", Price: 450, Qty: 2}},
		Subtotal: 900, Total: 900,
		CreatedAt: time.Now().Truncate(time.Microsecond),
	}
	if err := st.PlaceOrder(ctx, &o); err != nil || o.ID == 0 {
		t.Fatalf("PlaceOrder() id = %d, err = %v", o.ID, err)
	}

	// 其中一項庫存不足：整張訂單失敗，A 的庫存也不能被扣
	bad := shop.Order{Lines: []shop.Line{{SKU: "A", Qty: 1}, {SKU: "B", Qty: 2}}}
	err := st.PlaceOrder(ctx, &bad)
	if se, ok := errors.AsType[*shop.StockError](err); !ok || se.Have != 1 {
		t.Errorf("PlaceOrder() err = %v，want B 的 StockError", err)
	}
	if p, _ := st.Product(ctx, "A"); p.Stock != 8 {
		t.Errorf("A 的庫存 = %d，want 8", p.Stock)
	}

	if err := st.MarkPaid(ctx, o.ID, "貨到付款"); err != nil {
		t.Fatal(err)
	}
	if err := st.MarkPaid(ctx, o.ID, "貨到付款"); !errors.Is(err, shop.ErrPaid) {
		t.Errorf("重複付款 err = %v，want ErrPaid", err)
	}
	got, err := st.Order(ctx, o.ID)
	if err != nil || got.Status != shop.Paid || len(got.Lines) != 1 ||
		!got.CreatedAt.Equal(o.CreatedAt) {
		t.Errorf("Order() = %+v, %v", got, err)
	}
	if _, err := st.Order(ctx, 999); !errors.Is(err, shop.ErrNotFound) {
		t.Errorf("Order(999) err = %v，want ErrNotFound", err)
	}
	if orders, _ := st.Orders(ctx); len(orders) != 1 {
		t.Errorf("Orders() 有 %d 張，want 1", len(orders))
	}
}

func TestMemory(t *testing.T) {
	testStore(t, NewMemory())
}

// 設定環境變數才會執行，例如：
// GOSHOP_TEST_DSN="gouser:gopass@tcp(127.0.0.1:3306)/goshop_test" go test ./...
func TestMySQL(t *testing.T) {
	dsn := os.Getenv("GOSHOP_TEST_DSN")
	if dsn == "" {
		t.Skip("沒有設定 GOSHOP_TEST_DSN，略過 MySQL 測試")
	}
	st, err := OpenMySQL(t.Context(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	for _, table := range []string{"order_lines", "orders", "products"} {
		if _, err := st.db.ExecContext(t.Context(), "DELETE FROM "+table); err != nil {
			t.Fatal(err)
		}
	}
	testStore(t, st)
}

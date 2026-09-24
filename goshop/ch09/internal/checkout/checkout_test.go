package checkout

import (
	"errors"
	"testing"

	"goshop/internal/money"
	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

// newService 建立一個有兩項商品的測試用結帳服務
func newService() *Service {
	st := store.NewMemory(
		shop.Product{SKU: "A", Name: "咖啡豆", Price: 450, Stock: 10},
		shop.Product{SKU: "B", Name: "手沖壺", Price: 1280, Stock: 1},
	)
	return &Service{Store: st, Rules: []Discount{
		PercentOff(10), Threshold(2000, 300)}}
}

func TestBest(t *testing.T) {
	rules := []Discount{PercentOff(10), Threshold(2000, 300)}
	tests := []struct {
		subtotal, want money.Money
	}{
		{0, 0},
		{1000, 100}, // 9 折比較划算
		{2000, 300}, // 滿額折 300 比較划算
		{5000, 500}, // 9 折又贏回來
		{1999, 199},
	}
	for _, tt := range tests {
		if got := Best(tt.subtotal, rules...); got != tt.want {
			t.Errorf("Best(%d) = %d，want %d", tt.subtotal, got, tt.want)
		}
	}
}

func TestCheckout(t *testing.T) {
	svc := newService()
	var cart Cart
	cart.Add(Item{"A", 2}, Item{"B", 1})

	o, err := svc.Checkout(cart)
	if err != nil {
		t.Fatalf("Checkout() error = %v", err)
	}
	if o.ID != 1 || o.Total != 1880 {
		t.Errorf("訂單 = #%d %v，want #1 NT$1,880", o.ID, o.Total)
	}
	if p, _ := svc.Store.Product("A"); p.Stock != 8 {
		t.Errorf("A 的庫存 = %d，want 8", p.Stock)
	}
}

func TestCheckoutErrors(t *testing.T) {
	tests := []struct {
		name  string
		items []Item
		want  error
	}{
		{"空的購物車", nil, shop.ErrEmptyCart},
		{"商品不存在", []Item{{"Z", 1}}, shop.ErrNotFound},
		{"庫存不足", []Item{{"A", 1}, {"B", 2}}, shop.ErrOutOfStock},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService()
			_, err := svc.Checkout(Cart{Items: tt.items})
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v，want %v", err, tt.want)
			}
			// 結帳失敗時，庫存不能被扣掉
			if p, _ := svc.Store.Product("A"); p.Stock != 10 {
				t.Errorf("A 的庫存 = %d，want 10", p.Stock)
			}
		})
	}
}

func TestPay(t *testing.T) {
	svc := newService()
	o, _ := svc.Checkout(Cart{Items: []Item{{"A", 1}}})

	wallet := &payment.Wallet{Balance: 100}
	err := svc.Pay(o.ID, wallet, payment.CashOnDelivery{})
	if err != nil {
		t.Fatalf("Pay() error = %v", err)
	}
	got, _ := svc.Store.Order(o.ID)
	if got.Status != shop.Paid || got.PaidBy != "貨到付款" {
		t.Errorf("訂單 = %v／%s，want 已付款／貨到付款", got.Status, got.PaidBy)
	}
	if err := svc.Pay(o.ID, wallet); !errors.Is(err, shop.ErrPaid) {
		t.Errorf("重複付款 err = %v，want ErrPaid", err)
	}
}

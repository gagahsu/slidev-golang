package checkout

import (
	"errors"
	"testing"
	"time"

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

func TestCoupon(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	// 2026-09-25 是星期五，下午 3 點下單
	now := time.Date(2026, 9, 25, 15, 0, 0, 0, loc)

	svc := newService()
	svc.Now = func() time.Time { return now }
	svc.Coupons = map[string]Coupon{
		"FALL15": {"FALL15", 15, now.AddDate(0, 0, -7), now.AddDate(0, 0, 7)},
		"SUMMER": {"SUMMER", 20, now.AddDate(0, -3, 0), now.AddDate(0, -1, 0)},
	}

	o, err := svc.Checkout(Cart{Items: []Item{{"A", 2}}, Coupon: "FALL15"})
	if err != nil {
		t.Fatalf("Checkout() error = %v", err)
	}
	if o.Discount != 135 { // 900 的 15% 比 9 折多
		t.Errorf("Discount = %d，want 135", o.Discount)
	}
	// 星期五下單，2 個工作天後是下星期二
	if got := o.ShipBy.Weekday(); got != time.Tuesday {
		t.Errorf("ShipBy = %v，want Tuesday", got)
	}

	for _, code := range []string{"SUMMER", "NOPE"} {
		_, err := svc.Checkout(Cart{Items: []Item{{"A", 1}}, Coupon: code})
		if !errors.Is(err, shop.ErrCoupon) {
			t.Errorf("%s: err = %v，want ErrCoupon", code, err)
		}
	}
}

func TestCouponValid(t *testing.T) {
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	c := Coupon{Start: start, End: start.AddDate(0, 1, 0)}
	tests := []struct {
		t    time.Time
		want bool
	}{
		{start.Add(-time.Nanosecond), false},
		{start, true}, // 包含生效時間
		{c.End.Add(-time.Second), true},
		{c.End, false}, // 不包含失效時間
	}
	for _, tt := range tests {
		if got := c.Valid(tt.t); got != tt.want {
			t.Errorf("Valid(%v) = %v，want %v", tt.t, got, tt.want)
		}
	}
}

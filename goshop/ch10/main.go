package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"goshop/internal/checkout"
	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

var weekdays = [...]string{"日", "一", "二", "三", "四", "五", "六"}

// formatTime 以商店時區顯示，例如 2026-09-25 15:04（五）
func formatTime(t time.Time) string {
	t = t.In(shop.Location)
	return t.Format("2006-01-02 15:04") + "（" + weekdays[t.Weekday()] + "）"
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	st := store.NewMemory(
		shop.Product{SKU: "SKU-001", Name: "衣索比亞咖啡豆", Price: 450, Stock: 20},
		shop.Product{SKU: "SKU-002", Name: "濾掛咖啡（10入）", Price: 280, Stock: 50},
		shop.Product{SKU: "SKU-003", Name: "手沖壺", Price: 1280, Stock: 5},
		shop.Product{SKU: "SKU-004", Name: "馬克杯", Price: 350, Stock: 30},
	)
	now := time.Now()
	svc := &checkout.Service{
		Store: st,
		Rules: []checkout.Discount{
			checkout.PercentOff(10), checkout.Threshold(2000, 300)},
		Coupons: map[string]checkout.Coupon{
			// 本週限定 85 折（前後 7 天內有效）
			"WEEK15": {Code: "WEEK15", PercentOff: 15,
				Start: now.AddDate(0, 0, -7), End: now.AddDate(0, 0, 7)},
			// 上個月就結束的舊活動
			"SUMMER": {Code: "SUMMER", PercentOff: 20,
				Start: now.AddDate(0, -3, 0), End: now.AddDate(0, -1, 0)},
		},
	}

	carts := []checkout.Cart{
		{Items: []checkout.Item{{SKU: "SKU-002", Qty: 1}}, Coupon: "SUMMER"},
		{Items: []checkout.Item{{SKU: "SKU-001", Qty: 4}}, Coupon: "WEEK15"},
		{Items: []checkout.Item{{SKU: "SKU-003", Qty: 2}}},
	}
	for _, cart := range carts {
		o, err := svc.Checkout(cart)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := svc.Pay(o.ID, payment.CashOnDelivery{}); err != nil {
			fmt.Println(err)
		}
		o, _ = st.Order(o.ID)
		printReceipt(o)
	}
}

func printReceipt(o shop.Order) {
	fmt.Printf("===== 訂單 #%d（%v）=====\n", o.ID, o.Status)
	defer fmt.Println("===== 謝謝光臨 =====")
	fmt.Println("下單時間：", formatTime(o.CreatedAt))
	for _, l := range o.Lines {
		fmt.Printf("%s x %d\t%v\n", l.Name, l.Qty, l.Price.Times(l.Qty))
	}
	fmt.Println("小計：", o.Subtotal, "折扣：", o.Discount, o.Coupon)
	fmt.Println("應付：", o.Total, "／付款方式：", o.PaidBy)
	fmt.Println("預計出貨：", o.ShipBy.Format(time.DateOnly))
}

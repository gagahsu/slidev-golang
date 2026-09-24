package main

import (
	"fmt"
	"log/slog"
	"os"

	"goshop/internal/checkout"
	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

func main() {
	// 日誌輸出到 stderr，並顯示 Debug 等級以上的訊息
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelDebug})))

	st := store.NewMemory(
		shop.Product{SKU: "SKU-001", Name: "衣索比亞咖啡豆", Price: 450, Stock: 20},
		shop.Product{SKU: "SKU-002", Name: "濾掛咖啡（10入）", Price: 280, Stock: 50},
		shop.Product{SKU: "SKU-003", Name: "手沖壺", Price: 1280, Stock: 5},
		shop.Product{SKU: "SKU-004", Name: "馬克杯", Price: 350, Stock: 30},
	)
	svc := &checkout.Service{
		Store: st,
		Rules: []checkout.Discount{
			checkout.PercentOff(10), checkout.Threshold(2000, 300)},
	}
	wallet := &payment.Wallet{Balance: 1000}

	var cart checkout.Cart
	cart.Add(checkout.Item{SKU: "SKU-001", Qty: 2},
		checkout.Item{SKU: "SKU-003", Qty: 1})

	o, err := svc.Checkout(cart)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := svc.Pay(o.ID, wallet); err != nil {
		fmt.Println(err) // 錢包餘額不足，訂單保留為待付款
	}
	if err := svc.Pay(o.ID, payment.CashOnDelivery{}); err != nil {
		fmt.Println(err)
	}

	o, _ = st.Order(o.ID)
	printReceipt(o)
}

func printReceipt(o shop.Order) {
	fmt.Printf("===== 訂單 #%d（%v）=====\n", o.ID, o.Status)
	defer fmt.Println("===== 謝謝光臨 =====")
	for _, l := range o.Lines {
		fmt.Printf("%s x %d\t%v\n", l.Name, l.Qty, l.Price.Times(l.Qty))
	}
	fmt.Println("小計：", o.Subtotal, "折扣：", o.Discount)
	fmt.Println("應付：", o.Total, "／付款方式：", o.PaidBy)
}

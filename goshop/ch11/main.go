package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
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

	st, err := openStore("data")
	if err != nil {
		slog.Error("開啟資料失敗", "err", err)
		os.Exit(1)
	}
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

	// 把最後一張訂單編碼成排版過的 JSON
	orders, _ := st.Orders()
	if len(orders) > 0 {
		b, _ := json.MarshalIndent(orders[len(orders)-1], "", "  ")
		fmt.Println(string(b))
	}
	if err := saveStore(st, "data"); err != nil {
		slog.Error("儲存資料失敗", "err", err)
	}
}

// openStore 優先讀取上次存下的 gob 快照；
// 第一次執行時沒有快照，就從 products.json 匯入商品。
func openStore(dir string) (*store.Memory, error) {
	st := store.NewMemory()
	f, err := os.Open(filepath.Join(dir, "goshop.gob"))
	if err == nil {
		defer f.Close()
		return st, st.Load(f)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	f, err = os.Open(filepath.Join(dir, "products.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	products, err := store.DecodeProducts(f)
	if err != nil {
		return nil, err
	}
	return store.NewMemory(products...), nil
}

// saveStore 把目前的資料存成 gob 快照，下次啟動時接著用。
func saveStore(st *store.Memory, dir string) error {
	f, err := os.Create(filepath.Join(dir, "goshop.gob"))
	if err != nil {
		return err
	}
	defer f.Close()
	return st.Save(f)
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

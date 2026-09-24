package main

import "fmt"

// 商店名稱與免運門檻
const (
	shopName     = "GoShop"
	freeShipping = 1000
)

// 會員等級：每升一級多打 5% 折扣
const (
	Normal = iota // 一般會員
	Silver        // 銀卡會員
	Gold          // 金卡會員
)

// sell 賣出商品後扣掉庫存
func sell(stock *int, qty int) {
	*stock -= qty
}

func main() {
	name := "衣索比亞咖啡豆"
	price := 450
	stock := 20
	qty := 3
	level := Gold

	subtotal := price * qty
	discount := subtotal * level * 5 / 100
	total := subtotal - discount
	free := total >= freeShipping && qty > 0

	sell(&stock, qty)

	fmt.Println("===== 歡迎光臨", shopName, "=====")
	fmt.Println("商品：", name, "x", qty)
	fmt.Println("小計：", subtotal)
	fmt.Println("折扣：", discount)
	fmt.Println("總計：", total)
	fmt.Println("免運：", free)
	fmt.Println("剩餘庫存：", stock)
}

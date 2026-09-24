package main

import "fmt"

const (
	shopName     = "GoShop"
	freeShipping = 1000 // 滿 1000 元免運
	shippingFee  = 60   // 未滿免運門檻的運費
)

const (
	Normal = iota
	Silver
	Gold
)

// levelName 傳回會員等級的名稱
func levelName(level int) string {
	switch level {
	case Gold:
		return "金卡會員"
	case Silver:
		return "銀卡會員"
	default:
		return "一般會員"
	}
}

// discountPercent 傳回會員等級的折扣百分比
func discountPercent(level int) int {
	switch level {
	case Gold:
		return 10
	case Silver:
		return 5
	}
	return 0
}

// shipping 計算運費：滿額免運就提早 return
func shipping(total int) int {
	if total >= freeShipping {
		return 0
	}
	return shippingFee
}

func main() {
	name := "衣索比亞咖啡豆"
	price := 450
	stock := 4
	level := Silver

	fmt.Println("=====", shopName, "價格試算 =====")
	fmt.Println("商品：", name, "／", levelName(level))
	fmt.Println("數量   小計   折扣   運費   應付")
	for qty := range 6 {
		qty++ // 從 1 件開始
		if qty > stock {
			fmt.Println("庫存只剩", stock, "件，停止試算")
			break
		}
		subtotal := price * qty
		discount := subtotal * discountPercent(level) / 100
		total := subtotal - discount
		fee := shipping(total)
		fmt.Printf("%4d %6d %6d %6d %6d\n",
			qty, subtotal, discount, fee, total+fee)
	}
}

package main

import (
	"fmt"
	"strconv"
)

type Product struct {
	SKU   string
	Name  string
	Price int
	Stock int
}

type Item struct {
	SKU string
	Qty int
}

type Cart struct {
	Items []Item
}

// Add 可以一次放入多個品項（參數不定函式）
func (c *Cart) Add(items ...Item) {
	for _, item := range items {
		c.add(item)
	}
}

func (c *Cart) add(item Item) {
	for i := range c.Items {
		if c.Items[i].SKU == item.SKU {
			c.Items[i].Qty += item.Qty
			return
		}
	}
	c.Items = append(c.Items, item)
}

func (c Cart) Total(catalog map[string]Product) int {
	total := 0
	for _, it := range c.Items {
		total += catalog[it.SKU].Price * it.Qty
	}
	return total
}

// Discount 傳入小計，傳回可以折抵的金額
type Discount func(subtotal int) int

// PercentOff 建立「打 p 折扣」的折扣規則，例如 10 代表 9 折
func PercentOff(p int) Discount {
	return func(subtotal int) int {
		return subtotal * p / 100
	}
}

// Threshold 建立「滿 limit 元折 off 元」的折扣規則
func Threshold(limit, off int) Discount {
	return func(subtotal int) int {
		if subtotal >= limit {
			return off
		}
		return 0
	}
}

// Best 從多個折扣規則中挑出折最多的那一個
func Best(subtotal int, rules ...Discount) int {
	best := 0
	for _, rule := range rules {
		best = max(best, rule(subtotal))
	}
	return best
}

// formatNT 把金額加上千分位，例如 12345 → NT$12,345
func formatNT(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return "NT$" + s
}

// printReceipt 印出收據，並傳回應付金額
func printReceipt(cart Cart, catalog map[string]Product,
	rules ...Discount) (total int) {
	fmt.Println("========== GoShop 收據 ==========")
	defer fmt.Println("=========== 謝謝光臨 ===========")

	for _, it := range cart.Items {
		p := catalog[it.SKU]
		fmt.Printf("%s x %d\t%s\n", p.Name, it.Qty, formatNT(p.Price*it.Qty))
	}
	subtotal := cart.Total(catalog)
	discount := Best(subtotal, rules...)
	total = subtotal - discount
	fmt.Println("小計：", formatNT(subtotal))
	fmt.Println("折扣：", formatNT(discount))
	fmt.Println("應付：", formatNT(total))
	return total
}

func main() {
	catalog := map[string]Product{
		"SKU-001": {"SKU-001", "衣索比亞咖啡豆", 450, 20},
		"SKU-002": {"SKU-002", "濾掛咖啡（10入）", 280, 50},
		"SKU-003": {"SKU-003", "手沖壺", 1280, 5},
		"SKU-004": {"SKU-004", "馬克杯", 350, 30},
	}

	var cart Cart
	cart.Add(Item{"SKU-001", 2}, Item{"SKU-003", 1}, Item{"SKU-004", 2})

	// 週年慶：全館 9 折，或滿 2000 折 300，取優惠較多者
	rules := []Discount{PercentOff(10), Threshold(2000, 300)}
	total := printReceipt(cart, catalog, rules...)
	fmt.Println("本次消費：", formatNT(total))
}

package main

import "fmt"

// Product 是一項商品
type Product struct {
	SKU   string
	Name  string
	Price int
	Stock int
}

// Item 是購物車裡的一個品項
type Item struct {
	SKU string
	Qty int
}

// Cart 是購物車
type Cart struct {
	Items []Item
}

// Add 把商品放進購物車；同一個 SKU 會合併數量
func (c *Cart) Add(sku string, qty int) {
	for i := range c.Items {
		if c.Items[i].SKU == sku {
			c.Items[i].Qty += qty
			return
		}
	}
	c.Items = append(c.Items, Item{SKU: sku, Qty: qty})
}

// Total 用商品目錄查價，計算購物車總金額
func (c Cart) Total(catalog map[string]Product) int {
	total := 0
	for _, it := range c.Items {
		total += catalog[it.SKU].Price * it.Qty
	}
	return total
}

func main() {
	products := []Product{
		{"SKU-001", "衣索比亞咖啡豆", 450, 20},
		{"SKU-002", "濾掛咖啡（10入）", 280, 50},
		{"SKU-003", "手沖壺", 1280, 5},
		{"SKU-004", "馬克杯", 350, 30},
	}

	// 建立以 SKU 為鍵的商品目錄
	catalog := make(map[string]Product, len(products))
	for _, p := range products {
		catalog[p.SKU] = p
	}

	var cart Cart
	cart.Add("SKU-001", 2)
	cart.Add("SKU-003", 1)
	cart.Add("SKU-001", 1) // 同一商品會合併

	for _, sku := range []string{"SKU-004", "SKU-999"} {
		if _, ok := catalog[sku]; !ok { // comma ok 慣用法
			fmt.Println("找不到商品：", sku)
			continue
		}
		cart.Add(sku, 2)
	}

	fmt.Println("===== 購物車 =====")
	for _, it := range cart.Items {
		p := catalog[it.SKU]
		fmt.Printf("%s %s x %d = %d\n",
			p.SKU, p.Name, it.Qty, p.Price*it.Qty)
	}
	fmt.Println("品項數：", len(cart.Items))
	fmt.Println("總金額：", cart.Total(catalog))
}

package main

import (
	"errors"
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

type Discount func(subtotal int) int

func PercentOff(p int) Discount {
	return func(subtotal int) int {
		return subtotal * p / 100
	}
}

func Threshold(limit, off int) Discount {
	return func(subtotal int) int {
		if subtotal >= limit {
			return off
		}
		return 0
	}
}

func Best(subtotal int, rules ...Discount) int {
	best := 0
	for _, rule := range rules {
		best = max(best, rule(subtotal))
	}
	return best
}

func formatNT(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return "NT$" + s
}

// 哨兵錯誤：呼叫端可以用 errors.Is 判斷
var (
	ErrEmptyCart  = errors.New("購物車是空的")
	ErrUnknownSKU = errors.New("找不到商品")
)

// StockError 表示某項商品庫存不足
type StockError struct {
	SKU  string
	Want int // 想買幾件
	Have int // 實際剩幾件
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s 庫存不足：想買 %d 件，只剩 %d 件",
		e.SKU, e.Want, e.Have)
}

// Line 是訂單的一行明細，記下成交當下的品名和單價
type Line struct {
	SKU   string
	Name  string
	Price int
	Qty   int
}

// Order 是一張訂單
type Order struct {
	ID       int
	Lines    []Line
	Subtotal int
	Discount int
	Total    int
}

var lastOrderID = 0

// validate 檢查購物車的每個品項，把所有問題一次回報
func validate(catalog map[string]Product, cart Cart) error {
	if len(cart.Items) == 0 {
		return ErrEmptyCart
	}
	var errs []error
	for _, it := range cart.Items {
		p, ok := catalog[it.SKU]
		switch {
		case !ok:
			errs = append(errs, fmt.Errorf("%s: %w", it.SKU, ErrUnknownSKU))
		case p.Stock < it.Qty:
			errs = append(errs, &StockError{it.SKU, it.Qty, p.Stock})
		}
	}
	return errors.Join(errs...) // 沒有錯誤時傳回 nil
}

// Checkout 檢查購物車、扣庫存並建立訂單
func Checkout(catalog map[string]Product, cart Cart,
	rules ...Discount) (Order, error) {
	if err := validate(catalog, cart); err != nil {
		return Order{}, fmt.Errorf("結帳失敗：%w", err)
	}

	lastOrderID++
	o := Order{ID: lastOrderID}
	for _, it := range cart.Items {
		p := catalog[it.SKU]
		p.Stock -= it.Qty // 全部檢查通過才扣庫存
		catalog[it.SKU] = p
		o.Lines = append(o.Lines, Line{p.SKU, p.Name, p.Price, it.Qty})
		o.Subtotal += p.Price * it.Qty
	}
	o.Discount = Best(o.Subtotal, rules...)
	o.Total = o.Subtotal - o.Discount
	return o, nil
}

func printReceipt(o Order) {
	fmt.Printf("===== 訂單 #%d =====\n", o.ID)
	defer fmt.Println("===== 謝謝光臨 =====")
	for _, l := range o.Lines {
		fmt.Printf("%s x %d\t%s\n", l.Name, l.Qty, formatNT(l.Price*l.Qty))
	}
	fmt.Println("應付：", formatNT(o.Total))
}

func main() {
	catalog := map[string]Product{
		"SKU-001": {"SKU-001", "衣索比亞咖啡豆", 450, 20},
		"SKU-002": {"SKU-002", "濾掛咖啡（10入）", 280, 50},
		"SKU-003": {"SKU-003", "手沖壺", 1280, 5},
		"SKU-004": {"SKU-004", "馬克杯", 350, 30},
	}
	rules := []Discount{PercentOff(10), Threshold(2000, 300)}

	carts := []Cart{
		{},
		{Items: []Item{{"SKU-003", 9}, {"SKU-999", 1}}},
		{Items: []Item{{"SKU-001", 2}, {"SKU-003", 1}}},
	}
	for _, cart := range carts {
		o, err := Checkout(catalog, cart, rules...)
		if err != nil {
			fmt.Println(err)
			if se, ok := errors.AsType[*StockError](err); ok {
				fmt.Printf("→ 建議把 %s 改成 %d 件\n", se.SKU, se.Have)
			}
			if errors.Is(err, ErrUnknownSKU) {
				fmt.Println("→ 請移除已下架的商品")
			}
			continue
		}
		printReceipt(o)
	}
	fmt.Println("手沖壺剩餘庫存：", catalog["SKU-003"].Stock)
}

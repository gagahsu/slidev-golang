package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Product struct {
	SKU   string
	Name  string
	Price Money
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

type Discount func(subtotal Money) Money

func PercentOff(p int) Discount {
	return func(subtotal Money) Money {
		return subtotal * Money(p) / 100
	}
}

func Threshold(limit, off Money) Discount {
	return func(subtotal Money) Money {
		if subtotal >= limit {
			return off
		}
		return 0
	}
}

func Best(subtotal Money, rules ...Discount) Money {
	best := Money(0)
	for _, rule := range rules {
		best = max(best, rule(subtotal))
	}
	return best
}

// Money 是新台幣金額（元）
type Money int

// String 讓 Money 實作 fmt.Stringer，印出 NT$1,234
func (m Money) String() string {
	s := strconv.Itoa(int(m))
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
	Price Money
	Qty   int
}

// Order 是一張訂單
type Order struct {
	ID       int
	Lines    []Line
	Subtotal Money
	Discount Money
	Total    Money
	PaidBy   string // 付款方式，空字串代表尚未付款
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
			err := fmt.Errorf("%s: %w", it.SKU, ErrUnknownSKU)
			errs = append(errs, err)
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
		o.Subtotal += p.Price * Money(it.Qty)
	}
	o.Discount = Best(o.Subtotal, rules...)
	o.Total = o.Subtotal - o.Discount
	return o, nil
}

func printReceipt(o Order) {
	fmt.Printf("===== 訂單 #%d =====\n", o.ID)
	defer fmt.Println("===== 謝謝光臨 =====")
	for _, l := range o.Lines {
		fmt.Printf("%s x %d\t%v\n", l.Name, l.Qty, l.Price*Money(l.Qty))
	}
	fmt.Println("應付：", o.Total, "／付款方式：", o.PaidBy)
}

// ErrInsufficientFunds 表示錢包餘額不足
var ErrInsufficientFunds = errors.New("餘額不足")

// PaymentMethod 是所有付款方式都要具備的行為
type PaymentMethod interface {
	Name() string
	Pay(amount Money) error
}

// CreditCard 是信用卡，超過額度會被拒絕
type CreditCard struct {
	Last4 string
	Limit Money
}

func (c CreditCard) Name() string { return "信用卡 *" + c.Last4 }

func (c CreditCard) Pay(amount Money) error {
	if amount > c.Limit {
		return fmt.Errorf("超過信用額度 %v", c.Limit)
	}
	return nil
}

// Wallet 是 GoShop 錢包；付款會改變餘額，所以用指標接收器
type Wallet struct {
	Balance Money
}

func (w *Wallet) Name() string { return "GoShop 錢包" }

func (w *Wallet) Pay(amount Money) error {
	if w.Balance < amount {
		return fmt.Errorf("%w，只剩 %v", ErrInsufficientFunds, w.Balance)
	}
	w.Balance -= amount
	return nil
}

// CashOnDelivery 是貨到付款
type CashOnDelivery struct{}

func (CashOnDelivery) Name() string           { return "貨到付款" }
func (CashOnDelivery) Pay(amount Money) error { return nil }

// 在編譯時期確認三種付款方式都實作了 PaymentMethod
var (
	_ PaymentMethod = CreditCard{}
	_ PaymentMethod = (*Wallet)(nil)
	_ PaymentMethod = CashOnDelivery{}
)

// Pay 依序嘗試付款方式，直到其中一種成功
func (o *Order) Pay(methods ...PaymentMethod) error {
	var errs []error
	for _, m := range methods {
		err := m.Pay(o.Total)
		if err == nil {
			o.PaidBy = m.Name()
			return nil
		}
		errs = append(errs, fmt.Errorf("%s：%w", m.Name(), err))
	}
	return fmt.Errorf("訂單 #%d 付款失敗：%w", o.ID, errors.Join(errs...))
}

// describe 用型別 switch 針對不同付款方式顯示額外資訊
func describe(m PaymentMethod) string {
	switch v := m.(type) {
	case *Wallet:
		return fmt.Sprintf("%s（餘額 %v）", v.Name(), v.Balance)
	case CreditCard:
		return fmt.Sprintf("%s（額度 %v）", v.Name(), v.Limit)
	default:
		return m.Name()
	}
}

func main() {
	catalog := map[string]Product{
		"SKU-001": {"SKU-001", "衣索比亞咖啡豆", 450, 20},
		"SKU-002": {"SKU-002", "濾掛咖啡（10入）", 280, 50},
		"SKU-003": {"SKU-003", "手沖壺", 1280, 5},
		"SKU-004": {"SKU-004", "馬克杯", 350, 30},
	}
	rules := []Discount{PercentOff(10), Threshold(2000, 300)}

	wallet := &Wallet{Balance: 1000}
	card := CreditCard{Last4: "4242", Limit: 30000}
	methods := []PaymentMethod{wallet, card}

	for _, m := range methods {
		fmt.Println("付款方式：", describe(m))
	}

	carts := []Cart{
		{Items: []Item{{"SKU-002", 3}}},
		{Items: []Item{{"SKU-001", 2}, {"SKU-003", 1}}},
	}
	for _, cart := range carts {
		o, err := Checkout(catalog, cart, rules...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := o.Pay(methods...); err != nil {
			fmt.Println(err)
			continue
		}
		printReceipt(o)
	}
	fmt.Println("錢包餘額：", wallet.Balance)
}

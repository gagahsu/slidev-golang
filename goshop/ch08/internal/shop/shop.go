// Package shop 定義 GoShop 的核心資料：商品、訂單與錯誤。
package shop

import (
	"errors"
	"fmt"

	"goshop/internal/money"
)

// Product 是一項商品。
type Product struct {
	SKU   string
	Name  string
	Price money.Money
	Stock int
}

// Line 是訂單的一行明細，記下成交當下的品名和單價。
type Line struct {
	SKU   string
	Name  string
	Price money.Money
	Qty   int
}

// Status 是訂單狀態。
type Status int

const (
	Pending Status = iota // 待付款
	Paid                  // 已付款
)

func (s Status) String() string {
	if s == Paid {
		return "已付款"
	}
	return "待付款"
}

// Order 是一張訂單。
type Order struct {
	ID       int
	Lines    []Line
	Subtotal money.Money
	Discount money.Money
	Total    money.Money
	Status   Status
	PaidBy   string
}

var (
	ErrNotFound   = errors.New("找不到資料")
	ErrEmptyCart  = errors.New("購物車是空的")
	ErrPaid       = errors.New("訂單已付款")
	ErrOutOfStock = errors.New("庫存不足")
)

// StockError 表示某項商品庫存不足。
type StockError struct {
	SKU  string
	Want int
	Have int
}

func (e *StockError) Error() string {
	return fmt.Sprintf("%s 庫存不足：想買 %d 件，只剩 %d 件",
		e.SKU, e.Want, e.Have)
}

// Unwrap 讓 errors.Is(err, ErrOutOfStock) 成立。
func (e *StockError) Unwrap() error { return ErrOutOfStock }

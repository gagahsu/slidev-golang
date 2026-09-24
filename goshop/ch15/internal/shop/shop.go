// Package shop 定義 GoShop 的核心資料：商品、訂單與錯誤。
package shop

import (
	"errors"
	"fmt"
	"time"

	"goshop/internal/money"
)

// Product 是一項商品。
type Product struct {
	SKU   string      `json:"sku"`
	Name  string      `json:"name"`
	Price money.Money `json:"price"`
	Stock int         `json:"stock"`
}

// Line 是訂單的一行明細，記下成交當下的品名和單價。
type Line struct {
	SKU   string      `json:"sku"`
	Name  string      `json:"name"`
	Price money.Money `json:"price"`
	Qty   int         `json:"qty"`
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

// MarshalText 讓 JSON 輸出 "pending"／"paid"，而不是看不懂的數字。
func (s Status) MarshalText() ([]byte, error) {
	if s == Paid {
		return []byte("paid"), nil
	}
	return []byte("pending"), nil
}

// UnmarshalText 把 "pending"／"paid" 轉回 Status。
func (s *Status) UnmarshalText(b []byte) error {
	switch string(b) {
	case "pending":
		*s = Pending
	case "paid":
		*s = Paid
	default:
		return fmt.Errorf("未知的訂單狀態 %q", b)
	}
	return nil
}

// Order 是一張訂單。
type Order struct {
	ID        int         `json:"id"`
	Lines     []Line      `json:"lines"`
	Subtotal  money.Money `json:"subtotal"`
	Discount  money.Money `json:"discount"`
	Total     money.Money `json:"total"`
	Status    Status      `json:"status"`
	PaidBy    string      `json:"paid_by,omitzero"`
	Coupon    string      `json:"coupon,omitzero"` // 折價券代碼
	CreatedAt time.Time   `json:"created_at"`      // 下單時間
	ShipBy    time.Time   `json:"ship_by"`         // 預計出貨日
}

var (
	ErrNotFound   = errors.New("找不到資料")
	ErrEmptyCart  = errors.New("購物車是空的")
	ErrPaid       = errors.New("訂單已付款")
	ErrOutOfStock = errors.New("庫存不足")
	ErrCoupon     = errors.New("折價券無效或已過期")
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

// Location 是商店所在的時區，下單時間和出貨日都以它為準。
var Location, _ = time.LoadLocation("Asia/Taipei")

// ShipDate 計算預計出貨日：t 之後的第 n 個工作天（跳過週六、週日）。
func ShipDate(t time.Time, n int) time.Time {
	for n > 0 {
		t = t.AddDate(0, 0, 1)
		if wd := t.Weekday(); wd != time.Saturday && wd != time.Sunday {
			n--
		}
	}
	return t
}

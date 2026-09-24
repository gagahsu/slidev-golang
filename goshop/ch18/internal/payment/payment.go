// Package payment 定義付款方式。
package payment

import (
	"errors"
	"fmt"

	"goshop/internal/money"
)

// ErrInsufficientFunds 表示錢包餘額不足。
var ErrInsufficientFunds = errors.New("餘額不足")

// Method 是所有付款方式都要具備的行為。
type Method interface {
	Name() string
	Pay(amount money.Money) error
}

// CreditCard 是信用卡，超過額度會被拒絕。
type CreditCard struct {
	Last4 string
	Limit money.Money
}

func (c CreditCard) Name() string { return "信用卡 *" + c.Last4 }

func (c CreditCard) Pay(amount money.Money) error {
	if amount > c.Limit {
		return fmt.Errorf("超過信用額度 %v", c.Limit)
	}
	return nil
}

// Wallet 是 GoShop 錢包；付款會改變餘額，所以用指標接收器。
type Wallet struct {
	Balance money.Money
}

func (w *Wallet) Name() string { return "GoShop 錢包" }

func (w *Wallet) Pay(amount money.Money) error {
	if w.Balance < amount {
		return fmt.Errorf("%w，只剩 %v", ErrInsufficientFunds, w.Balance)
	}
	w.Balance -= amount
	return nil
}

// CashOnDelivery 是貨到付款。
type CashOnDelivery struct{}

func (CashOnDelivery) Name() string                 { return "貨到付款" }
func (CashOnDelivery) Pay(amount money.Money) error { return nil }

var (
	_ Method = CreditCard{}
	_ Method = (*Wallet)(nil)
	_ Method = CashOnDelivery{}
)

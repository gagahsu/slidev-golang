package checkout

import (
	"errors"
	"fmt"

	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

// Service 負責結帳與付款流程。
type Service struct {
	Store *store.Memory
	Rules []Discount
}

// Checkout 把購物車轉成訂單，並扣除庫存。
func (s *Service) Checkout(cart Cart) (shop.Order, error) {
	if len(cart.Items) == 0 {
		return shop.Order{}, shop.ErrEmptyCart
	}
	var o shop.Order
	var errs []error
	for _, it := range cart.Items {
		p, err := s.Store.Product(it.SKU)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		o.Lines = append(o.Lines, shop.Line{
			SKU: p.SKU, Name: p.Name, Price: p.Price, Qty: it.Qty})
		o.Subtotal += p.Price.Times(it.Qty)
	}
	if err := errors.Join(errs...); err != nil {
		return shop.Order{}, fmt.Errorf("結帳失敗：%w", err)
	}
	o.Discount = Best(o.Subtotal, s.Rules...)
	o.Total = o.Subtotal - o.Discount
	if err := s.Store.PlaceOrder(&o); err != nil {
		return shop.Order{}, fmt.Errorf("結帳失敗：%w", err)
	}
	return o, nil
}

// Pay 依序嘗試付款方式，直到其中一種成功。
func (s *Service) Pay(id int, methods ...payment.Method) error {
	o, err := s.Store.Order(id)
	if err != nil {
		return err
	}
	if o.Status == shop.Paid {
		return fmt.Errorf("訂單 #%d：%w", id, shop.ErrPaid)
	}
	var errs []error
	for _, m := range methods {
		if err := m.Pay(o.Total); err != nil {
			errs = append(errs, fmt.Errorf("%s：%w", m.Name(), err))
			continue
		}
		return s.Store.MarkPaid(id, m.Name())
	}
	return fmt.Errorf("訂單 #%d 付款失敗：%w", id, errors.Join(errs...))
}

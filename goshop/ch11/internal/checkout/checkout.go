package checkout

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

// Service 負責結帳與付款流程。
type Service struct {
	Store   *store.Memory
	Rules   []Discount
	Coupons map[string]Coupon
	Now     func() time.Time // 取得現在時間；測試時可以換掉
}

func (s *Service) now() time.Time {
	if s.Now == nil {
		return time.Now()
	}
	return s.Now()
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
	o.CreatedAt = s.now().In(shop.Location)
	o.ShipBy = shop.ShipDate(o.CreatedAt, 2)

	rules := s.Rules
	if code := cart.Coupon; code != "" {
		c, ok := s.Coupons[code]
		if !ok || !c.Valid(o.CreatedAt) {
			return shop.Order{}, fmt.Errorf("%s：%w", code, shop.ErrCoupon)
		}
		o.Coupon = code
		rules = append(rules[:len(rules):len(rules)], PercentOff(c.PercentOff))
	}
	o.Discount = Best(o.Subtotal, rules...)
	o.Total = o.Subtotal - o.Discount
	if err := s.Store.PlaceOrder(&o); err != nil {
		slog.Warn("結帳失敗", "err", err)
		return shop.Order{}, fmt.Errorf("結帳失敗：%w", err)
	}
	slog.Info("訂單成立", "id", o.ID, "total", o.Total)
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
			slog.Debug("付款方式失敗", "id", id, "method", m.Name(), "err", err)
			errs = append(errs, fmt.Errorf("%s：%w", m.Name(), err))
			continue
		}
		slog.Info("訂單付款", "id", id, "method", m.Name())
		return s.Store.MarkPaid(id, m.Name())
	}
	return fmt.Errorf("訂單 #%d 付款失敗：%w", id, errors.Join(errs...))
}

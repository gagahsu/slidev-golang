package store

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	"goshop/internal/shop"
)

// Memory 把資料存在記憶體中，可以搭配 Save／Load 存成 gob 快照。
type Memory struct {
	products map[string]shop.Product
	orders   map[int]shop.Order
	lastID   int
}

// NewMemory 建立一個記憶體儲存庫，並放入初始商品。
func NewMemory(products ...shop.Product) *Memory {
	m := &Memory{
		products: make(map[string]shop.Product),
		orders:   make(map[int]shop.Order),
	}
	for _, p := range products {
		m.products[p.SKU] = p
	}
	return m
}

// Products 傳回依 SKU 排序的所有商品。
func (m *Memory) Products(ctx context.Context) ([]shop.Product, error) {
	return slices.SortedFunc(maps.Values(m.products),
		func(a, b shop.Product) int { return cmp.Compare(a.SKU, b.SKU) }), nil
}

// Product 用 SKU 查詢一項商品。
func (m *Memory) Product(ctx context.Context, sku string) (shop.Product, error) {
	p, ok := m.products[sku]
	if !ok {
		return shop.Product{}, fmt.Errorf("商品 %s：%w", sku, shop.ErrNotFound)
	}
	return p, nil
}

// SaveProduct 新增商品；SKU 已存在時就更新它。
func (m *Memory) SaveProduct(ctx context.Context, p shop.Product) error {
	m.products[p.SKU] = p
	return nil
}

// PlaceOrder 檢查並扣除庫存，替訂單編號後存起來。
func (m *Memory) PlaceOrder(ctx context.Context, o *shop.Order) error {
	var errs []error
	for _, l := range o.Lines {
		p, err := m.Product(ctx, l.SKU)
		if err != nil {
			errs = append(errs, err)
		} else if p.Stock < l.Qty {
			errs = append(errs, &shop.StockError{SKU: l.SKU, Want: l.Qty, Have: p.Stock})
		}
	}
	if err := errors.Join(errs...); err != nil {
		return err
	}
	for _, l := range o.Lines { // 全部檢查通過才扣庫存
		p := m.products[l.SKU]
		p.Stock -= l.Qty
		m.products[l.SKU] = p
	}
	m.lastID++
	o.ID = m.lastID
	m.orders[o.ID] = *o
	return nil
}

// Order 用編號查詢訂單。
func (m *Memory) Order(ctx context.Context, id int) (shop.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return shop.Order{}, fmt.Errorf("訂單 #%d：%w", id, shop.ErrNotFound)
	}
	return o, nil
}

// Orders 傳回依編號排序的所有訂單。
func (m *Memory) Orders(ctx context.Context) ([]shop.Order, error) {
	return slices.SortedFunc(maps.Values(m.orders),
		func(a, b shop.Order) int { return cmp.Compare(a.ID, b.ID) }), nil
}

// MarkPaid 把訂單標記為已付款。
func (m *Memory) MarkPaid(ctx context.Context, id int, by string) error {
	o, err := m.Order(ctx, id)
	if err != nil {
		return err
	}
	if o.Status == shop.Paid {
		return fmt.Errorf("訂單 #%d：%w", id, shop.ErrPaid)
	}
	o.Status, o.PaidBy = shop.Paid, by
	m.orders[id] = o
	return nil
}

// Package store 負責保存商品與訂單。
package store

import (
	"context"

	"goshop/internal/shop"
)

// Store 是 GoShop 的資料存取介面，Memory 和 MySQL 都實作它。
type Store interface {
	Products(ctx context.Context) ([]shop.Product, error)
	Product(ctx context.Context, sku string) (shop.Product, error)
	SaveProduct(ctx context.Context, p shop.Product) error
	PlaceOrder(ctx context.Context, o *shop.Order) error
	Order(ctx context.Context, id int) (shop.Order, error)
	Orders(ctx context.Context) ([]shop.Order, error)
	MarkPaid(ctx context.Context, id int, by string) error
}

var (
	_ Store = (*Memory)(nil)
	_ Store = (*MySQL)(nil)
)

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"goshop/internal/shop"
)

// PlaceOrder 在一個交易中扣庫存、寫入訂單與明細；任何一步失敗就全部復原。
func (s *MySQL) PlaceOrder(ctx context.Context, o *shop.Order) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Commit 成功之後再 Rollback 不會有任何作用

	var errs []error
	for _, l := range o.Lines {
		// 用 stock >= ? 當條件：庫存不夠就不會更新到任何一列
		res, err := tx.ExecContext(ctx, `UPDATE products SET stock = stock - ?
			WHERE sku = ? AND stock >= ?`, l.Qty, l.SKU, l.Qty)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			errs = append(errs, s.stockError(ctx, tx, l))
		}
	}
	if err := errors.Join(errs...); err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO orders
		(subtotal, discount, total, status, paid_by, coupon, created_at, ship_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Subtotal, o.Discount, o.Total, o.Status, nullString(o.PaidBy),
		nullString(o.Coupon), o.CreatedAt, o.ShipBy)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, l := range o.Lines {
		_, err := tx.ExecContext(ctx, `INSERT INTO order_lines
			(order_id, sku, name, price, qty) VALUES (?, ?, ?, ?, ?)`,
			id, l.SKU, l.Name, l.Price, l.Qty)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	o.ID = int(id)
	return nil
}

// stockError 查出目前的庫存，說明為什麼扣不了庫存。
func (s *MySQL) stockError(ctx context.Context, tx *sql.Tx, l shop.Line) error {
	var have int
	err := tx.QueryRowContext(ctx,
		"SELECT stock FROM products WHERE sku = ?", l.SKU).Scan(&have)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("商品 %s：%w", l.SKU, shop.ErrNotFound)
	}
	if err != nil {
		return err
	}
	return &shop.StockError{SKU: l.SKU, Want: l.Qty, Have: have}
}

// nullString 把空字串存成資料庫的 NULL。
func nullString(s string) sql.Null[string] {
	return sql.Null[string]{V: s, Valid: s != ""}
}

const orderColumns = `SELECT id, subtotal, discount, total, status,
	paid_by, coupon, created_at, ship_by FROM orders`

// scanner 是 *sql.Row 和 *sql.Rows 共同的 Scan 方法。
type scanner interface {
	Scan(dest ...any) error
}

func scanOrder(row scanner) (shop.Order, error) {
	var o shop.Order
	var paidBy, coupon sql.Null[string]
	err := row.Scan(&o.ID, &o.Subtotal, &o.Discount, &o.Total, &o.Status,
		&paidBy, &coupon, &o.CreatedAt, &o.ShipBy)
	o.PaidBy, o.Coupon = paidBy.V, coupon.V // NULL 會變成空字串
	o.CreatedAt = o.CreatedAt.In(shop.Location)
	o.ShipBy = o.ShipBy.In(shop.Location)
	return o, err
}

// Order 用編號查詢訂單與它的明細。
func (s *MySQL) Order(ctx context.Context, id int) (shop.Order, error) {
	o, err := scanOrder(s.db.QueryRowContext(ctx, orderColumns+" WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return o, fmt.Errorf("訂單 #%d：%w", id, shop.ErrNotFound)
	}
	if err != nil {
		return o, err
	}
	o.Lines, err = s.lines(ctx, id)
	return o, err
}

func (s *MySQL) lines(ctx context.Context, orderID int) ([]shop.Line, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT sku, name, price, qty
		FROM order_lines WHERE order_id = ? ORDER BY sku`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lines []shop.Line
	for rows.Next() {
		var l shop.Line
		if err := rows.Scan(&l.SKU, &l.Name, &l.Price, &l.Qty); err != nil {
			return nil, err
		}
		lines = append(lines, l)
	}
	return lines, rows.Err()
}

// Orders 傳回依編號排序的所有訂單。
func (s *MySQL) Orders(ctx context.Context) ([]shop.Order, error) {
	rows, err := s.db.QueryContext(ctx, orderColumns+" ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []shop.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range orders { // rows 讀完時已自動關閉，這時再查每張訂單的明細
		if orders[i].Lines, err = s.lines(ctx, orders[i].ID); err != nil {
			return nil, err
		}
	}
	return orders, nil
}

// MarkPaid 把待付款的訂單標記為已付款。
func (s *MySQL) MarkPaid(ctx context.Context, id int, by string) error {
	res, err := s.db.ExecContext(ctx,
		"UPDATE orders SET status = ?, paid_by = ? WHERE id = ? AND status = ?",
		shop.Paid, by, id, shop.Pending)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		if _, err := s.Order(ctx, id); err != nil {
			return err // 訂單不存在
		}
		return fmt.Errorf("訂單 #%d：%w", id, shop.ErrPaid)
	}
	return nil
}

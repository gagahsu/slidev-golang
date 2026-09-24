package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"goshop/internal/shop"
)

// schema 是 GoShop 需要的資料表；IF NOT EXISTS 讓它可以重複執行。
var schema = []string{
	`CREATE TABLE IF NOT EXISTS products (
		sku   VARCHAR(32)  PRIMARY KEY,
		name  VARCHAR(100) NOT NULL,
		price INT          NOT NULL,
		stock INT          NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS orders (
		id         INT AUTO_INCREMENT PRIMARY KEY,
		subtotal   INT         NOT NULL,
		discount   INT         NOT NULL,
		total      INT         NOT NULL,
		status     TINYINT     NOT NULL DEFAULT 0,
		paid_by    VARCHAR(50) NULL,
		coupon     VARCHAR(32) NULL,
		created_at DATETIME(6) NOT NULL,
		ship_by    DATETIME(6) NOT NULL
	)`,
	`CREATE TABLE IF NOT EXISTS order_lines (
		order_id INT          NOT NULL,
		sku      VARCHAR(32)  NOT NULL,
		name     VARCHAR(100) NOT NULL,
		price    INT          NOT NULL,
		qty      INT          NOT NULL,
		PRIMARY KEY (order_id, sku),
		FOREIGN KEY (order_id) REFERENCES orders (id)
	)`,
}

// MySQL 把資料存在 MySQL 資料庫。
type MySQL struct {
	db *sql.DB
}

// OpenMySQL 連線到資料庫，並建立需要的資料表。
func OpenMySQL(ctx context.Context, dsn string) (*MySQL, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	cfg.ParseTime = true // DATETIME 自動轉成 time.Time
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("無法連線到資料庫：%w", err)
	}
	for _, stmt := range schema {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			db.Close()
			return nil, fmt.Errorf("建立資料表：%w", err)
		}
	}
	return &MySQL{db: db}, nil
}

// Close 關閉資料庫連線池。
func (s *MySQL) Close() error { return s.db.Close() }

// Products 傳回依 SKU 排序的所有商品。
func (s *MySQL) Products(ctx context.Context) ([]shop.Product, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT sku, name, price, stock FROM products ORDER BY sku")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []shop.Product
	for rows.Next() {
		var p shop.Product
		if err := rows.Scan(&p.SKU, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// Product 用 SKU 查詢一項商品。
func (s *MySQL) Product(ctx context.Context, sku string) (shop.Product, error) {
	var p shop.Product
	err := s.db.QueryRowContext(ctx,
		"SELECT sku, name, price, stock FROM products WHERE sku = ?", sku,
	).Scan(&p.SKU, &p.Name, &p.Price, &p.Stock)
	if errors.Is(err, sql.ErrNoRows) {
		return p, fmt.Errorf("商品 %s：%w", sku, shop.ErrNotFound)
	}
	return p, err
}

// SaveProduct 新增商品；SKU 已存在時就更新它。
func (s *MySQL) SaveProduct(ctx context.Context, p shop.Product) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO products (sku, name, price, stock)
		VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		name = VALUES(name), price = VALUES(price), stock = VALUES(stock)`,
		p.SKU, p.Name, p.Price, p.Stock)
	return err
}

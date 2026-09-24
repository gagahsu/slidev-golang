package store

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"goshop/internal/money"
	"goshop/internal/shop"
)

// ReadProductsCSV 讀取「sku,name,price,stock」格式的 CSV，第一行是標題。
// 格式錯誤的行會被略過，並一起回報在 error 中。
func ReadProductsCSV(r io.Reader) ([]shop.Product, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 4 // 每一行都必須剛好 4 個欄位
	records, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("讀取 CSV：%w", err)
	}

	var products []shop.Product
	var errs []error
	for i, rec := range records {
		if i == 0 {
			continue // 跳過標題列
		}
		price, err1 := strconv.Atoi(rec[2])
		stock, err2 := strconv.Atoi(rec[3])
		if err := errors.Join(err1, err2); err != nil {
			errs = append(errs, fmt.Errorf("第 %d 行：%w", i+1, err))
			continue
		}
		products = append(products, shop.Product{
			SKU: rec[0], Name: rec[1], Price: money.Money(price), Stock: stock})
	}
	return products, errors.Join(errs...)
}

// WriteOrdersCSV 把訂單寫成 CSV，方便用試算表軟體開啟。
func WriteOrdersCSV(w io.Writer, orders []shop.Order) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"id", "created_at", "items", "total", "status", "paid_by"})
	for _, o := range orders {
		cw.Write([]string{
			strconv.Itoa(o.ID),
			o.CreatedAt.Format(time.DateTime),
			strconv.Itoa(len(o.Lines)),
			strconv.Itoa(int(o.Total)),
			o.Status.String(),
			o.PaidBy,
		})
	}
	cw.Flush() // 把緩衝區的資料真正寫出去
	return cw.Error()
}

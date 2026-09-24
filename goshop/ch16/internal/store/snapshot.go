package store

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"

	"goshop/internal/shop"
)

// DecodeProducts 從 JSON 陣列讀取商品清單；不認得的欄位視為錯誤。
func DecodeProducts(r io.Reader) ([]shop.Product, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var products []shop.Product
	if err := dec.Decode(&products); err != nil {
		return nil, fmt.Errorf("解析商品 JSON：%w", err)
	}
	return products, nil
}

// snapshot 是寫進 gob 檔的完整資料；欄位要匯出 gob 才看得到。
type snapshot struct {
	Products map[string]shop.Product
	Orders   map[int]shop.Order
	LastID   int
}

// Save 把所有商品與訂單用 gob 格式寫到 w。
func (m *Memory) Save(w io.Writer) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	snap := snapshot{m.products, m.orders, m.lastID}
	return gob.NewEncoder(w).Encode(snap)
}

// Load 從 gob 格式讀回所有商品與訂單。
func (m *Memory) Load(r io.Reader) error {
	var snap snapshot
	if err := gob.NewDecoder(r).Decode(&snap); err != nil {
		return fmt.Errorf("讀取快照：%w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.products, m.orders, m.lastID = snap.Products, snap.Orders, snap.LastID
	if m.orders == nil { // gob 不會傳送空的 map
		m.orders = make(map[int]shop.Order)
	}
	return nil
}

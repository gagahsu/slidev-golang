// Package checkout 處理購物車、折扣與結帳流程。
package checkout

// Item 是購物車裡的一個品項。
type Item struct {
	SKU string `json:"sku" validate:"required"`
	Qty int    `json:"qty" validate:"min=1,max=99"`
}

// Cart 是購物車。
type Cart struct {
	Items  []Item `json:"items" validate:"min=1,max=20"`
	Coupon string `json:"coupon,omitzero"` // 折價券代碼，可以空白
}

// Add 把品項放進購物車；同一個 SKU 會合併數量。
func (c *Cart) Add(items ...Item) {
	for _, item := range items {
		c.add(item)
	}
}

func (c *Cart) add(item Item) {
	for i := range c.Items {
		if c.Items[i].SKU == item.SKU {
			c.Items[i].Qty += item.Qty
			return
		}
	}
	c.Items = append(c.Items, item)
}

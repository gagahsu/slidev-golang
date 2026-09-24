// Package webhook 用 HTTP POST 把事件通知其他系統（例如倉庫出貨系統）。
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Event 是送出去的通知內容。
type Event struct {
	Type string `json:"type"` // 事件種類，例如 order.paid
	Data any    `json:"data"` // 事件資料，例如訂單
}

// Send 把事件編碼成 JSON，POST 到 url；回應不是 2xx 就視為失敗。
func Send(ctx context.Context, c *http.Client, url string, e Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url,
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("送出 %s：%w", e.Type, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("送出 %s：對方回應 %s", e.Type, resp.Status)
	}
	return nil
}

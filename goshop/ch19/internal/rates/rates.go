// Package rates 向外部服務查詢外幣匯率。
package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// DefaultURL 是免費的匯率服務 ExchangeRate-API。
const DefaultURL = "https://open.er-api.com/v6"

// Client 是匯率服務的客戶端。
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New 建立一個 5 秒逾時的客戶端；baseURL 空白時使用 DefaultURL。
func New(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{Timeout: 5 * time.Second},
	}
}

// latest 對應服務回傳的 JSON，只取出需要的欄位。
type latest struct {
	Result string             `json:"result"`
	Rates  map[string]float64 `json:"rates"`
}

// Rate 查詢 1 單位的 from 可以換成多少 to，例如 TWD → USD。
func (c *Client) Rate(ctx context.Context, from, to string) (float64, error) {
	u := c.BaseURL + "/latest/" + url.PathEscape(from)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GoShop/1.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, fmt.Errorf("查詢匯率：%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("查詢匯率：伺服器回應 %s", resp.Status)
	}

	var data latest
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("解析匯率：%w", err)
	}
	rate, ok := data.Rates[to]
	if data.Result != "success" || !ok {
		return 0, fmt.Errorf("查不到 %s → %s 的匯率", from, to)
	}
	return rate, nil
}

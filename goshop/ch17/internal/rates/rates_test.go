package rates

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// httptest.NewServer 在本機啟動一個假的匯率服務，測試不需要連上網路
func TestRate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/latest/TWD" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			fmt.Fprint(w, `{"result":"success","rates":{"USD":0.03125}}`)
		}))
	defer srv.Close()

	c := New(srv.URL)
	got, err := c.Rate(t.Context(), "TWD", "USD")
	if err != nil || got != 0.03125 {
		t.Errorf("Rate(TWD, USD) = %v, %v", got, err)
	}
	if _, err := c.Rate(t.Context(), "TWD", "XYZ"); err == nil {
		t.Error("不存在的幣別應該回傳錯誤")
	}
	if _, err := c.Rate(t.Context(), "JPY", "USD"); err == nil {
		t.Error("404 應該回傳錯誤")
	}
}

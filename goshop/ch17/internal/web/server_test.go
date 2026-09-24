package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"goshop/internal/checkout"
	"goshop/internal/shop"
	"goshop/internal/store"
)

func newTestServer() http.Handler {
	st := store.NewMemory(
		shop.Product{SKU: "A", Name: "咖啡豆", Price: 450, Stock: 10},
		shop.Product{SKU: "B", Name: "手沖壺", Price: 1280, Stock: 1},
	)
	s := &Server{Store: st, Checkout: &checkout.Service{Store: st}}
	return s.Handler()
}

// do 送出一個請求，傳回回應的狀態碼與內容
func do(h http.Handler, method, target, body string) (int, string) {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if strings.HasPrefix(target, "/admin/") {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestAPI(t *testing.T) {
	h := newTestServer()
	tests := []struct {
		method, target, body string
		wantCode             int
		wantBody             string
	}{
		{"GET", "/api/products", "", 200, `"sku":"A"`},
		{"GET", "/api/products/B", "", 200, `"name":"手沖壺"`},
		{"GET", "/api/products/Z", "", 404, "找不到"},
		{"POST", "/api/orders", `{"items":[{"sku":"A","qty":2}]}`, 201, `"id":1`},
		{"POST", "/api/orders", `{"items":[{"sku":"B","qty":5}]}`, 409, "庫存不足"},
		{"POST", "/api/orders", `{"items":[]}`, 400, "購物車是空的"},
		{"POST", "/api/orders", `{"itemz":[]}`, 400, "JSON 格式錯誤"},
		{"GET", "/api/orders/1", "", 200, `"status":"pending"`},
		{"GET", "/api/orders/abc", "", 400, "正整數"},
		{"POST", "/api/orders/1/pay", `{"method":"cod"}`, 200, `"status":"paid"`},
		{"POST", "/api/orders/1/pay", `{"method":"cod"}`, 409, "已付款"},
		{"POST", "/api/orders/9/pay", `{"method":"cod"}`, 404, "找不到"},
		{"DELETE", "/api/products/A", "", 405, ""},
	}
	for _, tt := range tests {
		code, body := do(h, tt.method, tt.target, tt.body)
		if code != tt.wantCode || !strings.Contains(body, tt.wantBody) {
			t.Errorf("%s %s = %d %s，want %d …%s…",
				tt.method, tt.target, code, body, tt.wantCode, tt.wantBody)
		}
	}
}

func TestAdmin(t *testing.T) {
	h := newTestServer()
	form := url.Values{"sku": {"c"}, "name": {"馬克杯"}, "price": {"350"}, "stock": {"3"}}
	if code, _ := do(h, "POST", "/admin/products", form.Encode()); code != http.StatusSeeOther {
		t.Fatalf("新增商品狀態碼 = %d，want 303", code)
	}
	code, body := do(h, "GET", "/admin", "")
	if code != 200 || !strings.Contains(body, "馬克杯") || !strings.Contains(body, `class="low"`) {
		t.Errorf("GET /admin = %d\n%s", code, body)
	}
	form.Set("price", "abc")
	if code, _ := do(h, "POST", "/admin/products", form.Encode()); code != 400 {
		t.Errorf("價格錯誤時狀態碼 = %d，want 400", code)
	}

	code, body = do(h, "GET", "/api/products/C", "")
	var p shop.Product
	json.Unmarshal([]byte(body), &p)
	if code != 200 || p.Price != 350 {
		t.Errorf("GET /api/products/C = %d %+v", code, p)
	}
	if code, _ := do(h, "GET", "/static/style.css", ""); code != 200 {
		t.Errorf("GET /static/style.css = %d", code)
	}
}

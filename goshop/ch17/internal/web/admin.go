package web

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"goshop/internal/money"
	"goshop/internal/shop"
)

// adminData 是後台頁面模板需要的資料。
type adminData struct {
	Products []shop.Product
	Orders   []shop.Order
	Error    string
}

func (s *Server) adminPage(w http.ResponseWriter, r *http.Request) {
	s.renderAdmin(w, r, http.StatusOK, "")
}

func (s *Server) renderAdmin(w http.ResponseWriter, r *http.Request,
	status int, msg string) {
	products, err1 := s.Store.Products(r.Context())
	orders, err2 := s.Store.Orders(r.Context())
	if err1 != nil || err2 != nil {
		http.Error(w, "讀取資料失敗", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	data := adminData{Products: products, Orders: orders, Error: msg}
	if err := tmpl.ExecuteTemplate(w, "admin.html", data); err != nil {
		slog.Error("產生網頁失敗", "err", err)
	}
}

// saveProduct 處理後台的「新增／更新商品」表單。
func (s *Server) saveProduct(w http.ResponseWriter, r *http.Request) {
	price, err1 := strconv.Atoi(r.FormValue("price"))
	stock, err2 := strconv.Atoi(r.FormValue("stock"))
	p := shop.Product{
		SKU:   strings.ToUpper(strings.TrimSpace(r.FormValue("sku"))),
		Name:  strings.TrimSpace(r.FormValue("name")),
		Price: money.Money(price),
		Stock: stock,
	}
	if p.SKU == "" || p.Name == "" || err1 != nil || err2 != nil ||
		price <= 0 || stock < 0 {
		s.renderAdmin(w, r, http.StatusBadRequest, "請填寫完整的商品資料，價格要大於 0")
		return
	}
	if err := s.Store.SaveProduct(r.Context(), p); err != nil {
		s.renderAdmin(w, r, http.StatusInternalServerError, "儲存失敗")
		return
	}
	// 成功後重新導向（PRG 模式），重新整理頁面才不會重複送出表單
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

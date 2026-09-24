// Package web 提供 GoShop 的 RESTful JSON API 與後台管理網頁。
package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"goshop/internal/auth"
	"goshop/internal/checkout"
	"goshop/internal/store"
)

//go:embed templates static
var assets embed.FS

// 啟動時就解析模板；模板有語法錯誤時程式會直接 panic
var tmpl = template.Must(template.ParseFS(assets, "templates/*.html"))

// Server 保存處理請求時需要的相依物件。
type Server struct {
	Store    store.Store
	Checkout *checkout.Service
	Auth     *auth.Auth
}

// Handler 設定所有路由，傳回加上日誌中介軟體的 http.Handler。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/products", s.listProducts)
	mux.HandleFunc("GET /api/products/{sku}", s.getProduct)
	mux.HandleFunc("POST /api/orders", s.createOrder)
	mux.HandleFunc("GET /api/orders/{id}", s.getOrder)
	mux.HandleFunc("POST /api/orders/{id}/pay", s.payOrder)

	mux.HandleFunc("GET /login", s.loginPage)
	mux.HandleFunc("POST /login", s.login)
	mux.HandleFunc("POST /logout", s.logout)
	mux.HandleFunc("GET /admin", s.requireAdmin(s.adminPage))
	mux.HandleFunc("POST /admin/products", s.requireAdmin(s.saveProduct))
	static, _ := fs.Sub(assets, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.Handle("GET /{$}", http.RedirectHandler("/admin", http.StatusFound))
	return logging(mux)
}

// statusRecorder 記下處理器寫出的狀態碼。
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// logging 是記錄每個請求的中介軟體。
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path,
			"status", rec.status, "duration", time.Since(start))
	})
}

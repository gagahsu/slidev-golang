package web

import (
	"log/slog"
	"net/http"
	"time"
)

const cookieName = "goshop_admin"

// requireAdmin 是中介軟體：沒有有效的登入憑證就導向登入頁。
func (s *Server) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil || s.Auth.Verify(c.Value, time.Now()) != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	s.renderLogin(w, http.StatusOK, "")
}

func (s *Server) renderLogin(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "login.html", msg); err != nil {
		slog.Error("產生網頁失敗", "err", err)
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.Auth.CheckPassword(r.FormValue("password")) {
		slog.Warn("後台登入失敗", "ip", r.RemoteAddr)
		s.renderLogin(w, http.StatusUnauthorized, "密碼錯誤")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    s.Auth.NewToken(time.Now()),
		Path:     "/",
		MaxAge:   int(s.Auth.TTL.Seconds()),
		HttpOnly: true,         // JavaScript 讀不到，降低 XSS 竊取的風險
		Secure:   r.TLS != nil, // 使用 HTTPS 時，只在加密連線中傳送
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

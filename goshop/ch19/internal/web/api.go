package web

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"goshop/internal/checkout"
	"goshop/internal/money"
	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/validate"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("寫出 JSON 失敗", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// writeErr 依錯誤種類決定 HTTP 狀態碼。
func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shop.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, shop.ErrOutOfStock), errors.Is(err, shop.ErrPaid):
		writeError(w, http.StatusConflict, err.Error())
	case errors.As(err, new(*validate.FieldError)),
		errors.Is(err, shop.ErrEmptyCart), errors.Is(err, shop.ErrCoupon),
		errors.Is(err, payment.ErrInsufficientFunds):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		slog.Error("內部錯誤", "err", err) // 細節只寫進日誌，不回給客戶端
		writeError(w, http.StatusInternalServerError, "伺服器發生錯誤")
	}
}

func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.Store.Products(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	p, err := s.Store.Product(r.Context(), r.PathValue("sku"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	var cart checkout.Cart
	// 請求內容最多 1 MB，避免超大的請求吃光記憶體
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cart); err != nil {
		writeError(w, http.StatusBadRequest, "JSON 格式錯誤："+err.Error())
		return
	}
	if err := validate.Struct(cart); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.Checkout.Checkout(r.Context(), cart)
	if err != nil {
		writeErr(w, err)
		return
	}
	w.Header().Set("Location", "/api/orders/"+strconv.Itoa(o.ID))
	writeJSON(w, http.StatusCreated, o)
}

// orderID 從路徑取出訂單編號。
func orderID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return id, err == nil && id > 0
}

func (s *Server) getOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := orderID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "訂單編號必須是正整數")
		return
	}
	o, err := s.Store.Order(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

// payRequest 是付款 API 的請求內容。
type payRequest struct {
	Method string `json:"method"` // cod（貨到付款）或 card（信用卡）
	Last4  string `json:"last4"`  // 信用卡末 4 碼
}

func (s *Server) payOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := orderID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "訂單編號必須是正整數")
		return
	}
	var req payRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON 格式錯誤："+err.Error())
		return
	}
	var m payment.Method
	switch req.Method {
	case "cod":
		m = payment.CashOnDelivery{}
	case "card":
		m = payment.CreditCard{Last4: req.Last4, Limit: money.Money(50000)}
	default:
		writeError(w, http.StatusBadRequest, "不支援的付款方式："+req.Method)
		return
	}
	if err := s.Checkout.Pay(r.Context(), id, m); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.Store.Order(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

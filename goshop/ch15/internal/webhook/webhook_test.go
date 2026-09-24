package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSend(t *testing.T) {
	var got Event
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/fail" {
				http.Error(w, "倉庫系統維護中", http.StatusServiceUnavailable)
				return
			}
			if r.Method != http.MethodPost ||
				r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			json.NewDecoder(r.Body).Decode(&got)
			w.WriteHeader(http.StatusNoContent)
		}))
	defer srv.Close()

	err := Send(t.Context(), srv.Client(), srv.URL, Event{"order.paid", 42})
	if err != nil || got.Type != "order.paid" || got.Data != 42.0 {
		t.Errorf("Send() err = %v，對方收到 %+v", err, got)
	}
	err = Send(t.Context(), srv.Client(), srv.URL+"/fail", Event{"x", nil})
	if err == nil {
		t.Error("對方回應 503 時應該回傳錯誤")
	}
}

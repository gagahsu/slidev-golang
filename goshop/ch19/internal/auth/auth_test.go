package auth

import (
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	h1, _ := HashPassword("s3cret!")
	h2, _ := HashPassword("s3cret!")
	if h1 == h2 {
		t.Error("相同密碼的雜湊應該不同（有加鹽）")
	}
	a := &Auth{PasswordHash: []byte(h1)}
	if !a.CheckPassword("s3cret!") || a.CheckPassword("S3cret!") {
		t.Error("CheckPassword 結果錯誤")
	}
}

func TestToken(t *testing.T) {
	now := time.Date(2026, 9, 25, 10, 0, 0, 0, time.UTC)
	a := &Auth{Secret: []byte("0123456789abcdef0123456789abcdef"), TTL: time.Hour}
	token := a.NewToken(now)

	tests := []struct {
		name  string
		token string
		at    time.Time
		ok    bool
	}{
		{"有效", token, now.Add(59 * time.Minute), true},
		{"過期", token, now.Add(time.Hour), false},
		{"竄改到期時間", "9999999999" + token[10:], now, false},
		{"格式錯誤", "garbage", now, false},
	}
	for _, tt := range tests {
		if err := a.Verify(tt.token, tt.at); (err == nil) != tt.ok {
			t.Errorf("%s：Verify() err = %v", tt.name, err)
		}
	}
	other := &Auth{Secret: []byte("another-secret-another-secret-00"), TTL: time.Hour}
	if other.Verify(token, now) == nil {
		t.Error("用不同金鑰簽的憑證應該無效")
	}
}

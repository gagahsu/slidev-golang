// Package auth 處理後台登入：用 bcrypt 比對密碼，用 HMAC 簽署登入憑證。
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidToken 表示憑證被竄改、格式錯誤或已經過期。
var ErrInvalidToken = errors.New("登入憑證無效或已過期")

// Auth 保存管理員的密碼雜湊與簽章用的金鑰。
type Auth struct {
	PasswordHash []byte        // bcrypt 雜湊，不是明文密碼
	Secret       []byte        // HMAC 金鑰，至少 32 位元組
	TTL          time.Duration // 登入的有效時間
}

// HashPassword 產生 bcrypt 雜湊；每次結果都不同，因為會自動加鹽。
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h), err
}

// CheckPassword 比對密碼是否正確。
func (a *Auth) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(password))
	return err == nil
}

// NewToken 產生「到期時間.簽章」格式的憑證，例如 1790000000.q7Xk…
func (a *Auth) NewToken(now time.Time) string {
	exp := strconv.FormatInt(now.Add(a.TTL).Unix(), 10)
	return exp + "." + a.sign(exp)
}

// Verify 檢查簽章是否正確、是否還沒過期。
func (a *Auth) Verify(token string, now time.Time) error {
	exp, sig, ok := strings.Cut(token, ".")
	// hmac.Equal 花費的時間固定，避免被用「計時攻擊」猜出簽章
	if !ok || !hmac.Equal([]byte(sig), []byte(a.sign(exp))) {
		return ErrInvalidToken
	}
	unix, err := strconv.ParseInt(exp, 10, 64)
	if err != nil || now.Unix() >= unix {
		return ErrInvalidToken
	}
	return nil
}

func (a *Auth) sign(msg string) string {
	mac := hmac.New(sha256.New, a.Secret)
	mac.Write([]byte(msg))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

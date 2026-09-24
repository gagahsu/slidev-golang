package checkout

import "time"

// Coupon 是有使用期限的折價券。
type Coupon struct {
	Code       string
	PercentOff int       // 折扣百分比，例如 15 代表 85 折
	Start      time.Time // 生效時間（包含）
	End        time.Time // 失效時間（不包含）
}

// Valid 判斷時間 t 是否在有效期間 [Start, End) 內。
func (c Coupon) Valid(t time.Time) bool {
	return !t.Before(c.Start) && t.Before(c.End)
}

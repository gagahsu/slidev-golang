// Package money 處理新台幣金額。
package money

import "strconv"

// Money 是新台幣金額，單位是「元」。
type Money int

// Times 傳回單價乘上數量的金額。
func (m Money) Times(qty int) Money {
	return m * Money(qty)
}

// String 印出含千分位的金額，例如 NT$1,234、-NT$300。
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign, m = "-", -m // 先把負號拿掉，最後再補回去
	}
	s := strconv.Itoa(int(m))
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return sign + "NT$" + s
}

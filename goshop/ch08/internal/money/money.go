// Package money 處理新台幣金額。
package money

import "strconv"

// Money 是新台幣金額，單位是「元」。
type Money int

// Times 傳回單價乘上數量的金額。
func (m Money) Times(qty int) Money {
	return m * Money(qty)
}

// String 印出含千分位的金額，例如 NT$1,234。
func (m Money) String() string {
	s := strconv.Itoa(int(m))
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return "NT$" + s
}

package checkout

import "goshop/internal/money"

// Discount 傳入小計，傳回可以折抵的金額。
type Discount func(subtotal money.Money) money.Money

// PercentOff 建立「折 p%」的規則，例如 10 代表 9 折。
func PercentOff(p int) Discount {
	return func(subtotal money.Money) money.Money {
		return subtotal * money.Money(p) / 100
	}
}

// Threshold 建立「滿 limit 元折 off 元」的規則。
func Threshold(limit, off money.Money) Discount {
	return func(subtotal money.Money) money.Money {
		if subtotal >= limit {
			return off
		}
		return 0
	}
}

// Best 從多個規則中挑出折最多的金額。
func Best(subtotal money.Money, rules ...Discount) money.Money {
	var best money.Money
	for _, rule := range rules {
		best = max(best, rule(subtotal))
	}
	return best
}

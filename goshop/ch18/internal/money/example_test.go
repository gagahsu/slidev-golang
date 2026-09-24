package money_test

import (
	"fmt"

	"goshop/internal/money"
)

func ExampleMoney_String() {
	price := money.Money(1280)
	fmt.Println(price)
	fmt.Println(price.Times(3))
	fmt.Println(money.Money(-300))
	// Output:
	// NT$1,280
	// NT$3,840
	// -NT$300
}

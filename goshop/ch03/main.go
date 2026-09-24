package main

import (
	"fmt"
	"strconv"
	"strings"
)

// 從供應商收到的商品資料：SKU | 品名 | 單價 | 庫存
const data = `
sku-001 | 衣索比亞咖啡豆 | 450 | 20
sku-002 | 濾掛咖啡（10入） | 280 | 50
 sku-003| 手沖壺 | 1280 | 5
sku-004 | 馬克杯 | 三百五 | 30
`

// width 計算字串顯示時佔幾格：中文等全形字佔 2 格
func width(s string) int {
	w := 0
	for _, r := range s {
		if r < 128 {
			w++
		} else {
			w += 2
		}
	}
	return w
}

func main() {
	var sb strings.Builder
	count, value := 0, 0

	lines := strings.Split(strings.TrimSpace(data), "\n")
	for _, line := range lines {
		fields := strings.Split(line, "|")
		sku := strings.ToUpper(strings.TrimSpace(fields[0]))
		name := strings.TrimSpace(fields[1])
		price, err := strconv.Atoi(strings.TrimSpace(fields[2]))
		if err != nil {
			fmt.Println("略過價格格式錯誤的商品：", sku)
			continue
		}
		stock, _ := strconv.Atoi(strings.TrimSpace(fields[3]))

		count++
		value += price * stock
		pad := strings.Repeat(" ", 20-width(name))
		fmt.Fprintf(&sb, "%s  %s%s%5d 元\n", sku, name, pad, price)
	}

	fmt.Print(sb.String())
	fmt.Println("商品數：", count, "／庫存總值：", value, "元")
	avg := float64(value) / float64(count)
	fmt.Printf("平均每種商品庫存價值：%.1f 元\n", avg)
}

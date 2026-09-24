---
theme: penguin
class: text-center
highlighter: shiki
lineNumbers: true
drawings:
  persist: false
transition: slide-left
fonts:
  provider: none
title: 條件判斷與迴圈
routeAlias: ch02
style: |
  .slidev-layout p,
  .slidev-layout li,
  .slidev-layout td,
  .slidev-layout th,
  .slidev-layout div {
    font-size: max(16px, 1em);
  }
  table {
    width: 100%;
    margin: 1rem 0;
    border-collapse: collapse;
  }
  th, td {
    padding: 8px !important;
    border: 1px solid #e2e8f0 !important;
  }
  .index-table td {
    text-align: center;
    font-family: monospace;
  }
---

<div class="flex flex-col justify-center items-center h-full" style="background: #ffffff;">
  <p style="color: #5eada0; font-size: 1rem; font-weight: 600; letter-spacing: 0.2em; text-transform: uppercase; margin-bottom: 1.2rem;">Go Programming Masterclass</p>
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">條件判斷與迴圈</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「讓程式會做決定，也會重複做事」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第二章！

到目前為止，我們寫的程式都是從第一行一路執行到最後一行，像一條沒有岔路的直線。但真實世界的程式需要「做決定」：會員要打折、非會員不打折；也需要「重複做事」：把購物車裡的每一件商品價格加起來。

今天要學的 if、switch、for，就是讓程式能夠做決定和重複做事的工具。Go 在這部分的設計特別精簡：只有一種迴圈 for，但它能做到其他語言 while、do-while、foreach 做的所有事。
-->

---
layout: default
---

# Outline

- **回顧與前言** — 什麼是流程控制
- **if 敘述** — `if`、`else`、`else if`、起始賦值
- **switch 敘述** — 基本用法、無條件 switch、`fallthrough`
- **迴圈** — `for` 的三種型態、`for range`、`break` 與 `continue`
- **章節總結**

<!--
今天分成三大塊：if、switch、for。

if 和 switch 負責「做決定」，for 負責「重複做事」。學完之後，我們就能寫出像 FizzBuzz、成績等級判斷、九九乘法表這些經典的練習題。
-->

---

# 回顧：變數與算符

- 變數用 `var` 或 `:=` 宣告；`:=` 只能在函式內使用
- **比較算符**（`==`、`<`、`>=`…）的結果是 `bool`
- **邏輯算符** `&&`、`||` 可以組合多個條件，而且會**短路求值**
- 變數的作用範圍是它所在的 `{ }` 區塊；小心 `:=` 造成變數遮蔽

<!--
回顧一下上一章的重點。

我們學了怎麼宣告變數，也學了比較算符和邏輯算符，它們算出來的結果都是布林值 true 或 false。

這一章的 if 和 for，最需要的就是這些布林值：「如果這個條件是 true，就做某件事」、「只要這個條件是 true，就一直重複」。另外作用範圍的觀念，今天也會在 if 和 for 的區塊裡再次碰到。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# if 敘述
## Conditional Statements

<!--
第一個主題：if 敘述，讓程式根據條件決定要不要執行某段程式碼。
-->

---

# 什麼是 if 敘述？

「**if 敘述會檢查一個布林條件，條件為 `true` 時才執行大括號內的程式碼。**」

```go
package main

import "fmt"

func main() {
	temperature := 32
	if temperature > 30 {
		fmt.Println("好熱，開冷氣！")
	}
	fmt.Println("程式結束")
}
```

<!--
什麼是 if？就像我們出門前看天氣：「如果下雨，就帶傘」。條件成立就做，不成立就跳過。

執行後，因為 32 大於 30，會先印出「好熱，開冷氣！」，再印出「程式結束」。
-->

---

# Go 的 if 規則

| 規則 | 說明 |
| --- | --- |
| 條件**不用**小括號 | 寫 `if x > 0`，不寫 `if (x > 0)`（寫了 `gofmt` 也會拿掉） |
| 大括號**一定要寫** | 即使只有一行也不能省略 |
| 條件必須是 `bool` | `if 1 { }` 是編譯錯誤，不會把數字當成真假 |

```go
if (x > 0) return x    // 編譯錯誤：缺少大括號
if x > 0 { return x }  // ✅ 存檔後 gofmt 會自動展開成多行
```

<!--
Go 的 if 有三個跟其他語言不太一樣的規則。

第一，條件不需要小括號，寫了 gofmt 也會幫我們拿掉。第二，大括號一定要寫，就算只有一行也不能省略，這避免了很多「以為在 if 裡面、其實不在」的 bug。第三，條件一定要是布林值，不能像 C 或 JavaScript 一樣把 1 當成 true。
-->

---

# else 敘述

條件不成立時，執行 `else` 區塊：

```go
package main

import "fmt"

func main() {
	stock := 0
	if stock > 0 {
		fmt.Println("有庫存，可以下單")
	} else {
		fmt.Println("已售完")
	}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>格式規定：</b> <code>else</code> 必須和前一個 <code>}</code> 寫在<b>同一行</b>。因為 Go 會在行尾自動補分號，<code>}</code> 後面換行再寫 <code>else</code> 會編譯錯誤。
</div>

<!--
if 只能處理「條件成立」的情況，加上 else 就能處理「條件不成立」的情況，變成二選一。

這裡有一個 Go 特有的規定：else 一定要跟前面的右大括號寫在同一行。原因是 Go 編譯器會在行尾自動補分號，如果右大括號後面換行，if 敘述就被當成已經結束了，下一行單獨出現的 else 就變成語法錯誤。

不過大家不用特別記，存檔時 gofmt 會幫我們排好。
-->

---

# else if 敘述

有多個條件時，用 `else if` 串接，**由上往下檢查，第一個成立的條件就執行**：

```go
package main

import "fmt"

func main() {
	score := 78
	if score >= 90 {
		fmt.Println("A")
	} else if score >= 80 {
		fmt.Println("B")
	} else if score >= 70 {
		fmt.Println("C") // ✅ 印出 C
	} else {
		fmt.Println("需要加油")
	}
}
```

<!--
當選項超過兩個，就用 else if 串起來。

重點是「由上往下，第一個成立的就執行，其他都跳過」。78 分不大於等於 90，也不大於等於 80，到了大於等於 70 成立，印出 C，後面的 else 就不會再檢查。

所以條件的順序很重要：如果把 score >= 70 放在最上面，那 95 分也會被判斷成 C。寫範圍判斷的時候，要從最嚴格的條件開始寫。
-->

---

# if 敘述的起始賦值

`if` 可以在條件前面先執行一個**簡短敘述**，用分號隔開：

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	if n, err := strconv.Atoi("42"); err != nil {
		fmt.Println("轉換失敗：", err)
	} else {
		fmt.Println("轉換成功：", n*2) // 轉換成功： 84
	}
	// fmt.Println(n) // 編譯錯誤：n 只存在於 if/else 區塊內
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>好處：</b> 起始賦值宣告的變數，作用範圍只在 <code>if</code>／<code>else</code> 區塊內，用完就消失，不會汙染外面的程式碼。
</div>

<!--
這是 Go 很有特色的寫法：在 if 的條件前面，可以先寫一個簡短的敘述，通常是宣告變數，用分號隔開。

strconv.Atoi 會把字串轉成整數，它回傳兩個值：轉換結果和錯誤。我們在 if 的起始賦值裡接住這兩個值，然後馬上檢查 err 是不是 nil。

這個寫法的好處是：n 和 err 只活在 if 跟 else 的區塊裡，出了區塊就消失了。這讓程式碼更乾淨，也避免變數被誤用。

「if 起始賦值 + 檢查 err」是 Go 最常見的寫法之一，第 6 章錯誤處理會大量使用。
-->

---

# 使用 if 的注意事項：提早 return

Go 社群偏好**先處理例外情況並提早離開**，讓主要邏輯不要縮排太深：

```go
package main

import "fmt"

func discount(age int, isMember bool) float64 {
	if age < 0 {
		return 0 // 不合理的輸入，提早離開
	}
	if !isMember {
		return 1.0
	}
	return 0.8 // 主要邏輯保持在最外層
}

func main() {
	fmt.Println(discount(30, true), discount(30, false)) // 0.8 1
}
```

<!--
使用 if 的注意事項：Go 社群有一個很重要的慣例，叫做「提早 return」，也常被稱為 guard clause（守衛子句）。

意思是：先把不正常的情況處理掉、直接離開函式，這樣主要的邏輯就可以寫在最外層，不用一層一層的 else 包起來。

想像在機場的安檢：不合規定的旅客在門口就被擋下來，能進去的都是合格的旅客，裡面的流程就單純很多。

這個慣例在 Go 的程式碼裡隨處可見，特別是錯誤處理：if err != nil { return err }，錯誤先處理掉，正常流程往下走。
-->

---
layout: default
---

# 練習 1：運費計算
### 任務說明

寫一個程式，依照訂單金額 `amount` 決定運費：

| 條件 | 運費 |
| --- | --- |
| `amount` 小於等於 0 | 印出「金額錯誤」 |
| `amount` 小於 500 | 100 元 |
| `amount` 介於 500～999 | 60 元 |
| `amount` 大於等於 1000 | 免運 |

請用 `if` / `else if` / `else` 完成，分別用 `-10`、`300`、`800`、`1200` 測試。

<!--
這個練習是 else if 的經典應用：範圍判斷。

記得剛剛說的：條件的順序很重要。想想看，要從哪一個條件開始寫？
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	for _, amount := range []int{-10, 300, 800, 1200} {
		if amount <= 0 {
			fmt.Println(amount, "金額錯誤")
		} else if amount < 500 {
			fmt.Println(amount, "運費 100")
		} else if amount < 1000 {
			fmt.Println(amount, "運費 60")
		} else {
			fmt.Println(amount, "免運")
		}
	}
}
```

- 這裡先偷用了等一下會教的 `for range`，一次測試四個金額

<!--
先處理錯誤的情況，再由小到大判斷範圍。因為前面的條件已經排除了小於 500 的情況，第三個條件只要寫 amount < 1000 就好，不用寫成 amount >= 500 && amount < 1000。

為了一次測試四個金額，這裡先用了等一下會教的 for range，大家先把它當成「把每個金額拿出來跑一次」就好。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# switch 敘述
## Switch Statements

<!--
當選項很多的時候，一長串的 else if 會很難讀。這時候 switch 就派上用場了。
-->

---

# switch 敘述基礎

```go
package main

import "fmt"

func main() {
	day := "Sat"
	switch day {
	case "Mon", "Tue", "Wed", "Thu", "Fri": // 一個 case 可以有多個值
		fmt.Println("上班日")
	case "Sat", "Sun":
		fmt.Println("週末") // ✅
	default:
		fmt.Println("無效的日期")
	}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 的 switch 特色：</b> <b>自動 break</b>（不會掉到下一個 case）、一個 <code>case</code> 可以列<b>多個值</b>、字串／整數／運算式<b>都能用</b>。
</div>

<!--
什麼是 switch？想像自動販賣機：投幣之後按下按鈕，按鈕 1 掉出可樂、按鈕 2 掉出綠茶。switch 就是根據一個值，跳到對應的 case 執行。

Go 的 switch 比 C、Java 好用很多。最大的差別是「自動 break」：在 C 和 Java 裡，每個 case 結尾都要寫 break，忘記寫就會繼續執行下一個 case，這是非常經典的 bug。Go 直接反過來，預設就會自動結束。

另外一個 case 可以列出多個值，用逗號隔開，像這裡把週一到週五放在同一個 case。
-->

---

# switch 的不同用法：無條件 switch

`switch` 後面不接值，每個 `case` 寫一個布林條件 — 等同於 `if-else if` 鏈，但更易讀：

```go
package main

import "fmt"

func main() {
	bmi := 26.3
	switch {
	case bmi < 18.5:
		fmt.Println("過輕")
	case bmi < 24:
		fmt.Println("正常")
	case bmi < 27:
		fmt.Println("過重") // ✅
	default:
		fmt.Println("肥胖")
	}
}
```

<!--
switch 後面可以什麼都不接，這時候每個 case 就寫一個布林條件，從上往下檢查，第一個成立的就執行。

這個寫法在效果上跟 if-else if 完全一樣，但因為每個 case 都對齊，讀起來更清楚。Go 社群的慣例是：當 else if 超過兩三個的時候，改用無條件 switch。
-->

---

# switch 的不同用法：起始賦值

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	switch h := time.Now().Hour(); { // 起始賦值，h 只在 switch 內有效
	case h < 12:
		fmt.Println("早安")
	case h < 18:
		fmt.Println("午安")
	default:
		fmt.Println("晚安")
	}
}
```

<!--
switch 也有起始賦值，跟 if 一樣，用分號隔開。這個例子先取得目前的小時數，存到 h，然後用無條件 switch 判斷該說早安、午安還是晚安。

注意分號後面什麼都沒寫，代表這是一個「有起始賦值的無條件 switch」。h 只在 switch 區塊裡有效。
-->

---

# switch 的不同用法：fallthrough

```go
package main

import "fmt"

func main() {
	switch level := 2; level {
	case 2:
		fmt.Println("解鎖進階功能")
		fallthrough // 強制繼續執行下一個 case
	case 1:
		fmt.Println("解鎖基本功能")
	}
}
```

<!--
前面說過 Go 的 switch 會自動 break，如果真的需要「做完這個 case 繼續做下一個」，就要明確寫出 fallthrough。等級 2 的會員會先印出「解鎖進階功能」，然後 fallthrough 到 case 1，再印出「解鎖基本功能」。

注意 fallthrough 會無條件執行下一個 case，不會檢查下一個 case 的條件。實務上它很少用到，知道有這個東西就好。
-->

---
layout: default
---

# 練習 2：季節判斷
### 任務說明

給定月份 `month`（1～12），用 `switch` 印出季節：

| 月份 | 季節 |
| --- | --- |
| 3、4、5 | 春天 |
| 6、7、8 | 夏天 |
| 9、10、11 | 秋天 |
| 12、1、2 | 冬天 |
| 其他 | 月份錯誤 |

<!--
這個練習要用到「一個 case 列出多個值」和 default。

先自己試試看，寫完之後想想看：如果用 if-else if 寫，會長什麼樣子？哪一種比較好讀？
-->

---

# 練習 2：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	month := 11
	switch month {
	case 3, 4, 5:
		fmt.Println("春天")
	case 6, 7, 8:
		fmt.Println("夏天")
	case 9, 10, 11:
		fmt.Println("秋天")
	case 12, 1, 2:
		fmt.Println("冬天")
	default:
		fmt.Println("月份錯誤")
	}
}
```

<!--
每個 case 用逗號列出三個月份，不需要寫 break，Go 會自動結束。月份 11 會印出「秋天」。

如果用 if 寫，每個條件都要寫成 month == 3 || month == 4 || month == 5，是不是 switch 清楚多了？
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 迴圈
## Loops

<!--
接下來是重複做事的工具：迴圈。Go 只有一種迴圈，就是 for。
-->

---

# for 迴圈基礎：Go 只有 for

| 型態 | 語法 | 相當於其他語言的 |
| --- | --- | --- |
| **條件式** | `for 條件 { }` | `while` |
| **無限迴圈** | `for { }` | `while (true)` |
| **三段式** | `for 初始; 條件; 後置 { }` | 傳統 `for` |
| **範圍式** | `for i, v := range 集合 { }` | `foreach` |

<!--
什麼是迴圈？就是「重複執行同一段程式碼，直到某個條件不成立」。

其他語言通常有 for、while、do-while、foreach 好幾種迴圈，Go 只有 for 一個關鍵字，但它有四種寫法，可以涵蓋所有的情況。這又是「少即是多」的設計。
-->

---

# 條件式 for：相當於 while

```go
package main

import "fmt"

func main() {
	n := 1
	for n < 100 { // 條件式：相當於 while
		n *= 2
	}
	fmt.Println(n) // 128
}
```

<!--
範例是條件式的 for，效果跟其他語言的 while 一樣：只要 n 小於 100，就一直乘以 2。1、2、4、8…一直到 128 的時候條件不成立，迴圈結束，印出 128。
-->

---

# 無限迴圈

`for` 後面什麼都不寫，就是無限迴圈，通常搭配 `break` 或 `return` 離開：

```go
package main

import "fmt"

func main() {
	balance := 1000
	month := 0
	for {
		month++
		balance -= 300
		if balance < 300 {
			break // 餘額不足，離開迴圈
		}
	}
	// 第 3 個月餘額不足，剩下 100
	fmt.Println("第", month, "個月餘額不足，剩下", balance)
}
```

<!--
for 後面什麼都不接，就是無限迴圈。它不會自己停下來，必須在裡面用 break 或 return 離開。

聽起來很危險，但無限迴圈在實務上其實很常見：例如伺服器要不停地接收請求、程式要不停地讀取使用者輸入，直到使用者輸入「離開」為止。

範例是每個月扣 300 元，扣到餘額不足 300 就停下來。
-->

---

# for i 迴圈：三段式

語法：`for 初始敘述; 條件; 後置敘述 { }`

```go
package main

import "fmt"

func main() {
	for i := 1; i <= 5; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println() // 1 2 3 4 5

	for i := 10; i > 0; i -= 3 { // 也可以倒數、跳著數
		fmt.Print(i, " ")
	}
	fmt.Println() // 10 7 4 1
}
```

<!--
這是最傳統的三段式 for 迴圈，用分號分成三個部分：初始、條件、後置。跟 if 一樣，這三段不需要小括號。

第一個迴圈印出 1 到 5。第二個例子示範倒數、每次減 3，後置敘述可以是任何運算，印出 10 7 4 1。
-->

---

# 三段式 for 的執行順序

| 部分 | 何時執行 |
| --- | --- |
| `i := 1` | 迴圈開始前執行**一次** |
| `i <= 5` | **每一圈開始前**檢查，`false` 就結束 |
| `i++` | **每一圈結束後**執行 |

```text
i := 1 → 檢查 i <= 5 → 執行內容 → i++
       → 檢查 i <= 5 → 執行內容 → i++
       → …… → 條件為 false，結束
```

<!--
執行順序是：先執行一次初始敘述 i := 1；然後檢查條件，成立就執行迴圈內容；執行完做後置敘述 i++；再回去檢查條件，一直重複到條件不成立。

就像跑操場：起跑前先站到起跑線（初始，只做一次）；每跑一圈前看看還要不要跑（條件）；跑完一圈記一筆（後置）。
-->

---

# 補充：Go 1.22 起可以 range 一個整數

只是想「重複 N 次」的時候，最新的寫法更簡潔：

```go
package main

import "fmt"

func main() {
	for i := range 5 { // i 依序是 0, 1, 2, 3, 4
		fmt.Print(i, " ")
	}
	fmt.Println()

	for range 3 { // 不需要索引時，連變數都可以省略
		fmt.Println("Go!")
	}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>現代寫法：</b> <code>for i := range n</code> 等同於 <code>for i := 0; i &lt; n; i++</code>。從 Go 1.22 開始支援，新的程式碼建議優先使用。
</div>

<!--
這是 Go 1.22 加入的新語法，現在已經是主流寫法。

以前要重複 5 次，得寫 for i := 0; i < 5; i++，三個部分都要寫。現在可以直接寫 for i := range 5，i 會從 0 跑到 4。如果連 i 都用不到，連變數都可以省略，寫 for range 3 就好。

VS Code 的 gopls 甚至會提示我們：「這個三段式迴圈可以改寫成 range 整數」。
-->

---

# for range 迴圈

`for range` 可以走訪切片、字串、map 等集合，每一圈拿到**索引與值**：

```go
package main

import "fmt"

func main() {
	fruits := []string{"蘋果", "香蕉", "芭樂"}
	for i, f := range fruits {
		fmt.Println(i, f)
	}
	for _, f := range fruits { // 不需要索引時用 _ 忽略
		fmt.Print(f, " ")
	}
	fmt.Println()
}
```

<!--
for range 是 Go 最常用的迴圈，它會自動走訪一個集合裡的每一個元素，相當於其他語言的 foreach。

每一圈會拿到兩個值：索引和元素值。如果不需要索引，用底線 _ 把它丟掉。還記得 Go 規定「沒用到的變數是編譯錯誤」嗎？底線就是告訴 Go「這個值我故意不要」。
-->

---

# for range 能走訪的對象

| 走訪對象 | 第一個值 | 第二個值 |
| --- | --- | --- |
| 切片、陣列 | 索引 | 元素值 |
| 字串 | 位元組位置 | 字元（`rune`） |
| map | 鍵（key） | 值（value） |
| 整數 `n` | `0` 到 `n-1` | — |
| 通道（channel） | 收到的值 | — |
| 迭代器函式（Go 1.23+） | 依函式定義 | 依函式定義 |

<!--
表格列出了不同集合走訪時拿到的值。切片和 map 是第 4 章的內容，字串的走訪第 3 章會詳細說明，通道是第 16 章的內容；迭代器函式是 Go 1.23 加入的新功能，第 5 章學函式時會介紹。這裡先有印象就好。
-->

---

# for range 走訪字串與 map

```go
package main

import "fmt"

func main() {
	for i, ch := range "Go語言" { // 以字元（rune）為單位走訪
		fmt.Printf("%d:%c ", i, ch)
	}
	fmt.Println() // 0:G 1:o 2:語 5:言

	prices := map[string]int{"咖啡": 60, "紅茶": 30}
	for name, p := range prices { // ⚠️ map 的走訪順序不固定
		fmt.Println(name, p)
	}
}
```

<!--
這段程式碼有兩個值得注意的地方。

第一，走訪字串的時候，每一圈拿到的是一個「字元」，而不是一個位元組。「語」這個中文字在 UTF-8 裡佔 3 個位元組，所以下一個字「言」的位置是 5，不是 3。這個細節第 3 章講字串時會詳細說明。

第二，map 的走訪順序是不固定的，每次執行都可能不一樣。這是 Go 刻意的設計，避免我們寫出依賴順序的程式。如果需要固定順序，要先把 key 排序，第 4 章會教。
-->

---

# 使用 for 迴圈的注意事項：迴圈變數

**Go 1.22 起，每一圈的迴圈變數都是新的一份**，閉包或 goroutine 捕捉時不再踩坑：

```go
package main

import "fmt"

func main() {
	var prints []func()
	for i := range 3 {
		prints = append(prints, func() { fmt.Print(i, " ") })
	}
	for _, p := range prints {
		p()
	}
	fmt.Println() // Go 1.22+：0 1 2（舊版會印出 3 3 3）
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>版本差異：</b> 這個行為由 <code>go.mod</code> 的 <code>go</code> 版本決定。網路上舊文章寫的 <code>i := i</code> 技巧，在 Go 1.22 以後已經不需要了。
</div>

<!--
使用 for 迴圈的注意事項：迴圈變數的行為在 Go 1.22 有一個重大的改變。

這段程式碼在迴圈裡建立了三個函式，每個函式都印出 i。在 Go 1.21 以前，三個函式共用同一個 i，等到執行的時候 i 已經變成 3，所以會印出 3 3 3。這是 Go 史上最有名的坑之一。

從 Go 1.22 開始，每一圈的 i 都是一個新的變數，所以會正確印出 0 1 2。

為什麼要特別提？因為網路上很多舊文章會教大家在迴圈裡寫 i := i 來避開這個問題，現在這個技巧已經不需要了。看到舊程式碼有這種寫法，就知道它是為了相容舊版 Go。
-->

---

# break 和 continue 敘述

`break`：**立刻結束**整個迴圈｜`continue`：**跳過這一圈剩下的程式碼**，直接進入下一圈

```go
package main

import "fmt"

func main() {
	for i := range 10 {
		if i%2 == 0 {
			continue // 偶數跳過
		}
		if i > 7 {
			break // 大於 7 就結束
		}
		fmt.Print(i, " ")
	}
	fmt.Println() // 1 3 5 7
}
```

<!--
break 和 continue 可以改變迴圈的流程。

break 是「整個迴圈都不做了」，continue 是「這一圈不做了，直接進下一圈」。

用排隊買票比喻：continue 就像跳過某一位客人，直接服務下一位；break 就像直接關門，後面的客人都不服務了。

範例裡，偶數會被 continue 跳過，所以只印奇數；到了 9 的時候大於 7，break 結束迴圈。結果印出 1 3 5 7。
-->

---

# 補充：用標籤跳出多層迴圈

在巢狀迴圈中，`break` 只會跳出**最內層**；要跳出外層，可以加上**標籤（label）**：

```go
package main

import "fmt"

func main() {
	matrix := [][]int{{1, 2, 3}, {4, -1, 6}, {7, 8, 9}}
outer:
	for r, row := range matrix {
		for c, v := range row {
			if v < 0 {
				fmt.Printf("在 (%d,%d) 找到負數\n", r, c)
				break outer // 直接跳出外層迴圈
			}
		}
	}
}
```

<!--
這是補充內容。

在兩層迴圈裡，break 只會跳出它所在的那一層。如果我們在內層找到目標，想要連外層一起結束，就要用標籤。

在外層迴圈前面寫一個名字加冒號，例如 outer:，然後在內層寫 break outer，就能直接跳出外層迴圈。continue 也可以搭配標籤使用，意思是直接進入外層的下一圈。

執行後會印出「在 (1,1) 找到負數」。
-->

---
layout: default
---

# 練習 3：九九乘法表
### 任務說明

1. 用兩層 `for` 迴圈印出九九乘法表（1×1 到 9×9）
2. 每一列印出同一個被乘數，例如：`2x1=2 2x2=4 ... 2x9=18`
3. **進階**：只印出乘積是**偶數**的算式（使用 `continue`）
4. 請使用 Go 1.22 以後的 `range` 整數寫法

<!--
九九乘法表是巢狀迴圈最經典的練習。

提示：外層迴圈控制被乘數，內層迴圈控制乘數。range 9 會從 0 開始，要怎麼讓它從 1 開始？
-->

---

# 練習 3：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	for i := range 9 {
		for j := range 9 {
			a, b := i+1, j+1 // range 從 0 開始，所以加 1
			if (a*b)%2 != 0 {
				continue // 進階：跳過奇數乘積
			}
			fmt.Printf("%dx%d=%-3d", a, b, a*b)
		}
		fmt.Println()
	}
}
```

- `%-3d` 表示整數靠左對齊、至少佔 3 格，讓輸出排列整齊

<!--
range 9 會產生 0 到 8，所以我們各加 1 變成 1 到 9。也可以用傳統寫法 for i := 1; i <= 9; i++，兩種都對。

continue 會跳過乘積是奇數的算式。Printf 的 %-3d 是格式化技巧：減號代表靠左對齊，3 代表至少佔 3 個字元寬，這樣每一列的算式會排得很整齊。格式化輸出第 9 章會詳細介紹。
-->

---
layout: default
---

# 綜合練習：FizzBuzz
### 任務說明

印出 1 到 30，但是：

| 條件 | 印出 |
| --- | --- |
| 同時是 3 和 5 的倍數 | `FizzBuzz` |
| 是 3 的倍數 | `Fizz` |
| 是 5 的倍數 | `Buzz` |
| 其他 | 數字本身 |

要求：用 `for range` 整數迴圈 + **無條件 switch** 完成，最後印出 `FizzBuzz` 出現了幾次。

<!--
FizzBuzz 是程式設計面試裡最經典的題目，它考的就是今天學的迴圈加上條件判斷。

特別注意條件的順序：同時是 3 和 5 的倍數，要放在哪裡檢查？第 13 章資料庫的練習還會再用到 FizzBuzz，所以這題一定要會。
-->

---
zoom: 0.9
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	count := 0
	for n := range 30 {
		n++ // 讓 n 從 1 開始
		switch {
		case n%15 == 0:
			fmt.Println("FizzBuzz")
			count++
		case n%3 == 0:
			fmt.Println("Fizz")
		case n%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(n)
		}
	}
	fmt.Println("FizzBuzz 次數：", count) // 2
}
```

<!--
關鍵是條件的順序：15 的倍數一定要放在最前面檢查。如果先檢查 3 的倍數，15 會被判斷成 Fizz，因為 switch 只會執行第一個成立的 case。

n++ 這行讓 range 產生的 0 到 29 變成 1 到 30。因為 Go 1.22 之後每一圈的 n 都是新的變數，在迴圈裡修改它不會影響下一圈。

1 到 30 之間，15 和 30 是 FizzBuzz，所以次數是 2。
-->

---

# 章節總結

- **if**：條件不加括號、大括號必寫；`else` 要和 `}` 同一行；善用**起始賦值**與**提早 return**
- **switch**：自動 break；一個 `case` 可列多個值；無條件 `switch` 取代長串 `else if`；`fallthrough` 很少用
- **for 是唯一的迴圈**：條件式（while）、無限迴圈、三段式、`range`
- **現代寫法**：`for i := range n`（Go 1.22+）；每一圈的迴圈變數都是新的一份
- **break / continue**：結束迴圈／跳過這一圈；巢狀迴圈可搭配標籤

下一章我們會深入介紹 Go 的「核心型別」：布林、整數、浮點數、字串與 rune。

<!--
我們來整理今天學到的東西。

if 的部分，記得起始賦值和提早 return 這兩個 Go 特有的慣例。switch 會自動 break，無條件 switch 可以取代一長串的 else if。for 是 Go 唯一的迴圈，但有四種寫法，現在最主流的是 range 整數和 range 集合。

有了條件判斷和迴圈，我們已經可以寫出很多實用的小程式了。

下一章我們會回頭仔細看 Go 的核心型別：整數有哪幾種、浮點數為什麼會有誤差、字串和 rune 有什麼不一樣。今天範例裡走訪中文字串時位置跳了 3 格的原因，下一章就會揭曉。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天學的流程控制，是所有程式的骨架。

課後練習建議：試著把 FizzBuzz 改成用 if-else if 寫一次，比較兩種寫法的可讀性；再把九九乘法表改成印出三角形的格式。

有問題的同學現在可以提問！
-->

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
title: 函式
routeAlias: ch05
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">函式</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「把一段邏輯包起來、取個名字，就能一再重複使用」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第五章！

前面幾章我們其實已經寫了不少函式，main 就是一個函式，上一章也替 struct 寫了好幾個方法。今天我們要把函式完整、有系統地學一遍。

Go 的函式有幾個很有特色的地方：可以一次回傳多個值、函式本身可以當成值傳來傳去、還有一個叫做 defer 的關鍵字，可以讓某段程式碼「延後到函式結束時才執行」。這些特色會在後面的錯誤處理、檔案處理、HTTP 伺服器章節大量使用。
-->

---
layout: default
---

# Outline

- **函式基礎** — 宣告、參數、回傳值、Naked Returns
- **參數不定函式** — `...T` 接住任意數量的參數
- **匿名函式與閉包** — 沒有名字的函式、記住外部變數的函式
- **以函式為型別** — 自訂函式型別、當參數、當回傳值
- **補充：泛型與迭代器** — Go 1.18+ 泛型、Go 1.23+ `range` 函式
- **defer** — 延後執行、執行順序、變數值的副作用
- **章節總結**

<!--
今天的內容從最基本的函式宣告開始，一路學到閉包和函式型別，這兩個是比較進階、但實務上很常用的概念。

中間會補充兩個現代 Go 的功能：泛型和迭代器。最後是 defer，它是 Go 處理「收尾工作」的標準方式。
-->

---

# 回顧：複合型別

- **切片**：`s = append(s, v)` 一定要接住；共用底層陣列要小心
- **map**：用 `v, ok := m[k]` 判斷鍵是否存在
- **struct** + **方法**：`func (r *Rect) Scale(k float64)`，指標接收器可以修改原值
- **any** 與型別斷言：`v, ok := x.(T)`；型別 switch

<!--
回顧一下上一章。

我們學了切片、map、struct，也替 struct 寫了方法。方法其實就是「有接收器的函式」，今天學的所有函式規則，方法也都適用。

上一章我們一直在用 comma ok 這種「一次回傳兩個值」的寫法，今天就會知道怎麼自己寫出這樣的函式。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 函式
## Functions

<!--
先從函式的基本結構開始。
-->

---

# 函式的宣告和組成

```go
func add(a int, b int) int {
	return a + b
}
```

| 組成 | 範例 | 說明 |
| --- | --- | --- |
| **關鍵字** | `func` | 宣告函式 |
| **名稱** | `add` | 小寫開頭：套件內使用；大寫開頭：可被其他套件使用 |
| **參數列** | `(a int, b int)` | 名稱在前、型別在後 |
| **回傳型別** | `int` | 寫在參數列後面；沒有回傳值就省略 |
| **函式本體** | `{ return a + b }` | 用 `return` 回傳結果 |

<!--
什麼是函式？函式就是「把一段邏輯包起來，取一個名字」，之後只要呼叫名字，就能執行這段邏輯。就像咖啡機：放入咖啡豆和水（參數），按下按鈕，就會得到一杯咖啡（回傳值），我們不用每次都自己磨豆、煮水。

Go 的函式由五個部分組成：func 關鍵字、名稱、參數列、回傳型別、函式本體。跟變數宣告一樣，型別寫在名稱後面。

函式名稱的大小寫有特殊意義：大寫開頭的函式可以被其他套件使用，小寫開頭的只能在同一個套件裡使用。這個規則第 8 章會詳細說明。
-->

---

# 函式參數

同型別的連續參數可以**合併寫**；參數一律**傳值**（複製一份）

```go
package main

import "fmt"

func area(w, h float64) float64 { // w、h 都是 float64
	return w * h
}

func greet(name string, times int) {
	for range times {
		fmt.Println("Hi,", name)
	}
}

func main() {
	fmt.Println(area(3, 4)) // 12
	greet("Gopher", 2)
}
```

<!--
參數的部分有兩個重點。

第一，連續的同型別參數可以合併，w, h float64 等同於 w float64, h float64，這是 Go 程式碼裡很常見的寫法。

第二，Go 的參數一律是「傳值」，函式收到的是複製品。第一章我們學過，想要讓函式修改外面的變數，就要傳指標。

不過要注意：切片和 map 雖然也是傳值，但複製的是那個「小結構」（指標、長度、容量），底層的資料是共用的，所以函式裡修改切片的元素，外面看得到。這是上一章切片內部運作的延伸。
-->

---
zoom: 0.96
---

# 函式傳回值：多重回傳值

Go 的函式可以**一次回傳多個值**，最常見的是「結果 + 錯誤」

```go
package main

import (
	"errors"
	"fmt"
)

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除數不能是 0")
	}
	return a / b, nil
}

func main() {
	if q, err := divide(10, 4); err == nil {
		fmt.Println(q) // 2.5
	}
	_, err := divide(1, 0) // 不需要的回傳值用 _ 忽略
	fmt.Println(err)       // 除數不能是 0
}
```

<!--
多重回傳值是 Go 最有特色的功能之一。

回傳型別用小括號包起來，就能回傳多個值。最常見的模式是「結果 + 錯誤」：成功時回傳結果和 nil，失敗時回傳零值和一個錯誤。呼叫的人拿到之後，先檢查錯誤，再使用結果。

這就是 Go 不需要 try-catch 的原因：錯誤就是一個普通的回傳值。下一章錯誤處理會完整介紹這個模式。

不需要的回傳值，用底線忽略。但要注意：忽略 error 是很危險的習慣，除非真的確定不會出錯。
-->

---

# Naked Returns（具名回傳值）

回傳值可以取名字，它們會被初始化為零值；`return` 後面不寫東西就回傳它們

```go
package main

import "fmt"

func stats(nums []int) (sum int, avg float64) { // 具名回傳值
	for _, n := range nums {
		sum += n
	}
	if len(nums) > 0 {
		avg = float64(sum) / float64(len(nums))
	}
	return // naked return：自動回傳 sum 和 avg
}

func main() {
	s, a := stats([]int{80, 90, 100})
	fmt.Println(s, a) // 270 90
}
```

<!--
回傳值也可以取名字，這叫做「具名回傳值」。它們會在函式開始時被初始化為零值，可以在函式裡直接使用。

return 後面什麼都不寫，就會回傳這些具名變數的目前值，這叫做 naked return（裸回傳）。

具名回傳值的好處是：函式簽章本身就說明了每個回傳值的意義，一看就知道第一個是總和、第二個是平均。
-->

---

# 使用 Naked Returns 的注意事項

**注意事項之一：** 長函式裡的 naked return **很難讀**，看到 `return` 不知道回傳了什麼

**注意事項之二：** 具名回傳值**很適合當作文件**，但 `return` 時仍建議寫清楚

```go
// ✅ 推薦：具名回傳值當文件，return 時寫清楚
func splitName(full string) (first, last string) {
	first, last, _ = strings.Cut(full, " ")
	return first, last
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 官方建議：</b> naked return 只用在<b>很短</b>的函式。具名回傳值另一個重要用途是搭配 <code>defer</code> 修改回傳值（本章最後介紹）。
</div>

<!--
使用 naked return 的注意事項。

在短短幾行的函式裡，naked return 很方便；但如果函式有幾十行，讀到最後一行的 return，我們得往上找半天才知道到底回傳了什麼。所以 Go 官方的建議是：naked return 只用在很短的函式。

實務上比較推薦的做法是：用具名回傳值讓函式簽章更清楚，但 return 的時候還是把值寫出來。

另外取名時要小心：如果把回傳值取名叫 min、max，會遮蔽掉 Go 1.21 加入的內建函式 min、max，這是一個容易忽略的小細節。
-->

---
layout: default
---

# 練習 1：溫度統計
### 任務說明

寫一個函式 `tempStats(temps []float64) (low, high, avg float64, err error)`：

1. 如果切片是空的，回傳錯誤 `errors.New("沒有資料")`
2. 否則計算最低溫、最高溫、平均溫度
3. 在 `main` 中用 `[]float64{23.5, 28.1, 19.8, 31.2}` 和空切片各呼叫一次
4. 印出結果時，平均溫度只顯示到小數點後一位

<!--
這個練習會用到多重回傳值和具名回傳值，還有「結果 + 錯誤」的模式。

提示：Go 1.21 以後有內建的 min 和 max 函式，可以直接拿來比較兩個數字。但如果回傳值取名叫 low、high，就不會跟內建函式衝突。
-->

---
zoom: 0.79
---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
)

func tempStats(temps []float64) (low, high, avg float64, err error) {
	if len(temps) == 0 {
		return 0, 0, 0, errors.New("沒有資料")
	}
	low, high = temps[0], temps[0]
	sum := 0.0
	for _, t := range temps {
		low, high = min(low, t), max(high, t) // Go 1.21+ 內建 min / max
		sum += t
	}
	return low, high, sum / float64(len(temps)), nil
}

func main() {
	l, h, a, _ := tempStats([]float64{23.5, 28.1, 19.8, 31.2})
	fmt.Printf("低 %.1f 高 %.1f 平均 %.1f\n", l, h, a) // 低 19.8 高 31.2 平均 25.7
	if _, _, _, err := tempStats(nil); err != nil {
		fmt.Println("錯誤：", err)
	}
}
```

<!--
空切片的時候提早 return 錯誤。

計算最低最高溫時，用了 Go 1.21 加入的內建函式 min 和 max，它們可以接受任意數量的參數，比以前自己寫 if 比較方便很多。

注意呼叫 tempStats(nil) 也可以：nil 切片的長度是 0，所以會回傳錯誤。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 參數不定函式
## Variadic Functions

<!--
接下來是參數不定函式：可以接受任意數量參數的函式。
-->

---

# 參數不定函式

在最後一個參數的型別前加 `...`，就能接受**任意數量**的參數，函式內收到的是一個**切片**

```go
package main

import "fmt"

func sum(label string, nums ...int) int { // nums 的型別是 []int
	total := 0
	for _, n := range nums {
		total += n
	}
	fmt.Println(label, len(nums), "個數字")
	return total
}

func main() {
	fmt.Println(sum("A"))          // 0 個數字 → 0
	fmt.Println(sum("B", 1, 2, 3)) // 3 個數字 → 6
	scores := []int{90, 80}
	fmt.Println(sum("C", scores...)) // 用 ... 展開切片 → 170
}
```

<!--
什麼是參數不定函式？我們其實每天都在用：fmt.Println 可以印一個值，也可以印十個值，它就是參數不定函式。

在最後一個參數的型別前面加上三個點，就代表「這裡可以接受零個或多個 int」。在函式裡，nums 就是一個 []int 切片，可以用 for range 走訪。

如果手上已經有一個切片，想把它傳進去，就在切片後面加三個點把它「展開」。上一章 append(nums, more...) 就是同樣的語法，因為 append 本身就是一個參數不定函式。
-->

---

# 使用參數不定函式的注意事項

**注意事項之一：** `...` 只能用在**最後一個**參數

**注意事項之二：** 用 `s...` 展開傳入時，函式收到的是**同一個切片**，修改會影響外部

```go
package main

import "fmt"

func resetFirst(nums ...int) {
	if len(nums) > 0 {
		nums[0] = 0
	}
}

func main() {
	data := []int{5, 6, 7}
	resetFirst(data...)
	fmt.Println(data) // [0 6 7]：外面的切片被改了！

	resetFirst(8, 9) // 逐一傳入時，Go 會建立一個新切片
}
```

<!--
使用參數不定函式有兩個注意事項。

第一，三個點只能用在最後一個參數，因為 Go 需要知道前面的參數到哪裡結束。

第二，用切片加三個點傳進去時，Go 不會複製，函式拿到的就是同一個切片，共用同一個底層陣列，所以函式裡修改元素，外面也會看到。如果是一個一個傳進去，Go 會自己建立一個新的切片，就不會影響外面。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 匿名函式與閉包
## Anonymous Functions & Closures

<!--
接下來是匿名函式與閉包，這是 Go 裡函式「當成值」使用的開始。
-->

---

# 宣告匿名函式

沒有名字的函式，可以**存進變數**，或**宣告後立刻執行**

```go
package main

import "fmt"

func main() {
	square := func(n int) int { // 把匿名函式存進變數
		return n * n
	}
	fmt.Println(square(5)) // 25

	func(msg string) { // 宣告後立刻執行（IIFE）
		fmt.Println("立刻執行：", msg)
	}("Hello")

	fmt.Printf("%T\n", square) // func(int) int：函式也有型別
}
```

<!--
什麼是匿名函式？就是沒有名字的函式。

在 Go 裡，函式是「一等公民」，意思是函式跟 int、string 一樣是一種值，可以存進變數、當參數傳遞、當回傳值回傳。

第一種用法是把匿名函式存進變數，之後用變數名稱呼叫。第二種用法是宣告完馬上加括號執行，這在搭配 defer 和 goroutine 的時候很常見，第 16 章會大量看到 go func() { ... }() 這種寫法。

最後一行用 %T 印出 square 的型別：func(int) int，代表「接收一個 int、回傳一個 int 的函式」。函式的型別由參數和回傳值決定。
-->

---

# 什麼是閉包？

「**閉包（closure）是一個『記得』它被建立時周圍變數的函式**，即使那些變數的作用範圍已經結束。」

```go
package main

import "fmt"

func newCounter() func() int {
	count := 0          // 區域變數
	return func() int { // 回傳一個匿名函式，它「捕捉」了 count
		count++
		return count
	}
}

func main() {
	c1 := newCounter()
	fmt.Println(c1(), c1(), c1()) // 1 2 3
	c2 := newCounter()            // 新的閉包，有自己的 count
	fmt.Println(c2(), c1())       // 1 4
}
```

<!--
什麼是閉包？

想像一個背包：函式在建立的時候，會把它用到的外部變數裝進背包裡帶走。就算原本的函式已經結束了，背包裡的東西還在，而且每次呼叫都能繼續使用、修改。

newCounter 裡的 count 本來是區域變數，函式結束就該消失了。但因為回傳的匿名函式用到了 count，Go 會把 count 保留下來，放在閉包的「背包」裡。每次呼叫 c1，都是在修改同一個 count，所以會得到 1、2、3。

c2 是另一次呼叫 newCounter 產生的，它有自己的背包、自己的 count，所以從 1 開始，跟 c1 互不影響。
-->

---

# 建立閉包：實用範例

閉包很適合用來「產生客製化的函式」：

```go
package main

import (
	"fmt"
	"strings"
)

func makeGreeter(greeting string) func(string) string {
	return func(name string) string {
		return greeting + "，" + strings.ToUpper(name) + "！"
	}
}

func main() {
	hello := makeGreeter("你好")
	morning := makeGreeter("早安")
	fmt.Println(hello("go"))       // 你好，GO！
	fmt.Println(morning("gopher")) // 早安，GOPHER！
}
```

<!--
閉包最實用的地方，是「用一個函式產生很多個客製化的函式」。

makeGreeter 就像一個「問候語工廠」：給它「你好」，它就生產出一個會說「你好」的函式；給它「早安」，就生產出一個會說「早安」的函式。每個產出的函式都記得自己的 greeting。

這個模式在實務上很常見，例如第 15 章的 HTTP 中介軟體（middleware），就是用閉包包裝處理器；第 9 章的 logger 也可以用閉包加上前綴。
-->

---
layout: default
---

# 練習 2：累加器
### 任務說明

1. 寫一個函式 `makeAccumulator(start int) func(int) int`
2. 回傳的函式每次被呼叫時，會把參數加進總和，並回傳目前的總和
3. 建立兩個獨立的累加器：`a` 從 0 開始、`b` 從 100 開始
4. 分別呼叫 `a(10)`、`a(5)`、`b(1)`、`a(1)`，印出結果，驗證它們互不影響

<!--
這個練習是閉包的經典應用，跟剛剛的計數器很像，只是多了參數。

想想看：總和這個變數要宣告在哪裡？
-->

---

# 練習 2：解題提示
### 提示說明

```go
package main

import "fmt"

func makeAccumulator(start int) func(int) int {
	total := start // 被閉包捕捉的變數
	return func(n int) int {
		total += n
		return total
	}
}

func main() {
	a := makeAccumulator(0)
	b := makeAccumulator(100)
	fmt.Println(a(10), a(5), b(1), a(1)) // 10 15 101 16
}
```

<!--
total 宣告在 makeAccumulator 裡面、匿名函式外面，這樣它才能被閉包捕捉，而且每次呼叫 makeAccumulator 都會產生一個新的 total。

其實參數 start 本身也是區域變數，也可以直接用 start += n，效果一樣。

a 和 b 各自有自己的 total，所以結果是 10、15、101、16。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 以函式為型別的參數
## Function Types

<!--
既然函式是一種值，它也可以有自己的型別名稱，也可以當作參數和回傳值。
-->

---

# 自訂函式型別

用 `type` 替函式簽章取一個名字，讓程式碼更好讀：

```go
package main

import "fmt"

type Operation func(a, b int) int // 「接收兩個 int、回傳 int」的函式型別

func main() {
	var add Operation = func(a, b int) int { return a + b }
	mul := Operation(func(a, b int) int { return a * b })

	ops := map[string]Operation{"+": add, "*": mul}
	for _, sym := range []string{"+", "*"} {
		fmt.Println(3, sym, 4, "=", ops[sym](3, 4))
	}
}
```

<!--
上一章我們學過自訂型別，函式也可以有自訂型別。

type Operation func(a, b int) int 的意思是：「Operation 是一種函式，它接收兩個 int、回傳一個 int」。任何符合這個簽章的函式，都可以存進 Operation 型別的變數。

為什麼要取名字？當函式簽章很長、或者在很多地方重複出現時，取一個有意義的名字，程式碼會好讀很多。

範例裡把加法和乘法存進一個 map，用運算符號查出對應的函式，再呼叫它。這是一種簡單的「策略模式」，比寫一長串 switch 更有彈性。
-->

---
zoom: 0.88
---

# 使用自訂函式型別的參數

把函式當參數傳進去，讓呼叫的人決定「要怎麼做」：

```go
package main

import "fmt"

type Predicate func(int) bool

func filter(nums []int, keep Predicate) []int {
	var out []int
	for _, n := range nums {
		if keep(n) {
			out = append(out, n)
		}
	}
	return out
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6}
	isEven := func(n int) bool { return n%2 == 0 }
	fmt.Println(filter(nums, isEven))                            // [2 4 6]
	fmt.Println(filter(nums, func(n int) bool { return n > 3 })) // [4 5 6]
}
```

<!--
把函式當作參數，是一種非常強大的設計。

filter 函式負責「走訪切片、收集結果」這些固定的工作，但「哪些元素要留下來」由呼叫的人決定，透過 keep 這個函式參數傳進來。

同一個 filter，傳入 isEven 就篩出偶數，傳入「大於 3」的匿名函式就篩出大於 3 的數。標準函式庫裡很多函式都是這樣設計的，例如 slices.SortFunc、slices.IndexFunc、strings.FieldsFunc，第 15 章的 http.HandleFunc 也是把函式當參數傳進去。
-->

---

# 用自訂函式型別作為傳回值

```go
package main

import "fmt"

type Discount func(price float64) float64

func discountFor(level string) Discount {
	switch level {
	case "gold":
		return func(p float64) float64 { return p * 0.8 }
	case "silver":
		return func(p float64) float64 { return p * 0.9 }
	default:
		return func(p float64) float64 { return p }
	}
}

func main() {
	for _, lv := range []string{"gold", "silver", "normal"} {
		fmt.Println(lv, discountFor(lv)(1000)) // 800、900、1000
	}
}
```

<!--
函式也可以當作回傳值。

discountFor 根據會員等級，回傳一個對應的折扣函式。呼叫的人拿到這個函式之後，再把價格傳進去計算。

注意 discountFor(lv)(1000) 這個寫法：第一個括號呼叫 discountFor，拿到一個函式；第二個括號馬上呼叫這個函式，傳入 1000。

這樣的好處是：「決定用哪種折扣」和「計算折扣」這兩件事被分開了，之後新增等級只要改 discountFor 一個地方。
-->

---
zoom: 0.83
---

# 補充：泛型函式（Go 1.18+）

用**型別參數** `[T 約束]`，同一個函式可以處理多種型別，而且**保有型別安全**：

```go
package main

import (
	"cmp"
	"fmt"
)

func Filter[T any](s []T, keep func(T) bool) []T { // T 可以是任何型別
	var out []T
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func Largest[T cmp.Ordered](a, b T) T { // T 必須可以比大小
	return max(a, b)
}

func main() {
	fmt.Println(Filter([]string{"go", "java", "rust"}, func(s string) bool { return len(s) == 4 }))
	fmt.Println(Largest(3, 7), Largest("apple", "banana")) // 7 banana
}
```

<!--
這是補充內容：泛型。

剛剛的 filter 只能處理 []int，如果想要處理 []string，以前只能再寫一個，或是改用 any 失去型別檢查。Go 1.18 加入了泛型，在函式名稱後面用中括號宣告「型別參數」T，就能讓同一個函式處理多種型別。

中括號裡的 any 和 cmp.Ordered 叫做「約束」：any 代表任何型別都可以；cmp.Ordered 代表 T 必須是可以用大於小於比較的型別，例如數字和字串。

呼叫的時候，Go 會自動推斷 T 是什麼，不用自己寫出來。上一章用的 slices.Sort、slices.Contains，其實就是用泛型寫的。實務上建議：先寫具體型別的版本，真的需要重複時才改成泛型。
-->

---
zoom: 0.81
---

# 補充：迭代器函式（Go 1.23+）

符合 `func(yield func(T) bool)` 形式的函式，可以直接放進 `for range`：

```go
package main

import (
	"fmt"
	"iter"
)

func Countdown(n int) iter.Seq[int] { // iter.Seq[int] = func(yield func(int) bool)
	return func(yield func(int) bool) {
		for i := n; i > 0; i-- {
			if !yield(i) { // 呼叫端 break 時，yield 會回傳 false
				return
			}
		}
	}
}

func main() {
	for v := range Countdown(5) {
		if v == 2 {
			break
		}
		fmt.Print(v, " ") // 5 4 3
	}
	fmt.Println()
}
```

<!--
這也是補充內容：迭代器。

上一章我們用過 maps.Keys，它回傳的不是切片，而是一個「迭代器」，可以直接放進 for range。這是 Go 1.23 加入的功能。

迭代器本質上是一個函式，它接收一個叫做 yield 的函式參數，每產生一個值就呼叫一次 yield 把值「交出去」。如果呼叫端的迴圈 break 了，yield 會回傳 false，迭代器就應該停止。

這個設計結合了今天學的兩個概念：函式當參數（yield）、函式當回傳值（Countdown 回傳一個函式）。第一次看可能覺得繞，但會「用」迭代器比會「寫」重要，大家先能看懂 slices.Sorted(maps.Keys(m)) 這種寫法就好。
-->

---
layout: default
---

# 練習 3：計算機
### 任務說明

1. 定義函式型別 `type BinaryOp func(a, b float64) (float64, error)`
2. 建立一個 `map[string]BinaryOp`，包含 `+`、`-`、`*`、`/` 四種運算（`/` 除數為 0 時回傳錯誤）
3. 寫一個函式 `calc(a float64, op string, b float64) (float64, error)`：找不到運算子時回傳錯誤
4. 測試：`calc(6, "*", 7)`、`calc(1, "/", 0)`、`calc(1, "%", 2)`

<!--
這個練習把今天學的函式型別、多重回傳值、map 結合在一起。

提示：運算子找不到的時候，要用 comma ok 判斷 map 裡有沒有這個鍵。錯誤訊息可以用 fmt.Errorf 產生，它跟 Sprintf 很像，但回傳一個 error，下一章會正式介紹。
-->

---

# 練習 3：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
)

type BinaryOp func(a, b float64) (float64, error)

var ops = map[string]BinaryOp{
	"+": func(a, b float64) (float64, error) { return a + b, nil },
	"-": func(a, b float64) (float64, error) { return a - b, nil },
	"*": func(a, b float64) (float64, error) { return a * b, nil },
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, errors.New("除數不能是 0")
		}
		return a / b, nil
	},
}
```

<!--
ops 是一個套件層級的 map，每個運算子對應一個匿名函式。除法的函式多做了一個檢查：除數為 0 時回傳錯誤。
-->

---

# 練習 3：解題提示（續）
### 提示說明

```go
// 續上頁
func calc(a float64, op string, b float64) (float64, error) {
	f, ok := ops[op]
	if !ok {
		return 0, fmt.Errorf("不支援的運算子：%s", op)
	}
	return f(a, b)
}

func main() {
	fmt.Println(calc(6, "*", 7))
	fmt.Println(calc(1, "/", 0))
	fmt.Println(calc(1, "%", 2))
}
```

<!--
calc 先用 comma ok 查 map，找不到就回傳錯誤；找到了就直接 return f(a, b)，因為 f 的回傳值剛好也是 (float64, error)，可以直接回傳。

注意 fmt.Println 可以直接印出多重回傳值，它會把兩個值都印出來：42 <nil>、0 除數不能是 0、0 不支援的運算子：%。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# defer
## Deferred Function Calls

<!--
最後一個主題是 defer，Go 處理收尾工作的標準方式。
-->

---

# 用 defer 延後函式執行

「**`defer` 會把一個函式呼叫延後到『外層函式即將回傳』時才執行**，不論函式是正常結束、提早 return 還是 panic。」

```go
package main

import "fmt"

func process() {
	fmt.Println("1. 開啟資源")
	defer fmt.Println("3. 釋放資源（defer）")

	fmt.Println("2. 處理資料")
	if true {
		return // 提早 return，defer 仍然會執行
	}
	fmt.Println("這行不會執行")
}

func main() {
	process()
}
```

<!--
什麼是 defer？

想像我們去圖書館借書，最怕的就是看完忘了還。如果在借書的當下，就在手機上設一個「離開圖書館前還書」的提醒，就不會忘了。defer 就是這個提醒：在開啟資源的下一行，就寫好「函式結束時要釋放資源」。

defer 後面接一個函式呼叫，這個呼叫不會馬上執行，而是等到外層函式要結束的時候才執行。不管函式是正常結束、中途 return，甚至是 panic，defer 都保證會執行。

執行後會依序印出 1、2、3，第 3 行是 defer 的。
-->

---

# defer 的典型用途

| 情境 | 寫法 | 章節 |
| --- | --- | --- |
| 關閉檔案 | `f, err := os.Open(name)` → `defer f.Close()` | Ch 12 |
| 關閉 HTTP 回應 | `resp, err := http.Get(url)` → `defer resp.Body.Close()` | Ch 14 |
| 關閉資料庫連線／查詢結果 | `defer db.Close()`、`defer rows.Close()` | Ch 13 |
| 解除互斥鎖 | `mu.Lock()` → `defer mu.Unlock()` | Ch 16 |
| 從 panic 復原 | `defer func() { recover() }()` | Ch 6 |
| 測量執行時間 | `defer timeTrack(time.Now())` | Ch 10 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>慣用法：</b> 「<b>取得資源，檢查錯誤，馬上 defer 釋放</b>」三行寫在一起，之後就不用再煩惱忘記釋放。
</div>

<!--
defer 在實務上的用途都跟「收尾」有關：開了檔案要關、鎖了要解鎖、連線要關閉。

這張表列出了後面章節會遇到的典型用法。Go 的慣例是：取得資源、檢查錯誤、馬上 defer 釋放，三件事寫在一起。這樣不管後面的程式碼有多少個 return，資源都一定會被釋放。

這比其他語言的 try-finally 更簡潔，因為「開」和「關」寫在相鄰的兩行，一眼就能確認有沒有漏掉。
-->

---

# 多重 defer 的執行順序

多個 `defer` 會以**後進先出（LIFO）**的順序執行，像疊盤子一樣：

```go
package main

import "fmt"

func main() {
	fmt.Println("開始")
	for i := range 3 {
		defer fmt.Println("defer", i)
	}
	fmt.Println("結束")
}
```

```text
開始
結束
defer 2
defer 1
defer 0
```

<!--
如果一個函式裡有多個 defer，它們的執行順序是「後進先出」，也就是最後 defer 的最先執行。

就像疊盤子：先放的盤子在最下面，最後放的在最上面，拿的時候從最上面開始拿。

為什麼要這樣設計？因為資源通常有依賴關係。例如先開資料庫連線、再開查詢，關閉的時候就要先關查詢、再關連線，剛好是相反的順序。LIFO 讓這件事自動正確。
-->

---

# defer 對變數值的副作用

**注意事項之一：** `defer` 的**參數在 defer 那一行就先求值**，不是在執行時才求值

```go
package main

import "fmt"

func main() {
	x := 10
	defer fmt.Println("defer 參數：", x) // x 在這裡就被求值為 10
	defer func() {
		fmt.Println("defer 閉包：", x) // 閉包在執行時才讀取 x → 20
	}()
	x = 20
	fmt.Println("目前 x：", x)
}
```

```text
目前 x： 20
defer 閉包： 20
defer 參數： 10
```

<!--
使用 defer 的注意事項：參數什麼時候求值？

defer fmt.Println("...", x) 這一行，x 的值在「defer 這一行被執行時」就決定了，是 10。之後 x 改成 20 也不會影響。

但如果 defer 的是一個閉包，閉包裡讀取的 x 是在「真正執行時」才讀取，那時候 x 已經是 20 了。

這兩種行為的差別很重要，也是面試很愛考的題目。記住：參數立刻求值，閉包延後讀取。
-->

---
zoom: 0.88
---

# defer 修改具名回傳值

**注意事項之二：** `defer` 的閉包可以**讀取和修改具名回傳值**

```go
package main

import (
	"errors"
	"fmt"
)

func saveOrder(id int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("儲存訂單 %d 失敗：%w", id, err) // 在回傳前補充錯誤資訊
		}
	}()
	if id <= 0 {
		return errors.New("無效的 ID")
	}
	return nil
}

func main() {
	fmt.Println(saveOrder(0)) // 儲存訂單 0 失敗：無效的 ID
	fmt.Println(saveOrder(7)) // <nil>
}
```

<!--
defer 還有一個進階用法：搭配具名回傳值，在函式回傳之前修改回傳值。

return errors.New("無效的 ID") 這行，會先把錯誤指定給具名回傳值 err，然後執行 defer，最後才真正回傳。所以 defer 裡可以看到 err，還可以把它包裝成更詳細的錯誤訊息。

%w 是包裝錯誤的格式動詞，下一章會詳細說明。這個技巧在實務上常用來「統一補充錯誤資訊」，下一章學 recover 的時候也會用到同樣的模式。
-->

---
layout: default
---

# 綜合練習：中介函式（Middleware）
### 任務說明

1. 定義 `type Handler func(name string) string`
2. 寫一個 `hello` 處理器：回傳 `"Hello, " + name`
3. 寫一個 `withLog(h Handler) Handler`：回傳一個新的 Handler，執行前印出「開始：name」，並用 **defer** 在結束時印出「結束：name」
4. 寫一個 `withUpper(h Handler) Handler`：把結果轉成大寫
5. 組合成 `withLog(withUpper(hello))`，呼叫它並印出結果

<!--
這個綜合練習是第 15 章 HTTP 中介軟體的預告。

「一個函式接收函式、回傳函式」，這就是 middleware 的核心概念：像洋蔥一樣，一層一層把處理器包起來，每一層加上一點功能。

用到的觀念有：函式型別、函式當參數、函式當回傳值、閉包、defer。
-->

---
zoom: 0.81
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"strings"
)

type Handler func(name string) string

func hello(name string) string { return "Hello, " + name }

func withLog(h Handler) Handler {
	return func(name string) string {
		fmt.Println("開始：", name)
		defer fmt.Println("結束：", name)
		return h(name)
	}
}

func withUpper(h Handler) Handler {
	return func(name string) string { return strings.ToUpper(h(name)) }
}

func main() {
	h := withLog(withUpper(hello))
	fmt.Println(h("gopher")) // 開始 → 結束 → HELLO, GOPHER
}
```

<!--
withLog 和 withUpper 都是「接收一個 Handler，回傳一個新的 Handler」。回傳的匿名函式是一個閉包，它捕捉了參數 h，在適當的時機呼叫它。

組合的順序是由外往內：呼叫 h("gopher") 時，先進入 withLog，印出「開始」，然後呼叫 withUpper 包裝的函式，它再呼叫 hello，拿到結果轉大寫，回到 withLog，defer 印出「結束」，最後 main 印出 HELLO, GOPHER。

hello 可以直接傳給 withUpper，因為它的簽章跟 Handler 一樣，Go 會自動轉換。
-->

---

# 章節總結

- **函式基礎**：`func 名稱(參數) 回傳型別`；參數一律傳值；同型別參數可合併
- **多重回傳值**：慣例是「結果 + `error`」；具名回傳值可當文件，naked return 只用在短函式
- **參數不定函式**：`nums ...int` 收到切片；傳入切片要寫 `s...`
- **匿名函式與閉包**：函式是一等公民；閉包會記住並共用它捕捉的變數
- **函式型別**：`type Op func(int, int) int`；函式可當參數、當回傳值
- **泛型與迭代器**：`func F[T any]`（1.18+）；`iter.Seq` 可用於 `for range`（1.23+）
- **defer**：函式結束時執行、LIFO 順序；參數立刻求值；可修改具名回傳值

下一章我們會介紹 Go 的「錯誤處理」：`error` 介面、`panic` 與 `recover`。

<!--
我們來整理今天學到的東西。

函式的部分，最重要的是多重回傳值和「結果加錯誤」的慣例。閉包讓函式可以帶著狀態走，函式型別讓我們可以把「做法」當成參數傳遞。defer 是 Go 處理收尾的標準方式，記得它是後進先出、參數立刻求值。

今天我們已經多次看到 error 這個型別，也看到 errors.New、fmt.Errorf 和 %w。下一章會完整介紹 Go 的錯誤處理哲學：為什麼 Go 不用 try-catch？error 到底是什麼？什麼時候該用 panic？
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的閉包和函式型別是比較抽象的概念，第一次學覺得繞是正常的。

課後建議：把綜合練習的 middleware 範例，再加一個 withTimer，用 defer 印出執行花了多少時間（提示：time.Now() 和 time.Since()）。能寫出來，就代表今天的內容真的掌握了。

有問題的同學現在可以提問！
-->

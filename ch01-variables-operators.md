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
title: 變數與算符
routeAlias: ch01
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">變數與算符</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「先把資料放進盒子，程式才有東西可以算」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第一章！

上一章我們把開發環境準備好了，也跑出了第一支 Hello World。今天開始要正式學 Go 的語法。

任何程式都在做同一件事：把資料存起來、拿出來計算、再把結果存回去。「存起來」靠的是變數，「計算」靠的是算符（operator，也常翻成運算子）。所以這一章是整門課的地基，後面每一章都會用到今天的內容。

除了變數和算符，今天也會碰到 Go 很有特色的幾個概念：零值、指標、常數與 iota 列舉，還有變數的作用範圍。
-->

---
layout: default
---

# Outline

- **前言** — Go 語言簡介、Go 程式長什麼樣子
- **宣告變數** — `var`、型別推斷、短變數宣告 `:=`、多重宣告
- **更改變數值** — 單一賦值、多重賦值與交換
- **算符** — 算術、簡寫、比較與邏輯算符
- **零值** — Go 的變數永遠有初始值
- **值 vs. 指標** — `&`、`*`、以指標為參數的函式
- **常數與列舉** — `const`、`iota`
- **變數作用範圍** — 區塊、遮蔽（shadowing）
- **章節總結**

<!--
這章的內容比較多，但每一段都不難。

我們會先花一點時間認識 Go 這個語言的背景，然後從變數宣告開始，一路學到指標、常數、列舉和作用範圍。

其中「指標」是很多初學者害怕的主題，但 Go 的指標比 C 語言簡單很多，沒有指標運算，大家不用緊張，我們會用很生活化的例子說明。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## Introduction to Go

<!--
在寫程式之前，我們先花幾分鐘認識一下 Go 這個語言：它從哪裡來、為什麼要設計它、現在被用在哪些地方。
-->

---

# Go 語言簡介

| 項目 | 說明 |
| --- | --- |
| **誕生** | 2007 年在 Google 內部開始設計，2009 年公開，2012 年發布 Go 1.0 |
| **設計者** | Robert Griesemer、Rob Pike、Ken Thompson（Unix 與 C 語言的共同作者） |
| **設計動機** | 解決 Google 大型 C++ 專案「編譯慢、相依複雜、難以並行」的問題 |
| **語言特性** | 編譯式、靜態型別、垃圾回收（GC）、內建並行（goroutine） |
| **簡潔** | 只有 25 個關鍵字，沒有類別繼承、沒有例外（exception） |
| **代表作品** | Docker、Kubernetes、Terraform、Prometheus、etcd |

<!--
什麼是 Go？

Go 是 Google 在 2007 年開始設計的語言。當時 Google 內部有很多用 C++ 寫的大型系統，每次編譯要等很久，相依關係又很複雜。三位設計者（其中 Ken Thompson 就是 Unix 和 C 語言的作者之一）決定設計一個新語言：「像 C 一樣快，但像 Python 一樣好寫」。

Go 的特色可以用四個詞形容：編譯式、靜態型別、有垃圾回收、內建並行。而且它非常精簡，整個語言只有 25 個關鍵字，沒有類別繼承、也沒有 try-catch 例外處理，這些後面都會看到 Go 用什麼方式取代。

現在雲端領域的主流工具，像 Docker、Kubernetes，幾乎都是用 Go 寫的。所以學 Go，等於拿到進入雲端和後端開發的門票。
-->

---

# Go 語言的模樣

```go
package main

import "fmt"

func main() {
	name := "Gopher" // 宣告變數，型別自動推斷為 string
	scores := []int{90, 85, 77}

	total := 0
	for _, s := range scores { // 走訪切片
		total += s
	}
	fmt.Printf("%s 的平均分數：%.2f\n", name, float64(total)/float64(len(scores)))
}
```

執行後，console 會輸出：`Gopher 的平均分數：84.00`

<!--
這段程式碼的目的，是讓大家對 Go 的長相先有個整體印象，細節後面每一章都會詳細說明。

幾個觀察重點：第一，行尾不需要分號，Go 編譯器會自動幫我們補上。第二，:= 可以宣告變數並自動推斷型別，寫起來很像動態語言，但其實它是靜態型別。第三，迴圈只有 for 一種，沒有 while。第四，型別轉換要明確寫出來，像 float64(total)，Go 不會偷偷幫我們轉型。

執行後會印出 Gopher 的平均分數是 84.00。
-->

---

# Go 的設計哲學

| 哲學 | 在語法上的體現 |
| --- | --- |
| **明確勝過隱晦** | 不同型別不會自動轉換，必須寫 `float64(x)` |
| **少即是多** | 只有 `for` 一種迴圈；沒有三元運算子 `?:` |
| **程式碼是寫給人看的** | `gofmt` 統一格式；未使用的變數、import 是編譯錯誤 |
| **錯誤是值** | 函式用多重回傳值回傳 `error`，不使用 try-catch |
| **組合勝過繼承** | 沒有 `class` 與 `extends`，用 struct 內嵌與 interface 組合 |

<!--
Go 的語法看起來「少了很多東西」，這其實是刻意的設計。

例如 Go 沒有三元運算子，因為設計者認為 if-else 寫起來更容易閱讀。Go 也不會自動轉換型別，因為隱藏的轉型常常是 bug 的來源。

這張表大家先有印象就好，每個哲學在後面的章節都會遇到。當我們覺得「為什麼 Go 要這樣設計」的時候，回來看這張表，通常就能找到答案。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 宣告變數
## Variables

<!--
接下來進入第一個正式主題：變數。
-->

---

# 什麼是變數？

「變數就是一個**有名字、有型別**的盒子，用來存放資料。」

| 組成 | 說明 | 範例 |
| --- | --- | --- |
| **名稱** | 用來稱呼這個盒子 | `age` |
| **型別** | 決定盒子能放什麼資料 | `int`（整數） |
| **值** | 盒子裡目前放的東西 | `18` |

```go
var age int = 18
```

<!--
什麼是變數？

想像我們在搬家，把東西裝進紙箱，每個紙箱外面寫上「書」、「衣服」、「廚具」。變數就是這樣的紙箱：它有一個名字（紙箱上寫的字），有一個型別（這個箱子是專門裝書的，不能拿來裝湯），還有目前裡面裝的值。

Go 是靜態型別語言，意思是「箱子的用途一旦決定，就不能改」。宣告成 int 的變數，之後永遠只能放整數。這跟 Python、JavaScript 不一樣，一開始可能覺得綁手綁腳，但它能讓編譯器在執行之前就幫我們抓出很多錯誤。
-->

---

# 用 var 宣告變數

語法：`var 變數名稱 型別 = 值`

```go
package main

import "fmt"

var appName string = "GoShop" // 套件層級變數

func main() {
	var price int = 1200
	var rate float64 = 0.95
	fmt.Println(appName, price, rate)
}
```

執行後，console 會輸出：`GoShop 1200 0.95`

<!--
這是最完整的變數宣告寫法：var 關鍵字、變數名稱、型別、等號、值。

注意 Go 的型別寫在變數名稱「後面」，跟 Java、C 相反。Go 的設計者說這樣從左讀到右比較自然：「宣告一個變數 price，型別是 int，值是 1200」。

另外 var 可以寫在函式外面，這叫做套件層級變數，同一個套件裡的所有函式都能使用。
-->

---

# 用 var 一次宣告多個變數

用小括號把多個宣告包在一起，稱為**宣告區塊**：

```go
package main

import "fmt"

var (
	host    string = "localhost"
	port    int    = 8080
	debug   bool   = true
	timeout float64
)

func main() {
	fmt.Println(host, port, debug, timeout) // localhost 8080 true 0
}
```

<!--
當我們有一整組相關的變數，例如伺服器的設定值，可以用 var 加上小括號，一次宣告多個，這樣程式碼比較整齊。

注意 timeout 這個變數沒有給值，執行後印出來是 0。這就是 Go 的「零值」概念：沒有給初始值的變數，Go 會自動給一個預設值。這個概念很重要，等一下會有專門的小節說明。

另外大家看到等號有對齊，那不是我們手動排的，是存檔時 gofmt 自動對齊的。
-->

---

# 用 var 宣告時省略型別或賦值

| 寫法 | 範例 | 結果 |
| --- | --- | --- |
| 完整寫法 | `var a int = 10` | `a` 是 `int`，值為 `10` |
| 省略型別 | `var b = 10` | 由值推斷型別，`b` 是 `int` |
| 省略賦值 | `var c int` | `c` 是 `int`，值為零值 `0` |
| 兩者都省略 | `var d` | ❌ 編譯錯誤：無法得知型別 |

```go
var b = 3.14     // 推斷為 float64
var s = "hello"  // 推斷為 string
var ok = true    // 推斷為 bool
```

<!--
var 的宣告可以省略型別或省略值，但不能兩個都省略。

省略型別的時候，Go 會看等號右邊的值來「推斷」型別：整數推斷成 int，小數推斷成 float64，字串推斷成 string。

省略值的時候，變數會得到零值。

兩個都省略就不行了，因為 Go 完全不知道這個箱子要裝什麼東西。
-->

---

# 推斷型別發生問題的時候

型別推斷只看「值長什麼樣子」，有時候會推斷出**不是我們要的型別**：

```go
var price = 5            // 推斷為 int
total := price * 1.1     // 編譯錯誤：1.1 (untyped float constant) truncated to int
```

解法：**明確寫出型別**，或寫出正確形式的值

```go
var price float64 = 5    // 方法一：明確寫出型別
var price2 = 5.0         // 方法二：寫成小數，推斷為 float64
total := price * 1.1     // ✅ 5.5
```

<!--
型別推斷很方便，但它只看值「長什麼樣子」，不知道我們心裡的打算。

例如價格寫 5，Go 推斷成 int；之後想乘上 1.1 算稅金，就會編譯錯誤，因為整數不能跟小數直接相乘，而 Go 又不會自動轉型。

這時候有兩種解法：明確寫出型別 float64，或者把值寫成 5.0，讓 Go 推斷成 float64。

所以什麼時候該省略型別？原則是：當「值的樣子」剛好就是我們要的型別時才省略；如果不是，就老老實實寫出來。
-->

---

# 短變數宣告 :=

語法：`變數名稱 := 值`，等同於 `var 變數名稱 = 值`

```go
package main

import "fmt"

func main() {
	name := "Alice" // string
	age := 30       // int
	height := 1.68  // float64
	fmt.Println(name, age, height)
}
```

| 比較 | `var` | `:=` |
| --- | --- | --- |
| 能在函式外使用 | ✅ | ❌ 只能在函式內 |
| 能明確指定型別 | ✅ | ❌ 一律推斷 |

<!--
在函式裡面，Go 開發者最常用的其實是短變數宣告，冒號等於。

它的效果跟 var 加上型別推斷完全一樣，只是寫起來更短。實務上，函式裡面有九成的變數都是用 := 宣告的。

但 := 有兩個限制：只能在函式裡面用，函式外面的套件層級變數一定要用 var；而且它一律靠推斷，不能指定型別。所以如果要一個 float64 的 5，就要寫成 5.0，或改用 var。
-->

---

# 以短變數宣告建立多重變數

```go
package main

import "fmt"

func main() {
	x, y, label := 10, 20, "座標"
	fmt.Println(label, x, y)

	x, z := 99, 3 // x 已存在（重新賦值），z 是新變數 ✅
	fmt.Println(x, z)
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>規則：</b> <code>:=</code> 左邊<b>至少要有一個新變數</b>。如果全部都已經宣告過，會出現 <code>no new variables on left side of :=</code>，這時請改用 <code>=</code>。
</div>

<!--
:= 也可以一次宣告多個變數，左邊的名稱跟右邊的值一個對一個。

這裡有一個重要規則：:= 的左邊至少要有一個「新」變數。在第二次的 x, z := 99, 3 裡面，x 已經存在，所以對 x 來說只是重新賦值；z 是新的，所以這行合法。

這個規則在處理錯誤的時候超級常用，例如連續呼叫兩個都回傳 err 的函式，第二次就可以寫 result2, err := ...，重複使用同一個 err 變數。第 6 章會常常看到這個寫法。
-->

---

# 在單行程式內用 var 宣告多重變數

```go
package main

import "fmt"

func main() {
	var a, b, c int = 1, 2, 3           // 同型別，一次宣告
	var name, age, ok = "Bob", 25, true // 省略型別，各自推斷
	var w, h int                        // 都是零值 0
	fmt.Println(a, b, c, name, age, ok, w, h)
}
```

執行後，console 會輸出：`1 2 3 Bob 25 true 0 0`

<!--
var 也能在同一行宣告多個變數。

第一種是同型別的多個變數，型別只要寫一次。第二種是省略型別，每個變數會根據自己的值各自推斷，所以同一行裡可以有不同型別。第三種是只宣告不給值，全部都是零值。

實務上，一行宣告太多變數會降低可讀性，通常兩三個是上限，超過的話就改用 var 區塊。
-->

---

# 非英語的變數名稱

Go 的識別字可以使用 **Unicode 字母**，所以中文變數名稱是合法的：

```go
package main

import "fmt"

func main() {
	名字 := "地鼠"
	年齡 := 3
	fmt.Println(名字, 年齡) // 地鼠 3
}
```

| 命名規則 | 說明 |
| --- | --- |
| 可用字元 | 字母（含 Unicode）、數字、底線 `_`；不能以數字開頭 |
| 大小寫敏感 | `name` 和 `Name` 是兩個不同的變數 |
| 慣例 | 使用英文 **camelCase**，如 `userName`、`maxRetry` |

<!--
Go 的原始碼一律是 UTF-8 編碼，所以變數名稱可以用任何 Unicode 字母，包含中文。這段程式完全可以編譯執行。

但實務上，幾乎沒有人這樣做。原因之一是協作：不是每個同事都用同樣的輸入法。原因之二是 Go 的「大小寫」有特殊意義：名稱開頭大寫代表公開（exported），小寫代表私有。中文字沒有大小寫之分，所以中文開頭的名稱永遠是私有的。這個規則第 8 章講套件的時候會詳細說明。

所以結論是：知道可以，但請用英文 camelCase 命名。
-->

---
layout: default
---

# 練習 1：宣告個人資料
### 任務說明

在 `main()` 中宣告以下變數，並用 `fmt.Println` 印出來：

| 變數 | 型別 | 值 | 要求的宣告方式 |
| --- | --- | --- | --- |
| `name` | `string` | 你的名字 | `var` 完整寫法 |
| `age` | `int` | 你的年齡 | 短變數宣告 |
| `height` | `float64` | `170` | 短變數宣告（注意型別！） |
| `city`、`zip` | `string`、`int` | `"台北"`、`100` | 一行宣告兩個變數 |

<!--
這個練習要同時用到今天學的幾種宣告方式。

特別注意 height 那一格：值是 170，但型別要是 float64，用短變數宣告的時候要怎麼寫？這就是剛剛「推斷型別發生問題」那一頁的內容。
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	var name string = "小明"
	age := 20
	height := 170.0 // 寫成 170.0 才會推斷為 float64
	city, zip := "台北", 100
	fmt.Println(name, age, height, city, zip)
}
```

- 可以用 `fmt.Printf("%T\n", height)` 印出變數的型別，確認是 `float64`

<!--
關鍵在 height：寫 170 會被推斷成 int，要寫成 170.0 才會是 float64。

另外教大家一個確認型別的小技巧：fmt.Printf 搭配 %T，會印出變數的型別。當我們不確定 Go 推斷出什麼型別時，這個技巧很好用。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 更改變數值
## Assignment

<!--
變數宣告好之後，當然可以改變它的值，不然就不叫「變」數了。
-->

---

# 更改單一變數的值

語法：`變數名稱 = 新的值`（用 `=`，不是 `:=`）

```go
package main

import "fmt"

func main() {
	stock := 10
	stock = 8          // 賣出 2 件
	stock = stock + 5  // 進貨 5 件
	fmt.Println(stock) // 13

	// stock = "缺貨"  // 編譯錯誤：型別不同，不能放進去
}
```

<!--
更改變數的值用單純的等號。

注意兩個重點：第一，已經宣告過的變數要用 = 賦值，不是 :=。第二，新的值必須跟原本的型別一樣，stock 是 int，就不能放字串進去。這就是靜態型別的意思：箱子的用途決定了就不能改。

stock = stock + 5 這行在數學上看起來不合理，但在程式裡它的意思是：「先算出右邊 stock + 5 的結果，再放回 stock 這個箱子」。
-->

---

# 一次更改多個變數值

```go
package main

import "fmt"

func main() {
	a, b := 1, 2
	a, b = b, a       // 交換兩個變數的值，不需要暫存變數
	fmt.Println(a, b) // 2 1

	x, y := 0, 0
	x, y = x+10, x+20 // 右邊全部先算完，再一起賦值
	fmt.Println(x, y) // 10 20
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>多重賦值的規則：</b> 等號右邊的所有運算式會<b>先全部算完</b>，再一次指定給左邊的變數。
</div>

<!--
Go 可以一次改多個變數，最經典的用途就是「交換兩個變數的值」。

在很多語言裡，交換兩個變數需要第三個暫存變數，像是左手、右手要交換東西，得先把一樣東西放到桌上。Go 可以直接寫 a, b = b, a，一行搞定。

它的原理是：右邊的所有運算式會先全部算完，再一次指定給左邊。所以第二個例子裡，x+20 用到的 x 還是原本的 0，結果是 10 和 20，而不是 10 和 30。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 算符
## Operators

<!--
有了變數，接下來學怎麼對它們做運算。
-->

---

# 算符基礎：算術算符

| 算符 | 名稱 | 範例（`a := 17`、`b := 5`） | 結果 |
| --- | --- | --- | --- |
| `+` | 加法（也可串接字串） | `a + b` | `22` |
| `-` | 減法 | `a - b` | `12` |
| `*` | 乘法 | `a * b` | `85` |
| `/` | 除法（整數相除取商） | `a / b` | `3` |
| `%` | 取餘數 | `a % b` | `2` |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>整數除法：</b> 兩個整數相除，結果還是整數，小數部分直接捨去（往 0 取整），<code>17 / 5</code> 是 <code>3</code> 而不是 <code>3.4</code>。
</div>

<!--
算術算符就是加減乘除，大家都很熟，重點在除法和取餘數。

兩個整數相除，結果還是整數，小數直接被捨去。想像 17 顆蘋果分給 5 個人，每個人拿 3 顆，剩下的 2 顆就是餘數，用 % 算出來。

% 取餘數在實務上很常用，例如判斷奇偶數（n % 2 == 0），或者計算分頁（第幾頁、每頁幾筆）。
-->

---

# 算術算符 — 範例

```go
package main

import "fmt"

func main() {
	a, b := 17, 5
	fmt.Println(a/b, a%b)                // 3 2
	fmt.Println(17.0 / 5)                // 3.4（有一邊是浮點數）
	fmt.Println(float64(a) / float64(b)) // 3.4（變數要明確轉型）

	first, last := "Go", "pher"
	fmt.Println(first + last) // Gopher（字串串接）
}
```

<!--
這段程式碼的目的，是對照整數除法和浮點數除法的差別。

17.0 / 5 是 3.4，因為 17.0 是浮點數常數。但如果是兩個 int 變數，就必須用 float64() 明確轉型，Go 不會幫我們自動轉。

最後一行示範 + 也可以用來串接字串，Go 加 pher 就變成 Gopher，這也是 Go 吉祥物的名字。
-->

---

# 算符簡寫法

| 簡寫 | 等同於 | 簡寫 | 等同於 |
| --- | --- | --- | --- |
| `x += 5` | `x = x + 5` | `x %= 5` | `x = x % 5` |
| `x -= 5` | `x = x - 5` | `x++` | `x = x + 1` |
| `x *= 5` | `x = x * 5` | `x--` | `x = x - 1` |
| `x /= 5` | `x = x / 5` | `s += "!"` | 字串也適用 |

```go
count := 0
count++          // ✅ 1
count += 10      // ✅ 11
// y := count++  // 編譯錯誤：++ 是敘述，不是運算式
// ++count       // 編譯錯誤：Go 沒有前置遞增
```

<!--
算符簡寫法讓我們少打一點字：x += 5 就是 x = x + 5。

要特別注意的是 ++ 和 --。在 C、Java 裡，i++ 可以放在運算式裡面，還有 ++i 前置跟 i++ 後置的差別，常常讓人搞混。Go 直接把這個問題消除：++ 在 Go 裡是「敘述」不是「運算式」，只能單獨寫一行，而且只有後置寫法。

這又是「少即是多」的哲學：少一個容易出錯的地方。
-->

---

# 值的比較：比較算符

| 算符 | 意義 | 範例 | 結果 |
| --- | --- | --- | --- |
| `==` | 等於 | `3 == 3` | `true` |
| `!=` | 不等於 | `3 != 4` | `true` |
| `<`、`<=` | 小於、小於等於 | `3 < 3`、`3 <= 3` | `false`、`true` |
| `>`、`>=` | 大於、大於等於 | `"b" > "a"` | `true`（字串依位元組比較） |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>型別要一致：</b> <code>int</code> 不能和 <code>float64</code> 直接比較，也不能和 <code>int64</code> 比較，必須先轉成相同型別。
</div>

<!--
比較算符的結果一定是布林值：true 或 false。

字串也可以比大小，它是逐個位元組比較，所以 "b" 大於 "a"，跟字典的排序方式類似。

提醒大家，比較的兩邊型別必須一致。在 Go 裡 int 跟 int64 是不同型別，不能直接比較，要先轉型。這個規則一開始可能覺得麻煩，但它能避免很多隱藏的轉型錯誤。
-->

---

# 邏輯算符

| 算符 | 意義 | 說明 |
| --- | --- | --- |
| `&&` | AND（且） | 兩邊都是 `true` 才是 `true` |
| `\|\|` | OR（或） | 任一邊是 `true` 就是 `true` |
| `!` | NOT（非） | 反轉布林值 |

```go
package main

import "fmt"

func main() {
	age, hasTicket := 20, true
	fmt.Println(age >= 18 && hasTicket) // true：可以入場
	fmt.Println(age < 12 || age > 65)   // false：不適用優惠票
	fmt.Println(!hasTicket)             // false
}
```

<!--
邏輯算符用來組合多個條件。

&& 是「而且」：成年而且有票，才能入場。|| 是「或者」：小於 12 歲或者大於 65 歲，就適用優惠票。! 是「反轉」。

另外 && 和 || 有「短路求值」的特性：&& 左邊如果已經是 false，右邊就不會執行，因為結果一定是 false。這個特性下一章寫 if 判斷時很好用，例如先檢查指標不是 nil，再存取它的欄位。
-->

---

# 補充：位元算符與算符優先順序

| 優先順序（高 → 低） | 算符 |
| --- | --- |
| 5 | `*` `/` `%` `<<` `>>` `&` `&^` |
| 4 | `+` `-` `\|` `^` |
| 3 | `==` `!=` `<` `<=` `>` `>=` |
| 2 | `&&` |
| 1 | `\|\|` |

```go
fmt.Println(2 + 3*4)      // 14：先乘後加
fmt.Println((2 + 3) * 4)  // 20：用括號改變順序
fmt.Println(1 << 3)       // 8：左移 3 位 = 乘以 2³
```

<!--
這一頁是補充內容。

Go 的算符優先順序只有 5 層，比 C 語言的 15 層簡單很多。基本原則跟數學一樣：先乘除、後加減，比較之後才做邏輯運算。

表格裡的 <<、>>、&、| 這些是位元算符，直接對二進位的每個位元做運算，在一般應用程式比較少用到，第 3 章講位元組的時候會再碰到。

實務建議是：當順序不明顯的時候，直接加括號。多寫一對括號不會讓程式變慢，但能讓閱讀的人少想一步。
-->

---
layout: default
---

# 練習 2：計算購物車金額
### 任務說明

有一筆訂單：商品單價 `price := 399`、數量 `qty := 3`、折扣 `discount := 0.9`

1. 計算小計 `subtotal`（單價 × 數量）
2. 計算折扣後金額 `total`（小計 × 折扣），結果為 `float64`
3. 如果 `total` 大於等於 `1000` **而且**數量大於 `2`，免運費，用布林變數 `freeShipping` 表示
4. 用 `+=` 把運費 `60` 加進總額（如果不免運）— 這一步先直接寫 `total += 60`，下一章學完 `if` 再改寫

<!--
這個練習會用到算術算符、型別轉換、比較算符和邏輯算符。

注意第 2 步：subtotal 是 int，discount 是 float64，兩個不能直接相乘，要怎麼辦？
-->

---

# 練習 2：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	price, qty, discount := 399, 3, 0.9
	subtotal := price * qty                  // 1197
	total := float64(subtotal) * discount    // 1077.3
	freeShipping := total >= 1000 && qty > 2 // true
	fmt.Println(subtotal, total, freeShipping)
}
```

- `int` 和 `float64` 不能直接相乘，要用 `float64(subtotal)` 轉型

<!--
關鍵在 float64(subtotal)：Go 不會自動把 int 轉成 float64，必須自己寫出來。

freeShipping 那一行，右邊整個比較運算式的結果是布林值，直接存進變數。執行後會印出 1197、1077.3、true。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 零值
## Zero Values

<!--
接下來是 Go 一個很重要、也很貼心的設計：零值。
-->

---

# 什麼是零值？

「在 Go 裡，**宣告但沒有賦值的變數，會自動得到該型別的零值**，永遠不會是『未初始化』的狀態。」

| 型別 | 零值 |
| --- | --- |
| `int`、`float64` 等數字型別 | `0` |
| `bool` | `false` |
| `string` | `""`（空字串） |
| 指標、切片、map、通道、函式、介面 | `nil` |
| `struct` | 每個欄位都是各自的零值 |

<!--
什麼是零值？

在 C 語言裡，宣告變數沒給值，裡面可能是記憶體裡殘留的垃圾資料，這是很多 bug 的來源。Java 的區域變數沒給值則直接編譯錯誤。

Go 選了第三條路：沒給值，就自動給「零值」。數字是 0、布林是 false、字串是空字串，指標這類參考型別則是 nil（代表「什麼都沒有」）。

這個設計讓 Go 的變數永遠處在一個「可預期」的狀態。後面我們會看到，很多標準函式庫的型別都設計成「零值就能直接使用」，例如 sync.Mutex、bytes.Buffer，宣告完不用初始化就能用。
-->

---

# 零值 — 範例

```go
package main

import "fmt"

func main() {
	var i int
	var f float64
	var b bool
	var s string
	var p *int
	fmt.Printf("%v %v %v %q %v\n", i, f, b, s, p)
}
```

執行後，console 會輸出：`0 0 false "" <nil>`

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>格式化小技巧：</b> <code>%v</code> 印出預設格式；<code>%q</code> 會替字串加上雙引號，空字串才看得出來。
</div>

<!--
這段程式碼的目的，是親眼看看每種型別的零值。

我們用 fmt.Printf 搭配格式化動詞，%v 是「用預設格式印出值」，%q 會替字串加上雙引號。如果用 %v 印空字串，畫面上什麼都看不到；用 %q 才會看到兩個雙引號，確認它是空字串。

指標的零值是 nil，印出來會顯示 <nil>。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 值 vs. 指標
## Values & Pointers

<!--
接下來是很多初學者覺得最難的主題：指標。但別擔心，Go 的指標比 C 語言簡單很多。
-->

---

# 了解指標

「**指標（pointer）是一個存放『記憶體位址』的變數**，它指向另一個變數所在的位置。」

| 比喻 | 程式世界 |
| --- | --- |
| 一棟房子 | 變數（存放值的地方） |
| 房子裡的家具 | 變數的值 |
| 房子的地址 | 記憶體位址 |
| 寫著地址的紙條 | 指標 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>為什麼需要指標？</b> Go 的函式傳參數時一律<b>複製一份值</b>。想讓函式修改「原本那一份」資料，或避免複製很大的資料，就要傳指標。
</div>

<!--
什麼是指標？

想像我們要請朋友到家裡幫忙換燈泡。我們不可能把整棟房子複製一份給他，而是給他一張寫著地址的紙條，他照著地址找到我們家，換的就是我們家裡的燈泡。

指標就是那張紙條：它本身不是資料，而是記錄「資料放在哪裡」的地址。

為什麼需要指標？因為 Go 在呼叫函式的時候，參數一律是「複製一份」傳進去。如果我們希望函式能修改原本的變數，就要傳地址，也就是指標進去。這就是指標最主要的用途。
-->

---

# 取得指標：& 與指標型別

| 語法 | 意義 | 範例 |
| --- | --- | --- |
| `*T` | 「指向 T 的指標」型別 | `*int` 是指向 int 的指標 |
| `&x` | 取得變數 `x` 的位址 | `p := &x` |
| `new(T)` | 建立一個 T 的零值，回傳它的指標 | `p := new(int)` |

```go
package main

import "fmt"

func main() {
	count := 10
	var p *int = &count // p 存放 count 的位址
	q := new(int)       // q 指向一個值為 0 的 int
	fmt.Println(p, q)   // 印出兩個記憶體位址，例如 0xc000012080
}
```

<!--
取得指標有兩種方式。

第一種是用 & 取址符號：&count 的意思是「count 這個變數的地址」，把地址存進 p。p 的型別寫成 *int，讀作「指向 int 的指標」。

第二種是用內建函式 new：new(int) 會建立一個新的 int，值是零值 0，然後回傳它的地址。

直接印出指標，會看到一串 0x 開頭的十六進位數字，那就是記憶體位址。每次執行、每台電腦的結果都不一樣，不用在意它的實際數字。
-->

---

# 從指標取得值：解參考 *

在指標前面加 `*`，就能讀取或修改**指標指向的值**（稱為解參考，dereference）：

```go
package main

import "fmt"

func main() {
	count := 10
	p := &count

	fmt.Println(*p)    // 10：讀取 p 指向的值
	*p = 20            // 透過指標修改
	fmt.Println(count) // 20：原本的變數也改變了！
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>nil 指標：</b> 對值為 <code>nil</code> 的指標解參考，會在執行時期發生 <code>panic: invalid memory address or nil pointer dereference</code>。
</div>

<!--
有了地址，要怎麼拿到房子裡的東西？在指標前面加一個星號，這個動作叫做「解參考」。

*p 就是「p 指向的那個值」，所以印出 10。更重要的是 *p = 20 這行：我們透過地址，直接改了房子裡的東西，所以原本的 count 也變成 20 了。

使用指標的注意事項：如果指標是 nil，也就是紙條上沒有寫地址，這時候去解參考，程式會直接 panic 當掉。這是 Go 最常見的執行期錯誤之一，第 6 章會再詳細介紹 panic。
-->

---

# 採用指標的函式設計

```go
package main

import "fmt"

func doubleValue(n int) { n *= 2 }  // 收到的是「複製品」
func doublePtr(n *int)  { *n *= 2 } // 收到的是「地址」

func main() {
	score := 50
	doubleValue(score)
	fmt.Println(score) // 50：原本的值沒變

	doublePtr(&score)
	fmt.Println(score) // 100：透過指標修改了原本的值
}
```

<!--
這段程式碼的目的，是比較「傳值」和「傳指標」的差別，這是指標最重要的應用。

doubleValue 收到的是 score 的複製品，它把複製品乘以 2，但原本的 score 完全不受影響，就像我們影印了一份文件給別人，他在影本上塗改，原稿不會變。

doublePtr 收到的是地址，它透過 *n 修改的是原本那一份，所以 score 變成 100。

注意呼叫的時候要寫 &score，把地址傳進去。
-->

---

# 使用指標的注意事項

**注意事項之一：** Go 沒有指標運算，不能 `p++` 移到下一個位址（比 C 安全）

**注意事項之二：** 回傳區域變數的指標是**安全的**，Go 會自動把它放到 heap

```go
package main

import "fmt"

func newCounter() *int {
	c := 0
	return &c // ✅ 在 Go 是安全的（C 語言會出問題）
}

func main() {
	p := newCounter()
	*p++
	fmt.Println(*p) // 1
}
```

<!--
使用指標有兩個注意事項。

第一：Go 沒有指標運算。在 C 語言裡可以把指標加一，移動到下一個記憶體位置，這很強大但也很危險。Go 直接禁止，所以 Go 的指標安全很多。

第二：在 C 語言裡，回傳區域變數的地址是很危險的，因為函式結束後那塊記憶體就被回收了。但在 Go 裡完全安全，編譯器會做「逃逸分析」（escape analysis），發現這個變數在函式結束後還會被用到，就自動把它放到 heap 上，交給垃圾回收器管理。

補充一下：*p++ 在這裡是對「p 指向的值」加一，因為 Go 沒有指標運算，所以不會有歧義。
-->

---

# 補充：Go 1.26 起 new 可以接運算式

以前 `new` 只能接型別，想要「指向某個值的指標」得先宣告變數：

```go
package main

import "fmt"

func main() {
	age := 18
	p1 := &age // 傳統寫法：先有變數，再取址

	p2 := new(18)         // Go 1.26+：直接建立指向 18 的指標
	fmt.Println(*p1, *p2) // 18 18
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>實用場景：</b> 結構的欄位是指標型別（例如 JSON 的「可省略欄位」）時，<code>new(18)</code> 可以直接寫在結構常值裡，不用額外宣告暫存變數。
</div>

<!--
這是 Go 1.26 加入的新語法，大家有個印象就好。

以前如果想要一個「指向 18 的指標」，必須先宣告一個變數，再用 & 取地址，因為不能寫 &18（常數沒有地址）。現在 new 可以直接接一個運算式，new(18) 就會建立一個值為 18 的 int，並回傳它的指標。

這在第 11 章處理 JSON 的時候會很好用，屆時會再看到。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 常數與列舉
## Constants & Enums

<!--
有些值在程式執行期間永遠不會改變，例如圓周率、一週有幾天。這種值我們用常數來表示。
-->

---

# 常數 (constants)

語法：`const 名稱 型別 = 值`，常數的值在**編譯時期**就決定，不能修改

```go
package main

import "fmt"

const Pi = 3.14159       // 無型別常數（untyped constant）
const MaxUsers int = 100 // 有型別常數

const (
	AppName = "GoShop"
	Version = "1.0.0"
)

func main() {
	r := 2.0
	fmt.Println(Pi*r*r, MaxUsers, AppName, Version)
	// Pi = 3.14 // 編譯錯誤：cannot assign to Pi
}
```

<!--
常數用 const 宣告，它跟變數最大的差別是：值在編譯時期就決定了，執行期間不能修改。

常數適合放「不會變的東西」，例如數學常數、設定上限、版本號。好處是：第一，不會被不小心改掉；第二，一看到 const 就知道這個值是固定的，閱讀程式更輕鬆。

注意常數的值只能是編譯時期就能算出來的東西，像數字、字串、布林值，不能是函式呼叫的結果。
-->

---

# 無型別常數的彈性

「無型別常數（untyped constant）」在使用時才決定型別，因此可以搭配不同型別：

```go
package main

import "fmt"

const Big = 1 << 40 // 無型別常數，精確度比 int64 還高

func main() {
	var i int64 = Big   // 當成 int64 使用
	var f float64 = Big // 當成 float64 使用
	fmt.Println(i, f)   // 1099511627776 1.099511627776e+12

	const Rate = 0.05
	var price float32 = 100
	fmt.Println(price * Rate) // ✅ Rate 自動配合 float32
}
```

<!--
Go 的常數有一個很聰明的設計：沒有寫型別的常數叫做「無型別常數」，它像一個還沒決定要穿什麼衣服的人，要用的時候才換上適合的衣服。

所以同一個 Big 可以指定給 int64，也可以指定給 float64。Rate 可以跟 float32 相乘，也可以跟 float64 相乘，都不需要轉型。

這解釋了為什麼剛剛 price * 1.1 會出錯：1.1 是無型別的浮點數常數，它想要配合 int 型別的 price，但 1.1 沒辦法變成整數，所以編譯器報錯「truncated to int」。
-->

---

# 列舉 (enums)：iota

Go 沒有 `enum` 關鍵字，而是用 **const 區塊 + `iota`** 來實作列舉：

```go
package main

import "fmt"

type Weekday int // 自訂型別

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1（自動沿用上一行的運算式）
	Tuesday                  // 2
	Wednesday                // 3
)

func main() {
	day := Tuesday
	fmt.Println(day, day == Tuesday) // 2 true
}
```

<!--
什麼是列舉？列舉就是「一組有限、固定的選項」，例如星期一到星期日、訂單狀態（待付款、已出貨、已完成）。

Go 沒有 enum 這個關鍵字，而是用 const 區塊搭配 iota。iota 是一個特別的常數產生器，在 const 區塊裡，它從 0 開始，每往下一行就自動加 1。

而且後面幾行如果沒寫值，會自動沿用上一行的運算式，所以 Monday、Tuesday 不用重複寫 = iota。

另外我們先定義了一個自訂型別 Weekday，讓這些常數有自己的型別，這樣就不會跟一般的 int 搞混。自訂型別第 4 章會正式介紹。
-->

---

# iota 的進階用法

```go
package main

import "fmt"

type Size int

const (
	_       = iota             // 用 _ 跳過 0
	KB Size = 1 << (10 * iota) // 1 << 10
	MB                         // 1 << 20
	GB                         // 1 << 30
)

func main() {
	fmt.Println(KB, MB, GB) // 1024 1048576 1073741824
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>讓列舉從 1 開始：</b> 寫成 <code>Pending Status = iota + 1</code>，就能把 <code>0</code>（零值）保留給「尚未設定」。
</div>

<!--
iota 還可以搭配運算式使用。

第一個例子是儲存容量：用底線跳過第 0 個，接著 KB 是 1 左移 10 位，也就是 1024；MB 自動沿用運算式，iota 變成 2，就是 1 左移 20 位，以此類推。

第二個例子是讓列舉從 1 開始：iota + 1。為什麼要從 1 開始？因為 0 是 int 的零值，如果 Pending 是 0，我們就分不出「狀態是待處理」跟「根本沒設定狀態」。這是實務上很常見的技巧。
-->

---

# 補充：讓列舉印出名稱

替自訂型別加上 `String()` 方法，`fmt` 就會印出可讀的名稱：

```go
package main

import "fmt"

type Status int

const (
	Pending Status = iota + 1
	Shipped
	Done
)

func (s Status) String() string {
	return [...]string{"未知", "待處理", "已出貨", "已完成"}[s]
}

func main() {
	fmt.Println(Shipped) // 已出貨
}
```

<!--
直接印列舉常數，只會看到數字，不太好讀。

我們可以替 Status 這個型別加上一個叫做 String 的方法，fmt 在印出值的時候，會自動呼叫它。這裡用了一個陣列，依照數字對應到中文名稱。

func 後面那個 (s Status) 叫做接收器，代表這個函式是 Status 型別的方法。方法和介面是第 4 章、第 7 章的內容，這裡先看個效果就好。

補充一下：實務上常用官方工具 stringer 自動產生這個方法，不用自己手寫。
-->

---
layout: default
---

# 練習 3：訂單狀態列舉
### 任務說明

1. 定義自訂型別 `type Level int`
2. 用 `const` + `iota` 定義會員等級：`Bronze`、`Silver`、`Gold`，值從 **1** 開始
3. 定義常數 `BaseDiscount = 0.05`
4. 計算 `Gold` 會員的折扣：`float64(Gold) * BaseDiscount`，並印出結果

<!--
這個練習會用到常數、iota 和型別轉換。

想想看：Gold 的型別是 Level，要跟 float64 相乘需要做什麼？
-->

---

# 練習 3：解題提示
### 提示說明

```go
package main

import "fmt"

type Level int

const (
	Bronze Level = iota + 1 // 1
	Silver                  // 2
	Gold                    // 3
)

const BaseDiscount = 0.05

func main() {
	fmt.Printf("%.2f\n", float64(Gold)*BaseDiscount) // 0.15
}
```

<!--
Gold 是 Level 型別，底層是 int，所以要用 float64(Gold) 轉型後才能跟浮點數相乘。

BaseDiscount 是無型別常數，會自動配合 float64，所以不用轉型。

%.2f 表示印出小數點後兩位，執行後會印出 0.15。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 變數作用範圍
## Scope

<!--
最後一個主題：變數的作用範圍，也就是「一個變數在哪裡看得到、在哪裡看不到」。
-->

---

# 什麼是作用範圍？

「**變數的作用範圍，就是它被宣告的那個大括號 `{ }` 區塊**，出了區塊就看不到。」

| 層級 | 宣告位置 | 可見範圍 |
| --- | --- | --- |
| **宇宙區塊** | Go 內建（`int`、`true`、`len`…） | 所有地方 |
| **套件區塊** | 函式外的 `var` / `const` | 同套件的所有檔案 |
| **函式區塊** | 函式內 | 整個函式 |
| **區塊** | `if` / `for` / `switch` / `{ }` 內 | 只在該區塊內 |

<!--
什麼是作用範圍？

想像一棟公司大樓：大廳的公告欄（套件層級）每個人都看得到；某個部門辦公室裡的白板（函式層級），只有那個部門的人看得到；會議室裡的白板（區塊層級），只有在會議室裡的人看得到。

Go 的規則很簡單：變數的作用範圍就是它所在的大括號區塊。內層可以看到外層的變數，但外層看不到內層的變數。
-->

---

# 作用範圍 — 範例

```go
package main

import "fmt"

var company = "GoShop" // 套件層級

func main() {
	dept := "研發部" // 函式層級
	{
		room := "會議室 A"                  // 區塊層級
		fmt.Println(company, dept, room) // ✅ 都看得到
	}
	fmt.Println(company, dept) // ✅
	// fmt.Println(room)       // 編譯錯誤：undefined: room
}
```

<!--
這段程式碼的目的，是示範三個層級的變數在哪裡看得到。

在最內層的大括號裡，三個變數都看得到，因為內層可以看到外層。出了大括號之後，room 就消失了，再使用它會出現 undefined: room 的編譯錯誤。

實務上，單獨寫一對大括號的情況很少，更常見的是 if、for 的區塊，下一章就會大量看到。
-->

---

# 使用作用範圍的注意事項：變數遮蔽

在內層用 `:=` 宣告**同名變數**，會「遮蔽（shadow）」外層的變數：

```go
package main

import "fmt"

func main() {
	total := 100
	if true {
		total := 999       // ⚠️ 這是新的變數，不是修改外層的 total
		fmt.Println(total) // 999
	}
	fmt.Println(total) // 100：外層的 total 沒被改到！
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>避免方式：</b> 想修改外層變數，請用 <code>=</code> 而不是 <code>:=</code>。VS Code 的 gopls 與 <code>go vet</code> 的 shadow 分析器都能幫忙抓出這類問題。
</div>

<!--
使用作用範圍的注意事項：變數遮蔽。這是 Go 初學者最常踩到的坑之一。

在 if 區塊裡，我們本來想修改 total，但不小心寫成 :=，結果 Go 認為我們要「宣告一個新的 total」。這個新的 total 只活在 if 區塊裡，把外層的 total 遮住了。出了區塊，外層的 total 還是 100。

這個錯誤不會有任何編譯錯誤，程式照樣能跑，只是結果不對，所以特別難找。記住口訣：「修改用等號，宣告才用冒號等號」。
-->

---
layout: default
---

# 綜合練習：BMI 計算器
### 任務說明

1. 用 `const` 定義 BMI 的標準上下限：`Low = 18.5`、`High = 24.0`
2. 在 `main()` 用短變數宣告 `weight := 68.0`、`heightCm := 175`
3. 把身高換算成公尺（注意型別轉換），計算 `bmi := 體重 / 身高²`
4. 寫一個函式 `func roundTo1(p *float64)`，透過指標把 BMI 四捨五入到小數點後一位（提示：`math.Round(x*10) / 10`）
5. 用布林變數 `normal` 表示 BMI 是否介於 `Low` 和 `High` 之間，印出結果

<!--
這個綜合練習把整章的東西都串在一起：常數、短變數宣告、型別轉換、算術算符、指標、邏輯算符。

第 4 步是今天最有挑戰性的地方：用指標在函式裡修改外面的變數。先自己試試看。
-->

---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"math"
)

const Low, High = 18.5, 24.0

func roundTo1(p *float64) { *p = math.Round(*p*10) / 10 }

func main() {
	weight, heightCm := 68.0, 175
	h := float64(heightCm) / 100
	bmi := weight / (h * h)
	roundTo1(&bmi)
	normal := bmi >= Low && bmi <= High
	fmt.Println(bmi, normal) // 22.2 true
}
```

<!--
我們一步一步看。

heightCm 是 int，所以先用 float64 轉型再除以 100，得到 1.75 公尺。BMI 是體重除以身高的平方。

roundTo1 收到的是 bmi 的地址，*p 讀出原本的值，四捨五入後再寫回 *p，所以呼叫完之後，main 裡的 bmi 就變成 22.2 了。呼叫時記得寫 &bmi。

最後用 && 組合兩個比較，結果是 true。
-->

---

# 章節總結

- **宣告變數**：`var name T = v` 最完整；函式內最常用 `name := v`；`:=` 左邊至少要有一個新變數
- **型別推斷**：只看值的樣子，`5` 是 `int`、`5.0` 是 `float64`；Go **不會自動轉型**
- **算符**：整數相除會捨去小數；`++` 是敘述不是運算式；`&&`、`||` 會短路求值
- **零值**：沒給值的變數自動是 `0`、`false`、`""`、`nil`，永遠不會是垃圾值
- **指標**：`&x` 取址、`*p` 解參考；函式參數一律傳值，想修改原值就傳指標
- **常數與列舉**：`const` 編譯時期決定；用 `iota` 實作列舉
- **作用範圍**：以 `{ }` 區塊為界；小心 `:=` 造成變數遮蔽

下一章我們會介紹「條件判斷與迴圈」，讓程式可以做決定、重複做事。

<!--
我們來整理今天學到的東西。

變數宣告有 var 和 := 兩種，函式裡面最常用 :=。型別推斷很方便，但要注意 Go 不會自動轉型。算符的部分，記得整數除法會捨去小數，還有 ++ 只能單獨寫一行。

零值讓 Go 的變數永遠處在可預期的狀態。指標是這章最重要的觀念：Go 傳參數一律複製，想修改原本的值就傳指標。常數用 const，列舉用 iota。最後，小心 := 造成的變數遮蔽。

掌握這章之後，我們已經能存資料、算資料了。下一章會學 if、switch 和 for，讓程式能夠根據條件做不同的事，以及重複執行同樣的工作。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的內容是整門課的地基，特別是指標和零值，後面每一章都會用到。

課後建議大家把今天的範例都在 Go Playground 上跑一次，特別是指標那幾個範例，自己改改看數值，親眼確認結果，印象會最深刻。

有問題的同學現在可以提問！
-->

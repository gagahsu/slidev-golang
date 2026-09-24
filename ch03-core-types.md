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
title: 核心型別
routeAlias: ch03
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">核心型別</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「選對箱子的尺寸，資料才裝得剛剛好」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第三章！

前兩章我們一直在用 int、float64、string 這些型別，但還沒有仔細看過它們。今天就要來好好認識 Go 的核心型別：布林值、各種數字、字串，還有 Go 特有的 rune。

為什麼要花一整章講型別？因為型別選錯，程式會出現很難發現的 bug：整數會「溢位」變成負數、浮點數算錢會有誤差、中文字串的長度會算錯。這些都是實務上真的會發生的問題，今天會一個一個拆解。
-->

---
layout: default
---

# Outline

- **前言** — 為什麼型別很重要
- **布林值** — `true` / `false`
- **數字** — 整數、浮點數、溢位與繞回、大數值、位元組
- **字串** — 字串常值、常用操作、`rune` 與 UTF-8
- **nil 值** — 哪些型別的零值是 `nil`
- **章節總結**

<!--
今天的主題依照複雜度排列：布林最簡單，數字的種類最多，字串則牽涉到編碼，最需要花時間理解。

最後的 nil，是 Go 裡代表「什麼都沒有」的特殊值，它跟後面章節的切片、map、指標、錯誤處理都有關係。
-->

---

# 回顧：條件判斷與迴圈

- `if` 可以加**起始賦值**；Go 社群偏好**提早 return**
- `switch` 會自動 break；無條件 `switch` 可以取代長串 `else if`
- `for` 是唯一的迴圈；`for i := range n` 是 Go 1.22 以後的現代寫法
- `for range` 走訪字串時，中文字的**位置會一次跳 3 格** ← 今天揭曉原因

<!--
回顧一下上一章。

我們學了 if、switch、for 三種流程控制。上一章最後有一個懸念：用 for range 走訪「Go語言」這個字串的時候，「語」在位置 2，「言」卻在位置 5，中間跳了 3 格。今天講到字串和 rune 的時候，就會揭曉原因。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## Why Types Matter

<!--
在進入個別型別之前，我們先想想：為什麼 Go 要分這麼多種型別？
-->

---

# 什麼是型別？

「**型別決定了一個值佔用多少記憶體、能放什麼資料、能做哪些運算。**」

| 分類 | 型別 | 本章 |
| --- | --- | --- |
| **布林** | `bool` | ✅ |
| **整數** | `int`、`int8`～`int64`、`uint`、`uint8`～`uint64`、`uintptr` | ✅ |
| **浮點數** | `float32`、`float64` | ✅ |
| **複數** | `complex64`、`complex128` | 較少用，略過 |
| **字串與字元** | `string`、`byte`（= `uint8`）、`rune`（= `int32`） | ✅ |
| **複合型別** | 陣列、切片、map、struct、指標、函式、介面、通道 | Ch 4 起 |

<!--
什麼是型別？

回到第一章「變數是箱子」的比喻：型別就是箱子的規格。小箱子省空間但裝不了大東西，大箱子什麼都裝得下但比較佔位置。int8 就是只能裝 -128 到 127 的小箱子，int64 就是能裝九百京的大箱子。

Go 內建的基本型別就是這張表。複數型別在一般應用程式很少用到，我們略過；複合型別從下一章開始介紹。今天的重點是布林、整數、浮點數和字串。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 布林值
## true / false

<!--
先從最簡單的型別開始：布林值。
-->

---

# 布林值：bool

`bool` 只有兩個值：`true` 和 `false`，零值是 `false`

```go
package main

import "fmt"

func main() {
	isMember := true
	age := 17
	canBuyAlcohol := age >= 18 // 比較運算的結果就是 bool

	fmt.Println(isMember, canBuyAlcohol)    // true false
	fmt.Println(isMember && !canBuyAlcohol) // true
	// fmt.Println(isMember + 1) // 編譯錯誤：bool 不能做數學運算
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>和 C 不一樣：</b> Go 的 <code>bool</code> 不能和整數互相轉換，<code>int(true)</code> 或 <code>if 1 {}</code> 都是編譯錯誤。
</div>

<!--
bool 是最單純的型別，只有 true 和 false 兩個值。上一章所有的比較運算、邏輯運算，結果都是 bool。

提醒大家一個跟 C 語言不一樣的地方：Go 的 bool 跟整數完全不相通。在 C 語言裡 true 就是 1、false 就是 0，可以拿來做加法；在 Go 裡 true 就是 true，不能轉成數字。

命名的慣例上，布林變數常用 is、has、can 開頭，像 isMember、hasTicket、canBuyAlcohol，讀起來就像一個問句，一看就知道是 bool。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 數字
## Numbers

<!--
接下來是數字，這是今天種類最多的部分。
-->

---

# 整數型別

| 型別 | 大小 | 範圍 |
| --- | --- | --- |
| `int8` / `uint8` | 1 byte | -128 ～ 127 / 0 ～ 255 |
| `int16` / `uint16` | 2 bytes | -32,768 ～ 32,767 / 0 ～ 65,535 |
| `int32` / `uint32` | 4 bytes | 約 ±21 億 / 0 ～ 約 42 億 |
| `int64` / `uint64` | 8 bytes | 約 ±922 京 / 0 ～ 約 1,844 京 |
| `int` / `uint` | 依平台：64 位元系統為 8 bytes | 64 位元系統同 `int64` / `uint64` |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>該用哪一個？</b> 沒有特別理由就用 <b><code>int</code></b>。只有在需要固定大小（檔案格式、網路協定、資料庫欄位）或節省大量記憶體時，才選 <code>int32</code>、<code>int64</code> 等。
</div>

<!--
Go 的整數型別有兩個維度：大小（8、16、32、64 位元）和有沒有正負號（int 有號、uint 無號）。

u 開頭的 uint 是 unsigned，無號整數，只能存 0 和正數，但因為不用存正負號，正數的上限多了一倍。

int 和 uint 的大小依平台而定，現在的電腦幾乎都是 64 位元，所以 int 就等於 int64。

那到底該用哪一個？Go 官方的建議很明確：沒有特別理由就用 int。len() 回傳 int、迴圈的索引是 int，用 int 可以避免一堆型別轉換。
-->

---

# 整數 — 範例

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	var small int8 = 100
	var big int64 = 9_000_000_000 // 數字中可以用 _ 分隔，方便閱讀
	count := 42                   // 推斷為 int

	fmt.Println(small, big, count)
	fmt.Println(math.MaxInt8, math.MinInt8, math.MaxInt) // 各型別的上下限
	fmt.Println(0b1010, 0o17, 0xFF)                      // 二、八、十六進位：10 15 255

	// total := big + count      // 編譯錯誤：int64 和 int 是不同型別
	total := big + int64(count) // ✅ 明確轉型
	fmt.Println(total)
}
```

<!--
這段程式碼有幾個實用的小技巧。

第一，數字中間可以加底線，像 9_000_000_000，Go 會忽略底線，但人讀起來輕鬆很多，一眼就知道是 90 億。

第二，math 套件有各種型別的上下限常數，例如 math.MaxInt8 是 127。

第三，數字可以用 0b 開頭寫二進位、0o 開頭寫八進位、0x 開頭寫十六進位。

最後注意：即使在 64 位元的電腦上 int 跟 int64 一樣大，它們在 Go 裡仍然是「不同的型別」，不能直接相加，一定要轉型。
-->

---

# 浮點數型別

| 型別 | 大小 | 有效位數 | 用途 |
| --- | --- | --- | --- |
| `float32` | 4 bytes | 約 7 位 | 圖形處理、大量資料省記憶體 |
| `float64` | 8 bytes | 約 15～16 位 | **預設選擇**，小數推斷的型別 |

```go
package main

import "fmt"

func main() {
	var f32 float32 = 16_777_216
	fmt.Println(f32 + 1) // 1.6777216e+07：float32 精度不夠，+1 消失了

	price := 19.99
	fmt.Printf("%.1f %e\n", price, price) // 20.0 1.999000e+01
}
```

<!--
浮點數就是有小數點的數字，Go 有 float32 和 float64 兩種，預設用 float64。

float32 的有效位數只有大約 7 位，所以 16777216 加 1 之後，因為精度不夠，那個 1 直接消失了。這也是為什麼除非有特殊需求，一律用 float64。

Printf 的 %.1f 表示印到小數點後一位，會四捨五入；%e 是科學記號。
-->

---

# 使用浮點數的注意事項：精度誤差

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	a, b := 0.1, 0.2
	fmt.Println(a + b)      // 0.30000000000000004
	fmt.Println(a+b == 0.3) // false ⚠️

	const eps = 1e-9                     // 比較浮點數時，改用「差距夠小」來判斷
	fmt.Println(math.Abs(a+b-0.3) < eps) // true
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💰 <b>處理金額：</b> 不要用 <code>float64</code> 存錢！改用整數存「最小單位」（例如以「分」為單位的 <code>int64</code>），或使用第三方的 decimal 套件。
</div>

<!--
使用浮點數最重要的注意事項：浮點數有精度誤差。

0.1 加 0.2 竟然不等於 0.3！這不是 Go 的 bug，而是所有使用 IEEE 754 標準的語言都一樣，包括 Java、Python、JavaScript。原因是電腦用二進位存小數，0.1 在二進位是無限循環小數，就像 1/3 在十進位是 0.333… 一樣，只能存近似值。

所以兩個注意事項：第一，比較浮點數不要用 ==，要看兩者的差距是不是夠小。第二，處理金額絕對不要用 float64，改用整數存最小單位，例如新台幣用「元」、美金用「分」，或者使用專門的 decimal 套件。

補充一個有趣的細節：如果寫成 0.1 + 0.2 == 0.3 這種「常數」運算，結果會是 true，因為 Go 的無型別常數在編譯時期是用高精度計算的。誤差只會出現在變數的執行期運算。
-->

---

# 溢位和越界繞回

整數超出型別範圍時：**編譯時期**就能發現的會報錯；**執行時期**發生的則會「繞回」

```go
package main

import "fmt"

func main() {
	// var x int8 = 200   // 編譯錯誤：cannot use 200 (untyped int constant) as int8 value (overflows)

	var u uint8 = 255
	u++
	fmt.Println(u) // 0：超過最大值，繞回最小值

	var i int8 = -128
	i--
	fmt.Println(i) // 127：低於最小值，繞回最大值
}
```

<!--
什麼是溢位？就是數字超過了箱子能裝的範圍。

Go 在兩個時間點處理溢位。如果編譯時期就能看出來，例如直接把 200 放進 int8，編譯器會直接報錯。

但如果是執行時期才發生，例如 255 的 uint8 再加 1，Go 不會報錯，而是「繞回」到最小值 0。就像汽車的里程表，跑到 999999 之後再跑一公里，就變回 000000。

有號整數也一樣，-128 的 int8 再減 1，會繞回到 127。
-->

---

# 使用整數的注意事項：溢位不會報錯

溢位**不會 panic**，程式會帶著錯誤的數字繼續跑，要自己檢查：

```go
package main

import (
	"fmt"
	"math"
)

func safeAdd(a, b int32) (int32, bool) {
	if b > 0 && a > math.MaxInt32-b {
		return 0, false // 會溢位
	}
	return a + b, true
}

func main() {
	fmt.Println(safeAdd(2_100_000_000, 100_000_000)) // 0 false
	fmt.Println(safeAdd(1, 2))                       // 3 true
}
```

<!--
這是整數最危險的地方：溢位不會讓程式當掉，而是帶著一個錯誤的數字繼續執行。

歷史上有很多真實的災難都跟整數溢位有關，例如 1996 年亞利安 5 號火箭就是因為數值轉換溢位而爆炸。

如果我們的程式處理的數字可能很大，就要自己檢查。這個 safeAdd 函式在相加之前先判斷：a 是不是已經大於「上限減 b」，如果是，相加就會溢位。這裡用到了函式的多重回傳值，第 5 章會正式介紹。
-->

---

# 大數值：math/big

超過 `int64` / `uint64` 範圍的數字，使用標準函式庫的 `math/big`：

`big.Int` 任意大小的整數｜`big.Float` 任意精度的浮點數｜`big.Rat` 分數（有理數），例如 1/3

```go
package main

import (
	"fmt"
	"math/big"
)

func main() {
	f := new(big.Int).MulRange(1, 30) // 30! = 1 × 2 × … × 30
	fmt.Println(f)                    // 265252859812191058636308480000000

	p := new(big.Int).Exp(big.NewInt(2), big.NewInt(100), nil) // 2 的 100 次方
	fmt.Println(p)                                             // 1267650600228229401496703205376
}
```

<!--
如果連 int64 的九百京都不夠用呢？例如密碼學、天文計算、階乘，這時候就要用 math/big 套件。

big.Int 可以存任意大小的整數，只受記憶體限制。30 的階乘有 33 位數，int64 早就溢位了，big.Int 輕鬆算出來。

big.Int 的用法跟一般數字不一樣，不能用 + - * /，要呼叫方法。例如 Exp 是次方，MulRange 是連乘。另外注意它們都是用 new(big.Int) 建立，操作的是指標。第 18 章的 RSA 加密，背後就是用 big.Int 在做運算。
-->

---

# 位元組 (Byte)

`byte` 是 `uint8` 的**別名**，專門用來表示「一個位元組的原始資料」

```go
package main

import "fmt"

func main() {
	var b byte = 'A'          // 字元常值 'A' 的 ASCII 碼是 65
	fmt.Println(b, string(b)) // 65 A

	data := []byte("Go!") // 字串轉成位元組切片
	fmt.Println(data)     // [71 111 33]
	data[0] = 'g'
	fmt.Println(string(data)) // go!

	fmt.Printf("%08b\n", b) // 01000001：以二進位檢視
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>哪裡會用到？</b> 讀寫檔案、網路傳輸、加密、JSON 編碼……這些 API 幾乎都以 <code>[]byte</code> 作為輸入與輸出，第 11、12、14、18 章都會大量看到。
</div>

<!--
byte 是 uint8 的別名，兩個名字完全等價。那為什麼要取兩個名字？是為了「表達意圖」：寫 uint8 代表「一個 0 到 255 的小數字」，寫 byte 代表「一個位元組的原始資料」。

電腦裡所有的資料，最終都是一個一個的位元組。所以 []byte，也就是位元組切片，是 Go 裡處理原始資料的通用格式。讀檔案、網路傳輸、加密，拿到的都是 []byte。

字串可以轉成 []byte，轉過去之後就可以修改內容，再轉回字串。字串本身為什麼不能改？下一段會說明。
-->

---
layout: default
---

# 練習 1：年齡與存款
### 任務說明

1. 宣告 `age` 為 `uint8`，值為 `250`，連續加 10 次 1，印出結果並解釋為什麼
2. 宣告存款 `balance` 以「分」為單位存成 `int64`（例如 `12345` 代表 123.45 元）
3. 存入 `0.1` 元十次（每次加 `10` 分），印出 `balance` 並以 `%d.%02d` 格式顯示成「元」
4. 對照組：用 `float64` 把 `0.1` 加十次，印出和 `1.0` 比較的結果

<!--
這個練習把今天數字部分的兩個大坑都走一遍：整數溢位和浮點數誤差。

第 3 步的格式：123.45 元要怎麼從 12345 分算出來？提示：整數除法和取餘數。
-->

---
zoom: 0.9
---

# 練習 1：解題提示
### 提示說明

```go
package main

import "fmt"

func main() {
	var age uint8 = 250
	for range 10 {
		age++
	}
	fmt.Println(age) // 4：超過 255 後繞回 0 再往上加

	var balance int64 = 12345
	for range 10 {
		balance += 10
	}
	fmt.Printf("%d.%02d 元\n", balance/100, balance%100) // 124.45 元

	f := 0.0
	for range 10 {
		f += 0.1
	}
	fmt.Println(f, f == 1.0) // 0.9999999999999999 false
}
```

<!--
第一個結果是 4：250 加到 255 用了 5 次，第 6 次繞回 0，再加 4 次變成 4。

存款用整數「分」計算，結果完全精確；除以 100 得到元，取餘數得到分，%02d 表示不足兩位補 0。

而 float64 加十次 0.1，結果是 0.9999999999999999，不等於 1。這就是為什麼金額要用整數存。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 字串
## Strings & Runes

<!--
接下來是字串。Go 的字串設計跟很多語言不一樣，特別是處理中文的時候，一定要理解它的底層原理。
-->

---

# 字串與字串常值

| 寫法 | 名稱 | 特性 |
| --- | --- | --- |
| `"Hello\n"` | 直譯字串（interpreted） | 支援跳脫字元：`\n` 換行、`\t` Tab、`\"` 雙引號 |
| `` `C:\new\path` `` | 原始字串（raw） | 反引號包住，**所見即所得**，可以跨行，不處理跳脫字元 |

```go
package main

import "fmt"

func main() {
	s1 := "第一行\n第二行\t縮排"
	s2 := `C:\new\table
這是第二行，\n 不會換行`
	fmt.Println(s1)
	fmt.Println(s2)
}
```

<!--
Go 的字串常值有兩種寫法。

雙引號是「直譯字串」，裡面的反斜線有特殊意義，\n 會變成換行、\t 會變成 Tab。

反引號是「原始字串」，裡面寫什麼就是什麼，反斜線不會被處理，而且可以直接換行。它特別適合寫 Windows 的檔案路徑、正規表達式、SQL 指令、JSON 範本這類有很多反斜線或需要多行的內容。第 13 章寫 SQL 的時候就會常常用到。
-->

---

# 字串的本質：唯讀的位元組序列

「**Go 的字串是一段不可修改的位元組序列，內容通常是 UTF-8 編碼的文字。**」

```go
package main

import "fmt"

func main() {
	s := "Go語言"
	fmt.Println(len(s)) // 8：是「位元組數」，不是字數！
	fmt.Println(s[0])   // 71：索引取出的是一個 byte（'G' 的編碼）
	fmt.Println(s[0:2]) // Go：切片取出的是子字串

	// s[0] = 'g'        // 編譯錯誤：cannot assign to s[0] (字串不可修改)
	s = "go" + s[2:] // 要修改就建立新字串
	fmt.Println(s)   // go語言
}
```

<!--
這一頁是理解 Go 字串的關鍵。

Go 的字串，本質上是一串「位元組」，而且是唯讀的。「Go語言」這個字串，G 和 o 各佔 1 個位元組，「語」和「言」在 UTF-8 裡各佔 3 個位元組，所以 len 是 8，不是 4。

用索引取出的也是一個位元組，所以 s[0] 是 71，也就是 G 的編碼。

字串不能修改，s[0] = 'g' 會編譯錯誤。想要修改，就用字串相加或切片的方式「建立一個新的字串」。唯讀的好處是：字串可以安全地在多個地方共用，不用擔心被別人偷偷改掉。
-->

---
zoom: 0.9
---

# 常用的字串操作：strings 套件

| 函式 | 說明 | 範例 → 結果 |
| --- | --- | --- |
| `strings.Contains(s, sub)` | 是否包含 | `("golang", "go")` → `true` |
| `strings.HasPrefix` / `HasSuffix` | 開頭／結尾是否符合 | `("main.go", ".go")` → `true` |
| `strings.Index(s, sub)` | 第一次出現的位置 | `("chicken", "ken")` → `4` |
| `strings.Split(s, sep)` | 切割成切片 | `("a,b,c", ",")` → `[a b c]` |
| `strings.Join(elems, sep)` | 用分隔符號串接 | `([]string{"a","b"}, "-")` → `a-b` |
| `strings.ToUpper` / `ToLower` | 轉大寫／小寫 | `("Go")` → `GO` |
| `strings.TrimSpace(s)` | 去除前後空白 | `("  hi \n")` → `hi` |
| `strings.ReplaceAll(s, old, new)` | 全部取代 | `("a-b-c", "-", "+")` → `a+b+c` |
| `strings.Cut(s, sep)` | 在第一個分隔符號切成兩半 | `("k=v", "=")` → `k v true` |

<!--
字串處理是最常見的工作，Go 把常用的功能都放在 strings 套件裡。

這張表是最常用的九個函式，大家不用背，需要的時候回來查就好。特別推薦最後一個 strings.Cut，它是 Go 1.18 加入的，用來把「key=value」這種字串切成兩半，比先 Index 再切片方便很多，現在是主流寫法。
-->

---

# 常用的字串操作 — 範例

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	line := "  name=Gopher, age=15  "
	for part := range strings.SplitSeq(strings.TrimSpace(line), ",") { // Go 1.24+
		key, value, _ := strings.Cut(strings.TrimSpace(part), "=")
		fmt.Printf("%s → %s\n", key, value)
	}

	n, _ := strconv.Atoi("15") // 字串 → 整數
	s := strconv.Itoa(n + 1)   // 整數 → 字串
	fmt.Println(s + "歲")       // 16歲
	// fmt.Println(string(65))  // ⚠️ 得到 "A" 而不是 "65"，數字轉字串請用 strconv
}
```

<!--
這段程式碼的目的，是解析一行「key=value」格式的設定字串。

先用 TrimSpace 去掉前後空白，再用 SplitSeq 以逗號切開。SplitSeq 是 Go 1.24 新增的，它跟 Split 效果一樣，但不會先建立一整個切片，而是一個一個交給 for range，比較省記憶體。接著每一段再用 Cut 切成 key 和 value。

下半段是字串和數字的互相轉換，要用 strconv 套件：Atoi 把字串轉整數，Itoa 把整數轉字串。

注意最後一行註解：string(65) 不會得到 "65"，而是得到編碼 65 的字元 "A"。這是很常見的錯誤，go vet 也會警告。數字轉字串一定要用 strconv。
-->

---

# 補充：大量串接字串用 strings.Builder

字串不可修改，用 `+=` 在迴圈中串接會**反覆建立新字串**，效率很差：

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	var sb strings.Builder // 零值就能直接使用
	for i := range 5 {
		fmt.Fprintf(&sb, "第%d行;", i+1)
	}
	sb.WriteString("結束")
	fmt.Println(sb.String()) // 第1行;第2行;第3行;第4行;第5行;結束
}
```

<!--
因為字串不能修改，每次 s += "..." 都會建立一個全新的字串，把舊內容複製過去。迴圈跑一萬次，就複製了一萬次，越來越慢。

strings.Builder 內部用一個可以擴充的緩衝區，串接時不會一直複製，效率好很多。它的零值就可以直接使用，不需要初始化，這就是第一章說的「零值可用」的設計。

fmt.Fprintf 可以把格式化的結果寫進 Builder 裡。原則是：串接幾個字串用 + 就好，在迴圈裡大量串接才用 Builder。
-->

---

# Rune：一個 Unicode 字元

`rune` 是 `int32` 的別名，代表一個 **Unicode 碼位（code point）**，用**單引號**表示

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	var r rune = '語'
	fmt.Println(r, string(r))   // 35486 語
	fmt.Printf("%U %c\n", r, r) // U+8A9E 語

	s := "Go語言"
	fmt.Println(len(s), utf8.RuneCountInString(s)) // 8 4：位元組數 vs 字元數
	runes := []rune(s)
	fmt.Println(len(runes), string(runes[2])) // 4 語
}
```

<!--
什麼是 rune？rune 在英文是「符文」的意思，在 Go 裡代表「一個 Unicode 字元」。

Unicode 替全世界每個字元都編了一個號碼，「語」的號碼是 U+8A9E，也就是十進位的 35486。rune 就是用來存這個號碼的型別，它是 int32 的別名。

注意字元用單引號，字串用雙引號。'語' 是一個 rune，"語" 是一個字串。

想要知道字串有幾個「字」，要用 utf8.RuneCountInString，而不是 len。如果想用索引取出第幾個「字」，就先轉成 []rune，每個元素就是一個完整的字元。
-->

---
zoom: 0.95
---

# UTF-8 編碼：為什麼中文佔 3 個位元組？

| 字元 | Unicode 碼位 | UTF-8 位元組數 | UTF-8 編碼（十六進位） |
| --- | --- | --- | --- |
| `G` | U+0047 | 1 | `47` |
| `é` | U+00E9 | 2 | `C3 A9` |
| `語` | U+8A9E | 3 | `E8 AA 9E` |
| `🐹` | U+1F439 | 4 | `F0 9F 90 B9` |

```go
package main

import "fmt"

func main() {
	for i, ch := range "Go語言" { // range 字串：自動以 UTF-8 解碼成 rune
		fmt.Printf("位置 %d：%c\n", i, ch) // 位置 0、1、2、5
	}
}
```

<!--
現在揭曉上一章的懸念。

UTF-8 是一種「變長」編碼：英文字母只要 1 個位元組，歐洲語言的重音字母 2 個，中文 3 個，表情符號 4 個。這樣設計的好處是英文文件很省空間，而且跟古老的 ASCII 完全相容。

for range 走訪字串的時候，Go 會自動用 UTF-8 解碼，每一圈拿到一個完整的 rune，而索引 i 是這個字元「在位元組序列中的起始位置」。「語」從位置 2 開始，佔了 2、3、4 三個位元組，所以「言」從位置 5 開始。這就是為什麼位置會跳 3 格。

Go 的原始碼本身就規定是 UTF-8 編碼，Go 的設計者 Rob Pike 和 Ken Thompson 也正是 UTF-8 的發明者。
-->

---

# 使用字串的注意事項

**注意事項之一：** 要以「字」為單位處理中文時，用 `for range` 或 `[]rune`，**不要用索引 `s[i]`**

**注意事項之二：** 用位元組切片取子字串，可能會把中文字**切壞**

```go
package main

import "fmt"

func main() {
	s := "你好世界"
	fmt.Println(s[:4])                 // 你�：切到「好」的中間，變成亂碼
	fmt.Println(string([]rune(s)[:2])) // 你好：以 rune 為單位切
}
```

<!--
使用字串有兩個注意事項，都跟中文有關。

第一：s[i] 取出的是一個位元組，不是一個字。處理英文沒問題，處理中文就會出錯。要一個字一個字處理，用 for range 或先轉成 []rune。

第二：用切片取子字串時，數字代表的是位元組位置。「你好世界」的 s[:4]，會切到「好」這個字的第一個位元組，印出來就變成亂碼。正確做法是先轉成 []rune 再切。

轉成 []rune 需要額外的記憶體，所以只在真的需要以字元為單位操作時才轉。
-->

---
layout: default
---

# 練習 2：字串反轉
### 任務說明

寫一個函式 `reverse(s string) string`，把字串反轉：

| 輸入 | 輸出 |
| --- | --- |
| `"Hello"` | `"olleH"` |
| `"Go語言"` | `"言語oG"` |
| `"🐹地鼠"` | `"鼠地🐹"` |

提示：如果直接反轉位元組，中文會變成亂碼，要怎麼辦？

<!--
字串反轉是很經典的練習，但在 Go 裡處理中文的時候有一個陷阱。

先試試看用位元組的方式反轉，看看中文會發生什麼事，再想想看怎麼修正。
-->

---

# 練習 2：解題提示
### 提示說明

```go
package main

import "fmt"

func reverse(s string) string {
	r := []rune(s) // 先轉成 rune 切片，一個元素就是一個完整的字
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i] // 頭尾交換
	}
	return string(r)
}

func main() {
	fmt.Println(reverse("Hello"), reverse("Go語言"), reverse("🐹地鼠"))
}
```

<!--
關鍵是先轉成 []rune，這樣每個元素都是一個完整的字元，反轉就不會切壞中文。

迴圈用了兩個變數 i 和 j，一個從頭、一個從尾往中間走，每一圈交換兩邊的字元，直到相遇。這裡用到了第一章的多重賦值，一行就完成交換。

執行後會印出 olleH 言語oG 鼠地🐹。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# nil 值
## The Zero Value of Reference Types

<!--
最後一個主題：nil。
-->

---

# 什麼是 nil？

「**`nil` 是指標、切片、map、通道、函式、介面這六種型別的零值，代表『什麼都沒有』。**」

| 型別 | nil 時的行為 |
| --- | --- |
| 指標 `*T` | 解參考會 **panic** |
| 切片 `[]T` | 可以讀 `len`（為 0）、可以 `append`、可以 `range` ✅ |
| map | 可以讀（得到零值）、**寫入會 panic** |
| 通道 `chan T` | 收送都會永遠阻塞 |
| 函式 | 呼叫會 **panic** |
| 介面 | 呼叫方法會 **panic** |

<!--
什麼是 nil？

nil 就是「什麼都沒有」。第一章我們說過，每種型別都有零值：數字是 0、字串是空字串，而指標、切片、map、通道、函式、介面這六種型別的零值，就是 nil。

注意 int、string、bool、struct 這些型別不能是 nil，Go 的字串永遠不會是 null，這讓 Go 少了很多其他語言常見的 NullPointerException。

不同型別的 nil，行為不太一樣。最友善的是 nil 切片，它可以直接 append、可以 range，就像一個空切片；最危險的是 nil 指標和 nil map 的寫入，會直接讓程式 panic。這些型別後面的章節都會一一介紹，今天先建立一個全貌。
-->

---

# nil — 範例

```go
package main

import "fmt"

func main() {
	var p *int
	var s []int
	var m map[string]int

	fmt.Println(p == nil, s == nil, m == nil) // true true true

	s = append(s, 1)       // ✅ nil 切片可以直接 append
	fmt.Println(s, m["x"]) // [1] 0：讀取 nil map 得到零值

	if p != nil { // ✅ 使用指標前先檢查
		fmt.Println(*p)
	}
	// m["x"] = 1 // ⚠️ panic: assignment to entry in nil map
}
```

<!--
這段程式碼的目的，是示範 nil 的安全用法和危險用法。

三個變數都只宣告沒給值，所以都是 nil。nil 切片可以直接 append，這在實務上很常用：宣告一個空的切片變數，然後一個一個 append 進去。nil map 可以讀，會得到零值。

危險的操作是：對 nil 指標解參考，以及對 nil map 寫入，都會 panic。所以使用指標之前，養成先檢查是不是 nil 的習慣。map 要先用 make 建立才能寫入，下一章會教。
-->

---
layout: default
---

# 綜合練習：文字統計器
### 任務說明

寫一個函式 `stats(s string)`，統計字串中各類字元的數量並印出：

1. 總位元組數（`len`）與總字元數
2. 英文字母數、數字數、空白數、中文字數（提示：`unicode.IsLetter`、`unicode.IsDigit`、`unicode.IsSpace`、`unicode.Is(unicode.Han, r)`）
3. 把字串轉成大寫後印出（`strings.ToUpper`）
4. 用 `"Go 1.27 發布了 Hello 世界"` 測試

<!--
這個綜合練習把今天學的字串、rune、for range 都串起來。

unicode 套件提供了很多判斷字元種類的函式，unicode.Han 代表漢字。注意判斷英文字母時，IsLetter 對中文也會回傳 true，所以要先判斷是不是漢字。
-->

---
zoom: 0.76
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func stats(s string) {
	var letters, digits, spaces, han int
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Han, r): // 先判斷漢字
			han++
		case unicode.IsLetter(r):
			letters++
		case unicode.IsDigit(r):
			digits++
		case unicode.IsSpace(r):
			spaces++
		}
	}
	fmt.Println("位元組：", len(s), "字元：", utf8.RuneCountInString(s))
	fmt.Println("英文：", letters, "數字：", digits, "空白：", spaces, "中文：", han)
	fmt.Println(strings.ToUpper(s))
}

func main() { stats("Go 1.27 發布了 Hello 世界") }
```

<!--
用 for range 一個字元一個字元走訪，再用無條件 switch 分類。

順序很重要：漢字也算是 Letter，所以一定要先判斷漢字，再判斷英文字母。

執行後，位元組數 30、字元數 20，英文 7 個、數字 3 個、空白 4 個、中文 5 個。ToUpper 只會把英文字母轉大寫，中文不受影響。
-->

---

# 章節總結

- **bool**：只有 `true` / `false`，不能和整數互轉
- **整數**：沒特別理由就用 `int`；不同整數型別不能直接運算；執行期溢位會**繞回而不會報錯**
- **浮點數**：預設 `float64`；有精度誤差，**不要用 `==` 比較、不要用來存錢**
- **大數值**：超出 `int64` 用 `math/big`；`byte` 是 `uint8` 的別名，`[]byte` 是原始資料的通用格式
- **字串**：唯讀的 UTF-8 位元組序列；`len` 是位元組數；處理中文用 `for range` 或 `[]rune`
- **rune**：`int32` 的別名，代表一個 Unicode 字元
- **nil**：指標、切片、map、通道、函式、介面的零值；nil 切片可 `append`，nil map 不可寫入

下一章我們會介紹「複合型別」：陣列、切片、map、struct 與介面。

<!--
我們來整理今天學到的東西。

整數沒特別理由就用 int，要小心執行期溢位會默默繞回。浮點數有精度誤差，比較時看差距、存錢用整數。字串是唯讀的 UTF-8 位元組序列，len 回傳位元組數，處理中文要用 rune。nil 是六種參考型別的零值。

今天學的都是「單一的值」。下一章要學「複合型別」，也就是怎麼把很多個值組合在一起：陣列、切片、map 用來裝一堆同型別的資料，struct 用來把不同型別的欄位組成一筆資料。這些是 Go 程式裡最常用的資料結構。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的內容裡，最重要的是浮點數的精度誤差和字串的 UTF-8 編碼，這兩個在實務上都很常踩坑。

課後建議：試著用 len 和 utf8.RuneCountInString 量量看自己的名字，再試試看各種表情符號，體會一下 UTF-8 的變長編碼。

有問題的同學現在可以提問！
-->

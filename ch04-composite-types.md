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
title: 複合型別
routeAlias: ch04
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">複合型別</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「把一堆資料裝在一起，程式才有辦法處理真實世界」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第四章！

上一章我們認識了 Go 的核心型別：數字、字串、布林值，它們每一個都只能存「一個值」。但真實世界的資料很少是單獨一個值：一個班級有三十個學生的成績、一個商店有上百種商品、一個會員有姓名、電話、生日好幾個欄位。

今天要學的複合型別，就是把很多個值組合在一起的方法：陣列、切片、map 用來裝「一堆」資料，struct 用來把「不同欄位」組成一筆資料。其中切片和 map 是 Go 程式裡使用頻率最高的資料結構，幾乎每一支程式都會用到。
-->

---
layout: default
---

# Outline

- **集合型別** — 陣列、切片、map 的差別
- **陣列 (array)** — 固定長度、值語意
- **切片 (slice)** — `append`、切割、底層陣列、`slices` 套件
- **映射表 (map)** — 讀取、刪除、走訪、排序
- **自訂型別** — `type` 定義新型別
- **結構 (struct)** — 定義、比較、內嵌、方法
- **介面與型別檢查** — 型別轉換、型別斷言、型別 switch
- **GoShop 專案實作** — 第 4 步：商品結構、購物車與商品目錄
- **章節總結**

<!--
今天的內容是整門課最「密」的一章，因為複合型別是後面所有章節的基礎。

我們會先比較三種集合型別：陣列、切片、map。然後學自訂型別和 struct，這是 Go 用來組織資料的方式，相當於其他語言的「類別」。最後會初步認識介面和型別斷言，第 7 章會再深入。
-->

---

# 回顧：核心型別

- 整數沒特別理由就用 `int`；浮點數預設 `float64`，**不要用來存錢**
- 字串是**唯讀**的 UTF-8 位元組序列；`len` 回傳位元組數
- `rune` 代表一個 Unicode 字元；處理中文用 `for range` 或 `[]rune`
- `nil` 是指標、**切片**、**map**、通道、函式、介面的零值

<!--
回顧一下上一章。

我們學了數字、字串、rune 和 nil。上一章其實已經偷偷用過 []byte 和 []rune，它們就是今天要學的「切片」。另外 nil 的那張表裡，提到了 nil 切片可以 append、nil map 不能寫入，今天會正式說明原因。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 集合型別
## Collection Types

<!--
首先，我們來看看 Go 有哪三種集合型別，以及它們的差別。
-->

---

# 三種集合型別

| 型別 | 寫法 | 長度 | 以什麼存取 | 生活比喻 |
| --- | --- | --- | --- | --- |
| **陣列** | `[5]int` | 固定，寫在型別裡 | 索引 `0`～`n-1` | 雞蛋盒：6 格就是 6 格 |
| **切片** | `[]int` | 可以變長 | 索引 `0`～`n-1` | 購物清單：可以一直往下加 |
| **map** | `map[string]int` | 可以變長 | 鍵（key） | 電話簿：用名字查電話 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>實務上的使用頻率：</b> 切片 ≫ map ≫ 陣列。陣列很少直接使用，但它是切片的「底層」，理解陣列才能真正理解切片。
</div>

<!--
Go 有三種內建的集合型別。

陣列像雞蛋盒：買回來是 6 格就是 6 格，不能變多也不能變少。切片像購物清單：想到要買什麼就往下加一行。map 像電話簿：不是用「第幾個」查，而是用「名字」查電話。

實務上最常用的是切片，其次是 map，陣列很少直接使用。但是切片的底層其實就是陣列，所以我們還是從陣列開始學，後面理解切片的行為時會輕鬆很多。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 陣列
## Arrays

<!--
先從陣列開始。
-->

---

# 定義一個陣列

語法：`[長度]型別`，**長度是型別的一部分**，`[3]int` 和 `[4]int` 是不同的型別

```go
package main

import "fmt"

func main() {
	var scores [3]int // 零值：[0 0 0]
	primes := [5]int{2, 3, 5, 7, 11}
	days := [...]string{"一", "二", "三"} // ... 讓編譯器自己數長度

	// [0 0 0] [2 3 5 7 11] [一 二 三]
	fmt.Println(scores, primes, days)
	fmt.Println(len(primes), len(days)) // 5 3
}
```

<!--
陣列的型別寫法是中括號裡寫長度，後面接元素的型別。

有三種常見的建立方式：只宣告，會得到全部都是零值的陣列；用大括號直接給值；如果不想自己數有幾個元素，可以在中括號裡寫三個點，讓編譯器幫我們數。

重點是：長度是型別的一部分。[3]int 和 [4]int 在 Go 裡是完全不同的型別，不能互相指定，也不能傳給同一個函式。這就是陣列不太好用、實務上很少直接用的原因。
-->

---

# 透過索引鍵賦值

建立陣列時，可以用 `索引: 值` 的方式**只指定部分元素**：

```go
package main

import "fmt"

func main() {
	a := [5]int{1: 10, 3: 30} // 其他位置是零值
	fmt.Println(a)            // [0 10 0 30 0]

	b := [...]string{2: "c", 0: "a"} // 長度由最大索引決定
	fmt.Println(len(b), b)           // 3 [a  c]

	week := [7]bool{5: true, 6: true} // 標記週末
	// [false false false false false true true]
	fmt.Println(week)
}
```

<!--
建立陣列時，可以寫「索引: 值」，只指定某幾個位置，其他位置自動是零值。

這個寫法在陣列很大、但只有少數位置有值的時候很方便。第二個例子搭配三個點，長度會由最大的索引決定：最大索引是 2，所以長度是 3。

中間那個空的元素是空字串，所以印出來 a 和 c 中間有兩個空格。
-->

---

# 讀取、寫入與比較陣列

```go
package main

import "fmt"

func main() {
	temps := [3]float64{25.5, 27.0, 30.2}
	fmt.Println(temps[0], temps[len(temps)-1]) // 讀取：25.5 30.2

	temps[1] = 28.3    // 寫入
	fmt.Println(temps) // [25.5 28.3 30.2]

	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}
	fmt.Println(a == b) // true：長度相同、元素都相等

	// temps[3] = 0 // 編譯錯誤：index 3 out of bounds
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>陣列可以用 <code>==</code> 比較：</b> 條件是型別相同（長度與元素型別都一樣），且元素本身可以比較。
</div>

<!--
讀取和寫入陣列都用中括號加索引，索引從 0 開始，最後一個元素的索引是 len - 1。

陣列可以直接用 == 比較，Go 會逐個元素比對。但前提是型別要相同，[3]int 不能跟 [4]int 比較，因為它們是不同的型別。

索引超出範圍的話：如果編譯時期就看得出來，會直接編譯錯誤；如果是執行時期才發生，例如索引是變數，就會 panic。
-->

---

# 走訪陣列與陣列的「值語意」

```go
package main

import "fmt"

func main() {
	scores := [4]int{80, 92, 67, 75}
	sum := 0
	for _, s := range scores {
		sum += s
	}
	fmt.Println("平均：", sum/len(scores)) // 平均： 78

	copyOf := scores // ⚠️ 陣列指定時會「整個複製」一份
	copyOf[0] = 0
	fmt.Println(scores[0], copyOf[0]) // 80 0：原本的陣列不受影響
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>值語意：</b> 陣列指定給另一個變數、或傳進函式時，會<b>複製所有元素</b>。陣列很大時，這個複製的成本很高。
</div>

<!--
走訪陣列用上一章學的 for range。

這一頁的重點是下半段：陣列是「值」。把陣列指定給另一個變數，會把所有元素完整複製一份，兩個變數各自獨立，改其中一個不會影響另一個。傳進函式也一樣會複製。

這跟 Java、Python 很不一樣，在那些語言裡陣列是參考，指定只是多一個名字指向同一份資料。Go 的陣列比較像一張紙：影印一份給別人，別人在影本上畫圖，我們的原稿不受影響。

等一下學切片，會看到切片的行為剛好相反。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 切片
## Slices

<!--
接下來是今天的重頭戲：切片。這是 Go 裡使用頻率最高的資料結構。
-->

---

# 什麼是切片？

「**切片是一個『可以變長的陣列視窗』**：它本身不存資料，而是指向一個底層陣列的某一段。」

| 建立方式 | 範例 | 結果 |
| --- | --- | --- |
| 切片常值 | `s := []int{1, 2, 3}` | `[1 2 3]`，長度 3 |
| `make` | `s := make([]int, 3)` | `[0 0 0]`，長度 3 |
| `make` 預留容量 | `s := make([]int, 0, 10)` | `[]`，長度 0、容量 10 |
| 宣告 | `var s []int` | `nil` 切片，長度 0 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>型別寫法：</b> 陣列是 <code>[3]int</code>（中括號裡有長度），切片是 <code>[]int</code>（中括號是空的）。
</div>

<!--
什麼是切片？

我們可以把切片想像成一扇「窗戶」，窗戶後面是一整排的陣列。切片本身不存資料，它只記錄「從陣列的哪裡開始看、看幾格、最多可以看到哪裡」。

型別的寫法跟陣列只差一點：中括號裡面沒有長度。因為切片的長度可以變，所以長度不是型別的一部分。

建立切片最常見的方式是切片常值和 make。make 的第二個參數是長度，第三個參數是容量，容量的意義等一下會詳細說明。
-->

---
zoom: 0.88
---

# 使用切片：append、len、cap

```go
package main

import "fmt"

func main() {
	var cart []string // nil 切片
	cart = append(cart, "蘋果")
	cart = append(cart, "香蕉")
	fmt.Println(cart, len(cart)) // [蘋果 香蕉] 2

	cart[0] = "芒果" // 用索引修改
	for i, item := range cart {
		fmt.Println(i, item)
	}
}
```

| 函式 | 說明 |
| --- | --- |
| `append(s, v...)` | 在尾端加入元素，**回傳新的切片**，一定要接住：`s = append(s, v)` |
| `len(s)` | 目前的元素個數 |
| `cap(s)` | 底層陣列從切片起點算起的總容量 |

<!--
切片最常用的操作就是 append：在尾端加入元素。

特別注意 append 的寫法：cart = append(cart, "蘋果")。append 會「回傳」一個新的切片，我們一定要把結果接回原本的變數。為什麼？因為當容量不夠的時候，append 會換一個更大的底層陣列，回傳的切片指向新陣列，舊的變數還指向舊陣列。等一下講切片內部運作的時候會看得很清楚。

另外 nil 切片可以直接 append，不用先 make，這就是上一章說的「nil 切片很友善」。
-->

---

# 為切片附加多重元素

```go
package main

import "fmt"

func main() {
	nums := []int{1, 2}
	nums = append(nums, 3, 4, 5) // 一次加入多個元素
	fmt.Println(nums)            // [1 2 3 4 5]

	more := []int{6, 7}
	nums = append(nums, more...) // 用 ... 把切片「展開」
	fmt.Println(nums)            // [1 2 3 4 5 6 7]

	// 特例：字串可以展開附加到 []byte
	b := append([]byte("Go"), "語言"...)
	fmt.Println(string(b)) // Go語言
}
```

<!--
append 可以一次加入多個元素，用逗號隔開就好。

如果要把另一個切片整個接上去，要在切片後面加三個點，意思是「把這個切片展開，每個元素當成一個獨立的參數」。這三個點的語法第 5 章學「參數不定函式」時會再詳細說明。

最後一行是一個特例：[]byte 可以直接 append 一個字串展開，這在處理位元組資料時很方便。
-->

---

# 從切片和陣列建立新的切片

語法：`s[low:high]`，取出索引 `low` 到 `high-1` 的元素（**包含 low，不包含 high**）

```go
package main

import "fmt"

func main() {
	arr := [6]string{"a", "b", "c", "d", "e", "f"}
	s1 := arr[1:4] // [b c d]
	s2 := arr[:2]  // [a b]：省略 low 代表從 0 開始
	s3 := arr[3:]  // [d e f]：省略 high 代表到最後
	s4 := s1[1:]   // [c d]：切片也可以再切
	fmt.Println(s1, s2, s3, s4)
	fmt.Println(len(s1), cap(s1)) // 3 5
}
```

<!--
我們可以從陣列或切片「切」出一段，形成新的切片，語法是中括號裡寫起點和終點，用冒號隔開。

要記住的規則是「包含起點、不包含終點」，arr[1:4] 取的是索引 1、2、3 三個元素。這樣設計的好處是：長度剛好等於 high 減 low。

最後一行 cap(s1) 是 5，為什麼不是 3？因為容量是「從切片的起點算到底層陣列的尾端」，s1 從索引 1 開始，到陣列尾端還有 5 格。這個概念下一頁會畫圖說明。
-->

---

# 了解切片的內部運作

切片本身是一個很小的結構，只有三個欄位：

| 欄位 | 意義 |
| --- | --- |
| **指標** | 指向底層陣列中，切片的第一個元素 |
| **長度 `len`** | 切片目前能看到幾個元素 |
| **容量 `cap`** | 從起點到底層陣列尾端，最多能擴展到幾個元素 |

```text
arr:     [ a | b | c | d | e | f ]
              ↑
s1 = arr[1:4] 指標 → b，len = 3（b c d），cap = 5（b c d e f）
```

<!--
這一頁是理解切片最關鍵的一頁。

切片本身其實是一個很小的結構，只有三個欄位：一個指向底層陣列的指標、長度、容量。

用窗戶比喻：指標是窗戶的左邊框在哪裡、長度是窗戶現在開多寬、容量是這面牆最多能把窗戶開到多寬。

因為切片只是一個「窗戶」，所以多個切片可以開在同一面牆上，也就是共用同一個底層陣列。這帶來一個很重要的後果，下一頁來看。
-->

---

# 切片共用底層陣列

```go
package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	a := arr[1:3] // [2 3]，和 arr 共用底層陣列
	a[0] = 20
	fmt.Println(arr) // [1 20 3 4 5]：透過切片修改，陣列也變了

	a = append(a, 99) // 容量還夠，直接寫進底層陣列的下一格
	fmt.Println(arr)  // [1 20 3 99 5]：arr[3] 被覆蓋了！

	b := arr[1:3:3]     // 完整切片運算式 [low:high:max]，限制 cap = 2
	b = append(b, 100)  // 容量不夠，換新的底層陣列
	fmt.Println(arr, b) // [1 20 3 99 5] [20 3 100]
}
```

<!--
這就是切片「共用底層陣列」的後果。

切片 a 是 arr 的一段，修改 a[0]，arr[1] 也跟著變了，因為它們是同一塊記憶體。

更隱蔽的是 append：a 的長度是 2，但容量是 4，所以 append 的時候不需要換陣列，直接寫進底層陣列的下一格，結果把 arr[3] 的 4 蓋成了 99。這是切片最容易踩的坑。

解法是用「完整切片運算式」，加上第三個數字 max，限制容量。b 的容量被限制為 2，append 時容量不夠，就會換一個新的底層陣列，不會影響原本的 arr。
-->

---

# 切片的隱藏陣列置換

當 `append` 超過容量，Go 會**配置一個更大的新陣列，把資料複製過去**：

```go
package main

import "fmt"

func main() {
	s := make([]int, 0, 2)
	for i := range 5 {
		s = append(s, i)
		fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
	}
}
```

```text
len=1 cap=2 [0]
len=2 cap=2 [0 1]
len=3 cap=4 [0 1 2]        ← 容量不夠，換成 cap=4 的新陣列
len=4 cap=4 [0 1 2 3]
len=5 cap=8 [0 1 2 3 4]    ← 再換一次，cap=8
```

<!--
這段程式碼的目的，是觀察 append 什麼時候會換陣列。

一開始容量是 2，加入前兩個元素都不用換。加入第三個元素時，容量不夠了，Go 會偷偷配置一個容量更大的新陣列，把舊資料複製過去，再加入新元素。這就是「隱藏陣列置換」。

容量通常會翻倍成長（小切片時），這樣平均下來 append 的成本很低。

這也解釋了為什麼 append 一定要寫成 s = append(s, ...)：換了陣列之後，只有回傳的新切片才指向新陣列。
-->

---

# 使用切片的注意事項

**注意事項之一：** 已知大概數量時，用 `make([]T, 0, n)` **預先配置容量**，避免反覆置換

**注意事項之二：** 需要一份獨立的副本時，用 `slices.Clone` 或 `copy`

```go
package main

import (
	"fmt"
	"slices"
)

func main() {
	src := []int{1, 2, 3}
	dst := slices.Clone(src) // Go 1.21+：複製出獨立的切片
	dst[0] = 100
	fmt.Println(src, dst) // [1 2 3] [100 2 3]

	buf := make([]int, 2)
	n := copy(buf, src) // 傳統寫法：複製 min(len(buf), len(src)) 個元素
	fmt.Println(n, buf) // 2 [1 2]
}
```

<!--
使用切片有兩個注意事項。

第一，如果我們大概知道會放多少個元素，例如從資料庫讀出 1000 筆，就用 make 預先把容量設好，這樣 append 的時候就不會一直換陣列、一直複製，效能會好很多。

第二，因為切片會共用底層陣列，如果我們需要一份「完全獨立」的副本，不能直接指定，要用 slices.Clone。這是 Go 1.21 加入的 slices 套件提供的函式，現在是主流寫法。傳統的寫法是先 make 再 copy，copy 會回傳實際複製了幾個元素。
-->

---

# 補充：slices 套件的常用函式（Go 1.21+）

| 函式 | 說明 | 範例（`s := []int{3, 1, 2}`） |
| --- | --- | --- |
| `slices.Sort(s)` | 由小到大排序（原地修改） | `s` 變成 `[1 2 3]` |
| `slices.Contains(s, v)` | 是否包含 | `slices.Contains(s, 2)` → `true` |
| `slices.Index(s, v)` | 第一次出現的索引，找不到是 `-1` | `slices.Index(s, 1)` → `1` |
| `slices.Max(s)` / `Min(s)` | 最大／最小值 | `slices.Max(s)` → `3` |
| `slices.Reverse(s)` | 反轉（原地修改） | `s` 變成 `[2 1 3]` |
| `slices.Equal(a, b)` | 兩個切片內容是否相同 | 切片不能用 `==` 比較 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>以前要自己寫迴圈：</b> Go 1.21 之前，找元素、排序都要自己寫迴圈或用 <code>sort</code> 套件。現在 <code>slices</code> 套件是處理切片的主流工具。
</div>

<!--
Go 1.21 加入了 slices 套件，把切片最常用的操作都包裝好了。以前找一個元素在不在切片裡，要自己寫 for 迴圈，現在一行 slices.Contains 就搞定。

特別注意最後一個：切片不能用 == 比較（只能跟 nil 比較），要比較兩個切片的內容，就用 slices.Equal。

這張表是最常用的幾個，完整的函式列表可以用 go doc slices 查詢，第 17 章會教。
-->

---
layout: default
---

# 練習 1：成績處理
### 任務說明

有一組成績 `scores := []int{72, 95, 58, 88, 64, 100, 45}`

1. 建立一個**新切片** `passed`，只放及格（≥ 60）的成績（提示：`append`）
2. 印出 `passed` 的長度、最高分、最低分（提示：`slices.Max`、`slices.Min`）
3. 複製一份 `sorted := slices.Clone(scores)` 並排序，印出前三名
4. 確認原本的 `scores` 順序沒有被改變

<!--
這個練習會用到切片的 append、切割、以及 slices 套件。

第 3 步要注意：排序是原地修改，所以要先複製一份再排序，才不會動到原本的 scores。前三名要怎麼取？排序之後是由小到大，最後三個才是最高分。
-->

---
zoom: 0.88
---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"slices"
)

func main() {
	scores := []int{72, 95, 58, 88, 64, 100, 45}
	passed := make([]int, 0, len(scores))
	for _, s := range scores {
		if s >= 60 {
			passed = append(passed, s)
		}
	}
	// 5 100 64
	fmt.Println(len(passed), slices.Max(passed), slices.Min(passed))

	sorted := slices.Clone(scores)
	slices.Sort(sorted)
	top3 := sorted[len(sorted)-3:]
	slices.Reverse(top3)
	fmt.Println(top3, scores) // [100 95 88] [72 95 58 88 64 100 45]
}
```

<!--
passed 用 make 預先配置容量，因為最多就是 len(scores) 個。

排序前先 Clone，排序完用切片運算式取出最後三個，再用 Reverse 反轉成由大到小。注意 top3 跟 sorted 共用底層陣列，Reverse 也會改到 sorted，但因為 sorted 是複製出來的，所以原本的 scores 完全不受影響。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 映射表
## Maps

<!--
接下來是 map，Go 內建的「鍵值對」資料結構。
-->

---

# map 的基礎

「**map 是『鍵（key）→ 值（value）』的對照表**，用鍵來查詢值，查詢速度很快。」

語法：`map[鍵的型別]值的型別`

```go
package main

import "fmt"

func main() {
	prices := map[string]int{ // map 常值
		"咖啡": 60,
		"紅茶": 30, // 最後一個元素後面也要加逗號
	}
	stock := make(map[string]int) // 用 make 建立空的 map

	prices["綠茶"] = 35 // 新增
	prices["咖啡"] = 65 // 修改：鍵已存在就覆蓋
	stock["咖啡"] = 10
	// map[咖啡:65 紅茶:30 綠茶:35] 3 map[咖啡:10]
	fmt.Println(prices, len(prices), stock)
}
```

<!--
什麼是 map？就像電話簿：用名字查電話，名字是鍵、電話是值。其他語言叫做 dictionary、hash table 或 HashMap，Go 叫做 map。

型別寫法是 map，中括號裡寫鍵的型別，後面接值的型別。map[string]int 就是「用字串查整數」。

建立方式有兩種：map 常值，或是用 make 建立一個空的 map。新增和修改的語法一樣，都是 m[key] = value，鍵不存在就新增，存在就覆蓋。

注意 map 常值裡，最後一個元素後面也要加逗號，這是 Go 的語法要求，好處是之後新增一行時不用回頭修改上一行。
-->

---

# 從 map 讀取元素：comma ok 慣用法

讀取不存在的鍵**不會報錯**，而是得到零值 — 用第二個回傳值判斷鍵是否存在：

```go
package main

import "fmt"

func main() {
	stock := map[string]int{"咖啡": 0, "紅茶": 12}

	fmt.Println(stock["可樂"]) // 0：鍵不存在，得到零值

	if qty, ok := stock["咖啡"]; ok { // ok 為 true 代表鍵存在
		fmt.Println("咖啡庫存：", qty) // 咖啡庫存： 0
	}
	if _, ok := stock["可樂"]; !ok {
		fmt.Println("沒有販售可樂")
	}
}
```

<!--
從 map 讀取值用 m[key]，但這裡有個問題：如果鍵不存在，Go 不會報錯，而是回傳值型別的零值。

那我們怎麼分辨「咖啡的庫存是 0」和「根本沒賣咖啡」？答案是 comma ok 慣用法：讀取時用兩個變數接，第二個 ok 是布林值，代表鍵是否存在。

這個寫法搭配 if 的起始賦值，是 Go 最經典的慣用寫法之一。之後學型別斷言、通道的時候，還會再看到同樣的 comma ok 模式。
-->

---

# 從 map 刪除元素與走訪

```go
package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	scores := map[string]int{"Carol": 88, "Alice": 92, "Bob": 75}
	delete(scores, "Bob") // 刪除；鍵不存在也不會報錯
	scores["Dave"] = 60

	for name, s := range scores { // ⚠️ 每次執行的順序都可能不同
		fmt.Println(name, s)
	}
	// Go 1.23+：依鍵排序後走訪
	for _, name := range slices.Sorted(maps.Keys(scores)) {
		fmt.Print(name, "=", scores[name], " ")
	}
	fmt.Println() // Alice=92 Carol=88 Dave=60
}
```

<!--
刪除元素用內建函式 delete，傳入 map 和鍵。鍵不存在也不會報錯，所以不用先檢查。

走訪 map 用 for range，每一圈拿到鍵和值。但要特別注意：map 的走訪順序是隨機的，而且是 Go 刻意設計成隨機的，避免我們寫出依賴順序的程式。

如果需要固定順序，例如依照名字排序輸出，現代的寫法是 slices.Sorted(maps.Keys(scores))：maps.Keys 取出所有的鍵，slices.Sorted 把它們排序成一個切片。這是 Go 1.23 加入的迭代器寫法，以前要自己寫迴圈收集鍵再排序。
-->

---

# 使用 map 的注意事項

**注意事項之一：** nil map 可以讀，但**寫入會 panic** — 一定要先 `make` 或用常值建立

**注意事項之二：** map 是**參考型別**，指定或傳進函式時共用同一份資料

```go
package main

import "fmt"

func addItem(m map[string]int, name string) {
	m[name]++ // 函式內的修改，外面看得到
}

func main() {
	var bad map[string]int
	fmt.Println(bad["x"]) // 0：讀取 nil map 沒問題
	// bad["x"] = 1       // panic: assignment to entry in nil map

	cart := map[string]int{}
	addItem(cart, "蘋果")
	addItem(cart, "蘋果")
	fmt.Println(cart) // map[蘋果:2]
}
```

<!--
使用 map 有兩個注意事項。

第一：nil map 只能讀不能寫。只宣告 var m map[string]int，m 是 nil，這時候寫入會直接 panic。一定要用 make 或 map 常值建立之後才能寫入。

第二：map 跟陣列不一樣，它是參考型別。把 map 傳進函式，函式裡面的修改在外面看得到，因為它們指向同一份資料。這也是為什麼 addItem 不需要傳指標。

另外 m[name]++ 這個寫法很方便：鍵不存在時，讀到零值 0，加 1 之後寫回去，所以可以直接拿來計數。

補充：map 不是「並行安全」的，多個 goroutine 同時寫入會出錯，第 16 章會說明怎麼處理。
-->

---
layout: default
---

# 練習 2：單字計數器
### 任務說明

統計一段文字中每個單字出現的次數：

```go
text := "go is fun and go is fast and go is simple"
```

1. 用 `strings.Fields(text)` 把文字切成單字切片
2. 用 `map[string]int` 計數
3. **依照單字字母順序**印出每個單字和次數
4. 找出出現最多次的單字

<!--
單字計數是 map 最經典的應用。

strings.Fields 會依照空白把字串切成切片，跟 Split 很像，但它會自動處理連續的多個空白。計數的部分，還記得剛剛 m[name]++ 的技巧嗎？
-->

---
zoom: 0.88
---

# 練習 2：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	text := "go is fun and go is fast and go is simple"
	count := map[string]int{}
	for _, w := range strings.Fields(text) {
		count[w]++
	}
	best := ""
	for _, w := range slices.Sorted(maps.Keys(count)) {
		fmt.Println(w, count[w])
		if count[w] > count[best] {
			best = w
		}
	}
	fmt.Println("最多：", best) // 最多： go
}
```

<!--
計數只要一行 count[w]++，因為不存在的鍵會讀到 0。

依字母順序印出，用 slices.Sorted(maps.Keys(count))。找最多次的單字，用一個 best 變數記錄目前最多的，一開始是空字串，count[""] 會讀到 0，所以任何單字都比它多。

go 和 is 都出現 3 次，因為用了大於而不是大於等於，而且是依字母順序走訪，所以先遇到的 go 會被保留。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 自訂型別
## Custom Types

<!--
接下來學怎麼定義自己的型別。
-->

---

# 簡易自訂型別 (custom types)

語法：`type 新型別名稱 底層型別`

```go
package main

import "fmt"

type Celsius float64    // 攝氏溫度
type Fahrenheit float64 // 華氏溫度
type UserID int

func toF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func main() {
	var t Celsius = 36.5
	fmt.Println(toF(t)) // 97.7

	// var f Fahrenheit = t // 編譯錯誤：Celsius 和 Fahrenheit 是不同型別
	var id UserID = 1001
	fmt.Printf("%v %T\n", id, id) // 1001 main.UserID
}
```

<!--
什麼是自訂型別？就是基於現有的型別，定義一個新的型別名稱。

為什麼需要？想像一個溫度轉換程式，攝氏和華氏都是 float64，如果不小心把華氏溫度傳進只收攝氏的函式，程式不會報錯，只會算出錯的結果。定義成 Celsius 和 Fahrenheit 兩個型別之後，編譯器就會幫我們擋下這種錯誤。

1999 年 NASA 的火星氣候探測者號，就是因為一個團隊用公制、另一個團隊用英制，最後探測器墜毀。自訂型別就是用來避免這種「單位搞混」的問題。

另外，自訂型別可以加上方法，這是它最重要的用途，等一下講 struct 的時候會看到。
-->

---

# 補充：型別定義 vs. 型別別名

| 寫法 | 名稱 | 是不是新型別 | 範例 |
| --- | --- | --- | --- |
| `type A B` | 型別定義 | ✅ 是，和 B 不能互相指定 | `type Celsius float64` |
| `type A = B` | 型別別名 | ❌ 不是，只是另一個名字 | `byte` 是 `uint8`、`rune` 是 `int32`、`any` 是 `interface{}` |

```go
type MyInt int     // 新型別：MyInt 和 int 需要轉型
type Number = int  // 別名：Number 就是 int

var a MyInt = 1
var b Number = 2
var c int = b      // ✅ 別名可以直接指定
var d int = int(a) // 新型別要轉型
```

<!--
這是補充內容，多了一個等號，意義就完全不同。

type A B 是「型別定義」，會產生一個全新的型別，跟原本的型別不能直接互相指定。

type A = B 是「型別別名」，只是幫同一個型別取另一個名字，兩者完全等價。上一章的 byte 和 rune，以及等一下會看到的 any，都是別名。

實務上我們自己寫程式幾乎都用型別定義；型別別名主要用在大型專案重構，把型別從一個套件搬到另一個套件時保持相容。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 結構
## Structs

<!--
接下來是 struct，Go 用來組織資料最重要的工具。
-->

---

# 結構的定義

「**struct 是把多個不同型別的欄位，組合成一個新型別。**」

```go
package main

import "fmt"

type Product struct {
	ID    int
	Name  string
	Price float64
	Tags  []string
}

func main() {
	// 用欄位名稱指定（建議）
	p1 := Product{ID: 1, Name: "咖啡豆", Price: 450}
	p1.Tags = append(p1.Tags, "熱銷") // 用 . 存取欄位
	// 零值：每個欄位都是零值
	var p2 Product
	fmt.Printf("%+v\n%+v\n", p1, p2)
}
```

<!--
什麼是 struct？

想像一張商品標籤，上面有商品編號、名稱、價格、分類標籤，這些資料的型別都不一樣，但它們描述的是同一件商品。struct 就是把這些欄位組合在一起，變成一個新的型別。

如果學過 Java 或 Python，可以把 struct 想成「只有屬性的類別」。Go 沒有 class，struct 加上方法就能做到類別能做的事。

建立 struct 時，建議用「欄位名稱: 值」的寫法，沒寫到的欄位自動是零值。Printf 的 %+v 會連欄位名稱一起印出來，除錯時非常好用。
-->

---

# 匿名結構與指標

```go
package main

import "fmt"

type Point struct{ X, Y int }

func main() {
	// 匿名結構：只用一次、不需要命名的資料
	cfg := struct {
		Host string
		Port int
	}{"localhost", 8080}
	fmt.Println(cfg.Host, cfg.Port)

	p := &Point{1, 2} // 取得結構的指標
	p.X = 10          // 自動解參考，等同於 (*p).X = 10
	fmt.Println(*p)   // {10 2}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>自動解參考：</b> 透過結構指標存取欄位時，Go 允許直接寫 <code>p.X</code>，不必寫成 <code>(*p).X</code>。
</div>

<!--
struct 還有兩種常見用法。

第一是匿名結構：直接在使用的地方定義並建立，不給它取名字。適合只用一次的資料，例如測試資料、臨時的設定。第 9 章寫單元測試時會大量用到。

第二是結構指標：用 & 取得結構的地址。Go 很貼心，透過指標存取欄位時會自動解參考，直接寫 p.X 就好，不用寫成 (*p).X。

補充：Point 的欄位寫成 X, Y int，同型別的欄位可以寫在同一行。建立時 Point{1, 2} 是依照欄位順序給值，這種寫法只建議用在欄位很少、順序很明確的結構。
-->

---

# 結構的相互比較

所有欄位都可以比較時，結構就能用 `==` 比較（逐欄位比對）：

```go
package main

import "fmt"

type Point struct{ X, Y int }

type Order struct {
	ID    int
	Items []string // 切片不能比較
}

func main() {
	a, b := Point{1, 2}, Point{1, 2}
	fmt.Println(a == b) // true

	o1, o2 := Order{ID: 1}, Order{ID: 1}
	// o1 == o2 // 編譯錯誤：切片欄位使結構無法比較
	fmt.Println(o1.ID == o2.ID) // 改成比較需要的欄位
}
```

<!--
struct 可以用 == 比較，Go 會逐個欄位比對，全部相等才是 true。

但前提是：每個欄位都要是「可以比較」的型別。數字、字串、布林、指標、陣列都可以比較，但切片、map、函式不能比較。只要 struct 裡有任何一個欄位是切片，整個 struct 就不能用 ==。

這時候就要自己比較需要的欄位。第 19 章會介紹 reflect.DeepEqual，它可以深度比較任何值。
-->

---
zoom: 0.96
---

# 內嵌結構 (embedding)

把一個結構**不寫欄位名稱**放進另一個結構，它的欄位和方法會被「提升」到外層：

```go
package main

import "fmt"

type Address struct {
	City, Street string
}

type Customer struct {
	Name    string
	Address // 內嵌：沒有欄位名稱
}

func main() {
	c := Customer{
		Name:    "小明",
		Address: Address{City: "台北", Street: "信義路"},
	}
	fmt.Println(c.City)         // 台北：直接存取被提升的欄位
	fmt.Println(c.Address.City) // 台北：也可以寫完整路徑
}
```

<!--
什麼是內嵌？Go 沒有繼承，而是用「組合」的方式重複使用程式碼，內嵌就是組合的主要手段。

在 Customer 裡放一個 Address，但不寫欄位名稱，這就是內嵌。這樣做之後，Address 的欄位會被「提升」到 Customer，可以直接寫 c.City，不用寫 c.Address.City。

這很像「繼承」的效果，但本質上是組合：Customer「有一個」Address，而不是 Customer「是一個」Address。這就是第一章提到的「組合勝過繼承」。
-->

---
zoom: 0.97
---

# 替自訂型別加上方法 (method)

方法就是「綁定在某個型別上的函式」：`func (接收器) 方法名稱(參數) 回傳值`

```go
package main

import "fmt"

type Rect struct{ W, H float64 }

func (r Rect) Area() float64 { // 值接收器：拿到複製品，只讀取
	return r.W * r.H
}

func (r *Rect) Scale(k float64) { // 指標接收器：可以修改原本的值
	r.W *= k
	r.H *= k
}

func main() {
	r := Rect{W: 3, H: 4}
	fmt.Println(r.Area())    // 12
	r.Scale(2)               // Go 自動取址，等同於 (&r).Scale(2)
	fmt.Println(r, r.Area()) // {6 8} 48
}
```

<!--
什麼是方法？方法就是綁定在某個型別上的函式。func 關鍵字後面的括號叫做「接收器」，代表這個函式屬於哪個型別。

接收器有兩種。值接收器 (r Rect)：方法拿到的是複製品，適合只讀取、不修改的方法，像 Area。指標接收器 (r *Rect)：方法拿到的是地址，可以修改原本的值，像 Scale。

這跟第一章「傳值 vs. 傳指標」的道理完全一樣。

呼叫時 Go 很貼心：r 是一個值，但呼叫指標接收器的 Scale 時，Go 會自動幫我們取址。
-->

---

# 使用方法的注意事項

**注意事項之一：** 需要修改接收器、或結構很大時，用**指標接收器**

**注意事項之二：** 同一個型別的方法，接收器種類**保持一致**（有一個用指標，就全部用指標）

| 情境 | 建議 |
| --- | --- |
| 方法會修改欄位 | 指標接收器 `*T` |
| 結構很大（很多欄位） | 指標接收器，避免複製 |
| 結構包含 `sync.Mutex` 等不能複製的欄位 | 指標接收器 |
| 小而不可變的值（如 `Point`、`time.Time`） | 值接收器 `T` |

<!--
使用方法的注意事項：什麼時候用值接收器、什麼時候用指標接收器？

簡單的判斷原則：會修改欄位、結構很大、或者包含不能複製的東西（例如第 16 章的互斥鎖），就用指標接收器。小而不會變的值，例如座標點、時間，用值接收器。

另外 Go 官方建議：同一個型別的方法，接收器種類要一致。如果有一個方法需要指標接收器，其他方法也都用指標接收器。實務上，大部分的 struct 都用指標接收器。
-->

---
layout: default
---

# 練習 3：銀行帳戶
### 任務說明

1. 定義 `type Account struct`，欄位有 `Owner string`、`Balance int64`
2. 加上方法 `Deposit(amount int64)`：存款（金額要大於 0）
3. 加上方法 `Withdraw(amount int64) bool`：提款，餘額不足時回傳 `false`
4. 加上方法 `String() string`：回傳 `"小明 的餘額：500"` 這樣的字串
5. 建立帳戶，存 1000、提 300、再提 1000，印出每次結果

<!--
這個練習要用到 struct 和方法，特別要想清楚每個方法該用值接收器還是指標接收器。

String 方法我們在第一章看過：定義了它之後，fmt.Println 印出這個型別時就會自動呼叫。
-->

---
zoom: 0.95
---

# 練習 3：解題提示
### 提示說明

```go
package main

import "fmt"

type Account struct {
	Owner   string
	Balance int64
}

func (a *Account) Deposit(amount int64) {
	if amount > 0 {
		a.Balance += amount
	}
}

func (a *Account) Withdraw(amount int64) bool {
	if amount > a.Balance {
		return false
	}
	a.Balance -= amount
	return true
}
```

<!--
Deposit 和 Withdraw 會修改餘額，所以用指標接收器。Withdraw 先檢查餘額夠不夠，不夠就提早 return false，這就是上一章說的「提早 return」。
-->

---

# 練習 3：解題提示（續）
### 提示說明

```go
// 續上頁
func (a *Account) String() string {
	return fmt.Sprintf("%s 的餘額：%d", a.Owner, a.Balance)
}

func main() {
	acc := &Account{Owner: "小明"}
	acc.Deposit(1000)
	// true false 小明 的餘額：700
	fmt.Println(acc.Withdraw(300), acc.Withdraw(1000), acc)
}
```

<!--
String 雖然只讀取，但依照「接收器保持一致」的原則，也用指標接收器。

注意 acc 是用 &Account{...} 建立的指標，因為 String 是定義在 *Account 上，Println 收到指標時才會呼叫它。

fmt.Sprintf 跟 Printf 很像，但不是印出來，而是回傳格式化後的字串。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 介面與型別檢查
## Interfaces & Type Checks

<!--
最後一個主題：型別轉換、型別斷言和型別 switch。這一段是介面的預告，第 7 章會正式深入。
-->

---

# 型別轉換 (type conversion)

語法：`目標型別(值)` — Go **不會自動轉型**，一定要明確寫出來

| 轉換 | 寫法 | 注意事項 |
| --- | --- | --- |
| 整數 ↔ 浮點數 | `float64(i)`、`int(f)` | 浮點數轉整數會**截斷小數**（3.99 → 3） |
| 大整數 → 小整數 | `uint8(n)` | 超出範圍會**溢位繞回** |
| 自訂型別之間 | `Celsius(f)` | 底層型別相同才能轉 |
| 字串 ↔ `[]byte` / `[]rune` | `[]byte(s)`、`string(r)` | 會複製一份資料 |
| 數字 ↔ 字串 | `strconv.Itoa(i)`、`strconv.Atoi(s)` | **不能**用 `string(i)` |

<!--
型別轉換的語法是「型別名稱加上括號」，像呼叫函式一樣。這張表整理了最常見的五種轉換。

特別注意兩個會「默默出錯」的情況：浮點數轉整數會截斷小數，不是四捨五入；大整數轉小整數會溢位繞回。這兩種情況 Go 都不會報錯。

另外數字轉字串一定要用 strconv，上一章說過 string(65) 會得到 "A"。
-->

---

# 型別轉換 — 範例

```go
package main

import (
	"fmt"
	"strconv"
)

type Celsius float64

func main() {
	x, big := 3.99, 300
	m := int(x)          // 截斷小數：3
	u := uint8(big)      // 溢位繞回：300 - 256 = 44
	c := Celsius(x)      // 自訂型別：3.99
	s := strconv.Itoa(m) // 數字 → 字串："3"
	fmt.Println(m, u, c, s+"!")
	// n := int(3.99)    // 編譯錯誤：常數轉換不能損失精度
}
```

<!--
這段程式碼的目的，是實際看看表格裡的轉換結果。

3.99 轉成 int 變成 3，小數被截斷。300 轉成 uint8，因為 uint8 最大 255，會繞回變成 44。

最後一行註解是一個細節：如果對「常數」做轉換，Go 要求不能損失精度，所以 int(3.99) 會編譯錯誤；只有變數才會被截斷。

執行後會印出 3 44 3.99 3!。
-->

---

# 型別斷言與空介面 any

`any`（= `interface{}`）可以裝**任何型別**的值；要取回原本的型別，就用**型別斷言** `x.(T)`

```go
package main

import "fmt"

func main() {
	var box any = "Hello" // any 是 interface{} 的別名（Go 1.18+）

	s := box.(string) // 斷言 box 裡是 string
	fmt.Println(s)    // Hello

	n, ok := box.(int) // comma ok：斷言失敗不會 panic
	fmt.Println(n, ok) // 0 false

	// m := box.(int) // ⚠️ panic：interface {} is string, not int
}
```

<!--
什麼是空介面？介面是第 7 章的主題，這裡先認識一個最特別的介面：空介面，寫成 interface{}，Go 1.18 之後有一個更短的別名叫 any。

any 就像一個萬用的收納箱，什麼型別的值都可以放進去。但放進去之後，Go 只知道「裡面有個東西」，不知道它是什麼型別，也就不能直接拿來運算。

要把東西拿出來用，就要做「型別斷言」：box.(string) 的意思是「我斷定箱子裡是字串，把它拿出來」。如果猜錯了，程式會 panic。

安全的做法是用 comma ok：n, ok := box.(int)，猜錯的話 ok 是 false，n 是零值，程式不會當掉。又是 comma ok 慣用法！
-->

---
zoom: 0.88
---

# 型別 switch

需要依照 `any` 裡的實際型別做不同處理時，用**型別 switch**：

```go
package main

import "fmt"

func describe(v any) string {
	switch x := v.(type) { // 特殊語法：.(type) 只能用在 switch 中
	case int:
		return fmt.Sprintf("整數，兩倍是 %d", x*2)
	case string:
		return fmt.Sprintf("字串，長度 %d", len(x))
	case []int:
		return fmt.Sprintf("整數切片，有 %d 個元素", len(x))
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("其他型別：%T", x)
	}
}

func main() {
	for _, v := range []any{21, "Go", []int{1, 2}, nil, 3.14} {
		fmt.Println(describe(v))
	}
}
```

<!--
如果 any 裡可能是好幾種型別，一個一個做型別斷言太麻煩，這時候就用型別 switch。

語法是 switch x := v.(type)，括號裡寫的是關鍵字 type。每個 case 寫一個型別，進入某個 case 之後，x 就已經是那個型別了，可以直接運算：在 int 的 case 裡 x 是 int，可以乘以 2；在 string 的 case 裡 x 是 string，可以取長度。

執行後會依序印出：整數，兩倍是 42；字串，長度 2；整數切片，有 2 個元素；nil；其他型別：float64。

第 7 章會看到型別 switch 跟介面搭配的更多用法。
-->

---

# 使用 any 的注意事項

**注意事項之一：** `any` 會讓編譯器失去型別檢查的能力，**能用具體型別就不要用 `any`**

**注意事項之二：** 需要「同一套邏輯處理多種型別」時，現代 Go 優先考慮**泛型**（第 5 章補充）

| 適合使用 `any` 的場景 | 範例 |
| --- | --- |
| 處理內容未知的資料 | 解析結構不固定的 JSON（第 11 章） |
| 需要接受任何值的函式 | `fmt.Println(a ...any)` |
| 容器裡放不同型別的值 | `map[string]any` 設定檔 |

<!--
使用 any 的注意事項。

any 很方便，但代價是：放進 any 之後，編譯器就不知道它是什麼型別了，原本在編譯時期就能抓到的錯誤，變成要等到執行時期才會 panic。所以原則是：能用具體型別就不要用 any。

如果我們想要「一個函式能處理 int，也能處理 float64」，在 Go 1.18 以前只能用 any 加型別斷言；現在有了泛型，可以兼顧彈性和型別安全，第 5 章會補充介紹。

any 真正適合的場景是資料的型別真的無法預先知道，例如第 11 章解析結構不固定的 JSON。
-->

---
layout: default
---

# 綜合練習：購物車
### 任務說明

1. 定義 `type Item struct { Name string; Price int; Qty int }`
2. 定義 `type Cart struct { Owner string; Items []Item }`
3. 替 `*Cart` 加上方法：
   - `Add(item Item)`：如果同名商品已存在，就增加數量；否則加入切片
   - `Total() int`：計算總金額
   - `Summary() map[string]int`：回傳「商品名稱 → 小計」的 map
4. 加入「咖啡 60×2」、「蛋糕 120×1」、「咖啡 60×1」，印出總金額與依名稱排序的小計

<!--
這個綜合練習會用到今天學的切片、map、struct、方法。

Add 方法是最有挑戰的：要先檢查切片裡有沒有同名的商品。注意用 for range 走訪時，拿到的 item 是「複製品」，直接修改它不會影響切片裡的元素，要用索引修改。
-->

---
zoom: 0.8
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"maps"
	"slices"
)

type Item struct {
	Name       string
	Price, Qty int
}

type Cart struct {
	Owner string
	Items []Item
}

func (c *Cart) Add(item Item) {
	for i := range c.Items {
		if c.Items[i].Name == item.Name {
			c.Items[i].Qty += item.Qty // 用索引修改切片中的元素
			return
		}
	}
	c.Items = append(c.Items, item)
}
```

<!--
Add 方法裡，for i := range c.Items 只取索引，然後用 c.Items[i] 修改，這樣才會真正改到切片裡的元素。如果寫成 for _, it := range，it 只是複製品，改了也沒用。找到同名商品就 return，沒找到才 append。
-->

---
zoom: 0.8
---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func (c *Cart) Total() (sum int) {
	for _, it := range c.Items {
		sum += it.Price * it.Qty
	}
	return sum
}

func (c *Cart) Summary() map[string]int {
	m := make(map[string]int, len(c.Items))
	for _, it := range c.Items {
		m[it.Name] = it.Price * it.Qty
	}
	return m
}

func main() {
	cart := &Cart{Owner: "小明"}
	cart.Add(Item{"咖啡", 60, 2})
	cart.Add(Item{"蛋糕", 120, 1})
	cart.Add(Item{"咖啡", 60, 1})
	fmt.Println("總金額：", cart.Total()) // 300
	sum := cart.Summary()
	for _, name := range slices.Sorted(maps.Keys(sum)) {
		fmt.Println(name, sum[name])
	}
}
```

<!--
Total 用了一個小技巧：回傳值有名字 (sum int)，這叫「具名回傳值」，第 5 章會詳細說明。

Summary 用 make 建立 map，第二個參數是預估大小，可以減少 map 擴容的次數。

總金額是 300：咖啡 60×3 = 180，蛋糕 120。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 4 步：商品與購物車

<!--
回到 GoShop。前三步我們都用一個一個獨立的變數來存商品：name、price、stock。只有一項商品還好，如果有一百項商品，就需要三百個變數，這顯然不行。

今天學的結構、切片和 map，就是用來把相關的資料組織在一起的。
-->

---

# GoShop 第 4 步：商品與購物車
### 任務說明

1. 定義結構 `Product`（`SKU`、`Name`、`Price`、`Stock`）與 `Item`（`SKU`、`Qty`）
2. 定義 `Cart` 結構，裡面是 `Items []Item`
3. 方法 `(c *Cart) Add(sku string, qty int)`：同一個 SKU 要**合併數量**
4. 方法 `(c Cart) Total(catalog map[string]Product) int`：用商品目錄查價、算總金額
5. 建立以 SKU 為鍵的商品目錄 `map[string]Product`；加入購物車前先用 **comma ok** 確認商品存在

```text
找不到商品： SKU-999
===== 購物車 =====
SKU-001 衣索比亞咖啡豆 x 3 = 1350
SKU-003 手沖壺 x 1 = 1280
SKU-004 馬克杯 x 2 = 700
品項數： 3
總金額： 3330
```

<!--
這一步要替 GoShop 建立三個結構：商品 Product、購物車品項 Item、購物車 Cart。

購物車有兩個方法。Add 負責放商品進去，如果購物車裡已經有同一個商品，就把數量加上去，而不是多一行。Total 負責算總金額。

為什麼 Add 用指標接收器，Total 用值接收器？因為 Add 要修改購物車的內容，Total 只是讀取。這就是今天「使用方法的注意事項」講的重點。

最後，商品目錄用 map 存，SKU 當鍵，這樣用 SKU 查商品就不需要走訪整個切片。
-->

---

# GoShop 第 4 步：解題提示
### 結構與指標接收器

```go
// goshop/main.go
// Item 是購物車裡的一個品項
type Item struct {
	SKU string
	Qty int
}

// Cart 是購物車
type Cart struct {
	Items []Item
}

// Add 把商品放進購物車；同一個 SKU 會合併數量
func (c *Cart) Add(sku string, qty int) {
	for i := range c.Items {
		if c.Items[i].SKU == sku {
			c.Items[i].Qty += qty
			return
		}
	}
	c.Items = append(c.Items, Item{SKU: sku, Qty: qty})
}
```

<!--
Add 方法先走訪購物車，找找看有沒有同一個 SKU。

這裡有一個很重要的細節：迴圈寫成 for i := range c.Items，然後用 c.Items[i].Qty 修改。如果寫成 for _, it := range c.Items，再修改 it.Qty，改到的只是複本，購物車裡的數量不會變。這是切片最常見的陷阱之一。

找不到的話，就用 append 加一個新的品項。因為 append 可能會產生新的底層陣列，所以要把結果指定回 c.Items，而且接收器必須是指標，外面的購物車才會看到改變。
-->

---

# GoShop 第 4 步：解題提示（續）
### 用 map 查價、comma ok 慣用法

```go
// goshop/main.go
// Total 用商品目錄查價，計算購物車總金額
func (c Cart) Total(catalog map[string]Product) int {
	total := 0
	for _, it := range c.Items {
		total += catalog[it.SKU].Price * it.Qty
	}
	return total
}
```

```go
// goshop/main.go
	for _, sku := range []string{"SKU-004", "SKU-999"} {
		if _, ok := catalog[sku]; !ok { // comma ok 慣用法
			fmt.Println("找不到商品：", sku)
			continue
		}
		cart.Add(sku, 2)
	}
```

<!--
Total 用值接收器，因為它只需要讀取。用 SKU 從 map 拿出商品，單價乘以數量加起來就是總金額。

下面這段示範 comma ok 慣用法。map 查不到鍵的時候不會出錯，而是傳回零值，所以直接 catalog["SKU-999"] 會得到一個價格是 0 的空商品，這很危險。用第二個傳回值 ok 判斷鍵存不存在，才是正確的做法。

執行後 SKU-999 會被擋下來，購物車裡有三個品項，總金額 3330 元。
-->

---

# 章節總結

- **陣列**：長度固定且是型別的一部分；是**值**，指定時整個複製
- **切片**：指標＋長度＋容量；`s = append(s, v)` 一定要接住；**共用底層陣列**要小心
- **map**：用 `v, ok := m[k]` 判斷鍵是否存在；走訪順序隨機；nil map 不能寫入
- **現代工具**：`slices`、`maps` 套件（Go 1.21+）；`slices.Sorted(maps.Keys(m))`（Go 1.23+）
- **自訂型別**：`type Celsius float64` 讓編譯器幫忙擋住單位錯誤
- **struct**：組合欄位；內嵌實現「組合勝過繼承」；方法分值接收器與指標接收器
- **型別轉換與斷言**：`T(v)` 明確轉型；`any` 用 `x.(T)` 或型別 switch 取回原型別
- **GoShop**：用 `Product`、`Cart` 結構表示商品與購物車，以 `map` 做商品目錄，替購物車加上方法

下一章我們會深入「函式」：多重回傳值、閉包、函式型別與 defer。

<!--
我們來整理今天學到的東西。

三種集合型別裡，切片是最常用的，要記住它的三個欄位，以及 append 一定要接住、共用底層陣列要小心這兩個重點。map 要記住 comma ok 慣用法和走訪順序是隨機的。

struct 是 Go 組織資料的方式，搭配方法和內嵌，就能做到其他語言用類別做的事。最後的型別斷言和型別 switch，第 7 章學介面時會再深入。

GoShop 也終於有了像樣的資料結構：商品是一個 struct，購物車裡的品項是一個切片，商品目錄是以 SKU 為鍵的 map，購物車的 Add 和 Total 則是今天學的方法。

今天我們已經寫了很多函式和方法，下一章會正式、完整地介紹函式：多重回傳值、參數不定函式、匿名函式與閉包，以及 Go 特有的 defer。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的內容非常多，特別是切片的底層陣列，第一次聽可能會覺得抽象。

課後建議：把「切片共用底層陣列」那一頁的範例在 Playground 上跑一次，每一行都印出 len、cap 和陣列的內容，親眼看看底層陣列是怎麼被改變的。

有問題的同學現在可以提問！
-->

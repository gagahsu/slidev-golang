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
title: Go 語言的特殊套件：reflect 與 unsafe
routeAlias: ch19
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
  <h1 style="color: #1a5c5c; font-size: 3.2rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">reflect 與 unsafe</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「在執行時期看見型別，也看見記憶體」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到最後一章，第十九章！

一路走來，我們用了很多「很神奇」的功能：fmt.Println 可以印出任何型別的值、encoding/json 可以讀到 struct 的標籤、database/sql 的 Scan 可以把資料填進任何型別的指標。這些套件是怎麼在執行時期知道一個值是什麼型別、有哪些欄位的？答案就是今天的主角 reflect，反射。

另一個主角是 unsafe。Go 是一個型別安全、記憶體安全的語言，但標準函式庫的某些地方，為了效能需要直接操作記憶體，unsafe 就是那扇「後門」。

這兩個套件平常寫程式很少直接用到，但理解它們，能讓我們真正理解 Go 的型別系統和記憶體模型。
-->

---
layout: default
---

# Outline

- **反射 (reflection)** — `TypeOf` 與 `ValueOf`、修改指標指向的值、走訪結構欄位與標籤
- **練習：用 reflect 取代介面斷言**
- **DeepEqual** — 深度比較任意值
- **unsafe 套件** — `unsafe.Pointer`、`uintptr` 與記憶體位址、標準函式庫中的 unsafe
- **章節總結與課程回顧**

<!--
今天的內容分成兩大塊。

前半段是 reflect：怎麼在執行時期取得一個值的型別和內容、怎麼修改它、怎麼走訪 struct 的欄位和標籤。我們會用 reflect 寫出一個迷你版的「驗證器」，理解 JSON 套件的運作原理。

後半段是 unsafe：怎麼直接操作記憶體位址，以及標準函式庫在哪裡用了 unsafe。

最後，因為這是最後一章，我們會一起回顧整門課程的學習路線。
-->

---

# 回顧：加密安全

- 隨機值一律用 **`crypto/rand`**；密碼用 bcrypt / Argon2id / PBKDF2 儲存
- **AES-GCM** 對稱式加密、**RSA-OAEP** 非對稱式加密、**Ed25519** 數位簽章
- HTTPS = 憑證驗證身分 + 非對稱式交換金鑰 + 對稱式加密傳輸
- 第 11 章：`json.Marshal` 怎麼讀到 **`json:"name"` 標籤**？ ← 今天揭曉

<!--
回顧一下上一章。

我們學了雜湊、對稱與非對稱加密、數位簽章，以及 HTTPS 的運作原理。

今天的懸念要回到第 11 章：我們在 struct 欄位後面寫了 json:"name" 這樣的標籤，json.Marshal 就能正確地使用 name 當作鍵名稱。標籤只是一段字串，Go 的編譯器並不理解它，那 JSON 套件是怎麼讀到它的？今天就會揭曉。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 反射 (reflection)
## The reflect Package

<!--
先從反射開始。
-->

---
zoom: 0.95
---

# 什麼是反射？

「**反射（reflection）是程式在執行時期，檢查與操作『自己的型別和值』的能力。**」

| 概念 | 取得方式 | 回答的問題 |
| --- | --- | --- |
| `reflect.Type` | `reflect.TypeOf(v)` | 這是什麼型別？有哪些欄位、方法？ |
| `reflect.Value` | `reflect.ValueOf(v)` | 它的值是多少？可以修改嗎？ |
| `reflect.Kind` | `t.Kind()` / `v.Kind()` | 它**底層**是哪一類：`Int`、`String`、`Struct`、`Slice`、`Map`、`Ptr`… |

| 誰在用反射 | 用來做什麼 |
| --- | --- |
| `fmt` | 印出任何型別的值（`%v`、`%+v`） |
| `encoding/json`、`encoding/gob` | 讀取欄位與標籤，編碼／解碼任意結構 |
| `database/sql` 的 `Scan`、ORM 套件 | 把資料填進任意型別 |

<!--
什麼是反射？

平常寫程式的時候，型別是在「編譯時期」決定的：宣告 var x int，編譯器就知道 x 是整數。反射則是讓程式在「執行時期」，才去查看一個值是什麼型別、有哪些欄位，甚至修改它。

就像照鏡子：程式透過反射這面鏡子，看見了自己的樣子。

reflect 套件有兩個核心型別：Type 描述「型別」，Value 描述「值」。另外 Kind 代表型別的底層種類，例如第 4 章的自訂型別 Celsius，它的 Type 是 Celsius，但 Kind 是 Float64。

表格下半部是使用反射的標準函式庫，大家從第一章就在用了：fmt.Println 之所以能印出任何型別，就是因為它用反射查看傳進來的值。
-->

---
zoom: 0.91
---

# TypeOf() 和 ValueOf()

```go
package main

import (
	"fmt"
	"reflect"
)

type Celsius float64

func main() {
	var temp Celsius = 36.5
	t := reflect.TypeOf(temp)
	v := reflect.ValueOf(temp)

	fmt.Println(t, t.Name(), t.Kind()) // main.Celsius Celsius float64
	fmt.Println(v, v.Float())          // 36.5 36.5

	values := []any{42, "Go", []int{1}, map[string]int{}, &temp}
	for _, x := range values {
		xt := reflect.TypeOf(x)
		fmt.Printf("%-16v Kind=%v\n", xt, xt.Kind())
	}
	// Go 1.22+：不需要值就能取得型別
	fmt.Println(reflect.TypeFor[Celsius]())
}
```

<!--
reflect.TypeOf 回傳值的型別，reflect.ValueOf 回傳值本身的反射物件。

temp 的型別是 main.Celsius，名稱是 Celsius，但它的 Kind 是 float64，因為它的底層型別是 float64。v.Float() 可以用 float64 的形式取出它的值。

迴圈裡示範了幾種不同的值：整數、字串、切片、map、指標，它們的 Kind 分別是 int、string、slice、map、ptr。實務上寫反射的程式碼，通常就是用 Kind 做 switch，依照不同的種類做不同的處理。

最後一行的 reflect.TypeFor 是 Go 1.22 加入的泛型函式，不需要先有一個值，就能取得某個型別的 Type。
-->

---

# 取得指標值和修改之

要透過反射修改變數，必須傳入**指標**，再用 **`Elem()`** 取得指標指向的值

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	price := 100

	v := reflect.ValueOf(price)
	fmt.Println(v.CanSet()) // false：傳入的是複製品，改了也沒意義

	p := reflect.ValueOf(&price).Elem() // 指標 → Elem() → 指向的值
	fmt.Println(p.CanSet())             // true
	p.SetInt(80)
	fmt.Println(price) // 80：原本的變數被修改了

	// p.SetString("x") // panic：SetString on int Value
}
```

<!--
反射也可以修改值，但有一個規則：一定要傳入指標。

為什麼？第一章學過，Go 傳參數一律是複製。reflect.ValueOf(price) 拿到的是 price 的複製品，修改複製品沒有意義，所以 Go 直接禁止，CanSet 回傳 false。

要修改原本的變數，就傳入 &price，再呼叫 Elem() 取得「指標指向的值」，這時候 CanSet 就是 true 了，可以用 SetInt 修改。這就是第 1 章「值 vs. 指標」觀念在反射裡的延伸。

使用反射的注意事項：反射沒有編譯時期的型別檢查。對一個 int 呼叫 SetString，編譯器不會報錯，要等到執行時才會 panic。這就是反射的代價：失去了 Go 最大的優點，型別安全。

這也解釋了為什麼 json.Unmarshal 和 rows.Scan 都要傳指標：它們內部就是用反射修改我們的變數。
-->

---
zoom: 0.84
---

# 取得結構的欄位名稱、型別與其值

```go
package main

import (
	"fmt"
	"reflect"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price,omitempty"`
	cost  float64 // 未匯出
}

func main() {
	p := Product{ID: 7, Name: "咖啡豆", Price: 450, cost: 300}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	for i := range t.NumField() {
		f := t.Field(i) // reflect.StructField：名稱、型別、標籤
		if !f.IsExported() {
			fmt.Println(f.Name, "（未匯出，略過）")
			continue
		}
		fmt.Printf("%-5s %-7v json=%-15q 值=%v\n",
			f.Name, f.Type, f.Tag.Get("json"), v.Field(i).Interface())
	}
}
```

<!--
這就是第 11 章懸念的答案：JSON 套件怎麼讀到標籤？

NumField 取得 struct 有幾個欄位，Field(i) 取得第 i 個欄位的資訊，型別是 reflect.StructField，裡面有欄位名稱、型別，以及標籤。f.Tag.Get("json") 就能讀出 json 標籤的內容，例如 "name"。v.Field(i) 取得這個欄位的值，Interface() 把它轉回 any。

encoding/json 做的事情本質上就是這樣：走訪每個欄位、讀取標籤決定鍵名稱、取出值編碼成 JSON。

注意 cost 是小寫開頭的未匯出欄位，IsExported 回傳 false。反射可以「看到」未匯出的欄位，但不能讀取或修改它的值，呼叫 Interface() 會 panic。這就是為什麼第 11 章說「只有匯出的欄位會被 JSON 處理」。
-->

---
zoom: 0.88
---

# 補充：Go 1.26+ 用迭代器走訪欄位

`reflect.Value` 新增 **`Fields()`** 迭代器，搭配 `for range` 同時取得欄位資訊與值

```go
package main

import (
	"fmt"
	"reflect"
)

type Config struct {
	Host string `env:"APP_HOST"`
	Port int    `env:"APP_PORT"`
}

func main() {
	cfg := Config{Host: "localhost", Port: 8080}
	// Go 1.26+
	for field, value := range reflect.ValueOf(cfg).Fields() {
		env := field.Tag.Get("env")
		fmt.Printf("%s（%s）= %v\n", field.Name, env, value)
	}
}
```

```text
Host（APP_HOST）= localhost
Port（APP_PORT）= 8080
```

<!--
這是補充內容：Go 1.26 替 reflect.Value 加入了 Fields 迭代器。

以前走訪欄位要用 NumField 加上索引，分別從 Type 和 Value 取出欄位資訊和值。現在可以直接 for range，每一圈同時拿到 StructField 和 Value，程式碼簡潔很多。這就是第 5 章學的迭代器函式的實際應用。

這個範例讀取自訂的 env 標籤，很多設定管理套件就是用這個方式，把環境變數自動填進設定結構。
-->

---

# 使用反射的注意事項

| 缺點 | 說明 |
| --- | --- |
| **失去型別安全** | 錯誤要到執行時期才 panic，編譯器幫不上忙 |
| **效能較差** | 比直接存取慢數倍到數十倍 |
| **可讀性差** | 程式碼冗長，難以理解與維護 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 諺語：「Clear is better than clever. Reflection is never clear.」</b><br>
能用介面、型別 switch、<b>泛型</b>解決的問題，就不要用反射。反射適合寫「處理任意型別」的通用函式庫，例如序列化、驗證、ORM。
</div>

<!--
使用反射有三個缺點：失去型別安全、效能較差、可讀性差。

Go 社群有一句諺語：「清楚勝過聰明，而反射從來都不清楚」。所以能用第 7 章的介面和型別 switch、第 5 章的泛型解決的問題，就不要用反射。

反射真正適合的場景，是寫「處理任意型別」的通用函式庫：序列化（JSON、gob）、驗證器、ORM、依賴注入框架。一般的應用程式碼，很少需要直接使用反射。
-->

---
layout: default
---

# 練習：用 reflect 取代介面斷言
### 任務說明

第 7 章我們用**型別 switch** 處理 `any`；但型別 switch 只能列出**已知的型別**。請用反射改寫：

1. 寫函式 `describe(v any) string`，用 `reflect.ValueOf(v).Kind()` 做 switch，支援：
   - 所有整數種類（`Int`、`Int8`…`Int64`，**包含自訂型別**如 `type Age int`）→ `"整數 42"`
   - `String` → `"字串（長度 n）"`；`Slice` → `"切片（n 個元素）"`；`Map` → `"map（n 個鍵）"`
   - `Struct` → 列出所有**匯出欄位**的名稱與值
   - `Pointer` → 以 `Elem()` **遞迴**描述它指向的值
2. 測試：`42`、`Age(18)`、`"Gopher"`、`[]int{1,2}`、`map[string]bool{}`、`Product{...}`、`&Product{...}`

<!--
這個練習要體會反射跟型別 switch 的差別。

第 7 章的型別 switch 寫 case int，只能匹配到 int 這個型別；如果傳入自訂的 Age 型別，就不會匹配。反射用 Kind 判斷的是「底層種類」，Age 的 Kind 是 Int，所以不管自訂了多少種整數型別，都能統一處理。

指標的部分，用 Elem 取得指向的值，再遞迴呼叫 describe。
-->

---

# 練習：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Age int

type Product struct {
	Name  string
	Price int
	stock int
}
```

<!--
先準備測試用的型別：自訂整數型別 Age，以及有一個未匯出欄位 stock 的 Product 結構。
-->

---

# 練習：解題提示（續）
### 提示說明

```go
// 續上頁
func describe(v any) string {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return fmt.Sprintf("整數 %d", rv.Int())
	case reflect.String:
		return fmt.Sprintf("字串（長度 %d）", rv.Len())
	case reflect.Slice:
		return fmt.Sprintf("切片（%d 個元素）", rv.Len())
	case reflect.Map:
		return fmt.Sprintf("map（%d 個鍵）", rv.Len())
	case reflect.Pointer:
		return "指向 " + describe(rv.Elem().Interface())
	case reflect.Struct:
		return describeStruct(rv)
	}
	return "其他：" + rv.Type().String()
}
```

<!--
describe 用 rv.Kind() 做 switch。第一個 case 列出了五種整數的 Kind，Age 的 Kind 是 Int，所以會匹配到這裡。rv.Int() 以 int64 的形式取出整數值。

String、Slice、Map 都可以用 rv.Len() 取得長度。Pointer 的時候，用 Elem() 取得指向的值，Interface() 轉回 any，再遞迴呼叫 describe。
-->

---
zoom: 0.86
---

# 練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func describeStruct(rv reflect.Value) string {
	var parts []string
	for i := range rv.NumField() {
		f := rv.Type().Field(i)
		if !f.IsExported() {
			continue // 未匯出的欄位不能呼叫 Interface()
		}
		val := rv.Field(i).Interface()
		parts = append(parts, fmt.Sprintf("%s=%v", f.Name, val))
	}
	return rv.Type().Name() + "{" + strings.Join(parts, ", ") + "}"
}

func main() {
	p := Product{Name: "咖啡", Price: 60, stock: 3}
	for _, v := range []any{42, Age(18), "Gopher", []int{1, 2},
		map[string]bool{}, p, &p} {
		fmt.Println(describe(v))
	}
}
```

```text
整數 42 / 整數 18 / 字串（長度 6） / 切片（2 個元素） / map（0 個鍵）
Product{Name=咖啡, Price=60} / 指向 Product{Name=咖啡, Price=60}
```

<!--
describeStruct 走訪結構的每個欄位，跳過未匯出的欄位，把名稱和值組成字串。

main 準備了各種不同的值。Age(18) 被正確地辨識成整數，這是型別 switch 做不到的；stock 是未匯出的欄位，不會出現在輸出裡；&p 是指標，透過 Elem 遞迴描述它指向的結構。

這個練習寫出來的 describe，其實就是 fmt 套件處理 %v 的迷你版。
-->

---
zoom: 0.88
---

# DeepEqual：深度比較任意值

`==` 不能比較切片、map、含有切片的結構（Ch 4）；**`reflect.DeepEqual`** 會遞迴比較所有內容

```go
package main

import (
	"fmt"
	"reflect"
)

type Order struct {
	ID    int
	Items []string
	Meta  map[string]int
}

func main() {
	a := Order{1, []string{"咖啡"}, map[string]int{"qty": 2}}
	b := Order{1, []string{"咖啡"}, map[string]int{"qty": 2}}
	// a == b // 編譯錯誤：含有切片的結構無法比較
	fmt.Println(reflect.DeepEqual(a, b)) // true

	// false ⚠️ 空切片 ≠ nil 切片
	fmt.Println(reflect.DeepEqual([]int{}, []int(nil)))
	// false：型別不同
	fmt.Println(reflect.DeepEqual(1, int64(1)))
}
```

<!--
第 4 章說過，含有切片的 struct 不能用 == 比較，當時我們說第 19 章會介紹 reflect.DeepEqual，現在兌現了。

DeepEqual 用反射遞迴地比較兩個值的所有內容：struct 的每個欄位、切片的每個元素、map 的每組鍵值，全部相等才回傳 true。它在單元測試裡特別常用，用來比較「實際結果」和「期望結果」。

使用 DeepEqual 的注意事項：它非常嚴格。空切片和 nil 切片在大部分情況下行為一樣，但 DeepEqual 認為它們不相等；型別不同的值，即使數值一樣也不相等。

補充：如果只是比較切片或 map，現代的 Go 建議用 slices.Equal 和 maps.Equal，它們是泛型函式，有型別檢查，而且比 DeepEqual 快很多。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# unsafe 套件
## Bypassing the Type System

<!--
接下來是 unsafe 套件。它的名字就是一個警告：使用它的程式碼，Go 不再保證安全。
-->

---

# 什麼是 unsafe？

「**unsafe 讓我們繞過 Go 的型別系統，直接操作記憶體。**」

| 函式／型別 | 用途 |
| --- | --- |
| `unsafe.Sizeof(x)` | 值佔用幾個位元組 |
| `unsafe.Alignof(x)` / `unsafe.Offsetof(s.f)` | 記憶體對齊 / 欄位在結構中的偏移量 |
| `unsafe.Pointer` | **萬用指標**：可以和任何指標型別、`uintptr` 互相轉換 |
| `unsafe.Add(p, n)` | 指標往後移動 `n` 個位元組（Go 1.17+） |
| `unsafe.String` / `StringData` / `Slice` / `SliceData` | 字串、切片與指標之間**零複製**轉換（Go 1.20+） |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>使用 unsafe 的程式碼：</b> 不受 Go 1 相容性承諾保護、可能在不同平台或新版本上出錯、錯誤會造成記憶體損毀而不是乾淨的 panic。
</div>

<!--
什麼是 unsafe？

Go 是一個型別安全的語言：int 不能當成 string 使用、指標不能做運算（第 1 章）、陣列存取會檢查邊界。這些保護讓 Go 很少出現 C 語言那種記憶體損毀的 bug。

unsafe 套件就是繞過這些保護的「後門」。它可以查詢記憶體的大小和配置，最重要的是 unsafe.Pointer：一個可以和任何指標互相轉換的萬用指標，透過它就能把一塊記憶體「當成」任何型別來讀寫。

使用 unsafe 的程式碼，Go 不再提供任何保證：不受 Go 1 相容性承諾保護、換一個平台可能就出錯，而且出錯的時候不是一個乾淨的 panic，而是記憶體損毀、莫名其妙的錯誤結果。
-->

---
zoom: 0.86
---

# 記憶體配置：Sizeof、Alignof、Offsetof

```go
package main

import (
	"fmt"
	"unsafe"
)

type Bad struct {
	A bool  // 1 byte + 7 bytes 填充（讓 B 對齊 8）
	B int64 // 8 bytes
	C bool  // 1 byte + 7 bytes 填充
}

type Good struct {
	B int64 // 8 bytes
	A bool  // 1 byte
	C bool  // 1 byte + 6 bytes 填充
}

func main() {
	// 24 16
	fmt.Println(unsafe.Sizeof(Bad{}), unsafe.Sizeof(Good{}))
	// 8 8
	fmt.Println(unsafe.Offsetof(Bad{}.B), unsafe.Alignof(int64(0)))
	// 16 24
	fmt.Println(unsafe.Sizeof(""), unsafe.Sizeof([]int{}))
}
```

<!--
unsafe 最安全、也最有教育意義的用法，是查詢記憶體的配置。

CPU 讀取記憶體時，希望資料的位址是自己大小的倍數，這叫做「對齊」。int64 需要對齊 8，所以 Bad 結構裡，A 後面會被填充 7 個位元組，讓 B 從第 8 個位元組開始；C 後面又填充 7 個，總共 24 個位元組。

只要調整欄位的順序，把大的欄位放前面，Good 結構只需要 16 個位元組，省了三分之一。如果一個切片裡有一百萬個這種結構，就省下了 8MB 的記憶體。

最後一行：字串佔 16 個位元組，是一個指標加一個長度；切片佔 24 個位元組，是指標、長度、容量。這就是第 3 章和第 4 章說的字串和切片的內部結構，unsafe.Sizeof 證實了它。
-->

---
zoom: 0.87
---

# unsafe.Pointer 指標

`unsafe.Pointer` 可以把一個指標**轉換成任何其他型別的指標**，直接用另一種型別解讀同一塊記憶體

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	f := 1.0
	// *float64 → unsafe.Pointer → *uint64
	bits := *(*uint64)(unsafe.Pointer(&f))
	// 0x3ff0000000000000：1.0 的 IEEE 754 表示
	fmt.Printf("%#x\n", bits)

	// 標準函式庫提供了安全的做法，結果相同：math.Float64bits(f)
}
```

| 合法的轉換規則（節錄） |
| --- |
| `*T` → `unsafe.Pointer` → `*U`：兩種型別的記憶體配置必須相容 |
| `unsafe.Pointer` → `uintptr` → 運算 → `unsafe.Pointer`：必須在**同一個運算式**中完成 |

<!--
unsafe.Pointer 是 unsafe 套件的核心。

一般的指標有型別：*float64 只能指向 float64。unsafe.Pointer 是一個沒有型別的萬用指標，任何指標都能轉換成它，它也能轉換成任何指標。

這個範例把 float64 的指標，經過 unsafe.Pointer，轉換成 uint64 的指標，再解參考。同一塊 8 位元組的記憶體，用整數的方式解讀，就看到了 1.0 在 IEEE 754 浮點數格式下的二進位表示，這是第 3 章浮點數精度問題的根源。

注意這個轉換之所以可行，是因為 float64 和 uint64 剛好都是 8 個位元組。如果兩種型別的大小或配置不相容，結果就是未定義的。其實標準函式庫的 math.Float64bits 就是這樣實作的，實務上應該直接用它。
-->

---

# 以 uintptr 搭配 unsafe 存取記憶體位址

`uintptr` 是**可以做運算的整數**；搭配 `unsafe.Pointer` 就能模擬 C 語言的指標運算

```go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	arr := [3]int32{10, 20, 30}
	base := unsafe.Pointer(&arr[0])
	size := unsafe.Sizeof(arr[0]) // 4

	// 傳統寫法：轉成 uintptr 做運算，再轉回來（必須在同一個運算式中）
	second := (*int32)(unsafe.Pointer(uintptr(base) + size))
	// 現代寫法（Go 1.17+）：unsafe.Add，意圖更清楚
	third := (*int32)(unsafe.Add(base, 2*size))
	fmt.Println(*second, *third) // 20 30

	fmt.Printf("%#x\n", uintptr(base)) // 記憶體位址，例如 0xc000012080
}
```

<!--
第一章說過，Go 的指標不能做運算，不能 p++ 移到下一個位址。但透過 uintptr 和 unsafe，就能做到。

uintptr 是一個足以存放記憶體位址的整數型別。把 unsafe.Pointer 轉成 uintptr，就變成一個普通的整數，可以加減；加上元素的大小，就是下一個元素的位址；再轉回 unsafe.Pointer 和具體型別的指標，就能讀到陣列的第二個元素。

使用 uintptr 的注意事項：轉換成 uintptr、運算、轉回 unsafe.Pointer，這三步必須在「同一個運算式」裡完成。因為 uintptr 只是一個整數，垃圾回收器不知道它其實指向某個物件，如果分成好幾行寫，中間垃圾回收器可能搬動或回收了那個物件，位址就失效了。

Go 1.17 加入的 unsafe.Add 把這個模式包裝起來，意圖更清楚、也比較不容易寫錯，是現代推薦的寫法。

這個範例只是為了理解原理，一般程式直接寫 arr[1] 就好，而且還有邊界檢查的保護。
-->

---
zoom: 0.83
---

# 零複製轉換：unsafe.String 與 unsafe.Slice（Go 1.20+）

`string(b)` 與 `[]byte(s)` 會**複製**資料；在效能關鍵的地方，可以用 unsafe **共用同一塊記憶體**

```go
package main

import (
	"fmt"
	"unsafe"
)

func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b)) // 不複製
}

func main() {
	b := []byte("hello")
	s := bytesToString(b)
	fmt.Println(s) // hello

	b[0] = 'J'     // ⚠️ 修改 b，s 也跟著變了！
	fmt.Println(s) // Jello
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ 使用前提：轉換之後，<b>原本的 <code>[]byte</code> 絕對不能再被修改</b>。一般情況下請直接用 <code>string(b)</code>，Go 編譯器已經會在許多情況下自動省略複製。
</div>

<!--
Go 1.20 加入了四個函式，用來在字串、切片和指標之間做「零複製」的轉換：String、StringData、Slice、SliceData。

一般的 string(b) 會把位元組複製一份，因為第 3 章說過，字串是不可修改的，如果不複製，修改原本的 []byte 就會改到字串。在處理大量資料、效能非常關鍵的地方，這個複製可能成為瓶頸，unsafe.String 可以讓字串直接共用 []byte 的記憶體。

但代價就在範例的下半段：修改 b 之後，s 也跟著變了，一個「不可修改」的字串被改了！這會破壞很多程式碼的假設，例如用字串當 map 的鍵，鍵的內容突然變了，map 就會出錯。

所以使用前提是：轉換之後，原本的 []byte 絕對不能再修改。一般情況下，請直接用 string(b)。
-->

---
zoom: 0.91
---

# Go 語言標準套件中的 unsafe

| 套件 | 如何使用 unsafe |
| --- | --- |
| `strings.Builder` | `String()` 用 `unsafe.String` 直接回傳內部緩衝區，**不複製**（Ch 3） |
| `reflect` | 透過 `unsafe.Pointer` 讀寫任意型別的值 |
| `sync/atomic` | `atomic.Pointer[T]` 以 unsafe 實作無鎖的指標操作（Ch 16） |
| `math` | `Float64bits` / `Float64frombits` 在浮點數與整數之間轉換 |
| `runtime`、`syscall` | 與作業系統、記憶體配置器互動 |

```go
// strings.Builder.String() 的實作（簡化）
func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>結論：</b> unsafe 由 Go 團隊在標準函式庫中**小心地封裝**起來，提供安全的 API 給我們使用。應用程式碼幾乎不需要直接使用 unsafe。
</div>

<!--
unsafe 在標準函式庫裡其實用得很多，只是都被小心地包裝起來了。

最好的例子是第 3 章學的 strings.Builder：它的 String 方法就是用剛剛學的 unsafe.String，直接把內部的位元組緩衝區當成字串回傳，不需要複製。為什麼它可以這樣做？因為 Builder 保證「只會往後附加、不會修改已經寫入的內容」，所以不會違反字串不可修改的前提。這就是為什麼 Builder 比用加號串接字串快那麼多。

reflect 套件能讀寫任意型別的值，背後也是 unsafe.Pointer；atomic.Pointer、math.Float64bits 都一樣。

結論是：unsafe 是 Go 團隊用來打造安全工具的「原料」，由專家小心地封裝成安全的 API。我們寫應用程式的時候，幾乎不需要直接使用它；知道它的存在和原理，就能更深刻地理解 Go 為什麼這麼安全、又這麼快。
-->

---
layout: default
---

# 綜合練習：迷你驗證器
### 任務說明

用反射寫一個簡化版的結構驗證器（類似 `go-playground/validator`）：

```go
type SignUp struct {
	Name  string `validate:"required,max=10"`
	Email string `validate:"required"`
	Age   int    `validate:"min=18"`
}
```

1. 寫 `Validate(v any) error`：`v` 必須是 struct 或指向 struct 的指標，否則回傳錯誤
2. 走訪所有匯出欄位，解析 `validate` 標籤（用 `strings.Split` 以逗號分隔，`strings.Cut` 拆出 `max=10`）
3. 支援規則：`required`（不可為零值，`value.IsZero()`）、`max=N`（字串的**字元數** ≤ N）、`min=N`（整數 ≥ N）
4. 收集所有錯誤，用 `errors.Join` 回傳（Ch 6）
5. 測試：`SignUp{Name: "Gopher", Email: "g@go.dev", Age: 20}` 與 `SignUp{Name: "非常非常長的名字啊啊", Age: 16}`

<!--
這個綜合練習是整門課的收尾：用反射讀取 struct 標籤，實作一個迷你版的驗證器。

實務上，第 15 章的 API 要驗證使用者送來的 JSON，很多專案會使用 go-playground/validator 這個套件，只要在欄位上加標籤就能自動驗證，它的原理就是今天學的反射。

會用到的東西包括：反射、struct 標籤、字串處理、strconv、utf8、errors.Join，幾乎把整門課的知識都串起來了。
-->

---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type SignUp struct {
	Name  string `validate:"required,max=10"`
	Email string `validate:"required"`
	Age   int    `validate:"min=18"`
}
```

<!--
SignUp 結構的每個欄位都用 validate 標籤描述規則，多個規則用逗號分隔。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func Validate(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("Validate 需要 struct，收到 %s", rv.Kind())
	}
	var errs []error
	for i := range rv.NumField() {
		f := rv.Type().Field(i)
		if tag := f.Tag.Get("validate"); f.IsExported() && tag != "" {
			errs = append(errs, checkField(f.Name, rv.Field(i), tag)...)
		}
	}
	return errors.Join(errs...)
}
```

<!--
Validate 先處理指標：如果傳入的是指標，用 Elem 取得指向的結構。不是結構就回傳錯誤。

接著走訪每個欄位，只處理匯出且有 validate 標籤的欄位，把檢查交給 checkField，收集所有錯誤。最後用第 6 章學的 errors.Join 合併，沒有錯誤時會回傳 nil。
-->

---
zoom: 0.75
---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func checkField(name string, v reflect.Value, tag string) []error {
	var errs []error
	for _, rule := range strings.Split(tag, ",") {
		key, arg, _ := strings.Cut(rule, "=")
		n, _ := strconv.Atoi(arg)
		switch {
		case key == "required" && v.IsZero():
			errs = append(errs, fmt.Errorf("%s 為必填", name))
		case key == "max" && v.Kind() == reflect.String &&
			utf8.RuneCountInString(v.String()) > n:
			errs = append(errs, fmt.Errorf("%s 最多 %d 個字", name, n))
		case key == "min" && v.CanInt() && v.Int() < int64(n):
			errs = append(errs, fmt.Errorf("%s 不可小於 %d", name, n))
		}
	}
	return errs
}

func main() {
	ok := SignUp{Name: "Gopher", Email: "g@go.dev", Age: 20}
	bad := &SignUp{Name: "非常非常長的名字啊啊啊", Age: 16}
	fmt.Println(Validate(ok)) // <nil>
	fmt.Println(Validate(bad))
}
```

```text
Name 最多 10 個字
Email 為必填
Age 不可小於 18
```

<!--
checkField 用逗號切開標籤裡的規則，再用 strings.Cut 把 max=10 拆成規則名稱和參數，參數用 strconv.Atoi 轉成數字。

用無條件 switch 檢查三種規則：required 用 IsZero 判斷是不是零值，第一章學的零值在這裡派上用場；max 檢查字串的字元數，用第 3 章的 RuneCountInString，中文才會算對；min 檢查整數，CanInt 確認這個值是整數種類。

第一筆資料全部通過，回傳 nil；第二筆有三個錯誤，errors.Join 把它們合併，印出來每個錯誤一行。

恭喜大家，寫出了一個真正實用的反射程式！
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 19 步：欄位驗證器

<!--
來到 GoShop 的最後一步。

現在 GoShop 有三個地方會收到外部的資料：API 的下單請求、後台的新增商品表單、CSV 匯入。每個地方都要檢查「SKU 不能空白」「數量要大於 0」，而且檢查的程式碼都是一大串 if，寫法還各不相同。

encoding/json 用 struct tag 決定欄位名稱，我們也可以用同樣的方法，用 struct tag 描述驗證規則，再用反射讀出來。
-->

---

# GoShop 第 19 步：欄位驗證器
### 任務說明

1. 新增 `internal/validate`，`Struct(v any) error` 依照 **`validate` 標籤**檢查欄位：

| 規則 | 數字 | 字串 | 切片 |
| --- | --- | --- | --- |
| `required` | 不是 0 | 不是空字串 | 不是 nil |
| `min=N`／`max=N` | 值 ≥ N／≤ N | **字數** ≥ N／≤ N | 長度 ≥ N／≤ N |

2. 切片裡的結構要逐一檢查，錯誤訊息用 **JSON 欄位名稱**，例如 `items[0].qty`
3. 所有問題用 `errors.Join` 一次回報；每個問題是一個 `*FieldError`
4. 替 `Product`、`checkout.Item`、`checkout.Cart` 加上標籤，在 API、後台表單、CSV 匯入使用

```text
$ curl -X POST localhost:8080/api/orders -d '{"items":[{"sku":"","qty":0}]}'
{"error":"欄位 items[0].sku 不符合規則 required\n欄位 items[0].qty 不符合規則 min=1"}
```

<!--
這一步要寫一個小小的驗證套件，概念和很多網頁框架內建的驗證器一樣。

規則寫在 validate 標籤裡，用逗號分隔。min 和 max 會依照欄位的種類做不同的比較：數字比大小，字串比字數，切片比長度。字串要算的是字數而不是位元組數，第 3 章學過，中文一個字是 3 個位元組。

錯誤訊息要用 JSON 的欄位名稱，因為呼叫 API 的前端工程師看到的是 JSON，他不知道 Go 的欄位叫 Qty。
-->

---

# GoShop 第 19 步：解題提示
### 在標籤裡描述規則

```go
// goshop/internal/shop/shop.go
type Product struct {
	SKU   string      `json:"sku" validate:"required,max=32"`
	Name  string      `json:"name" validate:"required,max=100"`
	Price money.Money `json:"price" validate:"min=1"`
	Stock int         `json:"stock" validate:"min=0"`
}
```

```go
// goshop/internal/checkout/cart.go
type Item struct {
	SKU string `json:"sku" validate:"required"`
	Qty int    `json:"qty" validate:"min=1,max=99"`
}

// Cart 是購物車。
type Cart struct {
	Items  []Item `json:"items" validate:"min=1,max=20"`
```

<!--
先看怎麼使用。規則直接寫在結構的定義上，和 JSON 標籤放在一起，一眼就能看出這個欄位的限制。

SKU 必填、最多 32 個字；價格至少 1 元；每個品項最少買 1 件、最多 99 件；購物車最少 1 項、最多 20 項。

注意 Price 的型別是 money.Money，不是 int。但它的底層型別是 int，反射的 Kind 還是 reflect.Int，所以驗證器不需要特別處理自訂型別。
-->

---

# GoShop 第 19 步：解題提示（續）
### 走訪欄位、讀取標籤

```go
// goshop/internal/validate/validate.go
func check(rv reflect.Value, prefix string) []error {
	var errs []error
	for f, fv := range rv.Fields() { // Go 1.26：用迭代器走訪欄位
		if !f.IsExported() {
			continue
		}
		name := prefix + fieldName(f)
		for rule := range strings.SplitSeq(f.Tag.Get("validate"), ",") {
			if rule != "" && !ok(fv, rule) {
				errs = append(errs, &FieldError{Field: name, Rule: rule})
			}
		}
		if fv.Kind() == reflect.Slice { // 切片裡的結構也要檢查
			for i := range fv.Len() {
				// ...
					errs = append(errs, check(elem, fmt.Sprintf("%s[%d].", name, i))...)
				}
			}
		}
	}
	return errs
}
```

<!--
check 是驗證器的核心。它用本章補充的 Go 1.26 新寫法 rv.Fields()，一次拿到每個欄位的描述 f 和欄位的值 fv。

沒有匯出的欄位跳過，反射也不能讀取它們的值。

f.Tag.Get("validate") 取出標籤字串，用 strings.SplitSeq 切成一條一條規則，逐一檢查，不符合就記下一個 FieldError。

如果欄位是切片，就走訪每一個元素，元素是結構的話，遞迴呼叫 check，並且把 items[0]. 這樣的前綴傳下去，錯誤訊息就會是 items[0].qty。
-->

---

# GoShop 第 19 步：解題提示（續 2）
### 依 Kind 決定怎麼比較

```go
// goshop/internal/validate/validate.go
func ok(v reflect.Value, rule string) bool {
	if rule == "required" {
		return !v.IsZero()
	}
	key, arg, _ := strings.Cut(rule, "=")
	n, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		panic("validate: 規則寫錯了：" + rule) // 這是程式設計師的錯誤
	}
	var size int64
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		size = v.Int()
	case reflect.String:
		size = int64(utf8.RuneCountInString(v.String()))
	case reflect.Slice, reflect.Map:
		size = int64(v.Len())
	// ...
```

<!--
ok 判斷一個值是否符合一條規則。

required 最簡單，用 IsZero 判斷是不是零值。

min 和 max 先把規則切成名稱和數字，再依照值的 Kind 算出要比較的大小：整數就是值本身，字串是字數，切片和 map 是長度。

規則寫錯的時候，例如 min=abc，我們選擇 panic 而不是傳回 error。第 6 章的指導方針說過：這是寫程式的人的錯誤，不是使用者的錯誤，應該在開發階段就讓它爆出來。
-->

---

# GoShop 第 19 步：解題提示（續 3）
### 一個驗證器，三個地方使用

```go
// goshop/internal/web/api.go
	if err := validate.Struct(cart); err != nil {
		writeErr(w, err)
		return
	}
```

```go
// goshop/internal/store/csv.go
		p := shop.Product{
			SKU: rec[0], Name: rec[1], Price: money.Money(price), Stock: stock}
		if err := validate.Struct(p); err != nil {
			errs = append(errs, fmt.Errorf("第 %d 行：%w", i+1, err))
			continue
		}
```

```text
$ go run . -import bad.csv
已匯入 1 項商品
以下資料列被略過：
 第 2 行：欄位 name 不符合規則 required
欄位 price 不符合規則 min=1
```

<!--
驗證器寫好之後，三個地方都只要一行 validate.Struct 就能完成檢查：API 解碼完購物車之後、後台表單組好商品之後、CSV 每一行轉成商品之後。

API 的 writeErr 也加了一個 case：錯誤鏈裡有 *FieldError 的時候回 400。

以後要加新的規則，例如 SKU 只能是英文和數字，只要改驗證器和標籤，三個地方同時生效。這就是反射最適合的場景：寫一次通用的工具，給很多不同的型別使用。
-->

---

# GoShop 完成了！

```text
goshop/                            2,700 行 Go，10 個套件，全部有測試
├── main.go、doc.go、Makefile      命令列、伺服器、版本、跨平台編譯
└── internal/
    ├── money/      金額與千分位（Stringer、Example）              Ch 5、7、9、17
    ├── shop/       商品、訂單、錯誤、時區（JSON／Text 標籤）        Ch 4、6、10、11
    ├── checkout/   購物車、折扣、折價券、結帳與付款                Ch 5、8、10、16
    ├── payment/    付款方式介面                                  Ch 7
    ├── store/      Store 介面：記憶體＋gob、MySQL、CSV            Ch 11～13、16
    ├── rates/      匯率 API 客戶端                               Ch 14
    ├── webhook/    POST 通知、背景 worker pool                   Ch 14、16
    ├── web/        RESTful API、後台模板、登入                    Ch 15、18
    ├── auth/       bcrypt、HMAC 簽章憑證                          Ch 18
    └── validate/   反射驗證器                                    Ch 19
```

- **unsafe** 沒有出現在 GoShop 裡：一般的應用程式**不需要**它，這正是本章的結論

<!--
恭喜大家，GoShop 完成了！

從第 0 章的一行 Println 開始，GoShop 現在有兩千七百行 Go 程式碼、十個套件，每個套件都有測試。它可以用命令列操作、可以當網站服務、可以接 MySQL、可以處理很多人同時搶購，還有登入和 HTTPS。

右邊標註了每個套件用到的章節，大家可以看到，幾乎每一章的內容都在這個專案裡留下了痕跡。

大家可能發現，今天的 unsafe 沒有用在 GoShop 裡。這不是忘記了，而是刻意的：unsafe 是給標準函式庫和極少數效能關鍵的程式用的，一般的應用程式完全不需要它。知道什麼時候不該用一個工具，跟知道怎麼用它一樣重要。

建議大家把自己的 GoShop 放上 GitHub，它就是你學會 Go 最好的證明。
-->

---

# 章節總結

- **反射**：`reflect.TypeOf` / `ValueOf`；`Kind()` 是底層種類；`reflect.TypeFor[T]()`（1.22+）
- **修改值**：傳入**指標**，`Elem()` 後 `CanSet()` 才是 `true`；型別不符會在執行時期 panic
- **結構欄位**：`NumField`、`Field(i)`、`Tag.Get("json")`、`IsExported`；`Value.Fields()` 迭代器（1.26+）
- **DeepEqual**：深度比較任意值；切片 / map 優先用 `slices.Equal` / `maps.Equal`
- **反射的代價**：失去型別安全、較慢、難讀 → 優先考慮介面、型別 switch、**泛型**
- **unsafe**：`Sizeof` / `Alignof` / `Offsetof` 查記憶體配置；`unsafe.Pointer` 轉換指標；`uintptr` 運算要在同一運算式，或用 `unsafe.Add`
- **零複製**：`unsafe.String` / `Slice`（1.20+）；標準函式庫（`strings.Builder`、`reflect`、`atomic`）已安全封裝，應用程式**幾乎不需要**直接使用
- **GoShop**：用反射讀取 `validate` 標籤，一個驗證器同時檢查 API 請求、後台表單與 CSV 匯入

<!--
我們來整理今天學到的東西。

反射讓程式在執行時期查看和操作型別與值，JSON、fmt、database/sql 都靠它運作。修改值要傳指標，讀取 struct 標籤用 Tag.Get。但反射有失去型別安全、效能差、難讀的代價，能用介面和泛型就不要用反射。

unsafe 讓我們繞過型別系統直接操作記憶體，標準函式庫用它打造了很多高效能的工具，但應用程式碼幾乎不需要直接使用它。

GoShop 的最後一步用反射寫了一個欄位驗證器：讀取 struct tag、走訪欄位、依 Kind 判斷怎麼比較，一次檢查 API 請求、後台表單和 CSV 匯入。到這裡，GoShop 從一行 Println 長成了一個兩千多行、十個套件、有測試、有資料庫、有網站的完整系統。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 課程回顧
## Course Review

<!--
這是整門課的最後一章，讓我們一起回顧這段學習旅程。
-->

---

# Go 實戰開發：學習地圖

| 階段 | 章節 | 學會的能力 |
| --- | --- | --- |
| **語言基礎** | Ch 0～5 | 環境、變數、流程控制、型別、切片與 map、struct、函式與閉包 |
| **Go 的設計哲學** | Ch 6～8 | 錯誤是值、介面與隱性實作、套件與模組 |
| **工程實務** | Ch 9～12 | 格式化與日誌、單元測試、時間、JSON、檔案與命令列 |
| **後端開發** | Ch 13～15 | 資料庫、HTTP 客戶端、HTTP 伺服器與 RESTful API |
| **進階主題** | Ch 16～19 | 並行性運算、工具鏈、加密安全、reflect 與 unsafe |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
🚀 <b>下一步：</b> 用這門課的知識完成一個完整的專案 — 例如「MySQL + RESTful API + HTTPS + 單元測試 + 跨平台發布」的待辦事項服務；接著可以探索 gRPC、Docker / Kubernetes 部署、OpenTelemetry 監控。
</div>

<!--
我們一起走完了 20 章的內容。

一開始，我們安裝環境、學變數和流程控制，打好語言的基礎。接著認識了 Go 最有特色的設計：錯誤是值、隱性實作的介面、以套件組織程式碼。然後進入工程實務：日誌、測試、時間、JSON、檔案。再來是後端開發的核心：資料庫、HTTP 客戶端和伺服器。最後挑戰了進階主題：並行性運算、工具鏈、加密安全，以及今天的反射和 unsafe。

學完這門課，大家已經具備用 Go 開發一個完整後端服務的能力了。建議的下一步，是把這些知識整合起來，自己完成一個完整的專案：用 MySQL 存資料、提供 RESTful API、加上 HTTPS 和單元測試、跨平台編譯發布。這個過程中遇到的每一個問題，都會讓你對 Go 有更深的理解。

之後可以繼續探索 gRPC 微服務、Docker 和 Kubernetes 的部署、以及可觀測性的工具，這些都是 Go 最擅長的領域。
-->

---
layout: end
---

# 課程結束，謝謝大家！

Happy Gophering 🐹

<!--
這門課到這裡就全部結束了，謝謝大家一路的參與。

Go 的設計哲學是「簡單」：少即是多、明確勝過隱晦、清楚勝過聰明。希望大家在寫 Go 的過程中，也能體會到這種簡單帶來的力量。

有任何問題，現在都可以提問。祝大家 Happy Gophering！
-->

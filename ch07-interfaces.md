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
title: 介面
routeAlias: ch07
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">介面</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「不問你是誰，只問你會做什麼」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第七章！

上一章我們看到 error 其實是一個介面：只要有 Error 方法，就是 error。今天就要正式學習介面，它是 Go 最強大、也最有特色的功能。

很多語言都有介面，例如 Java 的 interface。但 Go 的介面有一個關鍵的不同：不需要寫 implements。只要一個型別「有」介面要求的方法，它就自動實作了這個介面。這個設計讓 Go 的程式碼非常有彈性，也是 Go 不需要類別繼承的原因。
-->

---
layout: default
---

# Outline

- **介面** — 認識介面、定義、實作、隱性實作的優點
- **鴨子定型和多型** — 會呱呱叫的就是鴨子；同一個呼叫、不同的行為
- **在函式中活用介面** — 介面當參數、當回傳值、空介面、型別斷言與型別 switch
- **章節總結**

<!--
今天的內容分成三大塊。

第一塊是介面的基本概念：怎麼定義、怎麼實作，以及 Go 的「隱性實作」為什麼是一個好設計。

第二塊是兩個重要的觀念：鴨子定型和多型，這是介面能帶來彈性的原因。

第三塊是實務應用：在函式中怎麼使用介面，以及第 4 章提過的型別斷言和型別 switch，搭配介面的進階用法。
-->

---

# 回顧：錯誤處理

- `error` 是一個**介面**：`type error interface { Error() string }`
- 任何有 `Error() string` 方法的型別，都可以當作 `error` 使用
- 自訂錯誤型別 + `errors.As` 可以取出錯誤的詳細資訊
- 第 4 章的 `any`（`interface{}`）也是介面：**沒有任何方法要求**的介面

<!--
回顧一下上一章。

我們自訂了 ValidationError，只是替它加上 Error 方法，它就能當作 error 回傳。我們從來沒有寫過「ValidationError 實作了 error」這樣的宣告，Go 自己就知道了。這就是今天要講的隱性實作。

另外第 4 章的 any，其實是一個「沒有任何方法」的介面，所以任何型別都符合它。今天學完，大家對 any 會有更完整的理解。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 介面
## Interfaces

<!--
先從介面的基本概念開始。
-->

---

# 認識介面

「**介面（interface）定義了一組『行為』（方法），任何擁有這些方法的型別，都可以被當成這個介面使用。**」

| 生活比喻 | 程式世界 |
| --- | --- |
| 牆上的電源插座規格 | 介面：定義「要有兩根扁平的插腳」 |
| 吹風機、手機充電器、電鍋 | 實作介面的各種型別 |
| 插座不在乎插上來的是什麼電器 | 使用介面的函式不在乎實際的型別 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>一句話：</b> struct 描述「資料長什麼樣子」，介面描述「能做什麼事」。
</div>

<!--
什麼是介面？

想像牆上的電源插座。插座規定了一個規格：插頭要有兩根扁平的插腳。只要符合這個規格，吹風機、手機充電器、電鍋都能插上去用。插座不需要知道插上來的是什麼電器，它只在乎「符不符合規格」。

介面就是這個規格。它定義一組方法，任何擁有這些方法的型別，都可以被當作這個介面使用。使用介面的程式碼，就像插座一樣，不需要知道背後實際的型別是什麼。

第 4 章學的 struct 描述的是「資料長什麼樣子」，今天學的介面描述的是「能做什麼事」。
-->

---

# 定義介面型別

語法：`type 介面名稱 interface { 方法簽章... }`

```go
type Shape interface {
	Area() float64      // 只寫方法的簽章，不寫實作
	Perimeter() float64
}

type Stringer interface { // 標準函式庫 fmt 套件中的介面
	String() string
}
```

| 慣例 | 說明 | 範例 |
| --- | --- | --- |
| **小介面** | 介面越小越好，常常只有 1～2 個方法 | `io.Reader`、`fmt.Stringer`、`error` |
| **-er 命名** | 單一方法的介面，名稱用「方法名 + er」 | `Read` → `Reader`、`Write` → `Writer` |

<!--
定義介面的語法是 type 名稱 interface，大括號裡列出方法的簽章，只有名稱、參數和回傳值，沒有實作。

Go 社群對介面有兩個重要的慣例。第一，介面越小越好。標準函式庫裡最常用的介面，像 io.Reader、fmt.Stringer、error，都只有一個方法。Go 的名言是：「介面越大，抽象越弱」。

第二，只有一個方法的介面，名稱通常是「方法名加上 er」：有 Read 方法的叫 Reader，有 Write 方法的叫 Writer，有 String 方法的叫 Stringer。
-->

---
zoom: 0.9
---

# 實作一個介面

**不需要寫 `implements`**：只要型別擁有介面的**所有方法**，就自動實作了該介面

```go
package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Rect struct{ W, H float64 }
type Circle struct{ R float64 }

func (r Rect) Area() float64   { return r.W * r.H }
func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

func main() {
	var s Shape = Rect{3, 4}       // Rect 有 Area()，所以是 Shape
	fmt.Println(s.Area())          // 12
	s = Circle{1}                  // Circle 也有 Area()，也是 Shape
	fmt.Printf("%.2f\n", s.Area()) // 3.14
}
```

<!--
實作介面在 Go 裡非常簡單：什麼都不用宣告，只要把方法寫出來就好。

Rect 有 Area 方法，Circle 也有 Area 方法，它們就都自動實作了 Shape 介面。我們可以把 Rect 或 Circle 放進 Shape 型別的變數裡，然後呼叫 s.Area()，Go 會自動呼叫實際型別的 Area 方法。

在 Java 裡，我們必須寫 class Rect implements Shape；在 Go 裡完全不需要，這叫做「隱性實作」（implicit implementation）。
-->

---

# 使用介面的注意事項：指標接收器

方法定義在**指標接收器** `*T` 上時，**只有 `*T` 實作了介面**，`T` 沒有：

```go
package main

import "fmt"

type Counter interface{ Inc() }

type Clicks struct{ n int }

func (c *Clicks) Inc() { c.n++ } // 指標接收器

func main() {
	c := &Clicks{}
	var ctr Counter = c // ✅ *Clicks 實作了 Counter
	ctr.Inc()
	fmt.Println(c.n) // 1

	// Clicks does not implement Counter
	// (method Inc has pointer receiver)
	// var bad Counter = Clicks{} // 編譯錯誤
}
```

<!--
使用介面的注意事項：指標接收器。

第 4 章我們學過，會修改資料的方法要用指標接收器。但這會影響介面的實作：如果 Inc 是定義在 *Clicks 上，那只有 *Clicks（指標）實作了 Counter 介面，Clicks（值）沒有。

為什麼？因為如果把一個值放進介面，介面裡存的是一份複製品，呼叫 Inc 修改的也是複製品，這通常不是我們要的，所以 Go 直接禁止。

錯誤訊息寫得很清楚：method Inc has pointer receiver。遇到這個錯誤，就把值改成指標，用 & 取址就好。
-->

---

# 隱性介面實作的優點

| 優點 | 說明 |
| --- | --- |
| **解耦** | 實作的一方不需要 import 介面所在的套件，兩邊互不依賴 |
| **事後抽象** | 先寫具體型別，需要時才定義介面，舊程式碼不用改 |
| **替別人的型別定義介面** | 標準函式庫的 `*os.File`、`*bytes.Buffer` 都能符合我們自己定義的介面 |
| **介面定義在使用端** | Go 的慣例：「由需要的一方定義介面」，而不是由實作的一方 |

```go
// 我們自己定義的小介面，*os.File、*strings.Builder 都自動符合
type StringWriter interface {
	WriteString(s string) (int, error)
}
```

<!--
隱性實作有什麼好處？

最重要的是「解耦」：實作介面的型別，完全不需要知道介面的存在。就像電器的製造商，不需要事先跟每一家插座廠商簽約，只要符合規格就能用。

這帶來一個很強大的能力：我們可以替「別人寫的型別」定義介面。例如標準函式庫的 *os.File 和 *strings.Builder 都有 WriteString 方法，我們自己定義一個 StringWriter 介面，它們就自動符合了，完全不用修改標準函式庫的程式碼。

這也衍生出 Go 的一個重要慣例：介面應該定義在「使用它的地方」，而不是實作它的地方。需要什麼行為，就在那裡定義一個剛好夠用的小介面。
-->

---

# 補充：在編譯時期確認有實作介面

隱性實作的缺點是「沒寫出來」，可以用這個慣用法讓編譯器幫忙檢查：

```go
package main

import "fmt"

type Shape interface{ Area() float64 }

type Square struct{ Side float64 }

func (s Square) Area() float64 { return s.Side * s.Side }

var _ Shape = Square{} // 如果 Square 沒有實作 Shape，這行會編譯錯誤

func main() {
	fmt.Println(Square{3}.Area()) // 9
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>常見寫法：</b> 指標接收器時寫成 <code>var _ Shape = (*Square)(nil)</code>，不會真的配置任何記憶體。
</div>

<!--
這是補充內容。

隱性實作雖然彈性，但也有一個小缺點：沒有明確寫出來，如果不小心把方法名稱打錯，型別就默默地不再實作介面了，要等到使用的地方才會報錯。

Go 社群有一個慣用法：寫一行 var _ Shape = Square{}，把 Square 指定給一個 Shape 型別的變數，變數名稱是底線，代表我們不會使用它。如果 Square 沒有實作 Shape，這行就會編譯錯誤。這相當於在程式碼裡寫下「我保證 Square 實作了 Shape」，而且由編譯器來驗證。
-->

---
layout: default
---

# 練習 1：通知服務
### 任務說明

1. 定義介面 `type Notifier interface { Send(to, msg string) error }`
2. 實作兩個型別：
   - `EmailNotifier`：印出 `[Email] 寄給 to：msg`
   - `SMSNotifier`（欄位 `Quota int`）：印出 `[SMS] 傳給 to：msg`，每次傳送 `Quota` 減 1，額度用完回傳錯誤
3. `SMSNotifier` 會修改欄位，請用**指標接收器**
4. 用 `var _ Notifier = ...` 確認兩個型別都實作了介面
5. 建立 `[]Notifier`，放入兩種通知器，走訪並傳送訊息給 `"Alice"`

<!--
這個練習是介面最常見的實際應用：通知服務。實務上可能有 Email、簡訊、LINE、Slack 各種通知方式，用介面把它們統一起來，呼叫端只要呼叫 Send 就好。

注意 SMSNotifier 用指標接收器，放進切片的時候要放什麼？
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

type Notifier interface {
	Send(to, msg string) error
}

type EmailNotifier struct{}

func (EmailNotifier) Send(to, msg string) error {
	fmt.Printf("[Email] 寄給 %s：%s\n", to, msg)
	return nil
}

type SMSNotifier struct{ Quota int }

func (s *SMSNotifier) Send(to, msg string) error {
	if s.Quota <= 0 {
		return errors.New("簡訊額度已用完")
	}
	s.Quota--
	fmt.Printf("[SMS] 傳給 %s：%s\n", to, msg)
	return nil
}
```

<!--
EmailNotifier 是一個空的 struct，因為它不需要任何欄位。方法的接收器沒有寫名字，因為用不到，這也是合法的寫法。

SMSNotifier 用指標接收器，因為要修改 Quota。額度不足就回傳錯誤。
-->

---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
var (
	_ Notifier = EmailNotifier{}
	_ Notifier = (*SMSNotifier)(nil)
)

func main() {
	notifiers := []Notifier{EmailNotifier{}, &SMSNotifier{Quota: 1}}
	for range 2 {
		for _, n := range notifiers {
			if err := n.Send("Alice", "您的訂單已出貨"); err != nil {
				fmt.Println("傳送失敗：", err)
			}
		}
	}
}
```

<!--
兩行 var _ 確認兩個型別都實作了 Notifier。SMSNotifier 是指標接收器，所以寫成 (*SMSNotifier)(nil)。

放進切片的時候，SMSNotifier 要用 & 取址。外層迴圈跑兩次，第二次傳簡訊時額度用完了，會印出「傳送失敗：簡訊額度已用完」。

這就是介面的威力：main 裡的迴圈完全不知道 n 是 Email 還是簡訊，它只知道 n 會 Send。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 鴨子定型和多型
## Duck Typing & Polymorphism

<!--
接下來看兩個重要的觀念：鴨子定型和多型。
-->

---

# 鴨子定型 (duck typing)

「**如果它走起來像鴨子、叫起來像鴨子，那它就是鴨子。**」

```go
package main

import "fmt"

type Quacker interface{ Quack() string }

type Duck struct{}
type RobotDuck struct{ Model string }

func (Duck) Quack() string        { return "呱呱！" }
func (r RobotDuck) Quack() string { return r.Model + "：嗶—呱呱！" }

func makeItQuack(q Quacker) { fmt.Println(q.Quack()) }

func main() {
	makeItQuack(Duck{})
	makeItQuack(RobotDuck{"RD-2"}) // 機器鴨不是鴨子，但會呱呱叫就行
}
```

<!--
什麼是鴨子定型？

這個名字來自一句話：「如果它走起來像鴨子、叫起來像鴨子，那它就是鴨子。」意思是：我們不在乎一個東西「是什麼」，只在乎它「能做什麼」。

makeItQuack 只要求參數會 Quack。真的鴨子會呱呱叫，機器鴨雖然不是鴨子，但它也會呱呱叫，所以也能傳進去。

Python、JavaScript 這類動態語言也有鴨子定型，但它們是在執行時期才檢查，萬一傳進去的東西不會呱呱叫，程式執行到那一行才會出錯。Go 是在「編譯時期」就檢查，兼顧了彈性和安全，所以常被稱為「靜態的鴨子定型」。
-->

---
zoom: 0.85
---

# 多型 (polymorphism)

「**同一個呼叫，依照實際的型別，產生不同的行為。**」

```go
package main

import (
	"fmt"
	"math"
)

type Shape interface{ Area() float64 }

type Rect struct{ W, H float64 }
type Circle struct{ R float64 }
type Triangle struct{ B, H float64 }

func (r Rect) Area() float64     { return r.W * r.H }
func (c Circle) Area() float64   { return math.Pi * c.R * c.R }
func (t Triangle) Area() float64 { return t.B * t.H / 2 }

func main() {
	shapes := []Shape{Rect{3, 4}, Circle{2}, Triangle{6, 5}}
	total := 0.0
	for _, s := range shapes {
		total += s.Area() // 同一行程式碼，呼叫三種不同的 Area
	}
	fmt.Printf("總面積：%.2f\n", total) // 總面積：39.57
}
```

<!--
什麼是多型？「多型」的意思是「多種型態」：同一行程式碼 s.Area()，因為 s 實際的型別不同，會執行不同的方法。

這個範例把三種形狀放進同一個 []Shape 切片，用一個迴圈計算總面積。迴圈裡的程式碼完全不知道每個形狀是什麼，它只知道「每個形狀都能算面積」。

多型最大的好處是「擴充」：如果之後要新增梯形、五邊形，只要替新型別寫一個 Area 方法，這個計算總面積的迴圈一行都不用改。

在 Java 裡，多型通常靠繼承實現；在 Go 裡，完全靠介面實現。
-->

---
zoom: 0.93
---

# 多型的實例：fmt.Stringer

`fmt` 在印出值時，會檢查它是否實作了 `fmt.Stringer`，有的話就呼叫 `String()`：

```go
package main

import "fmt"

type Money int64 // 以「分」為單位

func (m Money) String() string {
	return fmt.Sprintf("NT$%d.%02d", m/100, m%100)
}

type Temp float64

func (t Temp) String() string {
	return fmt.Sprintf("%.1f°C", float64(t))
}

func main() {
	// NT$123.45 36.5°C
	fmt.Println(Money(12345), Temp(36.55))
	// %d 印出原始數值
	fmt.Printf("%v | %s | %d\n", Money(500), Money(500), Money(500))
}
```

<!--
標準函式庫裡處處都是多型，最常見的就是 fmt.Stringer。

fmt.Println 在印出一個值的時候，會先檢查這個值有沒有 String 方法。如果有，就呼叫它，印出它回傳的字串；如果沒有，才用預設格式印出。第一章的列舉、第四章的銀行帳戶，都是利用這個機制。

所以我們替 Money 加上 String 方法，印出來就自動變成貨幣格式。注意 %d 會印出原始的整數值 500，因為 %d 明確要求整數格式，不會呼叫 String。

小細節：Temp 的 String 方法裡，先把 t 轉成 float64 再格式化。這是一個好習慣：如果在 String 方法裡用 %v 或 %s 直接印 t 本身，fmt 又會呼叫 String，造成無窮遞迴；轉成底層型別就能避免這個問題。
-->

---
layout: default
---

# 練習 2：員工薪資
### 任務說明

1. 定義介面 `type Payable interface { Pay() int }`
2. 實作三種員工：
   - `FullTime{Name string; Salary int}`：月薪
   - `PartTime{Name string; Hours, Rate int}`：時數 × 時薪
   - `Contractor{Name string; Fee int}`：固定報酬
3. 替三種員工都加上 `String()` 方法，印出 `"姓名(類型)"`
4. 寫函式 `payroll(staff []Payable) int`：印出每個人的薪水，回傳總額

<!--
這個練習是多型的經典應用：薪資系統。

三種員工的計薪方式不同，但都能 Pay。payroll 函式只需要知道「每個人都能 Pay」，不需要知道他是正職、兼職還是約聘。

印出每個人的時候，fmt 會自動呼叫 String 方法。想想看：payroll 的參數是 []Payable，印出時 fmt 怎麼知道要呼叫 String？
-->

---
zoom: 0.83
---

# 練習 2：解題提示
### 提示說明

```go
package main

import "fmt"

type Payable interface{ Pay() int }

type FullTime struct {
	Name   string
	Salary int
}
type PartTime struct {
	Name        string
	Hours, Rate int
}
type Contractor struct {
	Name string
	Fee  int
}

func (f FullTime) Pay() int   { return f.Salary }
func (p PartTime) Pay() int   { return p.Hours * p.Rate }
func (c Contractor) Pay() int { return c.Fee }

func (f FullTime) String() string   { return f.Name + "(正職)" }
func (p PartTime) String() string   { return p.Name + "(兼職)" }
func (c Contractor) String() string { return c.Name + "(約聘)" }
```

<!--
三種員工各自實作 Pay 和 String。注意它們之間沒有任何繼承關係，只是剛好都有這兩個方法。
-->

---

# 練習 2：解題提示（續）
### 提示說明

```go
// 續上頁
func payroll(staff []Payable) int {
	total := 0
	for _, s := range staff {
		// fmt 在執行時發現 s 也是 Stringer
		fmt.Printf("%v：%d\n", s, s.Pay())
		total += s.Pay()
	}
	return total
}

func main() {
	staff := []Payable{
		FullTime{"Alice", 60000},
		PartTime{"Bob", 80, 200},
		Contractor{"Carol", 30000},
	}
	fmt.Println("總額：", payroll(staff)) // 總額： 106000
}
```

<!--
payroll 裡，s 的型別是 Payable，編譯器只知道它有 Pay 方法。但 fmt.Printf 在執行時期會檢查「s 裡面實際的值有沒有 String 方法」，發現有，就呼叫它。這就是等一下會學到的「用型別斷言檢查是否實作了另一個介面」。

總額是 60000 + 16000 + 30000 = 106000。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 在函式中活用介面
## Interfaces in Functions

<!--
接下來看介面在函式中的實際應用。
-->

---
zoom: 0.79
---

# 以介面為參數的函式

參數用介面型別，函式就能接受**所有符合這個行為**的值：

```go
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// 只要求「能寫入」
func writeReport(w io.Writer, title string, items []string) error {
	if _, err := fmt.Fprintf(w, "== %s ==\n", title); err != nil {
		return err
	}
	for i, it := range items {
		fmt.Fprintf(w, "%d. %s\n", i+1, it)
	}
	return nil
}

func main() {
	items := []string{"咖啡", "蛋糕"}
	writeReport(os.Stdout, "訂單", items) // 寫到螢幕

	var sb strings.Builder
	writeReport(&sb, "備份", items) // 寫到記憶體中的字串
	fmt.Print(sb.String())
}
```

<!--
這是介面在實務上最重要的用法：函式的參數用介面型別。

io.Writer 是標準函式庫最重要的介面之一，它只有一個方法 Write。螢幕輸出 os.Stdout、檔案、字串緩衝區 strings.Builder、網路連線、HTTP 回應，全部都實作了 io.Writer。

writeReport 的參數是 io.Writer，所以同一個函式可以寫到螢幕、寫到檔案、寫到記憶體，甚至寫到網路上。第 12 章寫檔案、第 15 章寫 HTTP 回應，都會用到同一個 io.Writer。

這就是 Go 的名言：「接受介面，回傳具體型別」（Accept interfaces, return structs）的前半句。
-->

---

# 標準函式庫中最重要的介面

| 介面 | 方法 | 實作者 |
| --- | --- | --- |
| `io.Reader` | `Read(p []byte) (n int, err error)` | 檔案、網路連線、HTTP 請求本體、`strings.Reader` |
| `io.Writer` | `Write(p []byte) (n int, err error)` | 檔案、`os.Stdout`、HTTP 回應、`bytes.Buffer` |
| `fmt.Stringer` | `String() string` | 任何想要自訂印出格式的型別 |
| `error` | `Error() string` | 所有錯誤型別 |
| `http.Handler` | `ServeHTTP(w, r)` | HTTP 伺服器的處理器（Ch 15） |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>介面可以組合：</b> <code>io.ReadWriter</code> 就是把 <code>Reader</code> 和 <code>Writer</code> 內嵌在一起：<code>type ReadWriter interface { Reader; Writer }</code>
</div>

<!--
這張表列出了標準函式庫裡最重要的五個介面，後面的章節會一再遇到它們。

特別是 io.Reader 和 io.Writer，它們是 Go 處理「資料流」的共通語言。只要實作了 Reader，就能被任何讀取資料的函式使用；只要實作了 Writer，就能被任何寫入資料的函式使用。這讓 Go 的標準函式庫像樂高積木一樣可以自由組合。

另外介面也可以內嵌其他介面來組合，就像第 4 章的 struct 內嵌一樣。io.ReadWriter 就是把 Reader 和 Writer 組合在一起。
-->

---
zoom: 0.81
---

# 以介面為傳回值的函式

一般建議**回傳具體型別**；但需要**隱藏實作**或**依條件回傳不同實作**時，回傳介面

```go
package main

import "fmt"

type Storage interface {
	Save(key, value string)
	Load(key string) (string, bool)
}

// 小寫：外部看不到實作細節
type memStorage struct{ data map[string]string }

func (m *memStorage) Save(k, v string) { m.data[k] = v }
func (m *memStorage) Load(k string) (string, bool) {
	v, ok := m.data[k]
	return v, ok
}

func NewStorage() Storage { // 回傳介面：呼叫端只能用 Save / Load
	return &memStorage{data: map[string]string{}}
}

func main() {
	s := NewStorage()
	s.Save("lang", "Go")
	fmt.Println(s.Load("lang")) // Go true
}
```

<!--
函式也可以回傳介面。

Go 的一般建議是「回傳具體型別」，讓呼叫端可以使用所有的方法和欄位。但有兩種情況適合回傳介面：第一，想要隱藏實作細節；第二，依照條件可能回傳不同的實作。

這個範例裡，memStorage 是小寫開頭，其他套件看不到它，只能透過 NewStorage 取得一個 Storage 介面，只能使用 Save 和 Load。將來如果要把記憶體儲存換成資料庫儲存，只要改 NewStorage 回傳另一個實作，呼叫端的程式碼完全不用改。

最典型的例子就是 error：函式都回傳 error 介面，而不是具體的錯誤型別。
-->

---

# 空介面 interface{}（any）

**沒有任何方法要求**的介面，所以**任何型別都符合**

```go
package main

import "fmt"

func main() {
	var anything any // any 是 interface{} 的別名
	anything = 42
	anything = "hello"
	anything = []int{1, 2, 3}
	fmt.Println(anything) // [1 2 3]

	settings := map[string]any{ // 值的型別不固定的 map
		"port":  8080,
		"debug": true,
		"hosts": []string{"a.com", "b.com"},
	}
	fmt.Println(settings["port"], settings["debug"])
}
```

<!--
現在我們可以用介面的角度重新理解 any 了。

介面的規則是「擁有介面要求的所有方法，就實作了這個介面」。空介面沒有要求任何方法，所以任何型別都自動實作了它。這就是為什麼任何值都能放進 any。

第 4 章說過，any 會讓我們失去型別檢查，所以能用具體型別或有方法的介面時，就不要用 any。它真正適合的場景，是像這個設定 map 一樣，值的型別真的不固定，或者第 11 章解析結構未知的 JSON。
-->

---
zoom: 0.77
---

# 型別斷言與型別 switch：搭配介面

型別斷言不只能取出具體型別，也能**檢查是否實作了另一個介面**：

```go
package main

import (
	"fmt"
	"math"
)

type Shape interface{ Area() float64 }
type Perimeterer interface{ Perimeter() float64 }

type Rect struct{ W, H float64 }
type Circle struct{ R float64 }

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }
func (c Circle) Area() float64    { return math.Pi * c.R * c.R }

func describe(s Shape) {
	fmt.Printf("面積 %.1f", s.Area())
	if p, ok := s.(Perimeterer); ok { // 選擇性功能：有實作才使用
		fmt.Printf("，周長 %.1f", p.Perimeter())
	}
	fmt.Println()
}

func main() {
	describe(Rect{3, 4}) // 面積 12.0，周長 14.0
	describe(Circle{1})  // 面積 3.1
}
```

<!--
第 4 章學過型別斷言，可以把 any 裡的值取出來變成具體型別。其實型別斷言還有一個更強大的用法：檢查一個值「是否也實作了另一個介面」。

describe 收到的是 Shape，一定能算面積。但有些形狀還能算周長，有些不能。我們用 s.(Perimeterer) 檢查：如果這個形狀也實作了 Perimeterer，就額外印出周長。

這個模式叫做「選擇性介面」，標準函式庫大量使用。剛剛說 fmt 會檢查值有沒有 String 方法，就是用這個技巧：fmt 內部用型別斷言檢查值是不是 fmt.Stringer。
-->

---
zoom: 0.79
---

# 型別 switch 處理多種介面與型別

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

func toText(v any) string {
	switch x := v.(type) {
	case nil:
		return "（空）"
	case error: // 介面也可以當 case
		return "錯誤：" + x.Error()
	case fmt.Stringer: // time.Duration 有 String() 方法
		return "Stringer " + x.String()
	case string:
		return x
	default:
		return fmt.Sprint(x)
	}
}

func main() {
	for _, v := range []any{
		3 * time.Second, errors.New("逾時"), 42, nil,
	} {
		fmt.Println(toText(v))
	}
}
```

<!--
型別 switch 的 case 除了寫具體型別，也可以寫介面。這個 toText 函式依序檢查：是不是 nil、是不是 error、是不是 fmt.Stringer、是不是字串，其他型別都交給 default。

有兩個細節要注意。第一，case 的順序很重要，Go 會從上往下找第一個符合的 case。如果一個型別同時是 error 和 Stringer，會進入 error 的 case。

第二，一個 case 也可以列出多個型別，例如 case int, int64，但這時候 x 的型別不會被確定，仍然是 any。

time.Duration 是標準函式庫的時間長度型別，它有 String 方法，所以會進入 Stringer 的 case。執行後會印出：Stringer 3s 錯誤：逾時 42 （空）。
-->

---
zoom: 0.87
---

# 使用介面的注意事項：nil 介面陷阱

介面值包含「**型別**」和「**值**」兩部分，**兩者都是 nil** 時介面才等於 `nil`

```go
package main

import "fmt"

type MyErr struct{}

func (*MyErr) Error() string { return "出錯了" }

func check(fail bool) error {
	var p *MyErr // p 是 nil 指標
	if fail {
		p = &MyErr{}
	}
	// ⚠️ 回傳的 error 介面：型別是 *MyErr，值是 nil → 介面「不是」nil
	return p
}

func main() {
	err := check(false)
	fmt.Println(err == nil) // false！
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
✅ <b>正確寫法：</b> 沒有錯誤時，<b>直接 <code>return nil</code></b>，不要回傳一個「值為 nil 的具體型別指標」。
</div>

<!--
使用介面的注意事項：這是 Go 最有名的陷阱之一，連資深工程師都會踩到。

介面的值其實由兩個部分組成：「型別」和「值」。只有兩個部分都是 nil 的時候，介面才等於 nil。

check(false) 回傳的是一個 nil 的 *MyErr 指標。這個指標被放進 error 介面之後，介面的「型別」部分是 *MyErr，「值」部分是 nil。因為型別部分不是 nil，所以整個介面不等於 nil。呼叫端檢查 err == nil 得到 false，以為出錯了，但其實沒有。

避免方式很簡單：沒有錯誤的時候，直接 return nil，不要回傳一個具體型別的 nil 指標。
-->

---
zoom: 0.93
---

# 補充：介面也是泛型的「約束」

Go 1.18 起，介面除了描述方法，還能描述**型別集合**，用來限制泛型的型別參數：

```go
package main

import "fmt"

type Number interface {
	~int | ~int64 | ~float64 // 型別集合：底層型別是這三種之一
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

type Score int // 底層型別是 int，符合 ~int

func main() {
	fmt.Println(Sum([]int{1, 2, 3}), Sum([]float64{1.5, 2.5})) // 6 4
	fmt.Println(Sum([]Score{90, 85}))                          // 175
}
```

<!--
這是補充內容，連結第 5 章的泛型。

第 5 章我們看過泛型的約束，像 any 和 cmp.Ordered，它們其實都是介面。Go 1.18 之後，介面除了列出方法，還可以列出「型別集合」，用直線符號隔開，意思是「這幾種型別之一」。

前面的波浪號 ~int 表示「底層型別是 int 的所有型別」，所以我們自訂的 Score 型別也符合。

注意：含有型別集合的介面只能當作泛型的約束，不能用來宣告一般的變數。執行後會印出 6 4 175。
-->

---
layout: default
---

# 綜合練習：支付系統
### 任務說明

1. 定義介面 `type PaymentMethod interface { Pay(amount int) error }`
2. 實作 `CreditCard{Limit int}`（超過額度回傳錯誤）與 `Cash{}`；`CreditCard` 用指標接收器，每次付款扣除額度
3. 定義**選擇性介面** `type Refunder interface { Refund(amount int) }`，只有 `CreditCard` 實作（退款加回額度）
4. 寫函式 `checkout(m PaymentMethod, amount int) error`：付款失敗時，用 `%w` 包裝錯誤回傳
5. 寫函式 `refund(m PaymentMethod, amount int)`：用型別斷言檢查是否支援退款，不支援就印出「此付款方式不支援退款」
6. 用型別 switch 寫 `name(m PaymentMethod) string`，回傳「信用卡」或「現金」

<!--
這個綜合練習把今天學的所有觀念整合起來：介面、指標接收器、多型、選擇性介面、型別斷言、型別 switch，再加上上一章的錯誤包裝。

最有挑戰的是第 5 步：不是每種付款方式都能退款，要用型別斷言檢查。
-->

---
zoom: 0.94
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
)

type PaymentMethod interface{ Pay(amount int) error }
type Refunder interface{ Refund(amount int) }

type CreditCard struct{ Limit int }
type Cash struct{}

func (c *CreditCard) Pay(amount int) error {
	if amount > c.Limit {
		return errors.New("超過信用額度")
	}
	c.Limit -= amount
	return nil
}
func (c *CreditCard) Refund(amount int) { c.Limit += amount }
func (Cash) Pay(amount int) error       { return nil }
```

<!--
CreditCard 有 Pay 和 Refund 兩個方法，所以同時實作了 PaymentMethod 和 Refunder；Cash 只有 Pay，只實作了 PaymentMethod。

CreditCard 會修改額度，所以兩個方法都用指標接收器。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func checkout(m PaymentMethod, amount int) error {
	if err := m.Pay(amount); err != nil {
		return fmt.Errorf("%s 付款 %d 元失敗：%w", name(m), amount, err)
	}
	return nil
}

func refund(m PaymentMethod, amount int) {
	if r, ok := m.(Refunder); ok {
		r.Refund(amount)
		return
	}
	fmt.Println(name(m), "不支援退款")
}
```

<!--
checkout 用 %w 包裝錯誤，加上付款方式和金額的情境。

refund 用型別斷言 m.(Refunder) 檢查付款方式是否支援退款，這就是選擇性介面的用法。
-->

---
zoom: 0.97
---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func name(m PaymentMethod) string {
	switch m.(type) {
	case *CreditCard:
		return "信用卡"
	case Cash:
		return "現金"
	}
	return "未知"
}

func main() {
	card := &CreditCard{Limit: 1000}
	// <nil> 信用卡 付款 500 元失敗：超過信用額度
	fmt.Println(checkout(card, 800), checkout(card, 500))
	refund(card, 800)
	refund(Cash{}, 100)              // 現金 不支援退款
	fmt.Println("剩餘額度：", card.Limit) // 剩餘額度： 1000
}
```

<!--
name 用型別 switch 判斷實際型別。因為 CreditCard 是用指標實作的，case 要寫 *CreditCard；Cash 是值接收器，case 寫 Cash。在 case 裡不需要使用值，所以 switch 可以不寫 x :=。

額度 1000，付 800 剩 200，再付 500 失敗，退款 800 回到 1000。
-->

---

# 章節總結

- **介面**：定義「能做什麼」的一組方法；越小越好，單一方法的介面用 `-er` 命名
- **隱性實作**：不需要 `implements`，擁有全部方法就自動實作；可用 `var _ I = T{}` 讓編譯器檢查
- **指標接收器**：方法定義在 `*T` 上時，只有 `*T` 實作介面
- **鴨子定型與多型**：只問行為、不問身分；同一個呼叫依實際型別產生不同行為
- **活用介面**：「**接受介面，回傳具體型別**」；`io.Reader` / `io.Writer` 是資料流的共通語言
- **型別斷言**：可以檢查「是否也實作了另一個介面」（選擇性介面）；型別 switch 的 case 可以是介面
- **nil 陷阱**：沒有錯誤時直接 `return nil`

下一章我們會介紹「套件」：怎麼把程式碼拆分成模組、管理第三方套件。

<!--
我們來整理今天學到的東西。

介面描述的是「能做什麼」，Go 的介面是隱性實作的，只要方法齊全就自動實作。介面越小越好，接受介面、回傳具體型別。型別斷言除了取出具體型別，還能檢查選擇性介面。最後記得 nil 介面陷阱：沒有錯誤就直接 return nil。

到這一章為止，我們已經學完了 Go 語言本身的核心語法。但到目前為止，我們所有的程式碼都寫在同一個 main.go 裡。真實的專案會有成千上萬行程式碼，需要拆分成很多個套件，也需要使用別人寫好的第三方套件。下一章就要學 Go 的套件與模組系統。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
介面是 Go 最重要的抽象工具，今天的觀念會在後面每一章反覆出現：第 12 章的檔案是 io.Reader，第 14、15 章的 HTTP 處理器是 http.Handler，第 13 章的資料庫驅動程式也是透過介面實作的。

課後建議：打開 Go 官方文件，查查看 io.Reader 有哪些實作者，體會一下「小介面、大組合」的威力。

有問題的同學現在可以提問！
-->

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
title: 錯誤處理
routeAlias: ch06
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">錯誤處理</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「錯誤只是一個值 — 看見它、處理它、說清楚它」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第六章！

程式一定會出錯：使用者輸入了奇怪的資料、檔案不存在、網路斷線、資料庫連不上。一個好的程式，不是「不會出錯」，而是「出錯的時候知道怎麼處理」。

Go 處理錯誤的方式跟大部分語言都不一樣：它沒有 try-catch，而是把錯誤當成一個普通的回傳值。一開始大家可能會覺得「怎麼到處都是 if err != nil」，但學完今天的內容，就會理解這個設計背後的道理。

另外今天也會介紹 panic 和 recover：什麼時候程式會「當掉」，以及怎麼把它救回來。
-->

---
layout: default
---

# Outline

- **程式錯誤的類型** — 語法錯誤、執行期間錯誤、邏輯錯誤
- **其他語言的錯誤處理** — 例外（exception）vs. 錯誤值
- **error 介面** — `error` 型別、建立錯誤、包裝錯誤、`errors.Is` / `errors.As`
- **panic** — 什麼是 panic、`panic()` 函式
- **recover** — 在 `defer` 中復原
- **指導方針** — 什麼時候回傳 error、什麼時候 panic
- **章節總結**

<!--
今天的內容從「錯誤有哪些種類」開始，接著比較 Go 跟其他語言的做法，然後深入 Go 的 error 介面。

error 介面的部分是今天的重點，我們會學到怎麼建立錯誤、怎麼包裝錯誤，以及怎麼判斷一個錯誤是哪一種。

最後會學 panic 和 recover，以及 Go 社群公認的錯誤處理指導方針。
-->

---

# 回顧：函式

- 多重回傳值的慣例：**「結果 + `error`」**，成功時 `error` 為 `nil`
- `if v, err := f(); err != nil { ... }`：起始賦值 + 檢查錯誤
- `defer` 在函式結束時執行（**包含 panic 的時候**），順序是後進先出
- `defer` 的閉包可以修改**具名回傳值**

<!--
回顧一下上一章。

我們學了多重回傳值，而且已經看過很多次「結果加 error」的寫法。另外 defer 有兩個特性今天會用到：第一，就算發生 panic，defer 也會執行；第二，defer 的閉包可以修改具名回傳值。這兩個特性，就是等一下 recover 能運作的基礎。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 程式錯誤的類型
## Types of Errors

<!--
首先，我們來看看程式的錯誤可以分成哪幾類。
-->

---

# 三種程式錯誤

| 類型 | 何時發現 | 範例 | 誰來抓 |
| --- | --- | --- | --- |
| **語法錯誤** | 編譯時期 | 少了大括號、型別不符、未使用的變數 | 編譯器 |
| **執行期間錯誤** | 執行時期 | 索引超出範圍、nil 指標、除以 0 | 程式 panic |
| **邏輯錯誤** | 可能永遠不會發現 | 公式寫錯、條件寫反、單位搞混 | 測試、除錯 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>越早發現越便宜：</b> 編譯錯誤只要改一行就好；執行期間錯誤可能讓服務中斷；邏輯錯誤最難找，可能等到客戶抱怨才發現。
</div>

<!--
程式的錯誤可以分成三類，依照「什麼時候被發現」來區分。

語法錯誤在編譯時期就會被發現，程式根本跑不起來。執行期間錯誤要等到程式跑起來、執行到那一行才會發生。邏輯錯誤最麻煩：程式可以正常執行、不會當掉，但算出來的結果是錯的。

就像蓋房子：語法錯誤是設計圖畫錯，審圖時就退件了；執行期間錯誤是施工時才發現材料不夠；邏輯錯誤是房子蓋好了，住進去才發現門開錯方向。

Go 的設計哲學是「盡量讓錯誤在編譯時期被發現」，這也是為什麼 Go 對型別這麼嚴格。
-->

---

# 語法錯誤

程式**無法編譯**，`go build` / `go run` 會列出錯誤的檔案與行號

```go
package main

import "fmt"

func main() {
	count := "10"
	// 編譯錯誤：mismatched types string and untyped int
	total := count + 5
	fmt.Println(total)
}
```

```text
./main.go:7:11: invalid operation: count + 5
  (mismatched types string and untyped int)
```

<!--
語法錯誤是最容易處理的錯誤，因為編譯器會直接告訴我們錯在哪裡：哪個檔案、第幾行、第幾個字元，以及錯誤的原因。

這個例子是把字串跟數字相加，Go 不會自動轉型，所以編譯失敗。錯誤訊息說得很清楚：「mismatched types string and untyped int」，型別不符。

在 VS Code 裡，這類錯誤通常在存檔之前就會用紅色波浪線標出來，所以大部分的語法錯誤，我們在寫程式的當下就修掉了。
-->

---

# 執行期間錯誤

程式可以編譯，但**執行到某一行時 panic**，印出錯誤原因與呼叫堆疊後結束

```go
package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}
	i := 5
	fmt.Println(nums[i]) // 執行時 panic
}
```

```text
panic: runtime error: index out of range [5] with length 3

goroutine 1 [running]:
main.main()
	/home/user/demo/main.go:8 +0x1d
exit status 2
```

<!--
執行期間錯誤，程式可以通過編譯，但執行到某一行的時候，發生了無法繼續的狀況，Go 就會 panic。

panic 會印出錯誤原因：索引 5 超出了長度 3 的範圍；接著印出「呼叫堆疊」（stack trace），告訴我們是哪個 goroutine、在哪個函式、哪一行發生的。最後程式以狀態碼 2 結束。

讀懂 stack trace 是除錯的基本功：從上往下看，第一個出現我們自己程式碼的那一行，通常就是問題所在。
-->

---

# 邏輯錯誤／語意錯誤

程式**正常執行、不會報錯**，但結果不是我們要的

```go
package main

import "fmt"

func average(nums []int) float64 {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum / len(nums)) // ⚠️ 先做整數除法，小數已經被捨去
}

func main() {
	fmt.Println(average([]int{1, 2})) // 期望 1.5，實際印出 1
}
```

修正：`return float64(sum) / float64(len(nums))`

<!--
邏輯錯誤是最難找的，因為程式完全不會報錯。

這個 average 函式看起來很合理，但 sum / len(nums) 是整數除法，先把小數捨去了，再轉成 float64 已經來不及。期望 1.5，結果是 1。

這種錯誤編譯器幫不上忙，要靠「測試」來發現：寫一個測試案例，輸入 1 和 2，期望得到 1.5，執行測試就會發現不對。第 9 章會教大家怎麼寫單元測試。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 其它程式語言的錯誤處理方式
## Exceptions vs. Error Values

<!--
在介紹 Go 的做法之前，我們先看看其他語言是怎麼處理錯誤的。
-->

---

# 例外（exception）vs. 錯誤值（error value）

| 比較 | Java / Python / C#：例外 | Go：錯誤值 |
| --- | --- | --- |
| **錯誤如何傳遞** | `throw` 拋出，沿著呼叫堆疊往上跳 | 函式 `return` 一個 `error` |
| **如何處理** | `try { } catch (e) { }` | `if err != nil { }` |
| **看得出哪裡會出錯嗎** | 不一定，任何一行都可能拋出例外 | 看函式簽章就知道：有回傳 `error` 的才會出錯 |
| **忘記處理會怎樣** | 例外一路往上拋，程式可能中止 | 編譯器會提醒未使用的變數；`go vet` 等工具會檢查 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 的哲學：</b> 「錯誤是值（Errors are values）」— 錯誤是正常流程的一部分，應該被明確地看見、處理，而不是被「拋」到看不見的地方。
</div>

<!--
其他主流語言大多使用「例外」機制：出錯時 throw 一個例外，它會沿著呼叫堆疊一路往上跳，直到遇到 catch 為止。

例外的問題是「看不見」：讀程式碼的時候，我們不知道哪一行可能會拋出例外、會被誰接住。錯誤處理的流程是隱藏的。

Go 選擇了另一條路：錯誤就是一個普通的回傳值。一個函式會不會出錯，看它的簽章就知道；錯誤要怎麼處理，就寫在呼叫的下一行。雖然程式碼會多一些 if err != nil，但錯誤處理的流程完全攤在陽光下，任何人讀程式都一目了然。

Go 的設計者 Rob Pike 有一句名言：「Errors are values」，錯誤是值，我們可以像處理其他值一樣處理它。
-->

---

# 同一件事，兩種寫法

**Java：例外** — 可能出錯的程式碼放進 `try`，出錯時跳到 `catch`

```java
try {
    int n = Integer.parseInt(input);
    System.out.println(n * 2);
} catch (NumberFormatException e) {
    System.out.println("格式錯誤：" + e.getMessage());
}
```

**Go：錯誤值** — 錯誤是回傳值，下一行馬上檢查

```go
n, err := strconv.Atoi(input)
if err != nil {
	fmt.Println("格式錯誤：", err)
	return
}
fmt.Println(n * 2)
```

<!--
我們用同一件事來比較：把使用者輸入的字串轉成整數。

Java 把可能出錯的程式碼放進 try 區塊，出錯時跳到 catch。Go 則是呼叫 Atoi，拿到兩個值：結果和錯誤，下一行馬上檢查錯誤，有錯就處理並 return，沒錯就繼續往下走。

注意 Go 的寫法裡，正常流程一直保持在最左邊，錯誤處理則是往右縮排然後提早 return。這就是第二章學的「提早 return」慣例，也被稱為「讓快樂路徑靠左對齊」（happy path on the left）。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# error 介面
## The error Interface

<!--
接下來進入今天的重點：Go 的 error 介面。
-->

---

# Go 語言的 error 值

「**`error` 是 Go 內建的介面型別，任何有 `Error() string` 方法的型別都是 error。**」

```go
type error interface {
	Error() string
}
```

| 特性 | 說明 |
| --- | --- |
| **零值是 `nil`** | `err == nil` 代表「沒有錯誤」 |
| **只有一個方法** | `Error()` 回傳錯誤訊息的字串 |
| **是介面** | 可以用任何型別實作，帶著更多資訊（檔名、欄位、狀態碼…） |

<!--
什麼是 error？它是 Go 內建的一個介面型別，定義非常簡單：只要有一個 Error 方法、回傳字串，就是 error。

介面是第 7 章的主題，這裡先簡單理解為「一種規格」：只要符合這個規格的型別，都可以當作 error 使用。

error 的零值是 nil，所以 err == nil 代表沒有錯誤。Error 方法回傳的字串就是錯誤訊息，fmt.Println(err) 印出的就是這個字串。
-->

---
zoom: 0.94
---

# 建立 error 值：errors.New

```go
package main

import (
	"errors"
	"fmt"
)

func withdraw(balance, amount int) (int, error) {
	if amount <= 0 {
		return balance, errors.New("提款金額必須大於 0")
	}
	if amount > balance {
		return balance, errors.New("餘額不足")
	}
	return balance - amount, nil
}

func main() {
	b, err := withdraw(100, 300)
	if err != nil {
		fmt.Println("錯誤：", err) // 錯誤： 餘額不足
	}
	fmt.Println("餘額：", b) // 餘額： 100
}
```

<!--
建立錯誤最簡單的方式是 errors.New，傳入一段錯誤訊息，就得到一個 error。

withdraw 函式檢查兩種不合法的情況，各自回傳不同的錯誤。成功時，錯誤的位置回傳 nil。

錯誤訊息的慣例：Go 官方建議錯誤訊息用小寫開頭、結尾不加標點符號，因為錯誤訊息常常會被串接成更長的訊息，例如「提款失敗：餘額不足」。中文沒有大小寫問題，但結尾一樣不要加句號。
-->

---
zoom: 0.88
---

# 哨兵錯誤（sentinel errors）

把常用的錯誤宣告成**套件層級變數**，呼叫端就能判斷「是哪一種錯誤」：

```go
package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("找不到資料") // 慣例：以 Err 開頭

func findUser(id int) (string, error) {
	users := map[int]string{1: "Alice"}
	name, ok := users[id]
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}

func main() {
	_, err := findUser(99)
	if errors.Is(err, ErrNotFound) { // 判斷是不是這個錯誤
		fmt.Println("使用者不存在，顯示註冊頁面")
	}
}
```

<!--
有時候呼叫端需要知道「是哪一種錯誤」，才能做不同的處理。例如「找不到使用者」要導向註冊頁面，「資料庫連線失敗」要顯示系統維護中。

這時候可以把錯誤宣告成套件層級的變數，叫做「哨兵錯誤」。慣例是用 Err 開頭命名。標準函式庫裡有很多哨兵錯誤，例如 io.EOF 代表檔案讀完了，sql.ErrNoRows 代表查詢沒有結果，第 12、13 章都會用到。

判斷的時候，用 errors.Is(err, ErrNotFound)。為什麼不直接用 ==？因為錯誤常常會被「包裝」，等一下就會看到。
-->

---
zoom: 0.93
---

# 自訂 error 型別

需要讓錯誤**帶著更多資訊**時，定義一個實作 `Error()` 的型別：

```go
package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
	Value any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("欄位 %s 的值 %v 不合法", e.Field, e.Value)
}

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{Field: "age", Value: age}
	}
	return nil
}
```

<!--
如果錯誤需要帶著更多資訊，例如「哪個欄位出錯、值是什麼」，就自訂一個 error 型別。

ValidationError 是一個 struct，有 Field 和 Value 兩個欄位，再加上 Error 方法，它就符合 error 介面了。validateAge 回傳的是它的指標。
-->

---

# 自訂 error 型別：用 errors.As 取出

呼叫端用 `errors.As` 判斷錯誤是不是某個型別，是的話**取出來讀取欄位**：

```go
// 續上頁
func main() {
	err := validateAge(-5)
	var ve *ValidationError
	if errors.As(err, &ve) { // 取出特定型別的錯誤
		fmt.Println("請修正欄位：", ve.Field) // 請修正欄位： age
	}
}
```

<!--
呼叫端想取出這些額外資訊，就用 errors.As：傳入錯誤和一個目標變數的指標，如果錯誤是這個型別，errors.As 會把它放進目標變數並回傳 true。這樣就能讀取 ve.Field，知道是哪個欄位出錯。

標準函式庫裡的 *fs.PathError（檔案錯誤）、*strconv.NumError（數字轉換錯誤）都是這樣設計的。
-->

---
zoom: 0.81
---

# 使用 fmt.Errorf() 建立 error 值

`fmt.Errorf` 可以用格式化字串建立錯誤；使用 **`%w`** 時會「包裝」原本的錯誤

```go
package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parsePort(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		// 包裝錯誤，加上情境
		return 0, fmt.Errorf("解析 port %q 失敗：%w", s, err)
	}
	if n < 1 || n > 65535 {
		return 0, fmt.Errorf("port %d 超出範圍", n) // 沒有要包裝時，不用 %w
	}
	return n, nil
}

func main() {
	_, err := parsePort("80a")
	// 解析 port "80a" 失敗：strconv.Atoi: parsing "80a": invalid syntax
	fmt.Println(err)
	// true：包裝後仍能判斷原始錯誤
	fmt.Println(errors.Is(err, strconv.ErrSyntax))
}
```

<!--
fmt.Errorf 跟 Sprintf 很像，可以用格式化字串建立錯誤訊息。

它最重要的功能是 %w：用 %w 放入另一個錯誤時，會把原本的錯誤「包裝」起來。就像寄包裹：原本的錯誤是裡面的商品，外面再套一層寫著「解析 port 失敗」的包裝紙。

包裝的好處有兩個。第一，錯誤訊息會帶著「情境」，一路往上傳時，每一層都加上自己在做什麼，最後的訊息就像一條完整的線索：解析 port 失敗，因為 Atoi 解析錯誤。第二，包裝之後，errors.Is 和 errors.As 仍然可以「拆開包裝」檢查裡面的錯誤。這就是為什麼判斷錯誤要用 errors.Is，而不是 ==。
-->

---

# 錯誤的包裝鏈與檢查工具

| 函式 | 用途 | 版本 |
| --- | --- | --- |
| `fmt.Errorf("...: %w", err)` | 包裝錯誤，加上情境 | 1.13+ |
| `errors.Is(err, target)` | 包裝鏈中**是否有**某個錯誤值（哨兵錯誤） | 1.13+ |
| `errors.As(err, &target)` | 包裝鏈中**是否有**某個型別，有就取出 | 1.13+ |
| `errors.Unwrap(err)` | 拆開一層包裝 | 1.13+ |
| `errors.Join(e1, e2, ...)` | 把多個錯誤合併成一個 | 1.20+ |
| `errors.AsType[T](err)` | 泛型版的 `As`，不用先宣告變數 | 1.26+ |

```go
if pe, ok := errors.AsType[*fs.PathError](err); ok { // Go 1.26+
	fmt.Println("出錯的檔案：", pe.Path)
}
```

<!--
這張表整理了 errors 套件的工具。

最常用的是前三個：Errorf 加 %w 包裝錯誤、Is 判斷哨兵錯誤、As 取出特定型別的錯誤。

errors.Join 是 Go 1.20 加入的，可以把多個錯誤合併成一個，適合「驗證表單時一次回報所有欄位的錯誤」這種情境。

errors.AsType 是 Go 1.26 加入的新寫法，利用泛型，不用先宣告一個變數再傳指標，一行就能完成判斷和取出，是目前最新的主流寫法。
-->

---

# 使用 error 的注意事項

**注意事項之一：** 回傳錯誤時，**其他回傳值用零值**；呼叫端**先檢查錯誤**再用結果

**注意事項之二：** 錯誤**只處理一次**：要嘛記錄（log），要嘛往上回傳，不要兩個都做

```go
// ❌ 記錄了又回傳：同一個錯誤會在 log 裡出現好幾次
if err != nil {
	log.Println("讀取設定失敗：", err)
	return err
}

// ✅ 加上情境後往上回傳，由最上層統一記錄
if err != nil {
	return fmt.Errorf("讀取設定失敗：%w", err)
}
```

<!--
使用 error 有兩個注意事項。

第一，回傳錯誤的時候，其他回傳值用零值；呼叫的人拿到結果，一定要先檢查 err，確認是 nil 才能使用結果。

第二是實務上很常見的問題：錯誤只處理一次。很多人習慣每一層都 log 一次再回傳，結果同一個錯誤在 log 裡出現五六次，反而更難追查。正確做法是：中間層用 %w 加上情境往上回傳，由最上層（例如 main 或 HTTP handler）統一決定要記錄還是顯示給使用者。
-->

---
layout: default
---

# 練習 1：會員註冊驗證
### 任務說明

1. 宣告哨兵錯誤 `ErrEmptyName = errors.New("名稱不可為空")`
2. 定義自訂錯誤 `type AgeError struct { Age int }`，`Error()` 回傳 `"年齡 X 不合法"`
3. 寫 `validate(name string, age int) error`：
   - 名稱是空字串 → 回傳 `ErrEmptyName`
   - 年齡小於 0 或大於 150 → 回傳 `&AgeError{age}`
4. 寫 `register(name string, age int) error`：呼叫 `validate`，失敗時用 `%w` 包裝成 `"註冊 name 失敗：..."`
5. 在 `main` 中用 `errors.Is` / `errors.As` 分別判斷兩種錯誤，印出不同的提示

<!--
這個練習把今天 error 介面的內容全部走一遍：哨兵錯誤、自訂錯誤型別、%w 包裝、errors.Is、errors.As。

重點是第 5 步：register 回傳的錯誤已經被包裝過了，但 errors.Is 和 errors.As 仍然能找到裡面的原始錯誤。
-->

---
zoom: 0.88
---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
)

var ErrEmptyName = errors.New("名稱不可為空")

type AgeError struct{ Age int }

func (e *AgeError) Error() string {
	return fmt.Sprintf("年齡 %d 不合法", e.Age)
}

func validate(name string, age int) error {
	if name == "" {
		return ErrEmptyName
	}
	if age < 0 || age > 150 {
		return &AgeError{age}
	}
	return nil
}
```

<!--
先看第一部分：哨兵錯誤用 errors.New 宣告成套件層級變數；AgeError 是一個 struct，加上 Error 方法就成為 error。

validate 依序檢查名稱和年齡，都通過才回傳 nil。
-->

---
zoom: 0.91
---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
func register(name string, age int) error {
	if err := validate(name, age); err != nil {
		return fmt.Errorf("註冊 %q 失敗：%w", name, err)
	}
	return nil
}

func main() {
	for _, c := range []struct {
		name string
		age  int
	}{{"", 20}, {"Bob", 200}, {"Amy", 30}} {
		err := register(c.name, c.age)
		if ae, ok := errors.AsType[*AgeError](err); ok {
			fmt.Println("年齡有誤，請重新輸入：", ae.Age)
		} else if errors.Is(err, ErrEmptyName) {
			fmt.Println("請輸入名稱")
		} else if err == nil {
			fmt.Println(c.name, "註冊成功")
		}
	}
}
```

<!--
register 用 %w 包裝 validate 回傳的錯誤，加上「註冊誰失敗」的情境。

main 用了一個匿名結構的切片當作測試資料，這個技巧上一章學過。每一筆資料都呼叫 register，再用 errors.AsType 和 errors.Is 判斷錯誤種類。即使錯誤已經被包裝過一層，這兩個函式仍然能找到裡面的原始錯誤。

如果環境是 Go 1.25 以前，把 AsType 換成 var ae *AgeError 加上 errors.As(err, &ae) 就可以了。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# panic
## When Programs Crash

<!--
接下來看 panic：當程式遇到無法繼續的狀況時，會發生什麼事。
-->

---

# 何謂 panic？

「**panic 代表程式遇到了『不應該發生、無法繼續』的狀況**，會停止目前函式的執行、執行所有 defer、然後讓程式結束。」

| 常見的執行期 panic | 錯誤訊息 |
| --- | --- |
| 索引超出範圍 | `index out of range [5] with length 3` |
| nil 指標解參考 | `invalid memory address or nil pointer dereference` |
| 寫入 nil map | `assignment to entry in nil map` |
| 整數除以 0 | `integer divide by zero` |
| 型別斷言失敗（沒用 comma ok） | `interface conversion: interface {} is string, not int` |

<!--
什麼是 panic？

panic 的意思是「恐慌」：程式遇到了完全沒預料到、無法處理的狀況，只好緊急停下來。就像開車時突然發現煞車失靈，這已經不是「繞路」可以解決的問題了，只能緊急停車。

表格裡列出了最常見的五種執行期 panic，前面幾章我們都遇過了。這些都是「程式寫錯了」造成的，也就是 bug。

panic 發生時，Go 會停止目前函式、往上一層一層執行每個函式的 defer，最後印出錯誤訊息和 stack trace，讓程式結束。
-->

---
zoom: 0.88
---

# panic() 函式

我們也可以自己呼叫 `panic()`，用在「**程式設計上的錯誤**」或「**無法繼續的初始化失敗**」

```go
package main

import "fmt"

type Level int

func (l Level) Name() string {
	switch l {
	case 1:
		return "銅"
	case 2:
		return "銀"
	case 3:
		return "金"
	}
	// 理論上不該發生，代表程式有 bug
	panic(fmt.Sprintf("不存在的等級：%d", l))
}

func main() {
	defer fmt.Println("defer 仍然會執行")
	fmt.Println(Level(2).Name()) // 銀
	fmt.Println(Level(9).Name()) // panic: 不存在的等級：9
}
```

<!--
除了 Go 自己觸發的 panic，我們也可以呼叫 panic 函式，傳入任何值作為原因。

什麼時候該自己 panic？當遇到「理論上不可能發生」的狀況，代表程式本身有 bug。例如等級只有 1 到 3，如果出現 9，一定是某個地方寫錯了，這時候繼續執行只會讓錯誤擴大，不如直接停下來，讓開發者盡快發現。

注意 main 裡的 defer：即使發生 panic，defer 仍然會執行，會先印出「defer 仍然會執行」，再印出 panic 訊息。
-->

---
zoom: 0.92
---

# 標準函式庫中的 Must 慣例

Go 社群的慣例：名稱以 **`Must`** 開頭的函式，失敗時會 **panic** 而不是回傳 error

```go
package main

import (
	"fmt"
	"regexp"
)

// 套件層級初始化：正規表達式寫錯是「程式碼的 bug」，直接 panic 最合理
var emailRE = regexp.MustCompile(`^[\w.+-]+@[\w-]+\.[\w.]+$`)

func main() {
	fmt.Println(emailRE.MatchString("gopher@go.dev")) // true
	fmt.Println(emailRE.MatchString("not-an-email"))  // false
}
```

| 函式 | 一般版本（回傳 error） | Must 版本（panic） |
| --- | --- | --- |
| 正規表達式 | `regexp.Compile` | `regexp.MustCompile` |
| HTML 模板 | `template.New(...).Parse` | `template.Must(...)`（Ch 15） |

<!--
標準函式庫有一個慣例：以 Must 開頭的函式，失敗時會 panic。

為什麼需要這種函式？像正規表達式，如果是寫死在程式碼裡的，它寫錯就是程式的 bug，不是使用者的錯，而且在程式啟動時就會發現。這種情況下，回傳 error 反而讓程式碼變囉嗦，用 MustCompile 讓它在啟動時直接 panic，是最合理的選擇。

但如果正規表達式是使用者輸入的，就一定要用 Compile，回傳 error 讓使用者知道哪裡寫錯。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# recover (復原)
## Recovering from Panics

<!--
panic 會讓程式結束，但有時候我們不希望整個程式都掛掉，例如一個 HTTP 請求出錯，不應該讓整個伺服器停止服務。這時候就需要 recover。
-->

---

# 什麼是 recover？

「**`recover()` 可以攔截 panic，讓程式恢復正常執行**。它只有在 `defer` 的函式中呼叫才有效。」

```go
package main

import "fmt"

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil { // 有 panic 發生時，r 是 panic 的值
			err = fmt.Errorf("計算失敗：%v", r) // 把 panic 轉成 error
		}
	}()
	return a / b, nil // b 為 0 時會 panic
}

func main() {
	fmt.Println(safeDivide(10, 2)) // 5 <nil>
	// 0 計算失敗：runtime error: integer divide by zero
	fmt.Println(safeDivide(1, 0))
	fmt.Println("程式繼續執行")
}
```

<!--
什麼是 recover？

如果 panic 是緊急煞車，recover 就是安全氣囊：它不能阻止事故發生，但可以讓我們毫髮無傷地繼續上路。

recover 只能在 defer 的函式裡使用。當 panic 發生時，Go 會執行 defer，這時候呼叫 recover，就能拿到 panic 的值，而且 panic 會被「攔截」，不會再往上傳。

這個範例把上一章學的兩個技巧組合起來：defer 閉包加上具名回傳值。recover 攔截到 panic 之後，把它轉成一個 error，指定給具名回傳值 err。呼叫端就像處理一般錯誤一樣處理它，程式可以繼續執行。
-->

---

# 使用 recover 的注意事項

**注意事項之一：** `recover()` **只在 defer 的函式中直接呼叫**才有效，在其他地方呼叫永遠回傳 `nil`

**注意事項之二：** recover **只能攔截同一個 goroutine** 的 panic

**注意事項之三：** 不要用 panic / recover 取代正常的錯誤處理

| 適合用 recover 的地方 | 原因 |
| --- | --- |
| HTTP 伺服器的每個請求 | 一個請求出錯，不應該讓整個伺服器停止（`net/http` 內建就有做） |
| 背景工作的 goroutine | 避免一個工作的 bug 讓整個程式崩潰（第 16 章） |
| 套件的公開 API 邊界 | 把內部的 panic 轉成 error，不要讓 panic 跨出套件 |

<!--
使用 recover 有三個注意事項。

第一，recover 只有在 defer 的函式裡直接呼叫才有效。在一般的程式碼裡呼叫，永遠只會拿到 nil。

第二，recover 只能攔截同一個 goroutine 的 panic。第 16 章會學到 goroutine，如果在一個新的 goroutine 裡 panic，main 裡的 recover 攔不到，整個程式還是會結束。所以每個 goroutine 如果需要保護，要各自 defer recover。

第三，也是最重要的：panic 和 recover 不是 try-catch 的替代品。正常可預期的錯誤，例如檔案不存在、輸入格式錯誤，一律用 error 回傳。recover 只用在「程式的邊界」，防止一個 bug 讓整個服務停擺。
-->

---
layout: default
---

# 練習 2：安全執行器
### 任務說明

1. 寫一個函式 `safeRun(name string, task func()) (err error)`
2. 在 `safeRun` 中用 `defer` + `recover()`，把 task 裡發生的 panic 轉成 error：`"任務 name 發生 panic：原因"`
3. 準備三個任務：
   - 正常印出「Hello」
   - 存取 `nil map`（寫入）
   - 呼叫 `panic("自訂錯誤")`
4. 依序執行三個任務，印出每個任務的結果，確認程式最後有印出「全部完成」

<!--
這個練習是 recover 最實際的應用：一個「安全執行器」，不管傳進來的任務會不會 panic，都能把它轉成 error，程式繼續往下執行。

第 16 章寫背景工作的 goroutine 時，就會用到同樣的模式。
-->

---
zoom: 0.79
---

# 練習 2：解題提示
### 提示說明

```go
package main

import "fmt"

func safeRun(name string, task func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("任務 %s 發生 panic：%v", name, r)
		}
	}()
	task()
	return nil
}

func main() {
	tasks := map[string]func(){
		"hello": func() { fmt.Println("Hello") },
		"nilmap": func() {
			var m map[string]int
			m["a"] = 1
		},
		"custom": func() { panic("自訂錯誤") },
	}
	for _, name := range []string{"hello", "nilmap", "custom"} {
		fmt.Println(name, "→", safeRun(name, tasks[name]))
	}
	fmt.Println("全部完成")
}
```

<!--
safeRun 的結構跟剛剛的 safeDivide 一模一樣：defer 一個閉包，在裡面 recover，有 panic 就轉成 error 指定給具名回傳值。

任務存在 map 裡，因為 map 的走訪順序不固定，所以另外用一個切片決定執行順序。

執行後：hello 印出 Hello 並回傳 nil；nilmap 回傳「assignment to entry in nil map」；custom 回傳「自訂錯誤」；最後印出全部完成。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 處理 error 與 panic 的指導方針
## Guidelines

<!--
最後，我們把今天的內容整理成幾條 Go 社群公認的指導方針。
-->

---

# 處理 error 與 panic 的指導方針

| 情境 | 做法 |
| --- | --- |
| 可預期的失敗（檔案不存在、輸入錯誤、網路逾時） | 回傳 `error` |
| 程式的 bug、不可能發生的狀態 | `panic` |
| 啟動時的必要設定失敗（寫死的正規表達式、模板） | `Must` 函式 / `panic` |
| 需要在錯誤訊息加上情境 | `fmt.Errorf("做某事失敗：%w", err)` |
| 呼叫端需要判斷錯誤種類 | 哨兵錯誤 + `errors.Is`；自訂型別 + `errors.As` |
| 伺服器、goroutine 的邊界 | `defer` + `recover`，轉成 error 並記錄 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>絕對不要：</b> 用 <code>_</code> 忽略 error（<code>n, _ := strconv.Atoi(s)</code>）— 除非你<b>真的確定</b>不會出錯，並加上註解說明原因。
</div>

<!--
這張表是今天最重要的一頁，建議大家拍下來。

判斷的核心問題是：「這個錯誤是可預期的嗎？」檔案可能不存在、使用者可能輸入錯誤、網路可能逾時，這些都是正常會發生的事，用 error 回傳，讓呼叫端決定怎麼處理。

如果是「程式寫錯了」才會發生的狀況，就用 panic，讓問題盡快浮現。

最後，千萬不要用底線忽略 error。這是 Go 程式碼裡最常見的壞習慣，當程式出問題時，被忽略的錯誤會讓我們完全找不到原因。
-->

---
layout: default
---

# 綜合練習：設定檔解析器
### 任務說明

解析一組 `key=value` 格式的設定：`lines := []string{"port=8080", "debug=true", "timeout=abc", "name"}`

1. 宣告哨兵錯誤 `ErrBadFormat`（缺少 `=`）
2. 寫 `parseLine(line string) (key, value string, err error)`：用 `strings.Cut`，缺少 `=` 時回傳用 `%w` 包裝的 `ErrBadFormat`
3. 寫 `parseAll(lines []string) (map[string]string, error)`：解析每一行，**收集所有錯誤**，最後用 `errors.Join` 一起回傳
4. `timeout` 的值要能轉成整數，否則也算錯誤（包裝 `strconv.Atoi` 的錯誤）
5. 印出成功解析的設定，以及合併後的錯誤；用 `errors.Is` 確認其中包含 `ErrBadFormat`

<!--
這個綜合練習模擬真實世界的設定檔解析：不是遇到第一個錯誤就停下來，而是把所有錯誤都收集起來，一次告訴使用者哪裡要修改。

用到的觀念有：strings.Cut、哨兵錯誤、%w 包裝、errors.Join 合併、errors.Is 判斷。
-->

---
zoom: 0.91
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadFormat = errors.New("格式錯誤，缺少 =")

func parseLine(line string) (key, value string, err error) {
	key, value, ok := strings.Cut(line, "=")
	if !ok {
		return "", "", fmt.Errorf("第 %q 行：%w", line, ErrBadFormat)
	}
	if key == "timeout" {
		if _, err := strconv.Atoi(value); err != nil {
			return "", "", fmt.Errorf("timeout 必須是整數：%w", err)
		}
	}
	return key, value, nil
}
```

<!--
parseLine 用 strings.Cut 切開 key 和 value，第三個回傳值 ok 代表有沒有找到等號。找不到就回傳包裝過的 ErrBadFormat。

timeout 的值額外檢查能不能轉成整數，不能的話包裝 strconv 的錯誤回傳。注意這裡 if 起始賦值裡的 err 是一個新的區域變數，會遮蔽具名回傳值 err，這是第一章講的變數遮蔽，這裡是刻意的，因為我們最後是明確寫出回傳值。
-->

---
zoom: 0.94
---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func parseAll(lines []string) (map[string]string, error) {
	cfg := map[string]string{}
	var errs []error
	for _, line := range lines {
		k, v, err := parseLine(line)
		if err != nil {
			errs = append(errs, err) // 收集錯誤，繼續處理下一行
			continue
		}
		cfg[k] = v
	}
	return cfg, errors.Join(errs...) // 沒有錯誤時 Join 回傳 nil
}

func main() {
	lines := []string{"port=8080", "debug=true", "timeout=abc", "name"}
	cfg, err := parseAll(lines)
	fmt.Println(cfg)
	fmt.Println(err)
	fmt.Println("包含格式錯誤？", errors.Is(err, ErrBadFormat)) // true
}
```

<!--
parseAll 用一個 []error 切片收集錯誤，遇到錯誤就 append 然後 continue，繼續處理下一行。

最後用 errors.Join 把所有錯誤合併成一個。一個很貼心的設計：如果切片是空的，Join 會回傳 nil，所以不用另外判斷「有沒有錯誤」。

合併後的錯誤印出來，每個錯誤各佔一行。而且 errors.Is 可以穿透 Join，檢查裡面任何一個錯誤，所以能判斷出包含 ErrBadFormat。
-->

---

# 章節總結

- **三種錯誤**：語法錯誤（編譯器抓）、執行期間錯誤（panic）、邏輯錯誤（靠測試）
- **錯誤是值**：Go 不用例外，函式回傳 `error`，呼叫端 `if err != nil` 明確處理
- **error 介面**：只有 `Error() string`；用 `errors.New`、`fmt.Errorf` 建立
- **包裝與檢查**：`%w` 包裝加上情境；`errors.Is` 比對哨兵錯誤、`errors.As` / `errors.AsType` 取出錯誤型別；`errors.Join` 合併
- **panic**：程式的 bug、無法繼續的狀況；`Must` 函式慣例
- **recover**：只在 `defer` 中有效，只攔截同一個 goroutine；用在伺服器與 goroutine 的邊界
- **方針**：可預期的失敗回傳 error；不要忽略 error；錯誤只處理一次

下一章我們會正式介紹「介面」：Go 最強大的抽象工具。

<!--
我們來整理今天學到的東西。

Go 的錯誤處理核心就是一句話：錯誤是值。函式回傳 error，呼叫端明確處理。建立錯誤用 errors.New 和 fmt.Errorf，用 %w 包裝錯誤加上情境，用 errors.Is 和 errors.As 檢查錯誤種類。panic 留給真正的 bug，recover 用在程式的邊界。

今天我們看到 error 是一個「介面」：任何有 Error 方法的型別都是 error。這個「只要有某個方法，就符合某個介面」的概念，就是下一章的主題。介面是 Go 最強大、也最有特色的功能，它讓 Go 不需要繼承，也能寫出非常有彈性的程式。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
錯誤處理是 Go 程式設計師每天都在做的事，今天的內容一定要熟練。

課後建議：回頭看看前幾章的範例，找出所有用底線忽略 error 的地方，試著把它們改成正確的錯誤處理。

有問題的同學現在可以提問！
-->

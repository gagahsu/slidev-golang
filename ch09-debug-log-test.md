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
title: 程式除錯：格式化訊息、日誌與單元測試
routeAlias: ch09
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
  <h1 style="color: #1a5c5c; font-size: 3.2rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">程式除錯</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「格式化訊息、日誌與單元測試 — 讓 bug 無所遁形」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第九章！

前面八章我們學完了 Go 的語法和專案組織方式。從這一章開始進入實戰，而實戰的第一課，就是「怎麼找出程式的錯誤」。

寫程式的時間，其實有很大一部分是在除錯。今天會介紹三個工具：fmt 套件的格式化輸出，讓我們把變數看得清清楚楚；log 和 slog 套件，讓程式在正式環境留下紀錄；testing 套件的單元測試，讓電腦自動幫我們檢查程式對不對。
-->

---
layout: default
---

# Outline

- **前言** — 臭蟲的發生原因、除錯原則
- **以 fmt 套件做格式化輸出** — 格式化動詞、寬度與精度、浮點數格式、`strconv.FormatFloat`
- **使用 log 提供追蹤訊息／日誌** — `log` 套件、自訂 logger、結構化日誌 `log/slog`
- **撰寫單元測試** — `testing` 套件、表格驅動測試、覆蓋率、效能測試
- **GoShop 專案實作** — 第 9 步：替 GoShop 寫測試、加上日誌
- **章節總結**

<!--
今天的內容分成三大工具。

fmt 是我們從第一天就在用的套件，今天會把它的格式化功能完整學一遍。log 用來記錄程式執行的過程，除了傳統的 log 套件，還會介紹 Go 1.21 加入的結構化日誌 slog，這是現在的主流做法。

最後是單元測試，這是今天最重要的部分。會寫測試，才能放心地修改程式碼。
-->

---

# 回顧：套件

- 一個資料夾就是一個**套件**；**模組**的根目錄有 `go.mod`
- **大寫開頭匯出、小寫開頭私有**
- 常用指令：`go mod init`、`go get`、`go mod tidy`
- 測試檔 `xxx_test.go` 可以使用 `package xxx_test` ← 今天會用到

<!--
回顧一下上一章。

我們學了套件和模組，知道怎麼把程式碼拆成好幾個套件。上一章有提到一個例外：_test.go 結尾的測試檔。今天就會正式介紹測試檔的寫法，以及 go test 這個指令。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## Why Bugs Happen

<!--
先來聊聊臭蟲是怎麼來的，以及除錯的基本原則。
-->

---

# 臭蟲的發生原因

| 原因 | 範例 |
| --- | --- |
| **邏輯錯誤** | 條件寫反（`>` 寫成 `>=`）、整數除法忘記轉型 |
| **沒考慮到的輸入** | 空切片、負數、空字串、中文字、超大數字 |
| **誤解 API 的行為** | 以為 `append` 不會改到原本的陣列、以為 map 有順序 |
| **忽略錯誤** | `n, _ := strconv.Atoi(s)`，轉換失敗時 `n` 是 0 |
| **並行問題** | 多個 goroutine 同時修改同一個變數（Ch 16） |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
🐛 <b>「Bug」的由來：</b> 1947 年，工程師在哈佛 Mark II 電腦的繼電器裡找到一隻飛蛾，並把它貼在工作日誌上，寫下「First actual case of bug being found」。
</div>

<!--
臭蟲是怎麼產生的？

表格列出了最常見的五種原因。前四種我們在前幾章都遇過了：第 6 章的整數除法、第 3 章的中文字串長度、第 4 章的切片共用底層陣列、以及忽略錯誤。第五種並行問題會在第 16 章介紹。

至於為什麼叫「bug」？1947 年，工程師在一台早期電腦的繼電器裡，真的找到了一隻卡住的飛蛾，於是把它貼在工作日誌上，寫下「第一個真的找到蟲的案例」。從此以後，程式的錯誤就被叫做 bug，除錯就叫 debug。
-->

---

# 除錯原則

| 步驟 | 做法 |
| --- | --- |
| **1. 重現問題** | 找出「一定會出錯」的輸入與步驟，無法重現就無法確認已修好 |
| **2. 讀錯誤訊息** | 看清楚 panic 訊息、stack trace、錯誤包裝鏈，從**第一個自己的程式碼**開始看 |
| **3. 縮小範圍** | 用印出訊息、中斷點、二分法，找出哪一行開始不對 |
| **4. 提出假設、一次改一個地方** | 同時改很多地方，就不知道是哪一個修好的 |
| **5. 寫測試固定下來** | 修好之後，把這個案例寫成測試，避免同樣的 bug 再回來 |

<!--
除錯就像偵探辦案，有一套固定的步驟。

第一步是重現問題：找出「每次都會出錯」的條件。如果問題時好時壞，就很難確認修好了沒。

第二步是仔細讀錯誤訊息。很多人看到一長串紅字就慌了，其實 Go 的錯誤訊息都寫得很清楚，特別是第 6 章學的錯誤包裝鏈，會一層一層告訴我們在做什麼的時候出錯。

第三、四步是縮小範圍、一次只改一個地方。最後一步最重要：把這個 bug 寫成測試。這樣以後不管誰修改程式碼，只要同樣的 bug 又出現，測試就會立刻失敗。今天最後一段就會教怎麼寫測試。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 以 fmt 套件做格式化輸出
## Formatted Output

<!--
第一個工具：fmt 套件的格式化輸出。
-->

---

# fmt 套件的函式家族

| 函式 | 輸出到 | 說明 |
| --- | --- | --- |
| `Print` / `Println` / `Printf` | 標準輸出（螢幕） | `ln` 會加空白與換行；`f` 使用格式化字串 |
| `Sprint` / `Sprintln` / `Sprintf` | **回傳字串** | 組字串時使用 |
| `Fprint` / `Fprintln` / `Fprintf` | 任何 `io.Writer` | 寫到檔案、HTTP 回應、`strings.Builder` |
| `Errorf` | **回傳 error** | 支援 `%w` 包裝錯誤（Ch 6） |
| `Scan` / `Scanln` / `Scanf` | 從標準輸入讀取 | 讀取使用者輸入 |

<!--
fmt 套件的函式有一個很好記的命名規則。

前綴決定「輸出到哪裡」：沒有前綴是印到螢幕，S 開頭是 String，回傳字串，F 開頭是 File，寫到任何 io.Writer。

後綴決定「怎麼輸出」：沒有後綴是直接輸出，ln 會在值之間加空白並在最後換行，f 代表使用格式化字串。

所以 Fprintf 就是「用格式化字串寫到某個 Writer」。只要記住這個規則，整個 fmt 套件就不用背了。
-->

---

# fmt 的格式化動詞（verbs）

| 動詞 | 用途 | 範例 → 結果 |
| --- | --- | --- |
| `%v` | 預設格式（萬用） | `P{"Go", 16}` → `{Go 16}` |
| `%+v` | 加上欄位名稱 | → `{Name:Go Age:16}` |
| `%#v` | Go 語法的表示法 | → `main.P{Name:"Go", Age:16}` |
| `%T` | 型別 | → `main.P` |
| `%d` / `%x` / `%b` | 十進位／十六進位／二進位整數 | `255` → `255` / `ff` / `11111111` |
| `%s` / `%q` | 字串／加上雙引號的字串 | `"hi"` → `hi` / `"hi"` |
| `%t` / `%p` / `%%` | 布林／指標位址／百分比符號 | `true` → `true` |

<!--
格式化動詞是 Printf 裡以百分比符號開頭的佔位符。

最常用的是 %v，它是「萬用」的，任何型別都能印。除錯的時候，%+v 和 %#v 特別好用：%+v 會印出 struct 的欄位名稱，%#v 會用 Go 程式碼的格式印出來，連型別名稱和字串的引號都有，一眼就能看出每個欄位的值和型別。

%T 印出型別，第一章就用過了。%q 會替字串加上雙引號，可以看出字串前後有沒有多餘的空白。
-->

---

# 格式化動詞 — 除錯範例

```go
package main

import "fmt"

type Order struct {
	ID    int
	Item  string
	Price float64
	Tags  []string
}

func main() {
	o := Order{ID: 7, Item: "咖啡 ", Price: 60}
	fmt.Printf("%v\n", o) // {7 咖啡  60 []}
	// {ID:7 Item:咖啡  Price:60 Tags:[]}
	fmt.Printf("%+v\n", o)
	// main.Order{ID:7, Item:"咖啡 ", Price:60, Tags:[]string(nil)}
	fmt.Printf("%#v\n", o)
	fmt.Printf("%q %T\n", o.Item, o.Price) // "咖啡 " float64
}
```

<!--
這段程式碼的目的，是比較三種 %v 在除錯時的差別。

%v 只印出值，Item 後面多了一個空白，但很難看出來。%+v 加上欄位名稱，比較清楚。%#v 最詳細：字串有引號，一眼就能看出「咖啡」後面多了一個空白；Tags 顯示 []string(nil)，告訴我們它是 nil 切片，而不是空切片。

所以除錯的時候，印 struct 建議用 %+v 或 %#v。
-->

---

# 寬度、精度與旗標

語法：`%[旗標][寬度][.精度]動詞`

| 格式 | 意義 | `42` / `3.14159` 的結果 |
| --- | --- | --- |
| `%5d` | 寬度 5，靠右對齊 | `[   42]` |
| `%-5d` | 寬度 5，靠左對齊（`-` 旗標） | `[42   ]` |
| `%05d` | 寬度 5，左邊補 0（`0` 旗標） | `[00042]` |
| `%.2f` | 小數點後 2 位 | `[3.14]` |
| `%8.2f` | 總寬度 8、小數 2 位 | `[    3.14]` |
| `%+d` | 永遠顯示正負號 | `[+42]` |

<!--
格式化動詞前面還可以加上寬度、精度和旗標，用來控制排版。

寬度是「至少佔幾個字元」，不足的部分預設在左邊補空白，也就是靠右對齊。加上減號旗標會變成靠左對齊，加上 0 旗標會補 0 而不是空白。

精度寫在小數點後面，對浮點數來說是「小數點後幾位」。

這些在印出表格、報表的時候特別好用，第 2 章九九乘法表的 %-3d 就是這個用法。
-->

---

# 印出浮點數的進階格式化

```go
package main

import "fmt"

func main() {
	price := 1234.5678
	// [1234.567800]：預設 6 位小數
	fmt.Printf("[%f]\n", price)
	fmt.Printf("[%.1f]\n", price) // [1234.6]
	// [1.234568e+03]：科學記號
	fmt.Printf("[%e]\n", price)
	// [1.234567e+06]：自動選擇較短的格式
	fmt.Printf("[%g]\n", 1234567.0)
	// [   1234.57|1234.57   ]
	fmt.Printf("[%10.2f|%-10.2f]\n", price, price)

	fmt.Printf("%.2f %.0f %.0f\n", 2.675, 2.5, 3.5) // 2.67 2 4 ⚠️
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>四捨五入的陷阱：</b> <code>2.675</code> 在二進位中其實是 <code>2.67499999…</code>；而剛好在一半的 <code>2.5</code> 會採「<b>銀行家捨入</b>」到最近的偶數。需要精準的四捨五入時，用 <code>math.Round</code> 或整數運算。
</div>

<!--
浮點數有幾種不同的格式。%f 是一般的小數格式，預設 6 位小數；%e 是科學記號；%g 會自動選擇比較短的格式，數字很大或很小時用科學記號，其他時候用一般格式。

最後一行是一個很多人不知道的陷阱。2.675 取兩位，我們預期是 2.68，但印出來是 2.67，因為 2.675 在二進位裡其實是 2.674999…，第 3 章講過浮點數的精度問題。

另外 2.5 取整數印出 2，3.5 印出 4。這是因為剛好在一半的時候，Go 會捨入到最近的偶數，這叫做「銀行家捨入法」，可以避免大量計算時的累積誤差。所以印出金額時，不要依賴 Printf 來四捨五入。
-->

---

# 用 strconv.FormatFloat() 格式化浮點數

需要把浮點數轉成**字串**（而不是印出來）時，`strconv.FormatFloat` 比 `Sprintf` 更快

語法：`strconv.FormatFloat(值, 格式, 精度, 位元大小)`

| 呼叫 | 結果 | 說明 |
| --- | --- | --- |
| `FormatFloat(3.14159, 'f', 2, 64)` | `"3.14"` | 一般小數，2 位 |
| `FormatFloat(1234.5678, 'e', 3, 64)` | `"1.235e+03"` | 科學記號，3 位 |
| `FormatFloat(0.1, 'f', -1, 64)` | `"0.1"` | 精度 `-1`：**最少位數但能精確還原** |
| `FormatFloat(1e21, 'g', -1, 64)` | `"1e+21"` | 自動選擇格式 |

<!--
如果我們要的是一個字串，而不是印出來，例如要存進資料庫、寫進 CSV 檔，可以用 strconv.FormatFloat，它比 Sprintf 更快，因為不需要解析格式化字串。

四個參數分別是：要轉換的值、格式（'f' 一般小數、'e' 科學記號、'g' 自動）、精度、以及原本是 float32 還是 float64。

特別介紹精度 -1：它的意思是「用最少的位數，但保證轉回來的數字完全一樣」。0.1 會變成 "0.1"，而不是 "0.100000"。這在序列化資料時非常好用，第 11 章的 JSON 編碼內部就是用這個方式處理浮點數。
-->

---
layout: default
---

# 練習 1：商品報表
### 任務說明

印出如下格式的商品報表（名稱靠左 10 格、數量靠右 5 格、價格靠右 10 格且兩位小數）：

```text
名稱            數量       單價
----------------------------
Coffee        2      60.00
Cake          1     120.50
Sandwich     10      85.00
----------------------------
合計                1,152.50
```

提示：用 `%-10s`、`%5d`、`%10.2f`；千分位可以自己寫一個 `withComma` 函式處理整數部分

<!--
這個練習要用格式化動詞排出一個整齊的報表。

注意：中文字在終端機裡佔兩個字元寬，但 %-10s 是用「字元數」計算的，所以中文的對齊會跑掉。這就是為什麼題目的商品名稱用英文。

千分位 Go 的 fmt 沒有內建，需要自己處理：把整數部分轉成字串，從右邊每三位插入一個逗號。
-->

---
zoom: 0.84
---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func withComma(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64) // "1152.50"
	intPart, frac, _ := strings.Cut(s, ".")
	var b strings.Builder
	for i, ch := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	return b.String() + "." + frac
}

func main() {
	fmt.Printf("%-10s%5d%10.2f\n", "Coffee", 2, 60.0)
	// 合計           1,152.50
	fmt.Printf("%-10s%15s\n", "合計", withComma(1152.5))
}
```

<!--
withComma 先用 FormatFloat 轉成兩位小數的字串，再用 strings.Cut 切出整數部分和小數部分。

走訪整數部分的每個字元，當「剩下的位數」是 3 的倍數時，就插入一個逗號。例如 1152，走到索引 1 的時候，剩下 3 位，插入逗號，變成 1,152。

每一列用 %-10s 讓名稱靠左、%5d 讓數量靠右、%10.2f 讓價格靠右並保留兩位小數。完整的報表就是把每個商品都用同樣的格式印一次。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 使用 log 提供追蹤訊息／日誌
## Logging

<!--
第二個工具：日誌。
-->

---

# 印出追蹤訊息

最簡單的除錯方式：在關鍵位置**印出變數的值**，追蹤程式的執行流程

```go
package main

import "fmt"

func discount(total int, isMember bool) int {
	fmt.Printf("DEBUG discount() total=%d member=%t\n", total, isMember)
	if isMember && total > 1000 {
		total = total * 9 / 10
	}
	fmt.Printf("DEBUG discount() 回傳 %d\n", total)
	return total
}

func main() {
	fmt.Println(discount(1200, true)) // 1080
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>缺點：</b> 忘記刪除會混在正常輸出中；沒有時間、沒有檔名行號；無法依等級開關。正式的程式請改用 <b>log / slog</b>。
</div>

<!--
最簡單、也最常用的除錯方法，就是在關鍵的地方印出變數的值，看看程式執行到哪裡、變數是什麼。這個方法常被戲稱為「printf 除錯法」。

它很直覺，但有幾個缺點：除錯訊息跟正常的輸出混在一起，除完錯常常忘了刪；沒有時間戳記，不知道事情是什麼時候發生的；也沒辦法在正式環境關掉。

所以臨時除錯用 fmt 沒問題，但如果要在程式裡長期留下紀錄，就要用專門的日誌工具：log 或 slog。
-->

---

# 用 log 套件輸出日誌

`log` 會自動加上**時間戳記**，並輸出到**標準錯誤（stderr）**

| 函式 | 行為 |
| --- | --- |
| `log.Println` / `log.Printf` | 印出日誌訊息 |
| `log.Fatal` / `log.Fatalf` | 印出訊息後**立刻結束程式**（`os.Exit(1)`），**defer 不會執行** |
| `log.Panic` / `log.Panicf` | 印出訊息後 **panic**（defer 會執行、可被 recover） |
| `log.SetFlags(...)` | 設定前綴格式：日期、時間、微秒、檔名行號 |
| `log.SetPrefix("[app] ")` | 設定每一行的前綴 |

```go
log.SetFlags(log.LstdFlags | log.Lshortfile) // 日期時間 + 檔名:行號
log.Println("伺服器啟動，port =", 8080)
// 2026/09/24 10:30:00 main.go:12: 伺服器啟動，port = 8080
```

<!--
Go 的標準函式庫有一個 log 套件，用法跟 fmt 很像，但它會自動在每一行前面加上時間，而且輸出到標準錯誤 stderr，跟正常輸出分開。

log.Fatal 要特別注意：它印出訊息後會直接結束程式，而且 defer 不會執行！所以 Fatal 只適合用在 main 函式裡、程式啟動就失敗的時候，例如設定檔讀不到。在其他函式裡，請回傳 error。

SetFlags 可以設定要顯示哪些資訊，加上 Lshortfile 就會顯示檔名和行號，除錯時非常有用。
-->

---
zoom: 0.96
---

# 建立自訂 logger 物件

用 `log.New(輸出目的地, 前綴, 旗標)` 建立獨立的 logger，例如**不同模組、寫到檔案**

```go
package main

import (
	"log"
	"os"
)

func main() {
	flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	f, err := os.OpenFile("app.log", flag, 0o644)
	if err != nil {
		log.Fatal(err) // main 裡啟動失敗，才適合用 Fatal
	}
	defer f.Close()

	orderLog := log.New(f, "[order] ", log.LstdFlags|log.Lmsgprefix)
	payLog := log.New(os.Stderr, "[pay] ", log.LstdFlags|log.Lshortfile)

	orderLog.Println("建立訂單 #7")   // 寫進 app.log
	payLog.Printf("付款 %d 元", 350) // 印到終端機
}
```

<!--
log.New 可以建立一個獨立的 logger，有自己的輸出目的地、前綴和格式。

第一個參數是 io.Writer，第 7 章學過，任何實作了 Write 方法的東西都可以：檔案、螢幕、網路連線。這裡 orderLog 寫進 app.log 檔案，payLog 印到終端機。

開檔案的 os.OpenFile 是第 12 章的內容，這裡先知道它會開啟（或建立）一個檔案，並以附加的方式寫入。注意 Lmsgprefix 旗標會把前綴放在訊息前面、時間後面，比較好讀。
-->

---

# 結構化日誌：log/slog（Go 1.21+）

現代的服務會把日誌送到收集系統分析，**「鍵值對」格式**比純文字更容易搜尋與統計

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.Info("使用者登入", "user", "alice", "id", 42)
	// 2026/09/24 10:30:00 INFO 使用者登入 user=alice id=42

	opts := &slog.HandlerOptions{Level: slog.LevelDebug} // 最低輸出等級
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	logger.Debug("查詢資料庫", "sql", "SELECT 1", "ms", 3)
	logger.Error("付款失敗", "order", 7,
		slog.String("reason", "餘額不足"))
	// {"time":"...","level":"ERROR","msg":"付款失敗",
	//  "order":7,"reason":"餘額不足"}
}
```

<!--
Go 1.21 加入了 log/slog 套件，s 代表 structured，結構化。這是現在 Go 寫日誌的主流方式。

傳統的 log 輸出的是一行純文字，人讀起來很方便，但機器很難分析。正式環境的服務通常會把日誌送到收集系統，例如 Elasticsearch、Loki、雲端的日誌服務，這時候「鍵值對」的格式就方便多了：可以直接搜尋「user 等於 alice 的所有日誌」，或統計「每分鐘有幾個 ERROR」。

slog.Info 後面第一個參數是訊息，後面接著成對的鍵和值。用 NewJSONHandler 可以輸出成 JSON 格式，也可以用 NewTextHandler 輸出成 key=value 的純文字格式。
-->

---

# slog 的日誌等級

| 等級 | 函式 | 使用時機 |
| --- | --- | --- |
| `DEBUG` | `slog.Debug` | 開發除錯用的詳細資訊，正式環境通常關閉 |
| `INFO` | `slog.Info` | 正常的重要事件：啟動、登入、訂單建立 |
| `WARN` | `slog.Warn` | 不正常但還能處理：重試、使用預設值 |
| `ERROR` | `slog.Error` | 發生錯誤，需要有人注意 |

```go
slog.SetDefault(logger)                     // 設定全域預設的 logger
// 附加固定欄位，之後每筆日誌都會帶著
reqLog := logger.With("request_id", "a1b2")
reqLog.Info("處理請求")
```

<!--
slog 有四個等級：DEBUG、INFO、WARN、ERROR。設定最低等級之後，比它低的日誌就不會輸出。例如正式環境設成 INFO，所有的 DEBUG 日誌就自動關掉了，不需要改程式碼。

logger.With 是一個很實用的功能：它會回傳一個新的 logger，之後每一筆日誌都會自動帶上指定的欄位。例如在處理一個 HTTP 請求時，加上 request_id，這個請求的所有日誌都會帶著同一個 ID，出問題時就能把相關的日誌全部串起來。

SetDefault 可以把自訂的 logger 設成全域預設，之後呼叫 slog.Info 就會使用它。
-->

---
layout: default
---

# 練習 2：訂單處理日誌
### 任務說明

1. 建立一個輸出 **JSON 格式**、等級為 **DEBUG** 的 slog logger，並設為預設
2. 寫函式 `process(id, amount int) error`：
   - 開始時用 `Debug` 記錄 `id` 與 `amount`
   - 金額 ≤ 0 時用 `Warn` 記錄，並回傳錯誤
   - 成功時用 `Info` 記錄「訂單處理完成」
3. 在 `main` 中處理訂單 `(1, 300)`、`(2, 0)`；失敗時用 `Error` 記錄，欄位名稱為 `err`
4. 改成 `Level: slog.LevelInfo` 再執行一次，觀察 Debug 訊息消失

<!--
這個練習要大家體會 slog 的等級和結構化輸出。

完成之後，試著改變等級再執行一次，觀察輸出的變化。這就是正式環境和開發環境使用不同等級的做法。
-->

---
zoom: 0.79
---

# 練習 2：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"log/slog"
	"os"
)

func process(id, amount int) error {
	slog.Debug("開始處理訂單", "id", id, "amount", amount)
	if amount <= 0 {
		slog.Warn("金額不合法", "id", id, "amount", amount)
		return errors.New("金額必須大於 0")
	}
	slog.Info("訂單處理完成", "id", id)
	return nil
}

func main() {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	h := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(h))
	for _, o := range [][2]int{{1, 300}, {2, 0}} {
		if err := process(o[0], o[1]); err != nil {
			slog.Error("訂單失敗", "id", o[0], "err", err)
		}
	}
}
```

<!--
先建立 JSON handler，設定等級為 Debug，再用 SetDefault 設成預設 logger，之後直接呼叫 slog.Debug、slog.Info 就會使用它。

測試資料用了 [2]int 陣列的切片，每個元素是一對 id 和金額。

第二筆訂單的金額是 0，會先記錄一筆 WARN，回到 main 後再記錄一筆 ERROR。注意這裡其實違反了第 6 章的「錯誤只處理一次」，WARN 和 ERROR 記錄了同一件事。實務上可以只保留其中一個，這裡是為了示範不同的等級。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 撰寫單元測試
## Unit Testing

<!--
第三個工具，也是今天最重要的：單元測試。
-->

---

# 什麼是單元測試？

「**單元測試是用程式碼來檢查程式碼**：給定輸入，驗證輸出是否符合預期。」

| Go 測試的規則 | 說明 |
| --- | --- |
| **檔名** | 以 `_test.go` 結尾，例如 `calc_test.go`，`go build` 不會編譯它 |
| **函式名稱** | `func TestXxx(t *testing.T)`，`Test` 後面接大寫開頭的名稱 |
| **回報失敗** | `t.Errorf(...)`：記錄失敗、繼續執行；`t.Fatalf(...)`：記錄失敗、立刻停止 |
| **執行** | `go test ./...` 執行所有套件的測試；`-v` 顯示詳細結果 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>內建就有：</b> 不像其他語言需要另外安裝 JUnit、pytest、Jest，Go 的測試框架就在標準函式庫的 <code>testing</code> 套件裡。
</div>

<!--
什麼是單元測試？

想像我們開了一家飲料店，每次調整配方後，都要自己試喝一次看看味道對不對。如果有一百種飲料，每次改配方都要全部試喝一遍，很累人。單元測試就像是一台「自動試喝機」：我們先寫好「這杯飲料應該是什麼味道」，之後每次修改，按一個按鈕，機器就自動把一百種飲料全部試喝一遍，告訴我們哪幾杯味道不對。

Go 的測試規則很簡單：檔名以 _test.go 結尾，函式名稱以 Test 開頭，參數是 *testing.T。然後執行 go test 就好了。而且測試框架是內建的，不需要安裝任何東西。
-->

---
zoom: 0.86
---

# 第一個單元測試

```go
// 檔名：calc.go
package calc

// Average 回傳平均值；空切片回傳 0。
func Average(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums))
}
```

```go
// 檔名：calc_test.go
package calc

import "testing"

func TestAverage(t *testing.T) {
	got := Average([]int{1, 2})
	if got != 1.5 {
		t.Errorf("Average([1 2]) = %v，期望 1.5", got)
	}
}
```

<!--
這是一個最簡單的測試。

calc.go 裡有一個 Average 函式，就是第 6 章那個整數除法 bug 修好之後的版本。calc_test.go 是它的測試檔，跟 calc.go 放在同一個資料夾、同一個套件。

TestAverage 呼叫 Average，把結果存在 got，然後跟期望的值比較，不一樣就用 t.Errorf 回報失敗。Go 社群的慣例是用 got 和 want 當變數名稱，錯誤訊息的格式是「函式(輸入) = 實際值，期望 期望值」。
-->

---

# 執行測試：go test

```bash
go test            # 執行目前套件的測試
# ok  	example.com/calc	0.002s

go test -v         # 顯示每個測試的詳細結果
# === RUN   TestAverage
# --- PASS: TestAverage (0.00s)
# PASS

go test ./...      # 執行模組中所有套件的測試
go test -run Avg   # 只執行名稱符合 "Avg" 的測試
```

測試失敗時的輸出：

```text
--- FAIL: TestAverage (0.00s)
    calc_test.go:7: Average([1 2]) = 1，期望 1.5
FAIL
exit status 1
```

<!--
寫好測試之後，用 go test 執行。

沒有加任何參數時，全部通過只會顯示一行 ok。加上 -v 會顯示每個測試的名稱和結果。./... 是「目前資料夾以及所有子資料夾」的意思，會執行整個模組的測試。-run 可以用名稱篩選，只執行特定的測試。

如果測試失敗，會顯示 FAIL、測試檔的行號，以及我們在 t.Errorf 裡寫的錯誤訊息。這就是為什麼錯誤訊息要寫清楚：失敗的時候，看訊息就知道輸入是什麼、實際得到什麼、期望的是什麼。
-->

---
zoom: 0.85
---

# 表格驅動測試（Table-Driven Tests）

Go 社群最主流的測試寫法：**把測試案例放進切片，用迴圈 + `t.Run` 逐一執行**

```go
// 檔名：calc_table_test.go
package calc

import "testing"

func TestAverageTable(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want float64
	}{
		{"一般情況", []int{1, 2}, 1.5},
		{"單一元素", []int{7}, 7},
		{"空切片", nil, 0},
		{"負數", []int{-1, -3}, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // 子測試：每個案例獨立回報
			if got := Average(tt.nums); got != tt.want {
				t.Errorf("Average(%v) = %v，期望 %v",
					tt.nums, got, tt.want)
			}
		})
	}
}
```

<!--
這是 Go 社群最主流的測試寫法：表格驅動測試。

我們把所有的測試案例放進一個匿名結構的切片，每個案例有名稱、輸入和期望的輸出，就像一張表格。然後用迴圈走訪每個案例，用 t.Run 建立子測試。

這樣寫有兩個好處。第一，新增測試案例只要加一行，不用複製貼上整段程式碼。第二，t.Run 讓每個案例獨立回報結果，某一個失敗了，其他的還是會繼續執行，而且輸出會顯示是哪個案例失敗。

特別注意「空切片」和「負數」這種邊界案例，bug 最常躲在這些地方。
-->

---
zoom: 0.97
---

# 測試覆蓋率與效能測試

```bash
go test -cover                        # 顯示覆蓋率
# coverage: 100.0% of statements

# 用瀏覽器看哪些行沒被測到
go test -coverprofile=c.out && go tool cover -html=c.out
```

```go
// 檔名：calc_bench_test.go
package calc

import "testing"

func BenchmarkAverage(b *testing.B) {
	nums := make([]int, 1000)
	for b.Loop() { // Go 1.24+：自動決定要跑幾次
		Average(nums)
	}
}
```

```bash
go test -bench=. -benchmem
# BenchmarkAverage-8  2400000  495.3 ns/op  0 B/op  0 allocs/op
```

<!--
測試寫完之後，怎麼知道測得夠不夠完整？可以看「覆蓋率」：go test -cover 會告訴我們有多少比例的程式碼被測試執行到。加上 -coverprofile 再用 go tool cover -html，會開啟瀏覽器，用綠色和紅色標出哪些行有被測到、哪些沒有。

另外 testing 套件也支援效能測試：函式名稱以 Benchmark 開頭，參數是 *testing.B。Go 1.24 加入了 b.Loop()，它會自動決定要跑幾次才能得到穩定的結果，以前的寫法是 for i := 0; i < b.N; i++。

go test -bench=. 執行效能測試，結果會顯示每次執行花了多少奈秒，加上 -benchmem 還會顯示記憶體配置的次數。
-->

---

# 使用單元測試的注意事項

**注意事項之一：** 測試**公開行為**，不要測試內部實作細節，否則重構時測試會一直壞掉

**注意事項之二：** 測試要**獨立、可重複**：不依賴執行順序、不依賴外部網路、不依賴現在時間

| 常用的 testing 功能 | 用途 |
| --- | --- |
| `t.Helper()` | 標記輔助函式，失敗時回報呼叫端的行號 |
| `t.Cleanup(func())` | 註冊測試結束後的清理工作 |
| `t.TempDir()` | 建立測試用的暫存資料夾，結束後自動刪除 |
| `t.Parallel()` | 讓測試平行執行 |
| `go test -race` | 偵測並行的資料競爭（Ch 16） |

<!--
使用單元測試有兩個注意事項。

第一，測試「公開的行為」，也就是輸入和輸出，而不是內部怎麼實作的。這樣以後修改內部寫法時，只要行為不變，測試就不用改。

第二，測試要獨立、可重複。如果一個測試要連網路才能過、或是要看今天是星期幾，它就可能今天過、明天不過，這種「時好時壞」的測試最傷團隊的信任。

表格裡列出了一些常用的測試工具。t.TempDir 在第 12 章測試檔案操作時很好用；-race 是第 16 章並行性運算的重要工具。
-->

---
layout: default
---

# 綜合練習：密碼強度檢查
### 任務說明

1. 在套件 `password` 中寫函式 `Strength(pw string) (int, error)`：
   - 長度（以字元計算）< 8 回傳錯誤 `ErrTooShort`
   - 否則從 0 分開始：有小寫 +1、有大寫 +1、有數字 +1、有符號 +1、長度 ≥ 12 再 +1
2. 寫表格驅動測試 `TestStrength`，至少包含：太短、只有小寫、大小寫＋數字、全部條件、含中文的密碼
3. 使用 `errors.Is` 檢查錯誤案例
4. 執行 `go test -v -cover`，讓覆蓋率達到 100%

<!--
這個綜合練習要大家同時寫程式和測試。

建議的順序是：先寫測試，想清楚每種輸入期望的分數，再寫程式讓測試通過。這種「先寫測試」的開發方式叫做 TDD，測試驅動開發。

注意長度要用字元數計算，含中文的密碼要用 utf8.RuneCountInString，這是第 3 章的內容。
-->

---
zoom: 0.81
---

# 綜合練習：解題提示
### 提示說明

```go
// 檔名：password.go
package password

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var ErrTooShort = errors.New("密碼至少需要 8 個字元")

// classify 回傳密碼中是否包含小寫、大寫、數字、符號。
func classify(pw string) (lower, upper, digit, symbol bool) {
	for _, r := range pw {
		switch {
		case unicode.IsLower(r):
			lower = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsDigit(r):
			digit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			symbol = true
		}
	}
	return lower, upper, digit, symbol
}
```

<!--
先寫一個未匯出的輔助函式 classify，走訪每個字元，用 unicode 套件判斷種類，用具名回傳值記錄四個布林值。

注意中文字不是大寫也不是小寫，IsLower 和 IsUpper 都會回傳 false，所以中文不會加分，但會算進長度。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func Strength(pw string) (int, error) {
	n := utf8.RuneCountInString(pw) // 以字元計算長度
	if n < 8 {
		return 0, ErrTooShort
	}
	lower, upper, digit, symbol := classify(pw)
	score := 0
	for _, ok := range []bool{lower, upper, digit, symbol, n >= 12} {
		if ok {
			score++
		}
	}
	return score, nil
}
```

<!--
Strength 先用 RuneCountInString 計算字元數，不足 8 個就回傳哨兵錯誤。

接著呼叫 classify 取得四個布林值，再把五個條件放進一個布林切片，數有幾個是 true，就是分數。這個寫法比寫五個 if 簡潔。
-->

---

# 綜合練習：測試程式
### 提示說明

```go
// 檔名：password_test.go
package password

import (
	"errors"
	"testing"
)

var strengthTests = []struct {
	name    string
	pw      string
	want    int
	wantErr error
}{
	{"太短", "abc", 0, ErrTooShort},
	{"只有小寫", "abcdefgh", 1, nil},
	{"大小寫加數字", "Abcdefg1", 3, nil},
	{"全部條件", "Abcdefg1!xyz", 5, nil},
	{"含中文", "密碼abc123", 2, nil},
}
```

<!--
測試案例放在套件層級的變數 strengthTests，多了一個 wantErr 欄位，代表期望的錯誤。

「含中文」這個案例：「密碼abc123」共 8 個字元，剛好達到長度要求，有小寫和數字，所以是 2 分。
-->

---

# 綜合練習：測試程式（續）
### 提示說明

```go
// 續上頁
func TestStrength(t *testing.T) {
	for _, tt := range strengthTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Strength(tt.pw)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("錯誤 = %v，期望 %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Strength(%q) = %d，期望 %d",
					tt.pw, got, tt.want)
			}
		})
	}
}
```

<!--
errors.Is(err, nil) 在 err 也是 nil 的時候會回傳 true，所以正常的案例也能用同一個判斷。

錯誤不符合預期時用 t.Fatalf，因為錯誤都不對了，檢查分數也沒有意義，直接結束這個子測試。

執行 go test -v -cover，五個子測試都通過，覆蓋率 100%。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 9 步：測試與日誌

<!--
回到 GoShop。上一章我們做了一次大重構，把程式拆成五個套件。

問題來了：我們怎麼知道重構之後，結帳金額還是算對的？每次都手動執行、用眼睛看輸出，總有一天會漏看。今天學的單元測試，就是讓電腦幫我們檢查。
-->

---
class: code-sm
---

# GoShop 第 9 步：測試與日誌
### 任務說明

1. `internal/money/money_test.go`：用**表格驅動測試**檢查 `Money.String()`，**要包含負數**
2. `internal/checkout/checkout_test.go`：
   - `TestBest`：各種小計下，折扣規則是否挑對
   - `TestCheckout`：結帳成功時金額、編號、庫存都正確
   - `TestCheckoutErrors`：空購物車、商品不存在、庫存不足時，**庫存不能被扣**
3. 替 `Money.String()` 寫一個 benchmark（Go 1.24 的 `b.Loop`）
4. 在 `checkout` 用 `slog` 記錄「訂單成立」「訂單付款」；`main` 設定 `TextHandler`

```text
$ go test ./...
ok      goshop/internal/checkout    0.005s
--- FAIL: TestString/負數三位 (0.00s)
    money_test.go:21: Money(-300).String() = "NT$-,300"，want "-NT$300"
```

<!--
這一步要替 GoShop 寫測試。最重要的是結帳的測試，特別是「失敗的時候庫存不能被扣」，這種 bug 如果上線才發現，庫存資料就亂掉了。

另外請大家在 Money 的測試裡加上負數的案例，例如退款的時候金額可能是負的。

寫好之後執行 go test，你會發現……紅字了！負 300 元被印成了 NT$-,300。這是一個一直藏在我們程式裡的 bug，從第 5 章就在了，只是我們從來沒有用負數測試過。

這就是寫測試最大的價值：它會逼我們去想平常沒想到的情況。
-->

---
zoom: 0.91
---

# GoShop 第 9 步：解題提示
### 表格驅動測試

```go
// goshop/internal/money/money_test.go
func TestString(t *testing.T) {
	tests := []struct {
		name string
		in   Money
		want string
	}{
		{"零元", 0, "NT$0"},
		{"三位數", 999, "NT$999"},
		{"四位數", 1000, "NT$1,000"},
		// ...
		{"負數三位", -300, "-NT$300"},
		{"負數四位", -1500, "-NT$1,500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.want {
				t.Errorf("Money(%d).String() = %q，want %q",
					int(tt.in), got, tt.want)
			}
		})
	}
}
```

<!--
這是標準的表格驅動測試：先準備一張表，每一列是一個案例，有名稱、輸入和期望的結果，再用迴圈一個一個跑。

t.Run 會替每個案例建立子測試，失敗的時候訊息會顯示子測試的名稱，例如 TestString/負數三位，一眼就知道是哪個案例出問題。

錯誤訊息的格式建議寫成「函式呼叫 = 實際結果，want 期望結果」，這是 Go 社群的慣例，看的人不用翻程式碼就知道哪裡不對。
-->

---
zoom: 0.96
---

# GoShop 第 9 步：解題提示（續）
### 修正 bug：先處理負號

```go
// goshop/internal/money/money.go
// String 印出含千分位的金額，例如 NT$1,234、-NT$300。
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign, m = "-", -m // 先把負號拿掉，最後再補回去
	}
	s := strconv.Itoa(int(m))
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return sign + "NT$" + s
}
```

| 輸入 | 修正前 | 修正後 |
| --- | --- | --- |
| `-300` | `NT$-,300`（負號被當成一位數） | `-NT$300` |
| `-1500` | `NT$-1,500` | `-NT$1,500` |

<!--
bug 的原因是：負號也被算進字串長度了。-300 轉成字串是 4 個字元，程式以為它是四位數，就在負號後面插了一個逗號。

修正的方法很簡單：先把負號拿掉，用正數加千分位，最後再把負號補回去。

修好之後再執行一次 go test，全部通過。而且這個測試會一直留在專案裡，以後誰不小心改壞了，測試馬上就會失敗，這叫做「回歸測試」。
-->

---
zoom: 0.91
---

# GoShop 第 9 步：解題提示（續 2）
### 失敗時庫存不能被扣

```go
// goshop/internal/checkout/checkout_test.go
	tests := []struct {
		name  string
		items []Item
		want  error
	}{
		{"空的購物車", nil, shop.ErrEmptyCart},
		{"商品不存在", []Item{{"Z", 1}}, shop.ErrNotFound},
		{"庫存不足", []Item{{"A", 1}, {"B", 2}}, shop.ErrOutOfStock},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService()
			_, err := svc.Checkout(Cart{Items: tt.items})
			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v，want %v", err, tt.want)
			}
			// 結帳失敗時，庫存不能被扣掉
			if p, _ := svc.Store.Product("A"); p.Stock != 10 {
				t.Errorf("A 的庫存 = %d，want 10", p.Stock)
			}
		})
	}
```

<!--
這是結帳失敗的測試，一樣是表格驅動。

每個案例都建立一個全新的 Service，互不影響。newService 是我們寫的小工具函式，建立一個有 A、B 兩項商品的測試環境。

第三個案例最有意思：A 買 1 件沒問題，B 買 2 件庫存不足。我們不只檢查錯誤是 ErrOutOfStock，還檢查 A 的庫存是不是還是 10。如果程式邊檢查邊扣庫存，這個測試就會抓到。

errors.Is 能判斷 ErrOutOfStock，是因為 StockError 有一個 Unwrap 方法傳回 ErrOutOfStock，這是第 8 步拆套件時加上的。
-->

---

# GoShop 第 9 步：解題提示（續 3）
### benchmark 與結構化日誌

```go
// goshop/internal/money/money_test.go
func BenchmarkString(b *testing.B) {
	for b.Loop() {
		_ = Money(1234567).String()
	}
}
```

```go
// goshop/internal/checkout/checkout.go
	slog.Info("訂單成立", "id", o.ID, "total", o.Total)
```

```text
$ go test -bench . ./internal/money
BenchmarkString-4    7921296    144.0 ns/op
$ go run .
time=2026-09-24T07:10:57Z level=INFO msg=訂單成立 id=1 total=NT$1,880
time=2026-09-24T07:10:57Z level=INFO msg=訂單付款 id=1 method=貨到付款
```

<!--
benchmark 用 Go 1.24 的 b.Loop 寫法，迴圈會自動跑足夠多次，算出每次呼叫平均花多少時間。每次大約 144 奈秒，對一個字串處理函式來說已經很快了。

日誌的部分，我們在 checkout 裡用 slog.Info 記錄每一張成立的訂單。slog 是結構化日誌，每個欄位都是鍵值對，total 印出來是 NT$1,880，因為 slog 會呼叫 Money 的 String 方法。

main 一開始用 slog.SetDefault 設定 TextHandler 輸出到 stderr，這樣日誌和程式的正常輸出就分開了。
-->

---

# 章節總結

- **除錯原則**：重現 → 讀錯誤訊息 → 縮小範圍 → 一次改一處 → 寫成測試
- **fmt**：前綴決定輸出到哪（`S`、`F`）、後綴決定格式（`ln`、`f`）；除錯用 `%+v`、`%#v`、`%q`、`%T`
- **格式化數字**：`%8.2f` 寬度與精度；`%.0f` 採銀行家捨入；轉字串用 `strconv.FormatFloat`
- **log**：自動加時間；`log.Fatal` 會直接結束程式、**不執行 defer**；`log.New` 建立自訂 logger
- **slog**：結構化日誌（Go 1.21+）；`Debug`／`Info`／`Warn`／`Error` 四個等級；`With` 附加固定欄位
- **單元測試**：`xxx_test.go`、`func TestXxx(t *testing.T)`；**表格驅動測試** + `t.Run`；`-cover`、`-bench`
- **GoShop**：表格驅動測試抓到負數金額的 bug；`b.Loop` 效能測試；`slog` 記錄訂單與付款

下一章我們會介紹「時間處理」：`time` 套件的時間、格式化、時區與時間長度。

<!--
我們來整理今天學到的東西。

除錯有一套固定的步驟，最後一步一定是寫成測試。fmt 的函式有很好記的命名規則，除錯時用 %+v 和 %#v。浮點數格式化要注意銀行家捨入。日誌的部分，現代的 Go 程式建議用 slog 的結構化日誌。單元測試是 Go 內建的，表格驅動測試是主流寫法。

GoShop 這一步有了第一批單元測試，而且測試真的抓到了一個 bug：負數金額會印成 NT$-,300。修好之後，這個測試會一直守護著它，以後不管誰改了 String 方法，只要又壞掉，go test 馬上就會告訴我們。

從下一章開始，我們會學習 Go 標準函式庫裡最常用的幾個套件。第一個是 time 套件，時間處理幾乎是每個應用程式都會用到的功能，而且 Go 的時間格式化方式非常特別，一定會讓大家印象深刻。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的三個工具，會陪伴大家的整個 Go 開發生涯。

課後建議：回頭替前幾章的練習題補上表格驅動測試，例如第 2 章的 FizzBuzz、第 3 章的字串反轉。寫測試的過程中，很可能會發現當初沒想到的邊界案例。

有問題的同學現在可以提問！
-->

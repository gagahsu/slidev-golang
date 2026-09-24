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
title: 套件
routeAlias: ch08
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">套件</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「把程式碼分門別類，別人才用得到、自己才找得到」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第八章！

到上一章為止，我們已經學完了 Go 語言的核心語法。但是到目前為止，所有的程式碼都寫在同一個 main.go 裡。當程式越寫越大，幾千行、幾萬行都擠在同一個檔案，就像把整個家的東西都塞在一個房間裡，要找什麼都找不到。

今天要學的套件（package）和模組（module），就是 Go 用來組織程式碼的方式：怎麼把程式碼拆成好幾個套件、怎麼決定哪些東西要公開給別人用，以及怎麼下載、管理別人寫好的第三方套件。
-->

---
layout: default
---

# Outline

- **前言** — 何謂套件、運用套件的好處
- **使用套件** — 套件的命名、宣告、匯出規則
- **管理套件** — `GOROOT`、`GOPATH`、Go Modules、下載第三方模組
- **套件的呼叫與執行** — 套件別名、`init()` 函式、`internal` 目錄
- **GoShop 專案實作** — 第 8 步：把 GoShop 拆成多個套件
- **章節總結**

<!--
今天的內容分成四個部分。

首先認識什麼是套件、為什麼需要套件。接著學怎麼建立和使用自己的套件，最重要的是「匯出規則」：Go 用名稱的大小寫來決定公開或私有。

第三部分是管理套件，會從 Go 早期的 GOPATH 講到現在主流的 Go Modules，以及怎麼下載第三方模組。最後是一些細節：套件別名和 init 函式。
-->

---

# 回顧：介面

- 介面描述「**能做什麼**」；Go 的介面是**隱性實作**的
- 「**接受介面，回傳具體型別**」；`io.Reader` / `io.Writer` 是資料流的共通語言
- 回傳介面可以**隱藏實作**：`memStorage` 用小寫開頭，外部看不到
- 我們一直在用別人寫好的套件：`fmt`、`strings`、`errors`、`io`…

<!--
回顧一下上一章。

上一章有一個範例：memStorage 用小寫開頭，我們說「其他套件看不到它」。今天就會正式解釋，為什麼大小寫會決定「看不看得到」。

另外，從第一章開始，我們就一直在 import 別人寫好的套件：fmt、strings、errors。今天要反過來，學怎麼寫自己的套件，給別人 import。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## What Is a Package?

<!--
先來認識什麼是套件。
-->

---

# 何謂套件

「**套件（package）是同一個資料夾中、宣告了相同套件名稱的 `.go` 檔案集合。**」

```text
shop/                  ← 模組根目錄（有 go.mod）
├── go.mod
├── main.go            package main
└── pricing/           ← 一個資料夾 = 一個套件
    ├── pricing.go     package pricing
    └── discount.go    package pricing（同資料夾的檔案屬於同一個套件）
```

| 名詞 | 說明 |
| --- | --- |
| **套件（package）** | 一個資料夾，Go 程式碼組織與重複使用的**最小單位** |
| **模組（module）** | 一組一起發布、一起管理版本的套件，根目錄有 `go.mod` |

<!--
什麼是套件？

在 Go 裡，套件就是一個資料夾：同一個資料夾裡的所有 .go 檔，都屬於同一個套件，它們的第一行都要寫一樣的 package 名稱。同一個套件裡的檔案可以直接互相使用對方的函式和變數，就像在同一個檔案裡一樣。

模組則是比套件更大一層的單位：一個模組包含很多個套件，它們一起發布、一起管理版本。模組的根目錄有一個 go.mod 檔案，第零章我們用 go mod init 建立的就是模組。

用圖書館比喻：模組是整座圖書館，套件是一個一個的書架，.go 檔就是書架上的書。
-->

---

# 運用套件的好處

| 好處 | 說明 |
| --- | --- |
| **組織程式碼** | 相關的功能放在一起，找東西更快，例如 `pricing` 只處理價格 |
| **重複使用** | 寫一次，多個專案都能 import |
| **命名空間** | 不同套件可以有同名的函式：`pricing.Calc` 和 `shipping.Calc` 不會衝突 |
| **封裝** | 只公開必要的部分，內部實作可以隨時修改而不影響使用者 |
| **編譯速度** | Go 以套件為單位編譯，沒有修改的套件不需要重新編譯 |

<!--
使用套件有五個好處。

最重要的是「封裝」：套件可以只公開必要的函式，把內部的實作細節藏起來。就像餐廳的廚房：客人只看得到菜單（公開的函式），廚房裡怎麼做菜（內部實作）是看不到的。廚師想換一種做法，只要菜的味道一樣，客人完全不受影響。

另外 Go 以套件為單位編譯，而且有很好的快取機制，這也是 Go 編譯速度非常快的原因之一。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 使用套件
## Creating & Using Packages

<!--
接下來看怎麼建立和使用自己的套件。
-->

---

# 套件的宣告

每個 `.go` 檔的**第一行**（註解之後）都要宣告它屬於哪個套件：

```go
// Package pricing 負責商品價格的計算。
package pricing

const TaxRate = 0.05

func WithTax(price float64) float64 {
	return price * (1 + TaxRate)
}
```

| 規則 | 說明 |
| --- | --- |
| 同一個資料夾只能有**一個**套件名稱 | `_test.go` 測試檔可以例外使用 `pricing_test` |
| 套件名稱通常**等於資料夾名稱** | 資料夾 `pricing/` → `package pricing` |
| `package main` 是特例 | 代表可執行程式，必須有 `func main()` |

<!--
每一個 .go 檔的第一行，都要用 package 宣告它屬於哪個套件。

規則有三條：同一個資料夾只能有一個套件名稱；套件名稱通常跟資料夾名稱一樣；package main 是特例，代表這是一個可以執行的程式。

注意第一行上面的註解：「Package pricing 負責…」。這是 Go 的文件慣例，套件宣告上方的註解就是這個套件的說明文件，第 17 章的 go doc 工具會讀取它。
-->

---

# 套件的命名

| ✅ 好的套件名稱 | ❌ 避免的套件名稱 | 原因 |
| --- | --- | --- |
| `pricing` | `Pricing`、`price_calc` | 全部小寫、不用底線或大小寫混合 |
| `http`、`json`、`time` | `httputilities` | 簡短、名詞、清楚 |
| `user` | `util`、`common`、`helpers` | 太籠統，看不出放了什麼 |
| `strings` | `stringutils` | 標準函式庫的命名風格 |

```go
// 使用時會帶著套件名稱，所以函式名稱不需要重複套件名
pricing.WithTax(100)    // ✅ 讀起來很自然
pricing.PricingWithTax(100) // ❌ 重複（stutter）
```

<!--
套件的命名有幾個慣例。

第一，全部小寫，不要用底線、不要大小寫混合。第二，簡短、清楚，最好是一個名詞。第三，避免 util、common、helpers 這種籠統的名稱，因為看不出裡面放了什麼，最後常常變成什麼都丟進去的垃圾桶。

另外一個重要的觀念：使用套件時，前面一定會帶著套件名稱，例如 pricing.WithTax。所以函式名稱不需要再重複套件名稱，寫成 pricing.PricingWithTax 就很囉嗦，Go 社群稱這種情況為 stutter（口吃）。
-->

---

# 將套件的功能匯出

「**名稱以大寫字母開頭 → 匯出（exported），其他套件可以使用；小寫開頭 → 未匯出，只有同套件可以使用。**」

```go
package pricing

const TaxRate = 0.05      // ✅ 匯出：pricing.TaxRate
var defaultDiscount = 0.9 // 未匯出：只有 pricing 套件內可用

type Product struct {
	Name string  // ✅ 匯出的欄位
	cost float64 // 未匯出的欄位：外部看不到成本價
}

// ✅ 匯出
func WithTax(p float64) float64 { return round(p * (1 + TaxRate)) }

// 未匯出
func round(x float64) float64 { return float64(int(x*100+0.5)) / 100 }
```

<!--
這是 Go 最有特色的規則之一：用名稱的第一個字母大小寫，決定要不要公開。

大寫開頭的常數、變數、型別、函式、struct 欄位、方法，都會被「匯出」，其他套件可以使用；小寫開頭的就只有同一個套件裡看得到。

其他語言通常用 public、private 這些關鍵字，Go 直接用大小寫，一眼就能看出來，不需要額外的關鍵字。

struct 的欄位也一樣：Product 的 Name 是公開的，cost 是私有的，其他套件建立 Product 時無法設定或讀取成本價。這對第 11 章很重要：JSON 編碼只會處理「匯出」的欄位。
-->

---

# 在 main 套件中使用自己的套件

用 **「模組路徑 + 資料夾路徑」** 來 import：

```go
package main

import (
	"fmt"

	"example.com/shop/pricing" // go.mod 的 module 名稱 + 資料夾
)

func main() {
	fmt.Println(pricing.WithTax(100)) // 105
	fmt.Println(pricing.TaxRate)      // 0.05
	// pricing.round(1.5) // 編譯錯誤：undefined（未匯出）
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>import 分組慣例：</b> 標準函式庫一組、第三方與專案內的套件另一組，中間空一行。<code>goimports</code> / <code>gopls</code> 會自動幫我們整理。
</div>

<!--
在 main 套件裡使用自己的套件，import 的路徑是「go.mod 裡的模組名稱」加上「套件所在的資料夾路徑」。我們的模組叫 example.com/shop，pricing 資料夾在根目錄下，所以路徑是 example.com/shop/pricing。

import 之後，用「套件名稱.名稱」來使用。只能用匯出的名稱，小寫的 round 在外面是看不到的，會出現 undefined 的編譯錯誤。

另外注意 import 的分組：標準函式庫放一組，自己的套件和第三方套件放另一組，中間空一行。這是 Go 社群的慣例，VS Code 存檔時會自動整理。
-->

---
layout: default
---

# 練習 1：拆分計算機套件
### 任務說明

1. 建立模組 `go mod init example.com/calc`
2. 建立資料夾 `mathx/`，在 `mathx/mathx.go` 中宣告 `package mathx`：
   - 匯出函式 `Add(a, b int) int`、`Avg(nums ...int) float64`
   - 未匯出的輔助函式 `sum(nums []int) int`（`Avg` 會用到它）
3. 在根目錄的 `main.go` import `example.com/calc/mathx`，呼叫 `Add` 和 `Avg`
4. 試著在 `main.go` 呼叫 `mathx.sum`，觀察錯誤訊息

<!--
這個練習要大家自己動手建立一個多套件的專案。

重點是體會匯出規則：Add、Avg 是大寫，外面可以用；sum 是小寫，只有 mathx 套件自己可以用。試著在 main 裡呼叫 sum，看看編譯器怎麼說。
-->

---
zoom: 0.91
---

# 練習 1：解題提示
### 提示說明

```go
// mathx/mathx.go
package mathx

func Add(a, b int) int { return a + b }

func Avg(nums ...int) float64 {
	if len(nums) == 0 {
		return 0
	}
	return float64(sum(nums)) / float64(len(nums))
}

func sum(nums []int) int {
	t := 0
	for _, n := range nums {
		t += n
	}
	return t
}
```

```text
$ go run .
./main.go:9:20: undefined: mathx.sum   ← 呼叫未匯出函式時的錯誤
```

<!--
mathx.go 裡，Add 和 Avg 大寫開頭，會被匯出；sum 小寫，只能在 mathx 套件內使用。Avg 在套件內呼叫 sum 完全沒問題。

在 main.go 呼叫 mathx.sum 時，錯誤訊息是 undefined: mathx.sum，Go 會把未匯出的名稱當成「不存在」。

main.go 的寫法跟前一頁一樣：import "example.com/calc/mathx"，然後呼叫 mathx.Add(1, 2)、mathx.Avg(80, 90, 100)。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 管理套件
## GOROOT・GOPATH・Go Modules

<!--
接下來看 Go 怎麼管理套件：從早期的 GOPATH，到現在主流的 Go Modules。
-->

---

# GOROOT：Go 本身安裝的位置

`GOROOT` 是 Go 工具鏈和**標準函式庫**的位置

```bash
go env GOROOT
# /usr/local/go

ls $(go env GOROOT)/src
# bufio  bytes  context  crypto  encoding  errors
# fmt  io  net  os  strings  time ...
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>標準函式庫的原始碼都在這裡：</b> 想知道 <code>strings.ToUpper</code> 怎麼寫的？直接打開 <code>$GOROOT/src/strings/strings.go</code>，或在 VS Code 按 <code>F12</code> 跳到定義。<b>不需要、也不應該手動設定 GOROOT</b>。
</div>

<!--
GOROOT 是 Go 本身安裝的位置，裡面有編譯器，還有整個標準函式庫的原始碼。

我們 import "fmt" 的時候，Go 就是到 GOROOT/src/fmt 找這個套件。

一個很值得養成的習慣：看標準函式庫的原始碼。Go 的標準函式庫寫得非常好，是學習 Go 最好的範本。在 VS Code 裡對任何函式按 F12，就能直接跳到它的原始碼。

GOROOT 由安裝程式自動設定，現代的 Go 會自己找到它，不需要也不應該手動設定。
-->

---

# GOPATH：早期的工作區與現在的用途

| 時期 | GOPATH 的角色 |
| --- | --- |
| **Go 1.11 以前** | 所有專案**都必須**放在 `$GOPATH/src` 底下，所有專案共用同一份相依套件 |
| **現在（Go Modules）** | 專案可以放在**任何地方**；GOPATH 只剩下兩個用途 👇 |

| 現在的用途 | 路徑 |
| --- | --- |
| 下載的模組快取（唯讀） | `$GOPATH/pkg/mod` |
| `go install` 安裝的執行檔 | `$GOPATH/bin`（記得加進 `PATH`） |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>看到舊教學說「專案要放在 GOPATH/src」</b>，那是 2018 年以前的做法，現在已經不需要了。
</div>

<!--
GOPATH 是 Go 早期的「工作區」概念。

在 Go 1.11 之前，所有的 Go 專案都必須放在 GOPATH/src 資料夾底下，而且所有專案共用同一份相依套件。如果專案 A 需要某個套件的 1.0 版，專案 B 需要 2.0 版，就會打架。這是 Go 早期最被詬病的問題。

Go Modules 出現之後，專案可以放在電腦上的任何地方，每個專案各自記錄需要的版本。GOPATH 現在只剩兩個用途：存放下載的模組快取，以及 go install 安裝的執行檔。

網路上很多舊教學還在教 GOPATH 的做法，大家看到的時候知道那是舊時代的產物就好。
-->

---

# Go Modules：現代的套件管理方式

Go 1.11 推出、**Go 1.16 起成為預設**，以 `go.mod` 和 `go.sum` 管理相依模組

```text
module example.com/shop          ← 模組路徑：其他人 import 時使用

go 1.27.0                         ← 最低需要的 Go 版本

require (
	github.com/google/uuid v1.6.0 ← 直接相依的模組與版本
	github.com/go-sql-driver/mysql v1.10.1
)
```

| 檔案 | 用途 | 要不要 commit |
| --- | --- | --- |
| `go.mod` | 記錄模組路徑、Go 版本、相依模組與版本 | ✅ 要 |
| `go.sum` | 記錄每個模組版本的**雜湊值**，防止被竄改 | ✅ 要 |

<!--
Go Modules 是現在唯一主流的套件管理方式，從 Go 1.16 開始成為預設。

go.mod 就像是專案的「購物清單」：記錄這個模組叫什麼名字、需要哪個版本以上的 Go，以及需要哪些第三方模組、各是哪個版本。

go.sum 則像「收據」：記錄每個模組版本內容的雜湊值。下次下載的時候，Go 會比對雜湊值，如果內容被竄改過，就會拒絕使用。雜湊函式是第 18 章的內容，這裡先知道它是用來「驗證內容沒有被改過」就好。

這兩個檔案都要 commit 到 Git，這樣團隊裡每個人、每台 CI 機器，下載到的都是一模一樣的版本。
-->

---

# 常用的模組指令

| 指令 | 說明 |
| --- | --- |
| `go mod init <模組路徑>` | 建立新模組，產生 `go.mod` |
| `go get <模組>@<版本>` | 新增或升級相依模組，例如 `@latest`、`@v1.6.0` |
| `go get <模組>@none` | 移除相依模組 |
| `go mod tidy` | **自動整理**：補上缺少的、刪除沒用到的相依模組 |
| `go list -m all` | 列出所有相依模組（包含間接相依） |
| `go mod why <模組>` | 查詢為什麼需要這個模組 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>最常用的組合：</b> 在程式碼中寫好 <code>import</code> → 執行 <code>go mod tidy</code>，Go 會自動下載並更新 <code>go.mod</code> / <code>go.sum</code>。
</div>

<!--
這張表列出了最常用的模組指令。

實務上最常用的是 go mod tidy：我們先在程式碼裡寫好 import，然後執行 go mod tidy，Go 會自動分析程式碼用到了哪些模組，缺少的就下載、沒用到的就從 go.mod 裡刪掉。可以把它想成「自動整理購物清單」。

go get 則用來明確指定要新增或升級某個模組的版本，@latest 是最新版，也可以指定確切的版本號。
-->

---

# 下載第三方模組或套件

以產生 UUID 的 `github.com/google/uuid` 為例：

```bash
go get github.com/google/uuid@latest
# go: downloading github.com/google/uuid v1.6.0
# go: added github.com/google/uuid v1.6.0
```

```go
package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	// 例如 d6486408-e48f-41e4-aae2-9dee3d8da9c7
	fmt.Println("訂單編號：", uuid.NewString())
}
```

<!--
使用第三方模組的流程是：go get 下載，程式碼裡 import，然後就能使用了。

這裡以 Google 的 uuid 模組為例，UUID 是一種全世界不會重複的識別碼，常用來當訂單編號、資料的主鍵。

go get 執行後，go.mod 會多一行 require github.com/google/uuid v1.6.0，go.sum 也會多兩行雜湊值。

注意 import 的分組：fmt 是標準函式庫，一組；uuid 是第三方，另一組。
-->

---

# 如何挑選第三方模組：pkg.go.dev

到 [pkg.go.dev](https://pkg.go.dev) 搜尋模組，檢查以下資訊：

| 檢查項目 | 說明 |
| --- | --- |
| **Imported by** | 有多少其他模組使用它，越多代表越多人驗證過 |
| **版本與發布日期** | 是否有持續維護；`v1` 以上代表 API 穩定 |
| **License** | 授權條款是否允許商業使用（MIT、Apache-2.0、BSD 較寬鬆） |
| **文件與範例** | 有沒有清楚的說明和 Example |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 的文化：</b> 先看標準函式庫能不能做到。Go 的標準函式庫非常完整，HTTP 伺服器、JSON、加密、測試都內建，很多事情<b>不需要第三方套件</b>。
</div>

<!--
第三方模組要去哪裡找？官方的 pkg.go.dev 網站可以搜尋所有公開的 Go 模組，還會自動產生文件。

挑選的時候，建議檢查四件事：有多少人在用、有沒有持續維護、授權條款、文件是否清楚。

另外 Go 社群有一個文化：能用標準函式庫就用標準函式庫。Go 的標準函式庫非常完整，這門課後面的 JSON、檔案、HTTP 客戶端、HTTP 伺服器、加密，全部都只用標準函式庫就能完成，唯一的例外是第 13 章的 MySQL 驅動程式。
-->

---

# 補充：語意化版本與主版本號

| 版本 | 意義 |
| --- | --- |
| `v1.6.0` → `v1.6.1` | **Patch**：修正 bug，完全相容 |
| `v1.6.0` → `v1.7.0` | **Minor**：新增功能，向下相容 |
| `v1.x.x` → `v2.0.0` | **Major**：有不相容的修改 |

```go
import "github.com/go-chi/chi/v5" // v2 以上的主版本，模組路徑要加上 /vN
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>匯入相容規則：</b> Go 規定「相同的 import 路徑必須向下相容」，所以不相容的 v2 以上要換一個路徑（<code>/v2</code>、<code>/v5</code>），兩個主版本可以同時存在於同一個專案。
</div>

<!--
這是補充內容：版本號的意義。

Go 模組使用「語意化版本」，版本號分成三段：主版本、次版本、修訂版本。修訂版本是修 bug，次版本是加新功能，這兩種都保證相容；主版本改變，代表有不相容的修改。

Go 有一個獨特的規則：如果主版本號是 2 以上，import 路徑要加上 /v2、/v5 這樣的後綴。這樣 v1 和 v2 就是兩個不同的路徑，可以同時存在，升級的時候可以一個套件一個套件慢慢換。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 套件的呼叫與執行
## Aliases・init()・internal

<!--
最後來看套件使用上的一些細節。
-->

---
zoom: 0.95
---

# 套件別名

import 時可以替套件取一個別名：`import 別名 "路徑"`

```go
package main

import (
	crand "crypto/rand" // 別名：避免和 math/rand 衝突
	"fmt"
	"math/rand/v2"
)

func main() {
	fmt.Println(rand.IntN(100)) // math/rand/v2：一般用途的亂數
	// crypto/rand：密碼學安全的隨機字串（Go 1.24+）
	fmt.Println(crand.Text())
}
```

| 特殊寫法 | 意義 | 使用時機 |
| --- | --- | --- |
| `import _ "路徑"` | **空白匯入**：只執行套件的 `init()`，不使用它的名稱 | 註冊資料庫驅動程式（Ch 13） |
| `import . "路徑"` | 點匯入：不需要寫套件名稱 | **不建議使用**，會讓人看不出名稱來自哪裡 |

<!--
import 的時候可以替套件取別名，最常見的用途是解決名稱衝突。

crypto/rand 和 math/rand/v2 兩個套件的名稱都叫 rand，同時 import 會衝突。所以我們把 crypto/rand 取別名叫 crand。math/rand/v2 是 Go 1.22 推出的新版亂數套件，用法比舊版更簡單。crypto/rand 產生的是密碼學安全的亂數，第 18 章會用到。

表格裡有兩種特殊寫法。底線加路徑叫做「空白匯入」，意思是「我不會直接使用這個套件，只是要它的 init 函式執行」，第 13 章註冊 MySQL 驅動程式就會用到。點匯入可以省略套件名稱，但會讓程式碼難以閱讀，不建議使用。
-->

---

# init() 函式

每個套件可以有 `init()` 函式，在 **`main()` 之前自動執行**，用來做初始化

```go
package main

import "fmt"

var config = loadConfig() // 1. 套件層級變數先初始化

func loadConfig() map[string]string {
	fmt.Println("1. 初始化套件層級變數")
	return map[string]string{"env": "dev"}
}

func init() { // 2. 接著執行 init()（一個套件可以有多個）
	fmt.Println("2. init()：檢查設定", config["env"])
}

func main() { // 3. 最後才是 main()
	fmt.Println("3. main() 開始")
}
```

<!--
init 是一個特殊的函式：不能被呼叫、沒有參數也沒有回傳值，Go 會在程式啟動時自動執行它。

執行順序是：先初始化套件層級的變數，再執行 init，最後才執行 main。如果 main 套件 import 了其他套件，會先把被 import 的套件都初始化完（它們的變數和 init），才輪到 main 套件。

一個套件、甚至一個檔案裡，都可以有多個 init，它們會依照出現的順序執行。

執行後會依序印出 1、2、3。
-->

---

# 使用 init() 的注意事項

**注意事項之一：** 初始化順序是 **被 import 的套件 → 套件層級變數 → `init()` → `main()`**

**注意事項之二：** `init()` **不要做太多事**，它無法回傳錯誤、也讓程式行為變得不透明

| ✅ 適合放在 init | ❌ 不適合放在 init |
| --- | --- |
| 註冊驅動程式、編解碼器 | 連線資料庫、呼叫網路 API |
| 驗證寫死的設定、預先計算查詢表 | 讀取可能不存在的檔案 |
| 設定套件內的預設值 | 需要處理錯誤的任何工作 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>建議：</b> 需要處理錯誤的初始化，寫成一般函式（例如 <code>NewServer() (*Server, error)</code>）並在 <code>main</code> 中明確呼叫。
</div>

<!--
使用 init 有兩個注意事項。

第一，要記住初始化的順序：被 import 的套件先初始化，然後是自己套件的變數、init，最後是 main。

第二，init 不要做太多事。init 沒辦法回傳錯誤，失敗了只能 panic；而且它是自動執行的，讀程式碼的人不一定會注意到，程式的行為會變得不透明。

所以像連線資料庫這種可能失敗的工作，應該寫成一般的函式，回傳 error，在 main 裡面明確呼叫。init 只適合做「一定會成功」的簡單工作，例如註冊驅動程式。
-->

---

# 補充：internal 目錄與專案結構

放在 **`internal/`** 底下的套件，**只有同一個模組（父目錄）內**的程式碼可以 import

```text
shop/
├── go.mod                      module example.com/shop
├── cmd/
│   └── shop/main.go            ← 可執行程式的進入點（package main）
├── internal/
│   ├── pricing/pricing.go      ← 只有 shop 模組內可以 import
│   └── storage/storage.go
└── pkg/                        ← （選用）打算給外部使用的套件
```

```text
use of internal package example.com/shop/internal/pricing not allowed
```

<!--
這是補充內容：internal 目錄。

Go 有一個特殊規則：放在 internal 資料夾底下的套件，只有 internal 的父目錄底下的程式碼可以 import。其他模組想 import 的話，會出現「use of internal package not allowed」的編譯錯誤。

這是比大小寫更大一層的封裝：大小寫控制「套件外能不能看到」，internal 控制「模組外能不能看到」。

圖上是一個常見的 Go 專案結構：cmd 放可執行程式的進入點，internal 放專案內部的套件。不過 Go 官方並沒有規定一定要這樣排，小專案把所有檔案放在根目錄也完全沒問題，等專案長大了再拆分。
-->

---
layout: default
---

# 綜合練習：訂單系統模組
### 任務說明

建立模組 `example.com/orders`，結構如下：

```text
orders/
├── go.mod
├── cmd/orders/main.go
└── internal/
    └── order/order.go
```

1. `order` 套件：匯出 `type Order struct { ID string; Amount int }` 與 `func New(amount int) (*Order, error)`
2. `New` 使用第三方模組 `github.com/google/uuid` 產生 `ID`；金額 ≤ 0 時回傳錯誤
3. `order` 套件有一個 `init()`，印出「order 套件已載入」
4. `main.go` 用別名 `o` import `order` 套件，建立兩筆訂單（一筆金額為 0），印出結果
5. 用 `go mod tidy` 整理相依模組，觀察 `go.mod` 的變化

<!--
這個綜合練習把今天的內容全部串起來：多套件的專案結構、internal 目錄、匯出規則、第三方模組、init 函式、套件別名、go mod tidy。

建議大家真的在自己的電腦上動手做一次，觀察 go.mod 和 go.sum 的變化。
-->

---
zoom: 0.91
---

# 綜合練習：解題提示
### 提示說明

```go
// internal/order/order.go
package order

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Order struct {
	ID     string
	Amount int
}

func init() { fmt.Println("order 套件已載入") }

func New(amount int) (*Order, error) {
	if amount <= 0 {
		return nil, errors.New("金額必須大於 0")
	}
	return &Order{ID: uuid.NewString(), Amount: amount}, nil
}
```

<!--
order.go 宣告 package order，匯出 Order 型別和 New 函式。New 是 Go 慣用的「建構函式」命名：套件名稱加上 New，使用時寫成 order.New，讀起來就是「新的訂單」。

回傳 *Order 指標和 error，金額不合法時回傳 nil 和錯誤。
-->

---
zoom: 0.94
---

# 綜合練習：解題提示（續）
### 提示說明

```go
// cmd/orders/main.go
package main

import (
	"fmt"

	o "example.com/orders/internal/order" // 套件別名
)

func main() {
	for _, amt := range []int{350, 0} {
		ord, err := o.New(amt)
		if err != nil {
			fmt.Println("建立失敗：", err)
			continue
		}
		fmt.Printf("%+v\n", *ord)
	}
}
```

```bash
go mod tidy && go run ./cmd/orders
```

<!--
main.go 用別名 o 匯入 order 套件，所以呼叫時寫 o.New。

執行 go mod tidy，Go 會發現程式碼用到了 uuid 模組，自動下載並加到 go.mod。然後 go run ./cmd/orders 執行，注意路徑要指向 main 套件所在的資料夾。

執行後會先印出「order 套件已載入」（init 在 main 之前執行），然後印出第一筆訂單，第二筆印出建立失敗。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 8 步：拆成多個套件

<!--
回到 GoShop。經過七章的累積，main.go 已經快 300 行了：金額、商品、購物車、折扣、錯誤、訂單、付款，全部擠在同一個檔案裡。

找一個函式要捲很久，改一個地方還要擔心會不會影響到別的東西。這就是今天學套件的最好時機。
-->

---
zoom: 0.98
---

# GoShop 第 8 步：拆成多個套件
### 任務說明

依照「**一個套件只負責一件事**」的原則，把程式碼搬到 `internal/` 底下：

| 套件 | 負責的事 | 主要內容 |
| --- | --- | --- |
| `money` | 金額 | `Money`、`String()`、`Times()` |
| `shop` | 核心資料與錯誤 | `Product`、`Order`、`Status`、`ErrNotFound`、`StockError` |
| `store` | 保存資料 | `Memory`：商品、訂單、檢查並扣庫存 |
| `checkout` | 結帳流程 | `Cart`、`Discount`、`Service.Checkout`、`Service.Pay` |
| `payment` | 付款方式 | `Method` 介面、`CreditCard`、`Wallet`、`CashOnDelivery` |

- 訂單加上狀態 `Status`（`Pending`／`Paid`，用 `iota`），付款失敗時訂單保留為**待付款**
- `main.go` 只負責建立物件、把它們組起來

<!--
這是拆分的規劃。五個套件各自負責一件事，名稱都很短，而且是單數名詞，這是 Go 的套件命名慣例。

放在 internal 資料夾底下，代表這些套件只給 goshop 自己用，別的模組不能 import。這是今天補充介紹的 internal 目錄。

趁著重構，我們也順便調整了一點設計：訂單多了一個狀態欄位，結帳和付款拆成兩個步驟。這樣付款失敗的時候，訂單還在，只是狀態是「待付款」，顧客之後可以換一種方式再付一次，這也比較接近真實的電商。
-->

---

# GoShop 第 8 步：專案結構與相依關係

```text
goshop/
├── go.mod                  module goshop
├── main.go                 package main：組裝、執行
└── internal/
    ├── money/money.go      package money
    ├── shop/shop.go        package shop      → money
    ├── store/memory.go     package store     → shop
    ├── checkout/           package checkout  → shop、store、payment
    │   ├── cart.go
    │   ├── discount.go
    │   └── checkout.go
    └── payment/payment.go  package payment   → money
```

- 箭頭是 import 的方向；Go **不允許套件互相 import**（循環相依）
- 同一個資料夾的檔案屬於同一個套件，可以直接使用彼此**未匯出**的名稱

<!--
這是拆完之後的樣子。右邊的箭頭表示這個套件 import 了誰。

大家看一下方向：money 最底層，誰都可以用它；shop 用 money；store 用 shop；checkout 在最上面，用到 shop、store、payment。整個相依關係是一個由上往下的樹，沒有繞圈圈。

這很重要，因為 Go 不允許循環相依：如果 shop import 了 store，store 又 import shop，編譯就會失敗。規劃套件的時候，先想清楚誰依賴誰，就不會踩到這個坑。

checkout 套件有三個檔案，它們的第一行都是 package checkout，屬於同一個套件。
-->

---
zoom: 0.91
---

# GoShop 第 8 步：解題提示
### 匯出的名稱：大寫開頭

```go
// goshop/internal/money/money.go
// Package money 處理新台幣金額。
package money

import "strconv"

// Money 是新台幣金額，單位是「元」。
type Money int

// Times 傳回單價乘上數量的金額。
func (m Money) Times(qty int) Money {
	return m * Money(qty)
}
```

```go
// goshop/internal/shop/shop.go
// Product 是一項商品。
type Product struct {
	SKU   string
	Name  string
	Price money.Money
	Stock int
}
```

<!--
先看最底層的 money 套件。第一行的註解以「Package money」開頭，這是套件的說明文件，第 17 章的 go doc 會讀取它。

Money、Times、String 都是大寫開頭，所以其他套件可以使用。

在 shop 套件裡，就要寫 money.Money 來使用它，也就是「套件名稱.名稱」。Product 的欄位也全部大寫開頭，其他套件才能讀寫這些欄位。
-->

---
zoom: 0.97
---

# GoShop 第 8 步：解題提示（續）
### store：保存商品與訂單

```go
// goshop/internal/store/memory.go
// Memory 把資料存在記憶體中，程式結束就會消失。
type Memory struct {
	products map[string]shop.Product
	orders   map[int]shop.Order
	lastID   int
}

// NewMemory 建立一個記憶體儲存庫，並放入初始商品。
func NewMemory(products ...shop.Product) *Memory {
	m := &Memory{
		products: make(map[string]shop.Product),
		orders:   make(map[int]shop.Order),
	}
	for _, p := range products {
		m.products[p.SKU] = p
	}
	return m
}
```

- 欄位都是**小寫**：外部只能透過 `Product()`、`PlaceOrder()` 等方法存取

<!--
store 套件的 Memory 結構，欄位都是小寫開頭，其他套件看不到，只能透過它提供的方法操作。這就是封裝：資料要怎麼存是 store 自己的事，別人不需要知道，也不能亂改。

建立 Memory 要用 NewMemory 這個函式，它會把 map 初始化好。如果讓外面直接寫 store.Memory{}，map 會是 nil，寫入的時候就會 panic。提供 New 開頭的建構函式，是 Go 很常見的慣例。

PlaceOrder 的內容就是第 6 章的 validate 加上扣庫存，這裡就不重複列出了。
-->

---

# GoShop 第 8 步：解題提示（續 2）
### checkout：使用其他套件

```go
// goshop/internal/checkout/checkout.go
import (
	"errors"
	"fmt"

	"goshop/internal/payment"
	"goshop/internal/shop"
	"goshop/internal/store"
)

// Service 負責結帳與付款流程。
type Service struct {
	Store *store.Memory
	Rules []Discount
}
```

- 自己模組的套件用「模組名稱 + 路徑」import：`goshop/internal/store`
- 標準函式庫和自己的套件之間空一行，是 `goimports` 的排版慣例

<!--
checkout 套件的 import 區塊，上面是標準函式庫，下面是我們自己的套件。自己的套件路徑是 go.mod 裡的模組名稱 goshop，加上資料夾路徑。

Service 把結帳需要的東西都放在欄位裡：用哪個 Store、有哪些折扣規則。Checkout 和 Pay 就是它的兩個方法。

這裡的 Store 欄位型別是 *store.Memory，也就是寫死了一定要用記憶體版本。第 13 章換成 MySQL 的時候，我們會把它改成介面，到時候就會看到介面的好處。
-->

---
zoom: 0.88
---

# GoShop 第 8 步：解題提示（續 3）
### main：把零件組起來

```go
// goshop/main.go
	svc := &checkout.Service{
		Store: st,
		Rules: []checkout.Discount{
			checkout.PercentOff(10), checkout.Threshold(2000, 300)},
	}
	wallet := &payment.Wallet{Balance: 1000}

	var cart checkout.Cart
	cart.Add(checkout.Item{SKU: "SKU-001", Qty: 2},
		checkout.Item{SKU: "SKU-003", Qty: 1})

	o, err := svc.Checkout(cart)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := svc.Pay(o.ID, wallet); err != nil {
		fmt.Println(err) // 錢包餘額不足，訂單保留為待付款
	}
```

```text
訂單 #1 付款失敗：GoShop 錢包：餘額不足，只剩 NT$1,000
===== 訂單 #1（已付款）=====   ← 之後改用貨到付款成功
```

<!--
main 現在只做「組裝」的工作：建立 store、建立結帳服務、建立錢包，然後呼叫它們。

因為 main 在不同的套件，所有東西前面都要加上套件名稱，例如 checkout.Item、payment.Wallet。從別的套件使用結構的時候，Go 建議寫出欄位名稱，所以這裡寫 SKU: "SKU-001"，而不是只寫值。

執行後，錢包 1000 元不夠付 1880 元，訂單保留為待付款；接著改用貨到付款，訂單就變成已付款了。

最後提醒大家，執行多套件的專案要用 go run 點，而不是 go run main.go，這樣 Go 才會編譯整個模組。
-->

---

# 章節總結

- **套件**：一個資料夾就是一個套件；**模組**是一起管理版本的套件集合，根目錄有 `go.mod`
- **命名**：小寫、簡短、不用底線；避免 `util`、`common`；不要和套件名稱重複（stutter）
- **匯出規則**：**大寫開頭匯出、小寫開頭私有**，適用常數、變數、型別、函式、欄位、方法
- **GOROOT / GOPATH**：前者是 Go 安裝位置；後者現在只存模組快取與 `go install` 執行檔
- **Go Modules**：`go.mod` + `go.sum` 都要 commit；常用 `go get`、`go mod tidy`
- **import**：別名解決衝突；`_` 空白匯入只執行 `init()`；避免 `.` 點匯入
- **init()**：在 `main()` 前自動執行；只做一定會成功的簡單工作；`internal/` 限制模組外匯入
- **GoShop**：把 300 行的 `main.go` 拆成 `money`、`shop`、`store`、`checkout`、`payment` 五個 `internal` 套件

下一章我們會介紹「程式除錯」：格式化輸出、日誌與單元測試。

<!--
我們來整理今天學到的東西。

套件是 Go 組織程式碼的基本單位，一個資料夾就是一個套件。大寫開頭匯出、小寫開頭私有，這是 Go 最重要的封裝規則。Go Modules 是現代的套件管理方式，go.mod 和 go.sum 都要 commit，最常用的指令是 go mod tidy。

GoShop 在這一步完成了第一次大重構：300 行的 main.go 拆成五個 internal 套件，每個套件只負責一件事，main 只剩下「把零件組起來」的工作。之後每一章加功能，都只要動到相關的那一兩個套件。

到這裡，我們已經學完了 Go 語言的基礎，也知道怎麼組織一個專案。接下來的章節要進入實戰：下一章先學怎麼除錯，包括用 fmt 做格式化輸出、用 log 和 slog 記錄日誌，以及用 testing 套件寫單元測試，確保程式的正確性。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
套件和模組是 Go 專案的骨架，今天的內容大家一定要實際動手做一次，光看投影片很難體會。

課後建議：到 pkg.go.dev 搜尋幾個熱門的模組，例如 gin、cobra、zap，看看它們的 Imported by 數量、授權條款和文件，練習評估一個第三方模組。

有問題的同學現在可以提問！
-->

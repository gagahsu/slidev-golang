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
title: 運用 Go 語言工具
routeAlias: ch17
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">運用 Go 語言工具</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「一個 go 指令，從編譯到發布一手包辦」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十七章！

第 0 章我們說過，Go 最大的特色之一是「一個 go 指令做完所有事」。這一路上我們已經用過 go run、go build、go test、go mod tidy。今天要把 Go 的工具鏈完整地認識一遍：怎麼編譯出給不同作業系統用的執行檔、怎麼讓程式碼自動排版、怎麼找出潛在的 bug、怎麼查詢文件、怎麼管理工具。

這一章比較輕鬆，大多是指令的操作，但這些工具是 Go 開發體驗這麼好的原因，熟悉它們會讓大家的開發效率大幅提升。
-->

---
layout: default
---

# Outline

- **go build** — 編譯執行檔、編譯條件、跨平台編譯
- **go run** — 執行程式
- **gofmt 與 gopls** — 程式碼格式化、語言伺服器
- **go vet 與靜態分析** — 找出可疑的程式碼；staticcheck、golangci-lint；`go fix`
- **go doc** — 查詢與產生文件
- **go get / go install / go tool** — 下載模組與工具
- **章節總結**

<!--
今天的內容依照開發流程排列：寫程式時用 gofmt 和 gopls，寫完用 go vet 檢查，測試通過後用 go build 編譯，需要查資料時用 go doc，需要第三方工具時用 go get 和 go install。
-->

---

# 回顧：並行性運算

- `go f()` 啟動 goroutine；`sync.WaitGroup`（Go 1.25+ `wg.Go`）等待完成
- 資料競爭：用 **`go run -race`** / **`go test -race`** 偵測 ← 今天會看到更多工具旗標
- 通道：由送出方關閉；`select` 等待多個通道
- `context`：`WithTimeout` + `defer cancel()`；`go vet` 會檢查**忘記呼叫 cancel** ← 今天介紹 vet

<!--
回顧一下上一章。

我們學了 goroutine、通道和 context，也用過 -race 這個旗標偵測資料競爭。上一章也提到：go vet 會檢查 context 有沒有忘記呼叫 cancel、Mutex 有沒有被複製。今天就來正式認識 go vet 和其他工具。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# go build 工具：編譯可執行檔
## go build

<!--
先從最常用的 go build 開始。
-->

---

# 使用 go build

| 指令 | 說明 |
| --- | --- |
| `go build` | 編譯目前資料夾的套件；`main` 套件會產生執行檔 |
| `go build -o bin/app .` | 指定輸出的檔名與位置 |
| `go build ./...` | 編譯所有套件（檢查能不能編譯，不產生執行檔） |
| `go build -race` | 加入競爭偵測（Ch 16） |
| `go build -ldflags "-s -w"` | 移除除錯資訊，**縮小執行檔** |
| `go build -trimpath` | 移除執行檔中的本機路徑（發布時建議使用） |

```bash
go build -o bin/goshop ./cmd/goshop
# 約 8 MB，靜態連結，複製到同平台的電腦就能執行
ls -lh bin/goshop
go version -m bin/goshop     # 查看執行檔的 Go 版本與相依模組
```

<!--
go build 把 Go 程式編譯成執行檔。

-o 指定輸出的位置；./... 代表所有子資料夾的套件，常用在 CI 裡檢查整個專案能不能編譯。

-ldflags "-s -w" 會移除除錯用的符號表，執行檔可以縮小大約三成。-trimpath 會移除編譯時的本機路徑，發布給別人的執行檔建議加上，避免洩漏電腦的資料夾結構。

go version -m 是一個很實用的指令：它可以讀出任何 Go 執行檔是用哪個版本的 Go 編譯的、用了哪些模組的哪個版本。當線上的程式出問題時，可以用它確認執行的到底是哪一版。
-->

---

# 用 -ldflags 在編譯時注入版本資訊

```go
package main

import "fmt"

var (
	version = "dev" // 編譯時可以被 -X 覆蓋
	commit  = "none"
)

func main() {
	fmt.Printf("goshop %s (%s)\n", version, commit)
}
```

```bash
go run .
# goshop dev (none)

commit=$(git rev-parse --short HEAD)
go build -ldflags "-X main.version=1.2.3 -X main.commit=$commit" .
./goshop
# goshop 1.2.3 (a1b2c3d)
```

<!--
發布程式的時候，我們常常希望執行檔知道自己是哪一個版本，例如 goshop --version 印出版本號。

版本號不應該寫死在程式碼裡，否則每次發布都要改程式碼。Go 的做法是：在程式碼裡宣告一個套件層級的字串變數，編譯時用 -ldflags "-X 套件.變數=值" 把值注入進去。

這裡用 git rev-parse 取得目前的 commit 編號，一起注入。CI/CD 自動發布時，通常就是這樣把版本號和 commit 寫進執行檔。

注意：-X 只能設定字串型別的套件層級變數。
-->

---

# 編譯條件：選擇要編譯的檔案

**方法一：檔名後綴** — `_作業系統`、`_架構`、`_作業系統_架構`

```text
notify.go            ← 所有平台都會編譯
notify_windows.go    ← 只在 Windows 編譯
notify_darwin.go     ← 只在 macOS 編譯
notify_linux_arm64.go
```

**方法二：`//go:build` 建置限制**（寫在檔案第一行，`package` 之前）

```go
//go:build linux || darwin

package notify
```

```go
// 自訂標籤：go build -tags prod 時不編譯這個檔案
//go:build !prod

package config
```

<!--
有些程式碼只適用於特定的平台，例如 Windows 和 macOS 發送桌面通知的方式完全不同。Go 提供兩種方式選擇要編譯哪些檔案。

第一種是檔名後綴：檔名以 _windows.go 結尾，就只會在編譯 Windows 版本時被編進去。

第二種是 //go:build 建置限制，寫在檔案的第一行，後面空一行才是 package。它可以寫布林運算式：|| 是或、&& 是且、! 是非。除了作業系統，還可以用自訂的標籤，搭配 go build -tags 使用，例如區分開發版和正式版的設定。

補充：網路上舊的寫法是 // +build，從 Go 1.17 開始改成 //go:build，gofmt 會自動幫我們轉換。
-->

---

# 如何針對跨平台編譯

設定 `GOOS`（作業系統）與 `GOARCH`（CPU 架構）環境變數，**在任何電腦上都能編譯出其他平台的執行檔**

| 目標平台 | 指令 |
| --- | --- |
| Windows 64 位元 | `GOOS=windows GOARCH=amd64 go build -o app.exe .` |
| macOS Apple Silicon | `GOOS=darwin GOARCH=arm64 go build -o app-mac .` |
| Linux 伺服器 | `GOOS=linux GOARCH=amd64 go build -o app-linux .` |
| 樹莓派 | `GOOS=linux GOARCH=arm64 go build -o app-pi .` |

```bash
go tool dist list | head -3   # 列出所有支援的平台（約 47 種）
# aix/ppc64
# android/386
# android/amd64
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Windows PowerShell：</b> <code>$env:GOOS="linux"; $env:GOARCH="amd64"; go build .</code>；使用到 cgo 的程式需要對應平台的 C 編譯器，純 Go 程式可設定 <code>CGO_ENABLED=0</code>。
</div>

<!--
這是 Go 最讓人驚豔的功能之一：跨平台編譯。

只要設定兩個環境變數：GOOS 是目標作業系統，GOARCH 是目標 CPU 架構，就能在自己的電腦上，編譯出其他平台的執行檔。在 Mac 上編譯 Windows 的 exe、在 Windows 上編譯 Linux 伺服器用的程式，完全不需要安裝其他東西。

go tool dist list 會列出所有支援的平台，大約 47 種組合，從手機、伺服器到嵌入式裝置都有。

在 Windows 的 PowerShell 裡，設定環境變數的語法不一樣，寫法在下方的提示裡。另外，如果程式用到了 C 語言的程式碼（cgo），跨平台編譯會比較麻煩；純 Go 的程式設定 CGO_ENABLED=0，就能編譯出完全靜態、不依賴任何系統函式庫的執行檔，最適合放進 Docker。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# go run 工具：執行程式
## go run

<!--
接下來是我們從第 0 章就一直在用的 go run。
-->

---

# go run 工具：執行程式

| 指令 | 說明 |
| --- | --- |
| `go run .` | 編譯並執行目前資料夾的 `main` 套件 |
| `go run ./cmd/goshop -port 9000` | 執行指定套件，**後面的引數**傳給程式 |
| `go run -race .` | 開啟競爭偵測執行 |
| `go run golang.org/x/pkgsite/cmd/pkgsite@latest` | **直接執行遠端模組的程式**，不需要先安裝 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>go run 會快取編譯結果：</b> 程式碼沒有修改時，第二次執行幾乎是瞬間完成。暫存的執行檔放在系統暫存資料夾，執行完就刪除。
</div>

<!--
go run 把「編譯」和「執行」合併成一步，適合開發時使用。

套件路徑之後的所有引數，都會傳給程式本身，這是第 12 章學的 os.Args 和 flag。

最後一種用法很實用：go run 可以直接執行一個遠端模組的程式，加上 @版本號，Go 會自動下載、編譯、執行，不需要先安裝。等一下 go doc 的部分會用它啟動本機的文件網站。

另外 Go 有很好的編譯快取，程式碼沒改的話，第二次 go run 幾乎不需要時間。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# gofmt 工具：程式碼格式化
## gofmt & gopls

<!--
接下來是 gofmt，Go 程式碼統一格式的秘密。
-->

---

# gofmt：統一的程式碼格式

| 指令 | 說明 |
| --- | --- |
| `gofmt -l .` | **列出**格式不正確的檔案（CI 常用） |
| `gofmt -d main.go` | 顯示需要修改的差異 |
| `gofmt -w .` | **直接改寫**檔案 |
| `go fmt ./...` | 等同於 `gofmt -l -w`，格式化所有套件 |
| `gofmt -s -w .` | 額外**簡化**程式碼（例如 `[]T{T{1}}` → `[]T{{1}}`） |

```go
// 格式化前
func add(a int,b int)int{
return a+b}

// 格式化後：Tab 縮排、空白、大括號位置都由 gofmt 決定
func add(a int, b int) int {
	return a + b
}
```

<!--
gofmt 是 Go 官方的格式化工具，它決定了所有 Go 程式碼的長相：用 Tab 縮排、運算子前後加空白、左大括號不換行。

Go 的名言是：「gofmt 的風格不是每個人最喜歡的，但 gofmt 是每個人最喜歡的」。因為所有 Go 程式碼都長得一樣，讀別人的程式碼沒有障礙，Code Review 也不用再爭論格式。

實務上，VS Code 存檔時會自動格式化，所以很少需要手動執行。但在 CI 裡，常會用 gofmt -l 檢查：如果列出任何檔案，代表有人沒格式化就提交了，CI 就會失敗。

補充：這門課的投影片裡所有的 Go 範例，都用 gofmt 檢查過格式。
-->

---

# goimports 與 Go 語言伺服器 gopls

| 工具 | 功能 |
| --- | --- |
| `goimports` | gofmt 的所有功能 + **自動新增 / 刪除 import** |
| **`gopls`** | Go 官方的**語言伺服器（Language Server）**，提供編輯器所有「智慧」功能 |

| gopls 提供的功能 | VS Code 操作 |
| --- | --- |
| 自動完成、參數提示 | 輸入時自動出現 |
| 跳到定義 / 找所有參考 | `F12` / `Shift+F12` |
| 重新命名（整個專案） | `F2` |
| 即時錯誤與 `go vet` 警告 | 紅色 / 黃色波浪線 |
| 自動整理 import、現代化建議 | 存檔時 / 燈泡圖示 |

<!--
gofmt 只負責排版，goimports 更進一步：它會自動加上缺少的 import、刪除沒用到的 import。

gopls 是 Go 官方的語言伺服器，第 0 章安裝 VS Code 擴充套件時就裝好了。所謂語言伺服器，就是一個在背景執行、懂 Go 程式碼的程式，編輯器透過它提供各種智慧功能。所以不只 VS Code，Vim、Emacs、GoLand、Zed 等編輯器都能透過 gopls 得到一樣好的 Go 支援。

最實用的幾個功能：F12 跳到定義、Shift+F12 找所有使用的地方、F2 重新命名，整個專案裡所有用到這個名稱的地方都會一起改。gopls 還會提示「現代化」的寫法，例如建議把三段式迴圈改成 range 整數。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# go vet：程式靜態分析工具
## Static Analysis

<!--
接下來是 go vet：在執行之前，就找出可疑的程式碼。
-->

---

# go vet：找出可疑的程式碼

「**靜態分析**」：不執行程式，只分析程式碼，找出**編譯會過、但很可能是 bug** 的地方

<!-- check:skip -->
```go
package main

import (
	"context"
	"fmt"
	"sync"
)

type Counter struct {
	mu sync.Mutex
	n  int
}

func main() {
	fmt.Printf("%d 件\n", "三") // 格式與型別不符
	c := Counter{}
	c2 := c // 複製了含有 Mutex 的結構
	ctx, _ := context.WithCancel(context.Background())
	fmt.Println(c2.n, ctx)
}
```

```text
$ go vet .
./main.go:15:14: fmt.Printf format %d has arg "三" of wrong type
./main.go:17:8: assignment copies lock value to c2:
    …Counter contains sync.Mutex
./main.go:18:7: the cancel function returned by context.WithCancel
    should be called, not discarded, to avoid a context leak
```

<!--
什麼是靜態分析？就是「不執行程式，只讀程式碼」來找問題。go vet 專門找「編譯會通過，但很可能是 bug」的程式碼。

這個範例有三個 go vet 會抓到的問題，都是前幾章提過的：第 9 章的 Printf 格式與參數不符、第 16 章的複製了含有 Mutex 的結構、以及忘記呼叫 context 的 cancel 函式。

這些錯誤編譯器都不會報錯，程式也能執行，但結果會出錯，或是造成資源洩漏。go vet 在執行之前就把它們抓出來。

實務上，go test 在執行測試之前，會自動跑一部分的 go vet 檢查。VS Code 的 gopls 也會即時顯示 vet 的警告。
-->

---

# go vet 常見的檢查項目

| 分析器 | 找出的問題 | 相關章節 |
| --- | --- | --- |
| `printf` | `Printf` 格式動詞與參數不符 | Ch 9 |
| `copylocks` | 複製了 `sync.Mutex` 等不能複製的值 | Ch 16 |
| `lostcancel` | 忘記呼叫 `context.WithCancel` / `WithTimeout` 的 `cancel` | Ch 16 |
| `structtag` | struct 標籤格式錯誤，例如 `json:name`（少了引號） | Ch 11 |
| `unreachable` | 永遠不會執行到的程式碼 | Ch 6 |
| `stringintconv` | `string(65)` 這種數字轉字串的可疑寫法 | Ch 3 |
| `httpresponse` | 在檢查 `err` 之前使用 `resp` | Ch 14 |

<!--
這張表整理了 go vet 最常見的檢查項目，旁邊標註了相關的章節。大家會發現，很多前幾章提醒過的「注意事項」，其實 go vet 都能自動檢查。

所以一個好習慣是：commit 之前跑一次 go vet ./...，或是在 CI 裡自動執行。
-->

---

# 補充：更完整的檢查工具 — staticcheck 與 golangci-lint

| 工具 | 說明 | 安裝與使用 |
| --- | --- | --- |
| ~~golint~~ | 早期的風格檢查工具，**已於 2021 年停止維護** | 不建議使用 |
| **staticcheck** | 最主流的 Go 靜態分析工具，檢查項目遠多於 vet | `go install honnef.co/go/tools/cmd/staticcheck@latest` → `staticcheck ./...` |
| **golangci-lint** | 整合數十種檢查工具，**CI 中最常用** | 官方安裝腳本或 Homebrew → `golangci-lint run` |
| **govulncheck** | 檢查程式碼**實際呼叫到**的已知安全漏洞 | `go install golang.org/x/vuln/cmd/govulncheck@latest` → `govulncheck ./...` |

```text
$ staticcheck ./...
main.go:12:2: should use strings.Contains instead of … (S1003)
main.go:20:9: error strings should not be capitalized (ST1005)
```

<!--
這是補充內容。

很多舊的教材會介紹 golint，但它在 2021 年就已經停止維護了，官方建議改用 staticcheck。

staticcheck 是目前最主流的 Go 靜態分析工具，它的檢查項目比 go vet 多很多，包括程式碼簡化建議、效能問題、以及風格問題，例如第 6 章說的「錯誤訊息不要大寫開頭」。

golangci-lint 則是一個「整合工具」，它把 staticcheck、vet 和數十種其他檢查工具整合在一起，只要一個設定檔、一個指令，是 CI 裡最常見的選擇。

最後一個是 govulncheck，這是 Go 官方的安全漏洞檢查工具：它會比對漏洞資料庫，而且只回報程式碼「真的有呼叫到」的漏洞，不會一堆誤報。上線前一定要跑一次。
-->

---

# 補充：go fix 自動改寫成現代寫法（Go 1.26+）

`go fix` 內建多個「**現代化分析器（modernizers）**」，自動把舊寫法改成新語法

```bash
go fix -diff .   # 先預覽會修改哪些地方
go fix ./...     # 直接改寫
```

```diff
-	for i := 0; i < 3; i++ {
+	for i := range 3 {                        // Go 1.22 range 整數
-	m := 3
-	if 5 > m {
-		m = 5
-	}
+	m := max(5, 3)                            // Go 1.21 內建 max
-	for _, p := range strings.Split(s, ",") {
+	for p := range strings.SplitSeq(s, ",") { // Go 1.24 迭代器
-	for _, x := range xs { if x == s { return true } }; return false
+	return slices.Contains(xs, s)             // Go 1.21 slices 套件
```

<!--
這是補充內容：Go 1.26 全面翻新了 go fix 工具。

go fix 內建了一組「現代化分析器」，會找出可以用新語法改寫的舊程式碼，並且自動改寫。這門課教的很多新寫法：range 整數、內建的 min 和 max、strings.SplitSeq、slices.Contains，go fix 都能自動幫我們轉換。

建議先用 -diff 預覽會改哪些地方，確認沒問題再真正改寫。

這對維護舊專案特別有用：升級 Go 版本之後，跑一次 go fix，程式碼就自動跟上最新的寫法。這也呼應了第 0 章說的 Go 1 相容性承諾：舊程式碼不改也能跑，但工具會幫我們逐步現代化。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# go doc 工具：產生文件
## Documentation

<!--
接下來是文件：怎麼查詢文件、怎麼替自己的程式碼寫文件。
-->

---

# go doc：在終端機查詢文件

| 指令 | 說明 |
| --- | --- |
| `go doc strings` | 套件的說明與所有匯出的名稱 |
| `go doc strings.Cut` | 單一函式的說明 |
| `go doc -short strings` | 只列出簽章 |
| `go doc -all ./internal/order` | 自己專案的套件，顯示所有文件 |
| `go doc -src strings.Cut` | 顯示原始碼 |

```text
$ go doc strings.Cut
package strings // import "strings"

func Cut(s, sep string) (before, after string, found bool)
    Cut slices s around the first instance of sep, returning the text
    before and after sep. The found result reports whether sep appears
    in s. If sep does not appear in s, cut returns s, "", false.
```

<!--
go doc 可以在終端機直接查詢文件，不用打開瀏覽器。

後面接套件名稱，列出套件的說明和所有匯出的名稱；接「套件.函式」，顯示單一函式的說明；-src 還能直接看原始碼。

不只標準函式庫，自己專案的套件、下載的第三方模組都能查詢。

實際上，大部分的人還是習慣用瀏覽器看 pkg.go.dev，它就是用同樣的文件產生的，只是排版比較漂亮。
-->

---

# 撰寫文件註解（doc comments）

**匯出名稱前面的註解，就是它的文件**；以名稱開頭、完整的句子

```go
// Package pricing 負責商品價格與稅金的計算。
package pricing

// TaxRate 是營業稅率。
const TaxRate = 0.05

// WithTax 回傳含稅價格，四捨五入到小數點後兩位。
// 價格為負數時回傳 0。
func WithTax(price float64) float64 {
	if price < 0 {
		return 0
	}
	return float64(int(price*(1+TaxRate)*100+0.5)) / 100
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>本機預覽文件網站：</b> <code>go run golang.org/x/pkgsite/cmd/pkgsite@latest</code>，開啟 <code>http://localhost:8080</code> 就能看到跟 pkg.go.dev 一樣的文件頁面。
</div>

<!--
Go 的文件就是程式碼裡的註解，不需要額外的文件工具或特殊語法。

規則很簡單：寫在匯出名稱「正上方」、中間沒有空行的註解，就是它的文件。慣例是以名稱開頭，寫成完整的句子，例如「WithTax 回傳含稅價格」。套件的文件則以「Package 套件名稱」開頭，寫在 package 宣告的上方，第 8 章提過。

寫好之後，go doc 和 pkg.go.dev 都會自動讀取。想在本機預覽，可以用 go run 直接執行 pkgsite，它會產生一個跟 pkg.go.dev 一模一樣的網站。

staticcheck 等工具會檢查匯出的名稱有沒有文件，這是寫套件給別人用的基本禮貌。
-->

---

# 補充：Example 函式 — 可以執行的文件

在測試檔中寫 `ExampleXxx` 函式，**同時是文件、也是測試**

```go
// 檔名：pricing.go
package pricing

// WithTax 回傳含稅價格。
func WithTax(price float64) float64 { return price * 1.05 }
```

```go
// 檔名：example_test.go
package pricing

import "fmt"

func ExampleWithTax() {
	fmt.Println(WithTax(100))
	// Output: 105
}
```

<!--
這是補充內容：Example 函式。

在測試檔裡寫一個 Example 開頭的函式，最後加上 // Output: 註解，寫出預期的輸出。它有兩個作用：

第一，它是文件：go doc 和 pkg.go.dev 會把它顯示成使用範例，而且在 pkg.go.dev 上還能直接按 Run 執行。

第二，它是測試：go test 會執行它，並比對實際輸出和 Output 註解是不是一樣，不一樣就算測試失敗。這保證了文件裡的範例永遠是正確的，不會因為程式碼改了、文件忘了改而過時。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# go get 工具：下載模組或套件
## go get・go install・go tool

<!--
最後是模組與工具的管理。
-->

---

# go get：管理專案的相依模組

| 指令 | 說明 |
| --- | --- |
| `go get example.com/pkg` | 新增相依模組（最新版） |
| `go get example.com/pkg@v1.2.3` | 指定版本 |
| `go get -u ./...` | 升級所有相依模組到最新的 minor / patch 版本 |
| `go get example.com/pkg@none` | 移除相依模組 |
| `go get go@1.27` | 升級專案的 Go 版本（修改 `go.mod` 的 `go` 那一行） |
| `go mod tidy` | 整理：補上缺少的、刪除沒用到的 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>go get 不再用來安裝工具：</b> 從 Go 1.17 起，<code>go get</code> 只負責修改 <code>go.mod</code>；安裝執行檔請用 <code>go install</code>。網路上 <code>go get -u github.com/xxx/cmd/tool</code> 的舊教學已經不適用。
</div>

<!--
go get 用來管理專案的相依模組，第 8 章已經用過了。

這張表多列了幾個實用的用法：-u 升級相依模組；@none 移除；go get go@1.27 可以升級專案要求的 Go 版本。

特別提醒：在 Go 1.16 以前，go get 同時用來「下載套件」和「安裝工具」，網路上有很多舊教學寫 go get -u 某個工具。從 Go 1.17 開始，go get 只負責修改 go.mod，安裝工具要改用 go install，下一頁說明。
-->

---

# go install 與 go tool：安裝開發工具

| 方式 | 指令 | 特性 |
| --- | --- | --- |
| **全域安裝** | `go install honnef.co/go/tools/cmd/staticcheck@latest` | 安裝到 `$GOPATH/bin`，所有專案共用 |
| **專案工具**（Go 1.24+） | `go get -tool golang.org/x/tools/cmd/stringer` | 記錄在 `go.mod`，**版本跟著專案走** |
| 執行專案工具 | `go tool stringer -type=Status` | 團隊每個人用同一個版本 |

```text
// go.mod
module example.com/goshop

go 1.27.0

tool golang.org/x/tools/cmd/stringer   ← go get -tool 自動加入
```

<!--
安裝 Go 寫的工具有兩種方式。

第一種是 go install 加上 @版本，會把工具編譯後安裝到 GOPATH/bin，第 0 章我們把這個資料夾加進了 PATH，所以裝完就能直接執行。這適合個人常用的工具。

第二種是 Go 1.24 加入的「專案工具」：go get -tool 會把工具記錄在 go.mod 裡的 tool 指令，之後用 go tool 工具名稱執行。好處是工具的版本跟著專案走，團隊每個人、CI 用的都是同一個版本，不會發生「我這邊產生的程式碼跟你不一樣」的問題。

這裡的 stringer，就是第 1 章提過的、自動替列舉產生 String 方法的官方工具。
-->

---
layout: default
---

# 綜合練習：打造一個可發布的命令列工具
### 任務說明

1. 建立模組 `example.com/greet`，寫一個命令列工具，旗標 `-name` 與 `-version`
2. 宣告 `var version = "dev"`，`-version` 時印出版本號
3. 寫一個匯出函式 `Greeting(name string) string` 並加上**文件註解**，再寫一個 `ExampleGreeting` 測試
4. 依序執行：`gofmt -l .` → `go vet ./...` → `go test ./...`，全部通過
5. 用 `-ldflags "-s -w -X main.version=1.0.0"` 分別編譯出 **Windows、macOS（arm64）、Linux** 三個執行檔
6. 用 `go version -m` 檢查其中一個執行檔

<!--
這個綜合練習模擬一個工具從開發到發布的完整流程：寫程式、寫文件、格式化、靜態分析、測試、跨平台編譯。

實務上，這些步驟通常會寫成一個 Makefile 或 CI 腳本自動執行，每次 push 程式碼都會跑一遍。
-->

---

# 綜合練習：解題提示
### 提示說明

```go
// 檔名：main.go
package main

import (
	"flag"
	"fmt"
)

var version = "dev"

// Greeting 回傳對 name 的問候語；name 為空字串時使用「訪客」。
func Greeting(name string) string {
	if name == "" {
		name = "訪客"
	}
	return "你好，" + name + "！"
}

func main() {
	name := flag.String("name", "", "要問候的名字")
	showVer := flag.Bool("version", false, "顯示版本")
	flag.Parse()
	if *showVer {
		fmt.Println("greet", version)
		return
	}
	fmt.Println(Greeting(*name))
}
```

<!--
main.go 包含一個有文件註解的匯出函式 Greeting，以及處理旗標的 main。version 是套件層級的字串變數，等一下用 -ldflags 注入。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 檔名：main_test.go
package main

import "fmt"

func ExampleGreeting() {
	fmt.Println(Greeting("Gopher"))
	fmt.Println(Greeting(""))
	// Output:
	// 你好，Gopher！
	// 你好，訪客！
}
```

```bash
gofmt -l . && go vet ./... && go test ./...
flags='-s -w -X main.version=1.0.0'
GOOS=windows GOARCH=amd64 go build -ldflags "$flags" -o dist/ .
GOOS=darwin  GOARCH=arm64 go build -ldflags "$flags" -o dist/mac/ .
GOOS=linux   GOARCH=amd64 go build -ldflags "$flags" -o dist/linux/ .
go version -m dist/linux/greet
```

<!--
ExampleGreeting 的 Output 可以有多行，每一行都要跟實際輸出完全一樣，包括全形的驚嘆號。

下面的指令依序執行格式檢查、靜態分析、測試，用 && 串接，任何一步失敗就會停下來。接著用三組 GOOS、GOARCH 編譯出三個平台的執行檔，放在 dist 資料夾。最後用 go version -m 確認執行檔的資訊。

把這幾行寫成一個腳本，就是一個最簡單的發布流程。
-->

---

# 章節總結

- **go build**：`-o` 指定輸出；`-ldflags "-s -w"` 縮小；`-ldflags "-X main.version=…"` 注入版本；`-trimpath`
- **編譯條件**：檔名後綴 `_windows.go`；`//go:build linux || darwin`；自訂標籤 `-tags`
- **跨平台編譯**：`GOOS` + `GOARCH`；`go tool dist list`；純 Go 程式可 `CGO_ENABLED=0`
- **go run**：`go run .`；`go run 模組@版本` 直接執行遠端工具
- **格式與編輯器**：`gofmt -l -w` / `go fmt`；`goimports`；**gopls** 提供所有編輯器智慧功能
- **靜態分析**：`go vet ./...`；golint 已停止維護 → **staticcheck**、**golangci-lint**、**govulncheck**；`go fix` 現代化（1.26+）
- **文件**：`go doc`；文件註解以名稱開頭；`ExampleXxx` 是文件也是測試
- **模組與工具**：`go get` 只管 `go.mod`；工具用 `go install` 或 `go get -tool` + `go tool`（1.24+）

下一章我們會介紹「加密安全」：雜湊、對稱與非對稱加密、數位簽章與 HTTPS。

<!--
我們來整理今天學到的東西。

go build 可以用 ldflags 注入版本、用 GOOS 和 GOARCH 跨平台編譯。gofmt 讓所有 Go 程式碼長得一樣，gopls 讓編輯器變聰明。go vet 找出可疑的程式碼，staticcheck、golangci-lint、govulncheck 提供更完整的檢查，go fix 可以自動現代化程式碼。go doc 查詢文件，Example 函式是可以執行的文件。go get 管理相依模組，go install 和 go tool 管理工具。

下一章要進入一個很重要的主題：加密安全。密碼要怎麼存、資料要怎麼加密、怎麼確認資料沒有被竄改、HTTPS 是怎麼運作的。Go 的標準函式庫 crypto 套件非常完整，而且是由密碼學專家維護的，我們會學到怎麼正確地使用它。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的工具會陪伴大家的整個 Go 開發生涯，建議把 gofmt、go vet、go test 設定成存檔或 commit 時自動執行。

課後建議：替前幾章寫過的專案，跑一次 staticcheck 和 go fix -diff，看看工具會給出什麼建議。

有問題的同學現在可以提問！
-->

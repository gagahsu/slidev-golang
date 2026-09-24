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
title: 前置作業 — 開發環境與第一支程式
routeAlias: ch00
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">前置作業</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「工具準備好，第一行 Go 程式就跑得起來」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到 Go 實戰開發課程！

在正式學 Go 的語法之前，我們要先把「工具」準備好。就像要開始做菜之前，得先把爐子、鍋子、刀子都擺上流理台一樣，今天這一章就是在做這件事。

我們會安裝 Go 語言本身、安裝 VS Code 編輯器，接著用指令建立第一個 Go 專案並執行它。最後還會介紹一個不用安裝任何東西就能寫 Go 的網站：Go Playground。

學完這一章，大家就能在自己的電腦上寫出、並執行第一支 Go 程式。
-->

---
layout: default
---

# Outline

- **安裝 Go 語言** — 下載安裝、確認版本、認識 `go env`
- **安裝 Visual Studio Code** — 安裝編輯器與官方 Go 擴充套件
- **建立和執行專案** — `go mod init`、`go run`、`go build`
- **The Go Playground** — 在瀏覽器裡直接寫 Go
- **從 VS Code 執行 Go 程式** — 終端機執行、Run and Debug、中斷點
- **GoShop 專案實作** — 認識課程專案，建立 `goshop` 專案
- **章節總結**

<!--
今天的內容是一條直線：先裝 Go，再裝編輯器，然後建立專案、執行程式。

中間穿插一個 Go Playground，它很適合拿來快速測試一小段程式碼，後面的章節我們也會常常用到它。

最後會教大家怎麼在 VS Code 裡面直接按一個鍵就執行程式，甚至設中斷點一行一行看程式怎麼跑。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 安裝 Go 語言
## Installing Go

<!--
第一步，安裝 Go 語言本身。

這裡說的「Go 語言」，其實是一整套工具，官方稱為 Go toolchain（工具鏈），裡面包含編譯器、套件管理、格式化、測試等等工具，全部都整合在一個 go 指令裡面。
-->

---

# 什麼是 Go toolchain？

| 項目 | 說明 |
| --- | --- |
| **官方網站** | [go.dev](https://go.dev)，下載頁面為 [go.dev/dl](https://go.dev/dl) |
| **安裝內容** | 編譯器、標準函式庫、以及整合在一起的 `go` 指令 |
| **發布節奏** | 每年 2 月、8 月各發布一個大版本（例如 1.26、1.27） |
| **本課程版本** | 範例以 **Go 1.25 以上**的語法撰寫，並以 **Go 1.27** 驗證 |
| **相容承諾** | Go 1 相容性保證：舊程式在新版本上幾乎都能直接編譯 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>一個指令做完所有事：</b> 其他語言常常需要「編譯器 + 套件管理工具 + 格式化工具 + 測試框架」分別安裝；Go 全部整合在 <code>go</code> 這個指令裡（<code>go build</code>、<code>go test</code>、<code>go fmt</code>……）。
</div>

<!--
什麼是 Go toolchain？我們可以把它想像成一個「工具箱」，打開來裡面什麼都有。

其他語言常常要分別安裝好幾個工具：Java 要 JDK 加上 Maven 或 Gradle；JavaScript 要 Node.js 加上 npm、Prettier、Jest。Go 很不一樣，官方把這些都整合在一個 go 指令裡，這也是很多人喜歡 Go 的原因之一：「裝好就能用」。

版本的部分，Go 每半年發布一個大版本，2 月跟 8 月。本課程的範例使用 Go 1.25 以上的語法，並且用最新的 Go 1.27 實際編譯驗證過。大家安裝時直接裝官網最新的穩定版就可以了。

另外 Go 有一個很重要的承諾：「Go 1 相容性保證」。意思是用舊版本寫的程式，升級到新版本幾乎都能直接編譯，所以升級不用太擔心。
-->

---

# 在各作業系統安裝 Go

| 作業系統 | 安裝方式 | 預設安裝位置 |
| --- | --- | --- |
| **Windows** | 下載 `.msi` 安裝檔，一路按「Next」 | `C:\Program Files\Go` |
| **macOS** | 下載 `.pkg` 安裝檔，或 `brew install go` | `/usr/local/go` |
| **Linux** | 下載 `.tar.gz` 解壓縮到 `/usr/local` | `/usr/local/go` |

Linux 的安裝指令（以 amd64 為例，版本號請換成官網最新版）：

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.27.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

<!--
安裝方式依作業系統不同，但都很簡單。

Windows 跟 macOS 最簡單，下載安裝檔一路按下一步就好，安裝程式會自動幫我們把 go 指令加進 PATH 環境變數。macOS 如果有用 Homebrew，也可以直接 brew install go。

Linux 需要自己解壓縮。注意第一行 rm -rf，官方文件特別提醒：如果之前裝過舊版，一定要先刪掉舊的資料夾再解壓縮，不然新舊檔案混在一起會出問題。

最後兩行是把 go 指令加進 PATH。注意我們同時加了 $HOME/go/bin，這個資料夾是之後用 go install 安裝工具的地方，先加好之後就不用再煩惱。
-->

---

# 確認安裝成功：go version 與 go env

安裝完成後，**重新開啟終端機**，輸入：

```bash
go version
# go version go1.27.0 linux/amd64

go env GOROOT GOPATH
# /usr/local/go
# /home/user/go
```

| 環境變數 | 意義 |
| --- | --- |
| `GOROOT` | Go 本身安裝在哪裡（編譯器、標準函式庫） |
| `GOPATH` | 我們的工作區，下載的模組快取、`go install` 的執行檔放這裡 |

<!--
安裝完成後，第一件事就是確認有沒有裝好。

注意一定要「重新開啟」終端機，因為 PATH 環境變數只有在新開的終端機裡才會生效。這是初學者最常卡住的地方：明明裝好了，打 go version 卻說找不到指令，通常就是終端機沒重開。

go version 會印出版本號跟作業系統架構。go env 則可以查詢 Go 的環境設定，這裡我們先認識兩個：GOROOT 是 Go 本身裝在哪裡，GOPATH 是我們的工作區。

GOPATH 在早期的 Go 很重要，所有程式都要放在裡面；但現在有了 Go Modules，專案可以放在任何地方，GOPATH 主要只剩下「放快取跟工具」的用途。第 8 章講套件的時候會再詳細說明。
-->

---

# 補充：Go 會自動切換工具鏈版本

從 Go 1.21 開始，`go.mod` 裡寫的 `go` 版本如果比本機新，Go 會**自動下載對應版本**：

```bash
cat go.mod
# module example.com/hello
#
# go 1.27.0

go version   # 本機是 1.25，但專案要求 1.27
# go: downloading go1.27.0 (linux/amd64)
# go version go1.27.0 linux/amd64
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>GOTOOLCHAIN：</b> 這個行為由環境變數 <code>GOTOOLCHAIN</code> 控制，預設值是 <code>auto</code>。如果公司環境不能連網下載，可以設成 <code>local</code> 強制只用本機版本。
</div>

<!--
這是一個比較新的功能，大家先有印象就好。

以前如果同事的專案用了比較新的 Go 版本，我們得自己去官網下載升級。從 Go 1.21 開始，Go 會看 go.mod 裡面寫的版本，如果本機的版本比較舊，它會自動下載對應版本來用。

這就像手機 App 會自動更新一樣，讓團隊裡每個人用的 Go 版本都一致，減少「在我電腦上可以跑」的問題。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 安裝 Visual Studio Code
## Editor Setup

<!--
Go 裝好了，接下來要一個好用的編輯器。

寫 Go 的主流編輯器有兩個：免費的 VS Code，以及 JetBrains 付費的 GoLand。這門課我們用 VS Code，因為它免費、輕量，而且 Go 團隊官方維護了 VS Code 的擴充套件。
-->

---

# 安裝 VS Code 與 Go 擴充套件

| 步驟 | 動作 |
| --- | --- |
| **1. 下載 VS Code** | 到 [code.visualstudio.com](https://code.visualstudio.com) 下載並安裝 |
| **2. 安裝 Go 擴充套件** | 左側 Extensions（`Ctrl+Shift+X`）搜尋 **Go**，安裝發行者為 **Go Team at Google** 的版本 |
| **3. 安裝 Go 工具** | `Ctrl+Shift+P` → 輸入 **Go: Install/Update Tools** → 全選 → OK |
| **4. 確認** | 打開任何 `.go` 檔，右下角出現 Go 版本號即完成 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>認明發行者：</b> 市集上有很多名字含 Go 的套件，請安裝 ID 為 <code>golang.go</code>、由 Go Team at Google 發行的官方版本。
</div>

<!--
VS Code 本身安裝很簡單，下載一路下一步就好。

重點是第二步：安裝 Go 擴充套件。在市集裡搜尋 Go 會出現很多結果，要認明發行者是「Go Team at Google」，ID 是 golang.go，這是 Go 團隊官方維護的版本。

第三步是安裝擴充套件需要的工具，最重要的是 gopls（Go 的語言伺服器，負責自動完成、跳到定義、錯誤提示）跟 dlv（除錯器）。通常打開第一個 .go 檔時，VS Code 就會跳通知問要不要安裝，按 Install All 就可以了。
-->

---

# Go 擴充套件幫我們做了什麼？

| 功能 | 背後的工具 | 效果 |
| --- | --- | --- |
| 自動完成、跳到定義 | `gopls` | 打 `fmt.` 就列出所有函式 |
| 存檔自動排版 | `gopls`（呼叫 `gofmt`） | 縮排、空白自動整理成官方格式 |
| 自動 import | `gopls` | 用到 `strings.ToUpper` 自動補上 `import "strings"` |
| 錯誤即時提示 | `gopls` + `go vet` | 沒用到的變數，存檔前就畫紅線 |
| 除錯 | `dlv`（Delve） | 設中斷點、逐行執行、查看變數 |

<!--
安裝完之後，擴充套件會在背後幫我們做很多事情。

最有感的是「存檔自動排版」跟「自動 import」。在 Go 的世界裡，程式碼的格式是由官方工具 gofmt 統一決定的，所以不管誰寫的 Go 程式碼，長相都一模一樣。我們不需要跟同事爭論要用 tab 還是空白、大括號要不要換行，存檔的瞬間工具就幫我們排好了。

自動 import 也很方便：我們在程式裡用到 strings.ToUpper，存檔時 import "strings" 就會自動出現；反過來，沒用到的 import 也會被自動刪掉。這很重要，因為 Go 規定「沒用到的 import 是編譯錯誤」。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 建立和執行專案
## go mod init・go run・go build

<!--
工具都準備好了，現在來建立我們的第一個 Go 專案。
-->

---

# 建立第一個 Go 專案

```bash
mkdir hello          # 建立專案資料夾
cd hello
go mod init example.com/hello   # 初始化模組，產生 go.mod
code .               # 用 VS Code 開啟目前資料夾
```

執行 `go mod init` 後，資料夾裡會多出 `go.mod`：

```text
module example.com/hello

go 1.27.0
```

<!--
每一個 Go 專案，都從 go mod init 開始。

go mod init 後面接的是「模組路徑」，也就是這個專案的名字。慣例上會用網址的形式，例如 github.com/帳號/專案名稱；如果只是練習、沒有要發布，用 example.com/hello 這種寫法就可以了。

執行完之後會產生 go.mod 檔案。它就像是專案的「身分證」：第一行寫專案叫什麼名字，第二行寫這個專案需要哪個版本以上的 Go。之後如果用到第三方套件，也會記錄在這個檔案裡。
-->

---

# 寫下第一支 Go 程式：main.go

在專案資料夾新增 `main.go`：

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")
}
```

| 程式碼 | 意義 |
| --- | --- |
| `package main` | 這個檔案屬於 `main` 套件，代表它會被編譯成**可執行檔** |
| `import "fmt"` | 匯入標準函式庫的 `fmt` 套件（格式化輸出入） |
| `func main()` | 程式的進入點，執行時從這裡開始 |

<!--
這就是最經典的 Hello World。雖然只有短短幾行，但每一行都有它的意義。

第一行 package main：Go 的每個檔案都要宣告自己屬於哪個套件。main 是一個特別的名字，代表「這是一個可以執行的程式」，而不是給別人用的函式庫。

第三行 import "fmt"：匯入 fmt 套件，fmt 是 format 的縮寫，負責印出文字。

最後 func main()：程式的進入點。執行程式的時候，Go 會從 main 套件裡的 main 函式開始跑。

另外注意縮排：Go 官方格式用的是 Tab，不是空白。不用刻意記，存檔時 gofmt 會自動幫我們處理。
-->

---

# 執行程式：go run 與 go build

```bash
go run .        # 編譯並立刻執行（不留下執行檔）
# Hello, Go!

go build        # 編譯出執行檔
./hello         # Linux / macOS 執行
.\hello.exe     # Windows 執行
# Hello, Go!
```

| 指令 | 做了什麼 | 適合時機 |
| --- | --- | --- |
| `go run .` | 編譯到暫存資料夾並執行，跑完就刪掉 | 開發時快速測試 |
| `go build` | 編譯出一個獨立的執行檔 | 要部署、要交給別人使用 |

<!--
程式寫好了，有兩種執行方式。

go run 點，這個點代表「目前資料夾」。它會把程式編譯到一個暫存的地方，執行完就刪掉，所以資料夾裡不會多出任何檔案，很適合開發時一直改、一直跑。

go build 則會真的產生一個執行檔。Go 編譯出來的執行檔是「靜態連結」的，意思是這個檔案可以直接複製到另一台同樣作業系統的電腦上執行，對方完全不用安裝 Go。這是 Go 在部署上非常受歡迎的原因：一個檔案丟上伺服器就能跑。

執行後，console 會輸出 Hello, Go!
-->

---

# 使用 go run 的注意事項

**注意事項之一：** 建議用 `go run .` 而不是 `go run main.go`

```bash
go run main.go   # 只編譯 main.go，同資料夾其他 .go 檔不會被編進來
go run .         # 編譯整個套件（資料夾內所有 .go 檔）✅
```

**注意事項之二：** 沒用到的變數和 import 都是**編譯錯誤**

```go
package main

import "fmt" // 編譯錯誤：imported and not used

func main() {
	x := 10 // 編譯錯誤：declared and not used
}
```

<!--
使用 go run 有兩個注意事項。

第一個：很多網路教學會寫 go run main.go，但這樣只會編譯 main.go 這一個檔案。等到我們的程式拆成好幾個檔案時，就會出現「找不到函式」的錯誤。所以從現在開始就養成習慣，用 go run 點，編譯整個資料夾。

第二個是 Go 很有個性的地方：沒用到的變數、沒用到的 import，在其他語言頂多是警告，在 Go 是直接編譯失敗。一開始可能會覺得很囉嗦，但這個設計讓 Go 的程式碼永遠保持乾淨，不會堆滿沒用的東西。
-->

---

# 補充：常用的 go 指令一覽

| 指令 | 用途 | 本課程章節 |
| --- | --- | --- |
| `go mod init` | 建立新模組 | 本章、Ch 8 |
| `go run` / `go build` | 執行／編譯 | 本章、Ch 17 |
| `go fmt` | 格式化程式碼 | Ch 17 |
| `go vet` | 靜態分析，找出可疑的程式碼 | Ch 17 |
| `go test` | 執行單元測試 | Ch 9 |
| `go get` / `go mod tidy` | 下載、整理相依模組 | Ch 8、Ch 17 |
| `go doc` | 查詢文件 | Ch 17 |

<!--
這張表先讓大家有個全貌，不用現在就記起來。

每個指令都會在後面的章節正式登場，旁邊標註了對應的章節。今天只要會 go mod init、go run、go build 這三個就足夠了。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# The Go Playground
## go.dev/play

<!--
接下來介紹一個不用安裝任何東西就能寫 Go 的網站：Go Playground。
-->

---

# 什麼是 Go Playground？

| 項目 | 說明 |
| --- | --- |
| **網址** | [go.dev/play](https://go.dev/play) |
| **功能** | 在瀏覽器寫 Go、按 **Run** 執行、按 **Format** 排版 |
| **分享** | 按 **Share** 產生短網址，貼給別人就能看到同一份程式碼 |
| **版本** | 右上角可切換 Go 版本（最新版、前一版、開發版） |
| **限制** | 不能連網、不能讀寫本機檔案、執行時間有上限 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>時間是固定的：</b> Playground 裡的 <code>time.Now()</code> 永遠從 <b>2009-11-10 23:00:00 UTC</b> 開始（Go 公開發表的日子），所以跟時間有關的範例請在本機執行。
</div>

<!--
Go Playground 是官方提供的線上執行環境，打開網頁就能寫 Go。

它最好用的功能是 Share：寫好一段程式碼，按 Share 會產生一個短網址，把網址貼給同事或貼到論壇上，對方點開就能看到一模一樣的程式碼，還可以直接執行。在網路上問問題時，附上 Playground 連結是非常主流的做法。

它也有一些限制：因為是在官方的沙盒（sandbox）裡跑，所以不能連網路、不能讀寫檔案。另外有一個很有趣的細節：Playground 裡的時間是固定的，永遠從 2009 年 11 月 10 日開始，那是 Go 公開發表的日子。所以第 10 章講時間處理的時候，範例請在自己的電腦上跑。
-->

---

# 在 Playground 試試看

把下面的程式碼貼到 [go.dev/play](https://go.dev/play)，按下 **Run**：

```go
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Hello, Playground!")
	fmt.Println("Go 版本：", runtime.Version())
}
```

執行後，下方會輸出 `Hello, Playground!` 以及目前 Playground 使用的 Go 版本。

<!--
這段程式碼的目的，是確認 Playground 用的是哪個 Go 版本。

注意這裡 import 了兩個套件，所以用小括號把它們包起來，每行一個，這是 Go 匯入多個套件的標準寫法。

runtime.Version() 會回傳目前執行的 Go 版本。執行後，下方的輸出區會印出 Hello, Playground! 和版本號。大家可以試著切換右上角的版本選單，再按一次 Run，看看版本號有沒有變。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 從 VS Code 執行 Go 程式
## Run & Debug

<!--
最後，我們來看怎麼在 VS Code 裡面直接執行、除錯 Go 程式。
-->

---

# 在 VS Code 執行程式的三種方式

| 方式 | 操作 | 適合時機 |
| --- | --- | --- |
| **整合終端機** | `` Ctrl+` `` 開啟終端機，輸入 `go run .` | 最通用，跟實際部署時的操作一致 |
| **Run Without Debugging** | `Ctrl+F5` | 快速執行目前的程式 |
| **Start Debugging** | `F5` | 需要設中斷點、逐行檢查時 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>建議：</b> 初學階段先習慣用終端機 <code>go run .</code>，因為之後的工具（<code>go test</code>、<code>go build</code>）都在終端機操作。
</div>

<!--
在 VS Code 裡執行 Go 程式有三種方式。

第一種是整合終端機，按 Ctrl 加上反引號就能打開，然後輸入 go run 點。這是最通用的方式，我們也建議大家從這裡開始。

第二種是 Ctrl+F5，直接執行不除錯。第三種是 F5，啟動除錯模式，可以搭配中斷點使用，下一頁我們來看怎麼用。
-->

---

# 設定中斷點與除錯

1. 在行號左邊點一下，出現紅點就是**中斷點（breakpoint）**
2. 按 `F5` 啟動除錯，程式會停在中斷點
3. 左側 **VARIABLES** 面板可以看到目前所有變數的值
4. 用上方工具列控制執行：`F10` 逐行、`F11` 進入函式、`F5` 繼續

需要固定的執行設定時，可建立 `.vscode/launch.json`：

```json
{
  "version": "0.2.0",
  "configurations": [
    { "name": "Launch Package", "type": "go", "request": "launch",
      "mode": "auto", "program": "${workspaceFolder}" }
  ]
}
```

<!--
除錯器是找 bug 的神兵利器。

操作很簡單：在想停下來的那一行左邊點一下，出現紅點；按 F5 執行，程式跑到那一行就會停住。這時候左邊的面板會顯示所有變數目前的值，我們可以按 F10 一行一行往下執行，觀察變數怎麼變化。

如果每次都要用同樣的設定執行，例如要帶命令列參數，可以建立 launch.json。program 設成 workspaceFolder，代表執行整個專案資料夾，跟 go run 點的效果一樣。

背後負責除錯的工具叫做 Delve，也就是剛剛安裝工具時裝的 dlv。
-->

---
layout: default
---

# 練習 1：建立自我介紹程式
### 任務說明

1. 建立一個新資料夾 `intro`，並用 `go mod init example.com/intro` 初始化
2. 新增 `main.go`，印出三行文字：
   - 你的名字
   - 你想用 Go 做什麼
   - 目前使用的 Go 版本（使用 `runtime.Version()`）
3. 分別用 `go run .` 和 `go build` 執行一次
4. 把程式貼到 Go Playground，按 **Share** 取得分享網址

<!--
這個練習把今天學到的東西全部串起來：建立模組、寫程式、兩種執行方式，再加上 Playground 分享。

大家先自己動手試試看，卡住了再看下一頁的提示。
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("我是小明")
	fmt.Println("我想用 Go 寫後端 API")
	fmt.Println("Go 版本：", runtime.Version())
}
```

- `go build` 後，Linux/macOS 用 `./intro`，Windows 用 `.\intro.exe` 執行
- 如果出現 `go: cannot find main module`，代表忘了先執行 `go mod init`

<!--
參考答案就是 Hello World 的延伸，多印了幾行而已。

最常見的錯誤訊息是 go: cannot find main module，意思是找不到 go.mod，通常就是忘記 go mod init，或者終端機所在的資料夾不對。遇到這個錯誤，先用 pwd（Windows 用 cd）確認目前在哪個資料夾。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 0 步：建立專案

<!--
最後一個主題，是這門課的重頭戲：課程專案 GoShop。

學程式最怕的就是「每個語法都看得懂，但湊不起來」。所以這門課除了每章的練習題，還會帶大家從零開始，一章一章把一個完整的系統做出來。
-->

---

# 課程專案：GoShop 迷你電商後台

每一章的最後，我們都會替**同一個專案**加上新功能，上完整門課就完成一個完整的系統

| 功能 | 說明 | 用到的章節 |
| --- | --- | --- |
| 商品管理 | 商品目錄、CSV 匯入、後台網頁新增商品 | Ch 3、4、12、15 |
| 購物車與折扣 | 滿額折扣、會員折扣、有期限的折價券 | Ch 2、5、10 |
| 訂單與付款 | 庫存檢查、信用卡／錢包／貨到付款 | Ch 6、7 |
| 資料儲存 | gob 快照、JSON、MySQL 交易 | Ch 11、12、13 |
| 網站 | RESTful JSON API、後台管理網頁、HTTPS | Ch 15、18 |
| 品質與上線 | 單元測試、並行安全、跨平台編譯 | Ch 9、16、17 |

<!--
我們要做的專案叫做 GoShop，是一個迷你的電商後台。

選電商當題目的原因很簡單：大家都在網路上買過東西，不需要花時間解釋「購物車是什麼」「訂單為什麼要扣庫存」，可以把力氣都放在學 Go。

這張表是 GoShop 最後會有的功能。一開始它只是一個印出歡迎訊息的小程式，到第 19 章，它會有 REST API、後台網頁、MySQL 資料庫、登入機制，而且可以同時服務很多顧客而不會超賣。

每一個功能都對應到某幾章的內容，所以大家學完一章，馬上就會看到這章的語法在真實專案裡怎麼用。
-->

---

# GoShop 的成長路線（Ch 0 ～ 9）

| 章 | 這一步要做的事 |
| --- | --- |
| Ch 0 | 建立專案、印出歡迎訊息 |
| Ch 1 | 用變數、常數試算一筆訂單；用指標扣庫存 |
| Ch 2 | 會員折扣（`switch`）、滿額免運（`if`）、數量試算表（`for`） |
| Ch 3 | 解析供應商傳來的商品字串、對齊中文品名 |
| Ch 4 | 用結構表示商品和購物車，用 `map` 做商品目錄 |
| Ch 5 | 用閉包寫折扣規則，印出含千分位的收據 |
| Ch 6 | 庫存不足、商品不存在的錯誤處理 |
| Ch 7 | `Money` 型別與付款方式介面 |
| Ch 8 | 把一個大檔案拆成多個套件 |
| Ch 9 | 替結帳流程寫單元測試，加上結構化日誌 |

<!--
這是前半段的路線。

前 7 章我們都寫在同一個 main.go 裡面，因為這時候還沒學到套件。每一章都會在上一章的基礎上修改和擴充，所以大家做完一章，記得把程式碼保留下來，下一章接著用。

到第 8 章，main.go 會長到快 300 行，這時候我們就會很有感覺：「這個檔案太大了，應該要拆開」，剛好就是學套件的最好時機。第 9 章則會替它寫單元測試，而且測試還會幫我們抓到一個真正的 bug。
-->

---

# GoShop 的成長路線（Ch 10 ～ 19）

| 章 | 這一步要做的事 |
| --- | --- |
| Ch 10 | 有使用期限的折價券、計算預計出貨日 |
| Ch 11 | 從 JSON 匯入商品，用 gob 把資料存檔 |
| Ch 12 | 命令列旗標、CSV 匯入匯出、訂單日誌 |
| Ch 13 | 改用 MySQL，用交易（transaction）扣庫存 |
| Ch 14 | 呼叫匯率 API、用 webhook 通知倉庫 |
| Ch 15 | RESTful JSON API 與後台管理網頁 |
| Ch 16 | 用互斥鎖防止超賣，背景 worker 送通知 |
| Ch 17 | 注入版本號、跨平台編譯、撰寫文件 |
| Ch 18 | 後台登入（bcrypt、HMAC）與 HTTPS |
| Ch 19 | 用反射做一個欄位驗證器 |

<!--
後半段的路線，GoShop 會從一個命令列小程式，長成一個真正的網站服務。

第 13 章換成 MySQL，第 15 章變成網站，第 16 章處理「很多人同時搶購」的問題，第 18 章加上登入和 HTTPS。這些都是實際工作中一定會碰到的需求。

每一步的完整參考解答，都放在課程 repo 的 goshop 資料夾，例如 goshop/ch05 就是做完第 5 章之後的樣子。大家先自己試著寫，卡住了再去對照。
-->

---

# GoShop 第 0 步：建立專案
### 任務說明

1. 建立一個資料夾 `goshop`，在裡面用 `go mod init goshop` 初始化模組
2. 建立 `main.go`，執行後印出歡迎訊息
3. 用 `go run .` 執行，確認看到下面的結果
4. 用 VS Code 開啟 `goshop` 資料夾，之後每一章都在這個專案上修改

```text
$ go run .
歡迎光臨 GoShop！
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>參考解答：</b> 課程 repo 的 <code>goshop/ch00</code> ～ <code>goshop/ch19</code> 是每一章做完後的完整程式碼，卡住的時候可以對照。
</div>

<!--
第 0 步很簡單，就是用今天學到的指令，把專案建立起來。

這裡要注意模組名稱就叫 goshop，之後第 8 章拆套件的時候，import 路徑會是 goshop/internal/什麼什麼，名稱一致比較不會搞混。

建立好之後，用 VS Code 打開整個 goshop 資料夾，而不是只打開 main.go，這樣 Go 擴充套件才能正確認得這個模組。
-->

---

# GoShop 第 0 步：解題提示
### 提示說明

```bash
mkdir goshop
cd goshop
go mod init goshop    # 產生 go.mod，模組名稱是 goshop
```

```go
// goshop/main.go
package main

import "fmt"

func main() {
	fmt.Println("歡迎光臨 GoShop！")
}
```

- 程式碼區塊第一行的 `// goshop/main.go` 表示這段程式碼在專案裡的**檔案路徑**

<!--
先用 mkdir 建立資料夾，cd 進去，再用 go mod init goshop 產生 go.mod 檔。

main.go 的內容跟今天的 Hello World 幾乎一樣，只是把訊息換成歡迎光臨 GoShop。

之後每一章的解題提示，程式碼第一行都會用註解標出檔案路徑，例如 goshop/main.go、goshop/internal/store/memory.go，這樣大家就知道這段程式碼要放在專案的哪個檔案裡。

執行 go run . 之後，會印出「歡迎光臨 GoShop！」。恭喜大家，GoShop 正式開張了！
-->

---

# 章節總結

- **安裝 Go**：到 go.dev/dl 下載，安裝後**重開終端機**，用 `go version` 確認
- **VS Code**：安裝官方 Go 擴充套件（`golang.go`），並安裝 `gopls`、`dlv` 等工具
- **建立專案**：`go mod init <模組路徑>` 產生 `go.mod`；程式從 `package main` 的 `func main()` 開始
- **執行程式**：開發用 `go run .`，部署用 `go build` 產生單一執行檔
- **Go Playground**：go.dev/play，用 Share 分享程式碼；不能連網、時間固定
- **GoShop**：建立 `goshop` 專案，接下來每一章都會替它加上新功能

下一章我們會正式進入 Go 語法，從「變數與算符」開始。

<!--
我們來整理一下今天的重點。

第一，安裝 Go 之後要重開終端機，用 go version 確認。第二，VS Code 要安裝官方的 Go 擴充套件，它會幫我們自動排版、自動 import。第三，每個專案從 go mod init 開始，程式的進入點是 main 套件裡的 main 函式。第四，開發時用 go run 點，要部署時用 go build。第五，Go Playground 很適合快速測試和分享程式碼。

另外，我們也建立了課程專案 GoShop。從下一章開始，每一章的最後都會回到這個專案，把剛學到的東西用上去，整門課上完，就會得到一個真正能運作的電商後台。

到這裡，大家的開發環境已經完全準備好了。下一章我們會正式進入 Go 的語法，從最基本的變數宣告、各種運算子開始學起。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
環境安裝是最容易卡關的地方，每個人的電腦狀況都不太一樣。

如果 go version 跑不出來、VS Code 的工具裝不起來，現在就可以提問，我們一起把環境弄好，下一章才能順利開始寫程式。
-->

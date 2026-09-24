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
title: 系統與檔案
routeAlias: ch12
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">系統與檔案</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「和作業系統打交道：參數、訊號與檔案」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十二章！

到目前為止，我們的程式都是「自己跑自己的」：資料寫死在程式碼裡，執行完結果印在螢幕上就消失了。真實的程式需要跟作業系統互動：從命令列接收參數、在使用者按下 Ctrl+C 的時候優雅地結束、把資料存進檔案、下次執行時再讀回來。

今天就來學這些跟系統打交道的技能。學完之後，大家就能用 Go 寫出真正實用的命令列工具，這也是 Go 最擅長的領域之一，像 Docker、kubectl、GitHub CLI 都是用 Go 寫的命令列工具。
-->

---
layout: default
---

# Outline

- **命令列旗標與其引數** — `os.Args`、`flag` 套件
- **系統中斷訊號** — `os/signal`、優雅關閉
- **檔案存取權限** — `rwx`、八進位權限 `0o644`
- **建立與寫入檔案** — 建立、寫入、檢查存在、讀取整個檔案、逐行讀取、刪除
- **os.OpenFile()** — 最完整的檔案開啟方式
- **處理 CSV 格式檔案** — `encoding/csv`
- **章節總結**

<!--
今天的內容分成三大塊。

第一塊是程式的「輸入」：命令列參數和系統訊號。第二塊是今天的重點：檔案操作，包括權限、建立、寫入、讀取、刪除。第三塊是處理 CSV 這種最常見的資料檔案格式。
-->

---

# 回顧：編碼／解碼 JSON 資料

- `json.Unmarshal(data, &v)` / `json.Marshal(v)`；struct 標籤 `json:"name,omitempty"`
- `json.NewDecoder(r)` / `json.NewEncoder(w)` 搭配任何 `io.Reader` / `io.Writer`
- 上一章的綜合練習用 `bytes.Buffer` **模擬檔案** ← 今天換成真正的檔案
- 第 7 章：`*os.File` 同時實作了 `io.Reader` 與 `io.Writer`

<!--
回顧一下上一章。

我們學了 JSON 的編碼與解碼，Decoder 和 Encoder 可以搭配任何 Reader 和 Writer。上一章的綜合練習用 bytes.Buffer 模擬檔案，今天就要換成真正的檔案。

因為 *os.File 同時實作了 io.Reader 和 io.Writer，所以上一章寫的 Save 和 Load，今天可以直接傳入檔案使用，一行都不用改。這就是第 7 章「接受介面」的威力。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 命令列旗標與其引數
## Command-Line Flags

<!--
先來看怎麼從命令列接收參數。
-->

---
zoom: 0.93
---

# 最基本的方式：os.Args

`os.Args` 是一個字串切片：`os.Args[0]` 是程式本身，之後是使用者輸入的引數

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("程式：", os.Args[0])
	fmt.Println("引數：", os.Args[1:])
	if len(os.Args) < 2 {
		fmt.Println("用法：greet <名字>")
		os.Exit(1) // 非 0 的結束碼代表失敗
	}
	fmt.Println("Hello,", os.Args[1])
}
```

```bash
go run . Alice Bob
# 程式： /tmp/go-build.../exe/greet
# 引數： [Alice Bob]
# Hello, Alice
```

<!--
最基本的方式是 os.Args，它是一個字串切片。第 0 個元素是程式本身的路徑，後面才是使用者輸入的引數，以空白分隔。

如果使用者沒有輸入引數，我們印出用法說明，並用 os.Exit(1) 結束程式。結束碼 0 代表成功，非 0 代表失敗，這是作業系統的慣例，其他程式或腳本可以根據結束碼判斷我們的程式有沒有執行成功。

注意 os.Exit 跟 log.Fatal 一樣，會立刻結束程式，defer 不會執行。

os.Args 適合很簡單的情況；如果有很多選項，就要用 flag 套件。
-->

---
zoom: 0.91
---

# 使用 flag 套件定義旗標

```go
package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// 回傳 *int
	port := flag.Int("port", 8080, "伺服器的連接埠")
	// 回傳 *string
	env := flag.String("env", "dev", "執行環境：dev / prod")
	verbose := flag.Bool("v", false, "顯示詳細訊息")
	timeout := flag.Duration("timeout", 5*time.Second, "逾時時間")
	flag.Parse() // 一定要呼叫，才會解析命令列

	fmt.Println(*port, *env, *verbose, *timeout)
	fmt.Println("其他引數：", flag.Args()) // 旗標之後、非旗標的引數
}
```

```bash
go run . -port=9000 -env prod -v -timeout 30s a.txt b.txt
# 9000 prod true 30s
# 其他引數： [a.txt b.txt]
```

<!--
flag 套件是 Go 標準函式庫提供的命令列解析工具。

每一種型別都有對應的函式：flag.Int、flag.String、flag.Bool、flag.Duration。三個參數分別是旗標名稱、預設值、說明文字。它們回傳的是指標，所以使用的時候要加星號解參考。

定義完所有旗標之後，一定要呼叫 flag.Parse()，它才會真正去解析命令列。忘記呼叫的話，所有旗標都會是預設值。

旗標之後、不是旗標的引數，可以用 flag.Args() 取得，例如要處理的檔案名稱。

命令列上旗標的寫法很彈性：-port=9000 或 -port 9000 都可以，布林旗標只要寫 -v 就代表 true。
-->

---
zoom: 0.76
---

# 使用 flag 的注意事項

**注意事項之一：** `flag` 會自動產生 `-h` / `-help` 說明，列出所有旗標與預設值

**注意事項之二：** 想把值存進**既有的變數**（例如 struct 欄位），用 `XxxVar` 版本

```go
package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Host string
	Port int
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Host, "host", "localhost", "主機名稱")
	flag.IntVar(&cfg.Port, "port", 8080, "連接埠")
	flag.Parse()
	fmt.Printf("%+v\n", cfg)
}
```

```text
$ go run . -h
Usage of /tmp/go-build.../exe/app:
  -host string
    	主機名稱 (default "localhost")
  -port int
    	連接埠 (default 8080)
```

<!--
使用 flag 有兩個注意事項。

第一，flag 會自動提供 -h 說明，列出所有旗標的名稱、型別、說明和預設值，所以說明文字要寫清楚，使用者才知道每個旗標是做什麼的。

第二，如果想要把旗標的值直接存進某個變數，例如設定結構的欄位，就用 StringVar、IntVar 這些 Var 結尾的版本，第一個參數傳入變數的指標。這樣就不需要到處解參考，程式碼比較乾淨。

補充：大型的命令列工具，例如有子命令的 git commit、docker run，社群常用第三方套件 cobra，kubectl 和 GitHub CLI 都是用它寫的。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 系統中斷訊號
## OS Signals

<!--
接下來看系統訊號：當使用者按下 Ctrl+C，或是伺服器要關機的時候，程式怎麼知道？
-->

---

# 什麼是系統訊號？

「**訊號（signal）是作業系統通知程式『發生了某件事』的機制。**」

| 訊號 | 觸發方式 | 預設行為 | Go 中的常數 |
| --- | --- | --- | --- |
| `SIGINT` | 使用者按下 **Ctrl+C** | 結束程式 | `os.Interrupt` |
| `SIGTERM` | `kill <pid>`、Docker / Kubernetes 停止容器 | 結束程式 | `syscall.SIGTERM` |
| `SIGKILL` | `kill -9 <pid>` | **強制結束，無法攔截** | `os.Kill` |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>為什麼要攔截？</b> 程式被直接結束時，寫到一半的檔案可能會損壞、處理中的請求會中斷。攔截訊號後，我們可以先<b>收尾（graceful shutdown）</b>再結束。
</div>

<!--
什麼是系統訊號？

它是作業系統通知程式的一種方式，就像手機的推播通知。最常見的是 SIGINT，使用者在終端機按下 Ctrl+C 時送出；以及 SIGTERM，Docker、Kubernetes 在停止容器時會送出這個訊號。

這兩種訊號的預設行為都是「立刻結束程式」。問題是：如果程式正在寫檔案、正在處理使用者的付款請求，直接結束就可能造成資料損壞。

所以我們要攔截這些訊號，先把手上的工作做完、把資源關閉，再優雅地結束，這叫做 graceful shutdown，優雅關閉。

注意 SIGKILL 是無法攔截的，kill -9 就是「不管三七二十一，直接強制結束」。
-->

---

# 攔截訊號：signal.Notify

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigCh := make(chan os.Signal, 1) // 通道：用來接收訊號（Ch 16）
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	fmt.Println("程式執行中，按 Ctrl+C 結束…")
	sig := <-sigCh // 在這裡等待，直到收到訊號
	fmt.Println("\n收到訊號：", sig)

	fmt.Println("儲存資料、關閉連線中…")
	time.Sleep(500 * time.Millisecond) // 模擬收尾工作
	fmt.Println("已安全結束")
}
```

<!--
攔截訊號用 os/signal 套件的 Notify 函式。

這裡用到了一個還沒正式學過的東西：通道（channel），第 16 章會詳細介紹。現在先把它想成一個「信箱」：make 建立一個信箱，signal.Notify 告訴 Go「收到 SIGINT 或 SIGTERM 的時候，把訊號投進這個信箱」。

<-sigCh 這個寫法的意思是「從信箱取出一封信」，如果信箱是空的，程式就會停在這裡等待。所以程式會一直等，直到使用者按下 Ctrl+C，訊號被投進信箱，程式才繼續往下執行收尾工作。

執行後按 Ctrl+C，會印出收到訊號：interrupt，然後做完收尾才結束。
-->

---
zoom: 0.78
---

# 現代寫法：signal.NotifyContext

搭配 `context`（Ch 16），訊號到來時**取消 context**，是伺服器最主流的寫法

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(time.Second) // 每秒觸發一次
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("工作中…")
		case <-ctx.Done(): // 收到訊號時，ctx 會被取消
			fmt.Println("收到結束訊號：", context.Cause(ctx))
			return
		}
	}
}
```

<!--
Go 1.16 加入了 signal.NotifyContext，這是現在最主流的寫法，第 15 章的 HTTP 伺服器優雅關閉也會用它。

它回傳一個 context，當收到指定的訊號時，這個 context 就會被「取消」。context 是第 16 章的主題，現在先理解為「一個可以被取消的通知」。

程式用一個 for 迴圈持續工作，每秒印出一次「工作中」。select 會同時等待兩件事：計時器響了，或是 context 被取消了，哪一個先發生就執行哪一個。按下 Ctrl+C 後，ctx.Done() 會收到通知，程式印出原因後 return 結束。

select、ticker、context 都是第 16 章的內容，這裡先看懂整體的流程就好。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 檔案存取權限
## File Permissions

<!--
在建立檔案之前，我們要先了解檔案的權限。
-->

---

# 檔案存取權限：rwx 與八進位

Unix 系統的權限分為三組對象，每組有**讀（r=4）、寫（w=2）、執行（x=1）**三種權限

| 八進位 | 符號 | 擁有者 | 群組 | 其他人 | 常見用途 |
| --- | --- | --- | --- | --- | --- |
| `0o644` | `-rw-r--r--` | 讀寫 | 讀 | 讀 | 一般檔案 |
| `0o600` | `-rw-------` | 讀寫 | — | — | 私密檔案（金鑰、密碼） |
| `0o755` | `-rwxr-xr-x` | 讀寫執行 | 讀執行 | 讀執行 | 執行檔、資料夾 |
| `0o700` | `-rwx------` | 讀寫執行 | — | — | 私人資料夾 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>計算方式：</b> 6 = 4（讀）+ 2（寫）；7 = 4 + 2 + 1。Go 用 <code>0o</code> 前綴表示八進位（第 3 章），型別是 <code>fs.FileMode</code>。Windows 只會參考「唯讀」與否。
</div>

<!--
Linux 和 macOS 的檔案權限，分成三組對象：擁有者、同群組的人、其他所有人。每一組有三種權限：讀、寫、執行，分別用 4、2、1 代表，加起來就是這一組的權限。

所以 644 的意思是：擁有者 6（讀加寫）、群組 4（讀）、其他人 4（讀）。這是一般檔案最常用的權限。

存放密碼、私鑰的檔案要用 600，只有自己能讀寫，第 18 章產生的私鑰檔案就要用這個權限。執行檔和資料夾通常是 755，資料夾需要「執行」權限才能進入。

在 Go 裡用 0o 開頭表示八進位，型別是 fs.FileMode。Windows 的權限系統不一樣，Go 在 Windows 上只會參考是否唯讀。
-->

---
zoom: 0.94
---

# 查詢與修改權限

```go
package main

import (
	"fmt"
	"io/fs"
	"os"
)

func main() {
	os.WriteFile("secret.txt", []byte("api-key"), 0o644)

	info, err := os.Stat("secret.txt") // 取得檔案資訊
	if err != nil {
		panic(err)
	}
	// secret.txt 7 -rw-r--r--
	fmt.Println(info.Name(), info.Size(), info.Mode())

	os.Chmod("secret.txt", 0o600) // 修改權限
	info, _ = os.Stat("secret.txt")
	// -rw------- -rwxr-xr-x
	fmt.Println(info.Mode(), fs.FileMode(0o755))
	os.Remove("secret.txt")
}
```

<!--
os.Stat 可以取得檔案的資訊，包括名稱、大小、權限、修改時間、是不是資料夾。Mode 方法回傳權限，印出來就是 ls -l 看到的那種符號格式。

os.Chmod 可以修改權限。

另外提醒：建立檔案時指定的權限，會再被系統的 umask 設定過濾，通常會去掉群組和其他人的寫入權限，所以指定 0o666 實際上可能變成 0o644。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 建立與寫入檔案
## Creating & Writing Files

<!--
接下來是今天的重點：檔案的建立、寫入、讀取和刪除。
-->

---
zoom: 0.95
---

# 用 os 套件新建檔案

`os.Create(name)`：建立檔案（**已存在會被清空**），回傳 `*os.File`

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	// 權限 0o666（經 umask 後通常是 0o644）
	f, err := os.Create("hello.txt")
	if err != nil {
		fmt.Println("建立失敗：", err)
		return
	}
	defer f.Close() // 取得資源 → 檢查錯誤 → 馬上 defer 釋放

	fmt.Println("已建立：", f.Name())
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>一定要 Close：</b> 作業系統能同時開啟的檔案數量有限，忘記關閉會造成「檔案描述符耗盡（too many open files）」。
</div>

<!--
os.Create 會建立一個新檔案，回傳 *os.File 和錯誤。如果檔案已經存在，內容會被清空，這一點要特別小心。

拿到檔案之後，第一件事就是 defer f.Close()。這就是第 5 章說的 defer 慣用法：取得資源、檢查錯誤、馬上 defer 釋放，三行寫在一起。

為什麼一定要關閉？作業系統對每個程式能同時開啟的檔案數量有上限，通常是 1024 個。如果在迴圈裡開檔案卻忘了關，很快就會遇到 too many open files 的錯誤。
-->

---
zoom: 0.96
---

# 對檔案寫入字串

`*os.File` 實作了 `io.Writer`，可以用各種方式寫入

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("report.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("銷售報表\n")         // 寫入字串
	f.Write([]byte("======\n"))     // 寫入位元組切片
	fmt.Fprintf(f, "咖啡：%d 杯\n", 42) // 格式化寫入（io.Writer）
	if _, err := fmt.Fprintln(f, "合計：2520 元"); err != nil {
		fmt.Println("寫入失敗：", err)
	}
}
```

<!--
*os.File 實作了 io.Writer，所以有很多種寫入方式。

WriteString 寫入字串，Write 寫入位元組切片。最方便的是 fmt.Fprintf，第 9 章學過，F 開頭的函式可以寫到任何 io.Writer，當然也包括檔案。

每一個寫入方法都會回傳寫入的位元組數和錯誤。實務上寫入檔案失敗的情況不多，例如磁碟滿了，但重要的資料還是要檢查錯誤。
-->

---

# 一次完成建立檔案及寫入：os.WriteFile

檔案不大時，`os.WriteFile` 一行完成**建立、寫入、關閉**

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	data := []byte("host=localhost\nport=8080\n")
	if err := os.WriteFile("app.conf", data, 0o644); err != nil {
		fmt.Println("寫入失敗：", err)
		return
	}
	fmt.Println("設定檔已寫入")
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>舊寫法：</b> 網路上常見的 <code>ioutil.WriteFile</code>、<code>ioutil.ReadFile</code> 從 Go 1.16 起已<b>棄用</b>，請改用 <code>os.WriteFile</code>、<code>os.ReadFile</code>。
</div>

<!--
如果檔案內容不大，可以一次準備好，os.WriteFile 一行就能完成建立、寫入、關閉三件事，不用自己處理 Close。

三個參數是：檔名、內容的位元組切片、權限。檔案已存在時會被覆蓋。

提醒一下：網路上很多舊文章用的是 ioutil 套件，它從 Go 1.16 開始已經棄用了，功能都搬到 os 和 io 套件裡。看到 ioutil.ReadFile 就換成 os.ReadFile。
-->

---
zoom: 0.88
---

# 檢查檔案是否存在

用 `os.Stat` 取得資訊，再用 `errors.Is(err, fs.ErrNotExist)` 判斷

```go
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) { // 哨兵錯誤（Ch 6）
		return false, nil
	}
	return false, err // 其他錯誤，例如沒有權限
}

func main() {
	fmt.Println(exists("app.conf"))
	fmt.Println(exists("no-such-file.txt")) // false <nil>
}
```

<!--
Go 沒有直接的「檔案存在嗎」函式，標準做法是呼叫 os.Stat，看它回傳什麼錯誤。

沒有錯誤，代表檔案存在。錯誤是 fs.ErrNotExist，代表檔案不存在，這是第 6 章學的哨兵錯誤，用 errors.Is 判斷。其他錯誤，例如沒有權限讀取資料夾，就代表「不確定」，應該回傳錯誤讓呼叫端處理。

使用的注意事項：「先檢查存在、再開檔案」這種寫法有一個小問題：檢查完到開檔案之間，檔案可能被別的程式刪掉了。所以很多時候更好的做法是直接開檔案，再處理 ErrNotExist 錯誤。
-->

---
zoom: 0.96
---

# 一次讀取整個檔案內容：os.ReadFile

```go
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	data, err := os.ReadFile("app.conf") // 回傳 []byte
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Println("找不到設定檔，使用預設值")
		return
	} else if err != nil {
		fmt.Println("讀取失敗：", err)
		return
	}
	fmt.Printf("共 %d bytes\n%s", len(data), data)
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>注意檔案大小：</b> <code>ReadFile</code> 會把整個檔案讀進記憶體。幾 GB 的日誌檔請改用<b>逐行讀取</b>。
</div>

<!--
os.ReadFile 一次把整個檔案讀進記憶體，回傳 []byte。要當成字串使用就用 string() 轉換，或用 %s 印出來。

這裡示範了直接讀檔、再處理 ErrNotExist 的寫法：找不到設定檔就用預設值，其他錯誤才真正回報。

使用 ReadFile 的注意事項：它會把整個檔案讀進記憶體，設定檔、小型 JSON 檔沒問題，但如果是好幾 GB 的日誌檔，記憶體就爆了。這時候要用下一頁的逐行讀取。
-->

---
zoom: 0.88
---

# 一次讀取檔案中的一行字串：bufio.Scanner

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("app.conf") // 唯讀開啟
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f) // 包裝成逐行讀取的 Scanner
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text()) // 不含換行字元
		fmt.Printf("%d: %s\n", n, line)
	}
	if err := sc.Err(); err != nil { // 迴圈結束後檢查錯誤
		fmt.Println("讀取錯誤：", err)
	}
}
```

<!--
逐行讀取用 bufio.Scanner。os.Open 以唯讀方式開啟檔案，然後用 bufio.NewScanner 包裝。

sc.Scan() 每次讀一行，讀到就回傳 true，讀到檔案結尾或發生錯誤就回傳 false，所以很適合放在 for 迴圈的條件。sc.Text() 取出這一行的內容，不包含換行字元。

使用 Scanner 的注意事項：迴圈結束後，一定要呼叫 sc.Err() 檢查是「讀完了」還是「出錯了」。另外 Scanner 預設每行最多 64KB，超過會出錯，遇到超長的行可以用 sc.Buffer 調整。

逐行讀取每次只把一行放在記憶體裡，所以不管檔案多大都沒問題。
-->

---
zoom: 0.93
---

# 刪除檔案

| 函式 | 說明 |
| --- | --- |
| `os.Remove(path)` | 刪除**一個檔案或空資料夾**；不存在時回傳錯誤 |
| `os.RemoveAll(path)` | 刪除資料夾及**底下所有內容**；不存在時**不會**回傳錯誤 |
| `os.Rename(old, new)` | 重新命名或搬移 |
| `os.MkdirAll(path, 0o755)` | 建立多層資料夾 |

```go
err := os.Remove("app.conf")
if err != nil && !errors.Is(err, fs.ErrNotExist) {
	fmt.Println("刪除失敗：", err)
}
// data/2026/09
os.MkdirAll(filepath.Join("data", "2026", "09"), 0o755)
defer os.RemoveAll("data")
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>組合路徑請用 <code>filepath.Join</code>：</b> 它會依作業系統使用正確的分隔符號（Windows 是 <code>\</code>、Linux / macOS 是 <code>/</code>）。
</div>

<!--
刪除檔案用 os.Remove，它只能刪除檔案或「空的」資料夾。os.RemoveAll 可以刪除整個資料夾和底下的所有東西，要非常小心使用。

這張表也列出了兩個常用的函式：Rename 可以重新命名或搬移檔案；MkdirAll 可以一次建立多層資料夾，已經存在也不會出錯。

組合路徑的時候，請用 filepath.Join，不要自己用加號串接斜線。因為 Windows 的路徑分隔符號是反斜線，Linux 和 macOS 是斜線，filepath.Join 會自動處理。
-->

---
layout: default
---

# 練習 1：簡易記事本
### 任務說明

寫一個命令列記事本 `note`：

1. 旗標 `-file`（預設 `notes.txt`）指定檔案；旗標 `-list` 列出所有記事
2. 其他引數（`flag.Args()`）合併成一行記事，**附加**到檔案結尾，前面加上時間：
   `2026-09-24 14:05 | 買咖啡豆`
3. `-list` 時逐行讀出，加上編號印出；檔案不存在時印出「目前沒有記事」
4. 提示：附加寫入可以先用 `os.ReadFile` 讀出舊內容再 `os.WriteFile`（下一節會學更好的方法）

```bash
go run . 買咖啡豆
go run . -list
```

<!--
這個練習結合了今天前半段的內容：flag 旗標、寫入檔案、讀取檔案、檢查檔案是否存在。

附加寫入的部分，這裡先用「讀出舊內容、加上新內容、整個寫回去」的方法，下一節學了 os.OpenFile 之後，會有更好的做法。
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

func main() {
	file := flag.String("file", "notes.txt", "記事檔案")
	list := flag.Bool("list", false, "列出所有記事")
	flag.Parse()
```

<!--
先定義兩個旗標並解析：-file 指定檔案，預設是 notes.txt；-list 是布林旗標，有寫就代表要列出記事。
-->

---

# 練習 1：解題提示（續）
### 提示說明

```go
	// 續上頁
	old, err := os.ReadFile(*file)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		fmt.Println("讀取失敗：", err)
		os.Exit(1)
	}
	if *list {
		if len(old) == 0 {
			fmt.Println("目前沒有記事")
		}
		lines := strings.Split(strings.TrimSpace(string(old)), "\n")
		for i, l := range lines {
			if l != "" {
				fmt.Printf("%d. %s\n", i+1, l)
			}
		}
		return
	}
```

<!--
接著讀出舊的內容，如果錯誤是「檔案不存在」就當作空的，其他錯誤才結束程式。

-list 模式下，把內容切成一行一行加上編號印出。檔案不存在時 old 是空的，印出「目前沒有記事」。
-->

---

# 練習 1：解題提示（續 2）
### 提示說明

```go
	// 續上頁
	text := strings.Join(flag.Args(), " ")
	if text == "" {
		fmt.Println("用法：note [-file 檔名] <記事內容>")
		os.Exit(1)
	}
	line := time.Now().Format("2006-01-02 15:04") + " | " + text + "\n"
	data := append(old, line...)
	if err := os.WriteFile(*file, data, 0o644); err != nil {
		fmt.Println("寫入失敗：", err)
		os.Exit(1)
	}
	fmt.Println("已新增：", line)
}
```

<!--
新增記事的模式：把 flag.Args() 用空白串接成一行文字，沒有內容就印出用法。

接著在前面加上時間，用第 4 章學的 append 加上三個點，把字串附加到舊內容的 []byte 後面，再整個寫回檔案。

這個方法的缺點是：檔案越大，每次都要整個讀出來再寫回去，效率越差。下一節的 os.OpenFile 可以直接在檔案結尾附加。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 最完整的檔案開啟與建立功能
## os.OpenFile()

<!--
接下來看最完整的檔案開啟方式：os.OpenFile。
-->

---
zoom: 0.84
---

# os.OpenFile()：用旗標控制開啟方式

語法：`os.OpenFile(name, 旗標, 權限)`，旗標用 `|` 組合

| 旗標 | 意義 |
| --- | --- |
| `os.O_RDONLY` / `os.O_WRONLY` / `os.O_RDWR` | 唯讀／唯寫／讀寫（三選一） |
| `os.O_CREATE` | 檔案不存在時建立 |
| `os.O_APPEND` | 寫入時**附加到檔案結尾** |
| `os.O_TRUNC` | 開啟時**清空**檔案內容 |
| `os.O_EXCL` | 搭配 `O_CREATE`：檔案**已存在就回傳錯誤** |

| 常見組合 | 等同於 |
| --- | --- |
| `O_RDONLY` | `os.Open(name)` |
| `O_RDWR \| O_CREATE \| O_TRUNC`，`0o666` | `os.Create(name)` |
| `O_WRONLY \| O_CREATE \| O_APPEND`，`0o644` | **附加寫入日誌** |

<!--
os.Open 和 os.Create 其實都是 os.OpenFile 的簡化版。os.OpenFile 可以用旗標精確控制開啟的方式。

旗標分兩類：第一類是讀寫模式，唯讀、唯寫、讀寫三選一；第二類是額外的行為，可以用直線符號組合多個。

最常用的組合是 O_WRONLY、O_CREATE、O_APPEND：以唯寫模式開啟，不存在就建立，寫入時附加到結尾。這就是寫日誌檔的標準做法，第 9 章建立 logger 的時候已經用過了。

O_EXCL 很適合用在「不能覆蓋既有檔案」的情況，例如建立鎖定檔。
-->

---
zoom: 0.86
---

# os.OpenFile() — 附加寫入範例

```go
package main

import (
	"fmt"
	"os"
	"time"
)

func appendLog(path, msg string) error {
	flag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	f, err := os.OpenFile(path, flag, 0o644)
	if err != nil {
		return fmt.Errorf("開啟日誌失敗：%w", err)
	}
	defer f.Close()
	now := time.Now().Format(time.DateTime)
	_, err = fmt.Fprintf(f, "%s %s\n", now, msg)
	return err
}

func main() {
	for _, m := range []string{"服務啟動", "處理訂單 #7"} {
		if err := appendLog("app.log", m); err != nil {
			fmt.Println(err)
		}
	}
}
```

<!--
這段程式碼的目的，是寫一個可以重複使用的附加日誌函式。

每次呼叫 appendLog，都用附加模式開啟檔案、寫入一行、關閉。不管呼叫幾次，內容都會累加在檔案結尾，不會覆蓋之前的內容，也不需要把舊內容讀出來。

錯誤處理用了第 6 章的 %w 包裝，加上「開啟日誌失敗」的情境。

練習 1 的記事本，就可以改用這個方式附加記事。
-->

---

# 補充：os.Root 安全地存取資料夾（Go 1.24+）

處理**使用者提供的檔名**時，防止 `../` 跳出指定資料夾（路徑穿越攻擊）

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	os.MkdirAll("uploads", 0o755)
	root, err := os.OpenRoot("uploads") // 只能存取 uploads 資料夾
	if err != nil {
		fmt.Println(err)
		return
	}
	defer root.Close()

	_, err = root.Open("../etc/passwd") // 使用者惡意輸入的檔名
	// openat ../etc/passwd: path escapes from parent
	fmt.Println(err)
}
```

<!--
這是補充內容，跟資訊安全有關。

如果我們的程式會根據「使用者提供的檔名」開啟檔案，例如下載使用者上傳的檔案，就要小心「路徑穿越攻擊」：惡意使用者輸入 ../../etc/passwd，就可能讀到系統的敏感檔案。

Go 1.24 加入了 os.Root：OpenRoot 開啟一個資料夾當作「根」，之後透過 root 開啟的檔案，都只能在這個資料夾裡面，任何想跳出去的路徑都會被拒絕。

第 15 章的 HTTP 伺服器處理檔案上傳、下載時，這是很重要的安全措施。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 處理 CSV 格式檔案
## encoding/csv

<!--
最後來看 CSV 檔案的處理。
-->

---

# 什麼是 CSV？

「**CSV（Comma-Separated Values）以逗號分隔欄位、以換行分隔紀錄**，Excel、Google 試算表、資料庫都能匯入匯出。」

```text
name,category,price
咖啡豆,飲品,450
"蛋糕, 巧克力",甜點,120
綠茶,飲品,35
```

| 規則 | 說明 |
| --- | --- |
| 第一行 | 通常是**標題列**（欄位名稱） |
| 欄位內含逗號或換行 | 用**雙引號**包起來：`"蛋糕, 巧克力"` |
| 欄位內含雙引號 | 寫成兩個雙引號：`"他說""好"""` |

<!--
什麼是 CSV？它是一種最簡單的表格資料格式：每一行是一筆紀錄，欄位之間用逗號分隔。Excel、Google 試算表、資料庫都能匯入和匯出 CSV，所以它是交換表格資料最通用的格式。

CSV 看起來很簡單，用 strings.Split 以逗號切開好像就好了？其實不行：如果欄位本身就包含逗號，例如「蛋糕, 巧克力」，就要用雙引號包起來。自己處理這些規則很容易出錯，所以要用 Go 標準函式庫的 encoding/csv 套件。
-->

---
zoom: 0.84
---

# 走訪 CSV 檔內容

```go
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

func main() {
	data := `name,category,price
咖啡豆,飲品,450
"蛋糕, 巧克力",甜點,120
`
	r := csv.NewReader(strings.NewReader(data)) // 也可以傳入 *os.File
	for {
		record, err := r.Read() // 每次讀一筆，回傳 []string
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("格式錯誤：", err)
			return
		}
		fmt.Println(len(record), record)
	}
}
```


<!--
csv.NewReader 接收一個 io.Reader，實務上傳入 os.Open 開啟的檔案，這裡為了方便示範，用 strings.NewReader 把字串包裝成 Reader。

r.Read() 每次讀一筆紀錄，回傳一個字串切片，每個元素是一個欄位。讀到結尾時回傳 io.EOF。這個「用迴圈讀到 EOF」的模式，跟第 11 章 JSON Decoder 讀多筆資料一模一樣。

注意第三筆資料：「蛋糕, 巧克力」被雙引號包著，csv 套件正確地把它當成一個欄位。執行後三筆紀錄都印出 3 個欄位。
-->

---
zoom: 0.88
---

# 讀取每行資料各欄位的值

CSV 的欄位都是字串，需要用 `strconv` 轉型；用**標題列**找出欄位位置

```go
package main

import (
	"encoding/csv"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func main() {
	data := `name,category,price
咖啡豆,飲品,450
綠茶,飲品,35
蛋糕,甜點,x
`
	// 一次讀完
	rows, err := csv.NewReader(strings.NewReader(data)).ReadAll()
	if err != nil {
		panic(err)
	}
	header := rows[0]
	nameIdx := slices.Index(header, "name")
	priceIdx := slices.Index(header, "price")
```

<!--
CSV 讀出來的每個欄位都是字串，數字要用第 3 章學的 strconv 轉換。

ReadAll 一次把所有紀錄讀成一個二維的字串切片，適合檔案不大的情況。

實務上的一個好習慣：不要把欄位的位置寫死成 row[2]，而是從標題列用 slices.Index 找出「price 在第幾欄」。這樣就算有人在 Excel 裡調整了欄位順序，程式也不會出錯。
-->

---

# 讀取每行資料各欄位的值（續）

```go
	// 續上頁
	total := 0
	for i, row := range rows[1:] {
		price, err := strconv.Atoi(row[priceIdx])
		if err != nil {
			fmt.Printf("第 %d 筆 %s 價格錯誤：%v\n", i+2, row[nameIdx], err)
			continue
		}
		total += price
	}
	fmt.Println("總價：", total) // 總價： 485
}
```

<!--
走訪標題列之後的每一筆資料，把價格轉成整數加總。

價格轉換失敗的資料，印出第幾筆、哪個商品有問題，然後跳過，繼續處理其他資料。i+2 是因為要跳過標題列，而且行號從 1 開始算。
-->

---
zoom: 0.85
---

# 寫入 CSV：csv.Writer

```go
package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	f, err := os.Create("orders.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{"id", "item", "note"})
	w.Write([]string{"1", "咖啡", "少冰, 半糖"}) // 自動加上雙引號
	w.Flush()                              // 把緩衝區的資料寫入檔案
	if err := w.Error(); err != nil {
		fmt.Println("寫入失敗：", err)
	}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>給 Excel 開啟的 CSV：</b> Excel 預設不一定用 UTF-8 解讀，中文可能變亂碼；可以在檔案開頭寫入 BOM：<code>f.WriteString("﻿")</code>。
</div>

<!--
寫入 CSV 用 csv.NewWriter，傳入 io.Writer。每次呼叫 Write 寫入一筆紀錄，欄位裡有逗號的話，csv 套件會自動加上雙引號。

使用 csv.Writer 的注意事項：它內部有緩衝區，Write 不會馬上寫進檔案，最後一定要呼叫 Flush 把緩衝區的資料寫出去，再用 Error 檢查有沒有發生錯誤。忘記 Flush 的話，檔案可能是空的。

另外一個實務小技巧：如果產生的 CSV 要給 Excel 開啟，中文常常會變成亂碼，因為 Excel 預設不一定用 UTF-8 解讀。在檔案開頭寫入一個 BOM 字元，Excel 就會正確辨識。
-->

---
layout: default
---

# 綜合練習：銷售報表產生器
### 任務說明

寫一個命令列工具，讀取銷售紀錄 CSV，輸出分類統計的 JSON 報表：

```text
date,category,item,amount
2026-09-01,飲品,咖啡,120
2026-09-01,甜點,蛋糕,240
2026-09-02,飲品,綠茶,35
2026-09-02,飲品,咖啡,abc
```

1. 旗標：`-in`（輸入 CSV，必填）、`-out`（輸出 JSON，預設 `report.json`）
2. 用標題列找出 `category`、`amount` 的欄位位置；金額錯誤的資料印出警告並跳過
3. 統計每個分類的金額總和，用 `json.MarshalIndent` 寫入 `-out` 檔案
4. 程式開始時先用 `os.WriteFile` 建立範例 `sales.csv`，方便測試

<!--
這個綜合練習把今天學的東西全部串起來：flag 旗標、讀取 CSV、檔案寫入，再加上上一章的 JSON 編碼。

這種「讀取一種格式、轉換成另一種格式」的小工具，是 Go 在實務上非常常見的用途。
-->

---
zoom: 0.94
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func summarize(path string) (map[string]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) == 0 {
		return nil, fmt.Errorf("讀取 CSV 失敗：%w", err)
	}
```

<!--
summarize 負責讀取 CSV 並統計：開檔案、defer 關閉、ReadAll 讀出所有資料。資料是空的或讀取失敗，就回傳錯誤。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
	// 續上頁
	catIdx := slices.Index(rows[0], "category")
	amtIdx := slices.Index(rows[0], "amount")
	sum := map[string]int{}
	for i, row := range rows[1:] {
		amt, err := strconv.Atoi(row[amtIdx])
		if err != nil {
			fmt.Printf("警告：第 %d 行金額錯誤，已略過\n", i+2)
			continue
		}
		sum[row[catIdx]] += amt
	}
	return sum, nil
}
```

<!--
用標題列找出欄位位置，走訪每一筆資料累加到 map 裡，金額錯誤的資料印出警告並略過。

注意 map 的 += 寫法：鍵不存在時讀到零值 0，加上金額後寫回去，這是第 4 章學的計數技巧。
-->

---
zoom: 0.88
---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func main() {
	in := flag.String("in", "", "輸入的 CSV 檔（必填）")
	out := flag.String("out", "report.json", "輸出的 JSON 檔")
	flag.Parse()
	if *in == "" {
		flag.Usage() // 印出旗標說明
		os.Exit(2)
	}
	sum, err := summarize(*in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := json.MarshalIndent(sum, "", "  ")
	if err := os.WriteFile(*out, b, 0o644); err != nil {
		fmt.Println("寫入報表失敗：", err)
		os.Exit(1)
	}
	fmt.Printf("已產生 %s：%s\n", *out, b)
}
```

```bash
go run . -in sales.csv    # {"甜點": 240, "飲品": 155}
```

<!--
main 負責處理旗標和輸出。-in 是必填的，沒有給的話呼叫 flag.Usage() 印出說明，並以結束碼 2 結束，這是命令列工具「用法錯誤」的慣例。

統計結果用 MarshalIndent 轉成 JSON，再用 WriteFile 寫入檔案。map 編碼成 JSON 時鍵會自動排序，所以輸出順序是固定的。

範例 CSV 的建立，可以在 main 的最前面用 os.WriteFile 寫入，或是自己用文字編輯器建立 sales.csv。咖啡 abc 那一筆會印出警告並略過，所以飲品的總和是 155。
-->

---

# 章節總結

- **命令列**：`os.Args` 最基本；`flag.Int` / `String` / `Bool` / `Duration` + `flag.Parse()`；`flag.Args()` 取得其他引數
- **系統訊號**：`signal.Notify` 攔截 `SIGINT` / `SIGTERM`；**`signal.NotifyContext`** 是現代主流寫法；目標是優雅關閉
- **權限**：`0o644` 一般檔案、`0o600` 私密檔案、`0o755` 執行檔與資料夾
- **寫入**：`os.Create` + `defer f.Close()`；小檔案用 `os.WriteFile`；附加用 `os.OpenFile(..., O_APPEND ...)`
- **讀取**：`os.ReadFile` 讀整個檔案；大檔案用 `bufio.Scanner` 逐行讀取，結束後檢查 `sc.Err()`
- **存在與刪除**：`os.Stat` + `errors.Is(err, fs.ErrNotExist)`；`os.Remove` / `RemoveAll`；路徑用 `filepath.Join`
- **CSV**：`csv.NewReader(r).Read()` 讀到 `io.EOF`；`csv.NewWriter(w)` 寫完要 `Flush()`

下一章我們會介紹「SQL 與資料庫」：用 Go 連接 MySQL，完成資料的新增、查詢、更新。

<!--
我們來整理今天學到的東西。

命令列的部分，flag 套件可以輕鬆處理各種旗標。系統訊號的部分，目標是優雅關閉，現代寫法是 signal.NotifyContext。檔案的部分，記得開檔案之後馬上 defer Close，小檔案用 ReadFile、WriteFile，大檔案用 Scanner 逐行讀取。CSV 用 encoding/csv 處理，寫入記得 Flush。

今天把資料存在檔案裡，但檔案有很多限制：多個程式同時寫入會衝突、要找某一筆資料得從頭讀到尾。所以真正的應用程式，資料都存在資料庫裡。下一章就要學怎麼用 Go 連接 MySQL 資料庫，用 SQL 新增、查詢、更新資料。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天學的檔案操作和命令列旗標，已經足夠寫出很多實用的小工具了。

課後建議：把第 11 章的待辦清單綜合練習，改成真正存到檔案裡，再加上 -add、-done、-list 三個旗標，做成一個完整的命令列待辦清單工具。

有問題的同學現在可以提問！
-->

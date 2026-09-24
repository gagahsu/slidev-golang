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
title: 並行性運算
routeAlias: ch16
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">並行性運算</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「不要透過共享記憶體來溝通，而要透過溝通來共享記憶體」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十六章！

前面幾章我們好幾次提到 goroutine：HTTP 伺服器會為每個請求啟動一個 goroutine、優雅關閉要在背景啟動伺服器、多個請求同時修改資料要用互斥鎖。今天終於要正式學習 Go 最有特色、也是 Go 被稱為「雲端時代的語言」的關鍵：並行性運算。

封面上這句話是 Go 的格言：「不要透過共享記憶體來溝通，而要透過溝通來共享記憶體」。今天學完 goroutine 和通道，大家就會理解這句話的意思。
-->

---
layout: default
---

# Outline

- **前言** — 並行（concurrency）與平行（parallelism）
- **Goroutine 與 WaitGroup** — 啟動 goroutine、等待它們完成
- **解決記憶體資源競爭** — race condition、原子操作、互斥鎖
- **通道 (channel)** — 傳遞訊息、`select` 多重來源
- **並行性運算的流程控制** — 緩衝區與關閉、等待結束、取消信號、產生器、方向限制、結構方法
- **context 套件** — 逾時與取消
- **章節總結**

<!--
今天的內容比較多，也是整門課最有挑戰性的一章，但它是 Go 的精華，值得多花一點時間。

我們會先認識 goroutine，也就是 Go 的「輕量執行緒」；接著看多個 goroutine 同時修改資料時會發生什麼問題，以及怎麼用原子操作和互斥鎖解決；然後學 Go 最有特色的通道；最後學 context，這是 Go 處理逾時和取消的標準方式，前面幾章已經用過很多次了。
-->

---

# 回顧：HTTP 伺服器

- Go 的 HTTP 伺服器會為**每一個請求啟動一個 goroutine**
- 商品 API 的 `Store` 用 **`sync.Mutex`** 保護 `map`，避免多個請求同時修改
- 優雅關閉：`go func() { srv.ListenAndServe() }()` 在背景啟動伺服器
- `signal.NotifyContext` 回傳的 **`ctx.Done()`** 在收到訊號時關閉 ← 今天揭曉原理

<!--
回顧一下上一章。

上一章有好幾個地方「先照著寫」：Store 裡的 sync.Mutex、go func() 在背景啟動伺服器、<-ctx.Done() 等待關閉訊號。今天會一一解釋這些程式碼背後的原理。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 使用 Go 語言的並行性運算
## Goroutines & WaitGroup

<!--
先從 goroutine 開始。
-->

---

# 並行（concurrency）與平行（parallelism）

| 概念 | 意思 | 比喻 |
| --- | --- | --- |
| **並行（concurrency）** | 同時**處理**多件事（可以輪流做） | 一位廚師同時顧三個爐子，輪流翻炒 |
| **平行（parallelism）** | 同時**執行**多件事（真的同時） | 三位廚師各顧一個爐子 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 的做法：</b> 我們用 goroutine 描述「可以並行的工作」，Go 的排程器會自動把它們分配到所有 CPU 核心上平行執行（<code>runtime.NumCPU()</code> 個核心）。
</div>

<!--
先釐清兩個常被混用的名詞。

並行是「同時處理」多件事：一位廚師同時顧三個爐子，這個翻一下、那個攪一下，雖然同一瞬間只做一件事，但三道菜都在進行。平行是「同時執行」：三位廚師各顧一個爐子，真的同時在做。

Go 的設計是：我們負責把程式寫成「可以並行」的形式，也就是把工作拆成多個 goroutine；Go 的排程器負責把它們分配到電腦的所有 CPU 核心上，讓它們真正平行執行。我們不需要自己管理執行緒。
-->

---
zoom: 0.96
---

# Goroutine

在函式呼叫前面加上 **`go`**，就會啟動一個新的 goroutine，**不等它執行完**就繼續往下

```go
package main

import (
	"fmt"
	"time"
)

func cook(dish string) {
	for i := 1; i <= 3; i++ {
		fmt.Println(dish, "步驟", i)
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	go cook("炒飯") // 在新的 goroutine 執行
	go cook("湯麵")
	cook("煎蛋") // 在 main goroutine 執行
	// ⚠️ 暫時用 Sleep 等待，下一頁會改用 WaitGroup
	time.Sleep(50 * time.Millisecond)
}
```

<!--
什麼是 goroutine？它是 Go 的「輕量級執行緒」，在函式呼叫前面加上 go 關鍵字，這個函式就會在一個新的 goroutine 裡執行，而且 main 不會等它，會馬上繼續往下執行。

執行後，三道菜的步驟會交錯印出，每次執行的順序都可能不一樣，因為它們是同時進行的。

goroutine 非常輕量：一個作業系統的執行緒大約需要 1MB 的記憶體，一個 goroutine 一開始只需要幾 KB，所以一個程式同時跑幾十萬個 goroutine 都沒問題。這就是為什麼 Go 的 HTTP 伺服器可以放心地替每個請求開一個 goroutine。

注意最後一行的 Sleep：main 函式結束時，整個程式就結束了，其他 goroutine 不管有沒有做完都會被強制終止。這裡先用 Sleep 等待，但這是很差的做法，下一頁會學正確的方式。
-->

---
zoom: 0.97
---

# WaitGroup：等待一群 goroutine 完成

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	dishes := []string{"炒飯", "湯麵", "煎蛋"}
	var wg sync.WaitGroup // 零值就能使用

	for _, d := range dishes {
		wg.Add(1) // 1. 啟動前先把計數 +1
		go func() {
			defer wg.Done() // 2. 完成時計數 -1
			time.Sleep(10 * time.Millisecond)
			fmt.Println(d, "完成")
		}()
	}
	wg.Wait() // 3. 等到計數歸零
	fmt.Println("全部上菜！")
}
```

<!--
sync.WaitGroup 是等待一群 goroutine 完成的標準工具，它內部有一個計數器。

使用方式是三個步驟：啟動 goroutine 之前，Add(1) 把計數加一；goroutine 完成的時候，Done() 把計數減一，通常搭配 defer，確保一定會執行；main 呼叫 Wait()，會一直等到計數歸零。

這裡的 go func() { ... }() 是第 5 章學的「匿名函式宣告後立刻執行」，加上 go 就在新的 goroutine 裡執行。

注意 goroutine 裡用到了迴圈變數 d。第 2 章說過，Go 1.22 之後每一圈的 d 都是新的變數，所以每個 goroutine 拿到的是不同的菜名。在舊版 Go 裡，這是非常經典的 bug。
-->

---

# 補充：WaitGroup.Go（Go 1.25+）

`wg.Go(f)` 自動完成 `Add(1)`、`go`、`defer Done()` 三件事，**不會忘記、也不會寫錯位置**

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	results := make([]int, 5)
	for i := range 5 {
		wg.Go(func() { // Go 1.25+
			results[i] = i * i // 每個 goroutine 寫不同的索引，不會互相干擾
		})
	}
	wg.Wait()
	fmt.Println(results) // [0 1 4 9 16]
}
```

<!--
這是補充內容：Go 1.25 替 WaitGroup 加入了 Go 方法。

傳統寫法有三個步驟，很容易出錯：忘記 Add、把 Add 寫在 goroutine 裡面（這樣 Wait 可能在 Add 之前就執行了）、忘記 Done。wg.Go 把這三件事包在一起，傳入一個函式，它會自動處理計數，是目前最推薦的寫法。

這個範例讓每個 goroutine 把結果寫到切片的不同位置，因為每個 goroutine 寫的是不同的索引，所以不會互相干擾。如果多個 goroutine 寫的是「同一個」變數，就會發生下一節要講的資料競爭。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 解決記憶體資源競爭
## Race Conditions

<!--
接下來看多個 goroutine 同時修改同一個變數時，會發生什麼問題。
-->

---
zoom: 0.97
---

# 什麼是資料競爭（data race）？

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	count := 0
	var wg sync.WaitGroup
	for range 1000 {
		wg.Go(func() {
			count++ // ⚠️ 多個 goroutine 同時「讀取 → 加一 → 寫回」
		})
	}
	wg.Wait()
	fmt.Println(count) // 預期 1000，實際可能是 953、987…每次都不一樣
}
```

```bash
go run -race .   # 開啟競爭偵測器
# WARNING: DATA RACE
# Read at 0x00c000012345 by goroutine 8: ...
```

<!--
什麼是資料競爭？

count++ 看起來是一個動作，其實是三個步驟：讀取 count、加一、寫回去。如果兩個 goroutine 同時讀到 count 是 5，各自加一，都寫回 6，那就少算了一次。

就像兩個人同時用同一本存摺提款：兩個人都看到餘額 1000 元，各自提了 300 元，都寫下「餘額 700 元」，結果銀行少扣了 300 元。

所以這個程式跑 1000 次加一，結果常常不是 1000，而且每次都不一樣。這種 bug 最可怕的地方是：它不一定每次都發生，可能在開發時完全正常，上線之後才偶爾出錯。

Go 提供了一個神兵利器：-race 競爭偵測器。go run、go test 都可以加上 -race，它會在執行時偵測資料競爭，印出是哪兩個 goroutine、在哪一行同時存取了同一個變數。寫並行程式時，測試一定要加 -race。
-->

---
zoom: 0.8
---

# 原子操作 (atomic operation)

`sync/atomic` 提供**不可分割**的操作，適合簡單的計數器、旗標

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var count atomic.Int64 // Go 1.19+ 的型別化原子值
	var wg sync.WaitGroup
	for range 1000 {
		wg.Go(func() {
			count.Add(1) // 讀取、加一、寫回在一個不可分割的步驟完成
		})
	}
	wg.Wait()
	fmt.Println(count.Load()) // 1000，每次都正確
}
```

| 型別 | 常用方法 |
| --- | --- |
| `atomic.Int64` / `Int32` / `Uint64` | `Add`、`Load`、`Store`、`CompareAndSwap` |
| `atomic.Bool` / `atomic.Pointer[T]` | `Load`、`Store`、`Swap` |

<!--
解決資料競爭的第一種方法：原子操作。

「原子」的意思是不可分割：atomic 套件提供的操作，由 CPU 保證在一個步驟內完成，其他 goroutine 不可能在中間插隊。count.Add(1) 就是一個原子的加一操作，所以結果永遠是 1000。

Go 1.19 加入了型別化的原子值，例如 atomic.Int64，用法比舊的 atomic.AddInt64(&n, 1) 更安全、更好讀。

原子操作的效能很好，但它只適合非常簡單的情況：一個計數器、一個開關旗標。如果要保護的是多個變數、或是一段複雜的邏輯，就要用下一頁的互斥鎖。
-->

---
zoom: 0.77
---

# 互斥鎖 (mutex)

`sync.Mutex` 保證同一時間**只有一個 goroutine** 能執行被鎖住的程式碼區段

```go
package main

import (
	"fmt"
	"sync"
)

type Account struct {
	mu      sync.Mutex // 慣例：mutex 放在它所保護的欄位上方
	balance int
	history []string
}

func (a *Account) Deposit(n int) {
	a.mu.Lock()         // 取得鎖：其他 goroutine 會在這裡排隊等待
	defer a.mu.Unlock() // 函式結束時釋放鎖
	a.balance += n
	a.history = append(a.history, fmt.Sprintf("+%d", n))
}

func main() {
	acc := &Account{}
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() { acc.Deposit(10) })
	}
	wg.Wait()
	fmt.Println(acc.balance, len(acc.history)) // 1000 100
}
```

<!--
第二種方法是互斥鎖，Mutex 是 mutual exclusion 的縮寫。

想像一間只有一把鑰匙的廁所：要進去之前先拿鑰匙（Lock），用完把鑰匙還回去（Unlock）。鑰匙被拿走的時候，其他人只能在門口排隊等。

Deposit 方法一開始 Lock，用 defer Unlock，這樣整個方法的內容，同一時間只會有一個 goroutine 在執行。這裡同時保護了 balance 和 history 兩個欄位，這是原子操作做不到的。

這就是上一章商品 Store 的寫法。

補充：sync.Mutex 的零值就能直接使用，不需要初始化，這又是 Go「零值可用」的設計。
-->

---

# 使用互斥鎖的注意事項

**注意事項之一：** `Lock` 之後一定要 `Unlock`，建議用 **`defer`**；忘記解鎖會造成**死結（deadlock）**

**注意事項之二：** 含有 Mutex 的結構**不能被複製**，方法要用**指標接收器**（`go vet` 會檢查）

**注意事項之三：** 鎖住的範圍越小越好；讀多寫少時可以用 **`sync.RWMutex`**

| 型別 | 方法 | 適用 |
| --- | --- | --- |
| `sync.Mutex` | `Lock` / `Unlock` | 一般情況 |
| `sync.RWMutex` | `RLock` / `RUnlock`（讀）、`Lock` / `Unlock`（寫） | 讀取遠多於寫入，多個讀取可以同時進行 |
| `sync.Once` / `sync.OnceValue` | `Do(f)` / `OnceValue(f)()` | 只執行一次的初始化 |

<!--
使用互斥鎖有三個注意事項。

第一，Lock 之後一定要 Unlock。如果忘記解鎖，其他 goroutine 會永遠等下去，這叫做死結。用 defer Unlock 是最保險的做法。

第二，含有 Mutex 的結構不能被複製，因為複製之後會變成兩把不同的鎖，失去保護的效果。所以方法一定要用指標接收器，第 4 章說過的「包含不能複製的欄位要用指標接收器」就是指這個。go vet 會自動檢查這個錯誤。

第三，鎖住的範圍越小越好，不要在持有鎖的時候做很慢的事，例如呼叫網路 API。如果資料是「讀很多、寫很少」，例如設定檔、快取，可以用 RWMutex：多個讀取可以同時進行，只有寫入時才需要獨佔。
-->

---
layout: default
---

# 練習 1：並行下載計數器
### 任務說明

模擬同時下載 20 個檔案：

1. 寫函式 `download(id int) int`：`time.Sleep` 隨機 10～50 毫秒（`rand.IntN`），回傳檔案大小 `id * 100`
2. 用 `sync.WaitGroup`（`wg.Go`）同時啟動 20 個下載
3. 用 `atomic.Int64` 累計**總下載量**
4. 用 `sync.Mutex` 保護一個 `map[int]int`，記錄每個檔案的大小
5. 印出總下載量、map 的長度，以及總耗時（應該接近最慢的那一個，而不是全部加總）
6. 用 `go run -race .` 執行，確認沒有資料競爭

<!--
這個練習把 goroutine、WaitGroup、原子操作、互斥鎖全部用上。

第 5 步是並行的威力：20 個下載如果一個一個做，要花 20 × 平均 30 毫秒，大約 600 毫秒；同時做的話，總時間只會接近最慢的那一個，大約 50 毫秒。

rand.IntN 是第 8 章提過的 math/rand/v2 套件。
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

func download(id int) int {
	time.Sleep(time.Duration(10+rand.IntN(41)) * time.Millisecond)
	return id * 100
}
```

<!--
download 用 rand.IntN(41) 產生 0 到 40 的亂數，加上 10，就是 10 到 50 毫秒，再轉換成 time.Duration。回傳的檔案大小是 id 乘以 100。
-->

---
zoom: 0.91
---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
func main() {
	start := time.Now()
	var (
		wg    sync.WaitGroup
		total atomic.Int64
		mu    sync.Mutex
		sizes = map[int]int{}
	)
	for id := 1; id <= 20; id++ {
		wg.Go(func() {
			n := download(id)
			total.Add(int64(n))
			mu.Lock()
			sizes[id] = n
			mu.Unlock()
		})
	}
	wg.Wait()
	elapsed := time.Since(start).Round(time.Millisecond)
	fmt.Println(total.Load(), len(sizes), elapsed)
	// 21000 20 50ms
}
```

<!--
每個 goroutine 下載完成後，用原子操作累加總量，用互斥鎖保護 map 的寫入。這裡鎖住的範圍只有 map 那一行，下載本身不需要鎖，這就是「鎖住的範圍越小越好」。

總下載量是 100 + 200 + … + 2000 = 21000，總耗時大約 50 毫秒，接近最慢的那一個。

用 go run -race . 執行，沒有任何警告，代表沒有資料競爭。試試看把 mu.Lock() 和 Unlock() 註解掉，再用 -race 跑一次，就會看到 DATA RACE 的警告。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 通道 (channel)
## Communicating Between Goroutines

<!--
接下來是 Go 最有特色的功能：通道。
-->

---
zoom: 0.9
---

# 使用通道傳遞訊息

「**通道（channel）是 goroutine 之間傳遞資料的管道**，一邊送、一邊收。」

| 語法 | 意義 |
| --- | --- |
| `ch := make(chan string)` | 建立一個傳遞 `string` 的通道 |
| `ch <- "hi"` | **送出**：把值送進通道 |
| `v := <-ch` | **接收**：從通道取出一個值 |

```go
package main

import "fmt"

func main() {
	ch := make(chan string) // 無緩衝通道：送出和接收必須「同時」發生
	go func() {
		ch <- "咖啡好了" // 送出，直到有人接收之前會在這裡等待
	}()
	msg := <-ch // 接收，直到有人送出之前會在這裡等待
	fmt.Println(msg)
}
```

<!--
什麼是通道？

想像餐廳廚房和外場之間的出餐口：廚師把做好的菜放上出餐口，服務生從出餐口把菜拿走。通道就是 goroutine 之間的出餐口。

箭頭的方向代表資料流動的方向：ch <- 值，是把值送進通道；<-ch，是從通道取出值。

用 make(chan T) 建立的是「無緩衝通道」：出餐口一次只能放一道菜，而且廚師要等到服務生來拿，才能離開。送出的一方會等到有人接收，接收的一方會等到有人送出，兩邊在通道「交會」的瞬間完成交接。所以通道不只能傳遞資料，還能「同步」兩個 goroutine。

這就是封面那句話的意思：goroutine 之間不共享變數，而是透過通道把資料「交給」對方。
-->

---
zoom: 0.84
---

# 通道 — 收集多個 goroutine 的結果

```go
package main

import (
	"fmt"
	"time"
)

type result struct {
	url  string
	size int
}

func fetch(url string, ch chan result) {
	time.Sleep(20 * time.Millisecond) // 模擬網路請求
	ch <- result{url, len(url) * 100} // 把結果送回通道
}

func main() {
	urls := []string{"a.com", "bb.com", "ccc.com"}
	ch := make(chan result)
	for _, u := range urls {
		go fetch(u, ch)
	}
	for range urls { // 知道有幾個結果，就接收幾次
		r := <-ch
		fmt.Println(r.url, r.size) // 順序依完成的先後而定
	}
}
```

<!--
通道最常見的用法：啟動多個 goroutine 做事，再透過通道收集它們的結果。

每個 fetch 做完之後，把結果送進通道。main 知道總共有幾個網址，就從通道接收幾次。這裡不需要 WaitGroup，因為「接收到三個結果」本身就代表三個 goroutine 都做完了。

結果的順序是依照完成的先後，不一定跟網址的順序一樣。

另外注意：這裡沒有任何共享變數，也不需要互斥鎖。每個 goroutine 只是把結果「交出去」，由 main 統一處理。這就是「透過溝通來共享記憶體」。
-->

---

# 從通道讀取多重來源的資料：select

`select` 同時等待**多個通道**，哪一個先準備好就執行哪一個

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	orders := make(chan string)
	payments := make(chan string)
	go func() {
		time.Sleep(20 * time.Millisecond)
		orders <- "訂單 #7"
	}()
	go func() {
		time.Sleep(10 * time.Millisecond)
		payments <- "付款 #7"
	}()
```

<!--
如果要同時等待多個通道，就用 select。

先準備兩個通道：訂單和付款，各自由一個 goroutine 在不同的時間送出資料。付款比較快，10 毫秒就送出；訂單要 20 毫秒。
-->

---

# 從通道讀取多重來源的資料：select（續）

```go
	// 續上頁
	for range 2 {
		select {
		case o := <-orders:
			fmt.Println("收到", o)
		case p := <-payments:
			fmt.Println("收到", p)
		case <-time.After(time.Second): // 逾時
			fmt.Println("等太久了")
			return
		}
	}
}
```

<!--
select 的語法跟 switch 很像，每個 case 是一個通道的操作，哪一個通道先準備好，就執行哪一個 case。如果同時有多個準備好，會隨機選一個。

付款比較快，所以會先印出「收到付款」，再印出「收到訂單」。

time.After 會回傳一個通道，在指定的時間之後送出一個值，搭配 select 就能實現「逾時」：如果一秒內兩個通道都沒有資料，就執行逾時的 case。第 12 章的 signal.NotifyContext 範例，就是用 select 同時等待計時器和關閉訊號。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 並行性運算的流程控制
## Channel Patterns

<!--
接下來學幾種使用通道的常見模式。
-->

---
zoom: 0.82
---

# 通道緩衝區與通道關閉：close()

| 語法 | 說明 |
| --- | --- |
| `make(chan int, 3)` | **有緩衝**的通道：可以先放 3 個值，滿了才會等待 |
| `close(ch)` | **關閉**通道：代表「不會再送出任何值」 |
| `for v := range ch` | 一直接收，**直到通道被關閉且取完**為止 |
| `v, ok := <-ch` | `ok` 為 `false` 代表通道已關閉且沒有值了 |

```go
package main

import "fmt"

func main() {
	ch := make(chan int, 3) // 緩衝區大小 3
	ch <- 1
	ch <- 2 // 緩衝區還沒滿，送出不需要等待
	close(ch)
	for v := range ch { // 取出 1、2 之後，因為已關閉而結束迴圈
		fmt.Print(v, " ")
	}
	v, ok := <-ch
	fmt.Println(v, ok) // 0 false
}
```

<!--
通道可以有緩衝區：make 的第二個參數是緩衝區大小。就像出餐口可以放三道菜，廚師放上去就可以繼續做下一道，不用等服務生來拿，直到出餐口滿了才需要等待。

close 用來關閉通道，意思是「我不會再送出任何東西了」。關閉之後，接收的一方還是可以把剩下的值取完。

for range 可以走訪通道：一直接收，直到通道被關閉而且裡面的值都取完了，迴圈就會結束。這是第 2 章提過、for range 可以走訪的對象之一。

comma ok 的寫法可以判斷通道是否已經關閉：ok 是 false，代表通道已關閉而且沒有值了，v 會是零值。又是 comma ok 慣用法！
-->

---

# 使用通道的注意事項

| 操作 | nil 通道 | 已關閉的通道 |
| --- | --- | --- |
| 送出 `ch <- v` | 永遠阻塞 | **panic** |
| 接收 `<-ch` | 永遠阻塞 | 立刻回傳零值（`ok` 為 `false`） |
| 關閉 `close(ch)` | **panic** | **panic**（重複關閉） |

**原則：** 由**送出的一方**負責關閉通道；接收的一方不要關閉

```go
// fatal error: all goroutines are asleep - deadlock!
ch := make(chan int)
ch <- 1 // 沒有其他 goroutine 接收，main 永遠等在這裡
```

<!--
使用通道有幾個重要的注意事項，這張表建議大家記下來。

最重要的原則是：由送出的一方負責關閉通道。因為只有送出的一方知道「還有沒有東西要送」。如果接收的一方關閉了通道，送出的一方再送就會 panic。

另一個常見的錯誤是死結：在 main 裡對一個無緩衝通道送出，但沒有任何 goroutine 在接收，main 就會永遠卡住。Go 的執行環境很聰明，如果發現所有 goroutine 都在等待、沒有人能繼續，會直接報錯 all goroutines are asleep - deadlock!，而不是讓程式默默卡死。
-->

---

# 使用通道訊息來等待 Goroutine 結束

```go
package main

import (
	"fmt"
	"time"
)

func worker(done chan struct{}) {
	fmt.Println("開始工作…")
	time.Sleep(30 * time.Millisecond)
	fmt.Println("工作完成")
	close(done) // 關閉通道 = 廣播「我做完了」
}

func main() {
	// struct{} 不佔任何記憶體，只用來傳遞「訊號」
	done := make(chan struct{})
	go worker(done)
	<-done // 等到 done 被關閉
	fmt.Println("main 結束")
}
```

<!--
通道也可以用來等待一個 goroutine 結束，這是 WaitGroup 之外的另一種方式。

done 通道的元素型別是 struct{}，空結構，它不佔任何記憶體，因為我們只需要傳遞「一個訊號」，不需要傳遞任何資料。

worker 做完工作後關閉 done，main 的 <-done 就會收到零值、不再阻塞，繼續往下執行。

為什麼用 close 而不是送出一個值？因為關閉通道是一種「廣播」：送出一個值只能被一個接收者收到，但關閉通道之後，所有正在等待的接收者都會同時收到通知。下一頁的取消信號就是利用這個特性。
-->

---
zoom: 0.75
---

# 使用通道傳送取消信號

關閉通道 = **廣播給所有 goroutine**，讓它們停止工作

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, quit <-chan struct{}) {
	for {
		select {
		case <-quit: // 通道被關閉時，這個 case 會立刻被選中
			fmt.Println("工人", id, "收到停止信號")
			return
		default: // 沒收到信號就繼續工作
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	quit := make(chan struct{})
	var wg sync.WaitGroup
	for id := 1; id <= 3; id++ {
		wg.Go(func() { worker(id, quit) })
	}
	time.Sleep(30 * time.Millisecond)
	close(quit) // 一次通知所有工人
	wg.Wait()
}
```

<!--
利用「關閉通道是廣播」的特性，就能實作取消信號。

三個工人都在一個無限迴圈裡工作，每一圈用 select 檢查 quit 通道：如果 quit 被關閉了，<-quit 會立刻成功，工人就結束；否則執行 default，繼續工作。select 的 default 代表「如果其他 case 都還沒準備好，就執行這裡」，所以不會阻塞。

main 在 30 毫秒後關閉 quit，三個工人同時收到信號、各自結束。

這個模式非常重要，它就是等一下要學的 context 的原理：ctx.Done() 回傳的，就是一個會在取消時被關閉的通道。
-->

---

# 使用函式來產生通道及 Goroutine：產生器（generator）

函式內部**建立通道、啟動 goroutine**，然後**回傳通道**給呼叫者使用

```go
package main

import "fmt"

func countdown(n int) <-chan int { // 回傳「只能接收」的通道
	ch := make(chan int)
	go func() {
		defer close(ch) // 送完就關閉
		for i := n; i > 0; i-- {
			ch <- i
		}
	}()
	return ch
}
```

<!--
這是一個很常見的模式：產生器。一個函式在內部建立通道、啟動一個 goroutine 往通道送資料，然後把通道回傳給呼叫者。呼叫者只要用 for range 接收就好，完全不用管 goroutine 的細節。

countdown 產生 n 到 1 的數字，送完之後 defer close 關閉通道，這樣呼叫者的 for range 才會結束。
-->

---

# 產生器與管線（pipeline）

```go
// 續上頁
func squares(in <-chan int) <-chan int { // 管線（pipeline）：接上前一段
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * v
		}
	}()
	return out
}

func main() {
	for v := range squares(countdown(4)) {
		fmt.Print(v, " ") // 16 9 4 1
	}
	fmt.Println()
}
```

<!--
squares 接收一個通道、回傳另一個通道，把收到的每個數字平方後送出。這樣就能把多個函式像水管一樣接起來，叫做「管線」（pipeline）：countdown 產生數字，流進 squares 平方，再流到 main 印出來。每一段都在自己的 goroutine 裡同時進行。

執行後會印出 16 9 4 1。
-->

---

# 限制通道的收發方向

| 型別 | 意義 | 能做的操作 |
| --- | --- | --- |
| `chan T` | 雙向通道 | 送出、接收、關閉 |
| `chan<- T` | **只能送出**（箭頭指向 chan） | 送出、關閉 |
| `<-chan T` | **只能接收**（箭頭從 chan 出來） | 接收 |

```go
func producer(out chan<- int) { // 這個函式只能送出
	out <- 1
	// v := <-out // 編譯錯誤：cannot receive from send-only channel
	close(out)
}

func consumer(in <-chan int) { // 這個函式只能接收
	fmt.Println(<-in)
	// close(in) // 編譯錯誤：cannot close receive-only channel
}
```

<!--
通道的型別可以限制方向：chan<- 是「只能送出」，<-chan 是「只能接收」。記法是看箭頭的方向：箭頭指向 chan，代表資料流進通道，所以是送出；箭頭從 chan 出來，代表資料從通道流出，所以是接收。

雙向通道可以自動轉換成單向通道，所以把一個一般的通道傳給 producer 和 consumer，它們各自只能做被允許的操作。

限制方向的好處是：讓編譯器幫我們檢查。consumer 不小心關閉了通道、producer 不小心從通道接收，都會在編譯時期就被抓到。上一頁的產生器回傳 <-chan int，就是告訴呼叫者「你只能接收，不能送、也不能關」。
-->

---

# 將結構方法當成 Goroutine：工作池（worker pool）

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	jobs    chan int
	results chan string
	wg      sync.WaitGroup
}

func (p *Pool) worker(id int) { // 方法也可以用 go 啟動
	defer p.wg.Done()
	for j := range p.jobs {
		time.Sleep(10 * time.Millisecond) // 模擬處理時間
		p.results <- fmt.Sprintf("工人%d 處理工作%d", id, j)
	}
}
```

<!--
結構的方法也可以用 go 關鍵字啟動成 goroutine，go p.worker(id) 就是在新的 goroutine 裡執行 p 的 worker 方法。

Pool 結構包含工作通道、結果通道和 WaitGroup。worker 方法用 for range 從 jobs 通道不斷領取工作，直到 jobs 被關閉；處理完就把結果送到 results。
-->

---

# 工作池（續）

```go
// 續上頁
func main() {
	p := &Pool{jobs: make(chan int, 10), results: make(chan string, 10)}
	for id := 1; id <= 3; id++ { // 固定 3 個工人
		p.wg.Add(1)
		go p.worker(id)
	}
	for j := 1; j <= 6; j++ {
		p.jobs <- j
	}
	close(p.jobs) // 沒有新工作了
	// 工人都結束後關閉結果通道
	go func() { p.wg.Wait(); close(p.results) }()
	for r := range p.results {
		fmt.Println(r)
	}
}
```

<!--
這是實務上非常常見的「工作池」模式：固定啟動 3 個工人，從同一個 jobs 通道搶工作來做。6 個工作會被 3 個工人分著做完。為什麼要限制工人的數量？例如要呼叫一個 API，同時發出一萬個請求可能會被對方封鎖，用工作池就能控制「最多同時 3 個」。

最巧妙的是結尾：送完所有工作後關閉 jobs，工人的 for range 就會結束；另外開一個 goroutine 等所有工人結束後，再關閉 results，main 的 for range 才會結束。這展示了「由送出的一方關閉通道」的原則。
-->

---
layout: default
---

# 練習 2：並行網址檢查器
### 任務說明

用工作池檢查一組網址的「回應時間」（以 `time.Sleep` 模擬）：

1. 寫產生器 `gen(urls ...string) <-chan string`，把網址逐一送出後關閉
2. 寫 `check(url string) (time.Duration, error)`：網址包含 `"bad"` 時回傳錯誤，否則模擬 `len(url)*5` 毫秒的延遲
3. 啟動 **3 個工人**，從產生器接收網址，把結果（網址、耗時、錯誤）送到結果通道
4. 使用**單向通道**作為參數型別；結果通道在所有工人結束後關閉
5. 主程式印出每個結果，最後統計成功與失敗的數量

<!--
這個練習把今天學的通道模式組合起來：產生器、工作池、單向通道、關閉通道的時機。

建議先畫出資料流動的方向：產生器 → 網址通道 → 3 個工人 → 結果通道 → main。每一段誰送、誰收、誰負責關閉，想清楚之後再寫程式。
-->

---
zoom: 0.84
---

# 練習 2：解題提示
### 提示說明

```go
package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type result struct {
	url string
	dur time.Duration
	err error
}

func gen(urls ...string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, u := range urls {
			ch <- u
		}
	}()
	return ch
}
```

<!--
先看資料結構和產生器。

result 結構包含網址、耗時和錯誤。gen 是產生器，跟剛剛的 countdown 一樣的寫法：建立通道、在 goroutine 裡逐一送出網址、送完關閉。
-->

---

# 練習 2：解題提示（續）
### 提示說明

```go
// 續上頁
func check(url string) (time.Duration, error) {
	if strings.Contains(url, "bad") {
		return 0, errors.New("連線失敗")
	}
	d := time.Duration(len(url)*5) * time.Millisecond
	time.Sleep(d)
	return d, nil
}
```

<!--
check 模擬檢查網址：包含 bad 就回傳錯誤，否則依網址長度睡一段時間，再回傳耗時。
-->

---

# 練習 2：解題提示（續 2）
### 提示說明

```go
// 續上頁
func worker(in <-chan string, out chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for u := range in {
		d, err := check(u)
		out <- result{u, d, err}
	}
}
```

<!--
worker 的參數用了單向通道：in 只能接收，out 只能送出。WaitGroup 要傳指標，因為剛剛說過同步物件不能被複製，WaitGroup 也一樣。

每個工人從 in 不斷領取網址、檢查，再把結果送到 out，直到 in 被關閉。
-->

---
zoom: 0.91
---

# 練習 2：解題提示（續 3）
### 提示說明

```go
// 續上頁
func main() {
	urls := gen("go.dev", "bad.example", "pkg.go.dev", "github.com")
	results := make(chan result)
	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go worker(urls, results, &wg)
	}
	go func() { wg.Wait(); close(results) }()

	ok, fail := 0, 0
	for r := range results {
		if r.err != nil {
			fail++
			fmt.Println("✗", r.url, r.err)
			continue
		}
		ok++
		fmt.Println("✓", r.url, r.dur)
	}
	fmt.Printf("成功 %d、失敗 %d\n", ok, fail) // 成功 3、失敗 1
}
```

<!--
main 啟動 3 個工人，都從同一個 urls 通道接收。另外開一個 goroutine 等工人都結束後關閉 results，main 就能用 for range 接收所有結果，最後統計成功與失敗的數量。

注意這裡用的是傳統的 Add、go、Done，因為 worker 是一個獨立的函式，需要自己呼叫 Done。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# context 套件
## Cancellation & Timeouts

<!--
最後學 context 套件，前面好幾章已經用過它了，今天來看它的完整面貌。
-->

---
zoom: 0.91
---

# 什麼是 context？

「**context 在呼叫鏈中傳遞『取消信號』與『截止時間』**，讓整條鏈上的工作可以一起停下來。」

| 函式 | 用途 |
| --- | --- |
| `context.Background()` | 最上層的空 context，在 `main`、測試中使用 |
| `context.WithCancel(parent)` | 回傳可以**手動取消**的 context 與 `cancel` 函式 |
| `context.WithTimeout(parent, d)` | 經過 `d` 之後**自動取消** |
| `context.WithDeadline(parent, t)` | 到了時間點 `t` **自動取消** |
| `context.WithValue(parent, key, val)` | 附帶**請求範圍**的值（例如 request ID） |

| `ctx` 的方法 | 意義 |
| --- | --- |
| `ctx.Done()` | 取消時會被**關閉**的通道（就是前面的取消信號模式！） |
| `ctx.Err()` | 取消的原因：`context.Canceled` 或 `context.DeadlineExceeded` |

<!--
什麼是 context？

想像一個 HTTP 請求進來，處理器要查資料庫、呼叫兩個外部 API。如果使用者等不及關掉了瀏覽器，這些後續的工作都應該停下來，不然就是在浪費資源。context 就是用來傳遞這種「取消信號」的：從處理器一路傳給資料庫查詢、API 呼叫，一旦取消，整條鏈上的工作都會收到通知。

context 的原理其實就是剛剛學的「關閉通道廣播」：ctx.Done() 回傳一個通道，取消的時候這個通道會被關閉。

建立 context 的函式都是 With 開頭，從一個父 context 衍生出子 context：WithCancel 可以手動取消，WithTimeout 和 WithDeadline 會在時間到的時候自動取消。父 context 被取消時，所有的子 context 也會一起被取消。
-->

---
zoom: 0.81
---

# 使用 context 控制逾時

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func slowQuery(ctx context.Context) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond): // 模擬需要 200ms 的查詢
		return "查詢結果", nil
	case <-ctx.Done(): // context 被取消或逾時
		return "", ctx.Err()
	}
}

func main() {
	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, 100*time.Millisecond)
	defer cancel() // ⚠️ 一定要呼叫 cancel，釋放計時器等資源

	_, err := slowQuery(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		// 查詢逾時： context deadline exceeded
		fmt.Println("查詢逾時：", err)
	}
}
```

<!--
這是 context 最常見的用法：控制逾時。

slowQuery 模擬一個需要 200 毫秒的查詢，它用 select 同時等待兩件事：查詢完成，或是 ctx.Done() 被關閉。哪一個先發生就執行哪一個。

main 設定了 100 毫秒的逾時，所以 ctx.Done() 會先被關閉，slowQuery 回傳 ctx.Err()，也就是 context.DeadlineExceeded。

使用 context 的注意事項：With 開頭的函式都會回傳一個 cancel 函式，一定要呼叫它，通常用 defer cancel()。即使逾時已經自動發生，呼叫 cancel 也能盡早釋放內部的計時器資源。go vet 會檢查有沒有忘記呼叫。

第 13 章的資料庫、第 14 章的 HTTP 請求，都是用這個方式設定逾時的：database/sql 和 net/http 內部都會監聽 ctx.Done()。
-->

---

# 使用 context 的注意事項

| ✅ 建議 | ❌ 避免 |
| --- | --- |
| 當作函式的**第一個參數**，命名為 `ctx` | 存在 struct 欄位裡 |
| 一路**往下傳遞**給資料庫、HTTP 等呼叫 | 在中間改用 `context.Background()`，切斷取消信號 |
| `defer cancel()` | 忘記呼叫 `cancel`（資源洩漏） |
| `WithValue` 只放**請求範圍**的資料（request ID、使用者） | 用 `WithValue` 傳遞一般的函式參數 |
| HTTP 處理器使用 **`r.Context()`** | 自己建立新的 `Background()` |

```go
func (a *API) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // 使用者斷線時會自動取消
	p, err := a.store.Get(ctx, id) // 一路傳給資料庫查詢
	// …
}
```

<!--
使用 context 有幾個 Go 社群公認的慣例。

第一，context 永遠是函式的第一個參數，名稱叫 ctx。第 13 章 Store 的每個方法都是這樣寫的。

第二，一路往下傳遞，不要在中間換成 Background，否則取消信號就斷了。

第三，WithValue 只用來傳遞「跟這次請求有關」的資料，例如 request ID、登入的使用者，不要拿它來取代一般的函式參數，因為它沒有型別檢查，會讓程式碼變得難以理解。

最後一個實務重點：HTTP 處理器裡，用 r.Context() 取得請求的 context。當使用者關掉瀏覽器、或是伺服器優雅關閉時，這個 context 會自動被取消，再把它傳給資料庫查詢，就能把整條鏈一起停下來。
-->

---
layout: default
---

# 綜合練習：有逾時的並行查價
### 任務說明

同時向 3 家供應商查詢咖啡豆的價格，**只要最便宜的結果，而且最多等 150 毫秒**：

1. 寫 `quote(ctx context.Context, vendor string, delay time.Duration, price int) (int, error)`：
   用 `select` 等待 `time.After(delay)` 或 `ctx.Done()`
2. 在 `main` 用 `context.WithTimeout` 建立 150ms 逾時的 context
3. 同時查詢：`A`（50ms, 480 元）、`B`（100ms, 450 元）、`C`（300ms, 400 元）
4. 用通道收集結果，**逾時的供應商視為失敗**
5. 印出每家的結果與最便宜的報價；並用 `-race` 確認沒有資料競爭

<!--
這個綜合練習模擬比價網站的場景：同時向多家供應商查價，但不能讓使用者等太久，所以設定一個總逾時。

用到了今天的所有重點：goroutine、通道收集結果、select、context 逾時。

想想看：C 最便宜，但它要 300 毫秒才會回應，超過了 150 毫秒的上限，最後會選到哪一家？
-->

---
zoom: 0.91
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"context"
	"fmt"
	"time"
)

type offer struct {
	vendor string
	price  int
	err    error
}

func quote(ctx context.Context, vendor string,
	delay time.Duration, price int) (int, error) {
	select {
	case <-time.After(delay):
		return price, nil
	case <-ctx.Done():
		return 0, fmt.Errorf("%s 未在時限內回應：%w", vendor, ctx.Err())
	}
}
```

<!--
quote 用 select 同時等待兩件事：供應商回應（time.After），或是 context 逾時。逾時的時候，用 %w 包裝 ctx.Err()，加上是哪一家供應商的情境。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func main() {
	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, 150*time.Millisecond)
	defer cancel()
	vendors := []struct {
		name  string
		delay time.Duration
		price int
	}{
		{"A", 50 * time.Millisecond, 480},
		{"B", 100 * time.Millisecond, 450},
		{"C", 300 * time.Millisecond, 400},
	}
```

<!--
main 建立 150 毫秒逾時的 context，再用匿名結構的切片準備三家供應商的測試資料：名稱、回應時間、報價。
-->

---
zoom: 0.91
---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
	// 續上頁
	// 緩衝區：逾時後 goroutine 也不會卡住
	ch := make(chan offer, len(vendors))
	for _, v := range vendors {
		go func() {
			p, err := quote(ctx, v.name, v.delay, v.price)
			ch <- offer{v.name, p, err}
		}()
	}
	best := offer{price: -1}
	for range vendors {
		o := <-ch
		if o.err != nil {
			fmt.Println("✗", o.err)
			continue
		}
		fmt.Println("✓", o.vendor, o.price)
		if best.price < 0 || o.price < best.price {
			best = o
		}
	}
	fmt.Println("最便宜：", best.vendor, best.price) // 最便宜： B 450
}
```

<!--
同時啟動三個 goroutine 查價，結果送進一個有緩衝區的通道。

為什麼要用緩衝區？如果 main 因為某些原因提早 return，不再接收結果，無緩衝通道會讓 goroutine 永遠卡在送出那一行，造成 goroutine 洩漏。緩衝區大小等於供應商數量，保證每個 goroutine 都能送出結果後正常結束。

A 在 50 毫秒、B 在 100 毫秒回應，C 要 300 毫秒，超過 150 毫秒的上限，所以收到逾時錯誤。最便宜的有效報價是 B 的 450 元。用 go run -race . 執行，沒有任何資料競爭，因為所有的結果都是透過通道傳遞。
-->

---

# 章節總結

- **goroutine**：`go f()` 啟動輕量執行緒；`main` 結束時所有 goroutine 都會結束
- **WaitGroup**：`Add` / `Done` / `Wait`；Go 1.25+ 用 **`wg.Go(f)`**
- **資料競爭**：用 **`go run -race`** / `go test -race` 偵測；簡單計數用 `atomic.Int64`，複雜狀態用 `sync.Mutex`（`defer Unlock`、指標接收器）
- **通道**：`ch <- v` 送出、`<-ch` 接收；無緩衝通道同步兩端；`select` 等待多個通道，`time.After` 做逾時
- **流程控制**：由**送出方關閉**通道；`for range ch` 收到關閉為止；關閉通道 = **廣播**；產生器、管線、工作池；單向通道 `chan<-` / `<-chan`
- **context**：`WithTimeout` / `WithCancel` + **`defer cancel()`**；`ctx.Done()` 是取消時被關閉的通道；當第一個參數一路往下傳

下一章我們會介紹「Go 語言工具」：`go build`、跨平台編譯、`gofmt`、`go vet`、`go doc`。

<!--
我們來整理今天學到的東西。

goroutine 是 Go 的輕量執行緒，用 go 關鍵字啟動，用 WaitGroup 等待它們完成。多個 goroutine 存取同一個變數時會有資料競爭，用 -race 偵測，用 atomic 或 Mutex 保護。通道是 goroutine 之間溝通的管道，select 可以同時等待多個通道。關閉通道是一種廣播，這是取消信號和 context 的原理。context 在整條呼叫鏈中傳遞逾時和取消。

並行性運算是 Go 最強大的武器，但也最容易寫出難以除錯的 bug，所以一定要養成用 -race 測試的習慣。

下一章會比較輕鬆，我們要認識 Go 的工具鏈：怎麼編譯出給 Windows、macOS、Linux 用的執行檔，怎麼格式化程式碼、做靜態分析、查詢文件。這些工具是 Go 開發體驗這麼好的原因。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天是整門課最有挑戰性的一章，第一次學覺得抽象是正常的。

課後建議：把第 14 章的 GitHub 客戶端，改成同時查詢 5 個 repository 的星星數，用工作池限制最多同時 2 個請求，並用 context 設定 3 秒的總逾時。這個練習會把 HTTP 客戶端和今天的並行性運算結合起來。

有問題的同學現在可以提問！
-->

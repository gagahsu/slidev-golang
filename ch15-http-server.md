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
title: 建立 HTTP 伺服器程式
routeAlias: ch15
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
  <h1 style="color: #1a5c5c; font-size: 3.2rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">建立 HTTP 伺服器程式</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「標準函式庫就能撐起正式環境的 Web 服務」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十五章！

上一章我們是「呼叫 API」的那一方，今天要換到另一端：自己建立 HTTP 伺服器，提供網頁和 API 給別人使用。

在其他語言裡，寫 Web 伺服器通常要先安裝一個框架，例如 Java 的 Spring Boot、Python 的 Django。Go 很不一樣：標準函式庫的 net/http 就是一個可以直接上正式環境的 HTTP 伺服器，而且從 Go 1.22 開始，內建的路由功能已經支援 HTTP 方法和路徑參數，很多專案完全不需要第三方框架。

今天會從最簡單的 Hello World 伺服器開始，一路做到一個完整的 RESTful JSON API。
-->

---
layout: default
---

# Outline

- **打造最基本的伺服器** — 處理器（handler）、路由、多重路徑
- **解讀網址參數** — 查詢參數、路徑參數 `{id}`
- **使用模板產生網頁** — `html/template`
- **使用靜態網頁資源** — 靜態檔案、`embed`、模板檔案
- **表單與 POST** — 讀取表單、更新伺服器資料
- **簡易 RESTful API** — 交換 JSON 資料、中介軟體、優雅關閉
- **GoShop 專案實作** — 第 15 步：RESTful API 與後台網頁
- **章節總結**

<!--
今天的內容由淺入深。

先寫出最基本的伺服器，理解「處理器」和「路由」這兩個核心概念。接著學怎麼讀取網址裡的參數，以及用模板產生動態網頁、提供靜態檔案。然後處理表單的 POST 請求。

最後是今天的重頭戲：一個完整的 RESTful JSON API，會用到前幾章學的 JSON、錯誤處理，還會加上日誌中介軟體和優雅關閉。
-->

---

# 回顧：HTTP 客戶端

- 請求 = **方法 + 路徑 + 標頭 + 本體**；回應 = **狀態碼 + 標頭 + 本體**
- 客戶端：`client.Get` / `client.Post` / `client.Do(req)`，記得**檢查狀態碼**與 `defer resp.Body.Close()`
- 狀態碼：`2xx` 成功、`4xx` 客戶端錯、`5xx` 伺服器錯
- 今天換到另一端：**讀取請求、寫出回應**

<!--
回顧一下上一章。

我們學了 HTTP 客戶端：組出請求、送出、讀取回應。今天的伺服器剛好相反：收到請求、讀取它的內容，然後寫出回應。上一章看到的方法、路徑、標頭、本體、狀態碼，今天都會從「伺服器的角度」再看一次。

建議大家今天開兩個終端機：一個跑伺服器，另一個用 curl 或上一章寫的 Go 客戶端來測試。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 打造最基本的伺服器
## Handlers & Routing

<!--
先從最基本的伺服器開始。
-->

---
zoom: 0.82
---

# 使用 HTTP 請求處理器 (handler)

「**處理器（handler）是負責『收到請求、寫出回應』的函式**。」

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", r.RemoteAddr) // 寫到 w，就是寫出回應
}

func main() {
	http.HandleFunc("/", hello) // 把路徑 "/" 對應到 hello
	log.Println("伺服器啟動：http://localhost:8080")
	// 開始接收請求（會一直執行）
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

| 參數 | 型別 | 用途 |
| --- | --- | --- |
| `w` | `http.ResponseWriter` | **寫出回應**：狀態碼、標頭、本體（它是 `io.Writer`） |
| `r` | `*http.Request` | **讀取請求**：方法、網址、標頭、本體 |

<!--
什麼是處理器？它就是一個函式，專門負責處理請求。

處理器的函式簽章是固定的：第一個參數 w 是 http.ResponseWriter，用來寫出回應；第二個參數 r 是 *http.Request，裡面有請求的所有資訊。

w 實作了 io.Writer，所以可以直接用 fmt.Fprintf 寫進去，寫進去的東西就是回應的本體。這又是第 7 章「接受介面」的威力。

main 裡面，HandleFunc 把路徑 "/" 對應到 hello 函式，然後 ListenAndServe 在 8080 連接埠開始接收請求。ListenAndServe 會一直執行，除非發生錯誤才會回傳，所以用 log.Fatal 包起來。

執行後打開瀏覽器，輸入 localhost:8080，就能看到 Hello 加上我們的 IP 位址。
-->

---
zoom: 0.9
---

# http.Handler 介面

`HandleFunc` 背後，其實是 Go 最重要的介面之一：

```go
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

```go
package main

import (
	"fmt"
	"net/http"
)

type Counter struct{ n int } // 任何有 ServeHTTP 方法的型別都是 Handler

func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.n++
	fmt.Fprintf(w, "你是第 %d 位訪客\n", c.n)
}

func main() {
	http.Handle("/count", &Counter{}) // Handle 接收 Handler 介面
	http.ListenAndServe(":8080", nil)
}
```

<!--
HandleFunc 背後，其實是 http.Handler 介面：只要有 ServeHTTP 方法的型別，就是一個 Handler。這是第 7 章學的隱性實作。

Counter 是一個 struct，有 ServeHTTP 方法，所以它就是 Handler，可以用 http.Handle 註冊到某個路徑。用 struct 當處理器的好處是可以帶著狀態，例如這個計數器。

HandleFunc 其實是一個方便的包裝：它把一般的函式轉換成 Handler。

補充一個重要的觀念：Go 的 HTTP 伺服器會為每一個請求啟動一個 goroutine，所以多個請求會同時執行 ServeHTTP，這個 Counter 其實有「資料競爭」的問題，第 16 章會學怎麼用互斥鎖解決。
-->

---
zoom: 0.87
---

# 簡單的 routing（路由）控制：ServeMux

**路由**就是「依照請求的方法與路徑，決定交給哪個處理器」

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 建立自己的路由器，不要用全域的 DefaultServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "首頁") // {$} 代表「剛好是 /」
		})
	mux.HandleFunc("GET /about",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "關於我們")
		})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>Go 1.22 起的路由樣式：</b> <code>"方法 路徑"</code>，例如 <code>"GET /about"</code>；方法不符會自動回傳 <code>405 Method Not Allowed</code>。
</div>

<!--
什麼是路由？就像大樓的總機：來電的人說要找業務部，總機就轉給業務部；要找客服，就轉給客服。路由器根據請求的方法和路徑，把請求轉給對應的處理器。

Go 的路由器叫做 ServeMux。建議用 http.NewServeMux 建立自己的，而不是用全域的預設路由器，這樣比較不容易被其他套件意外註冊路徑。

Go 1.22 大幅強化了 ServeMux：路徑前面可以寫 HTTP 方法，例如 "GET /about" 只接受 GET 請求，其他方法會自動回傳 405。以前這些都要自己在處理器裡寫 if 判斷，或是使用第三方的路由套件。

"/{$}" 是一個特殊寫法，代表「剛好是根路徑」。如果只寫 "/"，它會匹配所有沒有被其他規則匹配的路徑。
-->

---

# 修改程式來應付多重路徑請求

| 樣式 | 匹配 | 不匹配 |
| --- | --- | --- |
| `"/about"` | `/about`（任何方法） | `/about/team` |
| `"GET /about"` | `GET /about`（也包含 HEAD） | `POST /about` → 405 |
| `"/static/"` | `/static/`、`/static/a.css`、`/static/img/b.png` | `/staticx` |
| `"GET /items/{id}"` | `/items/42`、`/items/abc` | `/items/42/x` |
| `"GET /files/{path...}"` | `/files/a/b/c.txt` | — |
| `"/"` | **所有**沒被其他樣式匹配的路徑 | — |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>最精確的樣式優先：</b> 同時符合多個樣式時，ServeMux 會選擇<b>最具體</b>的那一個，與註冊順序無關。
</div>

<!--
這張表整理了 ServeMux 的路徑樣式規則。

沒有斜線結尾的路徑是「精確匹配」；斜線結尾的路徑是「前綴匹配」，底下的所有路徑都會匹配，適合提供靜態檔案。

大括號是「路徑參數」，例如 {id} 可以匹配任何一段文字，等一下會學怎麼讀取它。{path...} 加上三個點，可以匹配剩下的所有路徑。

另外有一個很重要的規則：如果一個請求同時符合多個樣式，ServeMux 會選最具體的那一個，跟註冊的順序無關。所以 "/" 永遠是最後的「預設」處理器，常用來回傳 404 頁面。
-->

---

# 多重路徑 — 範例

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("POST /products", createProduct) // 同路徑、不同方法
	mux.HandleFunc("/", notFound)                   // 其他所有路徑
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<!--
這個範例註冊了四個路由。

注意 /products 註冊了兩次：GET 對應 listProducts，POST 對應 createProduct。同一個網址，依照方法做不同的事，這就是等一下 RESTful API 的核心概念。最後的 "/" 接住所有沒有被匹配的路徑，回傳 404。
-->

---

# 多重路徑 — 範例（續）

```go
// 續上頁
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "首頁")
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "商品列表")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintln(w, "已新增商品")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "找不到頁面："+r.URL.Path, http.StatusNotFound)
}
```

<!--
四個處理器都很簡單，重點在 createProduct 和 notFound。

createProduct 用 w.WriteHeader 設定狀態碼 201 Created。使用 WriteHeader 的注意事項：它一定要在寫入本體「之前」呼叫，因為 HTTP 回應是先送狀態碼和標頭、再送本體，一旦開始寫本體，狀態碼就已經送出去了，改不了。沒有呼叫 WriteHeader 的話，預設是 200。

http.Error 是一個方便的函式，一次設定狀態碼和錯誤訊息。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 解讀網址參數來動態產生網頁
## Query & Path Parameters

<!--
接下來學怎麼讀取網址裡的參數，根據參數產生不同的內容。
-->

---
zoom: 0.85
---

# 查詢參數與路徑參數

| 種類 | 網址範例 | 讀取方式 |
| --- | --- | --- |
| **查詢參數** | `/search?q=咖啡&page=2` | `r.URL.Query().Get("q")` |
| **路徑參數**（Go 1.22+） | `/products/42` 對應 `"GET /products/{id}"` | `r.PathValue("id")` |

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello",
		func(w http.ResponseWriter, r *http.Request) {
			name := r.URL.Query().Get("name") // 沒有這個參數時是 ""
			if name == "" {
				name = "訪客"
			}
			fmt.Fprintf(w, "你好，%s！\n", name)
		})
```

<!--
網址裡的參數有兩種。

查詢參數是問號後面的 key=value，用 r.URL.Query().Get 讀取，這跟上一章組查詢參數用的 url.Values 是同一個型別。參數不存在時回傳空字串，所以要自己處理預設值。

測試方式：瀏覽器輸入 localhost:8080/hello?name=Gopher。
-->

---

# 查詢參數與路徑參數（續）

```go
	// 續上頁
	mux.HandleFunc("GET /products/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Error(w, "id 必須是數字", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "商品 #%d 的詳細資料\n", id)
		})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<!--
路徑參數是路徑的一部分，例如 /products/42 裡的 42。在路由樣式裡用大括號定義，再用 r.PathValue 讀取，這是 Go 1.22 加入的功能。

讀到的參數都是字串，要轉成數字就用 strconv.Atoi。使用參數的注意事項：參數是使用者可以任意輸入的，一定要驗證，轉換失敗要回傳 400 Bad Request，不能讓程式 panic。

測試方式：localhost:8080/products/42 和 /products/abc。
-->

---
layout: default
---

# 練習 1：BMI 計算 API
### 任務說明

建立一個伺服器，提供以下路由：

| 路由 | 功能 |
| --- | --- |
| `GET /{$}` | 回傳「BMI 計算器：請使用 /bmi?h=身高cm&w=體重kg」 |
| `GET /bmi` | 讀取查詢參數 `h`、`w`，回傳 `BMI = 22.9（正常）` |
| `GET /bmi/{h}/{w}` | 同上，但改用路徑參數 |

1. 參數缺少或不是數字時，回傳 `400 Bad Request` 與錯誤訊息
2. 把計算與分類寫成一個共用的函式 `bmi(hCm, wKg float64) string`
3. 用 `strconv.ParseFloat` 轉換參數

<!--
這個練習要大家寫出第一個真正有功能的伺服器。

重點是兩種參數的讀取方式，以及參數驗證：使用者可能不給參數、給了不是數字的參數，這些都要回傳 400。

分類標準可以沿用第 2 章 switch 的 BMI 練習：小於 18.5 過輕、小於 24 正常、小於 27 過重、其他肥胖。
-->

---
zoom: 0.94
---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func bmi(hCm, wKg float64) string {
	v := wKg / (hCm / 100 * hCm / 100)
	level := "肥胖"
	switch {
	case v < 18.5:
		level = "過輕"
	case v < 24:
		level = "正常"
	case v < 27:
		level = "過重"
	}
	return fmt.Sprintf("BMI = %.1f（%s）", v, level)
}
```

<!--
先把計算邏輯寫成一個獨立的函式 bmi，它完全不知道 HTTP 的存在，只負責計算。這樣的好處是：兩個路由可以共用，而且很容易寫單元測試。

把「商業邏輯」和「HTTP 處理」分開，是寫伺服器很重要的原則。
-->

---
zoom: 0.79
---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
func handle(w http.ResponseWriter, hs, ws string) {
	h, err1 := strconv.ParseFloat(hs, 64)
	wt, err2 := strconv.ParseFloat(ws, 64)
	if err1 != nil || err2 != nil || h <= 0 || wt <= 0 {
		http.Error(w, "h 與 w 必須是正數", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, bmi(h, wt))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "BMI 計算器：請使用 /bmi?h=身高cm&w=體重kg")
		})
	mux.HandleFunc("GET /bmi",
		func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			handle(w, q.Get("h"), q.Get("w"))
		})
	mux.HandleFunc("GET /bmi/{h}/{w}",
		func(w http.ResponseWriter, r *http.Request) {
			handle(w, r.PathValue("h"), r.PathValue("w"))
		})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<!--
handle 是一個共用的輔助函式：接收兩個字串參數，轉換、驗證、計算、寫出回應。兩個路由只差在「參數從哪裡來」，一個從查詢參數，一個從路徑參數。

參數缺少時 Get 回傳空字串，ParseFloat 會失敗；負數或 0 也不合理，一起回傳 400。

測試：/bmi?h=175&w=70 和 /bmi/175/70 都會回傳 BMI = 22.9（正常）。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 使用模板產生網頁
## html/template

<!--
接下來學怎麼產生 HTML 網頁。
-->

---

# 使用模板產生網頁：html/template

模板是「**有空格可以填資料的 HTML**」，用 `{{ }}` 標記要填入的地方

```go
package main

import (
	"html/template"
	"log"
	"net/http"
)

var page = template.Must(template.New("page").Parse(`
<h1>{{.Title}}</h1>
<ul>
{{range .Items}}<li>{{.Name}}：{{.Price}} 元</li>
{{else}}<li>目前沒有商品</li>{{end}}
</ul>`))
```

<!--
用 fmt.Fprintf 組 HTML 很痛苦，而且容易出錯。Go 提供了 html/template 模板套件。

模板就像一張「填空的表格」：HTML 裡用兩個大括號標記要填入資料的地方，{{.Title}} 代表「填入資料的 Title 欄位」，點代表目前的資料。{{range}} 走訪切片，每個元素重複一次；{{else}} 在切片是空的時候顯示。

template.Must 是第 6 章學的 Must 慣例：模板是寫死在程式碼裡的，寫錯就是 bug，直接 panic。模板只需要解析一次，所以宣告成套件層級的變數。
-->

---

# 使用模板產生網頁：html/template（續）

```go
// 續上頁
type Item struct {
	Name  string
	Price int
}

func main() {
	http.HandleFunc("GET /",
		func(w http.ResponseWriter, r *http.Request) {
			data := map[string]any{"Title": "今日菜單",
				"Items": []Item{{"咖啡", 60}, {"<b>蛋糕</b>", 120}}}
			if err := page.Execute(w, data); err != nil {
				log.Println(err)
			}
		})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

<!--
Execute 把資料填進模板，寫到 w。資料可以是 struct，也可以是 map，這裡用 map 放標題和商品切片。

注意第二個商品的名稱是 <b>蛋糕</b>，下一頁會說明它為什麼不會變成粗體。
-->

---

# 模板的常用語法

| 語法 | 意義 |
| --- | --- |
| `{{.Name}}` | 輸出目前資料的 `Name` 欄位（或 map 的鍵） |
| `{{if .IsVIP}} … {{else}} … {{end}}` | 條件判斷 |
| `{{range .Items}} … {{end}}` | 走訪切片或 map，區塊內的 `.` 是每一個元素 |
| `{{len .Items}}` / `{{printf "%.1f" .Price}}` | 呼叫內建函式 |
| `{{/* 註解 */}}` | 模板註解，不會輸出 |
| `{{template "header" .}}` | 引用另一個模板 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
🔐 <b>自動跳脫：</b> <code>html/template</code> 會依位置自動跳脫特殊字元，上一頁的 <code>&lt;b&gt;蛋糕&lt;/b&gt;</code> 會原樣顯示成文字，而不是變成粗體，藉此防止 <b>XSS 攻擊</b>。請不要改用 <code>text/template</code> 產生 HTML。
</div>

<!--
這張表整理了模板最常用的語法。

模板可以做條件判斷、走訪切片、呼叫函式，但建議讓模板保持簡單，複雜的邏輯寫在 Go 程式碼裡，模板只負責「顯示」。

最重要的是安全性：html/template 會自動跳脫特殊字元。上一頁的蛋糕名稱是 <b>蛋糕</b>，如果直接輸出，瀏覽器會把它當成 HTML 標籤。html/template 會把角括號轉換成 &lt; 和 &gt;，瀏覽器就只會顯示成文字。

這可以防止 XSS 跨站腳本攻擊：惡意使用者在留言裡寫一段 JavaScript，如果沒有跳脫，其他人瀏覽時就會執行那段程式碼。Go 還有一個 text/template 套件，它不會做 HTML 跳脫，所以產生 HTML 一定要用 html/template。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 使用靜態網頁資源
## Static Files & embed

<!--
接下來看怎麼提供靜態檔案：HTML、CSS、圖片、JavaScript。
-->

---
zoom: 0.97
---

# 讀取靜態 HTML 網頁：http.ServeFile

```text
myweb/
├── main.go
└── static/
    ├── index.html
    ├── style.css
    └── img/logo.png
```

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/index.html") // 回傳單一檔案
		})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<!--
提供單一檔案，用 http.ServeFile。它會讀取檔案、根據副檔名設定正確的 Content-Type，還會處理快取相關的標頭，比自己用 os.ReadFile 再寫出去完整很多。

這裡只把首頁對應到 static/index.html。但一個網站通常有很多靜態檔案，一個一個註冊太麻煩，下一頁看怎麼一次提供整個資料夾。
-->

---

# 在伺服器上提供多重靜態資源：http.FileServer

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// 把 static 資料夾變成 Handler
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fs))
	// /assets/style.css    → static/style.css
	// /assets/img/logo.png → static/img/logo.png
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>StripPrefix：</b> 請求的路徑是 <code>/assets/style.css</code>，但檔案在 <code>static/style.css</code>，所以要先把 <code>/assets/</code> 前綴去掉，再交給 FileServer。
</div>

<!--
http.FileServer 可以把一整個資料夾變成一個 Handler，資料夾裡的所有檔案都可以透過網址存取。

我們把它註冊在 "/assets/" 這個前綴路徑，所以 /assets/ 開頭的請求都會交給它。但是 FileServer 會用請求的完整路徑去找檔案，/assets/style.css 會去找 static/assets/style.css，這就不對了。所以要用 http.StripPrefix 先把 /assets/ 這個前綴去掉。

StripPrefix 本身也是一個 Handler，它包裝了另一個 Handler，這就是第 5 章學的 middleware 概念。

FileServer 會自動防止路徑穿越攻擊，請求 /assets/../main.go 是拿不到原始碼的。
-->

---
zoom: 0.93
---

# 把靜態資源編譯進執行檔：embed

用 `//go:embed` 把檔案**嵌入執行檔**，部署時只需要一個檔案（Go 1.16+）

```go
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var staticFiles embed.FS // 編譯時把 static 資料夾嵌入

func main() {
	sub, err := fs.Sub(staticFiles, "static") // 去掉最外層的 static/
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("GET /assets/",
		http.StripPrefix("/assets/", http.FileServerFS(sub))) // Go 1.22+
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<!--
第 0 章說過，Go 最大的優點之一是「編譯出一個執行檔，丟上伺服器就能跑」。但如果網站有 HTML、CSS、圖片這些靜態檔案，部署時就要連同資料夾一起複製，很容易漏掉。

Go 1.16 加入了 embed 功能：在變數上面寫一行 //go:embed 註解，編譯時就會把指定的檔案或資料夾「嵌入」執行檔裡。部署時只要一個檔案，所有靜態資源都在裡面。

注意 //go:embed 和 go 之間不能有空白，而且這個變數必須是套件層級的。

embed.FS 實作了 fs.FS 介面，Go 1.22 加入的 http.FileServerFS 可以直接接收它。fs.Sub 用來去掉最外層的 static 資料夾。

這段程式需要專案裡真的有 static 資料夾才能編譯。
-->

---
zoom: 0.81
---

# 使用模板檔案產生動態網頁

模板放在獨立的 `.html` 檔，用 `ParseFS` 從嵌入的檔案系統解析

```html
<!-- templates/products.html -->
<!DOCTYPE html>
<html>
<head><link rel="stylesheet" href="/assets/style.css"></head>
<body>
  <h1>商品列表（{{len .}} 項）</h1>
  <table>
    {{range .}}
    <tr><td>{{.Name}}</td><td>{{.Price}} 元</td></tr>
    {{end}}
  </table>
</body>
</html>
```

```go
//go:embed templates
var tmplFS embed.FS

var tmpl = template.Must(template.ParseFS(tmplFS, "templates/*.html"))

func productsPage(w http.ResponseWriter, r *http.Request) {
	items := []Item{{"咖啡", 60}, {"蛋糕", 120}}
	err := tmpl.ExecuteTemplate(w, "products.html", items)
	if err != nil {
		http.Error(w, "頁面產生失敗", http.StatusInternalServerError)
	}
}
```

<!--
模板寫在 Go 程式碼裡的字串中，不好編輯，也沒有編輯器的 HTML 語法提示。實務上會把模板放在獨立的 .html 檔案裡。

template.ParseFS 可以從嵌入的檔案系統解析模板，第二個參數是萬用字元，一次載入 templates 資料夾下的所有 .html 檔。ExecuteTemplate 用檔名指定要使用哪一個模板。

這個 products.html 模板裡，資料本身就是切片，所以 {{len .}} 取長度，{{range .}} 走訪。注意 CSS 的連結指向上一頁設定的 /assets/ 路徑。

這樣就組成了一個完整的網站：模板負責動態內容、FileServer 負責靜態資源，全部嵌入在一個執行檔裡。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 用表單和 POST 方法更新伺服器資料
## Forms & POST

<!--
接下來處理表單：讓使用者透過網頁把資料送到伺服器。
-->

---

# 讀取表單資料：r.FormValue

```go
package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var (
	messages []string // 留言（暫存在記憶體中）
	board    = template.Must(template.New("b").Parse(`
<form method="POST" action="/messages">
  <input name="text" placeholder="寫下留言"> <button>送出</button>
</form>
<ul>{{range .}}<li>{{.}}</li>{{end}}</ul>`))
)
```

<!--
這是一個簡單的留言板。

首頁用模板顯示一個表單和所有留言。表單的 method 是 POST，action 是 /messages，按下送出時，瀏覽器會把表單資料用 POST 送到 /messages。這跟上一章 http.PostForm 送出的格式一模一樣。
-->

---

# 讀取表單資料：r.FormValue（續）

```go
// 續上頁
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}",
		func(w http.ResponseWriter, r *http.Request) {
			board.Execute(w, messages)
		})
	mux.HandleFunc("POST /messages",
		func(w http.ResponseWriter, r *http.Request) {
			text := strings.TrimSpace(r.FormValue("text")) // 讀取表單欄位
			if text != "" {
				messages = append(messages, text)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther) // 303：導回首頁
		})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

<!--
伺服器用 r.FormValue 讀取表單欄位，參數是欄位的 name。它會自動解析表單，欄位不存在時回傳空字串。

處理完之後，用 http.Redirect 把使用者導回首頁，狀態碼是 303 See Other。這個模式叫做 PRG（Post/Redirect/Get）：如果不導向，使用者按重新整理，瀏覽器會重新送出一次 POST，留言就會重複。

補充：messages 被多個請求同時修改，會有資料競爭的問題，第 16 章會學怎麼用互斥鎖保護它。
-->

---

# 使用表單的注意事項

**注意事項之一：** 一律**驗證**使用者的輸入：長度、格式、範圍；**不要信任前端的檢查**

**注意事項之二：** 限制請求本體的大小，避免惡意的超大請求

**注意事項之三：** 輸出使用者的資料時，使用 `html/template` 自動跳脫（防止 XSS）

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 本體最多 1 MB
if err := r.ParseForm(); err != nil {
	http.Error(w, "表單格式錯誤", http.StatusBadRequest)
	return
}
name := r.PostFormValue("name") // 只讀 POST 本體，不讀網址的查詢參數
if n := utf8.RuneCountInString(name); n == 0 || n > 20 {
	http.Error(w, "名稱需為 1～20 個字", http.StatusBadRequest)
	return
}
```

<!--
使用表單有三個注意事項，都跟安全有關。

第一，一律驗證使用者的輸入。前端網頁可以用 JavaScript 檢查，但使用者可以繞過前端，直接用 curl 送出任何資料，所以後端一定要再驗證一次。

第二，限制請求本體的大小。http.MaxBytesReader 包裝請求本體，超過上限就會讀取失敗，防止有人送出好幾 GB 的請求把伺服器記憶體灌爆。

第三，輸出使用者的資料時，一定要用 html/template 自動跳脫。

另外 r.PostFormValue 只會讀取 POST 本體裡的欄位，不會讀網址的查詢參數，比 FormValue 更明確。驗證長度時用 RuneCountInString，這是第 3 章學的，中文字要以字元計算。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 簡易 RESTful API：交換 JSON 資料
## RESTful JSON API

<!--
最後是今天的重頭戲：RESTful JSON API。
-->

---

# 什麼是 RESTful API？

「**用網址表示『資源』，用 HTTP 方法表示『動作』，用 JSON 交換資料。**」

| 方法 + 路徑 | 動作 | 成功時的狀態碼 | 對應 SQL（Ch 13） |
| --- | --- | --- | --- |
| `GET /api/products` | 取得全部商品 | `200 OK` | `SELECT` |
| `GET /api/products/{id}` | 取得一個商品 | `200 OK`（找不到 `404`） | `SELECT … WHERE id = ?` |
| `POST /api/products` | 新增商品 | `201 Created` | `INSERT` |
| `PUT /api/products/{id}` | 更新商品 | `200 OK` | `UPDATE` |
| `DELETE /api/products/{id}` | 刪除商品 | `204 No Content` | `DELETE` |

<!--
什麼是 RESTful API？它是一種設計 API 的風格，核心概念是：用網址表示「資源」，例如 /api/products 代表商品這種資源；用 HTTP 方法表示「動作」，GET 取得、POST 新增、PUT 更新、DELETE 刪除。

這樣設計的好處是「一看就懂」：看到 DELETE /api/products/42，不用看文件就知道是刪除 42 號商品。

表格的最後一欄，是第 13 章學的 SQL 指令。大家會發現，RESTful API 的五個操作，剛好對應到資料庫的 CRUD。實務上，API 就是把資料庫操作包裝成 HTTP 介面，讓前端或其他服務呼叫。

今天為了專注在 HTTP 的部分，資料先存在記憶體裡；把它換成第 13 章的 Store，就是一個真正的後端服務了。
-->

---
zoom: 0.91
---

# API 的資料層：記憶體中的商品儲存庫

```go
// 檔名：store.go
package main

import (
	"errors"
	"sync"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var ErrNotFound = errors.New("找不到商品")

type Store struct {
	mu     sync.Mutex // 保護 items，允許多個請求同時存取（Ch 16）
	items  map[int]Product
	nextID int
}

func NewStore() *Store {
	return &Store{items: map[int]Product{}, nextID: 1}
}
```

<!--
先寫資料層。為了專注在 HTTP，資料存在記憶體的 map 裡。

這裡提前用了第 16 章的 sync.Mutex，互斥鎖。因為 HTTP 伺服器會同時處理多個請求，如果兩個請求同時修改 map，程式會出錯。互斥鎖保證同一時間只有一個請求能存取 items。下一章會詳細說明，現在先照著寫。

Product 加上 JSON 標籤，因為它會直接編碼成 API 的回應。
-->

---
zoom: 0.79
---

# API 的資料層（續）

```go
// 續上頁
func (s *Store) List() []Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 空的時候輸出 [] 而不是 null
	list := make([]Product, 0, len(s.items))
	for _, p := range s.items {
		list = append(list, p)
	}
	return list
}

func (s *Store) Get(id int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.items[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (s *Store) Create(p Product) Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = s.nextID
	s.nextID++
	s.items[p.ID] = p
	return p
}
```

<!--
Store 的每個方法都先 Lock、再 defer Unlock，這是第 5 章學的 defer 慣用法。

List 用 make 建立一個長度 0 的切片，而不是 nil 切片，這樣沒有商品的時候，JSON 會輸出空陣列而不是 null，這是第 11 章提過的前端友善做法。

Get 找不到商品時回傳哨兵錯誤 ErrNotFound。Create 指定新的 ID 並存入 map。
-->

---

# API 的資料層（續 2）

```go
// 續上頁
func (s *Store) Update(id int, p Product) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return Product{}, ErrNotFound
	}
	p.ID = id
	s.items[id] = p
	return p, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
```

<!--
Update 和 Delete 也都先檢查商品存不存在，不存在就回傳 ErrNotFound。

資料層寫完了。注意這裡完全沒有任何 HTTP 的程式碼，它只負責「怎麼存資料」。如果之後要換成 MySQL，只要把這些方法的內容換成第 13 章的 SQL 操作，HTTP 的部分完全不用改。
-->

---
zoom: 0.97
---

# 回傳 JSON 的輔助函式

```go
// 檔名：api.go
package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status) // 標頭要在 WriteHeader 之前設定
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("寫出 JSON 失敗", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
```

<!--
接著寫兩個輔助函式，讓處理器寫起來更簡潔。

writeJSON 做三件事：設定 Content-Type 標頭，告訴客戶端回應是 JSON；寫出狀態碼；把資料編碼成 JSON 寫出去。

注意順序：標頭一定要在 WriteHeader 之前設定，WriteHeader 一呼叫，標頭就送出去了，之後再設定都沒用。

writeError 統一錯誤的回應格式：{"error": "錯誤訊息"}。整個 API 的錯誤格式一致，前端處理起來就很方便。
-->

---

# 處理器：取得與新增

```go
// 續上頁
type API struct{ store *Store }

func (a *API) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.List())
}

func (a *API) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id 必須是數字")
		return
	}
	p, err := a.store.Get(id)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}
```

<!--
處理器定義成 API 型別的方法，API 帶著 store，這樣處理器就能存取資料層，不需要全域變數。這是實務上很常見的「依賴注入」寫法。

list 最簡單：取出所有商品，回傳 200 和 JSON。

get 先讀取路徑參數並轉成數字，失敗回傳 400；接著查詢商品，找不到回傳 404。這裡把資料層的 ErrNotFound 對應到 HTTP 的 404，這就是第 13 章說的「用自己的哨兵錯誤，讓上層不需要知道底層細節」。
-->

---
zoom: 0.97
---

# 處理器：新增、更新與刪除

```go
// 續上頁
func decodeProduct(w http.ResponseWriter,
	r *http.Request) (Product, bool) {
	var p Product
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		msg := "JSON 格式錯誤：" + err.Error()
		writeError(w, http.StatusBadRequest, msg)
		return p, false
	}
	if p.Name == "" || p.Price < 0 {
		writeError(w, http.StatusBadRequest, "name 必填，price 不可為負數")
		return p, false
	}
	return p, true
}

func (a *API) create(w http.ResponseWriter, r *http.Request) {
	if p, ok := decodeProduct(w, r); ok {
		writeJSON(w, http.StatusCreated, a.store.Create(p))
	}
}
```

<!--
新增和更新都需要讀取請求本體的 JSON，所以寫一個共用的 decodeProduct。

它用 MaxBytesReader 限制本體大小、用 Decoder 的嚴格模式解析，這是第 11 章學的；解析成功後再驗證欄位：名稱必填、價格不能是負數。任何一步失敗都回傳 400，並回傳 false 讓呼叫端知道要停止。

create 只要三行：解析成功就新增，回傳 201 Created 和新增後的商品，包含伺服器指定的 ID。
-->

---
zoom: 0.86
---

# 處理器：新增、更新與刪除（續）

```go
// 續上頁
func (a *API) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id 必須是數字")
		return
	}
	p, ok := decodeProduct(w, r)
	if !ok {
		return
	}
	if p, err = a.store.Update(id, p); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (a *API) remove(w http.ResponseWriter, r *http.Request) {
	// 轉換失敗時 id 為 0，會回傳 404
	id, _ := strconv.Atoi(r.PathValue("id"))
	if err := a.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204：成功，沒有本體
}
```

<!--
update 結合了 get 和 create 的邏輯：讀取路徑參數、解析 JSON、更新資料，找不到商品回傳 404。

remove 刪除成功時回傳 204 No Content，代表「成功了，但沒有東西要回傳」，所以只寫狀態碼、不寫本體。

這裡為了精簡，remove 忽略了 Atoi 的錯誤：轉換失敗時 id 是 0，而 0 號商品不存在，一樣會回傳 404。實務上建議跟 get 一樣回傳 400。
-->

---

# 中介軟體（middleware）：記錄每個請求

```go
// 續上頁
type statusRecorder struct {
	http.ResponseWriter // 內嵌：其他方法直接沿用（Ch 4）
	status              int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r) // 交給下一層處理
		slog.Info("request", "method", r.Method, "path", r.URL.Path,
			"status", rec.status, "duration", time.Since(start))
	})
}
```

<!--
中介軟體就是第 5 章綜合練習寫過的 middleware：接收一個 Handler，回傳一個新的 Handler，在呼叫原本的 Handler 前後加上額外的功能。

logging 記錄每個請求的方法、路徑、狀態碼和花費的時間，使用第 9 章學的 slog 結構化日誌、第 10 章的 time.Since。

為了取得狀態碼，我們需要「攔截」WriteHeader 的呼叫。statusRecorder 內嵌了原本的 ResponseWriter，所以它自動擁有 Header、Write 這些方法；我們只覆寫 WriteHeader，把狀態碼記下來再交給原本的方法。這是第 4 章內嵌、第 7 章介面的綜合應用。
-->

---
zoom: 0.77
---

# 組合路由並啟動伺服器（含優雅關閉）

```go
// 檔名：main.go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	api := &API{store: NewStore()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/products", api.list)
	mux.HandleFunc("POST /api/products", api.create)
	mux.HandleFunc("GET /api/products/{id}", api.get)
	mux.HandleFunc("PUT /api/products/{id}", api.update)
	mux.HandleFunc("DELETE /api/products/{id}", api.remove)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           logging(mux), // 用中介軟體包住整個路由器
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
```

<!--
main 把所有東西組合起來。

先建立 API，註冊五個路由，對應 RESTful 的五個操作，方法和路徑一目了然。

接著建立 http.Server，而不是直接呼叫 http.ListenAndServe。原因是正式環境一定要設定逾時：ReadHeaderTimeout 限制讀取標頭的時間，可以防止一種叫做 Slowloris 的攻擊，攻擊者故意非常慢地送出請求，佔住伺服器的連線。ReadTimeout、WriteTimeout 限制整個讀寫的時間，IdleTimeout 限制閒置連線保留多久。

Handler 設定成 logging(mux)，用中介軟體包住整個路由器，所有請求都會被記錄。
-->

---
zoom: 0.97
---

# 組合路由並啟動伺服器（含優雅關閉）（續）

```go
	// 續上頁
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() { // 在另一個 goroutine 啟動伺服器（Ch 16）
		slog.Info("伺服器啟動", "addr", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("伺服器錯誤", "err", err)
			stop()
		}
	}()

	<-ctx.Done() // 等待 Ctrl+C 或 SIGTERM
	slog.Info("收到關閉訊號，處理中的請求完成後結束…")
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil { // 優雅關閉
		slog.Error("關閉失敗", "err", err)
	}
}
```

<!--
最後是優雅關閉，這是第 12 章 signal.NotifyContext 的實際應用。

伺服器在另一個 goroutine 裡啟動，go 關鍵字是第 16 章的內容，現在先理解為「在背景執行」。主程式則停在 <-ctx.Done()，等待使用者按 Ctrl+C 或系統送出 SIGTERM。

收到訊號後，呼叫 srv.Shutdown：它會停止接收新的請求，並等待處理中的請求完成，最多等 10 秒。這樣正在付款的使用者不會被突然中斷。

ListenAndServe 在 Shutdown 之後會回傳 http.ErrServerClosed，這是正常的結束，不是錯誤，所以要排除。
-->

---

# 測試 API：curl

```bash
# 新增
curl -i -X POST localhost:8080/api/products \
  -H 'Content-Type: application/json' -d '{"name":"咖啡","price":60}'
# HTTP/1.1 201 Created
# {"id":1,"name":"咖啡","price":60}

# [{"id":1,"name":"咖啡","price":60}]
curl localhost:8080/api/products
curl localhost:8080/api/products/99       # {"error":"找不到商品"}
curl -X PUT localhost:8080/api/products/1 \
  -d '{"name":"咖啡","price":65}'
# HTTP/1.1 204 No Content
curl -i -X DELETE localhost:8080/api/products/1
curl -X POST localhost:8080/api/products -d '{"nme":"x"}'
# {"error":"JSON 格式錯誤：json: unknown field \"nme\""}
```

伺服器的日誌：

```text
2026/09/24 10:30:00 INFO request method=POST path=/api/products
                    status=201 duration=197.971µs
```

<!--
伺服器啟動後，用 curl 測試每一個 API。-i 會顯示回應的狀態碼和標頭，-X 指定方法，-H 設定標頭，-d 設定本體。

大家可以看到，每個操作都回傳了正確的狀態碼：新增 201、找不到 404、刪除 204，打錯欄位名稱的 JSON 被嚴格模式擋下來，回傳 400。

伺服器的終端機裡，logging 中介軟體會記錄每一個請求。

除了 curl，也可以用 VS Code 的 REST Client 擴充套件、Postman，或是上一章寫的 Go 客戶端來測試。
-->

---
zoom: 0.91
---

# 補充：用 httptest 替處理器寫單元測試

```go
// 檔名：api_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	api := &API{store: NewStore()}
	body := strings.NewReader(`{"name":"咖啡","price":60}`)
	req := httptest.NewRequest(http.MethodPost, "/api/products", body)
	rec := httptest.NewRecorder() // 假的 ResponseWriter，會記錄回應

	api.create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("狀態碼 = %d，期望 %d", rec.Code, http.StatusCreated)
	}
	if !strings.Contains(rec.Body.String(), `"id":1`) {
		t.Errorf("回應 = %s", rec.Body.String())
	}
}
```

<!--
這是補充內容：怎麼替 HTTP 處理器寫單元測試。

net/http/httptest 套件提供了測試用的工具：NewRequest 建立一個假的請求，NewRecorder 建立一個假的 ResponseWriter，它會把處理器寫出的狀態碼、標頭、本體都記錄下來。

直接呼叫處理器，再檢查 recorder 的內容，就能測試處理器的行為，完全不需要真的啟動伺服器、也不需要網路。這就是第 9 章學的單元測試，應用在 HTTP 處理器上。

執行 go test，就能自動驗證新增商品的 API 是否正確。
-->

---
layout: default
zoom: 0.96
---

# 綜合練習：待辦事項 API
### 任務說明

參考商品 API，完成一個待辦事項（Todo）的 RESTful API：

| 路由 | 功能 |
| --- | --- |
| `GET /api/todos` | 取得全部；支援查詢參數 `?done=true` / `?done=false` 篩選 |
| `POST /api/todos` | 新增（`title` 必填、最多 50 字） |
| `PATCH /api/todos/{id}/done` | 標記為完成 |
| `DELETE /api/todos/{id}` | 刪除 |

1. `Todo` 結構：`ID int`、`Title string`、`Done bool`、`CreatedAt time.Time`
2. 使用 `sync.Mutex` 保護資料；錯誤一律回傳 `{"error": "..."}`
3. 加上 `logging` 中介軟體與優雅關閉
4. 替 `POST /api/todos` 寫一個 `httptest` 單元測試

<!--
這個綜合練習要大家自己完成一個完整的 API。

大部分的程式碼可以參考商品 API，需要自己思考的地方有：查詢參數的篩選、PATCH 方法只修改部分欄位、標題長度的驗證要用字元數計算。

完成之後，試著用 curl 把每一個 API 都呼叫一次，確認狀態碼都正確。
-->

---
zoom: 0.85
---

# 綜合練習：解題提示
### 提示說明

```go
// 查詢參數篩選
func (a *API) listTodos(w http.ResponseWriter, r *http.Request) {
	list := a.store.List()
	if v := r.URL.Query().Get("done"); v != "" {
		want, err := strconv.ParseBool(v)
		if err != nil {
			msg := "done 必須是 true 或 false"
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		list = slices.DeleteFunc(list, func(t Todo) bool {
			return t.Done != want
		})
	}
	writeJSON(w, http.StatusOK, list)
}

// 路由
mux.HandleFunc("GET /api/todos", api.listTodos)
mux.HandleFunc("POST /api/todos", api.createTodo)
mux.HandleFunc("PATCH /api/todos/{id}/done", api.markDone)
mux.HandleFunc("DELETE /api/todos/{id}", api.removeTodo)
```

- 標題驗證：`utf8.RuneCountInString(t.Title)` 介於 1～50
- `markDone`：`store.MarkDone(id)` 回傳 `ErrNotFound` 時回應 `404`，成功回傳更新後的 Todo

<!--
這裡提示最有挑戰的部分：查詢參數篩選。

用 strconv.ParseBool 把字串轉成布林值，它接受 true、false、1、0 等寫法，不合法就回傳 400。篩選用 Go 1.21 加入的 slices.DeleteFunc，刪除不符合條件的元素，剩下的就是我們要的。

路由的部分，PATCH /api/todos/{id}/done 的路徑參數在中間，Go 1.22 的 ServeMux 完全支援。

其他部分跟商品 API 大同小異，大家可以把商品 API 的程式碼複製過來修改。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 15 步：API 與後台網頁

<!--
回到 GoShop。到目前為止，GoShop 都是在終端機裡操作的命令列工具。

但顧客不會打開終端機買東西。今天學了 HTTP 伺服器，我們終於可以讓 GoShop 變成一個網站：給前端和手機 App 用的 JSON API，加上給店長用的後台網頁。
-->

---
zoom: 0.91
---

# GoShop 第 15 步：API 與後台網頁
### 任務說明

| 方法與路徑 | 功能 | 成功 | 常見錯誤 |
| --- | --- | --- | --- |
| `GET /api/products` | 商品清單 | 200 | |
| `GET /api/products/{sku}` | 單一商品 | 200 | 404 |
| `POST /api/orders` | 下單（JSON 購物車） | 201 + `Location` | 400、409 |
| `GET /api/orders/{id}` | 查詢訂單 | 200 | 400、404 |
| `POST /api/orders/{id}/pay` | 付款（`cod`／`card`） | 200 | 404、409 |
| `GET /admin` | 後台：商品與訂單列表 | 200 | |
| `POST /admin/products` | 後台表單：新增或更新商品 | 303 | 400 |

- 模板與 CSS 用 **`embed`** 編譯進執行檔；每個請求用 `slog` 記錄
- `go run . -http :8080`，按 Ctrl+C 時**優雅關閉**，記憶體版本關閉前存檔

<!--
這是 GoShop 網站的所有路由，前五個是 JSON API，後兩個是後台網頁。

設計 API 的時候，狀態碼要用對：新增成功是 201 Created，並且在 Location 標頭告訴客戶端新訂單的網址；找不到是 404；庫存不足、重複付款這種「和目前狀態衝突」的情況是 409 Conflict；客戶端送來的資料有問題是 400。

後台表單送出成功之後回 303，讓瀏覽器重新導向到後台首頁，等一下會解釋為什麼。
-->

---

# GoShop 第 15 步：預期結果

```text
$ curl -i -X POST localhost:8080/api/orders \
    -d '{"items":[{"sku":"SKU-001","qty":2}],"coupon":"WEEK15"}'
HTTP/1.1 201 Created
Location: /api/orders/1
{"id":1,"lines":[…],"subtotal":900,"discount":135,"total":765,
 "status":"pending","coupon":"WEEK15",…}

$ curl -X POST localhost:8080/api/orders/1/pay \
    -d '{"method":"card","last4":"4242"}'
{"id":1,…,"status":"paid","paid_by":"信用卡 *4242",…}

$ curl -X POST localhost:8080/api/orders \
    -d '{"items":[{"sku":"SKU-003","qty":99}]}'
{"error":"結帳失敗：SKU-003 庫存不足：想買 99 件，只剩 5 件"}
```

<!--
這是用 curl 測試 API 的結果。

第一個請求下單，回應 201 和 Location 標頭，訂單狀態是 pending 待付款。第二個請求用信用卡付款，狀態變成 paid。第三個請求想買 99 個手沖壺，得到 409 和清楚的錯誤訊息。

大家也可以用 VS Code 的 REST Client 擴充套件，或是 Postman 這類工具來測試，會比 curl 方便。
-->

---

# GoShop 第 15 步：後台網頁

<img src="/img/goshop/admin.png" class="h-110 mx-auto border rounded shadow" alt="GoShop 後台">

<!--
這是瀏覽器打開 localhost:8080/admin 看到的後台網頁。

上面是商品列表，手沖壺只剩 2 件，庫存少於 5 件的商品會用紅字提醒店長補貨。中間是新增商品的表單，SKU 已經存在的話就是更新。下面是訂單列表，可以看到每張訂單的金額和付款狀態。

整個頁面只用了 html/template 和一個很短的 CSS 檔，沒有任何前端框架。對後台管理這種內部工具來說，這樣就已經很夠用了。
-->

---
class: code-sm
zoom: 0.91
---

# GoShop 第 15 步：解題提示
### 路由：Go 1.22 的 ServeMux 樣式

```go
// goshop/internal/web/server.go
//go:embed templates static
var assets embed.FS

// 啟動時就解析模板；模板有語法錯誤時程式會直接 panic
var tmpl = template.Must(template.ParseFS(assets, "templates/*.html"))

// Server 保存處理請求時需要的相依物件。
type Server struct {
	Store    store.Store
	Checkout *checkout.Service
}

// Handler 設定所有路由，傳回加上日誌中介軟體的 http.Handler。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/products", s.listProducts)
	mux.HandleFunc("GET /api/products/{sku}", s.getProduct)
	mux.HandleFunc("POST /api/orders", s.createOrder)
	// ...
	return logging(mux)
}
```

<!--
web 套件的 Server 結構保存處理請求需要的東西：Store 和結帳服務。處理器都寫成 Server 的方法，就能直接使用這些欄位，不需要全域變數。

路由用 Go 1.22 的新樣式：方法加路徑，路徑裡的 {sku} 是路徑參數，在處理器裡用 r.PathValue("sku") 取出。方法不符合的時候，ServeMux 會自動回 405。

templates 和 static 兩個資料夾用 go:embed 嵌入成一個 embed.FS。模板在程式啟動時就用 template.Must 解析好，寫錯的話程式一啟動就會 panic，而不是等到有人打開網頁才發現。

最後把整個 mux 包上本章寫過的 logging 中介軟體。
-->

---
class: code-sm
---

# GoShop 第 15 步：解題提示（續）
### 依錯誤種類決定狀態碼

```go
// goshop/internal/web/api.go
// writeErr 依錯誤種類決定 HTTP 狀態碼。
func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shop.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, shop.ErrOutOfStock), errors.Is(err, shop.ErrPaid):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, shop.ErrEmptyCart), errors.Is(err, shop.ErrCoupon),
		errors.Is(err, payment.ErrInsufficientFunds):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		slog.Error("內部錯誤", "err", err) // 細節只寫進日誌，不回給客戶端
		writeError(w, http.StatusInternalServerError, "伺服器發生錯誤")
	}
}
```

- 第 6 章設計的哨兵錯誤，在這裡決定了 API 的狀態碼

<!--
writeErr 是整個 API 錯誤處理的核心。它用 errors.Is 判斷錯誤的種類，決定要回哪一個狀態碼。

這裡就看得出第 6 章設計哨兵錯誤的價值了：store 和 checkout 完全不知道 HTTP 的存在，它們只傳回 ErrNotFound、ErrOutOfStock 這些錯誤；web 套件再把這些錯誤翻譯成 HTTP 的語言。各層各司其職。

default 的情況是我們沒預料到的錯誤，例如資料庫斷線。這時候詳細的錯誤只寫進日誌，回給客戶端的只有「伺服器發生錯誤」，避免把內部資訊洩漏出去。
-->

---
class: code-sm
zoom: 0.92
---

# GoShop 第 15 步：解題提示（續 2）
### 處理器：解碼 JSON、回傳 201

```go
// goshop/internal/web/api.go
func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	var cart checkout.Cart
	// 請求內容最多 1 MB，避免超大的請求吃光記憶體
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cart); err != nil {
		writeError(w, http.StatusBadRequest, "JSON 格式錯誤："+err.Error())
		return
	}
	o, err := s.Checkout.Checkout(r.Context(), cart)
	if err != nil {
		writeErr(w, err)
		return
	}
	w.Header().Set("Location", "/api/orders/"+strconv.Itoa(o.ID))
	writeJSON(w, http.StatusCreated, o)
}
```

- `checkout.Item`、`checkout.Cart` 加上 JSON 標籤，就能直接當成請求的格式
- `r.Context()`：客戶端斷線時會被取消，資料庫查詢也跟著停止

<!--
createOrder 把請求的 JSON 直接解碼成第 8 章的 checkout.Cart，只要替 Cart 和 Item 加上 JSON 標籤就行了。

MaxBytesReader 限制請求最多 1 MB，避免有人送一個超大的請求把伺服器的記憶體吃光。DisallowUnknownFields 讓欄位名稱打錯的請求直接得到 400，前端工程師很快就能發現問題。

結帳時傳入 r.Context()，這是每個請求自帶的 context。客戶端斷線的時候它會被取消，第 13 章 MySQL 的查詢也會跟著停止，不會浪費資料庫的資源。

成功之後先設定 Location 標頭，再用 writeJSON 回傳 201 和訂單內容。
-->

---
class: code-sm
---

# GoShop 第 15 步：解題提示（續 3）
### 後台模板與表單

```html
<!-- goshop/internal/web/templates/admin.html -->
    {{range .Products}}
    <tr class="{{if lt .Stock 5}}low{{end}}">
      <td>{{.SKU}}</td><td>{{.Name}}</td>
      <td>{{.Price}}</td><td>{{.Stock}}</td>
    </tr>
    {{end}}
```

```go
// goshop/internal/web/admin.go
	if err := s.Store.SaveProduct(r.Context(), p); err != nil {
		s.renderAdmin(w, r, http.StatusInternalServerError, "儲存失敗")
		return
	}
	// 成功後重新導向（PRG 模式），重新整理頁面才不會重複送出表單
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
```

- `{{.Price}}` 會呼叫 `Money.String()`，網頁上直接顯示 `NT$1,280`

<!--
模板用 range 走訪商品，用 if lt .Stock 5 判斷庫存是不是少於 5 件，是的話加上 low 這個 CSS class，讓整列變成紅字。

{{.Price}} 印出來是 NT$1,280，因為 html/template 也會呼叫 Stringer。第 7 章寫的 String 方法，在網頁上也派上用場了。

表單處理完成後，用 303 See Other 重新導向回後台首頁，這叫做 PRG 模式：Post、Redirect、Get。如果直接回傳網頁，使用者按重新整理，瀏覽器會問要不要重新送出表單，一不小心就重複新增了。
-->

---
zoom: 0.97
---

# GoShop 第 15 步：解題提示（續 4）
### 啟動伺服器與優雅關閉

```go
// goshop/main.go
	app := &web.Server{Store: st, Checkout: newService(st)}
	srv := &http.Server{
		Addr:              addr,
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	slog.Info("伺服器啟動", "addr", addr)

	select {
	case err := <-errCh: // 例如連接埠已被佔用
		return err
	case <-ctx.Done():
	}
	slog.Info("收到結束訊號，關閉伺服器中…")
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx) // 等進行中的請求處理完才關閉
```

<!--
最後是啟動伺服器。這段和本章「組合路由並啟動伺服器」的寫法一樣：用 signal.NotifyContext 等待 Ctrl+C，在另一個 goroutine 裡執行 ListenAndServe，收到訊號後呼叫 Shutdown，等進行中的請求都處理完才關閉。

這裡多了一個 errCh 通道：如果連接埠已經被佔用，ListenAndServe 會馬上失敗，我們要把這個錯誤傳回來，而不是傻傻地等 Ctrl+C。select 同時等待兩件事，哪個先發生就處理哪個。goroutine、通道和 select 是下一章的主題。

serve 結束後，回到 run 函式，記憶體版本會把資料存成快照，所以伺服器重新啟動後資料都還在。
-->

---

# 章節總結

- **處理器**：`func(w http.ResponseWriter, r *http.Request)`；`http.Handler` 介面只有 `ServeHTTP`
- **路由**：`http.NewServeMux()`；Go 1.22+ 樣式 `"GET /products/{id}"`、`r.PathValue("id")`、`"/{$}"`
- **參數**：查詢參數 `r.URL.Query().Get`；表單 `r.FormValue` / `r.PostFormValue`；一律驗證、限制大小
- **回應**：先設定標頭 → `w.WriteHeader(狀態碼)` → 寫本體；`http.Error`、`http.Redirect`（PRG）
- **網頁**：`html/template` 自動跳脫防 XSS；`http.FileServer` + `StripPrefix`；`//go:embed` 打包靜態資源
- **RESTful API**：資源用網址、動作用方法；`writeJSON` / `writeError`；中介軟體記錄日誌
- **正式環境**：`http.Server` 設定逾時；`signal.NotifyContext` + `srv.Shutdown` 優雅關閉；`httptest` 寫測試
- **GoShop**：`web` 套件提供商品與訂單的 RESTful JSON API、`html/template` 後台網頁與表單；`embed` 打包模板和 CSS；優雅關閉

下一章我們會介紹「並行性運算」：goroutine、通道與 context。

<!--
我們來整理今天學到的東西。

HTTP 伺服器的核心是處理器和路由。Go 1.22 之後的 ServeMux 支援方法和路徑參數，大部分的專案不需要第三方框架。產生網頁用 html/template，它會自動防止 XSS；提供靜態檔案用 FileServer，搭配 embed 可以打包成單一執行檔。RESTful API 用網址表示資源、用方法表示動作，搭配 JSON 交換資料。正式環境記得設定逾時和優雅關閉。

GoShop 這一步變成了一個網站：前台可以用 JSON API 查商品、下單、付款；店長可以在後台網頁看庫存和訂單、用表單新增商品。錯誤依照種類對應到 404、409、400 等狀態碼；模板和 CSS 用 embed 編譯進執行檔，部署時只要一個檔案。

今天我們好幾次提到「每個請求都在自己的 goroutine 裡執行」、「多個請求同時修改資料要用互斥鎖」。下一章就要正式學習 Go 最有特色的功能：並行性運算。goroutine、通道、互斥鎖、context，這些是 Go 被稱為「雲端時代的語言」的關鍵。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天的內容很多，但這一章學完，大家已經能用 Go 寫出一個真正的後端服務了。

課後建議：把今天的商品 API 的記憶體 Store，換成第 13 章的 MySQL Store，完成一個真正存在資料庫裡的 RESTful API。這個練習會把第 11、13、15 章的內容全部串起來。

有問題的同學現在可以提問！
-->

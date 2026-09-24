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
title: 使用 Go 的 HTTP 客戶端
routeAlias: ch14
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
  <h1 style="color: #1a5c5c; font-size: 3.2rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">使用 Go 的 HTTP 客戶端</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「發出請求、接住回應 — 呼叫全世界的 API」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十四章！

現代的應用程式很少是單打獨鬥的：要查天氣會呼叫氣象局的 API、要收款會呼叫金流的 API、要登入可能會呼叫 Google 的 API。這些 API 幾乎都是透過 HTTP 協定溝通，而負責「發出請求」的那一方，就叫做 HTTP 客戶端。

今天要學的是 Go 標準函式庫 net/http 的客戶端功能：發送 GET 和 POST 請求、解析 JSON 回應、上傳檔案、加上自訂標頭。完全不需要安裝任何第三方套件。
-->

---
layout: default
---

# Outline

- **前言** — HTTP 的請求與回應
- **Go 語言的 HTTP 客戶端** — `http.Client`、逾時設定
- **傳送 GET 請求** — `http.Get`、讀取回應、解析 JSON、查詢參數
- **傳送 POST 請求** — 送出 JSON、上傳檔案（multipart）
- **自訂標頭** — `http.NewRequestWithContext`、`client.Do`
- **章節總結**

<!--
今天的內容會先快速複習 HTTP 協定的基本觀念，接著認識 Go 的 http.Client，然後依序學 GET、POST、上傳檔案，最後學怎麼自訂請求的標頭。

這一章的範例會呼叫兩個網站：GitHub 的公開 API，以及 httpbin.org。httpbin 是一個專門給開發者測試 HTTP 請求的網站，我們送什麼過去，它就把收到的內容原封不動回傳給我們，很適合用來觀察請求長什麼樣子。
-->

---

# 回顧：SQL 與資料庫

- `*sql.DB` 是**連線池**：建立一次、全程共用
- 所有操作使用 **`XxxContext`** 版本，搭配 `context.WithTimeout` 避免無限等待
- 查詢完一定要 **`defer rows.Close()`**，否則連線不會歸還
- 第 11 章：`json.NewDecoder(r).Decode(&v)` 可以直接從 `io.Reader` 解碼 ← 今天的 HTTP 回應就是 `io.Reader`

<!--
回顧一下上一章。

我們學了資料庫的操作，有幾個觀念今天會以不同的形式再出現：*sql.DB 要重複使用，今天的 http.Client 也是；資料庫操作要設定逾時，HTTP 請求也要；rows 用完要 Close，HTTP 回應的 Body 也要 Close。

另外第 11 章的 json.NewDecoder 可以從任何 io.Reader 解碼，今天的 HTTP 回應本體就是一個 io.Reader，所以解析 API 回傳的 JSON 只要一行。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## HTTP Requests & Responses

<!--
先快速複習 HTTP 協定。
-->

---

# HTTP 的請求與回應

```text
客戶端（Go 程式）                               伺服器（API）
  ── 請求 Request ──────────────────────────▶
     GET /repos/golang/go HTTP/1.1              ← 方法、路徑、版本
     Host: api.github.com                       ← 標頭（headers）
     Accept: application/json
                                                ← 本體（GET 通常沒有）
  ◀────────────────────────── 回應 Response ──
     HTTP/1.1 200 OK                            ← 狀態碼
     Content-Type: application/json             ← 標頭
     {"full_name":"golang/go", ...}             ← 本體
```

<!--
HTTP 是一問一答的協定：客戶端送出請求，伺服器回傳回應。

請求由三個部分組成：第一行是方法（GET、POST）和路徑；接著是標頭，一行一個「名稱: 值」，用來傳遞額外的資訊，例如想要什麼格式的回應；最後是本體，放要送出的資料，GET 請求通常沒有本體。

回應也有三個部分：狀態碼、標頭、本體。狀態碼告訴我們結果如何，本體就是我們要的資料，通常是 JSON。

用寄信比喻：請求是我們寄出的信，標頭是信封上的資訊，本體是信的內容；回應就是對方的回信。
-->

---

# HTTP 方法與狀態碼

| 方法 | 用途 | | 狀態碼 | 意義 |
| --- | --- | --- | --- | --- |
| `GET` | 取得資料 | | `200 OK` | 成功 |
| `POST` | 新增資料 | | `201 Created` | 新增成功 |
| `PUT` / `PATCH` | 更新資料（整筆 / 部分） | | `400 Bad Request` | 請求格式錯誤 |
| `DELETE` | 刪除資料 | | `401` / `403` | 未登入 / 沒有權限 |
| | | | `404 Not Found` | 找不到資源 |
| | | | `500` / `503` | 伺服器錯誤 / 暫時無法服務 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>狀態碼的分類：</b> <code>2xx</code> 成功、<code>3xx</code> 重新導向、<code>4xx</code> 客戶端的錯、<code>5xx</code> 伺服器的錯。
</div>

<!--
HTTP 方法代表「想做什麼事」：GET 取得、POST 新增、PUT 和 PATCH 更新、DELETE 刪除。這四種剛好對應上一章資料庫的 CRUD，下一章寫 RESTful API 的時候會看到它們的對應關係。

狀態碼是三位數字，第一個數字代表分類：2 開頭是成功，4 開頭是客戶端的錯，例如網址打錯、沒有登入；5 開頭是伺服器的錯。

記住這個分類很有用：看到 4xx，先檢查自己送出的請求；看到 5xx，問題在對方的伺服器。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# Go 語言的 HTTP 客戶端
## http.Client

<!--
接下來認識 Go 的 HTTP 客戶端。
-->

---

# http.Client

| 用法 | 說明 |
| --- | --- |
| `http.Get(url)` / `http.Post(...)` | 便利函式，使用 `http.DefaultClient` |
| `client := &http.Client{Timeout: 10 * time.Second}` | **自訂客戶端**，設定逾時 |
| `client.Do(req)` | 送出自己建立的 `*http.Request`（可以設定方法、標頭） |

```go
var client = &http.Client{ // 套件層級：建立一次、全程共用
	Timeout: 10 * time.Second, // 整個請求（連線 + 傳送 + 讀取）的時間上限
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>DefaultClient 沒有逾時！</b> 對方伺服器不回應時，程式會永遠卡住。正式的程式請一律使用設定了 <code>Timeout</code> 的自訂 client。
</div>

<!--
Go 的 HTTP 客戶端是 http.Client。

最簡單的用法是 http.Get、http.Post 這些便利函式，它們使用內建的 DefaultClient。但這裡有一個大陷阱：DefaultClient 沒有設定逾時！如果對方的伺服器當掉、一直不回應，我們的程式就會永遠卡在那裡。

所以正式的程式，一定要自己建立一個 http.Client，設定 Timeout。Timeout 是整個請求的時間上限，包含建立連線、送出請求、讀取回應。

另外 http.Client 跟上一章的 *sql.DB 一樣，內部有連線池，會重複使用 TCP 連線，所以應該建立一次、全程共用，而且可以安全地被多個 goroutine 同時使用。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 對伺服器傳送 GET 請求
## GET Requests

<!--
先從最常用的 GET 請求開始。
-->

---

# 使用 http.Get() 發送 GET 請求

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://api.github.com/repos/golang/go")
	if err != nil { // 網路錯誤：DNS 查不到、連線失敗、逾時…
		fmt.Println("請求失敗：", err)
		return
	}
	defer resp.Body.Close() // ⚠️ 一定要關閉回應本體

	fmt.Println(resp.StatusCode, resp.Status) // 200 200 OK
	fmt.Println(resp.Header.Get("Content-Type"))
	body, err := io.ReadAll(resp.Body) // Body 是 io.ReadCloser
	if err != nil {
		fmt.Println("讀取失敗：", err)
		return
	}
	fmt.Println(len(body), "bytes")
}
```

<!--
http.Get 送出一個 GET 請求，回傳 *http.Response 和錯誤。

這裡的 err 只代表「網路層面的錯誤」：網址不存在、連不上、逾時。如果伺服器有回應，即使是 404 或 500，err 也是 nil。這一點等一下會特別說明。

拿到回應之後，第一件事是 defer resp.Body.Close()。這跟上一章 rows.Close() 的道理一樣：不關閉的話，底層的 TCP 連線就不會被歸還，久了會耗盡資源。

resp.StatusCode 是數字的狀態碼，resp.Status 是文字的狀態；resp.Header.Get 讀取回應標頭；resp.Body 是一個 io.ReadCloser，可以用 io.ReadAll 一次讀完。
-->

---

# 使用 http.Get 的注意事項：檢查狀態碼

**`err == nil` 不代表成功！** 伺服器回傳 404、500 時，`err` 仍然是 `nil`

```go
package main

import (
	"fmt"
	"net/http"
)

func fetch(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("請求 %s 失敗：%w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { // 用常數，不要寫魔術數字 200
		return fmt.Errorf("請求 %s 失敗：狀態碼 %s", url, resp.Status)
	}
	return nil
}

func main() {
	err := fetch("https://api.github.com/repos/golang/no-such-repo")
	fmt.Println(err) // 請求 … 失敗：狀態碼 404 Not Found
}
```

<!--
使用 http.Get 最重要的注意事項：err 是 nil，不代表請求成功。

Go 的設計是：只要伺服器有回應，就不算錯誤，不管狀態碼是多少。這跟很多語言的 HTTP 函式庫不一樣，初學者常常以為 err 是 nil 就可以放心解析資料，結果拿到的是一段錯誤訊息的 HTML。

所以一定要自己檢查 resp.StatusCode。建議用 net/http 提供的常數，例如 http.StatusOK、http.StatusNotFound，比直接寫 200、404 更好讀。

查詢一個不存在的 repository，GitHub 會回傳 404，我們把它轉換成一個清楚的錯誤訊息。
-->

---

# 取得並解析伺服器的 JSON 資料

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Repo struct {
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	PushedAt    time.Time `json:"pushed_at"`
}

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/golang/go")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var repo Repo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		panic(err)
	}
	fmt.Printf("%s ⭐ %d（%s）\n",
		repo.FullName, repo.Stars, repo.Description)
}
```

<!--
呼叫 API 最常見的情況，是回應的本體是 JSON，要把它解析成 Go 的 struct。

我們只定義需要的欄位：GitHub 回傳的 JSON 有上百個欄位，但 struct 裡只放四個，第 11 章說過，多出來的鍵會被自動忽略。pushed_at 是 RFC3339 格式，直接用 time.Time 接住。

解析的部分一行就好：json.NewDecoder(resp.Body).Decode(&repo)。因為 resp.Body 是 io.Reader，可以直接交給 Decoder，不需要先 ReadAll。這就是第 7 章「小介面、大組合」的威力。

這裡改用自訂的 client，設定了 10 秒逾時。為了投影片精簡，省略了狀態碼檢查，實際寫程式時要加上。
-->

---

# 加上查詢參數：url.Values

用 `url.Values` 組合查詢字串，會自動處理**編碼**（空白、中文、特殊符號）

```go
package main

import (
	"fmt"
	"net/url"
)

func main() {
	q := url.Values{}
	q.Set("q", "language:go stars:>1000") // 空白和 > 會被正確編碼
	q.Set("sort", "stars")
	q.Set("per_page", "3")

	u := url.URL{
		Scheme:   "https",
		Host:     "api.github.com",
		Path:     "/search/repositories",
		RawQuery: q.Encode(), // 鍵會依字母排序
	}
	fmt.Println(u.String())
	// https://api.github.com/search/repositories?per_page=3
	//   &q=language%3Ago+stars%3A%3E1000&sort=stars
}
```

<!--
很多 API 會用「查詢參數」接收條件，也就是網址問號後面的 key=value。

不要自己用字串串接查詢參數！因為參數值可能包含空白、中文、& 這些特殊字元，必須經過「URL 編碼」才能放進網址。自己串接很容易出錯，甚至造成安全問題。

url.Values 是 Go 處理查詢參數的標準工具，Set 設定參數，Encode 產生編碼過的查詢字串。再搭配 url.URL 組出完整的網址，每個部分都清清楚楚。

注意輸出裡，冒號變成 %3A，空白變成加號，大於符號變成 %3E，這就是 URL 編碼。
-->

---
layout: default
---

# 練習 1：查詢熱門的 Go 專案
### 任務說明

呼叫 GitHub 搜尋 API：`https://api.github.com/search/repositories`

1. 建立 `http.Client`，設定 **10 秒逾時**
2. 用 `url.Values` 組合查詢參數：`q=language:go`、`sort=stars`、`per_page=5`
3. 檢查狀態碼，不是 `200` 就回傳錯誤
4. 回應 JSON 的結構如下，定義對應的 struct 並解析：

```json
{ "total_count": 1234567,
  "items": [
    { "full_name": "golang/go", "stargazers_count": 131245 }, …
  ] }
```

5. 印出前 5 名的專案名稱與星星數

<!--
這個練習把 GET 請求的所有步驟串起來：建立有逾時的 client、組合查詢參數、檢查狀態碼、解析巢狀的 JSON。

回應的 JSON 有兩層：外層有 total_count 和 items，items 是一個陣列，每個元素是一個專案。想想看 struct 要怎麼定義？第 11 章學過解碼到複合結構。
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SearchResult struct {
	Total int `json:"total_count"`
	Items []struct {
		FullName string `json:"full_name"`
		Stars    int    `json:"stargazers_count"`
	} `json:"items"`
}

var client = &http.Client{Timeout: 10 * time.Second}
```

<!--
SearchResult 的 Items 用了匿名結構的切片，只取需要的兩個欄位。

client 宣告成套件層級的變數，整個程式共用同一個。
-->

---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
func searchGo(n int) (*SearchResult, error) {
	q := url.Values{"q": {"language:go"}, "sort": {"stars"}}
	q.Set("per_page", fmt.Sprint(n))
	u := "https://api.github.com/search/repositories?" + q.Encode()
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("搜尋失敗：%s", resp.Status)
	}
	var r SearchResult
	return &r, json.NewDecoder(resp.Body).Decode(&r)
}

func main() {
	r, err := searchGo(5)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, it := range r.Items {
		fmt.Printf("%d. %-24s ⭐ %d\n", i+1, it.FullName, it.Stars)
	}
}
```

<!--
url.Values 本身是 map[string][]string，所以可以直接用 map 常值建立，也可以用 Set 設定。per_page 需要字串，用 fmt.Sprint 轉換。

searchGo 依序做：送出請求、defer 關閉、檢查狀態碼、解析 JSON。最後一行 return &r, json.NewDecoder(...).Decode(&r) 是一個簡潔的寫法：Decode 的錯誤直接當作回傳值。

main 用 %-24s 讓專案名稱靠左對齊，這是第 9 章學的格式化技巧。

補充：GitHub API 未登入時每小時限制 60 次請求，練習時不要寫無限迴圈一直呼叫。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 用 POST 請求傳送資料給伺服器
## POST Requests

<!--
接下來學 POST 請求：把資料送給伺服器。
-->

---

# 送出 POST 請求並接收回應

語法：`http.Post(url, contentType, body io.Reader)`

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Order struct {
	Item string `json:"item"`
	Qty  int    `json:"qty"`
}

func main() {
	// 1. 編碼成 JSON
	payload, _ := json.Marshal(Order{Item: "咖啡", Qty: 2})
	resp, err := http.Post("https://httpbin.org/post",
		// 2. 本體是 io.Reader
		"application/json", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var echo struct {
		// httpbin 會把收到的 JSON 放在 "json" 欄位回傳
		JSON Order `json:"json"`
	}
	json.NewDecoder(resp.Body).Decode(&echo)
	fmt.Println(resp.Status, echo.JSON) // 200 OK {咖啡 2}
}
```

<!--
POST 請求要送出資料，資料放在請求的本體裡。

http.Post 需要三個參數：網址、內容類型、本體。送 JSON 的時候，內容類型是 application/json，告訴伺服器「本體是 JSON 格式」。本體的型別是 io.Reader，所以先用 json.Marshal 編碼成 []byte，再用 bytes.NewReader 包裝成 Reader。

httpbin.org/post 會把收到的內容原封不動回傳，送出的 JSON 會放在回應的 json 欄位裡。我們解析出來，確認伺服器確實收到了我們送的訂單。

這裡用了一個匿名結構 echo，只取回應中需要的欄位。
-->

---

# 送出表單資料：http.PostForm

網頁表單的格式是 `application/x-www-form-urlencoded`，用 `url.Values` 組成

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	form := url.Values{}
	form.Set("username", "gopher")
	form.Set("message", "你好 & 歡迎")
	resp, err := http.PostForm("https://httpbin.org/post", form)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var echo struct {
		Form map[string]any `json:"form"`
	}
	json.NewDecoder(resp.Body).Decode(&echo)
	fmt.Println(echo.Form["message"]) // 你好 & 歡迎
}
```

<!--
除了 JSON，另一種常見的 POST 格式是「表單」，也就是網頁上 form 送出的格式。

http.PostForm 接收一個 url.Values，會自動編碼成表單格式，並設定正確的內容類型。表單的編碼方式跟查詢參數一樣，中文和 & 符號都會被正確處理。

httpbin 會把收到的表單放在 form 欄位回傳。下一章寫 HTTP 伺服器的時候，會學到伺服器這一端怎麼讀取表單資料。
-->

---

# 用 POST 請求上傳檔案：multipart/form-data

```go
package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

func main() {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField("title", "月報")                        // 一般欄位
	part, _ := w.CreateFormFile("file", "report.csv")  // 檔案欄位
	strings.NewReader("月份,金額\n9,2520\n").WriteTo(part) // 寫入檔案內容
	// ⚠️ 一定要 Close，才會寫入結尾的分隔線
	w.Close()

	resp, err := http.Post("https://httpbin.org/post",
		w.FormDataContentType(), &body) // 含 boundary 的 Content-Type
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status) // 200 OK
}
```

<!--
上傳檔案要用另一種格式：multipart/form-data。它可以在一個請求裡同時放入一般的欄位和檔案，每個部分之間用一條特殊的「分隔線」（boundary）隔開。

Go 用 mime/multipart 套件產生這種格式。NewWriter 寫到一個 bytes.Buffer；WriteField 寫入一般欄位；CreateFormFile 建立一個檔案欄位，回傳一個 io.Writer，把檔案內容寫進去。實務上會用 io.Copy(part, f) 把 os.Open 開啟的檔案複製進去，這裡為了方便示範用字串代替。

使用 multipart 的注意事項：寫完一定要呼叫 w.Close()，它會寫入結尾的分隔線，少了它伺服器會解析失敗。另外 Content-Type 要用 w.FormDataContentType()，裡面包含了分隔線的字串。
-->

---
layout: default
---

# 練習 2：回報訂單
### 任務說明

1. 定義 `type Report struct { OrderID int; Items []string; Total int }`，加上 JSON 標籤
2. 寫函式 `postJSON(client *http.Client, url string, v any) (map[string]any, error)`：
   - 把 `v` 編碼成 JSON 送出
   - 狀態碼不是 `2xx` 時回傳錯誤
   - 把回應的 JSON 解析成 `map[string]any` 回傳
3. 用 `postJSON` 把報告送到 `https://httpbin.org/post`，印出回應中的 `json` 欄位
4. 再送到 `https://httpbin.org/status/500`，確認會得到錯誤

<!--
這個練習要大家寫一個可以重複使用的 postJSON 函式，這在實務上非常常見：幾乎每個專案都會有一個類似的工具函式，統一處理 JSON 的送出和錯誤。

httpbin.org/status/500 會固定回傳 500 狀態碼，很適合用來測試錯誤處理。判斷 2xx 可以用 resp.StatusCode/100 == 2。
-->

---

# 練習 2：解題提示
### 提示說明

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Report struct {
	OrderID int      `json:"order_id"`
	Items   []string `json:"items"`
	Total   int      `json:"total"`
}

func postJSON(c *http.Client, url string,
	v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("POST %s：%s", url, resp.Status)
	}
	var out map[string]any
	return out, json.NewDecoder(resp.Body).Decode(&out)
}
```

<!--
postJSON 依序做：編碼成 JSON、送出、defer 關閉、檢查狀態碼、解析回應。

判斷 2xx 的技巧：狀態碼除以 100，整數除法會捨去餘數，200 到 299 除以 100 都是 2。

回應解析成 map[string]any，因為我們不知道每個 API 的回應結構，這是第 11 章「處理內容未知的 JSON」的應用。
-->

---

# 練習 2：解題提示（續）
### 提示說明

```go
// 續上頁
func main() {
	c := &http.Client{Timeout: 10 * time.Second}
	r := Report{OrderID: 7, Items: []string{"咖啡", "蛋糕"}, Total: 300}

	out, err := postJSON(c, "https://httpbin.org/post", r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out["json"])
	// map[items:[咖啡 蛋糕] order_id:7 total:300]

	_, err = postJSON(c, "https://httpbin.org/status/500", r)
	// POST https://httpbin.org/status/500：500 Internal Server Error
	fmt.Println(err)
}
```

<!--
main 建立一個有逾時的 client，把報告送到 httpbin，印出它回傳的 json 欄位，確認伺服器收到的內容跟我們送的一樣。注意 order_id 是 JSON 的數字，解析成 any 會變成 float64，印出來是 7。

第二次送到 status/500，postJSON 發現狀態碼不是 2xx，回傳錯誤。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 在客戶端使用自訂標頭做為請求選項
## Custom Headers & Requests

<!--
最後學怎麼自訂請求：設定標頭、設定逾時、使用其他 HTTP 方法。
-->

---

# 建立自訂請求：http.NewRequestWithContext

需要**自訂標頭**、**使用 PUT / DELETE**、或**控制逾時與取消**時，自己建立 `*http.Request`

| 步驟 | 程式碼 |
| --- | --- |
| 1. 建立請求 | `req, err := http.NewRequestWithContext(ctx, "GET", url, nil)` |
| 2. 設定標頭 | `req.Header.Set("Authorization", "Bearer "+token)` |
| 3. 送出 | `resp, err := client.Do(req)` |

| 常見標頭 | 用途 |
| --- | --- |
| `Authorization` | 身分驗證，例如 `Bearer <token>`、API 金鑰 |
| `Accept` | 希望收到的格式，例如 `application/json` |
| `Content-Type` | 本體的格式（POST / PUT 時） |
| `User-Agent` | 客戶端的名稱與版本 |

<!--
http.Get 和 http.Post 很方便，但它們不能設定標頭，也只能用 GET 和 POST。需要更多控制的時候，就自己建立一個 *http.Request。

步驟是三步：用 NewRequestWithContext 建立請求，指定 context、方法、網址、本體；用 req.Header.Set 設定標頭；最後用 client.Do 送出。

表格列出了最常用的四種標頭。Authorization 是最重要的：幾乎所有需要登入的 API，都要在這個標頭放入金鑰或 token，第 18 章會再提到 token 的安全性。
-->

---

# 自訂標頭 — 範例

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://httpbin.org/headers", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "goshop-client/1.0")
	req.Header.Set("X-Request-Id", "a1b2c3")

	resp, err := http.DefaultClient.Do(req) // 逾時由 ctx 控制
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var echo struct{ Headers map[string]any }
	json.NewDecoder(resp.Body).Decode(&echo)
	fmt.Println(echo.Headers["User-Agent"])   // goshop-client/1.0
	fmt.Println(echo.Headers["X-Request-Id"]) // a1b2c3
}
```

<!--
這段程式碼的目的，是送出一個帶有自訂標頭的請求，並透過 httpbin 確認伺服器真的收到了這些標頭。

NewRequestWithContext 的第一個參數是 context，這裡設定了 5 秒逾時。方法用 http.MethodGet 常數，比直接寫字串 "GET" 不容易打錯。GET 沒有本體，所以最後一個參數是 nil。

接著用 Header.Set 設定三個標頭，X- 開頭的是自訂標頭的慣例寫法。然後用 client.Do 送出。這裡用 DefaultClient 也沒關係，因為逾時已經由 context 控制了。

httpbin.org/headers 會把收到的標頭回傳，印出來可以看到 goshop-client/1.0 和 a1b2c3。
-->

---

# 使用自訂請求的注意事項

**注意事項之一：** 需要取消或逾時時，用 **`NewRequestWithContext`**（`NewRequest` 沒有 context）

**注意事項之二：** 回應本體**一定要讀完並關閉**，連線才能被重複使用；讀取時加上**大小上限**

```go
const maxBody = 1 << 20 // 1 MB
// 避免被超大回應灌爆記憶體
body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
```

**注意事項之三：** `client.Timeout` 與 `context` 逾時同時存在時，**先到的先生效**

| 設定 | 範圍 |
| --- | --- |
| `http.Client{Timeout: ...}` | 這個 client 的**所有請求**的上限 |
| `context.WithTimeout(...)` | **單一請求**的上限，也能手動取消（Ch 16） |

<!--
使用自訂請求有三個注意事項。

第一，建立請求一律用 NewRequestWithContext，而不是舊的 NewRequest，這樣才能控制逾時和取消。第 16 章學 context 的時候會看到，一個 HTTP 伺服器處理請求時，如果使用者關掉了瀏覽器，可以透過 context 把後續的 API 呼叫一起取消。

第二，回應本體一定要關閉。另外，如果呼叫的是不信任的外部 API，讀取時加上 io.LimitReader 限制大小，避免對方回傳一個好幾 GB 的回應，把我們的記憶體灌爆。

第三，client 的 Timeout 和 context 的逾時可以同時使用，先到的先生效。client 的 Timeout 是一個全域的保險，context 用來控制個別的請求。
-->

---
layout: default
---

# 綜合練習：API 客戶端封裝
### 任務說明

寫一個可以重複使用的 GitHub API 客戶端：

1. 定義 `type GitHubClient struct { http *http.Client; token string }` 與 `NewGitHubClient(token string) *GitHubClient`（逾時 10 秒）
2. 寫方法 `(c *GitHubClient) get(ctx context.Context, path string, v any) error`：
   - 組出 `https://api.github.com` + `path`
   - 設定 `Accept: application/vnd.github+json`；`token` 不是空字串時設定 `Authorization: Bearer <token>`
   - 狀態碼是 `404` 時回傳哨兵錯誤 `ErrNotFound`，其他非 `200` 回傳一般錯誤
   - 解析 JSON 到 `v`
3. 寫方法 `Repo(ctx, owner, name string) (*Repo, error)`
4. 查詢 `golang/go` 與不存在的 `golang/nope`，用 `errors.Is` 判斷後者

<!--
這個綜合練習是實務上非常常見的設計：把某個外部 API 封裝成一個客戶端型別。

好處是：設定標頭、檢查狀態碼、解析 JSON 這些重複的工作，全部集中在 get 這個私有方法裡；公開的 Repo 方法只需要一行呼叫 get。之後要加上新的 API，例如查詢使用者、查詢 issue，都只要再寫一個短短的方法。

這個設計跟上一章的 Store 很像，都是把「怎麼跟外部系統溝通」封裝起來。
-->

---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("找不到資源")

type Repo struct {
	FullName string `json:"full_name"`
	Stars    int    `json:"stargazers_count"`
}

type GitHubClient struct {
	http  *http.Client
	token string
}

func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		http:  &http.Client{Timeout: 10 * time.Second},
		token: token,
	}
}
```

<!--
先定義哨兵錯誤 ErrNotFound、Repo 結構，以及 GitHubClient。

NewGitHubClient 是建構函式，建立一個設定好逾時的 http.Client。欄位都是小寫，外部只能透過方法使用，這是第 8 章的封裝。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func (c *GitHubClient) get(ctx context.Context,
	path string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.github.com"+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return fmt.Errorf("GET %s：%w", path, ErrNotFound)
	case resp.StatusCode != http.StatusOK:
		return fmt.Errorf("GET %s：%s", path, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}
```

<!--
get 是私有方法，負責所有共通的工作：建立請求、設定標頭、送出、檢查狀態碼、解析 JSON。

token 不是空字串才設定 Authorization 標頭。狀態碼用無條件 switch 判斷：404 回傳包裝過的 ErrNotFound，其他非 200 回傳一般錯誤。

參數 v 的型別是 any，呼叫端傳入任何 struct 的指標，都能解析進去。
-->

---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func (c *GitHubClient) Repo(ctx context.Context,
	owner, name string) (*Repo, error) {
	var r Repo
	if err := c.get(ctx, "/repos/"+owner+"/"+name, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func main() {
	gh := NewGitHubClient("") // 公開資料不需要 token
	ctx := context.Background()
	if r, err := gh.Repo(ctx, "golang", "go"); err == nil {
		fmt.Printf("%s ⭐ %d\n", r.FullName, r.Stars)
	}
	_, err := gh.Repo(ctx, "golang", "nope")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("查無此專案：", err)
	}
}
```

<!--
Repo 方法只需要組出路徑，交給 get 處理。之後要新增其他 API，例如查詢使用者，也只要寫一個類似的短方法。

main 查詢 golang/go，印出名稱和星星數；查詢不存在的 golang/nope，用 errors.Is 判斷是 ErrNotFound，即使它已經被 fmt.Errorf 包裝過一層。

實務上，token 應該從環境變數讀取，例如 os.Getenv("GITHUB_TOKEN")，不要寫死在程式碼裡。
-->

---

# 章節總結

- **HTTP 基礎**：請求 = 方法 + 路徑 + 標頭 + 本體；回應 = 狀態碼 + 標頭 + 本體；`2xx` 成功、`4xx` 客戶端錯、`5xx` 伺服器錯
- **http.Client**：**DefaultClient 沒有逾時**，自訂 `&http.Client{Timeout: ...}`，建立一次、全程共用
- **GET**：`client.Get(url)` → `defer resp.Body.Close()` → **檢查 `StatusCode`** → `json.NewDecoder(resp.Body).Decode(&v)`
- **查詢參數**：用 `url.Values` + `Encode()`，不要自己串接
- **POST**：`client.Post(url, "application/json", bytes.NewReader(b))`；表單用 `PostForm`；檔案用 `mime/multipart`（記得 `w.Close()`）
- **自訂請求**：`http.NewRequestWithContext` + `req.Header.Set` + `client.Do`；用 `io.LimitReader` 限制回應大小

下一章我們會換到另一端：用 Go 建立 **HTTP 伺服器**，提供網頁與 RESTful API。

<!--
我們來整理今天學到的東西。

HTTP 客戶端的標準流程是：建立有逾時的 client、送出請求、defer 關閉本體、檢查狀態碼、解析 JSON。記住兩個最常見的陷阱：DefaultClient 沒有逾時，以及 err 是 nil 不代表成功。需要自訂標頭或使用其他方法時，用 NewRequestWithContext 加上 client.Do。

今天我們是「呼叫別人的 API」的那一方。下一章要換到另一端，用 Go 建立 HTTP 伺服器，自己提供網頁和 API 給別人呼叫。Go 的 net/http 伺服器效能非常好，很多公司直接用標準函式庫就能撐起正式環境的服務。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
今天學的 HTTP 客戶端，是串接各種外部服務的基礎。

課後建議：到政府資料開放平台或中央氣象署開放資料平台申請一組 API 金鑰，用今天學的方法寫一個查詢天氣的小程式，練習在標頭或查詢參數中帶入金鑰。

有問題的同學現在可以提問！
-->

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
title: 編碼／解碼 JSON 資料
routeAlias: ch11
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
  <h1 style="color: #1a5c5c; font-size: 3.2rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">編碼／解碼 JSON 資料</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「struct 加上標籤，就能和全世界交換資料」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十一章！

現代的程式幾乎都在交換資料：手機 App 跟後端伺服器、前端網頁跟 API、微服務跟微服務之間。而這些資料交換，最主流的格式就是 JSON。

今天要學的是 Go 的 encoding/json 套件：怎麼把收到的 JSON 解碼成 Go 的 struct，以及怎麼把 Go 的 struct 編碼成 JSON。這一章學完，第 14 章的 HTTP 客戶端和第 15 章的 RESTful API 就會非常順手。

最後還會介紹 Go 1.27 正式推出的新版 JSON 套件 encoding/json/v2，以及 Go 自己的二進位編碼格式 gob。
-->

---
layout: default
---

# Outline

- **前言** — 什麼是 JSON、編碼與解碼
- **解碼 JSON 為 Go 結構** — `Unmarshal`、struct 標籤、複合結構
- **將 Go 結構編碼為 JSON** — `Marshal`、略過欄位、排版輸出
- **使用 Decoder / Encoder** — 以串流的方式處理 JSON
- **處理內容未知的 JSON** — 解碼成 `map[string]any`
- **gob** — Go 自有的編碼格式
- **補充：encoding/json/v2** — Go 1.27 的新版 JSON 套件
- **GoShop 專案實作** — 第 11 步：JSON 匯入商品、gob 存檔
- **章節總結**

<!--
今天的內容從解碼開始，因為實務上最常遇到的情況是「收到一段 JSON，要把它變成 Go 的資料」。接著學反方向的編碼。

然後會學 Decoder 和 Encoder，它們可以直接從檔案、網路連線讀寫 JSON。最後處理結構不固定的 JSON，以及 Go 自己的 gob 格式。
-->

---

# 回顧：時間處理

- 參考時間 `2006-01-02 15:04:05`；API 傳遞時間用 **RFC3339**
- 內部一律使用 **UTC**，顯示時才用 `t.In(loc)` 轉換
- 比較時間用 `Equal`；`time.Duration` 表示時間長度
- 第 8 章：**大寫開頭的欄位才會匯出** ← 今天 JSON 編碼的關鍵規則

<!--
回顧一下上一章。

我們學了時間的格式化，特別提到 API 傳遞時間要用 RFC3339 格式。今天就會看到，time.Time 在 JSON 裡預設就是 RFC3339 格式。

另外請大家回想第 8 章的匯出規則：大寫開頭的名稱才會匯出。這條規則在今天非常重要，因為 JSON 套件只看得到匯出的欄位。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## What Is JSON?

<!--
先來認識 JSON。
-->

---

# 什麼是 JSON？

「**JSON（JavaScript Object Notation）是一種以純文字表示結構化資料的格式**，人和機器都容易讀寫。」

```json
{
  "id": 7,
  "name": "咖啡豆",
  "price": 450.5,
  "in_stock": true,
  "tags": ["熱銷", "新品"],
  "supplier": { "name": "好豆商行", "phone": null }
}
```

| JSON 型別 | 對應的 Go 型別 |
| --- | --- |
| 物件 `{}` | `struct`、`map[string]T` |
| 陣列 `[]` | 切片、陣列 |
| 字串 / 數字 / 布林 / `null` | `string` / `int`、`float64` / `bool` / `nil`（指標、切片、map） |

<!--
什麼是 JSON？

JSON 是一種純文字的資料格式，它最早來自 JavaScript，但現在幾乎所有語言都支援。它只有六種型別：物件（大括號）、陣列（中括號）、字串、數字、布林值和 null。

JSON 的物件對應到 Go 的 struct 或 map，陣列對應到切片。表格裡列出了對應關係。

我們把「Go 的資料轉成 JSON 文字」叫做編碼（encode），或稱為序列化（marshal）；反過來「JSON 文字轉成 Go 的資料」叫做解碼（decode），或反序列化（unmarshal）。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 解碼 JSON 為 Go 結構
## Unmarshal

<!--
先學解碼：把 JSON 變成 Go 的 struct。
-->

---
zoom: 0.93
---

# Unmarshal()

語法：`json.Unmarshal(data []byte, v any) error`，`v` 必須是**指標**

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	ID    int
	Name  string
	Price float64
}

func main() {
	data := []byte(`{"id": 7, "name": "咖啡豆", "price": 450.5}`)
	var p Product
	if err := json.Unmarshal(data, &p); err != nil { // 傳入指標
		fmt.Println("解碼失敗：", err)
		return
	}
	fmt.Printf("%+v\n", p) // {ID:7 Name:咖啡豆 Price:450.5}
}
```

<!--
json.Unmarshal 接收兩個參數：JSON 的位元組切片，以及一個「指標」，指向要存放結果的變數。

為什麼要傳指標？第一章學過：Go 傳參數一律複製一份，如果傳的是值，Unmarshal 修改的是複製品，我們的 p 不會有任何變化。傳指標，Unmarshal 才能把資料寫進我們的 p。

JSON 裡的鍵是小寫的 id、name，struct 的欄位是大寫的 ID、Name，為什麼能對應起來？因為 encoding/json 在比對名稱時「不分大小寫」。不過實務上不建議依賴這個行為，下一頁會用標籤明確指定。
-->

---
zoom: 0.96
---

# 加上結構 JSON 標籤（struct tag）

用**反引號**在欄位後面加上標籤，指定 JSON 中的**鍵名稱**與選項

```go
type Product struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	// Go 用 CamelCase，JSON 用 snake_case
	InStock  bool     `json:"in_stock"`
	Tags     []string `json:"tags,omitempty"` // 空值時編碼會省略這個欄位
	Password string   `json:"-"`              // 完全忽略這個欄位
	cost     float64  // 小寫開頭：未匯出，JSON 套件看不到
}
```

| 標籤 | 意義 |
| --- | --- |
| `json:"name"` | 指定 JSON 的鍵名稱 |
| `json:"name,omitempty"` | 值是零值、空切片、空 map 時，**編碼時省略** |
| `json:"-"` | 編碼、解碼都忽略 |

<!--
struct 標籤是寫在欄位型別後面、用反引號包起來的字串。它本身不影響程式的執行，而是讓其他套件「讀取」的附加資訊，JSON 套件就會讀取 json 這個標籤。第 19 章學 reflect 的時候，會知道套件是怎麼讀到這些標籤的。

最常用的是指定 JSON 的鍵名稱。Go 的慣例是 CamelCase，但 JSON 常用 snake_case，例如 in_stock，用標籤就能對應起來。

omitempty 讓空值的欄位在編碼時省略；減號代表完全忽略，常用在密碼這種不應該輸出的欄位。

最後注意 cost：小寫開頭，沒有匯出，JSON 套件完全看不到它，不管加什麼標籤都沒用。這是初學者最常犯的錯誤。
-->

---
zoom: 0.75
---

# 解碼 JSON 到複合結構

JSON 的巢狀物件與陣列，對應到 Go 的**巢狀 struct 與切片**

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}

type Order struct {
	ID       int      `json:"id"`
	Customer struct { // 匿名巢狀結構
		Name string `json:"name"`
	} `json:"customer"`
	Items []Item `json:"items"`
}

func main() {
	data := `{"id":1,"customer":{"name":"小明"},
		"items":[{"name":"咖啡","qty":2},{"name":"蛋糕","qty":1}]}`
	var o Order
	if err := json.Unmarshal([]byte(data), &o); err != nil {
		panic(err)
	}
	// 小明 2 咖啡
	fmt.Println(o.Customer.Name, len(o.Items), o.Items[0].Name)
}
```

<!--
真實世界的 JSON 通常是巢狀的：訂單裡面有顧客資料，還有一個商品的陣列。

Go 的做法很直覺：JSON 怎麼巢狀，struct 就怎麼巢狀。顧客是一個物件，對應到一個 struct；商品是一個陣列，對應到 []Item 切片。

Customer 用了第 4 章學的匿名結構，適合只在這裡用一次的結構。如果其他地方也會用到，就定義成獨立的型別。
-->

---
zoom: 0.85
---

# 使用 Unmarshal 的注意事項

**注意事項之一：** JSON 中**多出來的鍵會被忽略**；struct 中**缺少的鍵保持零值**

**注意事項之二：** **型別不符**會回傳錯誤；JSON 數字解碼到 `any` 時一律是 `float64`

```go
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID  int  `json:"id"`
	Age *int `json:"age"` // 用指標區分「沒有給」(nil) 和「給了 0」
}

func main() {
	var u User
	err := json.Unmarshal([]byte(`{"id":"abc"}`), &u)
	// json: cannot unmarshal string into
	// Go struct field User.id of type int
	fmt.Println(err)

	json.Unmarshal([]byte(`{"id":1,"extra":true}`), &u) // extra 被忽略
	// 1 true：沒有 age
	fmt.Println(u.ID, u.Age == nil)
}
```

<!--
使用 Unmarshal 有兩個注意事項。

第一，JSON 裡有、struct 裡沒有的鍵，會被默默忽略；struct 裡有、JSON 裡沒有的欄位，會保持零值。這個寬鬆的行為很方便，但也可能讓打錯字的鍵名稱被默默忽略。等一下會看到怎麼讓它變嚴格。

第二，型別不符會回傳錯誤，例如 JSON 裡的 id 是字串，struct 裡是 int。

另外有一個實用的技巧：如果需要區分「沒有給這個欄位」和「給了零值」，例如年齡 0 歲和沒填年齡，就把欄位宣告成指標。沒給的時候是 nil，給了 0 的時候是指向 0 的指標。
-->

---
layout: default
---

# 練習 1：解析天氣資料
### 任務說明

解析以下 JSON，印出每個城市的名稱、溫度，以及最高溫的城市：

```json
{
  "updated_at": "2026-09-24T14:00:00+08:00",
  "cities": [
    {"name": "台北", "temp_c": 31.5, "rain": true},
    {"name": "台中", "temp_c": 33.2, "rain": false},
    {"name": "高雄", "temp_c": 32.8}
  ]
}
```

1. 定義 `Report`、`City` 兩個 struct，使用 JSON 標籤
2. `updated_at` 解碼成 `time.Time`（JSON 套件預設支援 RFC3339）
3. 印出更新時間（格式 `01/02 15:04`）、每個城市，以及最高溫的城市

<!--
這個練習要大家自己定義 struct 並加上標籤。

注意兩個地方：updated_at 可以直接解碼成 time.Time，因為 time.Time 內建支援 RFC3339 格式的 JSON；高雄沒有 rain 這個鍵，解碼後會是什麼值？
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type City struct {
	Name  string  `json:"name"`
	TempC float64 `json:"temp_c"`
	Rain  bool    `json:"rain"`
}

type Report struct {
	UpdatedAt time.Time `json:"updated_at"`
	Cities    []City    `json:"cities"`
}
```

<!--
City 和 Report 兩個 struct，每個欄位都加上 JSON 標籤，對應 snake_case 的鍵。

UpdatedAt 的型別是 time.Time，JSON 套件會自動用 RFC3339 格式解析，而且保留 +08:00 的時區。
-->

---
zoom: 0.97
---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
func main() {
	data := []byte(`{"updated_at":"2026-09-24T14:00:00+08:00","cities":[
		{"name":"台北","temp_c":31.5,"rain":true},
		{"name":"台中","temp_c":33.2},
		{"name":"高雄","temp_c":32.8}]}`)
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		panic(err)
	}
	// 09/24 14:00
	fmt.Println("更新時間：", r.UpdatedAt.Format("01/02 15:04"))
	hottest := r.Cities[0]
	for _, c := range r.Cities {
		fmt.Printf("%s %.1f°C 下雨：%t\n", c.Name, c.TempC, c.Rain)
		if c.TempC > hottest.TempC {
			hottest = c
		}
	}
	fmt.Println("最高溫：", hottest.Name) // 台中
}
```

<!--
main 把 JSON 解碼到 Report，因為保留了 +08:00 的時區，格式化後是 14:00。

高雄沒有 rain 這個鍵，解碼後是 bool 的零值 false。最高溫的城市用一個變數記錄，走訪時比較，結果是台中。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 將 Go 結構編碼為 JSON
## Marshal

<!--
接下來是反方向：把 Go 的資料編碼成 JSON。
-->

---
zoom: 0.83
---

# Marshal()

語法：`json.Marshal(v any) ([]byte, error)`

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	p := Product{ID: 7, Name: "咖啡豆", Price: 450,
		CreatedAt: time.Date(2026, 9, 24, 14, 5, 0, 0, time.UTC)}
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// {"id":7,"name":"咖啡豆","price":450,
	//  "created_at":"2026-09-24T14:05:00Z"}
}
```

<!--
json.Marshal 把任何 Go 的值編碼成 JSON，回傳位元組切片和錯誤。要印出來或當成字串使用，用 string() 轉換。

欄位的順序跟 struct 定義的順序一樣，鍵名稱依照標籤。time.Time 自動轉成 RFC3339 格式的字串。

Marshal 什麼時候會回傳錯誤？例如值裡面有通道、函式這種 JSON 無法表示的型別，或是浮點數是無限大、NaN。一般的 struct 很少出錯，但還是要檢查。
-->

---

# 將有多重欄位的結構轉為 JSON

| Go 的值 | 編碼成 JSON |
| --- | --- |
| `nil` 切片 `[]string(nil)` | `null` |
| 空切片 `[]string{}` | `[]` |
| `map[string]int{"b": 2, "a": 1}` | `{"a":1,"b":2}`（**鍵會自動排序**） |
| `[]byte("hi")` | `"aGk="`（Base64 字串） |
| `"<a&b>"` | `"\u003ca\u0026b\u003e"`（預設跳脫 HTML 字元） |
| 內嵌 struct | 欄位被**提升**到外層物件 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>前端最常抱怨的問題：</b> 沒有資料的切片輸出 <code>null</code> 而不是 <code>[]</code>，JavaScript 的 <code>.length</code> 就會出錯。需要 <code>[]</code> 時，請初始化為空切片：<code>Tags: []string{}</code>。
</div>

<!--
這張表整理了幾個常見型別的編碼結果，有幾個細節要特別注意。

第一個是前端工程師最常抱怨的問題：nil 切片會編碼成 null，空切片才會編碼成空陣列。如果前端拿到 null 去讀 length，JavaScript 就會出錯。所以 API 回傳列表的時候，沒有資料要記得初始化成空切片。

另外 map 編碼時鍵會自動排序，[]byte 會變成 Base64 字串，HTML 的特殊字元預設會被跳脫，避免 JSON 直接放進網頁時造成安全問題。

內嵌的 struct，欄位會被提升到外層，跟第 4 章內嵌的行為一致。
-->

---
zoom: 0.89
---

# 略過欄位：omitempty 與 omitzero

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Profile struct {
	Name     string   `json:"name"`
	Nickname string   `json:"nickname,omitempty"` // "" 時省略
	Tags     []string `json:"tags,omitempty"`     // nil 或空切片時省略
	// Go 1.24+：零值時省略
	Birthday time.Time `json:"birthday,omitzero"`
	Password string    `json:"-"` // 永遠不輸出
}

func main() {
	b, _ := json.Marshal(Profile{Name: "Alice", Password: "secret"})
	fmt.Println(string(b)) // {"name":"Alice"}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>為什麼需要 omitzero？</b> <code>omitempty</code> 對 struct 無效（<code>time.Time</code> 零值會輸出 <code>"0001-01-01T00:00:00Z"</code>）；<code>omitzero</code> 會判斷<b>任何型別</b>是否為零值，或呼叫它的 <code>IsZero()</code> 方法。
</div>

<!--
略過欄位有三種標籤。

omitempty 是最傳統的：值是零值、空字串、空切片、空 map 的時候，編碼時省略這個欄位。

但 omitempty 有一個長年被詬病的缺點：它對 struct 無效。time.Time 是一個 struct，就算是零值也會輸出「西元 1 年 1 月 1 日」，非常奇怪。Go 1.24 加入了 omitzero，它對任何型別都有效，如果型別有 IsZero 方法，還會呼叫它來判斷。所以時間欄位請用 omitzero。

減號則是永遠不輸出，像密碼這種欄位一定要加上，避免不小心透過 API 洩漏出去。
-->

---
zoom: 0.75
---

# 有排版的 JSON 編碼結果：MarshalIndent

語法：`json.MarshalIndent(v, 前綴, 縮排字串)`

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Host  string   `json:"host"`
	Port  int      `json:"port"`
	Debug bool     `json:"debug"`
	Hosts []string `json:"allowed_hosts"`
}

func main() {
	cfg := Config{"localhost", 8080, true, []string{"a.com", "b.com"}}
	b, _ := json.MarshalIndent(cfg, "", "  ") // 每一層縮排兩個空白
	fmt.Println(string(b))
}
```

```json
{
  "host": "localhost",
  "port": 8080,
  "debug": true,
  "allowed_hosts": [
    "a.com",
    "b.com"
  ]
}
```

<!--
Marshal 輸出的 JSON 是擠成一行的，適合在網路上傳輸，因為比較小。但如果要給人看，例如設定檔、除錯輸出，就用 MarshalIndent，它會自動換行和縮排。

第二個參數是每一行的前綴，通常是空字串；第三個參數是每一層的縮排，常用兩個空白或 Tab。
-->

---
layout: default
---

# 練習 2：產生 API 回應
### 任務說明

1. 定義 API 回應的結構：
   - `Response`：`Success bool`、`Data any`（`omitempty`）、`Error string`（`omitempty`）
   - `User`：`ID int`、`Name string`、`Email string`（`omitempty`）、`Password string`（不輸出）、`CreatedAt time.Time`（`omitzero`）
2. 寫函式 `toJSON(v any) string`，用 `MarshalIndent` 回傳排版後的字串（錯誤時回傳錯誤訊息）
3. 產生兩種回應並印出：
   - 成功：`Data` 是兩個使用者的切片（其中一個沒有 Email）
   - 失敗：`Error` 是「找不到使用者」，`Data` 是 `nil`

<!--
這個練習模擬第 15 章要寫的 API 回應。

實務上很多團隊會定義一個統一的回應格式，成功時放 data，失敗時放 error。用 omitempty 讓不需要的欄位不要出現。

注意 Data 的型別是 any，可以放任何東西，這是第 7 章空介面的實際應用。
-->

---
zoom: 0.77
---

# 練習 2：解題提示
### 提示說明

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}

func toJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "編碼失敗：" + err.Error()
	}
	return string(b)
}
```

<!--
Response 的 Data 和 Error 都加上 omitempty，成功時沒有 error、失敗時沒有 data。

User 的 Password 用減號永遠不輸出，CreatedAt 用 omitzero，沒設定時間就不輸出。

toJSON 包裝了 MarshalIndent，編碼失敗時回傳錯誤訊息，這樣呼叫端寫起來比較簡潔。
-->

---

# 練習 2：解題提示（續）
### 提示說明

```go
// 續上頁
func main() {
	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Password: "x",
			CreatedAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Bob"},
	}
	fmt.Println(toJSON(Response{Success: true, Data: users}))
	fmt.Println(toJSON(Response{Success: false, Error: "找不到使用者"}))
}
```

```json
{
  "success": false,
  "error": "找不到使用者"
}
```

<!--
成功的回應，Data 放入使用者切片。Bob 沒有 Email 和建立時間，所以這兩個欄位不會出現；兩個人的 Password 都不會出現。

失敗的回應只有 success 和 error 兩個欄位，data 因為是 nil 被省略了。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 使用 Decoder/Encoder 處理 JSON 資料
## Streaming JSON

<!--
接下來看 Decoder 和 Encoder，它們可以直接從 io.Reader 讀、寫到 io.Writer。
-->

---
zoom: 0.97
---

# 使用 Decoder / Encoder 處理 JSON 資料

| 方式 | 輸入／輸出 | 適合的情境 |
| --- | --- | --- |
| `json.Unmarshal` / `json.Marshal` | `[]byte` | 資料已經在記憶體中 |
| `json.NewDecoder(r).Decode(&v)` | 任何 `io.Reader` | **HTTP 請求本體**、檔案、網路連線 |
| `json.NewEncoder(w).Encode(v)` | 任何 `io.Writer` | **HTTP 回應**、檔案、`os.Stdout` |

```go
package main

import (
	"encoding/json"
	"os"
)

func main() {
	enc := json.NewEncoder(os.Stdout) // 直接寫到螢幕，不需要先轉成 []byte
	enc.SetIndent("", "  ")
	enc.Encode(map[string]any{"status": "ok", "count": 3})
}
```

<!--
Marshal 和 Unmarshal 處理的是記憶體中的 []byte。但很多時候，JSON 資料是從檔案、網路連線「流」進來的，這時候用 Decoder 和 Encoder 更方便。

NewDecoder 接收任何 io.Reader，NewEncoder 接收任何 io.Writer，這就是第 7 章說的「接受介面」的威力：同一個 Encoder，可以寫到螢幕、寫到檔案、寫到 HTTP 回應。

第 15 章寫 HTTP 伺服器時，讀取請求的 JSON 用 json.NewDecoder(r.Body).Decode，回傳 JSON 用 json.NewEncoder(w).Encode，這是最主流的寫法。

注意 Encode 會在最後自動加上換行。
-->

---

# Decoder 的嚴格模式：DisallowUnknownFields

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Event struct {
	Type string `json:"type"`
	User string `json:"user"`
}

func main() {
	bad := `{"type":"login","usr":"alice"}` // user 打成 usr
	dec := json.NewDecoder(strings.NewReader(bad))
	dec.DisallowUnknownFields() // 遇到未知的鍵就回傳錯誤
	var e Event
	fmt.Println(dec.Decode(&e)) // json: unknown field "usr"
}
```

<!--
Decoder 有兩個 Unmarshal 做不到的功能，第一個是嚴格模式。

DisallowUnknownFields 開啟嚴格模式：遇到 struct 裡沒有的鍵，就回傳錯誤。剛剛說過 Unmarshal 會默默忽略多出來的鍵，打錯字也不知道；嚴格模式可以抓到這種錯誤，很適合用在 API 的請求驗證。這裡故意把 user 打成 usr，就被抓到了。
-->

---
zoom: 0.77
---

# Decoder 讀取多筆 JSON（JSON Lines）

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Event struct {
	Type string `json:"type"`
	User string `json:"user"`
}

func main() {
	logs := `{"type":"login","user":"a"}
{"type":"logout","user":"b"}`
	dec := json.NewDecoder(strings.NewReader(logs))
	for {
		var ev Event
		err := dec.Decode(&ev)
		if errors.Is(err, io.EOF) { // 讀完了
			break
		} else if err != nil {
			fmt.Println("格式錯誤：", err)
			break
		}
		fmt.Println(ev.Type, ev.User)
	}
}
```

<!--
Decoder 的第二個功能：可以連續讀取多筆 JSON。

很多日誌系統用 JSON Lines 格式，一行一筆 JSON，第 9 章 slog 的 JSON 輸出就是這種格式。用迴圈呼叫 Decode，讀到 io.EOF 代表資料讀完了。io.EOF 就是第 6 章說的哨兵錯誤。

執行後會印出 login a 和 logout b。
-->
---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 處理內容未知的 JSON 資料
## Dynamic JSON

<!--
有時候我們不知道 JSON 長什麼樣子，沒辦法事先定義 struct，這時候怎麼辦？
-->

---
zoom: 0.97
---

# 將 JSON 格式解碼成 map

解碼到 `map[string]any`，再用**型別斷言**取出值

| JSON 型別 | 解碼到 `any` 時的 Go 型別 |
| --- | --- |
| 物件 | `map[string]any` |
| 陣列 | `[]any` |
| 數字 | **`float64`**（不是 `int`！） |
| 字串 / 布林 / `null` | `string` / `bool` / `nil` |

```go
var m map[string]any
raw := `{"name":"Go","version":1.27,"tags":["fast"]}`
json.Unmarshal([]byte(raw), &m)
name := m["name"].(string)
ver := m["version"].(float64)
tags := m["tags"].([]any)
fmt.Println(name, ver, tags[0]) // Go 1.27 fast
```

<!--
當 JSON 的結構不固定，例如第三方 API 的回應欄位會變動、使用者自訂的設定，就沒辦法事先定義 struct。這時候可以解碼到 map[string]any。

表格是解碼到 any 時的型別對應。要特別注意：所有的數字都會變成 float64，即使 JSON 裡寫的是整數。

取出值的時候，要用第 4 章學的型別斷言，把 any 轉回具體的型別。巢狀的物件是 map[string]any，陣列是 []any，一層一層斷言下去。

可以看得出來，這種寫法比用 struct 囉嗦很多，而且失去了型別檢查，所以只在真的不知道結構的時候才使用。
-->

---
zoom: 0.75
---

# 走訪內容未知的 JSON

```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func walk(v any, indent int) {
	pad := strings.Repeat("  ", indent)
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			fmt.Printf("%s%s:\n", pad, k)
			walk(val, indent+1)
		}
	case []any:
		for i, val := range x {
			fmt.Printf("%s[%d]\n", pad, i)
			walk(val, indent+1)
		}
	default:
		fmt.Printf("%s%v (%T)\n", pad, x, x)
	}
}

func main() {
	var data any
	raw := `{"user":{"name":"Amy","roles":["admin"]},"age":30}`
	json.Unmarshal([]byte(raw), &data)
	walk(data, 0)
}
```

<!--
這段程式碼的目的，是走訪一個完全未知結構的 JSON，把每一層的內容印出來。

walk 函式用型別 switch 判斷目前的值：是物件就走訪每個鍵，是陣列就走訪每個元素，都遞迴呼叫自己；其他的就是字串、數字、布林這些基本值，直接印出值和型別。

這是型別 switch 和遞迴的組合應用。因為 map 的走訪順序不固定，每次執行的輸出順序可能不一樣。
-->

---
zoom: 0.9
---

# 將 map 編碼成 JSON 格式

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	event := map[string]any{
		"type":    "purchase",
		"user_id": 42,
		"items":   []string{"咖啡", "蛋糕"},
		"meta":    map[string]any{"coupon": nil, "channel": "app"},
	}
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// {"items":["咖啡","蛋糕"],"meta":{"channel":"app","coupon":null},
	//  "type":"purchase","user_id":42}
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 map 編碼時，<b>鍵會依字母順序排序</b>，所以輸出是固定的；<code>nil</code> 會變成 <code>null</code>。
</div>

<!--
反過來，map[string]any 也可以直接編碼成 JSON，這在組合動態的資料時很方便，例如記錄事件、組合查詢條件。

注意 map 編碼時，鍵會依照字母順序排序。雖然 map 本身的走訪順序是隨機的，但 JSON 的輸出是固定的，這讓測試和比對結果變得容易。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# gob：Go 自有的編碼格式
## encoding/gob

<!--
最後介紹 Go 自有的編碼格式：gob。
-->

---

# gob 與 JSON 的比較

`encoding/gob` 是 Go 專用的**二進位**編碼格式，適合 **Go 程式之間**交換資料

| 比較 | JSON | gob |
| --- | --- | --- |
| **格式** | 純文字，人看得懂 | 二進位，人看不懂 |
| **跨語言** | ✅ 所有語言都支援 | ❌ 只有 Go |
| **型別資訊** | 數字都是 number | 保留 Go 的型別（`int`、`float32`、`[]byte`…） |
| **需要標籤嗎** | 需要 `json:"..."` 對應鍵名稱 | 不需要，直接用欄位名稱 |
| **適合情境** | API、設定檔、跨語言 | Go 程式之間的 RPC、快取、存檔 |

<!--
gob 是 Go 自己設計的二進位編碼格式。

跟 JSON 比較：gob 是二進位的，人看不懂，也只有 Go 能讀，但它會完整保留 Go 的型別資訊，而且不需要標籤。

什麼時候用 gob？當資料只在 Go 程式之間傳遞，例如兩個 Go 服務之間的內部通訊，或是把 Go 的資料暫存到檔案、快取裡。如果資料需要給其他語言、給人看，就用 JSON。

實務上，gob 的使用率遠低於 JSON；跨服務通訊現在更主流的是 Protocol Buffers，這裡知道有 gob 這個選項就好。
-->

---
zoom: 0.91
---

# gob — 範例

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Session struct {
	UserID int
	Roles  []string
	Scores map[string]float32
}

func main() {
	var buf bytes.Buffer // 實作了 io.Writer 與 io.Reader
	in := Session{
		UserID: 7,
		Roles:  []string{"admin"},
		Scores: map[string]float32{"go": 9.5},
	}
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		panic(err)
	}
	fmt.Println("編碼後大小：", buf.Len(), "bytes")
```

<!--
gob 的用法跟 JSON 的 Encoder、Decoder 幾乎一樣：NewEncoder 接收 io.Writer，NewDecoder 接收 io.Reader。

這裡用 bytes.Buffer 當作中介，它同時實作了 Writer 和 Reader，就像一個記憶體中的檔案：先寫進去，再讀出來。Session 沒有任何標籤，gob 直接用欄位名稱對應。
-->

---

# gob — 範例（續）

```go
	// 續上頁
	var out Session
	if err := gob.NewDecoder(&buf).Decode(&out); err != nil {
		panic(err)
	}
	// {UserID:7 Roles:[admin] Scores:map[go:9.5]}
	fmt.Printf("%+v\n", out)
}
```

<!--
接著從同一個 Buffer 解碼回 Session。解碼之後，float32 還是 float32，不會像 JSON 一樣變成 float64，這就是 gob 會保留 Go 型別資訊的好處。
-->

---
zoom: 0.95
---

# 補充：encoding/json/v2（Go 1.27 正式推出）

新版 JSON 套件修正了 v1 多年來的問題，**預設行為更安全、更一致**

| 行為 | `encoding/json`（v1） | `encoding/json/v2` |
| --- | --- | --- |
| 鍵名稱比對 | 不分大小寫 | **區分大小寫** |
| `nil` 切片 / map | `null` | **`[]` / `{}`** |
| 重複的鍵 | 允許，後者覆蓋 | **回傳錯誤** |
| 讀寫 `io.Reader` / `Writer` | `NewDecoder` / `NewEncoder` | `json.UnmarshalRead` / `json.MarshalWrite` |
| 選項 | 方法設定 | 函式參數：`json.RejectUnknownMembers(true)` |

```go
import "encoding/json/v2"

// {"id":1,"name":"Alice","tags":[]}
b, err := json.Marshal(User{ID: 1, Name: "Alice"})
err = json.UnmarshalRead(r.Body, &u, json.RejectUnknownMembers(true))
```

<!--
這是補充內容：Go 1.27 正式推出了新版的 JSON 套件 encoding/json/v2。

v1 的 encoding/json 從 Go 1.0 就存在，有一些多年來被詬病的行為，但因為相容性承諾，不能直接修改。所以 Go 團隊推出了 v2，import 路徑多了 /v2，這就是第 8 章說的「主版本號改變，路徑也改變」。

v2 的預設行為更嚴格、更安全：鍵名稱區分大小寫、nil 切片輸出空陣列（前端工程師最開心）、重複的鍵會回傳錯誤。

目前兩個版本可以並存，既有的程式碼繼續用 v1 完全沒問題，新專案可以考慮使用 v2。本章的觀念，struct 標籤、omitzero 這些，在 v2 都一樣適用。
-->

---
layout: default
---

# 綜合練習：待辦清單存檔
### 任務說明

1. 定義 `Todo struct`：`ID int`、`Title string`、`Done bool`、`DueAt time.Time`（`omitzero`）
2. 定義 `TodoList struct`：`Owner string`、`Items []Todo`（沒有資料時要輸出 `[]` 而不是 `null`）
3. 寫方法 `(l *TodoList) Save(w io.Writer) error`：用 `json.NewEncoder` 排版輸出
4. 寫函式 `Load(r io.Reader) (*TodoList, error)`：用 `json.NewDecoder`，**開啟嚴格模式**，錯誤時用 `%w` 包裝
5. 在 `main` 中：建立清單 → `Save` 到 `bytes.Buffer` → `Load` 讀回 → 印出未完成的項目
6. 再用 `strings.NewReader` 載入一段有錯字（`"titel"`）的 JSON，印出錯誤

<!--
這個綜合練習模擬「把資料存成 JSON 檔、再讀回來」的完整流程，第 12 章學完檔案操作後，把 bytes.Buffer 換成真正的檔案就可以了。

重點是 Save 和 Load 的參數用 io.Writer 和 io.Reader，而不是檔案，這樣測試的時候可以用 bytes.Buffer，正式使用時換成檔案，這就是第 7 章的「接受介面」。
-->

---
zoom: 0.94
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Todo struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Done  bool      `json:"done"`
	DueAt time.Time `json:"due_at,omitzero"`
}

type TodoList struct {
	Owner string `json:"owner"`
	Items []Todo `json:"items"`
}
```

<!--
Todo 的 DueAt 用 omitzero，沒設定期限就不輸出。TodoList 的 Items 是 Todo 的切片。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func (l *TodoList) Save(w io.Writer) error {
	if l.Items == nil {
		l.Items = []Todo{} // 確保輸出 [] 而不是 null
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(l)
}

func Load(r io.Reader) (*TodoList, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var l TodoList
	if err := dec.Decode(&l); err != nil {
		return nil, fmt.Errorf("載入待辦清單失敗：%w", err)
	}
	return &l, nil
}
```

<!--
Save 方法先檢查 Items 是不是 nil，是的話換成空切片，確保輸出 [] 而不是 null。然後用 NewEncoder 寫到傳入的 Writer。

Load 用 NewDecoder 讀取，開啟嚴格模式，解碼失敗時用 %w 包裝錯誤，加上「載入待辦清單失敗」的情境。
-->

---
zoom: 0.86
---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func main() {
	list := &TodoList{Owner: "小明", Items: []Todo{
		{ID: 1, Title: "寫作業", Done: true},
		{ID: 2, Title: "買咖啡",
			DueAt: time.Date(2026, 9, 25, 9, 0, 0, 0, time.UTC)},
	}}
	var buf bytes.Buffer
	if err := list.Save(&buf); err != nil {
		panic(err)
	}
	loaded, err := Load(&buf)
	if err != nil {
		panic(err)
	}
	for _, t := range loaded.Items {
		if !t.Done {
			due := t.DueAt.Format(time.DateTime)
			fmt.Println("未完成：", t.Title, due)
		}
	}
	typo := `{"owner":"x","items":[{"id":1,"titel":"typo"}]}`
	_, err = Load(strings.NewReader(typo))
	fmt.Println(err) // 載入待辦清單失敗：json: unknown field "titel"
}
```

<!--
main 先建立清單、存到 bytes.Buffer，再從同一個 Buffer 載入回來，印出未完成的項目：買咖啡，以及它的期限。

最後故意載入一段把 title 打成 titel 的 JSON，嚴格模式抓到了這個錯誤。如果沒有開嚴格模式，titel 會被默默忽略，Title 變成空字串，這種 bug 非常難找。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 11 步：JSON 與 gob 存檔

<!--
回到 GoShop。到目前為止，GoShop 有一個很大的問題：程式一結束，所有訂單和庫存的變化就全部消失了，下次執行又從頭開始。

另外，商品資料還寫死在 main.go 裡，要新增商品就得改程式碼。今天學的 JSON 和 gob，剛好可以解決這兩個問題。
-->

---

# GoShop 第 11 步：JSON 與 gob 存檔
### 任務說明

1. 替 `Product`、`Line`、`Order` 加上 **JSON 標籤**（`sku`、`paid_by,omitzero`…）
2. `store.DecodeProducts(r io.Reader)`：從 JSON 陣列讀取商品，**不認得的欄位視為錯誤**
3. `Memory.Save(w)`／`Memory.Load(r)`：用 **gob** 存取所有商品、訂單與編號
4. `main`：有 `data/goshop.gob` 就載入；沒有就從 `data/products.json` 匯入；結束前存檔
5. `Status` 實作 `MarshalText`，讓 JSON 輸出 `"paid"`，而不是看不懂的數字 `1`

```text
$ go run .          ← 第一次執行：訂單 #1、#2
$ go run .          ← 第二次執行：訂單編號接著是 #3、#4，庫存也接著扣
```

<!--
這一步要讓 GoShop 能夠存檔。

商品資料改放在 data/products.json，這是一個人可以直接打開來編輯的 JSON 檔。讀取的時候開啟嚴格模式，欄位名稱打錯就會報錯，而不是默默忽略。

存檔用的是 gob。gob 是 Go 專用的二進位格式，速度快、檔案小，而且 time.Time、map 這些型別都能直接存，很適合拿來做程式自己的存檔。

驗收的方式很簡單：連續執行兩次，第二次的訂單編號應該接著第一次，庫存也應該接著扣。
-->

---

# GoShop 第 11 步：解題提示
### JSON 標籤

```go
// goshop/internal/shop/shop.go
type Order struct {
	ID        int         `json:"id"`
	Lines     []Line      `json:"lines"`
	Subtotal  money.Money `json:"subtotal"`
	Discount  money.Money `json:"discount"`
	Total     money.Money `json:"total"`
	Status    Status      `json:"status"`
	PaidBy    string      `json:"paid_by,omitzero"`
	Coupon    string      `json:"coupon,omitzero"` // 折價券代碼
	CreatedAt time.Time   `json:"created_at"`      // 下單時間
	ShipBy    time.Time   `json:"ship_by"`         // 預計出貨日
}
```

- JSON 的慣例是 `snake_case`；Go 的欄位是 `PascalCase`，用標籤對應
- `omitzero`（Go 1.24+）：零值（空字串）時不輸出這個欄位

<!--
JSON 標籤寫在欄位型別的後面，用反引號包起來。它告訴 encoding/json：這個欄位在 JSON 裡叫什麼名字。

JSON 的欄位名稱慣例是小寫加底線，Go 的欄位要大寫開頭才能匯出，兩邊的慣例不一樣，就用標籤來對應。

PaidBy 和 Coupon 加上了 omitzero，沒有付款、沒有用折價券的時候，JSON 裡就不會出現這兩個欄位，輸出比較乾淨。time.Time 會自動編碼成 RFC 3339 格式的字串，例如 2026-09-24T15:12:55+08:00。
-->

---

# GoShop 第 11 步：解題提示（續）
### 嚴格解碼與 gob 快照

```go
// goshop/internal/store/snapshot.go
// DecodeProducts 從 JSON 陣列讀取商品清單；不認得的欄位視為錯誤。
func DecodeProducts(r io.Reader) ([]shop.Product, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var products []shop.Product
	if err := dec.Decode(&products); err != nil {
		return nil, fmt.Errorf("解析商品 JSON：%w", err)
	}
	return products, nil
}

// snapshot 是寫進 gob 檔的完整資料；欄位要匯出 gob 才看得到。
type snapshot struct {
	Products map[string]shop.Product
	Orders   map[int]shop.Order
	LastID   int
}
```

<!--
DecodeProducts 接收 io.Reader，而不是檔名。這樣它可以讀檔案、讀網路請求，測試時也可以用 strings.NewReader 直接傳一段字串進去，不需要真的建立檔案。

DisallowUnknownFields 是今天學的嚴格模式。如果有人把 price 打成 prize，就會得到錯誤，而不是讀出價格是 0 的商品。

snapshot 結構是要存進 gob 的資料。Memory 的欄位是小寫的，gob 看不到，所以我們另外定義一個欄位大寫的結構，存檔時把資料搬進去。
-->

---

# GoShop 第 11 步：解題提示（續 2）
### 存檔與讀檔

```go
// goshop/internal/store/snapshot.go
// Save 把所有商品與訂單用 gob 格式寫到 w。
func (m *Memory) Save(w io.Writer) error {
	return gob.NewEncoder(w).Encode(snapshot{m.products, m.orders, m.lastID})
}
```

```go
// goshop/main.go
func openStore(dir string) (*store.Memory, error) {
	st := store.NewMemory()
	f, err := os.Open(filepath.Join(dir, "goshop.gob"))
	if err == nil {
		defer f.Close()
		return st, st.Load(f)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	// ...
```

- 沒有快照時，改從 `products.json` 匯入商品（第一次執行）

<!--
Save 只有一行：建立 gob 編碼器，把 snapshot 編碼寫出去。Load 則是反過來解碼，再把資料放回 Memory。

main 的 openStore 先試著開啟快照檔。開啟成功就載入；如果錯誤是「檔案不存在」，代表這是第一次執行，就改從 products.json 匯入。用 errors.Is 搭配 fs.ErrNotExist 判斷，是檢查檔案不存在的標準寫法。

其他的錯誤，例如沒有讀取權限，就直接回報，不能當成第一次執行。
-->

---

# GoShop 第 11 步：解題提示（續 3）
### 自訂 JSON 的輸出：MarshalText

```go
// goshop/internal/shop/shop.go
// MarshalText 讓 JSON 輸出 "pending"／"paid"，而不是看不懂的數字。
func (s Status) MarshalText() ([]byte, error) {
	if s == Paid {
		return []byte("paid"), nil
	}
	return []byte("pending"), nil
}
```

```text
{
  "id": 2,
  "total": 2260,
  "status": "paid",
  "paid_by": "貨到付款",
  "created_at": "2026-09-24T15:12:55.427501531+08:00",
  ...
}
```

<!--
Status 是一個整數，直接編碼會變成 "status": 1，別人看了根本不知道 1 是什麼意思。

只要替 Status 實作 encoding.TextMarshaler 介面，也就是 MarshalText 方法，encoding/json 就會改用它的結果，輸出 "status": "paid"。反過來，實作 UnmarshalText 就能把 "paid" 解碼回 Status。這又是一個介面的應用：encoding/json 只認得介面，不用知道 Status 是什麼。

下面是用 MarshalIndent 印出的訂單 JSON，status 變成了看得懂的文字，沒有用折價券的 coupon 欄位也因為 omitzero 被省略了。
-->

---

# 章節總結

- **JSON 與 Go**：物件 ↔ struct / map、陣列 ↔ 切片；只有**匯出的欄位**會被處理
- **解碼**：`json.Unmarshal(data, &v)` 傳指標；多出的鍵忽略、缺少的保持零值；指標欄位區分「沒給」與「零值」
- **編碼**：`json.Marshal(v)`、`MarshalIndent`；map 的鍵會排序；**nil 切片輸出 `null`**
- **struct 標籤**：`json:"name"`、`omitempty`、**`omitzero`（1.24+）**、`json:"-"`
- **Decoder / Encoder**：搭配 `io.Reader` / `io.Writer`；`DisallowUnknownFields` 嚴格模式
- **未知結構**：解碼到 `map[string]any`，數字一律是 `float64`，用型別斷言取值
- **gob** 是 Go 專用的二進位格式；**`encoding/json/v2`**（1.27）預設更安全
- **GoShop**：替資料加上 JSON 標籤、從 `products.json` 匯入商品，用 gob 快照讓資料在重新執行後還在

下一章我們會介紹「系統與檔案」：命令列旗標、系統訊號，以及檔案的讀寫。

<!--
我們來整理今天學到的東西。

JSON 是現代程式交換資料的共通語言。Go 用 struct 標籤對應 JSON 的鍵，Unmarshal 解碼、Marshal 編碼，記得只有匯出的欄位會被處理。omitempty 和 omitzero 用來省略空值，Decoder 和 Encoder 可以直接讀寫 io.Reader 和 io.Writer，嚴格模式可以抓到打錯字的鍵。

GoShop 這一步終於「記得住」東西了：第一次執行從 products.json 匯入商品，結束前把所有資料存成 gob 快照，下次執行接著用。JSON 給人看、給別的系統交換資料；gob 給 Go 程式自己存檔，兩者各有用途。

今天的綜合練習用 bytes.Buffer 模擬檔案。下一章就要學真正的檔案操作：建立、讀取、寫入、刪除檔案，還有處理 CSV 格式。另外也會學命令列旗標和系統訊號，讓我們寫出真正實用的命令列工具。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
JSON 是後面章節的基礎，第 14 章呼叫 API、第 15 章寫 API 都會大量使用。

課後建議：找一個公開的 JSON API，例如政府資料開放平台，把回應的 JSON 貼到網路上的「JSON to Go」工具，看看它自動產生的 struct，再跟自己寫的比較。

有問題的同學現在可以提問！
-->

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
title: 時間處理
routeAlias: ch10
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">時間處理</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「記住一個日子：2006 年 1 月 2 日下午 3 點 4 分 5 秒」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十章！

時間是每個應用程式都會用到的東西：訂單的建立時間、會員的生日、優惠活動的截止日、API 的回應時間。但時間處理也是最容易出錯的地方：時區搞錯、月底加一個月變成下下個月、格式化字串寫錯。

Go 的 time 套件功能很完整，但它的時間格式化方式非常特別，跟其他語言都不一樣。封面上那句話就是今天的關鍵：記住 2006 年 1 月 2 日下午 3 點 4 分 5 秒，Go 的時間格式化就學會一半了。
-->

---
layout: default
---

# Outline

- **建立時間資料** — 取得系統時間、建立指定時間、取出年月日時分秒
- **時間值的格式化** — 參考時間 `2006-01-02 15:04:05`、`Format`、`Parse`
- **時間值的管理** — 增減時間、設定時區
- **比較與時間長度** — `Before` / `After` / `Equal`、`time.Duration`、測量執行時間
- **章節總結**

<!--
今天的內容分成四個部分。

先學怎麼取得和建立時間，以及怎麼取出年、月、日這些資訊。接著是 Go 最特別的格式化方式。第三部分是時間的加減和時區，最後是時間的比較，以及用 Duration 表示「一段時間長度」。
-->

---

# 回顧：程式除錯

- fmt 格式化：`%v`、`%+v`、`%#v`、`%T`；寬度與精度 `%8.2f`
- 日誌：`log` 自動加時間；`log/slog` 結構化日誌，`With` 附加固定欄位
- 單元測試：表格驅動測試 + `t.Run`；`go test -v -cover`
- 效能測試：`for b.Loop() { }`，測量每次執行花多少時間 ← 今天學怎麼自己測量時間

<!--
回顧一下上一章。

我們學了 fmt 的格式化、log 和 slog 的日誌，以及單元測試。日誌的每一行都有時間戳記，效能測試會告訴我們每次執行花了幾奈秒，這些背後都是 time 套件。今天我們就來學怎麼自己處理時間。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 建立時間資料
## Creating Time Values

<!--
先從怎麼取得和建立時間開始。
-->

---

# 取得系統時間：time.Now()

`time.Time` 代表**某一個時間點**，精確到奈秒，並帶有時區資訊

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	// 2026-09-24 14:05:09.123456789 +0800 CST m=+0.000012345
	fmt.Println(now)
	fmt.Println(now.Unix())      // 1790229909：Unix 時間戳記（秒）
	fmt.Println(now.UnixMilli()) // 毫秒時間戳記，常用於 JavaScript / API

	var zero time.Time         // 零值：西元 1 年 1 月 1 日
	fmt.Println(zero.IsZero()) // true
}
```

<!--
time.Now() 會回傳目前的系統時間，型別是 time.Time。

直接印出來，可以看到日期、時間、奈秒、時區。最後面的 m=+0.000012345 是「單調時鐘」的讀數，是 Go 用來精確測量時間間隔的，等一下會說明。

Unix 時間戳記是「從 1970 年 1 月 1 日 UTC 開始經過的秒數」，資料庫、API 常常用它來存時間。UnixMilli 是毫秒版本，JavaScript 的 Date.now() 就是毫秒。

time.Time 的零值是西元 1 年 1 月 1 日，可以用 IsZero 判斷一個時間有沒有被設定過。
-->

---

# 建立指定的時間：time.Date()

語法：`time.Date(年, 月, 日, 時, 分, 秒, 奈秒, 時區)`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(2026, time.September, 24, 14, 5, 9, 0, time.UTC)
	fmt.Println(t) // 2026-09-24 14:05:09 +0000 UTC

	loc := time.Local // 電腦設定的本地時區
	launch := time.Date(2026, 12, 25, 20, 0, 0, 0, loc)
	fmt.Println(launch.Year(), launch.Month()) // 2026 December
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>月份是 <code>time.Month</code> 型別：</b> 可以寫 <code>time.September</code>，也可以直接寫數字 <code>9</code>（無型別常數會自動轉換）。月份從 <b>1</b> 開始，不像 JavaScript 從 0 開始。
</div>

<!--
time.Date 用來建立一個指定的時間，參數依序是年、月、日、時、分、秒、奈秒，最後是時區。

月份的型別是 time.Month，它是第一章學過的列舉，可以寫 time.September，也可以直接寫數字 9。好消息是 Go 的月份從 1 開始，不像 JavaScript 的一月是 0，那是 JavaScript 最有名的坑之一。

時區可以用 time.UTC，或 time.Local 代表電腦設定的本地時區。
-->

---
zoom: 0.91
---

# 取得時間資料中的特定項目

| 方法 | `2026-09-24 14:05:09` 的結果 | 型別 |
| --- | --- | --- |
| `t.Year()` / `t.Month()` / `t.Day()` | `2026` / `September` / `24` | `int` / `time.Month` / `int` |
| `t.Hour()` / `t.Minute()` / `t.Second()` | `14` / `5` / `9` | `int` |
| `t.Weekday()` | `Thursday` | `time.Weekday` |
| `t.YearDay()` | `267`（一年中的第幾天） | `int` |
| `t.Date()` | `2026 September 24`（一次取三個值） | `(int, Month, int)` |
| `t.ISOWeek()` | `2026 39`（ISO 週數） | `(int, int)` |

```go
y, m, d := t.Date()
// 2026 年 9 月 24 日（Thursday）
fmt.Printf("%d 年 %d 月 %d 日（%s）\n", y, m, d, t.Weekday())
```

<!--
time.Time 有很多方法可以取出其中的某一部分。

年、月、日、時、分、秒各有一個方法。Weekday 回傳星期幾，YearDay 回傳這是一年中的第幾天。Date 方法可以一次拿到年月日三個值，這是第 5 章學的多重回傳值。

注意 Month 和 Weekday 都是自訂型別，而且有 String 方法，所以用 %s 或 Println 印出來是英文名稱；用 %d 印出來則是數字。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 時間值的格式化
## Formatting & Parsing

<!--
接下來是 Go 最特別的部分：時間的格式化。
-->

---

# Go 的參考時間（reference time）

其他語言用 `YYYY-MM-DD` 這種代號；**Go 用一個特定的時間當範本**：

<div class="text-center text-3xl font-bold my-6" style="color:#1a5c5c">
Mon Jan 2 15:04:05 MST 2006
</div>

| 元素 | 代號 | 記憶法 | 元素 | 代號 | 記憶法 |
| --- | --- | --- | --- | --- | --- |
| 月 | `01` / `1` / `Jan` | **1** | 秒 | `05` / `5` | **5** |
| 日 | `02` / `2` | **2** | 年 | `2006` / `06` | **6** |
| 時 | `15`（24 時制）/ `03`（12 時制） | **3**（下午 3 點） | 時區 | `-0700` / `MST` | **7** |
| 分 | `04` / `4` | **4** | 星期 | `Mon` / `Monday` | — |

<!--
這是 Go 時間格式化最特別的地方。

大部分語言用 YYYY 代表年、MM 代表月、DD 代表日這種代號。Go 不一樣，它用一個「參考時間」當範本：2006 年 1 月 2 日，下午 3 點 4 分 5 秒，時區是 -0700。

為什麼是這個時間？因為用美式的順序寫出來，剛好是 1、2、3、4、5、6、7：1 月、2 日、下午 3 點、4 分、5 秒、06 年、-07 時區。

我們想要什麼格式，就把這個參考時間「寫成那個格式」。例如想要「年-月-日」，就寫 2006-01-02。一開始會覺得很怪，但習慣之後會發現它很直覺：格式字串本身就長得跟輸出一模一樣。
-->

---

# 將時間轉成指定格式的字串：Format

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(2026, 9, 24, 14, 5, 9, 0, time.UTC)
	fmt.Println(t.Format("2006-01-02 15:04:05")) // 2026-09-24 14:05:09
	fmt.Println(t.Format("2006/1/2 3:04PM"))     // 2026/9/24 2:05PM
	fmt.Println(t.Format("Monday, 02-Jan-06"))   // Thursday, 24-Sep-26
	fmt.Println(t.Format("20060102"))            // 20260924：常用於檔名
	// 2026 年 09 月 24 日
	fmt.Println(t.Format("2006 年 01 月 02 日"))
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>最常見的錯誤：</b> 把格式寫成 <code>"2006-01-02 12:04:05"</code>（用 12 代表小時），結果每次都輸出「12」。24 小時制要用 <b><code>15</code></b>，12 小時制用 <b><code>03</code></b> 或 <b><code>3</code></b>。
</div>

<!--
Format 方法接收一個格式字串，回傳格式化後的字串。

大家看看每一行：格式字串長什麼樣子，輸出就長什麼樣子。想要斜線分隔就寫斜線，想要中文就直接寫中文。01 代表補零的月份，1 代表不補零的月份。

最常見的錯誤是把小時寫成 12：很多人直覺覺得「小時就是 12」，但 Go 看到 12 不會把它當成小時的代號，就原封不動輸出 12。記住：24 小時制是 15，也就是下午 3 點。
-->

---

# 內建的格式常數

| 常數 | 格式 | 輸出範例 |
| --- | --- | --- |
| `time.DateTime` | `"2006-01-02 15:04:05"` | `2026-09-24 14:05:09` |
| `time.DateOnly` | `"2006-01-02"` | `2026-09-24` |
| `time.TimeOnly` | `"15:04:05"` | `14:05:09` |
| `time.RFC3339` | `"2006-01-02T15:04:05Z07:00"` | `2026-09-24T14:05:09Z` |
| `time.RFC1123` | `"Mon, 02 Jan 2006 15:04:05 MST"` | `Thu, 24 Sep 2026 14:05:09 UTC` |
| `time.Kitchen` | `"3:04PM"` | `2:05PM` |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>建議：</b> API 與 JSON 交換時間一律用 <b>RFC3339</b>（ISO 8601）；<code>DateTime</code>、<code>DateOnly</code>、<code>TimeOnly</code> 是 Go 1.20 加入的常數，不用再自己寫格式字串。
</div>

<!--
time 套件內建了很多常用的格式常數，不用每次都自己寫。

DateTime、DateOnly、TimeOnly 是 Go 1.20 加入的，涵蓋了最常用的三種格式。

RFC3339 是國際標準的時間格式，也就是 ISO 8601，最後的 Z 代表 UTC 時區，如果是台灣時間會變成 +08:00。API 和 JSON 交換時間的時候，一律建議用 RFC3339，因為它包含時區資訊，不會有誤解。第 11 章的 JSON 編碼，time.Time 預設就會轉成 RFC3339 格式。
-->

---

# 將特定格式的時間字串轉成時間值：Parse

語法：`time.Parse(格式, 字串) (time.Time, error)`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t, err := time.Parse("2006-01-02 15:04", "2026-12-25 20:30")
	fmt.Println(t, err) // 2026-12-25 20:30:00 +0000 UTC <nil>

	_, err = time.Parse(time.DateOnly, "2026-13-01")
	fmt.Println(err) // parsing time "2026-13-01": month out of range

	loc, _ := time.LoadLocation("Asia/Taipei")
	s := "2026-09-24 09:00:00"
	tw, _ := time.ParseInLocation(time.DateTime, s, loc)
	// 2026-09-24 09:00:00 +0800 CST = 2026-09-24 01:00:00 +0000 UTC
	fmt.Println(tw, "=", tw.UTC())
}
```

<!--
Parse 是 Format 的反向操作：把字串解析成 time.Time。格式字串一樣用參考時間的寫法。

Parse 會回傳錯誤，因為使用者輸入的字串可能格式不對、或日期不存在，例如 13 月。第 6 章學的錯誤處理在這裡派上用場。

使用 Parse 的注意事項：如果字串裡沒有時區資訊，Parse 會把它當成 UTC。在台灣，使用者輸入「早上 9 點」通常指的是台灣時間，這時候要用 ParseInLocation，指定時區來解析。台灣早上 9 點，等於 UTC 凌晨 1 點。
-->

---
layout: default
---

# 練習 1：活動倒數
### 任務說明

1. 用 `time.Parse` 解析活動開始時間字串 `"2026/12/31 23:59"`（台灣時間，用 `ParseInLocation`）
2. 印出活動時間的三種格式：
   - `2026 年 12 月 31 日 23:59`
   - `Thursday`（星期幾）
   - RFC3339 格式
3. 解析一個錯誤的字串 `"2026/02/30 10:00"`，印出錯誤訊息

<!--
這個練習要大家練習參考時間的寫法。

注意題目的日期用斜線分隔，格式字串也要用斜線。2 月 30 日不存在，Parse 會怎麼處理？
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
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	const layout = "2006/01/02 15:04"
	event, err := time.ParseInLocation(layout, "2026/12/31 23:59", loc)
	if err != nil {
		fmt.Println("解析失敗：", err)
		return
	}
	fmt.Println(event.Format("2006 年 01 月 02 日 15:04"))
	fmt.Println(event.Weekday())
	fmt.Println(event.Format(time.RFC3339)) // 2026-12-31T23:59:00+08:00

	_, err = time.ParseInLocation(layout, "2026/02/30 10:00", loc)
	// parsing time "2026/02/30 10:00": day out of range
	fmt.Println(err)
}
```

<!--
格式字串用 const 宣告，因為解析兩次都要用同一個格式。

RFC3339 格式的結尾是 +08:00，代表這是 UTC 加 8 小時的時間，也就是台灣時間。2026 年 12 月 31 日是星期四。

2 月 30 日不存在，Parse 會回傳 day out of range 的錯誤。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 時間值的管理
## Adding Time & Time Zones

<!--
接下來學怎麼對時間做加減，以及怎麼處理時區。
-->

---

# 建立和增減時間值

| 方法 | 用途 | 範例 |
| --- | --- | --- |
| `t.Add(d)` | 加上一段**時間長度**（時、分、秒） | `t.Add(90 * time.Minute)` |
| `t.AddDate(年, 月, 日)` | 加上**年月日** | `t.AddDate(0, 1, 0)` 加一個月 |
| `t.Truncate(d)` | 捨去到某個單位 | `t.Truncate(time.Hour)` 整點 |

```go
t := time.Date(2026, 9, 24, 14, 5, 9, 0, time.UTC)
fmt.Println(t.Add(90 * time.Minute))  // 2026-09-24 15:35:09 +0000 UTC
// 2026-09-23 14:05:09：負數代表往前
fmt.Println(t.Add(-24 * time.Hour))
// 2026-10-31 14:05:09：一個月又七天後
fmt.Println(t.AddDate(0, 1, 7))
```

<!--
時間的加減有兩個方法。

Add 加上一段「時間長度」，也就是 time.Duration，適合加減時、分、秒。負數代表往前。time.Minute、time.Hour 這些是 Duration 常數，乘上數字就是幾分鐘、幾小時。

AddDate 加上年、月、日，適合「下個月」、「明年」這種以日曆為單位的計算。

為什麼要分兩個方法？因為「一天」不一定是 24 小時：在有日光節約時間的地方，一年會有一天只有 23 小時、一天有 25 小時。AddDate 是依照日曆計算，Add 是依照實際經過的時間計算。
-->

---

# 使用 AddDate 的注意事項：月底問題

`AddDate` 會把「不存在的日期」**自動往後進位（正規化）**

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	jan31 := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	// 2026-03-03 ⚠️ 不是 2 月底！
	fmt.Println(jan31.AddDate(0, 1, 0).Format(time.DateOnly))

	// 想取得「下個月的最後一天」：下下個月的第 0 天
	lastDay := time.Date(2026, 2+1, 0, 0, 0, 0, 0, time.UTC)
	fmt.Println(lastDay.Format(time.DateOnly)) // 2026-02-28
}
```

<!--
使用 AddDate 有一個很重要的注意事項：月底問題。

1 月 31 日加一個月，我們直覺會覺得是 2 月 28 日。但 Go 的做法是：先把月份加 1，變成「2 月 31 日」，2 月沒有 31 日，就往後進位，2 月 28 日再過 3 天，變成 3 月 3 日。

這在處理「每月扣款日」、「會員到期日」的時候很容易出錯。

如果想要「某個月的最後一天」，有一個技巧：用 time.Date 建立「下個月的第 0 天」，第 0 天會被正規化成上個月的最後一天。所以 3 月 0 日就是 2 月 28 日。
-->

---
zoom: 0.86
---

# 設定時區來取得新時間值

`time.Time` 內含時區；用 `t.In(loc)` 取得**同一個時間點**在另一個時區的表示

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	meeting := time.Date(2026, 9, 24, 14, 0, 0, 0, time.UTC)

	zones := []string{
		"Asia/Taipei", "America/New_York", "Europe/London",
	}
	for _, name := range zones {
		loc, err := time.LoadLocation(name) // IANA 時區名稱
		if err != nil {
			fmt.Println(err)
			continue
		}
		local := meeting.In(loc).Format("01/02 15:04 MST")
		fmt.Printf("%-18s %s\n", name, local)
	}
	utc8 := time.FixedZone("UTC+8", 8*60*60)
	fmt.Println(meeting.In(utc8).Format(time.RFC3339))
}
```

<!--
同一個時間點，在不同時區會有不同的「時鐘時間」。台灣下午 10 點開會，美國東岸是早上 10 點。

LoadLocation 用 IANA 時區資料庫的名稱載入時區，格式是「洲/城市」，例如 Asia/Taipei。In 方法會回傳同一個時間點在另一個時區的表示，時間點本身沒有變，只是換一個角度看。

FixedZone 可以建立一個固定偏移量的時區，參數是名稱和相對於 UTC 的秒數。它不會處理日光節約時間，所以有日光節約的地區還是要用 LoadLocation。

執行後會看到：台北 22:00 CST、紐約 10:00 EDT、倫敦 15:00 BST。
-->

---

# 使用時區的注意事項

**注意事項之一：** 伺服器內部與資料庫**一律用 UTC 儲存**，只在顯示給使用者時才轉換時區

**注意事項之二：** `LoadLocation` 需要系統的時區資料庫；在精簡的 Docker 映像檔中可能找不到

```go
import _ "time/tzdata" // 把時區資料庫編譯進執行檔（約增加 450 KB）
```

| 情境 | 建議 |
| --- | --- |
| 儲存、比較、計算 | `t.UTC()` |
| 顯示給使用者 | `t.In(userLoc).Format(...)` |
| 解析使用者輸入 | `time.ParseInLocation(layout, s, userLoc)` |
| API 傳遞 | RFC3339 格式（包含時區偏移） |

<!--
使用時區有兩個注意事項。

第一，這是業界的共識：伺服器內部和資料庫一律用 UTC。因為使用者可能在世界各地，伺服器也可能部署在不同的時區，統一用 UTC 就不會混亂。只有在「顯示給使用者看」的最後一步，才轉換成使用者的時區。

第二，LoadLocation 需要作業系統提供的時區資料庫。在一般的電腦上沒問題，但如果把 Go 程式放進非常精簡的 Docker 映像檔，可能會找不到時區資料，出現錯誤。解法是空白匯入 time/tzdata，把時區資料直接編譯進執行檔。這就是上一章學的空白匯入的實際用途。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 時間值的比較與時間長度處理
## Comparing Times & Durations

<!--
最後來看時間的比較，以及時間長度 Duration。
-->

---

# 比較時間

| 方法 | 意義 |
| --- | --- |
| `a.Before(b)` | `a` 是否在 `b` 之前 |
| `a.After(b)` | `a` 是否在 `b` 之後 |
| `a.Equal(b)` | 是否為**同一個時間點**（不同時區也能正確比較） |
| `a.Compare(b)` | 回傳 `-1`、`0`、`+1`（Go 1.20+，適合排序） |

```go
utc := time.Date(2026, 9, 24, 14, 0, 0, 0, time.UTC)
// 同一個時間點，台灣時間 22:00
tw := utc.In(time.FixedZone("UTC+8", 8*3600))

fmt.Println(utc.Equal(tw)) // true：同一個時間點
fmt.Println(utc == tw)     // false ⚠️：時區不同，結構內容不同
```

<!--
比較兩個時間，用 Before、After、Equal 這三個方法。Go 1.20 加入的 Compare 回傳 -1、0、1，很適合搭配 slices.SortFunc 排序。

使用時間比較的注意事項：不要用 == 比較時間！

utc 和 tw 是同一個時間點，只是時區不同。Equal 會比較「是不是同一個時間點」，回傳 true；但 == 會比較整個結構的內容，包含時區，所以回傳 false。這是很常見的 bug。記住：比較時間一律用 Equal。
-->

---

# 時間長度：time.Duration

`time.Duration` 是 `int64` 的自訂型別，單位是**奈秒**

| 常數 | 值 | 範例 |
| --- | --- | --- |
| `time.Nanosecond` | 1 | — |
| `time.Millisecond` | 1,000,000 | `500 * time.Millisecond` |
| `time.Second` | 1,000,000,000 | `3 * time.Second` |
| `time.Minute` / `time.Hour` | 60 秒 / 60 分 | `90 * time.Minute` |

```go
d, _ := time.ParseDuration("1h30m15s") // 從字串解析
fmt.Println(d, d.Minutes())            // 1h30m15s 90.25
fmt.Println(d.Round(time.Hour))        // 2h0m0s
timeout := 5 * time.Second             // 常用於設定逾時（Ch 14、16）
```

<!--
time.Duration 代表一段時間長度，例如 3 秒、90 分鐘。它其實是一個 int64 的自訂型別，單位是奈秒，所以 1 秒就是 10 億。

time 套件定義了一組常數，用乘法組合出想要的長度：3 * time.Second 就是 3 秒。這個寫法的可讀性非常好，一看就知道是幾秒。

ParseDuration 可以從字串解析，例如設定檔裡寫 "30s"、"1h30m"。Duration 有 String 方法，印出來就是這種人類可讀的格式。

Duration 在後面的章節會一直出現：第 14 章 HTTP 客戶端的逾時設定、第 16 章並行性運算的等待時間，都是用 Duration。
-->

---

# 用時間長度來改變時間、計算間隔

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Date(2026, 9, 24, 9, 0, 0, 0, time.UTC)
	// 用 Duration 改變時間
	end := start.Add(8*time.Hour + 30*time.Minute)
	fmt.Println(end.Format(time.TimeOnly)) // 17:30:00

	worked := end.Sub(start)            // 兩個時間相減得到 Duration
	fmt.Println(worked, worked.Hours()) // 8h30m0s 8.5

	xmas := time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC)
	days := int(xmas.Sub(start).Hours() / 24)
	fmt.Println("距離聖誕節還有", days, "天") // 距離聖誕節還有 91 天
}
```

<!--
Duration 和 Time 可以互相運算。

Time 加上 Duration 得到新的 Time，這個我們已經看過了。兩個 Time 用 Sub 相減，得到的是 Duration，也就是兩個時間點之間經過了多久。

Duration 有 Hours、Minutes、Seconds 等方法，把它轉成浮點數的時數、分鐘數。要算天數，就用小時數除以 24。

注意 8*time.Hour + 30*time.Minute，兩個 Duration 可以直接相加，因為它們是同一個型別。
-->

---
zoom: 0.93
---

# 測量時間長度：time.Since

```go
package main

import (
	"fmt"
	"time"
)

func timeTrack(start time.Time, name string) {
	// time.Since(t) = time.Now().Sub(t)
	fmt.Printf("%s 花了 %v\n", name, time.Since(start))
}

func slowTask() {
	// 參數在 defer 時就求值 → 記錄開始時間
	defer timeTrack(time.Now(), "slowTask")
	time.Sleep(120 * time.Millisecond) // 模擬耗時的工作
}

func main() {
	slowTask() // slowTask 花了 120.123456ms
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>單調時鐘：</b> <code>time.Now()</code> 內含單調時鐘讀數，<code>Since</code>、<code>Sub</code> 會用它計算，即使中途系統時間被校正（NTP）也不會算錯。
</div>

<!--
time.Since 是 time.Now().Sub(t) 的簡寫，用來計算「從某個時間點到現在經過了多久」，最常用來測量程式的執行時間。

這裡用了第 5 章學的 defer 技巧：defer timeTrack(time.Now(), "slowTask")。還記得 defer 的參數會「立刻求值」嗎？所以 time.Now() 在 defer 這一行就被執行了，記錄的是開始時間；等到函式結束時才呼叫 timeTrack，計算經過的時間。一行就完成了執行時間的測量。

time.Sleep 會讓程式暫停一段時間，這裡用來模擬一個很慢的工作。

最後的補充：time.Now() 回傳的時間內含一個「單調時鐘」，它只會往前走，不受系統時間校正的影響。所以用 Since 測量時間間隔，永遠是準確的。
-->

---
layout: default
---

# 綜合練習：會員到期提醒
### 任務說明

有一組會員資料（到期日以 `time.DateOnly` 格式的字串儲存）：

```go
members := map[string]string{
	"Alice": "2026-10-01", "Bob": "2026-09-20", "Carol": "2027-01-15",
}
```

以 `today := time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC)` 為今天：

1. 解析每個會員的到期日，解析失敗時印出錯誤
2. 計算距離到期還有幾天（已過期為負數）
3. **依到期日排序**（提示：`slices.SortFunc` + `time.Time.Compare`）後印出：
   - 已過期：`Bob 已過期 4 天`
   - 30 天內到期：`Alice 將於 7 天後到期（Thursday）`
   - 其他：`Carol 到期日 2027/01/15`

<!--
這個綜合練習把今天學的東西串起來：Parse 解析字串、Sub 計算間隔、Compare 排序、Format 格式化。

排序的部分，slices.SortFunc 接收一個比較函式，回傳負數、0、正數，剛好就是 Compare 的回傳值。
-->

---
zoom: 0.77
---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"fmt"
	"slices"
	"time"
)

type member struct {
	name   string
	expire time.Time
}

func parseAll(raw map[string]string) []member {
	var list []member
	for name, s := range raw {
		t, err := time.Parse(time.DateOnly, s)
		if err != nil {
			fmt.Println(name, "日期格式錯誤：", err)
			continue
		}
		list = append(list, member{name, t})
	}
	slices.SortFunc(list, func(a, b member) int {
		return a.expire.Compare(b.expire) // 依到期日由早到晚
	})
	return list
}
```

<!--
parseAll 先把 map 轉成 member 結構的切片，因為 map 沒有順序，要排序一定要先轉成切片。解析失敗的會員印出錯誤後跳過。

slices.SortFunc 的比較函式直接回傳 a.expire.Compare(b.expire)，就會依到期日由早到晚排序。
-->

---
zoom: 0.94
---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func main() {
	today := time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC)
	list := parseAll(map[string]string{
		"Alice": "2026-10-01",
		"Bob":   "2026-09-20",
		"Carol": "2027-01-15",
	})
	for _, m := range list {
		days := int(m.expire.Sub(today).Hours() / 24)
		switch {
		case days < 0:
			fmt.Printf("%s 已過期 %d 天\n", m.name, -days)
		case days <= 30:
			fmt.Printf("%s 將於 %d 天後到期（%s）\n",
				m.name, days, m.expire.Weekday())
		default:
			fmt.Printf("%s 到期日 %s\n",
				m.name, m.expire.Format("2006/01/02"))
		}
	}
}
```

<!--
main 呼叫 parseAll 取得排序好的切片，再用 Sub 計算到期日和今天相差的時間，除以 24 小時得到天數。

用無條件 switch 分成三種情況：負數代表已過期，30 天以內提醒即將到期並顯示星期幾，其他的顯示到期日。

因為 Parse 沒有指定時區，會解析成 UTC，跟 today 同一個時區，所以相減的結果剛好是整數天。執行後依序印出 Bob 已過期 4 天、Alice 將於 7 天後到期、Carol 到期日 2027/01/15。
-->

---

# 章節總結

- **建立時間**：`time.Now()`、`time.Date(...)`；用 `Year()`、`Month()`、`Weekday()` 等方法取出資訊
- **參考時間**：`Mon Jan 2 15:04:05 MST 2006`（1 2 3 4 5 6 7）；24 小時制是 **`15`**
- **格式化與解析**：`t.Format(layout)`、`time.Parse(layout, s)`；使用者輸入用 `ParseInLocation`
- **格式常數**：`time.DateTime`、`DateOnly`、`TimeOnly`（1.20+）；API 傳遞用 `RFC3339`
- **增減時間**：`Add(Duration)`、`AddDate(年, 月, 日)`；小心**月底進位**
- **時區**：內部一律 UTC，顯示時 `t.In(loc)`；精簡映像檔用 `import _ "time/tzdata"`
- **比較與長度**：用 `Equal` 不要用 `==`；`Sub` 得到 `Duration`；`time.Since` 測量執行時間

下一章我們會介紹「JSON 的編碼與解碼」：Go 結構和 JSON 資料的互相轉換。

<!--
我們來整理今天學到的東西。

Go 的時間格式化用參考時間 2006 年 1 月 2 日下午 3 點 4 分 5 秒，記住 1234567 的口訣。解析使用者輸入時要注意時區，內部一律用 UTC。加減時間用 Add 和 AddDate，小心月底進位。比較時間用 Equal，不要用 ==。Duration 代表時間長度，用 time.Since 測量執行時間。

下一章要學 JSON。現代的程式幾乎都透過 JSON 交換資料：前端和後端、服務和服務之間、設定檔。Go 用 encoding/json 套件，搭配 struct 標籤，就能輕鬆地在 Go 結構和 JSON 之間互相轉換。今天學的 time.Time，在 JSON 裡就會變成 RFC3339 格式的字串。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
時間處理看似簡單，實務上卻處處是坑：時區、月底、日光節約時間。

課後建議：寫一個小程式，列出你的朋友或同事所在城市的目前時間，練習 LoadLocation 和 In。如果身邊沒有國外的朋友，就列出東京、紐約、倫敦、雪梨。

有問題的同學現在可以提問！
-->

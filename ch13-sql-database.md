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
title: SQL 與資料庫
routeAlias: ch13
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">SQL 與資料庫</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「database/sql 一套介面，連接各種資料庫」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十三章！

上一章我們學會了把資料存進檔案。但檔案有很多限制：多個程式同時寫入會互相覆蓋、要找某一筆資料得從頭讀到尾、資料之間的關係很難維護。所以真實世界的應用程式，資料幾乎都存在資料庫裡。

今天我們要學的是：怎麼用 Go 連接 MySQL 資料庫，完成資料的新增、查詢、更新和刪除，也就是常說的 CRUD。Go 的標準函式庫提供了 database/sql 套件，它定義了一套通用的介面，搭配不同的驅動程式，就能連接 MySQL、PostgreSQL、SQLite 等各種資料庫。
-->

---
layout: default
---

# Outline

- **前言** — `database/sql` 與驅動程式的關係
- **安裝 MySQL 資料庫** — 安裝、建立使用者與資料庫、下載驅動程式
- **以 Go 語言連接資料庫** — `sql.Open`、`Ping`、連線池
- **建立、清空和移除資料表** — `ExecContext`
- **插入資料** — 佔位符 `?`、交易（transaction）與預備敘述
- **查詢資料** — `QueryContext`、`QueryRowContext`、`NULL` 的處理
- **更新既有資料** — `UPDATE`、`DELETE`、`RowsAffected`
- **練習：FizzBuzz 統計表**
- **章節總結**

<!--
今天的內容依照實際開發的順序安排：先把 MySQL 安裝好，建立使用者和資料庫；接著用 Go 連線，建立資料表；然後依序學新增、查詢、更新、刪除。

最後的綜合練習，會把第 2 章的 FizzBuzz 存進資料庫，再用 SQL 做統計。
-->

---

# 回顧：系統與檔案

- 命令列旗標：`flag.String` / `flag.Int` + `flag.Parse()`
- 檔案：`os.Create` + `defer f.Close()`；`os.ReadFile` / `os.WriteFile`；`bufio.Scanner` 逐行讀取
- CSV：`csv.NewReader(r).Read()` 讀到 `io.EOF`
- `signal.NotifyContext` 回傳的 **`context`** ← 今天每個資料庫操作都會用到它

<!--
回顧一下上一章。

我們學了命令列旗標、系統訊號和檔案操作。上一章在系統訊號的地方，第一次看到了 context，今天每一個資料庫操作都會帶著 context，用來控制逾時和取消。第 16 章會完整介紹 context，今天先記住一個用法就好：context.WithTimeout 可以設定「最多等多久」。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## database/sql & Drivers

<!--
先來了解 Go 連接資料庫的架構。
-->

---

# database/sql 與驅動程式

```text
你的程式碼
    │  使用統一的 API：db.QueryContext、db.ExecContext …
    ▼
database/sql（標準函式庫）  ← 定義介面、管理連線池、交易
    │  透過 driver 介面呼叫
    ▼
資料庫驅動程式（第三方）    ← go-sql-driver/mysql、pgx、sqlite …
    │  用各資料庫的通訊協定溝通
    ▼
MySQL / PostgreSQL / SQLite 伺服器
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>介面的威力（Ch 7）：</b> <code>database/sql</code> 只定義介面，不知道底下是哪一種資料庫。換資料庫時，大部分程式碼都不需要修改，只要換驅動程式和 SQL 語法差異。
</div>

<!--
Go 連接資料庫分成兩層。

上層是標準函式庫的 database/sql 套件，它提供統一的 API：查詢、執行、交易，並且幫我們管理「連線池」。下層是資料庫的驅動程式，由第三方提供，負責用各種資料庫的通訊協定溝通。

這就像萬用充電器：database/sql 是充電器本體，驅動程式是各種轉接頭。換一種資料庫，只要換一個轉接頭，我們的程式碼大部分都不用改。這就是第 7 章學的介面的威力：database/sql 定義介面，驅動程式實作介面。

MySQL 最主流的驅動程式是 go-sql-driver/mysql，這也是本課程唯一需要的第三方套件。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 安裝 MySQL 資料庫
## Installing MySQL

<!--
先把 MySQL 資料庫安裝好。
-->

---

# 安裝 MySQL Server

| 方式 | 指令／步驟 | 適合 |
| --- | --- | --- |
| **Docker（推薦）** | 見下方指令 | 不想弄亂電腦、要快速建立與刪除 |
| **Windows / macOS** | 到 [dev.mysql.com/downloads](https://dev.mysql.com/downloads/mysql/) 下載安裝程式 | 不熟悉 Docker |
| **macOS Homebrew** | `brew install mysql && brew services start mysql` | Mac 使用者 |
| **Ubuntu** | `sudo apt install mysql-server` | Linux 伺服器 |

```bash
docker run -d --name goshop-db -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=rootpass mysql:8.4
docker exec -it goshop-db mysql -uroot -p   # 進入 MySQL 命令列
```

<!--
安裝 MySQL 有好幾種方式，這門課推薦用 Docker。

用 Docker 的好處是：一行指令就能啟動一個乾淨的 MySQL，用完可以整個刪掉，不會在電腦裡留下任何東西，而且團隊裡每個人用的版本都一樣。-p 3306:3306 是把容器的 3306 連接埠對應到電腦上，MYSQL_ROOT_PASSWORD 設定 root 的密碼。mysql:8.4 是 MySQL 目前的長期支援版本。

如果不熟悉 Docker，也可以到 MySQL 官網下載安裝程式，一路下一步，過程中會要求設定 root 密碼。

安裝完成後，用 mysql 命令列工具，或是 MySQL Workbench、DBeaver 這類圖形化工具連進去。
-->

---

# 新增資料庫使用者

**不要用 root 帳號連線**；替應用程式建立專屬的使用者，只給需要的權限

```sql
-- 以 root 登入後執行
CREATE USER 'gouser'@'%' IDENTIFIED BY 'gopass';

-- 只授權 goshop 資料庫的所有權限
GRANT ALL PRIVILEGES ON goshop.* TO 'gouser'@'%';

-- 確認權限
SHOW GRANTS FOR 'gouser'@'%';
```

| 語法 | 意義 |
| --- | --- |
| `'gouser'@'%'` | 使用者 `gouser`，可以從**任何主機**連線（`localhost` 則只限本機） |
| `goshop.*` | `goshop` 資料庫中的**所有資料表** |

<!--
新增一個專屬應用程式的使用者，是資料庫安全的基本功。

為什麼不要用 root？root 擁有所有權限，可以刪除任何資料庫。如果程式有漏洞被入侵，攻擊者拿到的就是 root 權限。替應用程式建立專屬的使用者，只給它需要的權限，這叫做「最小權限原則」。

CREATE USER 建立使用者，@ 後面是允許從哪裡連線，百分比符號代表任何主機，用 Docker 時需要這樣設定。GRANT 授予權限，goshop.* 代表 goshop 資料庫裡的所有資料表。

實務上，正式環境通常不會給 ALL PRIVILEGES，而是只給 SELECT、INSERT、UPDATE、DELETE。
-->

---

# 建立一個 MySQL 資料庫

```sql
CREATE DATABASE goshop
  CHARACTER SET utf8mb4          -- 完整支援中文與表情符號
  COLLATE utf8mb4_unicode_ci;    -- 排序與比較規則

SHOW DATABASES;                  -- 確認建立成功
USE goshop;                      -- 切換到 goshop 資料庫
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>一定要用 utf8mb4：</b> MySQL 早期的 <code>utf8</code> 最多只支援 3 個位元組，存不了 🐹 這種 4 個位元組的字元（第 3 章的 UTF-8）。MySQL 8 的預設值已經是 <code>utf8mb4</code>。
</div>

<!--
接著建立資料庫。資料庫就像一個資料夾，裡面會放很多張資料表。

CHARACTER SET 要設定成 utf8mb4。這是 MySQL 一個很有名的歷史包袱：MySQL 早期的 utf8 其實是「閹割版」，每個字元最多只支援 3 個位元組。第 3 章學過，表情符號在 UTF-8 裡要 4 個位元組，所以用舊的 utf8 存表情符號會出錯。utf8mb4 才是完整的 UTF-8。

MySQL 8 以後預設就是 utf8mb4，但明確寫出來是好習慣。COLLATE 是排序和比較的規則，ci 代表不分大小寫。
-->

---

# 下載 Go 語言的 MySQL 驅動程式

```bash
mkdir goshop && cd goshop
go mod init example.com/goshop
go get github.com/go-sql-driver/mysql@latest
# go: added github.com/go-sql-driver/mysql v1.10.1
```

驅動程式透過**空白匯入**，在 `init()` 中向 `database/sql` 註冊自己（Ch 8）：

```go
import (
	"database/sql"

	// 只需要它的 init()：註冊名為 "mysql" 的驅動程式
	_ "github.com/go-sql-driver/mysql"
)
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 本章範例會直接使用驅動程式的 <code>mysql.Config</code> 產生連線字串，此時是一般的 import，同樣會執行 <code>init()</code> 完成註冊。
</div>

<!--
建立專案、用 go get 下載驅動程式，這是第 8 章學過的流程。

驅動程式的 import 方式很特別：前面加一個底線，這就是第 8 章說的「空白匯入」。我們不會直接呼叫驅動程式的任何函式，只需要它的 init 函式執行，init 裡會呼叫 sql.Register，向 database/sql 註冊一個名為 "mysql" 的驅動程式。之後 sql.Open("mysql", ...) 就能找到它。

本章的範例會使用驅動程式提供的 mysql.Config 來產生連線字串，這時候是一般的 import，一樣會執行 init，所以不需要另外寫空白匯入。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 以 Go 語言連接資料庫
## Connecting to MySQL

<!--
驅動程式準備好了，現在用 Go 連接資料庫。
-->

---

# 連線字串（DSN）

`sql.Open(驅動程式名稱, DSN)`；MySQL 驅動程式的 DSN 格式：

```text
使用者:密碼@tcp(主機:連接埠)/資料庫?參數1=值&參數2=值
gouser:gopass@tcp(127.0.0.1:3306)/goshop?parseTime=true
```

| 常用參數 | 說明 |
| --- | --- |
| `parseTime=true` | 把 `DATETIME` / `TIMESTAMP` 轉成 `time.Time`（否則是 `[]byte`） |
| `loc=Local` | 解析時間時使用的時區（預設 UTC） |
| `timeout=5s` | 建立連線的逾時時間 |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
🔐 <b>不要把密碼寫死在程式碼裡：</b> 實務上從環境變數（<code>os.Getenv("DB_PASS")</code>）或設定檔讀取，避免密碼被 commit 到 Git。
</div>

<!--
連接資料庫需要一個「連線字串」，也叫 DSN（Data Source Name），裡面包含使用者、密碼、主機、連接埠、資料庫名稱和參數。

最重要的參數是 parseTime=true：沒有它的話，資料庫的時間欄位會被讀成 []byte，而不是第 10 章學的 time.Time。這是初學者最常忘記的設定。

另外提醒一個資安觀念：密碼不要寫死在程式碼裡。程式碼會被 commit 到 Git，密碼就跟著外洩了。實務上從環境變數讀取，本章為了方便示範，才直接寫在程式碼裡。
-->

---

# 開啟資料庫連線：openDB

```go
// 檔名：db.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

func openDB(ctx context.Context) (*sql.DB, error) {
	cfg := mysql.NewConfig() // 用 Config 組出 DSN，不用手動拼字串
	cfg.User, cfg.Passwd = "gouser", "gopass"
	cfg.Net, cfg.Addr = "tcp", "127.0.0.1:3306"
	cfg.DBName = "goshop"
	cfg.ParseTime = true // DATETIME 自動轉成 time.Time
```

<!--
我們把開啟連線寫成一個函式 openDB，放在 db.go 檔案裡。這一章的範例會分成好幾個檔案，都屬於同一個 main 套件，最後組成一個完整的程式。

組 DSN 的時候，推薦用驅動程式提供的 mysql.NewConfig，一個欄位一個欄位設定，再用 FormatDSN 產生字串。這樣比自己拼字串安全，密碼裡有特殊符號也會自動處理。
-->

---

# 開啟資料庫連線：openDB（續）

```go
	// 續上頁
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)                 // 最多同時開啟 10 條連線
	db.SetMaxIdleConns(5)                  // 最多保留 5 條閒置連線
	db.SetConnMaxLifetime(5 * time.Minute) // 連線最長使用時間

	if err := db.PingContext(ctx); err != nil { // 真正嘗試連線
		db.Close()
		return nil, fmt.Errorf("無法連線到資料庫：%w", err)
	}
	return db, nil
}
```

<!--
sql.Open 回傳一個 *sql.DB。

這裡有一個非常重要的觀念：sql.Open「不會」真的去連線，它只是準備好設定。所以就算密碼打錯，sql.Open 也不會回傳錯誤。要確認能不能連線，一定要呼叫 PingContext，它才會真正建立連線。

中間三行是連線池的設定。*sql.DB 不是「一條連線」，而是一個「連線池」：它會自動管理多條連線，需要時借出、用完歸還。設定最大連線數，可以避免程式把資料庫的連線用光。
-->

---

# 使用 *sql.DB 的注意事項

**注意事項之一：** `sql.Open` **不會真的連線**，一定要 `PingContext` 確認

**注意事項之二：** `*sql.DB` 是**連線池**，應該在程式啟動時建立**一次**、全程共用；它可以安全地被多個 goroutine 同時使用

**注意事項之三：** 所有操作都使用 **`XxxContext`** 版本，搭配逾時避免無限等待

```go
bg := context.Background()
ctx, cancel := context.WithTimeout(bg, 5*time.Second)
defer cancel()
db, err := openDB(ctx)
if err != nil {
	log.Fatal(err) // 啟動時連不上資料庫，直接結束
}
defer db.Close() // 程式結束時才關閉
```

<!--
使用 *sql.DB 有三個注意事項。

第一，sql.Open 不會真的連線，一定要 Ping。

第二，*sql.DB 是連線池，不是一條連線。很多初學者會在每次查詢前 Open、查完就 Close，這樣完全失去了連線池的好處，效能會很差。正確的做法是程式啟動時建立一次，整個程式共用，程式結束時才關閉。

第三，所有的資料庫操作都有兩個版本：Query 和 QueryContext、Exec 和 ExecContext。請一律使用 Context 版本，搭配 context.WithTimeout 設定逾時。這樣如果資料庫卡住了，程式最多等幾秒就會放棄，而不是無限等待。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 建立、清空和移除資料表
## CREATE・TRUNCATE・DROP

<!--
連線成功之後，第一件事是建立資料表。
-->

---

# 資料表的結構

```sql
CREATE TABLE IF NOT EXISTS products (
  id         BIGINT AUTO_INCREMENT PRIMARY KEY,  -- 主鍵，自動遞增
  name       VARCHAR(100) NOT NULL,              -- 不可為空
  price      INT NOT NULL,
  note       VARCHAR(255) NULL,                  -- 可以是 NULL
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

| MySQL 型別 | Go 型別 |
| --- | --- |
| `BIGINT` / `INT` | `int64` / `int` |
| `VARCHAR` / `TEXT` | `string` |
| `DATETIME`（需 `parseTime=true`） | `time.Time` |
| 可以是 `NULL` 的欄位 | `sql.Null[T]`（Go 1.22+）或指標 `*T` |

<!--
這是我們要建立的商品資料表，有五個欄位。

id 是主鍵，AUTO_INCREMENT 代表新增資料時自動產生遞增的編號。name 和 price 是 NOT NULL，一定要有值。note 允許 NULL，代表可以沒有備註。created_at 預設是新增資料當下的時間。

表格是 MySQL 型別和 Go 型別的對應。要特別注意可以是 NULL 的欄位：Go 的 string 不能是 nil，所以要用 sql.Null[string] 或 *string 來接，等一下查詢的時候會詳細說明。
-->

---

# 用 ExecContext 建立、清空、移除資料表

```go
// 檔名：store.go
package main

import (
	"context"
	"database/sql"
)

type Store struct{ db *sql.DB } // 把資料庫操作集中在 Store

func (s *Store) CreateTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id         BIGINT AUTO_INCREMENT PRIMARY KEY,
			name       VARCHAR(100) NOT NULL,
			price      INT NOT NULL,
			note       VARCHAR(255) NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`)
	return err
}
```

<!--
不會回傳資料列的 SQL，例如 CREATE、INSERT、UPDATE、DELETE，都用 ExecContext 執行。

我們定義一個 Store 型別，把所有資料庫操作都寫成它的方法。這是實務上很常見的做法，叫做 Repository 模式：把「怎麼存取資料」集中在一個地方，其他程式碼只要呼叫 Store 的方法就好，不需要知道 SQL 怎麼寫。第 15 章寫 API 的時候會再用到這個設計。

SQL 用第 3 章學的反引號原始字串，可以直接換行，寫起來清楚很多。IF NOT EXISTS 讓程式重複執行也不會出錯。
-->

---

# 清空與移除資料表

```go
// 續上頁
func (s *Store) Truncate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "TRUNCATE TABLE products")
	return err
}

func (s *Store) Drop(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DROP TABLE IF EXISTS products")
	return err
}
```

| SQL | 效果 | 可以復原嗎 |
| --- | --- | --- |
| `DELETE FROM products` | 刪除所有資料列，一筆一筆刪 | 可在交易中 Rollback |
| `TRUNCATE TABLE products` | **清空**資料，自動編號重設為 1 | ❌ |
| `DROP TABLE products` | **移除整張資料表**（結構與資料） | ❌ |

<!--
清空和移除資料表也都用 ExecContext。

表格比較了三種「刪除」的差別。DELETE 是一筆一筆刪除資料，可以加上 WHERE 條件只刪部分資料，也可以在交易裡復原。TRUNCATE 是直接清空整張表，速度快，自動編號會重新從 1 開始，但不能復原。DROP 最徹底，連資料表本身都移除了。

TRUNCATE 和 DROP 在正式環境要非常小心，通常只在測試環境、或初始化資料時使用。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 插入資料
## INSERT

<!--
資料表建好了，接下來新增資料。
-->

---
zoom: 0.84
---

# 插入資料：使用佔位符 ?

```go
// 檔名：product.go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Product struct {
	ID        int64
	Name      string
	Price     int
	Note      sql.Null[string] // 可以是 NULL 的欄位（Go 1.22+）
	CreatedAt time.Time
}

func (s *Store) Add(ctx context.Context, p Product) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		"INSERT INTO products (name, price) VALUES (?, ?)",
		p.Name, p.Price)
	if err != nil {
		return 0, fmt.Errorf("新增 %s 失敗：%w", p.Name, err)
	}
	return res.LastInsertId() // 取得自動產生的 id
}
```

<!--
新增資料用 INSERT，一樣是 ExecContext。

注意 SQL 裡的問號，這叫做「佔位符」。實際的值寫在後面的參數，依照順序對應到每一個問號。驅動程式會安全地把值放進去。

ExecContext 回傳一個 sql.Result，LastInsertId 可以取得資料庫自動產生的 id，RowsAffected 可以取得影響了幾筆資料。

Product 的 Note 欄位用了 sql.Null[string]，這是 Go 1.22 加入的泛型型別，用來表示「可能是 NULL 的值」，等一下查詢的時候會說明。
-->

---

# 使用 SQL 的注意事項：SQL Injection

**絕對不要用字串串接組 SQL**，一律使用佔位符 `?`

```go
name := "x'); DROP TABLE products; --" // 惡意的使用者輸入

// ❌ 危險：使用者輸入直接變成 SQL 的一部分
query := fmt.Sprintf(
	"INSERT INTO products (name, price) VALUES ('%s', 0)", name)

// ✅ 安全：值和 SQL 分開傳送，資料庫不會把它當成指令
db.ExecContext(ctx,
	"INSERT INTO products (name, price) VALUES (?, ?)", name, 0)
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
🔐 <b>SQL Injection</b> 常年位居 OWASP 十大網站安全風險。佔位符只能放「值」，資料表名稱、欄位名稱不能用 <code>?</code>，需要動態指定時請用<b>白名單</b>檢查。
</div>

<!--
這是資料庫程式設計最重要的安全觀念：SQL Injection，SQL 注入攻擊。

如果用 Sprintf 把使用者的輸入直接串進 SQL，惡意的使用者可以輸入一段精心設計的字串，把原本的 SQL 提前結束，再接上自己的指令，例如 DROP TABLE，整張資料表就被刪掉了。

使用佔位符的時候，SQL 指令和值是分開傳送給資料庫的，不管值裡面有什麼奇怪的字元，資料庫都只會把它當成「資料」，不會當成「指令」執行。

所以規則很簡單：永遠不要用字串串接組 SQL，一律用問號佔位符。
-->

---
zoom: 0.93
---

# 交易（transaction）與預備敘述

一次新增多筆資料時，用**交易**確保「**全部成功，或全部不做**」

```go
// 續上頁
func (s *Store) AddMany(ctx context.Context, items []Product) error {
	tx, err := s.db.BeginTx(ctx, nil) // 開始交易
	if err != nil {
		return err
	}
	defer tx.Rollback() // 中途 return 時自動復原；Commit 後呼叫不會有作用

	stmt, err := tx.PrepareContext(ctx, // 預備敘述：SQL 只解析一次
		"INSERT INTO products (name, price, note) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range items {
		_, err := stmt.ExecContext(ctx, p.Name, p.Price, p.Note)
		if err != nil {
			return fmt.Errorf("新增 %s 失敗：%w", p.Name, err)
		}
	}
	return tx.Commit() // 全部成功才提交
}
```

<!--
什麼是交易？想像銀行轉帳：從 A 帳戶扣錢、在 B 帳戶加錢，這兩個動作必須「全部成功」或「全部不做」，不能扣了錢卻沒加到。交易就是用來保證這件事的。

BeginTx 開始一個交易，之後的操作都透過 tx 執行。全部成功就 Commit 提交；中途有任何錯誤就 Rollback 復原，好像什麼都沒發生過。

這裡用了一個很漂亮的技巧：defer tx.Rollback()。如果中途出錯提早 return，defer 會自動復原；如果順利執行到 Commit，之後的 Rollback 不會有任何作用。這樣就不用在每個錯誤處理的地方都寫 Rollback。

PrepareContext 建立「預備敘述」：SQL 只解析一次，之後重複執行很多次，適合大量新增資料。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 查詢資料
## SELECT

<!--
資料新增好了，接下來查詢資料。
-->

---

# 查詢並印出整個資料表內容：QueryContext

```go
// 續上頁
func (s *Store) List(ctx context.Context) ([]Product, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, price, note, created_at "+
			"FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // ⚠️ 一定要關閉，否則連線不會歸還連線池

	var list []Product
	for rows.Next() { // 逐筆讀取
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Note, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err() // 檢查走訪過程中是否發生錯誤
}
```

<!--
查詢多筆資料用 QueryContext，它回傳 *sql.Rows，可以想成一個「游標」，用 rows.Next() 一筆一筆往下讀。

每一筆資料用 rows.Scan 讀出來，參數是指標，依照 SELECT 的欄位順序對應。欄位數量和順序一定要跟 SELECT 一樣，否則會出錯。

有兩個一定要做的事：第一，defer rows.Close()，否則這條連線會一直被佔用，不會歸還連線池，久了連線就會被用光。第二，迴圈結束後檢查 rows.Err()，確認是「讀完了」而不是「中途出錯」。這跟上一章 bufio.Scanner 的 sc.Err() 是同樣的模式。
-->

---

# 處理 NULL：sql.Null[T]

資料庫的 `NULL` 無法放進 Go 的 `string`、`int`，需要能表示「沒有值」的型別

| 寫法 | 有值時 | 是 NULL 時 |
| --- | --- | --- |
| `sql.Null[string]`（Go 1.22+） | `{V: "限量", Valid: true}` | `{V: "", Valid: false}` |
| `sql.NullString`（舊寫法） | `{String: "限量", Valid: true}` | `{String: "", Valid: false}` |
| `*string` | 指向值的指標 | `nil` |

```go
for _, p := range list {
	note := "（無）"
	if p.Note.Valid { // 先檢查是不是 NULL
		note = p.Note.V
	}
	created := p.CreatedAt.Format(time.DateTime)
	fmt.Println(p.ID, p.Name, p.Price, note, created)
}
```

<!--
資料庫的 NULL 代表「沒有值」，但 Go 的 string 不能是 nil。如果把 NULL 掃描進 string，會出現錯誤。

解法有三種。Go 1.22 加入的 sql.Null[T] 是泛型版本，任何型別都能用，V 是值，Valid 代表是不是有值。舊的寫法是 sql.NullString、sql.NullInt64 這些，每種型別一個，功能一樣。也可以用指標，NULL 的時候是 nil。

使用的時候，先檢查 Valid，是 true 才讀取 V。這跟第 4 章 map 的 comma ok 慣用法很像：先確認有沒有，再使用。
-->

---

# 查詢單筆資料：QueryRowContext

```go
// 續上頁
var ErrNotFound = errors.New("找不到商品")

func (s *Store) Get(ctx context.Context, id int64) (Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, price FROM products WHERE id = ?", id,
	).Scan(&p.ID, &p.Name, &p.Price)
	if errors.Is(err, sql.ErrNoRows) { // 查無資料
		return p, ErrNotFound
	}
	return p, err
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>QueryRowContext 不需要 Close：</b> <code>Scan</code> 完成後會自動釋放連線。查無資料時，<code>Scan</code> 回傳哨兵錯誤 <code>sql.ErrNoRows</code>。
</div>

<!--
只需要查詢一筆資料的時候，用 QueryRowContext，後面直接串 Scan。它比 QueryContext 簡單：不用迴圈、不用 Close，Scan 完成後連線會自動歸還。

查不到資料的時候，Scan 會回傳 sql.ErrNoRows，這是第 6 章學的哨兵錯誤。我們把它轉換成自己定義的 ErrNotFound，這樣呼叫 Store 的程式碼就不需要知道底層用的是 database/sql，只要判斷 ErrNotFound 就好。第 15 章寫 API 的時候，ErrNotFound 就會對應到 HTTP 404。
-->

---

# 查詢符合條件的資料

```go
// 續上頁
func (s *Store) FindCheaper(ctx context.Context,
	max int) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT name FROM products WHERE price <= ? ORDER BY price", max)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		names = append(names, n)
	}
	return names, rows.Err()
}
```

<!--
加上 WHERE 條件，就能查詢符合條件的資料。條件的值一樣用問號佔位符，絕對不要用字串串接。

這個方法查詢價格小於等於 max 的商品名稱，依價格排序。結構跟 List 一模一樣：QueryContext、defer Close、迴圈 Scan、檢查 rows.Err()。

實務上這樣的程式碼會寫很多次，模式都一樣，熟悉之後寫起來會很快。
-->

---
layout: default
---

# 練習 1：商品查詢
### 任務說明

替 `Store` 新增方法 `Search(ctx context.Context, keyword string) ([]Product, error)`：

1. 查詢名稱**包含**關鍵字的商品（SQL：`WHERE name LIKE ?`）
2. 關鍵字前後加上 `%`：`"%" + keyword + "%"`（在 Go 裡組，不要串進 SQL）
3. 依價格由高到低排序，最多回傳 10 筆（`ORDER BY price DESC LIMIT 10`）
4. `note` 欄位是 NULL 時，印出「（無備註）」
5. 在 `main` 中搜尋「茶」，印出結果

<!--
這個練習要大家自己寫一個查詢方法。

注意第 2 步：LIKE 的萬用字元百分比符號，要加在「值」上面，而不是寫在 SQL 字串裡跟問號串接。這樣才能維持佔位符的安全性。
-->

---
zoom: 0.84
---

# 練習 1：解題提示
### 提示說明

```go
// 檔名：search.go
package main

import "context"

func (s *Store) Search(ctx context.Context,
	keyword string) ([]Product, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, price, note, created_at FROM products
		WHERE name LIKE ? ORDER BY price DESC LIMIT 10`,
		"%"+keyword+"%") // 萬用字元加在值上
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price,
			&p.Note, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}
```

<!--
Search 的結構跟 List 一樣，差別在 WHERE 條件和排序。

LIKE 的百分比符號加在參數上，"%" + keyword + "%"，驅動程式會把整個字串當成一個值安全地傳送。

印出結果時，檢查 p.Note.Valid，是 false 就印「（無備註）」。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 更新既有資料
## UPDATE & DELETE

<!--
最後是更新和刪除資料。
-->

---

# 更新既有資料：UPDATE

```go
// 續上頁
func (s *Store) UpdatePrice(ctx context.Context,
	id int64, price int) error {
	res, err := s.db.ExecContext(ctx,
		"UPDATE products SET price = ? WHERE id = ?", price, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected() // 實際被修改的筆數
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ <b>UPDATE / DELETE 一定要有 WHERE：</b> 少了 WHERE，整張資料表的資料都會被修改或刪除！
</div>

<!--
更新資料用 UPDATE，一樣是 ExecContext。

RowsAffected 回傳實際被修改的筆數。如果是 0，代表找不到這個 id 的商品，我們回傳 ErrNotFound。

使用 UPDATE 最重要的注意事項：一定要有 WHERE 條件。少了 WHERE，資料表裡的每一筆資料都會被更新成同一個價格。這是資料庫操作最可怕的失誤之一，很多公司都發生過。

補充一個 MySQL 的細節：如果更新的新值跟舊值一樣，MySQL 的 RowsAffected 會回傳 0，因為實際上沒有修改任何資料。
-->

---

# 刪除資料與錯誤處理

```go
// 續上頁
func (s *Store) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM products WHERE id = ?", id)
	return err
}
```

判斷 MySQL 特定的錯誤，用 `errors.As` 取出驅動程式的錯誤型別：

```go
var me *mysql.MySQLError
if errors.As(err, &me) {
	switch me.Number {
	case 1062: // Duplicate entry：違反唯一性限制
		return ErrDuplicate
	case 1146: // Table doesn't exist
		return fmt.Errorf("資料表不存在：%w", err)
	}
}
```

<!--
刪除資料用 DELETE，一樣要記得 WHERE 條件。

另外介紹一個實務上很常用的技巧：判斷 MySQL 特定的錯誤。驅動程式的錯誤型別是 *mysql.MySQLError，它有一個 Number 欄位，是 MySQL 的錯誤代碼。用第 6 章學的 errors.As 取出來，就能根據代碼做不同的處理。

最常見的是 1062：違反唯一性限制，例如註冊時 email 重複了。這時候應該回傳「這個 email 已經被註冊」，而不是「系統錯誤」。
-->

---
zoom: 0.94
---

# 把所有功能串起來：main.go

```go
// 檔名：main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func main() {
	bg := context.Background()
	ctx, cancel := context.WithTimeout(bg, 10*time.Second)
	defer cancel()
	db, err := openDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	s := &Store{db: db}
	if err := s.CreateTable(ctx); err != nil {
		log.Fatal(err)
	}
```

<!--
最後把前面寫的所有方法在 main.go 裡串起來。

先建立一個 10 秒逾時的 context，開啟資料庫連線，失敗就結束程式。然後建立 Store，建立資料表。
-->

---
zoom: 0.91
---

# 把所有功能串起來：main.go（續）

```go
	// 續上頁
	s.Truncate(ctx) // 每次執行前清空，方便重複練習
	s.Add(ctx, Product{Name: "咖啡豆", Price: 450})
	s.AddMany(ctx, []Product{
		{Name: "綠茶", Price: 35},
		{Name: "蛋糕", Price: 120,
			Note: sql.Null[string]{V: "限量", Valid: true}},
	})
	s.UpdatePrice(ctx, 2, 40)
	list, err := s.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range list {
		fmt.Println(p.ID, p.Name, p.Price, p.Note.V)
	}
	_, err = s.Get(ctx, 99)
	fmt.Println(err)                     // 找不到商品
	fmt.Println(s.FindCheaper(ctx, 200)) // [綠茶 蛋糕] <nil>
}
```

```text
1 咖啡豆 450
2 綠茶 40
3 蛋糕 120 限量
```

<!--
清空資料表之後，新增一筆、再用交易新增兩筆，接著把綠茶的價格改成 40。

List 列出所有商品，Get 查詢一個不存在的 id，得到「找不到商品」；FindCheaper 查詢 200 元以下的商品。

注意：為了讓投影片精簡，這裡有幾個呼叫沒有檢查錯誤，實際寫程式時每個錯誤都要處理。

執行 go run . 就能看到結果。因為所有檔案都是 main 套件，go run 點會把它們一起編譯。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 練習：FizzBuzz 統計表
## Putting It All Together

<!--
最後是本章的綜合練習：把第 2 章的 FizzBuzz 存進資料庫，再用 SQL 做統計。
-->

---
layout: default
---

# 綜合練習：FizzBuzz 統計表
### 任務說明

1. 建立資料表 `fizzbuzz`：`n INT PRIMARY KEY`、`word VARCHAR(10) NOT NULL`
2. 寫函式 `fizzBuzz(n int) string`（第 2 章的規則）
3. 在**一個交易**中，用**預備敘述**把 1～100 的結果存進資料表（先清空舊資料）
4. 用一個 SQL 查詢統計每一種結果的次數，數字統一歸類為「數字」：

```text
數字        53
Fizz      27
Buzz      14
FizzBuzz   6
```

提示：`CASE WHEN word REGEXP '^[0-9]+$' THEN '數字' ELSE word END` + `GROUP BY`

<!--
這個綜合練習把今天學的所有東西都用上了：建立資料表、交易、預備敘述、查詢、Scan。

統計的部分，把工作交給 SQL：GROUP BY 可以依照某個欄位分組計數。因為數字每一個都不一樣，直接 GROUP BY 會有 53 組，所以用 CASE WHEN 把所有的數字歸類成同一組「數字」。
-->

---

# 綜合練習：解題提示
### 提示說明

```go
// 檔名：fizzbuzz.go
package main

import (
	"context"
	"fmt"
	"strconv"
)

func fizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	}
	return strconv.Itoa(n)
}
```

<!--
fizzBuzz 函式就是第 2 章的邏輯，改寫成回傳字串。注意 15 的倍數要最先判斷。

最後的 return strconv.Itoa(n) 把數字轉成字串，第 3 章學過不能用 string(n)。
-->

---
zoom: 0.75
---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func (s *Store) SaveFizzBuzz(ctx context.Context, max int) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS fizzbuzz (
			n INT PRIMARY KEY, word VARCHAR(10) NOT NULL)`)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "DELETE FROM fizzbuzz")
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO fizzbuzz (n, word) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for n := 1; n <= max; n++ {
		if _, err := stmt.ExecContext(ctx, n, fizzBuzz(n)); err != nil {
			return err
		}
	}
	return tx.Commit()
}
```

<!--
SaveFizzBuzz 先建立資料表，然後開始交易。在交易裡先用 DELETE 清空舊資料，再用預備敘述新增 1 到 100 的結果，全部成功才 Commit。

這裡用 DELETE 而不是 TRUNCATE，因為 DELETE 在交易裡可以被復原；而 TRUNCATE 在 MySQL 裡會自動提交交易，破壞交易的完整性。

100 筆資料在一個交易裡新增，比一筆一筆自動提交快很多。
-->

---
zoom: 0.97
---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func (s *Store) FizzBuzzStats(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT CASE WHEN word REGEXP '^[0-9]+$' THEN '數字'
		            ELSE word END AS kind,
		       COUNT(*) AS cnt
		FROM fizzbuzz GROUP BY kind ORDER BY cnt DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var kind string
		var cnt int
		if err := rows.Scan(&kind, &cnt); err != nil {
			return err
		}
		fmt.Printf("%-8s %3d\n", kind, cnt)
	}
	return rows.Err()
}
```

<!--
統計的部分全部交給 SQL：CASE WHEN 把純數字的 word 歸類成「數字」，其他保持原樣；GROUP BY kind 依類別分組，COUNT(*) 計數，ORDER BY cnt DESC 依數量由多到少排序。

Go 這邊只需要讀出兩個欄位並印出來。在 main 裡呼叫 s.SaveFizzBuzz(ctx, 100) 和 s.FizzBuzzStats(ctx) 就完成了。

結果是數字 53 個、Fizz 27 個、Buzz 14 個、FizzBuzz 6 個，加起來剛好 100。
-->

---

# 章節總結

- **架構**：`database/sql` 定義介面與連線池，驅動程式（`go-sql-driver/mysql`）負責連線；用**空白匯入**或一般匯入註冊驅動程式
- **準備**：建立專屬使用者（最小權限）；資料庫使用 **`utf8mb4`**；DSN 加上 **`parseTime=true`**
- **連線**：`sql.Open` 不會連線，要 `PingContext`；`*sql.DB` 是連線池，**建立一次、全程共用**
- **執行**：`ExecContext` 用於 CREATE / INSERT / UPDATE / DELETE；`LastInsertId`、`RowsAffected`
- **安全**：**一律使用佔位符 `?`**，杜絕 SQL Injection
- **查詢**：`QueryContext` + `defer rows.Close()` + `rows.Next()` + `Scan` + `rows.Err()`；單筆用 `QueryRowContext`，查無資料是 `sql.ErrNoRows`
- **NULL 與交易**：`sql.Null[T]`（1.22+）；`BeginTx` + `defer tx.Rollback()` + `Commit`

下一章我們會介紹「HTTP 客戶端」：用 Go 呼叫網路上的 API。

<!--
我們來整理今天學到的東西。

Go 用 database/sql 加上驅動程式連接資料庫，*sql.DB 是連線池，程式啟動時建立一次。所有操作用 Context 版本，執行用 ExecContext、查詢用 QueryContext 和 QueryRowContext。最重要的安全觀念：一律使用佔位符，杜絕 SQL Injection。查詢完記得 Close 和檢查 rows.Err()，需要「全部成功或全部不做」的時候用交易。

到這裡，我們已經能把資料存在資料庫裡了。接下來兩章要進入網路的世界：下一章學 HTTP 客戶端，用 Go 呼叫網路上的 API、取得 JSON 資料；第 15 章學 HTTP 伺服器，把今天的資料庫操作包裝成 API，讓別人呼叫。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
資料庫是後端開發的核心，今天的內容建議大家一定要自己動手做一次，特別是安裝 MySQL 和建立使用者的部分。

課後建議：替 products 資料表加上一個 category 欄位，寫一個方法統計每個分類的商品數量和平均價格，練習 GROUP BY 和 AVG。

有問題的同學現在可以提問！
-->

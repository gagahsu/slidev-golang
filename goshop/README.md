# GoShop — 課程專案參考解答

「Go 實戰開發」每一章的最後，都會在同一個專案 **GoShop（迷你電商後台）** 上加入新功能。
這個資料夾存放每一章完成後的**完整參考解答**：`ch05/` 就是做完第 5 章「GoShop 專案實作」後的樣子。

| 步驟 | 新增的功能 | 主要檔案 |
| --- | --- | --- |
| ch00 | 建立專案、印出歡迎訊息 | `main.go` |
| ch01 | 以變數、常數與 iota 試算一筆訂單；用指標扣庫存 | `main.go` |
| ch02 | 會員折扣（switch）、滿額免運（if）、數量試算表（for） | `main.go` |
| ch03 | 解析供應商的商品字串（strings、strconv）、對齊中文品名 | `main.go` |
| ch04 | `Product`、`Cart` 結構；以 map 做商品目錄；方法 | `main.go` |
| ch05 | 閉包折扣規則、參數不定函式、千分位金額、defer 收據 | `main.go` |
| ch06 | 哨兵錯誤、自訂 `StockError`、`errors.Join` 一次回報所有問題 | `main.go` |
| ch07 | `Money` 實作 `fmt.Stringer`；`PaymentMethod` 介面與三種付款方式 | `main.go` |
| ch08 | 拆成 `money`、`shop`、`store`、`checkout`、`payment` 套件 | `internal/…` |
| ch09 | 表格驅動測試（抓到負數金額的 bug）、benchmark、`slog` | `*_test.go` |
| ch10 | 有使用期限的折價券、預計出貨日（跳過週末）、時區 | `checkout/coupon.go` |
| ch11 | JSON 標籤、從 `products.json` 匯入、gob 快照、`MarshalText` | `store/snapshot.go` |
| ch12 | 命令列旗標、CSV 匯入匯出、檔案權限、JSON Lines 訂單日誌 | `main.go`、`store/csv.go` |
| ch13 | `Store` 介面、MySQL 實作、交易扣庫存、合約測試 | `store/mysql*.go` |
| ch14 | 匯率 API 客戶端、webhook 通知（POST） | `rates/`、`webhook/` |
| ch15 | RESTful JSON API、後台模板網頁、表單、靜態資源、優雅關閉 | `web/` |
| ch16 | 互斥鎖防止超賣、背景通知 worker pool | `store/memory.go`、`webhook/notifier.go` |
| ch17 | `-version`（ldflags）、Makefile 跨平台編譯、Example 測試、文件註解 | `Makefile`、`doc.go` |
| ch18 | bcrypt 管理員密碼、HMAC 簽章的登入憑證、HTTPS（TLS 1.3） | `auth/`、`web/login.go` |
| ch19 | 用反射讀取 `validate` 標籤的驗證器 | `validate/` |

## 執行

```bash
cd goshop/ch19
go run . -list                          # 列出商品（第一次執行會從 data/products.json 匯入）
go run . -buy SKU-001:2,SKU-003:1 -coupon WEEK15
go test ./...                           # 執行所有測試

# 啟動網站（第 18 章起需要管理員密碼）
export GOSHOP_ADMIN_HASH=$(echo 'my-password' | go run . -hash-password)
go run . -http :8080                    # API：http://localhost:8080/api/products
                                        # 後台：http://localhost:8080/admin
```

第 13 章起可以改用 MySQL：

```sql
CREATE DATABASE goshop_app CHARACTER SET utf8mb4;
GRANT ALL ON goshop_app.* TO 'gouser'@'localhost';
```

```bash
export GOSHOP_DSN="gouser:gopass@tcp(127.0.0.1:3306)/goshop_app"
go run . -import data/new-products.csv
```

## 檢查

```bash
pnpm check:project        # 每一步都要 gofmt、go vet、go test 通過，且投影片上的程式碼和這裡一致
```

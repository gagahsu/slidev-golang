# Go 實戰開發（Go Programming Masterclass）

以 [Slidev](https://sli.dev) 製作的 Go 語言教學投影片，共 20 章（前置作業 + Chapter 1～19），
所有範例以 **Go 1.27** 撰寫並實際編譯、執行驗證，採用目前主流的寫法與標準函式庫。

| 章 | 主題 | 章 | 主題 |
| --- | --- | --- | --- |
| 0 | 前置作業：安裝 Go / VS Code、建立專案、Playground | 10 | 時間處理 |
| 1 | 變數與算符 | 11 | 編碼／解碼 JSON（含 `encoding/json/v2`、gob） |
| 2 | 條件判斷與迴圈 | 12 | 系統與檔案：旗標、訊號、檔案、CSV |
| 3 | 核心型別 | 13 | SQL 與資料庫（MySQL） |
| 4 | 複合型別 | 14 | HTTP 客戶端 |
| 5 | 函式（含泛型、迭代器補充） | 15 | HTTP 伺服器與 RESTful API |
| 6 | 錯誤處理 | 16 | 並行性運算：goroutine、channel、context |
| 7 | 介面 | 17 | Go 語言工具：build、vet、fix、doc、tool |
| 8 | 套件與 Go Modules | 18 | 加密安全：SHA-2/3、AES-GCM、RSA-OAEP、Ed25519、TLS |
| 9 | 除錯：fmt、log / slog、單元測試 | 19 | reflect 與 unsafe |

## 課程專案：GoShop 迷你電商後台

每一章的最後都有一節「**GoShop 專案實作**」，把這一章學到的東西用在同一個專案上：
第 0 章建立專案，前 7 章寫在單一 `main.go`，第 8 章拆成套件，之後陸續加上測試、JSON／gob 存檔、
命令列工具、MySQL、REST API 與後台網頁、並行安全、HTTPS 登入，到第 19 章完成一個約 2,700 行、
10 個套件、全部有測試的系統。

每一步的完整參考解答在 [`goshop/ch00`](goshop/) ～ [`goshop/ch19`](goshop/)，說明見 [goshop/README.md](goshop/README.md)。

## 使用方式

```bash
pnpm install
pnpm dev          # http://localhost:3030 ，目錄頁可點選各章節
pnpm run ch05     # 只開啟單一章節
pnpm build        # 產生靜態網站到 dist/
```

每張投影片都附有講稿（Presenter Mode：按 `p` 或開啟 `/presenter`）。

## 驗證投影片中的 Go 程式碼

```bash
GOTOOLCHAIN=go1.27.0 pnpm check:go            # 編譯 + go vet 所有可獨立編譯的範例
GOTOOLCHAIN=go1.27.0 pnpm check:go ch16 --run # 逐一執行，對照投影片上標註的輸出
WITH_DEPS=1 pnpm check:go ch13                # 一併檢查使用第三方模組（MySQL 驅動程式）的範例
GOTOOLCHAIN=go1.27.0 pnpm check:project        # GoShop：每一步 gofmt／vet／test，並比對投影片摘錄
```

詳細的撰寫規範請見 [CLAUDE.md](CLAUDE.md)。

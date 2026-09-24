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
```

詳細的撰寫規範請見 [CLAUDE.md](CLAUDE.md)。

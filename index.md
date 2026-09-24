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
title: Go Programming Masterclass
routeAlias: home
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

<style>
.chapter-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  width: 100%;
  max-width: 960px;
  margin-top: 1.2rem;
}
.chapter-card {
  display: block;
  background: #f0faf9;
  border: 2px solid #5eada0;
  border-radius: 12px;
  padding: 1.2rem 0.8rem;
  text-decoration: none !important;
  color: #1a5c5c !important;
  transition: all 0.2s ease;
}
.chapter-card:hover {
  background: #5eada0;
  color: white !important;
  transform: translateY(-3px);
  box-shadow: 0 6px 16px rgba(94, 173, 160, 0.35);
}
.chapter-card:hover .chapter-subtitle {
  color: rgba(255,255,255,0.85) !important;
}
.chapter-num {
  font-size: 1.6rem;
  font-weight: 900;
  margin-bottom: 0.3rem;
}
.chapter-subtitle {
  font-size: max(13px, 0.88rem);
  color: #4a7c7c;
  margin-top: 0.3rem;
}
.chapter-card-adv {
  border-color: #e0a96d;
}
.chapter-card-adv .chapter-badge {
  color: #c97b2c;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.1em;
}
</style>

<div class="flex flex-col items-center h-full" style="background: #ffffff; overflow-y: auto; padding: 1.5rem 0;">
  <p style="color: #5eada0; font-size: 1rem; font-weight: 600; letter-spacing: 0.2em; text-transform: uppercase; margin-bottom: 1rem;">Go Programming Masterclass</p>
  <h1 style="color: #1a5c5c; font-size: 2.8rem; font-weight: 900; line-height: 1.2; margin-bottom: 0.5rem;">Go 實戰開發・課程目錄</h1>
  <div style="height: 4px; width: 240px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 0.5rem;"></div>
  <p style="color: #9dc4c4; font-size: 0.9rem; margin-bottom: 0;">點擊章節卡片開始學習｜範例以 Go 1.27 撰寫並驗證</p>
  <div class="chapter-grid">
    <Link to="ch00" class="chapter-card">
      <div class="chapter-num">Ch 0</div>
      <div>前置作業</div>
      <div class="chapter-subtitle">Setup &amp; First Program</div>
    </Link>
    <Link to="ch01" class="chapter-card">
      <div class="chapter-num">Ch 1</div>
      <div>變數與算符</div>
      <div class="chapter-subtitle">Variables &amp; Operators</div>
    </Link>
    <Link to="ch02" class="chapter-card">
      <div class="chapter-num">Ch 2</div>
      <div>條件判斷與迴圈</div>
      <div class="chapter-subtitle">Conditionals &amp; Loops</div>
    </Link>
    <Link to="ch03" class="chapter-card">
      <div class="chapter-num">Ch 3</div>
      <div>核心型別</div>
      <div class="chapter-subtitle">Core Types</div>
    </Link>
    <Link to="ch04" class="chapter-card">
      <div class="chapter-num">Ch 4</div>
      <div>複合型別</div>
      <div class="chapter-subtitle">Composite Types</div>
    </Link>
    <Link to="ch05" class="chapter-card">
      <div class="chapter-num">Ch 5</div>
      <div>函式</div>
      <div class="chapter-subtitle">Functions</div>
    </Link>
    <Link to="ch06" class="chapter-card">
      <div class="chapter-num">Ch 6</div>
      <div>錯誤處理</div>
      <div class="chapter-subtitle">Error Handling</div>
    </Link>
    <Link to="ch07" class="chapter-card">
      <div class="chapter-num">Ch 7</div>
      <div>介面</div>
      <div class="chapter-subtitle">Interfaces</div>
    </Link>
    <Link to="ch08" class="chapter-card">
      <div class="chapter-num">Ch 8</div>
      <div>套件</div>
      <div class="chapter-subtitle">Packages &amp; Modules</div>
    </Link>
    <Link to="ch09" class="chapter-card">
      <div class="chapter-num">Ch 9</div>
      <div>程式除錯</div>
      <div class="chapter-subtitle">fmt, Logging &amp; Testing</div>
    </Link>
    <Link to="ch10" class="chapter-card">
      <div class="chapter-num">Ch 10</div>
      <div>時間處理</div>
      <div class="chapter-subtitle">Time &amp; Duration</div>
    </Link>
    <Link to="ch11" class="chapter-card">
      <div class="chapter-num">Ch 11</div>
      <div>編碼／解碼 JSON</div>
      <div class="chapter-subtitle">JSON &amp; gob</div>
    </Link>
    <Link to="ch12" class="chapter-card">
      <div class="chapter-num">Ch 12</div>
      <div>系統與檔案</div>
      <div class="chapter-subtitle">Flags, Signals &amp; Files</div>
    </Link>
    <Link to="ch13" class="chapter-card">
      <div class="chapter-num">Ch 13</div>
      <div>SQL 與資料庫</div>
      <div class="chapter-subtitle">database/sql &amp; MySQL</div>
    </Link>
    <Link to="ch14" class="chapter-card">
      <div class="chapter-num">Ch 14</div>
      <div>HTTP 客戶端</div>
      <div class="chapter-subtitle">net/http Client</div>
    </Link>
    <Link to="ch15" class="chapter-card">
      <div class="chapter-num">Ch 15</div>
      <div>HTTP 伺服器</div>
      <div class="chapter-subtitle">net/http Server</div>
    </Link>
    <Link to="ch16" class="chapter-card">
      <div class="chapter-num">Ch 16</div>
      <div>並行性運算</div>
      <div class="chapter-subtitle">Goroutines &amp; Channels</div>
    </Link>
    <Link to="ch17" class="chapter-card">
      <div class="chapter-num">Ch 17</div>
      <div>Go 語言工具</div>
      <div class="chapter-subtitle">The Go Toolchain</div>
    </Link>
    <Link to="ch18" class="chapter-card">
      <div class="chapter-num">Ch 18</div>
      <div>加密安全</div>
      <div class="chapter-subtitle">Cryptography &amp; TLS</div>
    </Link>
    <Link to="ch19" class="chapter-card">
      <div class="chapter-num">Ch 19</div>
      <div>reflect 與 unsafe</div>
      <div class="chapter-subtitle">Reflection &amp; unsafe</div>
    </Link>
  </div>
</div>

<!--
大家好，歡迎來到 Go 實戰開發課程！

這一頁是整門課的目錄，一共 20 個章節：第 0 章先把開發環境裝好，第 1 到第 8 章打好語言基礎（變數、流程控制、型別、函式、錯誤處理、介面、套件），第 9 到第 15 章進入實戰（除錯與測試、時間、JSON、檔案、資料庫、HTTP 客戶端與伺服器），最後第 16 到第 19 章是 Go 最有特色的並行性運算、工具鏈、加密安全，以及 reflect／unsafe 這兩個特殊套件。

點任何一張卡片就可以跳到該章節，每一章的封面也都有「返回目錄」的連結。
-->

---
src: ./ch00-setup.md
---

---
src: ./ch01-variables-operators.md
---

---
src: ./ch02-flow-control.md
---

---
src: ./ch03-core-types.md
---

---
src: ./ch04-composite-types.md
---

---
src: ./ch05-functions.md
---

---
src: ./ch06-error-handling.md
---

---
src: ./ch07-interfaces.md
---

---
src: ./ch08-packages.md
---

---
src: ./ch09-debug-log-test.md
---

---
src: ./ch10-time.md
---

---
src: ./ch11-json.md
---

---
src: ./ch12-system-files.md
---

---
src: ./ch13-sql-database.md
---

---
src: ./ch14-http-client.md
---

---
src: ./ch15-http-server.md
---

---
src: ./ch16-concurrency.md
---

---
src: ./ch17-go-tools.md
---

---
src: ./ch18-crypto.md
---

---
src: ./ch19-reflect-unsafe.md
---

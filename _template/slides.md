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
title: 投影片主標題
routeAlias: chXX
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">投影片主標題</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「一句話描述」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第 N 章！（講稿：先回顧上一章，再說明這章要解決什麼問題）
-->

---
layout: default
---

# Outline

- **主題一** — 一句話說明
- **主題二** — 一句話說明
- **章節總結**

---

# 回顧：上一章的重點

- 重點一
- 重點二

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 主題一
## English Subtitle

---

# 什麼是 XXX？

| 項目 | 說明 |
| --- | --- |
| **用途** | 說明 |

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>補充：</b> 說明文字
</div>

<!--
講稿：先情境、再定義；程式碼前說目的，程式碼後說執行結果。
-->

---
layout: default
---

# 練習 1：題目名稱
### 任務說明

題目描述...

---

# 練習 1：解題提示
### 提示說明

1. 步驟一
2. 步驟二

---

# 章節總結

- 重點一
- 重點二

下一章我們會介紹...

---
layout: end
---

# Q & A

有任何問題嗎？

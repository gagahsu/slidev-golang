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
title: 加密安全
routeAlias: ch18
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
  <h1 style="color: #1a5c5c; font-size: 3.4rem; font-weight: 900; line-height: 1.15; margin-bottom: 1.5rem;">加密安全</h1>
  <div style="height: 4px; width: 320px; background: linear-gradient(90deg, #5eada0, #a7d9d0); border-radius: 2px; margin-bottom: 1.5rem;"></div>
  <p style="color: #4a7c7c; font-size: 1.15rem; font-style: italic;">「不要自己發明密碼學，正確地使用標準函式庫」</p>
  <Link to="home" style="color: #9dc4c4; font-size: 0.85rem; margin-top: 2rem; text-decoration: none; letter-spacing: 0.05em;">← 返回目錄</Link>
</div>

<!--
大家好，歡迎來到第十八章！

前面幾章我們寫了資料庫、HTTP 伺服器、API，這些程式都會處理使用者的資料：密碼、個人資料、付款資訊。如果這些資料被偷走或竄改，後果非常嚴重。今天要學的就是保護資料的工具：雜湊、加密、數位簽章，以及 HTTPS。

封面這句話是資安界的鐵律：不要自己發明密碼學演算法。即使是專家設計的演算法，也要經過多年的公開檢驗才敢使用。我們要做的，是正確地使用 Go 標準函式庫提供的工具。Go 的 crypto 套件由專門的密碼學團隊維護，是業界公認品質最好的實作之一。
-->

---
layout: default
---

# Outline

- **前言** — 資訊安全的三個目標、密碼學工具的全貌
- **雜湊函式** — MD5、SHA-2、SHA-3、HMAC、密碼儲存
- **對稱式加密** — AES-GCM
- **非對稱式加密** — RSA-OAEP
- **數位簽章** — Ed25519
- **HTTPS / TLS 與 X.509 憑證** — 自簽署憑證、ECDSA
- **GoShop 專案實作** — 第 18 步：後台登入與 HTTPS
- **章節總結**

<!--
今天的內容依照工具的類型排列。

雜湊用來確認資料有沒有被改過；對稱式加密用同一把金鑰加密和解密；非對稱式加密用一對公鑰和私鑰；數位簽章用來證明「這份資料確實是我發出的」；最後把這些工具組合起來，就是保護網站連線的 HTTPS。
-->

---

# 回顧：Go 語言工具

- `go build -ldflags "-X main.version=…"` 注入版本；`GOOS` / `GOARCH` 跨平台編譯
- `go vet` 靜態分析；**`govulncheck`** 檢查程式實際呼叫到的已知漏洞 ← 資安的第一步
- 第 8 章：`go.sum` 記錄每個模組的**雜湊值**，防止被竄改 ← 今天學雜湊的原理
- 第 13 章：**SQL Injection**；第 15 章：**XSS**、限制請求大小

<!--
回顧一下上一章。

我們學了 Go 的工具鏈，其中 govulncheck 是資安的第一步：確保我們用的套件沒有已知的漏洞。

其實這門課一路上已經碰過很多資安觀念：第 8 章的 go.sum 用雜湊值防止模組被竄改、第 13 章的 SQL Injection、第 15 章的 XSS。今天要學的是這些防護背後的密碼學工具。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 前言
## Cryptography Overview

<!--
先來看密碼學工具的全貌。
-->

---

# 資訊安全的三個目標與對應的工具

| 目標 | 意思 | 工具 | Go 套件 |
| --- | --- | --- | --- |
| **機密性** | 只有該看的人看得到 | 對稱式 / 非對稱式加密 | `crypto/aes`、`crypto/rsa` |
| **完整性** | 資料沒有被竄改 | 雜湊、HMAC | `crypto/sha256`、`crypto/hmac` |
| **真實性** | 確認資料來源 | 數位簽章、憑證 | `crypto/ed25519`、`crypto/x509` |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
🎲 <b>所有的金鑰、nonce、salt 都要用 <code>crypto/rand</code> 產生</b>，絕對不要用 <code>math/rand</code>（可預測）。Go 1.24 起 <code>crypto/rand.Read</code> 保證不會失敗。
</div>

<!--
資訊安全有三個基本目標，簡稱 CIA。

機密性：資料只有該看的人看得到，靠加密達成。完整性：資料在傳輸或儲存的過程中沒有被竄改，靠雜湊達成。真實性：確認資料真的是某個人發出的，靠數位簽章達成。

另外一個非常重要的提醒：所有跟安全有關的隨機值，例如金鑰、nonce、salt，一定要用 crypto/rand 產生。第 8 章提過的 math/rand 是給遊戲、模擬用的，它產生的數字是可以被預測的，用在安全的地方就等於把鑰匙交給攻擊者。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 雜湊函式
## MD5・SHA-2・SHA-3

<!--
先從雜湊函式開始。
-->

---
zoom: 0.91
---

# 什麼是雜湊函式？

「**雜湊函式把任意長度的資料，轉換成固定長度的『指紋』**，而且無法從指紋還原原本的資料。」

| 特性 | 說明 |
| --- | --- |
| **固定長度** | 不管輸入多大，SHA-256 的輸出永遠是 32 bytes（64 個十六進位字元） |
| **單向** | 無法從雜湊值反推原始資料 |
| **雪崩效應** | 輸入只改一個字，輸出就完全不同 |
| **抗碰撞** | 很難找到兩份不同的資料，產生相同的雜湊值 |

| 演算法 | 輸出長度 | 狀態 |
| --- | --- | --- |
| MD5 | 128 bits | ❌ **已被攻破**，只能用於非安全用途（例如快取的鍵） |
| SHA-256 / SHA-512（SHA-2） | 256 / 512 bits | ✅ 目前的主流 |
| SHA3-256（SHA-3） | 256 bits | ✅ 新一代標準（Go 1.24 起在標準函式庫） |

<!--
什麼是雜湊函式？

它就像人的指紋：每個人的指紋都不一樣，看到指紋就能確認是誰，但你沒辦法從指紋「還原」出一個人。雜湊函式把任何資料，不管是一個字還是一部電影，都轉換成一段固定長度的指紋。

雜湊有四個重要特性：固定長度、單向、雪崩效應、抗碰撞。雪崩效應的意思是，輸入只要改一點點，輸出就完全不一樣，所以很適合用來檢查資料有沒有被改過。第 8 章的 go.sum 就是用雜湊值確認下載的模組沒有被竄改。

MD5 已經被攻破了，可以刻意製造出雜湊值相同的兩份資料，所以不能再用在安全的用途。目前的主流是 SHA-2 家族，SHA-3 是新一代的標準。
-->

---
zoom: 0.84
---

# 使用 MD5、SHA-2、SHA-3

```go
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
)

func main() {
	data := []byte("hello")
	m := md5.Sum(data)        // [16]byte
	s2 := sha256.Sum256(data) // [32]byte
	s3 := sha3.Sum256(data)   // [32]byte，Go 1.24+
	fmt.Println("MD5    ", hex.EncodeToString(m[:]))
	fmt.Println("SHA-256", hex.EncodeToString(s2[:]))
	fmt.Println("SHA3   ", hex.EncodeToString(s3[:]))

	h := sha256.Sum256([]byte("hellO")) // 只改一個字母
	// 雪崩效應：結果完全不同
	fmt.Println("SHA-256", hex.EncodeToString(h[:]))
}
```

```text
MD5     5d41402abc4b2a76b9719d911017c592
SHA-256 2cf24dba5fb0a30e26e83b2ac5b9e29e
        1b161e5c1fa7425e73043362938b9824
```

<!--
Go 的雜湊函式用法非常一致：套件名稱加上 Sum，傳入 []byte，回傳一個固定長度的位元組陣列。

回傳的是陣列，第 4 章學過陣列跟切片不一樣，要轉成切片才能傳給 hex.EncodeToString，所以寫 m[:]。hex.EncodeToString 把位元組轉成十六進位的字串，這是雜湊值最常見的表示方式，也可以直接用 Printf 的 %x。

最後一行把 hello 的最後一個字母改成大寫，SHA-256 的結果就完全不一樣了，這就是雪崩效應。
-->

---
zoom: 0.83
---

# 計算檔案的雜湊值：hash.Hash 介面

大檔案不要一次讀進記憶體，用 **`io.Copy`** 串流計算

```go
package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New() // hash.Hash 實作了 io.Writer
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func main() {
	os.WriteFile("demo.txt", []byte("hello"), 0o644)
	fmt.Println(fileSHA256("demo.txt")) // 2cf24dba…9824 <nil>
}
```

<!--
如果要計算一個好幾 GB 的檔案的雜湊值，不能用 os.ReadFile 一次讀進記憶體。

sha256.New() 回傳一個 hash.Hash，它實作了 io.Writer 介面，所以可以用 io.Copy 把檔案「倒進去」，一邊讀一邊計算。最後呼叫 Sum(nil) 取得結果。這又是第 7 章 io.Reader、io.Writer 組合的威力。

實務上，軟體下載頁面常常會附上 SHA-256 值，下載完用這個方法計算一次，比對一樣，就代表檔案完整、沒有被竄改。
-->

---
zoom: 0.79
---

# HMAC：加上金鑰的雜湊

一般雜湊**任何人都能算**；HMAC 需要**金鑰**，可以確認「**資料沒被改、而且來自持有金鑰的人**」

```go
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var secret = []byte("server-secret-key") // 只有伺服器知道

func sign(msg string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func verify(msg, sig string) bool {
	expected, _ := hex.DecodeString(sign(msg))
	given, _ := hex.DecodeString(sig)
	return hmac.Equal(expected, given) // ⚠️ 用 hmac.Equal，不要用 ==
}

func main() {
	sig := sign("user=alice&role=member")
	fmt.Println(verify("user=alice&role=member", sig)) // true
	fmt.Println(verify("user=alice&role=admin", sig))  // false：被竄改了
}
```

<!--
一般的雜湊有一個問題：任何人都能計算。攻擊者可以修改資料，再重新計算雜湊值，接收者就分辨不出來了。

HMAC 是「加上金鑰的雜湊」，只有持有金鑰的人才能算出正確的值。伺服器發給使用者一段資料和它的 HMAC，使用者就算把 role=member 改成 role=admin，也算不出對應的 HMAC，伺服器一驗證就會發現被竄改了。

網站的 Cookie 簽章、JWT（HS256）、Webhook 的驗證，都是用 HMAC。

使用 HMAC 的注意事項：比對的時候一定要用 hmac.Equal，不要用 ==。hmac.Equal 是「固定時間比較」，不管在第幾個位元組不同，花的時間都一樣；用 == 的話，攻擊者可以透過量測時間差，一個位元組一個位元組地猜出正確的值，這叫做時序攻擊。
-->

---

# 使用雜湊的注意事項：儲存密碼

**不要**直接用 SHA-256 存密碼（太快，容易被暴力破解）；使用**專門的密碼雜湊函式**

| 方式 | 評價 |
| --- | --- |
| 明文儲存 | ❌ 資料庫外洩 = 所有密碼外洩 |
| `sha256(password)` | ❌ 速度太快，每秒可試數十億次；相同密碼有相同雜湊 |
| **bcrypt** / **Argon2id**（`golang.org/x/crypto`） | ✅ 刻意設計得很慢、自動加鹽，業界主流 |
| **PBKDF2**（`crypto/pbkdf2`，Go 1.24+） | ✅ 標準函式庫內建，迭代次數建議 ≥ 600,000 |

```go
import "golang.org/x/crypto/bcrypt"

// 存進資料庫
hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
// 登入時比對，nil 代表正確
err := bcrypt.CompareHashAndPassword(hash, []byte(input))
```

<!--
使用雜湊最重要的注意事項：儲存密碼。

很多人以為把密碼做 SHA-256 雜湊再存起來就安全了，其實不然。SHA-256 的設計目標是「快」，現代的顯示卡每秒可以計算幾十億次，攻擊者拿到雜湊值之後，把常見的密碼一個一個算過去，很快就能找出原本的密碼。而且相同的密碼會得到相同的雜湊值，一眼就能看出哪些人用了一樣的密碼。

所以密碼要用「專門的密碼雜湊函式」：bcrypt、Argon2id、PBKDF2。它們刻意設計得很慢，而且會自動加入隨機的「鹽」（salt），讓相同的密碼得到不同的結果。

bcrypt 是業界最常見的選擇，在 golang.org/x/crypto 套件裡，用法只有兩個函式：註冊時 GenerateFromPassword，登入時 CompareHashAndPassword。
-->

---
layout: default
---

# 練習 1：檔案完整性檢查
### 任務說明

寫一個命令列工具 `checksum`：

1. 旗標 `-algo` 可選 `md5`、`sha256`、`sha3`（預設 `sha256`）
2. 對 `flag.Args()` 裡的每個檔案，用**串流**方式計算雜湊值，印出 `雜湊值  檔名`（跟 `sha256sum` 指令的格式一樣）
3. 旗標 `-verify 雜湊值` 存在時，比對第一個檔案的雜湊值，印出「✓ 相符」或「✗ 不符」
4. 提示：三種演算法都有 `New()` 函式，回傳 `hash.Hash`，可以用 `map[string]func() hash.Hash` 選擇

<!--
這個練習結合了今天的雜湊、第 12 章的命令列旗標，以及第 5 章的函式型別。

重點在第 4 步：md5.New、sha256.New、sha3.New256 的回傳值都是 hash.Hash 介面，所以可以放進同一個 map 裡，用演算法名稱查出對應的函式。這就是第 7 章介面的威力。
-->

---

# 練習 1：解題提示
### 提示說明

```go
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha3"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
)

var algos = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha256": sha256.New,
	"sha3":   func() hash.Hash { return sha3.New256() },
}
```

<!--
algos 是一個 map，把演算法名稱對應到建立 hash.Hash 的函式。md5.New 和 sha256.New 的型別剛好就是 func() hash.Hash，可以直接放進去；sha3.New256 回傳的是 *sha3.SHA3，所以用一個匿名函式包裝一下。
-->

---

# 練習 1：解題提示（續）
### 提示說明

```go
// 續上頁
func sum(newHash func() hash.Hash, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := newHash()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
```

<!--
sum 函式接收「建立雜湊的函式」和檔案路徑，開啟檔案、用 io.Copy 串流計算，最後用 %x 格式化成十六進位字串。
-->

---
zoom: 0.84
---

# 練習 1：解題提示（續 2）
### 提示說明

```go
// 續上頁
func main() {
	algo := flag.String("algo", "sha256", "md5 / sha256 / sha3")
	want := flag.String("verify", "", "要比對的雜湊值")
	flag.Parse()
	newHash, ok := algos[*algo]
	if !ok {
		fmt.Println("不支援的演算法：", *algo)
		os.Exit(2)
	}
	for i, path := range flag.Args() {
		got, err := sum(newHash, path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%s  %s\n", got, path)
		if i == 0 && *want != "" {
			if got == *want {
				fmt.Println("✓ 相符")
			} else {
				fmt.Println("✗ 不符")
			}
		}
	}
}
```

<!--
main 解析旗標，用 comma ok 從 map 查出演算法，找不到就結束。接著走訪每個檔案計算雜湊值，格式跟 Linux 的 sha256sum 指令一樣：雜湊值、兩個空白、檔名。

有 -verify 的時候，比對第一個檔案的結果。這裡比對的是公開的檔案雜湊值，不是秘密，所以用 == 就可以；如果比對的是 HMAC 這種跟秘密有關的值，就要用 hmac.Equal。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 加密法
## Symmetric & Asymmetric Encryption

<!--
接下來是加密：讓資料只有該看的人看得到。
-->

---

# 對稱式 vs. 非對稱式加密

| 比較 | 對稱式加密 | 非對稱式加密 |
| --- | --- | --- |
| **金鑰** | 加密、解密用**同一把**金鑰 | **公鑰**加密、**私鑰**解密 |
| **比喻** | 家裡的門鎖：鎖門開門同一把鑰匙 | 郵筒：誰都能投信，只有郵差能打開 |
| **速度** | 非常快 | 很慢（約慢上千倍） |
| **難題** | 怎麼安全地把金鑰交給對方？ | 資料量大時太慢 |
| **代表演算法** | **AES-GCM**、ChaCha20-Poly1305 | **RSA-OAEP**、ECDH |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>實務上兩者搭配使用：</b> 用非對稱式加密安全地交換一把「對稱式金鑰」，之後的大量資料都用對稱式加密。HTTPS 就是這樣運作的。
</div>

<!--
加密分成兩大類。

對稱式加密：加密和解密用同一把金鑰，就像家裡的門鎖。它非常快，但有一個難題：我們要怎麼把金鑰安全地交給對方？如果在網路上直接傳送金鑰，被攔截就完了。

非對稱式加密：有一對金鑰，公鑰公開給所有人，私鑰自己保管。用公鑰加密的資料，只有私鑰才能解開。就像郵筒：任何人都能把信投進去，但只有拿著鑰匙的郵差能打開。它解決了金鑰交換的問題，但速度很慢。

所以實務上兩者搭配：先用非對稱式加密交換一把對稱式的金鑰，之後的資料都用對稱式加密。今天最後講的 HTTPS 就是這樣運作的。
-->

---
zoom: 0.79
---

# 對稱式加密法：AES + GCM

**AES** 是加密演算法，**GCM** 模式同時提供**加密**與**防竄改**（authenticated encryption）

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func main() {
	key := make([]byte, 32) // 32 bytes = AES-256
	rand.Read(key)          // 金鑰一定要用 crypto/rand 產生

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// Go 1.24+：自動處理 nonce
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		panic(err)
	}
	ciphertext := gcm.Seal(nil, nil, []byte("信用卡：4111-1111"), nil)
	fmt.Printf("密文 %d bytes：%x…\n", len(ciphertext), ciphertext[:8])

	plain, err := gcm.Open(nil, nil, ciphertext, nil)
	fmt.Println(string(plain), err) // 信用卡：4111-1111 <nil>
}
```

<!--
AES 是目前最主流的對稱式加密演算法，全世界的銀行、政府、網站都在用。金鑰長度可以是 16、24、32 位元組，分別是 AES-128、AES-192、AES-256。

AES 本身一次只能加密 16 個位元組，所以要搭配一種「模式」來加密任意長度的資料。GCM 是目前最推薦的模式，它除了加密，還會附上一個驗證標籤：密文只要被改動一個位元，解密時就會失敗，這叫做「認證加密」。

GCM 需要一個「nonce」，每次加密都必須不一樣。Go 1.24 加入的 NewGCMWithRandomNonce 會自動產生隨機的 nonce 並附在密文前面，解密時自動取出，我們完全不用自己處理，這是目前最推薦、最不容易出錯的寫法。

Seal 加密、Open 解密。密文比明文多了 28 個位元組：12 個是 nonce，16 個是驗證標籤。
-->

---
zoom: 0.87
---

# 使用 AES-GCM 的注意事項

**注意事項之一：** 同一把金鑰，**nonce 絕對不能重複使用**，否則加密會被破解（`NewGCMWithRandomNonce` 已自動處理）

**注意事項之二：** 解密失敗一律視為**資料被竄改**，不要使用任何部分結果

**注意事項之三：** 金鑰不能寫死在程式碼中，要從環境變數或金鑰管理服務（KMS）取得

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

func main() {
	key := make([]byte, 32)
	rand.Read(key)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCMWithRandomNonce(block)
	ct := gcm.Seal(nil, nil, []byte("轉帳 100 元"), nil)

	ct[len(ct)-1] ^= 0x01 // 攻擊者竄改了最後一個位元
	_, err := gcm.Open(nil, nil, ct, nil)
	fmt.Println(err) // cipher: message authentication failed
}
```

<!--
使用 AES-GCM 有三個注意事項。

第一，nonce 絕對不能重複。如果同一把金鑰用同一個 nonce 加密了兩份資料，攻擊者就能推算出內容。以前的寫法要自己用 crypto/rand 產生 nonce，很容易出錯；用 NewGCMWithRandomNonce 就不用擔心了。

第二，GCM 會驗證資料的完整性。這個範例模擬攻擊者竄改了密文的一個位元，解密時就會回傳 message authentication failed 的錯誤。遇到這個錯誤，一律當作資料被竄改，不能使用。

第三，金鑰的保管。加密再強，金鑰被偷就沒用了。金鑰不能寫在程式碼裡，實務上從環境變數讀取，大型系統會用雲端的金鑰管理服務（KMS）。
-->

---
zoom: 0.81
---

# 非對稱式加密法：RSA-OAEP

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func main() {
	// 私鑰（至少 2048 bits）
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	pub := &priv.PublicKey // 公鑰可以公開給任何人

	msg := []byte("AES 金鑰：a1b2c3…")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader,
		pub, msg, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("密文長度：", len(ciphertext)) // 256：等於金鑰長度

	plain, err := rsa.DecryptOAEP(sha256.New(), nil, priv,
		ciphertext, nil)
	fmt.Println(string(plain), err) // 只有私鑰能解密
}
```

<!--
RSA 是最經典的非對稱式加密演算法。

rsa.GenerateKey 產生一對金鑰，長度至少要 2048 位元，Go 1.24 以後小於 1024 位元的金鑰會直接被拒絕。私鑰裡包含了公鑰，用 &priv.PublicKey 取出來。

加密要用 OAEP 模式，這是目前推薦的填充方式；舊的 PKCS#1 v1.5 加密有已知的攻擊方式，不要使用。EncryptOAEP 用公鑰加密，DecryptOAEP 用私鑰解密。

注意 RSA 能加密的資料量很小，2048 位元的金鑰大約只能加密 190 個位元組，所以它通常只用來加密一把 AES 金鑰，而不是加密整份資料。這就是剛剛說的「兩者搭配使用」。
-->

---
zoom: 0.81
---

# 補充：把金鑰存成 PEM 檔案

金鑰通常以 **PEM**（Base64 + 標頭）格式存成文字檔，方便交換與保存

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKCS8PrivateKey(priv) // 轉成標準的 DER 位元組
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: der}
	// ⚠️ 私鑰只有自己能讀
	os.WriteFile("key.pem", pem.EncodeToMemory(block), 0o600)

	pubDER, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}
	pubPEM := pem.EncodeToMemory(pubBlock)
	os.WriteFile("pub.pem", pubPEM, 0o644) // 公鑰可以公開
}
```

```text
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7…
-----END PRIVATE KEY-----
```

<!--
這是補充內容：金鑰怎麼存成檔案。

金鑰在記憶體裡是 Go 的結構，要存成檔案，先用 x509 套件轉成標準的 DER 二進位格式，再用 pem 套件轉成文字格式。PEM 就是我們常看到的 BEGIN PRIVATE KEY、END PRIVATE KEY 包起來的那種文字檔，裡面是 Base64 編碼。

注意檔案權限：私鑰檔案要用 0o600，只有自己能讀寫，這是第 12 章學的檔案權限。公鑰可以公開，用 0o644。

讀取的時候反過來：pem.Decode 解出 DER，再用 x509.ParsePKCS8PrivateKey 解析成金鑰。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# 數位簽章
## Ed25519

<!--
接下來是數位簽章：證明「這份資料確實是我發出的，而且沒有被改過」。
-->

---
zoom: 0.93
---

# 數位簽章：Ed25519

**私鑰簽章、公鑰驗證**（跟加密相反）— 證明**資料來源**且**沒被竄改**

```go
package main

import (
	"crypto/ed25519"
	"fmt"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(nil) // nil：使用 crypto/rand
	if err != nil {
		panic(err)
	}
	msg := []byte("版本 1.2.3 的 SHA-256：2cf24dba…")

	sig := ed25519.Sign(priv, msg) // 只有私鑰能產生簽章
	fmt.Println("簽章長度：", len(sig)) // 64

	// true：任何人都能用公鑰驗證
	fmt.Println(ed25519.Verify(pub, msg, sig))
	msg[0] = 'X'
	fmt.Println(ed25519.Verify(pub, msg, sig)) // false：內容被改了
}
```

<!--
數位簽章跟非對稱式加密剛好相反：用私鑰「簽章」，用公鑰「驗證」。

就像現實中的簽名：只有本人能簽出自己的簽名（私鑰），但任何人都能拿著樣本比對簽名是不是真的（公鑰）。

Ed25519 是現代最推薦的簽章演算法：速度快、金鑰短（公鑰只有 32 位元組）、而且 API 設計得很難用錯。GenerateKey 傳入 nil 會自動使用 crypto/rand。Sign 產生 64 位元組的簽章，Verify 回傳布林值。

資料被改了一個字，驗證就失敗。這就是軟體更新的安全機制：開發者用私鑰替新版本簽章，使用者的程式用內建的公鑰驗證，確保下載的更新真的來自官方、沒有被植入惡意程式。SSH 金鑰、Git commit 簽章也都常用 Ed25519。
-->

---

# 雜湊、HMAC、數位簽章的比較

| 工具 | 需要金鑰 | 誰能產生 | 誰能驗證 | 用途 |
| --- | --- | --- | --- | --- |
| **雜湊**（SHA-256） | 不需要 | 任何人 | 任何人 | 檢查檔案完整性、go.sum |
| **HMAC** | 一把共享金鑰 | 持有金鑰者 | 持有金鑰者 | Cookie、Webhook、JWT（HS256） |
| **數位簽章**（Ed25519） | 私鑰 / 公鑰 | **只有私鑰持有者** | **任何有公鑰的人** | 軟體更新、憑證、JWT（EdDSA） |

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
💡 <b>怎麼選？</b> 發出者和驗證者是<b>同一方</b>（例如伺服器驗證自己發的 Cookie）→ HMAC；驗證者是<b>第三方</b>或<b>很多人</b> → 數位簽章。
</div>

<!--
這張表比較了三種「保護完整性」的工具。

雜湊不需要金鑰，任何人都能算，所以只能防止「意外」的損壞，擋不住刻意的竄改。HMAC 需要一把共享的金鑰，適合同一方自己產生、自己驗證的情況。數位簽章用一對金鑰，只有私鑰持有者能簽，但任何人都能驗證，適合要讓很多人驗證的情況。

選擇的原則寫在下方：看驗證的人是誰。
-->

---
layout: default
---

# 練習 2：加密的訊息信封
### 任務說明

Alice 要傳一段**很長**的機密訊息給 Bob，並證明是自己發的：

1. Bob 產生 RSA 金鑰對；Alice 產生 Ed25519 金鑰對（公鑰都公開）
2. Alice：產生隨機的 AES-256 金鑰 → 用 **AES-GCM** 加密訊息 → 用 **Bob 的 RSA 公鑰**加密 AES 金鑰 → 用**自己的 Ed25519 私鑰**對「加密後的訊息」簽章
3. 把三樣東西放進 `Envelope{EncKey, Ciphertext, Signature []byte}`
4. Bob：先用 **Alice 的公鑰**驗證簽章 → 用**自己的 RSA 私鑰**解出 AES 金鑰 → 用 AES-GCM 解密訊息
5. 竄改 `Ciphertext` 的一個位元組，確認 Bob 會拒絕

<!--
這個練習把今天學的三種工具組合起來，這其實就是很多加密通訊軟體的基本架構：用對稱式加密處理大量資料，用非對稱式加密保護對稱金鑰，用數位簽章確認身分。

這種「用 RSA 保護 AES 金鑰」的做法叫做「信封加密」（envelope encryption），雲端服務的金鑰管理也是這樣運作的。
-->

---
zoom: 0.84
---

# 練習 2：解題提示
### 提示說明

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
)

type Envelope struct{ EncKey, Ciphertext, Signature []byte }

func newGCM(key []byte) cipher.AEAD {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		panic(err)
	}
	return gcm
}
```

<!--
先定義信封結構，三個欄位都是 []byte。

newGCM 是一個輔助函式，把 AES 金鑰轉換成 GCM 物件。金鑰長度是我們自己控制的，不會出錯，所以錯誤直接 panic。
-->

---
zoom: 0.84
---

# 練習 2：解題提示（續）
### 提示說明

```go
// 續上頁
func seal(msg []byte, bobPub *rsa.PublicKey,
	alicePriv ed25519.PrivateKey) Envelope {
	aesKey := make([]byte, 32)
	rand.Read(aesKey)
	ct := newGCM(aesKey).Seal(nil, nil, msg, nil)
	encKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader,
		bobPub, aesKey, nil)
	if err != nil {
		panic(err)
	}
	return Envelope{encKey, ct, ed25519.Sign(alicePriv, ct)}
}

func open(e Envelope, bobPriv *rsa.PrivateKey,
	alicePub ed25519.PublicKey) ([]byte, error) {
	if !ed25519.Verify(alicePub, e.Ciphertext, e.Signature) {
		return nil, errors.New("簽章驗證失敗：不是 Alice 發的或已被竄改")
	}
	aesKey, err := rsa.DecryptOAEP(sha256.New(), nil, bobPriv,
		e.EncKey, nil)
	if err != nil {
		return nil, err
	}
	return newGCM(aesKey).Open(nil, nil, e.Ciphertext, nil)
}
```

<!--
seal 是 Alice 做的事：產生隨機的 AES 金鑰、用它加密訊息、用 Bob 的公鑰加密 AES 金鑰、再用自己的私鑰對密文簽章。

open 是 Bob 做的事，順序剛好相反：先驗證簽章，確認是 Alice 發的、而且沒被竄改，驗證失敗就立刻拒絕；接著用自己的私鑰解出 AES 金鑰，最後解密訊息。

先驗證簽章再解密，是一個好的習慣：不信任的資料，在確認來源之前不要做任何處理。
-->

---

# 練習 2：解題提示（續 2）
### 提示說明

```go
// 續上頁
func main() {
	bobPriv, _ := rsa.GenerateKey(rand.Reader, 2048)
	alicePub, alicePriv, _ := ed25519.GenerateKey(nil)

	msg := []byte("下午三點，老地方見。")
	env := seal(msg, &bobPriv.PublicKey, alicePriv)
	got, err := open(env, bobPriv, alicePub)
	fmt.Println(string(got), err) // 下午三點，老地方見。 <nil>

	env.Ciphertext[0] ^= 0xFF // 竄改
	_, err = open(env, bobPriv, alicePub)
	fmt.Println(err) // 簽章驗證失敗：不是 Alice 發的或已被竄改
}
```

<!--
main 產生兩組金鑰，Alice 封裝訊息，Bob 打開，成功讀到內容。

接著竄改密文的第一個位元組，Bob 在驗證簽章的階段就發現了，直接拒絕。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# HTTPS/TLS 與 X.509 憑證
## TLS & Certificates

<!--
最後，我們把今天學的所有工具組合起來，看看 HTTPS 是怎麼運作的。
-->

---

# HTTPS 是怎麼運作的？

```text
瀏覽器 / Go 客戶端                                   伺服器
  ── ① 我想連線，我支援 TLS 1.3 ─────────────────────▶
  ◀──────── ② 這是我的憑證（公鑰 + CA 的數位簽章）──
  ③ 用內建的 CA 公鑰驗證憑證的簽章 → 確認對方真的是 goshop.com
  ── ④ 用金鑰交換（ECDH）協商出一把對稱式金鑰 ──────▶
  ◀═══════ ⑤ 之後的資料都用 AES-GCM 等對稱式加密傳輸 ═══════▶
```

| 元件 | 用到今天學的 |
| --- | --- |
| **X.509 憑證** | 伺服器的公鑰 + 憑證機構（CA）的**數位簽章** |
| **金鑰交換** | **非對稱式**密碼學，安全地協商出對稱式金鑰 |
| **資料傳輸** | **對稱式加密**（AES-GCM）+ 完整性保護 |

<!--
HTTPS 就是 HTTP 加上 TLS 加密，它把今天學的所有工具都組合在一起。

連線的時候，伺服器會先出示一張「憑證」。憑證裡有伺服器的公鑰，還有憑證機構（CA）的數位簽章，證明這個公鑰確實屬於 goshop.com。瀏覽器和作業系統內建了可信任的 CA 公鑰，用它驗證憑證的簽章，就能確認對方不是冒充的。

確認身分之後，雙方用非對稱式的金鑰交換演算法，協商出一把只有彼此知道的對稱式金鑰，之後所有的資料都用 AES-GCM 這類對稱式加密傳輸。

這就是剛剛說的「非對稱式加密交換金鑰，對稱式加密傳輸資料」。
-->

---

# 給伺服器使用自簽署憑證：產生憑證

開發與測試時，可以產生**自簽署憑證**（自己當 CA）；正式環境請使用 Let's Encrypt 等 CA 簽發的憑證

```bash
# Go 內建的產生工具（預設 RSA）
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go \
  --host localhost,127.0.0.1

# ★ 改用 ECDSA 簽章演算法（P-256 曲線，金鑰更短、速度更快）
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go \
  --host localhost,127.0.0.1 --ecdsa-curve P256
# wrote cert.pem
# wrote key.pem          ← 權限自動設為 0600

openssl x509 -in cert.pem -noout -text | grep Algorithm
#   Signature Algorithm: ecdsa-with-SHA256
#   Public Key Algorithm: id-ecPublicKey
```

<!--
正式的網站，憑證要由 CA 簽發，現在最主流的是免費的 Let's Encrypt。但開發和測試的時候，可以自己產生一張「自簽署憑證」：自己當 CA，替自己簽名。

Go 的原始碼裡附了一個產生憑證的小工具 generate_cert.go，用 go run 直接執行就好。--host 指定這張憑證適用的主機名稱和 IP。

它預設使用 RSA 金鑰。加上 --ecdsa-curve P256，就會改用 ECDSA 簽章演算法。ECDSA 是基於橢圓曲線的演算法，256 位元的 ECDSA 金鑰，安全性相當於 3072 位元的 RSA，但金鑰更短、速度更快，現在主流的網站憑證大多已經改用 ECDSA。

執行後會產生 cert.pem 和 key.pem 兩個檔案，私鑰的權限自動設為 0600。
-->

---
zoom: 0.84
---

# 用 Go 程式產生 ECDSA 自簽署憑證

```go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// ★ ECDSA P-256
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"GoShop Dev"}},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0), // 有效期一年
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
```

<!--
也可以用 Go 程式自己產生憑證，這樣更能理解憑證裡有什麼。

ecdsa.GenerateKey 產生一把 P-256 曲線的 ECDSA 金鑰。x509.Certificate 是憑證的範本：序號、組織名稱、適用的網域名稱和 IP、有效期間、用途。
-->

---

# 用 Go 程式產生 ECDSA 自簽署憑證（續）

```go
	// 續上頁
	// 自簽署：簽發者（第 3 個參數）就是自己
	der, _ := x509.CreateCertificate(rand.Reader, tmpl, tmpl,
		&key.PublicKey, key)
	kb, _ := x509.MarshalPKCS8PrivateKey(key)
	certPEM := &pem.Block{Type: "CERTIFICATE", Bytes: der}
	keyPEM := &pem.Block{Type: "PRIVATE KEY", Bytes: kb}
	os.WriteFile("cert.pem", pem.EncodeToMemory(certPEM), 0o644)
	os.WriteFile("key.pem", pem.EncodeToMemory(keyPEM), 0o600)
}
```

<!--
CreateCertificate 的第二和第三個參數都是 tmpl，意思是「簽發者就是自己」，這就是「自簽署」的意思。正式的憑證，簽發者會是 CA。最後用 CA 的私鑰（這裡就是自己的私鑰）替憑證簽章。

最後跟剛剛的 RSA 金鑰一樣，轉成 PEM 格式存檔，私鑰用 0o600 權限。
-->

---
zoom: 0.76
---

# HTTPS 伺服器

```go
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /",
		func(w http.ResponseWriter, r *http.Request) {
			v := tls.VersionName(r.TLS.Version)
			fmt.Fprintf(w, "Hello over %s\n", v)
		})
	srv := &http.Server{
		Addr:              ":8443",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		// 禁用舊版 TLS
		TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	log.Println("https://localhost:8443")
	// 指定憑證與私鑰
	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

<div class="mt-4 p-3 bg-blue-50 border-l-4 border-blue-400 text-gray-700 text-sm text-left">
⚠️ 瀏覽器打開會出現「不安全」的警告，因為這張憑證不是可信任的 CA 簽發的 — 這正是憑證驗證在保護我們。
</div>

<!--
有了憑證，把第 15 章的伺服器改成 HTTPS 只需要一個改變：把 ListenAndServe 換成 ListenAndServeTLS，傳入憑證和私鑰的檔案路徑。

TLSConfig 設定 MinVersion 為 TLS 1.2，禁用有安全漏洞的舊版本。Go 預設就會優先使用 TLS 1.3，並且只啟用安全的加密套件，這是 Go 的 TLS 實作被廣泛稱讚的地方：預設值就是安全的。

處理器裡用 r.TLS.Version 取得這次連線使用的 TLS 版本。

用瀏覽器打開 https://localhost:8443，會看到「不安全」的警告，因為瀏覽器不信任我們自己簽的憑證。這個警告不是錯誤，而是憑證驗證機制在正常運作：它不認識這個簽發者，所以提醒我們可能被冒充了。
-->

---

# HTTPS 客戶端：信任自簽署憑證

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	caPEM, err := os.ReadFile("cert.pem")
	if err != nil {
		panic(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caPEM) // 把自簽署憑證加入「信任清單」
```

<!--
Go 的 HTTP 客戶端預設也會驗證憑證，連線到自簽署憑證的伺服器會出現 certificate signed by unknown authority 的錯誤。

正確的做法是：把這張自簽署憑證加入客戶端的「信任清單」。x509.NewCertPool 建立一個憑證池，AppendCertsFromPEM 把憑證加進去。
-->

---

# HTTPS 客戶端：信任自簽署憑證（續）

```go
	// 續上頁
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
		},
	}}
	resp, err := client.Get("https://localhost:8443/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Print(string(body)) // Hello over TLS 1.3
}
```

<!--
再把憑證池設定到 TLSClientConfig 的 RootCAs。這樣客戶端就只信任這張憑證，其他的驗證照常進行。

連線成功後，伺服器回傳 Hello over TLS 1.3。公司內部的服務之間互相呼叫，如果使用公司自己的 CA，也是用同樣的方式設定。
-->

---

# 使用 TLS 的注意事項

**注意事項之一：** 絕對不要在正式環境使用 `InsecureSkipVerify: true`（等於關掉驗證，任何人都能冒充）

**注意事項之二：** 私鑰權限 `0o600`、不要 commit 到 Git；憑證有**到期日**，要自動更新

**注意事項之三：** 正式環境使用 CA 簽發的憑證

| 情境 | 建議 |
| --- | --- |
| 公開網站 | Let's Encrypt（免費），或雲端負載平衡器代管憑證 |
| Go 程式自動申請憑證 | `golang.org/x/crypto/acme/autocert` |
| 公司內部服務 | 內部 CA + `RootCAs` 設定 |
| 單元測試 | `httptest.NewTLSServer`，自動產生測試用憑證 |

<!--
使用 TLS 有三個注意事項。

第一，網路上很多教學遇到憑證錯誤，就教大家設定 InsecureSkipVerify: true 關掉驗證。這在正式環境是非常危險的：關掉驗證之後，任何人都能冒充伺服器，HTTPS 的保護就完全失效了。正確的做法是像剛剛一樣，把憑證加入信任清單。

第二，私鑰要妥善保管，權限設為 0600，而且絕對不要 commit 到 Git。憑證都有到期日，Let's Encrypt 的憑證只有 90 天，要設定自動更新。

第三，正式環境使用 CA 簽發的憑證。Go 有一個 autocert 套件，可以讓程式自動向 Let's Encrypt 申請和更新憑證，非常方便。寫單元測試的時候，httptest.NewTLSServer 會自動產生測試用的憑證。
-->

---
layout: default
---

# 綜合練習：安全的設定檔
### 任務說明

寫一個工具，把含有密碼的設定檔**加密後儲存**，並能確認檔案沒被竄改：

1. 從環境變數 `CONFIG_KEY` 讀取 64 個十六進位字元的金鑰（`hex.DecodeString` → 32 bytes）；沒有設定時印出錯誤並結束
2. `encryptFile(key []byte, in, out string) error`：讀取 JSON 設定檔 → AES-GCM 加密 → 寫入 `out`（權限 `0o600`）
3. `decryptFile(key []byte, path string) (map[string]any, error)`：解密並解析 JSON；解密失敗時回傳 `"設定檔已損毀或金鑰錯誤"`
4. 另外輸出 `out + ".sha256"`，內容是加密檔的 SHA-256，讓維運人員不需要金鑰也能檢查檔案完整性
5. 提示：產生金鑰 `openssl rand -hex 32`

<!--
這個綜合練習模擬實務上很常見的需求：設定檔裡有資料庫密碼、API 金鑰這些機密，不能用明文存在伺服器上。

用到了今天的 AES-GCM、SHA-256，以及第 11 章的 JSON、第 12 章的檔案操作和權限。

注意第 4 步：為什麼除了 GCM 本身的完整性保護，還要另外輸出 SHA-256？因為 GCM 要有金鑰才能驗證，而 SHA-256 不需要金鑰，維運人員可以在不接觸機密的情況下，檢查檔案有沒有損壞。
-->

---

# 綜合練習：解題提示
### 提示說明

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func gcmFor(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCMWithRandomNonce(block)
}
```

<!--
gcmFor 把金鑰轉成 GCM 物件。金鑰長度不對時，aes.NewCipher 會回傳錯誤。
-->

---

# 綜合練習：解題提示（續）
### 提示說明

```go
// 續上頁
func encryptFile(key []byte, in, out string) error {
	plain, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	gcm, err := gcmFor(key)
	if err != nil {
		return err
	}
	ct := gcm.Seal(nil, nil, plain, nil)
	sum := sha256.Sum256(ct)
	digest := hex.EncodeToString(sum[:]) + "\n"
	os.WriteFile(out+".sha256", []byte(digest), 0o644)
	return os.WriteFile(out, ct, 0o600)
}
```

<!--
encryptFile 讀取明文的設定檔、加密、計算密文的 SHA-256 寫到 .sha256 檔，最後把密文寫入，權限 0600。
-->

---

# 綜合練習：解題提示（續 2）
### 提示說明

```go
// 續上頁
func decryptFile(key []byte, path string) (map[string]any, error) {
	ct, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	gcm, err := gcmFor(key)
	if err != nil {
		return nil, err
	}
	plain, err := gcm.Open(nil, nil, ct, nil)
	if err != nil {
		return nil, errors.New("設定檔已損毀或金鑰錯誤")
	}
	var cfg map[string]any
	return cfg, json.Unmarshal(plain, &cfg)
}
```

<!--
decryptFile 讀取密文、解密、解析 JSON。解密失敗時，不管是檔案被改過還是金鑰錯誤，都回傳同一個籠統的錯誤訊息，不要透露太多細節給可能的攻擊者。
-->

---

# 綜合練習：解題提示（續 3）
### 提示說明

```go
// 續上頁
func main() {
	key, err := hex.DecodeString(os.Getenv("CONFIG_KEY"))
	if err != nil || len(key) != 32 {
		fmt.Println("請設定 CONFIG_KEY（64 個十六進位字元）")
		os.Exit(1)
	}
	sample := []byte(`{"db_pass":"s3cret","port":8080}`)
	os.WriteFile("app.json", sample, 0o600)
	if err := encryptFile(key, "app.json", "app.json.enc"); err != nil {
		panic(err)
	}
	cfg, err := decryptFile(key, "app.json.enc")
	fmt.Println(cfg["port"], err) // 8080 <nil>
}
```

<!--
main 從環境變數讀取金鑰，用 hex.DecodeString 轉成位元組，長度不是 32 就結束。接著建立一個範例設定檔、加密、再解密讀回來。

執行方式：export CONFIG_KEY=$(openssl rand -hex 32)，再 go run .。
-->

---
layout: section
class: flex flex-col justify-center items-center text-center
---

# GoShop 專案實作
## 第 18 步：後台登入與 HTTPS

<!--
回到 GoShop。現在的後台有一個很嚴重的問題：任何人只要知道網址 /admin，就能看到所有訂單、隨意修改商品價格。

而且資料是用 HTTP 明文傳送的，同一個咖啡廳 Wi-Fi 的人都看得到。今天學的雜湊、HMAC 和 TLS，剛好可以解決這兩個問題。
-->

---

# GoShop 第 18 步：後台登入與 HTTPS
### 任務說明

1. `internal/auth`：
   - `HashPassword`／`CheckPassword`：用 **bcrypt** 產生與比對密碼雜湊
   - `NewToken`／`Verify`：「到期時間.簽章」格式的憑證，用 **HMAC-SHA256** 簽章
2. `web`：`GET/POST /login`、`POST /logout`；`requireAdmin` 中介軟體保護 `/admin`
3. 登入成功後設定 cookie：`HttpOnly`、`Secure`（HTTPS 時）、`SameSite=Lax`
4. `-hash-password`：從標準輸入讀取密碼，印出 bcrypt 雜湊，設定到 `GOSHOP_ADMIN_HASH`
5. `-tls-cert`、`-tls-key`：改用 HTTPS，**最低 TLS 1.3**

```bash
echo 'goshop-admin' | go run . -hash-password   # 印出 bcrypt 雜湊
export GOSHOP_ADMIN_HASH='$2a$10$xpeieGHIjk…'   # 貼上剛剛的結果
go run $(go env GOROOT)/src/crypto/tls/generate_cert.go \
    -host localhost,127.0.0.1 -ecdsa-curve P256
go run . -http :8443 -tls-cert cert.pem -tls-key key.pem
```

<!--
這一步要替後台加上登入機制，並且改用 HTTPS。

密碼的部分，就是本章「使用雜湊的注意事項」講的：絕對不能存明文，也不能用 SHA-256，要用 bcrypt 這種刻意設計得很慢的函式。

登入之後，我們要發給瀏覽器一個「通行證」，之後的每個請求都帶著它。這個通行證用 HMAC 簽章，只有知道金鑰的伺服器能產生，別人改了一個字，簽章就對不上了。

憑證的部分，用 Go 內建的 generate_cert.go 產生 ECDSA 的自簽憑證，在本機測試 HTTPS。
-->

---
class: code-sm
---

# GoShop 第 18 步：解題提示
### bcrypt 與 HMAC 簽章憑證

```go
// goshop/internal/auth/auth.go
// CheckPassword 比對密碼是否正確。
func (a *Auth) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(password))
	return err == nil
}

// NewToken 產生「到期時間.簽章」格式的憑證，例如 1790000000.q7Xk…
func (a *Auth) NewToken(now time.Time) string {
	exp := strconv.FormatInt(now.Add(a.TTL).Unix(), 10)
	return exp + "." + a.sign(exp)
}
// ...
func (a *Auth) sign(msg string) string {
	mac := hmac.New(sha256.New, a.Secret)
	mac.Write([]byte(msg))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
```

<!--
CheckPassword 用 bcrypt 的 CompareHashAndPassword 比對，傳回 nil 就代表密碼正確。bcrypt 的雜湊裡已經包含了鹽和成本參數，所以比對時不需要另外存鹽。

NewToken 產生登入憑證，格式是「到期時間點簽章」。到期時間是 Unix 秒數，簽章是用 HMAC-SHA256 對到期時間簽出來的，再用 URL 安全的 Base64 編碼，這樣放進 cookie 不會有特殊字元的問題。

有人想把到期時間改成 9999999999 讓自己永遠登入？沒用，因為他算不出新的簽章，他沒有金鑰。
-->

---
zoom: 0.79
---

# GoShop 第 18 步：解題提示（續）
### 驗證憑證：hmac.Equal

```go
// goshop/internal/auth/auth.go
// Verify 檢查簽章是否正確、是否還沒過期。
func (a *Auth) Verify(token string, now time.Time) error {
	exp, sig, ok := strings.Cut(token, ".")
	// hmac.Equal 花費的時間固定，避免被用「計時攻擊」猜出簽章
	if !ok || !hmac.Equal([]byte(sig), []byte(a.sign(exp))) {
		return ErrInvalidToken
	}
	unix, err := strconv.ParseInt(exp, 10, 64)
	if err != nil || now.Unix() >= unix {
		return ErrInvalidToken
	}
	return nil
}
```

| 測試案例 | 結果 |
| --- | --- |
| 59 分鐘後 | ✅ 有效 |
| 剛好 1 小時後 | ❌ 過期 |
| 竄改到期時間 | ❌ 簽章不符 |
| 用別的金鑰簽的 | ❌ 簽章不符 |

<!--
Verify 先用 strings.Cut 把憑證切成到期時間和簽章兩段，再用同一把金鑰重新簽一次到期時間，比對兩個簽章是否相同。

比對一定要用 hmac.Equal，不能用 ==。== 比對字串時，遇到第一個不同的字元就會停下來，攻擊者可以量測回應時間，一個字元一個字元猜出正確的簽章。hmac.Equal 不管哪裡不同，花的時間都一樣。

簽章正確之後，才檢查有沒有過期。下面的表格是 auth_test.go 裡的測試案例，每一種竄改都會被擋下來。
-->

---
class: code-sm
---

# GoShop 第 18 步：解題提示（續 2）
### 登入中介軟體與安全的 cookie

```go
// goshop/internal/web/login.go
// requireAdmin 是中介軟體：沒有有效的登入憑證就導向登入頁。
func (s *Server) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil || s.Auth.Verify(c.Value, time.Now()) != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
```

```go
// goshop/internal/web/login.go
		HttpOnly: true,         // JavaScript 讀不到，降低 XSS 竊取的風險
		Secure:   r.TLS != nil, // 使用 HTTPS 時，只在加密連線中傳送
		SameSite: http.SameSiteLaxMode,
```

<!--
requireAdmin 是一個中介軟體，它接收一個處理器，傳回一個加上登入檢查的新處理器。沒有 cookie，或是 cookie 裡的憑證無效，就導向登入頁；通過檢查才交給原本的處理器。

在路由的地方，把後台的兩個處理器包起來：s.requireAdmin(s.adminPage)。API 是給顧客用的，不需要登入，所以不包。

cookie 的三個安全設定：HttpOnly 讓網頁上的 JavaScript 讀不到它；Secure 讓它只在 HTTPS 連線中傳送，r.TLS 不是 nil 就代表這是 HTTPS 請求；SameSite=Lax 讓別的網站發起的 POST 請求不會帶上這個 cookie，可以防範 CSRF 攻擊。
-->

---
zoom: 0.88
---

# GoShop 第 18 步：解題提示（續 3）
### HTTPS：TLS 1.3

```go
// goshop/main.go
	srv := &http.Server{
		Addr:              opt.addr,
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS13},
	}
	errCh := make(chan error, 1)
	go func() {
		if opt.tlsCert != "" {
			errCh <- srv.ListenAndServeTLS(opt.tlsCert, opt.tlsKey)
		} else {
			errCh <- srv.ListenAndServe()
		}
	}()
```

```text
$ curl --cacert cert.pem https://localhost:8443/admin
→ 303 導向 /login
$ curl --cacert cert.pem -d password=goshop-admin \
    https://localhost:8443/login
→ Set-Cookie: goshop_admin=1790263391.u…; HttpOnly; Secure
$ curl --tls-max 1.2 --cacert cert.pem https://localhost:8443/
→ 連線失敗（伺服器只接受 TLS 1.3）
```

<!--
改用 HTTPS 只需要兩個改變：TLSConfig 設定最低版本是 TLS 1.3，以及把 ListenAndServe 換成 ListenAndServeTLS，傳入憑證和私鑰檔。

用 curl 測試：沒登入打開 /admin 會被導向登入頁；用正確的密碼登入，拿到一個有 HttpOnly 和 Secure 的 cookie；最後故意只允許 TLS 1.2 連線，伺服器直接拒絕。

簽章用的金鑰從環境變數 GOSHOP_SECRET 讀取。沒有設定的時候，程式用 Go 1.24 的 crypto/rand.Text 產生隨機金鑰，缺點是伺服器重新啟動後，所有人都要重新登入。正式環境請一定要設定 GOSHOP_SECRET。
-->

---

# 章節總結

- **原則**：不要自己發明密碼學；隨機值一律用 **`crypto/rand`**
- **雜湊**：`sha256.Sum256` / `sha3.Sum256`（1.24+）；MD5 已不安全；大檔案用 `io.Copy` 串流
- **HMAC**：加上金鑰的雜湊，比對用 **`hmac.Equal`**；密碼用 **bcrypt / Argon2id / PBKDF2**，不要用 SHA-256
- **對稱式加密**：**AES-256-GCM**；`cipher.NewGCMWithRandomNonce`（1.24+）自動處理 nonce；解密失敗 = 被竄改
- **非對稱式加密**：**RSA-OAEP**，金鑰 ≥ 2048 bits；只用來加密小資料（例如 AES 金鑰）
- **數位簽章**：**Ed25519**，私鑰簽、公鑰驗；證明來源與完整性
- **TLS**：憑證 = 公鑰 + CA 簽章；ECDSA P-256 更短更快；`ListenAndServeTLS`；信任自簽憑證用 `RootCAs`，**不要 `InsecureSkipVerify`**
- **GoShop**：bcrypt 存管理員密碼、HMAC 簽署登入憑證、安全的 cookie 設定；ECDSA 自簽憑證 + TLS 1.3

下一章是最後一章：Go 語言的特殊套件 **reflect 與 unsafe**。

<!--
我們來整理今天學到的東西。

密碼學的第一原則是不要自己發明，第二原則是隨機值用 crypto/rand。雜湊用來檢查完整性，HMAC 加上金鑰，密碼要用專門的密碼雜湊函式。對稱式加密用 AES-GCM，非對稱式加密用 RSA-OAEP，數位簽章用 Ed25519。HTTPS 把這些工具組合在一起：憑證證明身分、非對稱式加密交換金鑰、對稱式加密傳輸資料。

GoShop 這一步補上了安全性：管理員密碼只存 bcrypt 雜湊；登入後發一個用 HMAC 簽章的憑證，放在 HttpOnly、Secure 的 cookie 裡，比對簽章用 hmac.Equal；伺服器加上 TLS 1.3，用 ECDSA 的自簽憑證在本機測試 HTTPS。

下一章是這門課的最後一章，我們要看 Go 的兩個特殊套件：reflect 和 unsafe。第 11 章的 JSON 套件是怎麼讀到 struct 標籤的？答案就是 reflect。unsafe 則是讓我們繞過 Go 型別系統的「後門」，了解它才知道為什麼 Go 平常這麼安全。
-->

---
layout: end
---

# Q & A

有任何問題嗎？

<!--
資訊安全是一個很深的領域，今天只是入門，但只要記住「用對工具、不要自己發明」，就能避開大部分的陷阱。

課後建議：把第 15 章的待辦事項 API 改成 HTTPS，再加上一個用 HMAC 簽章的 API 金鑰驗證中介軟體，練習把今天的工具應用在真實的服務上。

有問題的同學現在可以提問！
-->

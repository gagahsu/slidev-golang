// 檢查投影片內容是否超出版面（需先啟動 dev server）
// 用法：node scripts/check-overflow.mjs [port=3030] [起始頁] [結束頁]
import { chromium } from 'playwright-chromium'

const port = process.argv[2] ?? '3030'
const from = parseInt(process.argv[3] ?? '1', 10)
const to = process.argv[4] ? parseInt(process.argv[4], 10) : Infinity

const browser = await chromium.launch(process.env.CHROMIUM_PATH ? { executablePath: process.env.CHROMIUM_PATH } : {})
const page = await browser.newPage({ viewport: { width: 1280, height: 720 } })
await page.goto(`http://localhost:${port}/1`, { waitUntil: 'networkidle' })
const total = await page.evaluate(() => window.__slidev__?.nav?.total ?? 0)
let bad = 0
for (let n = from; n <= Math.min(total, to); n++) {
  await page.evaluate(n => window.__slidev__.nav.go(n), n)
  await page.waitForTimeout(700)
  const r = await page.evaluate((n) => {
    const layout = document.querySelector(`.slidev-page-${n} .slidev-layout`)
    if (!layout) return null
    const box = layout.getBoundingClientRect()
    let maxBottom = box.top
    for (const el of layout.querySelectorAll('*')) {
      const b = el.getBoundingClientRect()
      if (b.height > 0 && b.width > 0) maxBottom = Math.max(maxBottom, b.bottom)
    }
    const title = layout.querySelector('h1')?.textContent?.trim() ?? ''
    // 程式碼行被自動換行（太長）的數量
    let wrapped = 0
    for (const line of layout.querySelectorAll('pre code .line')) {
      const lh = parseFloat(getComputedStyle(line).lineHeight) || 20
      if (line.getBoundingClientRect().height > lh * 1.5) wrapped++
    }
    return { over: Math.round(maxBottom - box.bottom), h: Math.round(box.height), title, wrapped }
  }, n)
  if (r && r.over > 2) {
    bad++
    console.log(`第 ${n} 頁 超出 ${r.over}px（版面高 ${r.h}px）：${r.title}`)
  }
  if (r && r.wrapped > 0) console.log(`第 ${n} 頁 有 ${r.wrapped} 行程式碼過長而換行：${r.title}`)
}
console.log(`檢查 ${from}～${Math.min(total, to)} 頁，${bad} 頁超出版面`)
await browser.close()

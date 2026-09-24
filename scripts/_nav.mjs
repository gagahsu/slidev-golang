import { chromium } from 'playwright-chromium'
const browser = await chromium.launch({ executablePath: '/opt/pw-browsers/chromium' })
const page = await browser.newPage({ viewport: { width: 1280, height: 720 } })
await page.goto('http://localhost:3060/1', { waitUntil: 'networkidle' })
const total = await page.evaluate(() => window.__slidev__.nav.total)
console.log('total slides', total)
for (const n of ['00','05','13','19']) {
  await page.goto('http://localhost:3060/1', { waitUntil: 'networkidle' })
  await page.click(`a[href$="/ch${n}"], .chapter-card >> text="Ch ${parseInt(n)}"`)
  await page.waitForTimeout(1200)
  const r = await page.evaluate(() => ({ no: window.__slidev__.nav.currentPage, h1: document.querySelector(`.slidev-page-${window.__slidev__.nav.currentPage} h1`)?.textContent }))
  console.log(n, r)
  // back link
  await page.click('.slidev-page-' + r.no + ' >> text=← 返回目錄')
  await page.waitForTimeout(800)
  console.log('  back ->', await page.evaluate(() => window.__slidev__.nav.currentPage))
}
await browser.close()

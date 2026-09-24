// 檢查投影片裡的 Go 程式碼區塊能不能編譯（go vet）
//
// 規則：
//   - 只檢查以 `package` 開頭（前面可有註解）、可獨立編譯的 ```go 區塊（片段程式碼不檢查）
//   - 區塊前一行有 `<!-- check:skip -->` 的區塊略過（例如刻意示範 go vet 警告）
//   - 有 `// 編譯錯誤` 標記（行尾或獨立一行）的區塊略過（刻意示範錯誤的程式碼；註解掉的程式碼行不算）
//   - import 了第三方模組或範例模組（github.com／golang.org／gopkg.in／example.com）的區塊略過
//   - import 了 "testing" 的區塊會存成 x_test.go
//   - 第一行是 `// 檔名：xxx.go` 的區塊，會和同檔案中前一個「也有檔名」的區塊放在同一個套件資料夾（多檔案專案、程式 + 測試檔）
//   - 第一行是 `// 續上頁` 的區塊，會接在同檔案上一個區塊後面一起檢查（跨頁的長程式）
//   - 第一行是 `// goshop/…` 的區塊是 GoShop 專案的摘錄，不在這裡檢查（見 check-project.mjs）
//
// 用法：pnpm check:go            （檢查全部章節）
//       pnpm check:go ch05       （只檢查 ch05）
//       pnpm check:go ch05 --run （檢查後逐一執行，印出輸出）
// 可用 GOTOOLCHAIN 環境變數指定 Go 版本，例如 GOTOOLCHAIN=go1.27.0 pnpm check:go
// WITH_DEPS=1：一併檢查使用第三方模組的區塊（會執行 go mod tidy 下載相依模組）
import { readdirSync, readFileSync, mkdirSync, writeFileSync, rmSync } from 'fs'
import { spawnSync } from 'child_process'
import { join } from 'path'
import { fileURLToPath } from 'url'

const root = fileURLToPath(new URL('..', import.meta.url))
const work = join(root, '.snippet-check')
const args = process.argv.slice(2)
const run = args.includes('--run')
const filter = args.find(a => !a.startsWith('--'))

if (!process.env.KEEP) rmSync(work, { recursive: true, force: true })
mkdirSync(work, { recursive: true })
writeFileSync(join(work, 'go.mod'), 'module snippets\n\ngo 1.27\n')

const files = readdirSync(root)
  .filter(f => /^ch\d+.*\.md$/.test(f))
  .filter(f => !filter || f.startsWith(filter))
  .sort()

const origin = new Map()
let count = 0
let skipped = 0

// REWRITE="https://httpbin.org=http://127.0.0.1:8081"：驗證時把網址換成本機的測試伺服器
const rewrites = (process.env.REWRITE ?? '').split(',').filter(Boolean).map(r => r.split('='))
const rewrite = code => rewrites.reduce((c, [from, to]) => c.replaceAll(from, to), code)

const blocks = [] // { files: [{ name, code }], where }
const withDeps = process.env.WITH_DEPS === '1'
const skipRe = withDeps
  ? /import[\s\S]*?"example\.com\//
  : /import[\s\S]*?"(github\.com|golang\.org|gopkg\.in|example\.com)\//
const stripComments = code => code.replace(/^(\s*\/\/[^\n]*\n)+/, '')

for (const file of files) {
  const lines = readFileSync(join(root, file), 'utf8').split('\n')
  let last = null // 同一檔案中上一個可檢查的區塊（供「續上頁」「檔名」接續）
  for (let i = 0; i < lines.length; i++) {
    if (!/^```go\b/.test(lines[i])) continue
    const start = i + 1
    let prev = i - 1
    while (prev >= 0 && lines[prev].trim() === '') prev--
    const skipMarked = prev >= 0 && lines[prev].includes('check:skip') // 刻意示範 go vet 警告等情況
    let end = start
    while (end < lines.length && !/^```\s*$/.test(lines[end])) end++
    const code = lines.slice(start, end).join('\n')
    i = end
    if (/^\/\/ goshop\//.test(code)) { // GoShop 專案摘錄：由 check-project.mjs 對照 goshop/ 資料夾檢查
      last = null
      continue
    }
    if (/^\s*\/\/ 續上頁/.test(code)) {
      if (last?.files) last.files.at(-1).code += (code.startsWith('\t') ? '\n' : '\n\n') + code // 以 Tab 縮排的續上頁：函式本體中途接續
      continue
    }
    if (!/^\s*package \w+/.test(stripComments(code))) continue
    const named = code.match(/^\/\/ 檔名：(\S+\.go)/)
    if (skipMarked || /^\s*([^\s/].*)?\/\/\s*編譯錯誤/m.test(code) || skipRe.test(code) ||
        (named && last?.skippedProject)) {
      skipped++
      // 多檔案專案中有檔案被略過時，同一專案的其他檔案也一併略過
      last = named ? { skippedProject: true } : null
      continue
    }
    const name = named ? named[1].split('/').pop()
      : code.includes('"testing"') ? 'x_test.go' : 'main.go'
    const pkgOf = c => stripComments(c).match(/^\s*package (\w+)/)?.[1]
    if (named && last?.named && last.files && pkgOf(last.files[0].code) === pkgOf(code) && !last.files.some(f => f.name === name)) {
      last.files.push({ name, code }) // 同一個套件的另一個檔案（例如測試檔）
      continue
    }
    last = { files: [{ name, code }], where: `${file}:${start}`, named: Boolean(named) }
    blocks.push(last)
  }
}

for (const b of blocks) {
  const dir = `s${String(++count).padStart(4, '0')}`
  mkdirSync(join(work, dir))
  for (const f of b.files) writeFileSync(join(work, dir, f.name), rewrite(f.code) + '\n')
  // //go:embed 需要的檔案：建立空的佔位檔案，讓範例可以編譯
  for (const f of b.files) {
    for (const [, pat] of f.code.matchAll(/^\/\/go:embed (\S+)/gm)) {
      const target = join(work, dir, pat.includes('*') ? pat.slice(0, pat.lastIndexOf('/')) : pat)
      mkdirSync(target, { recursive: true })
      const ext = pat.includes('*.') ? pat.slice(pat.lastIndexOf('.')) : '.html'
      writeFileSync(join(target, 'placeholder' + ext), '')
    }
  }
  origin.set(dir, b.where)
}

console.log(`檢查 ${count} 個程式碼區塊（略過 ${skipped} 個）...`)
if (withDeps) {
  const tidy = spawnSync('go', ['mod', 'tidy'], { cwd: work, encoding: 'utf8' })
  if (tidy.status !== 0) console.error(tidy.stderr)
}
const res = spawnSync('go', ['vet', './...'], { cwd: work, encoding: 'utf8' })
const out = (res.stdout + res.stderr).replace(/(?:\.\/)?(s\d{4})\/(\w+\.go)?/g,
  (m, dir) => origin.has(dir) ? `[${origin.get(dir)}] ` : m)
if (res.status !== 0) {
  console.error(out)
  process.exit(1)
}
console.log('✔ 全部通過')

// gofmt：列出排版不符合官方格式的區塊（只提示，不算失敗）
const fmtRes = spawnSync('gofmt', ['-l', '.'], { cwd: work, encoding: 'utf8' })
const unformatted = fmtRes.stdout.trim().split('\n').filter(Boolean)
for (const f of unformatted) {
  const dir = f.split('/')[0]
  console.log(`⚠ gofmt 格式不符：${origin.get(dir) ?? f}`)
}

// --run：實際執行每個 main 程式並印出輸出，方便對照投影片上標註的結果
if (run) {
  for (const [dir, where] of origin) {
    const names = readdirSync(join(work, dir)).filter(n => n.endsWith('.go'))
    const src = names.map(n => readFileSync(join(work, dir, n), 'utf8')).join('\n')
    const hasTest = names.some(n => n.endsWith('_test.go'))
    if (!hasTest && (!/^package main\b/m.test(src) || !/func main\(\)/.test(src))) continue
    // 先編譯再執行，逾時（例如伺服器程式）時直接結束執行檔，不留下背景程序
    const bin = join(work, dir, 'app.bin')
    const r = hasTest
      ? spawnSync('go', ['test', '-v', `./${dir}`], { cwd: work, encoding: 'utf8', timeout: 60000 })
      : (spawnSync('go', ['build', '-o', bin, `./${dir}`], { cwd: work }),
        spawnSync(bin, [], { cwd: join(work, dir), encoding: 'utf8', timeout: 5000 }))
    console.log(`\n── ${where}`)
    console.log((r.stdout + r.stderr).trimEnd())
  }
}
if (!process.env.KEEP) rmSync(work, { recursive: true, force: true })

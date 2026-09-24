// 檢查投影片裡的 Go 程式碼區塊能不能編譯（go vet）
//
// 規則：
//   - 只檢查以 `package` 開頭、可獨立編譯的 ```go 區塊（片段程式碼不檢查）
//   - 程式碼行尾有 `// 編譯錯誤` 標記者略過（刻意示範錯誤的程式碼；註解掉的行不算）
//   - import 了第三方模組（github.com／golang.org／gopkg.in）的區塊略過
//   - import 了 "testing" 的區塊會存成 x_test.go
//   - 第一行是 `// 續上頁` 的區塊，會接在同檔案上一個區塊後面一起檢查（跨頁的長程式）
//
// 用法：pnpm check:go            （檢查全部章節）
//       pnpm check:go ch05       （只檢查 ch05）
//       pnpm check:go ch05 --run （檢查後逐一執行，印出輸出）
// 可用 GOTOOLCHAIN 環境變數指定 Go 版本，例如 GOTOOLCHAIN=go1.27.0 pnpm check:go
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

const blocks = [] // { code, where }

for (const file of files) {
  const lines = readFileSync(join(root, file), 'utf8').split('\n')
  let last = null // 同一檔案中上一個可檢查的區塊（供「續上頁」接續）
  for (let i = 0; i < lines.length; i++) {
    if (!/^```go\b/.test(lines[i])) continue
    const start = i + 1
    let end = start
    while (end < lines.length && !/^```\s*$/.test(lines[end])) end++
    const code = lines.slice(start, end).join('\n')
    i = end
    if (/^\s*\/\/ 續上頁/.test(code) && last) {
      last.code += '\n\n' + code
      continue
    }
    if (!/^\s*package \w+/.test(code)) continue
    if (/^\s*[^\s/].*\/\/\s*編譯錯誤/m.test(code) || /import[\s\S]*?"(github\.com|golang\.org|gopkg\.in)\//.test(code)) {
      skipped++
      last = null
      continue
    }
    last = { code, where: `${file}:${start}` }
    blocks.push(last)
  }
}

for (const b of blocks) {
  const dir = `s${String(++count).padStart(4, '0')}`
  mkdirSync(join(work, dir))
  const name = b.code.includes('"testing"') ? 'x_test.go' : 'main.go'
  writeFileSync(join(work, dir, name), b.code + '\n')
  origin.set(dir, b.where)
}

console.log(`檢查 ${count} 個程式碼區塊（略過 ${skipped} 個）...`)
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
    const r = spawnSync('go', ['run', `./${dir}`], { cwd: work, encoding: 'utf8', timeout: 20000 })
    console.log(`\n── ${where}`)
    console.log((r.stdout + r.stderr).trimEnd())
  }
}
if (!process.env.KEEP) rmSync(work, { recursive: true, force: true })

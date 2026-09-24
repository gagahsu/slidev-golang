// 檢查 GoShop 課程專案（goshop/ch00 … goshop/ch19）
//
// 1. 每一步的參考解答都要：gofmt 排版正確、go vet 通過、go test 通過
// 2. 投影片中第一行是「// goshop/<路徑>」的程式碼區塊（HTML 用 <!-- goshop/… -->、
//    Makefile 用 # goshop/…），內容必須和 goshop/chNN/<路徑> 檔案一致：
//    - 區塊內容（去掉第一行）必須依序出現在檔案中
//    - 單獨一行的 `// ...`（或 `# ...`、`<!-- ... -->`）代表省略，前後兩段分開比對
//
// 用法：pnpm check:project            （全部）
//       pnpm check:project ch13       （只檢查 ch13）
//       pnpm check:project --no-test  （只比對投影片與檔案，不跑 go test）
// 設定 GOSHOP_TEST_DSN 時，ch13 以後的 MySQL 測試也會執行。
import { existsSync, readdirSync, readFileSync } from 'fs'
import { spawnSync } from 'child_process'
import { join } from 'path'
import { fileURLToPath } from 'url'

const root = fileURLToPath(new URL('..', import.meta.url))
const args = process.argv.slice(2)
const filter = args.find(a => !a.startsWith('--'))
const runTests = !args.includes('--no-test')

const decks = readdirSync(root).filter(f => /^ch\d+.*\.md$/.test(f)).sort()
const header = /^(?:\/\/|#|<!--)\s*goshop\/(\S+?)(?:\s*-->)?\s*$/
const elision = /^\s*(?:\/\/|#|<!--)\s*\.\.\.\s*(?:-->)?\s*$/
let failed = 0
let excerpts = 0

const fail = msg => {
  failed++
  console.error('✘ ' + msg)
}

for (const deck of decks) {
  const ch = deck.slice(0, 4)
  if (filter && ch !== filter) continue
  const dir = join(root, 'goshop', ch)
  const lines = readFileSync(join(root, deck), 'utf8').split('\n')

  // 1. 投影片摘錄與檔案內容比對
  for (let i = 0; i < lines.length; i++) {
    if (!/^```\w*/.test(lines[i])) continue
    let end = i + 1
    while (end < lines.length && !/^```\s*$/.test(lines[end])) end++
    const body = lines.slice(i + 1, end)
    const where = `${deck}:${i + 2}`
    i = end
    const m = body[0]?.match(header)
    if (!m) continue
    excerpts++
    const file = join(dir, m[1])
    if (!existsSync(file)) {
      fail(`${where} 找不到檔案 goshop/${ch}/${m[1]}`)
      continue
    }
    const content = readFileSync(file, 'utf8')
    const chunks = [[]]
    for (const l of body.slice(1)) {
      if (elision.test(l)) chunks.push([])
      else chunks.at(-1).push(l)
    }
    let pos = 0
    for (const chunk of chunks) {
      while (chunk.length && chunk[0].trim() === '') chunk.shift()
      while (chunk.length && chunk.at(-1).trim() === '') chunk.pop()
      if (!chunk.length) continue
      const text = chunk.join('\n')
      const at = content.indexOf(text, pos)
      if (at < 0) {
        // 找出第一行對不上的地方，方便修正
        const bad = chunk.find(l => !content.includes(l)) ?? chunk[0]
        fail(`${where} 與 goshop/${ch}/${m[1]} 不一致：\n    ${bad}`)
        break
      }
      pos = at + text.length
    }
  }

  // 2. 參考解答本身要能通過檢查
  if (!runTests || !existsSync(dir)) continue
  const go = (...a) => spawnSync('go', a, { cwd: dir, encoding: 'utf8' })
  const fmt = spawnSync('gofmt', ['-l', '.'], { cwd: dir, encoding: 'utf8' })
  if (fmt.stdout.trim()) fail(`goshop/${ch} gofmt 格式不符：${fmt.stdout.trim()}`)
  for (const cmd of [['vet', './...'], ['test', '-count=1', './...']]) {
    const r = go(...cmd)
    if (r.status !== 0) {
      fail(`goshop/${ch} go ${cmd.join(' ')}\n${(r.stdout + r.stderr).trim()}`)
      break
    }
  }
  console.log(`✔ goshop/${ch}`)
}

console.log(`比對 ${excerpts} 個投影片摘錄`)
if (failed) {
  console.error(`共 ${failed} 個問題`)
  process.exit(1)
}
console.log('✔ 全部通過')

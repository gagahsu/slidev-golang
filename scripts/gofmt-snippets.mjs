// 用 gofmt 重新排版投影片中以 `package` 開頭的 ```go 區塊，直接寫回 .md
// 以 `// 續上頁` 開頭的區塊會和前一個區塊合併後一起排版，再拆回原本的投影片
// 用法：node scripts/gofmt-snippets.mjs [ch05]
import { readdirSync, readFileSync, writeFileSync } from 'fs'
import { spawnSync } from 'child_process'
import { join } from 'path'
import { fileURLToPath } from 'url'

const root = fileURLToPath(new URL('..', import.meta.url))
const filter = process.argv[2]
const files = readdirSync(root).filter(f => /^ch\d+.*\.md$/.test(f) && (!filter || f.startsWith(filter)))
const MARK = '// 續上頁'

for (const file of files) {
  const lines = readFileSync(join(root, file), 'utf8').split('\n')

  // 找出所有 go 區塊的位置
  const blocks = []
  for (let i = 0; i < lines.length; i++) {
    if (!/^```go\b/.test(lines[i])) continue
    let end = i + 1
    while (end < lines.length && !/^```\s*$/.test(lines[end])) end++
    blocks.push({ start: i + 1, end, code: lines.slice(i + 1, end).join('\n') })
    i = end
  }

  // 組成「鏈」：package 區塊 + 後續的續上頁區塊
  const chains = []
  for (const b of blocks) {
    if (b.code.startsWith(MARK) && chains.length && chains.at(-1).open) {
      chains.at(-1).parts.push(b)
    } else if (/^\s*package \w+/.test(b.code.replace(/^(\s*\/\/[^\n]*\n)+/, ''))) {
      chains.push({ parts: [b], open: true })
    } else if (chains.length) {
      chains.at(-1).open = false
    }
  }

  let changed = 0
  const replacements = new Map() // start → 新的程式碼
  for (const chain of chains) {
    const joined = chain.parts.map(p => p.code).join('\n\n')
    const r = spawnSync('gofmt', [], { input: joined + '\n', encoding: 'utf8' })
    if (r.status !== 0) continue
    const out = r.stdout.replace(/\n$/, '').split('\n')
    const pieces = [[]]
    for (const l of out) {
      if (l === MARK) pieces.push([])
      pieces.at(-1).push(l)
    }
    if (pieces.length !== chain.parts.length) continue
    pieces.forEach((p, k) => {
      while (p.length && p.at(-1) === '') p.pop()
      const code = p.join('\n')
      if (code !== chain.parts[k].code) {
        replacements.set(chain.parts[k].start, { end: chain.parts[k].end, code })
        changed++
      }
    })
  }

  if (!changed) continue
  const out = []
  for (let i = 0; i < lines.length; i++) {
    const rep = replacements.get(i)
    if (rep) {
      out.push(...rep.code.split('\n'))
      i = rep.end - 1
      continue
    }
    out.push(lines[i])
  }
  writeFileSync(join(root, file), out.join('\n'))
  console.log(`${file}：重新排版 ${changed} 個區塊`)
}

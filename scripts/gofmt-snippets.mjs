// 用 gofmt 重新排版投影片中以 `package` 開頭的 ```go 區塊，直接寫回 .md
// 用法：node scripts/gofmt-snippets.mjs [ch05]
import { readdirSync, readFileSync, writeFileSync } from 'fs'
import { spawnSync } from 'child_process'
import { join } from 'path'
import { fileURLToPath } from 'url'

const root = fileURLToPath(new URL('..', import.meta.url))
const filter = process.argv[2]
const files = readdirSync(root).filter(f => /^ch\d+.*\.md$/.test(f) && (!filter || f.startsWith(filter)))

for (const file of files) {
  const lines = readFileSync(join(root, file), 'utf8').split('\n')
  const out = []
  let changed = 0
  for (let i = 0; i < lines.length; i++) {
    out.push(lines[i])
    if (!/^```go\b/.test(lines[i])) continue
    let end = i + 1
    while (end < lines.length && !/^```\s*$/.test(lines[end])) end++
    const code = lines.slice(i + 1, end).join('\n')
    let result = code
    if (/^\s*package \w+/.test(code)) {
      const r = spawnSync('gofmt', [], { input: code + '\n', encoding: 'utf8' })
      if (r.status === 0) result = r.stdout.replace(/\n$/, '')
    }
    if (result !== code) changed++
    out.push(...result.split('\n'))
    i = end - 1
  }
  if (changed) {
    writeFileSync(join(root, file), out.join('\n'))
    console.log(`${file}：重新排版 ${changed} 個區塊`)
  }
}

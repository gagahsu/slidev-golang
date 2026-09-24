# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Slidev-based teaching presentation** for the course「Go 實戰開發」(Go Programming Masterclass, 20 chapters: Ch 0–19).
It is **not** a Go application — the content is about Go, but the repository itself is a Node.js/Slidev project.
All Go examples target **Go 1.27** (the latest release) and use modern idioms (range-over-int, `slices`/`maps`,
`errors.Join`/`errors.AsType`, `log/slog`, Go 1.22 `ServeMux` patterns, `sync.WaitGroup.Go`, `encoding/json/v2`, …).

## Commands

```bash
pnpm install          # Install dependencies (pnpm only, not npm/yarn)
pnpm dev              # Dev server at localhost:3030 with all chapters (index.md)
pnpm run ch16         # Start only one chapter deck (ch00 … ch19)
pnpm build            # Build to dist/ for deployment
pnpm run export:all   # Export chapter decks to dist/*.pdf (accepts "ch14" or "14-19")

pnpm check:go                 # Compile + go vet every standalone Go snippet in the slides
pnpm check:go ch05 --run      # …only ch05, and run each program / test to compare with the slides
pnpm fmt:go                   # gofmt every standalone Go snippet in place
pnpm check:width              # List code lines that are likely too wide for a slide
pnpm check:overflow 3030      # With a dev server running: report slides whose content overflows
```

The `.npmrc` sets `shamefully-hoist=true`, required by Slidev.

## Architecture

- `index.md` — Portal page (課程目錄) with chapter cards; imports every deck via `src:`.
- `chNN-<slug>.md` — One deck per chapter, `routeAlias: chNN`, dense 0-indexed numbering (`ch00` = 前置作業).
  The deck number, `routeAlias`, the `<Link to="chNN">` card and the `Ch N` label must always agree.
- `_template/` — Blueprint for new chapters (frontmatter, cover, callout, exercise pages).
- `global-bottom.vue` — Page `X / Y` footer; `style.css` — code-block and inline-code styling.
- `scripts/` — `export-all.mjs` (PDF export), `check-go-snippets.mjs`, `gofmt-snippets.mjs`,
  `check-line-width.py`, `check-overflow.mjs`, `set-zoom.py`.

## Slide Authoring Conventions

- **Language:** Traditional Chinese (zh-TW); English for code identifiers and technical terms.
- **Persona for speaker notes:** the 古古 (kucw.io) style used in `slidev-springboot`:
  address students as「我們／大家／同學」(never「您／各位」); 先情境再定義; everyday analogies that come back to code;
  say what a snippet does before it and what it prints after it; each chapter opens with a「回顧」slide and
  ends with「章節總結」+「下一章我們會介紹…」; section patterns:「什麼是 XXX？」「使用 XXX 的注意事項」「補充：」.
- **Every slide has presenter notes** (the last `<!-- -->` block of the slide).
- **Page types:** cover (HTML) → `layout: default` Outline → 回顧 → `layout: section` dividers → content →
  練習（任務說明 + 解題提示）per section → 綜合練習 → 章節總結 → `layout: end` Q&A.
- **Tables + code, not bullet walls.** Callouts use the blue `bg-blue-50 border-l-4` div.
- **Slide size (1280×720):** ≤ ~20 code lines on a code-only slide, fewer with a table or callout;
  code lines ≤ ~70 half-width columns (CJK ≈ 1.8). Prefer splitting a long solution into「（續）」slides;
  use per-slide `zoom:` (≥ 0.75) only for small overflows. Never use two-column layouts for code.

## Go Snippet Rules (enforced by `pnpm check:go`)

- Blocks starting with `package` (optionally after comment lines) must compile and pass `go vet`.
- Mark intentionally broken code with `// 編譯錯誤` (trailing, or on its own line), or put
  `<!-- check:skip -->` on the line before the code fence (e.g. a deliberate `go vet` demo).
- A block whose first line is `// 續上頁` continues the previous block (use a leading tab when it continues
  inside a function body).
- `// 檔名：xxx.go` blocks in the same package form one multi-file project (e.g. `store.go` + `api.go` + `main_test.go`).
- `WITH_DEPS=1` also checks blocks importing third-party modules;
  `REWRITE="https://httpbin.org=http://127.0.0.1:8081"` redirects URLs to a local test server when running.
- Keep snippets gofmt-clean (`pnpm fmt:go`); comments that make a line too long go on the line above.

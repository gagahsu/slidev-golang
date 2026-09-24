#!/usr/bin/env python3
"""為指定標題的投影片加上 zoom 設定：python3 scripts/set-zoom.py 檔案.md 0.9 "標題1" "標題2" ..."""
import re, sys

path, zoom, titles = sys.argv[1], sys.argv[2], sys.argv[3:]
lines = open(path, encoding='utf-8').read().split('\n')
kv = re.compile(r'^[\w-]+:\s*.*$')

for t in titles:
    try:
        idx = lines.index(f'# {t}')
    except ValueError:
        sys.exit(f'找不到投影片：{t}')
    j = idx - 1
    while j >= 0 and lines[j].strip() == '':
        j -= 1
    if lines[j] != '---':
        sys.exit(f'投影片標題前不是分隔線：{t}')
    # 往上找 frontmatter：連續的 key: value 行，再往上是 ---
    k = j - 1
    while k >= 0 and kv.match(lines[k]):
        k -= 1
    if k < j - 1 and k >= 0 and lines[k] == '---':
        block = [l for l in lines[k + 1:j] if not l.startswith('zoom:')] + [f'zoom: {zoom}']
        lines[k + 1:j] = block
    else:
        lines[j:j + 1] = ['---', f'zoom: {zoom}', '---']
open(path, 'w', encoding='utf-8').write('\n'.join(lines))

#!/usr/bin/env python3
"""列出投影片程式碼區塊中「顯示寬度」過長、會在投影片上自動換行的行。
用法：python3 scripts/check-line-width.py [寬度上限=70] [檔名前綴]
寬度計算：全形字（中日韓文字、全形標點）算 1.7，Tab 算 2，其他算 1。
實測 1280×720 的投影片，程式碼區塊一行約可容納 70 個半形字元。"""
import sys, glob, unicodedata

limit = float(sys.argv[1]) if len(sys.argv) > 1 else 70
prefix = sys.argv[2] if len(sys.argv) > 2 else 'ch'

def width(s):
    w = 0
    for ch in s:
        if ch == '\t':
            w += 2
        elif unicodedata.east_asian_width(ch) in ('W', 'F'):
            w += 1.7  # 全形字在程式碼字型中約為半形的 1.7 倍寬
        else:
            w += 1
    return w

total = 0
for path in sorted(glob.glob(f'{prefix}*.md')):
    lines = open(path, encoding='utf-8').read().split('\n')
    in_code = False
    for i, line in enumerate(lines, 1):
        if line.startswith('```'):
            in_code = not in_code
            continue
        if in_code and width(line) > limit:
            total += 1
            print(f'{path}:{i}: 寬度 {width(line):.0f}：{line.strip()[:60]}')
print(f'共 {total} 行超過 {limit}')

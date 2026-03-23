#!/usr/bin/env python3
"""CLI 工具：可獨立測試注音解碼"""

import sys
from .decoder import decode, detect_and_decode


def main():
    if len(sys.argv) > 1:
        text = ' '.join(sys.argv[1:])
    else:
        text = sys.stdin.read().strip()

    if not text:
        print("用法: zhuyin-decode <text>")
        print("範例: zhuyin-decode su3cl3")
        sys.exit(1)

    result = detect_and_decode(text)
    if result['is_zhuyin']:
        print(f"偵測到注音誤打！")
        print(f"原始: {result['original']}")
        print(f"注音: {result['zhuyin']}")
    else:
        print(f"未偵測到注音誤打")
        print(f"強制解碼: {decode(text)}")


if __name__ == '__main__':
    main()

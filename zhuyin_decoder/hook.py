#!/usr/bin/env python3
"""
Claude Code user-prompt-submit hook

讀取使用者輸入，偵測注音誤打，若偵測到則透過 additionalContext 注入注音解讀提示。

stdin: {"session_id": ..., "prompt": "使用者輸入", ...}
stdout: {"hookSpecificOutput": {"hookEventName": "UserPromptSubmit", "additionalContext": "..."}}
"""

import json
import sys

from .decoder import detect_and_decode


def main():
    input_data = json.loads(sys.stdin.read())
    prompt = input_data.get("prompt", "")

    result = detect_and_decode(prompt)

    if result["is_zhuyin"]:
        context = result["hint"]
        output = {
            "hookSpecificOutput": {
                "hookEventName": "UserPromptSubmit",
                "additionalContext": context,
            }
        }
        print(json.dumps(output, ensure_ascii=False))
    else:
        # 不修改，輸出空 JSON
        print(json.dumps({}))


if __name__ == "__main__":
    main()

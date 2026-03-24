# elvish-decoder

忘記切換輸入法？沒關係，AI 照樣看得懂。

`elvish-decoder` 是一個 Claude Code hook，自動偵測注音鍵盤誤打（精靈文）並轉換為注音符號，讓 AI 能正確理解你的意思。

## 問題

使用注音輸入法的人一定遇過這個情境：打了一整段才發現忘記切輸入法，送出去的是一串英文亂碼。

```
你想打：你好嗎
實際送出：su3cl3a87
```

AI 看到 `su3cl3a87` 完全無法理解。這個工具解決這個問題。

## 功能

- 自動偵測注音鍵盤誤打（標準注音鍵盤配置）
- 支援混合輸入（一半中文一半忘記切）
- 支援一聲（空白鍵）和二三四聲（聲調鍵 3467）
- 不會誤判正常英文、shell 命令、檔案路徑、程式碼
- Go 編譯，啟動時間約 10ms，幾乎無感

## 安裝

### 前置需求

- [Go](https://go.dev/dl/) 1.26 以上

### 1. 編譯 binary

```bash
git clone https://github.com/ting9093wei/elvish-decoder.git
cd elvish-decoder/go
go build -o elvish-hook .
```

### 2. 放到固定路徑

```bash
# 放在你喜歡的位置，例如：
cp elvish-hook ~/.local/bin/elvish-hook
# 或直接用專案內的路徑
```

### 3. 設定 Claude Code hook

編輯 `~/.claude/settings.json`，在 `hooks` 中加入：

```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/path/to/elvish-hook",
            "timeout": 5,
            "statusMessage": "偵測注音輸入..."
          }
        ]
      }
    ]
  }
}
```

把 `/path/to/elvish-hook` 換成你的實際路徑。

### 4. 重啟 Claude Code 或開啟 `/hooks` 讓設定生效

## 使用範例

設定完成後，正常使用 Claude Code 即可。忘記切輸入法時：

```
你輸入：su3cl3
AI 收到提示：使用者可能忘記切換輸入法。「su3cl3」對應的注音是「ㄋㄧˇ ㄏㄠˇ」
AI 回覆：你好！
```

混合輸入也支援：

```
你輸入：我想要cj04284g4
AI 收到：[中文] 我想要 [注音誤打]「cj04284g4」→「ㄏㄨㄢˋ ㄉㄚˋ ㄕˋ」
```

正常英文輸入不受影響：

```
你輸入：git status
AI 收到：git status（原封不動）
```

## 鍵盤對應表

標準注音鍵盤配置：

```
1ㄅ  2ㄉ  3ˇ  4ˋ  5ㄓ  6ˊ  7˙  8ㄚ  9ㄞ  0ㄢ  -ㄦ
 qㄆ  wㄊ  eㄍ  rㄐ  tㄔ  yㄗ  uㄧ  iㄛ  oㄟ  pㄣ
  aㄇ  sㄋ  dㄎ  fㄑ  gㄕ  hㄘ  jㄨ  kㄜ  lㄠ  ;ㄤ
   zㄈ  xㄌ  cㄏ  vㄒ  bㄖ  nㄙ  mㄩ  ,ㄝ  .ㄡ  /ㄥ
```

聲調：`3`=三聲ˇ、`4`=四聲ˋ、`6`=二聲ˊ、`7`=輕聲˙、空白鍵=一聲

## 偵測策略

1. **有聲調鍵（強信號）**：數字 3/4/6/7 混在字母中 → 高機率是注音誤打
2. **無聲調鍵（一聲）**：多個短 token + 含注音數字鍵 + 非英文 → 可能是注音
3. **排除機制**：常見英文字典、shell 命令、檔案路徑、程式碼語法

## 已知限制

- 僅支援標準注音鍵盤（不支援倉頡、嘸蝦米等）
- 輸出為注音符號，非中文字（Claude 可直接讀懂注音）
- 極短的一聲輸入（如單一 `t`）難以判斷，可能漏偵測
- 偵測基於啟發式規則，非 100% 準確

## 技術細節

- Go 編譯，單一 binary，無外部依賴
- 冷啟動約 10ms（Python 版約 80ms）
- 透過 Claude Code 的 `UserPromptSubmit` hook 整合
- Hook 不修改原始輸入，僅透過 `additionalContext` 注入解讀提示

## License

MIT

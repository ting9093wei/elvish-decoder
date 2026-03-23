package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type HookInput struct {
	Prompt string `json:"prompt"`
}

type HookOutput struct {
	HookSpecificOutput *HookSpecific `json:"hookSpecificOutput,omitempty"`
}

type HookSpecific struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext"`
}

func buildHint(segments []Segment) string {
	var parts []string
	for _, seg := range segments {
		if seg.IsChinese {
			parts = append(parts, fmt.Sprintf("[中文] %s", seg.Text))
		} else if seg.IsZhuyin {
			parts = append(parts, fmt.Sprintf("[注音誤打] 「%s」→「%s」", strings.TrimSpace(seg.Text), seg.Zhuyin))
		} else {
			trimmed := strings.TrimSpace(seg.Text)
			if trimmed != "" {
				parts = append(parts, fmt.Sprintf("[英文] %s", trimmed))
			}
		}
	}

	hint := "使用者可能忘記切換輸入法（混合輸入）。解讀如下：\n"
	for _, p := range parts {
		hint += p + "\n"
	}
	hint += "請根據以上資訊理解使用者意圖並回覆。"
	return hint
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read stdin: %v\n", err)
		os.Exit(1)
	}

	var input HookInput
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Fprintf(os.Stderr, "parse json: %v\n", err)
		os.Exit(1)
	}

	segments, hasZhuyin := analyzeInput(input.Prompt)
	if hasZhuyin {
		hint := buildHint(segments)
		output := HookOutput{
			HookSpecificOutput: &HookSpecific{
				HookEventName:     "UserPromptSubmit",
				AdditionalContext: hint,
			},
		}
		json.NewEncoder(os.Stdout).Encode(output)
	} else {
		fmt.Println("{}")
	}
}

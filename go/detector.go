package main

import (
	"strings"
	"unicode"
)

// 常見 shell 命令（排除誤判）
var shellCommands = map[string]bool{
	"ls": true, "cd": true, "cp": true, "mv": true, "rm": true,
	"cat": true, "pwd": true, "df": true, "du": true, "ps": true,
	"ln": true, "wc": true, "id": true, "su": true, "bg": true,
	"fg": true, "vi": true, "bc": true, "nl": true, "od": true,
}

// 常見英文字
var commonWords = map[string]bool{
	"a": true, "an": true, "the": true, "is": true, "it": true,
	"in": true, "on": true, "at": true, "to": true, "of": true,
	"i": true, "we": true, "he": true, "me": true, "my": true,
	"no": true, "do": true, "go": true, "so": true, "or": true,
	"if": true, "up": true, "as": true, "be": true, "by": true,
	"ok": true, "hi": true, "am": true, "us": true,
	"and": true, "for": true, "but": true, "not": true, "you": true,
	"all": true, "can": true, "had": true, "her": true, "was": true,
	"one": true, "our": true, "out": true, "get": true, "has": true,
	"him": true, "his": true, "how": true, "its": true, "let": true,
	"may": true, "new": true, "now": true, "old": true, "see": true,
	"way": true, "who": true, "did": true, "got": true, "put": true,
	"say": true, "she": true, "too": true, "use": true, "run": true,
	"fix": true, "bug": true, "add": true, "set": true, "try": true,
	"why": true, "yes": true, "yet": true,
	"this": true, "that": true, "with": true, "have": true,
	"from": true, "they": true, "been": true, "call": true,
	"code": true, "file": true, "test": true, "help": true,
	"make": true, "will": true, "what": true, "when": true,
	"your": true, "each": true, "some": true, "then": true,
	"them": true, "look": true, "only": true, "come": true,
	"could": true, "would": true, "should": true, "about": true,
	"after": true, "their": true, "which": true, "there": true,
	"where": true, "other": true, "these": true, "those": true,
	"think": true, "first": true,
	"class": true, "function": true, "return": true, "import": true,
	"print": true, "while": true, "break": true, "continue": true,
	"except": true, "raise": true, "yield": true, "async": true,
	"await": true, "const": true, "export": true, "default": true,
	"delete": true, "typeof": true, "switch": true,
	"git": true, "npm": true, "pip": true, "ssh": true,
	"api": true, "url": true, "cli": true, "env": true,
	"readme": true, "config": true, "setup": true, "index": true,
	"main": true, "init": true, "info": true, "data": true,
	"list": true, "item": true, "name": true, "type": true,
	"value": true, "error": true, "true": true, "false": true,
	"null": true, "none": true, "self": true, "void": true,
	"src": true, "bin": true, "lib": true, "var": true, "tmp": true,
	"log": true, "cmd": true, "pkg": true, "mod": true, "sum": true,
}

// Segment 代表輸入中的一個區段
type Segment struct {
	Text      string
	IsChinese bool
	IsZhuyin  bool
	Zhuyin    string // 解碼後的注音（僅當 IsZhuyin=true）
}

// analyzeInput 分析輸入，拆分中文/非中文區段，偵測注音誤打
// 回傳 segments 和是否包含任何注音誤打
func analyzeInput(text string) ([]Segment, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, false
	}

	// 拆分成中文和非中文區段
	rawSegments := splitByChineseRuns(text)
	hasAnyZhuyin := false

	var segments []Segment
	for _, seg := range rawSegments {
		if seg.IsChinese {
			segments = append(segments, seg)
		} else {
			// 對非中文區段做注音偵測
			nonChinese := strings.TrimSpace(seg.Text)
			if nonChinese == "" {
				segments = append(segments, seg)
				continue
			}
			if isLikelyZhuyinPart(nonChinese) {
				seg.IsZhuyin = true
				seg.Zhuyin = decode(nonChinese)
				hasAnyZhuyin = true
			}
			segments = append(segments, seg)
		}
	}
	return segments, hasAnyZhuyin
}

// splitByChineseRuns 將文字拆分為中文和非中文的連續區段
func splitByChineseRuns(text string) []Segment {
	var segments []Segment
	var current []rune
	inChinese := false

	for _, r := range text {
		isCJK := (r >= 0x4e00 && r <= 0x9fff) ||
			(r >= 0x3400 && r <= 0x4dbf) ||
			(r >= 0xf900 && r <= 0xfaff)
		// 標點符號跟著前一個區段的類型
		isPunct := (r >= 0x3000 && r <= 0x303f) || (r >= 0xff00 && r <= 0xffef)

		if isCJK || (isPunct && inChinese) {
			if !inChinese && len(current) > 0 {
				segments = append(segments, Segment{Text: string(current), IsChinese: false})
				current = nil
			}
			inChinese = true
			current = append(current, r)
		} else {
			if inChinese && len(current) > 0 {
				segments = append(segments, Segment{Text: string(current), IsChinese: true})
				current = nil
			}
			inChinese = false
			current = append(current, r)
		}
	}
	if len(current) > 0 {
		segments = append(segments, Segment{Text: string(current), IsChinese: inChinese})
	}
	return segments
}

// isLikelyZhuyinPart 判斷一段非中文文字是否為注音誤打（原 isLikelyZhuyin 邏輯）
func isLikelyZhuyinPart(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}

	// 像程式碼/路徑 → 跳過
	if looksLikeCodeOrPath(text) {
		return false
	}

	// 像 shell 命令 → 跳過
	if looksLikeShellCommand(text) {
		return false
	}

	tokens := strings.Fields(text)
	if len(tokens) == 0 {
		return false
	}

	hasTone := hasToneDigitMixed(text)

	if hasTone {
		return checkWithTone(tokens)
	}
	return checkNoTone(tokens)
}

// isLikelyZhuyin 保留向後相容（純文字判斷）
func isLikelyZhuyin(text string) bool {
	_, hasZhuyin := analyzeInput(text)
	return hasZhuyin
}

func hasToneDigitMixed(text string) bool {
	lower := strings.ToLower(text)
	for i := 0; i < len(lower)-1; i++ {
		a, b := lower[i], lower[i+1]
		if isLetter(a) && isToneDigit(b) {
			return true
		}
		if isToneDigit(a) && isLetter(b) {
			return true
		}
	}
	return false
}

func checkWithTone(tokens []string) bool {
	zhuyinCount := 0
	for _, token := range tokens {
		if tokenLooksLikeZhuyin(strings.ToLower(token)) {
			zhuyinCount++
		}
	}
	return float64(zhuyinCount)/float64(len(tokens)) >= 0.5
}

func checkNoTone(tokens []string) bool {
	if len(tokens) < 2 {
		return false
	}

	// 平均長度 ≤ 4
	totalLen := 0
	for _, t := range tokens {
		totalLen += len(t)
	}
	if float64(totalLen)/float64(len(tokens)) > 4 {
		return false
	}

	zhuyinCount := 0
	englishCount := 0
	for _, token := range tokens {
		t := strings.ToLower(token)
		if looksLikeEnglish(t) {
			englishCount++
		} else if tokenLooksLikeFirstTone(t) {
			zhuyinCount++
		}
	}

	if float64(englishCount) > float64(len(tokens))*0.3 {
		return false
	}
	return float64(zhuyinCount)/float64(len(tokens)) >= 0.5
}

func tokenLooksLikeZhuyin(token string) bool {
	if len(token) <= 1 {
		return false
	}
	if looksLikeEnglish(token) {
		return false
	}

	hasTone := false
	zhuyinChars := 0
	hasDigitMixed := false

	for i := 0; i < len(token); i++ {
		c := token[i]
		if allZhuyinKeys[c] {
			zhuyinChars++
		}
		if i < len(token)-1 {
			a, b := c, token[i+1]
			if isLetter(a) && isToneDigit(b) || isToneDigit(a) && isLetter(b) {
				hasTone = true
			}
			if isLetter(a) && isDigit(b) || isDigit(a) && isLetter(b) {
				hasDigitMixed = true
			}
		}
	}

	coverage := float64(zhuyinChars) / float64(len(token))

	if hasTone && coverage > 0.7 {
		return true
	}
	if hasDigitMixed && coverage > 0.8 {
		return true
	}
	// 短 token + 高覆蓋率 + 非純字母
	if len(token) <= 4 && coverage >= 0.9 && !isAllAlpha(token) {
		return true
	}
	return false
}

func tokenLooksLikeFirstTone(token string) bool {
	if len(token) > 5 || len(token) < 1 {
		return false
	}

	// 含注音數字鍵 (1,2,5,8,9,0) 混字母
	hasZhuyinDigit := false
	zhuyinChars := 0
	for i := 0; i < len(token); i++ {
		c := token[i]
		if allZhuyinKeys[c] {
			zhuyinChars++
		}
		if i < len(token)-1 {
			a, b := c, token[i+1]
			if (isLetter(a) && zhuyinDigitKeys[b]) || (zhuyinDigitKeys[a] && isLetter(b)) {
				hasZhuyinDigit = true
			}
		}
	}

	coverage := float64(zhuyinChars) / float64(len(token))

	if hasZhuyinDigit && coverage >= 0.8 {
		return true
	}
	// 純字母 1-2 字元且非英文
	if isAllAlpha(token) && len(token) <= 2 && !looksLikeEnglish(token) {
		return true
	}
	return false
}

func looksLikeCodeOrPath(text string) bool {
	// 檔案路徑 (.xxx 副檔名)
	for i := 0; i < len(text)-1; i++ {
		if text[i] == '.' && i+1 < len(text) && isLetter(text[i+1]) {
			// 有 .ext 的 pattern
			extLen := 0
			for j := i + 1; j < len(text) && isLetter(text[j]); j++ {
				extLen++
			}
			if extLen >= 1 && extLen <= 5 {
				return true
			}
		}
	}
	// URL
	if strings.Contains(text, "http://") || strings.Contains(text, "https://") {
		return true
	}
	// 程式碼符號
	for _, r := range text {
		switch r {
		case '(', ')', '{', '}', '[', ']', '=', '<', '>':
			return true
		}
	}
	return false
}

func looksLikeShellCommand(text string) bool {
	tokens := strings.Fields(strings.ToLower(text))
	if len(tokens) == 0 {
		return false
	}
	// 第一個 token 是 shell 命令
	if shellCommands[tokens[0]] {
		return true
	}
	// 帶 flag (--xxx 或 -x)
	for _, t := range tokens {
		if strings.HasPrefix(t, "--") || (strings.HasPrefix(t, "-") && len(t) <= 3 && len(t) >= 2) {
			return true
		}
	}
	return false
}

func looksLikeEnglish(token string) bool {
	clean := strings.NewReplacer(".", "", "_", "", "-", "", "/", "").Replace(token)

	if commonWords[clean] || commonWords[token] {
		return true
	}
	// 純字母、有母音、≥ 3 字元
	if isAllAlpha(clean) && len(clean) >= 3 {
		for _, c := range clean {
			if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
				return true
			}
		}
	}
	return false
}

// helpers
func isLetter(c byte) bool  { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isDigit(c byte) bool   { return c >= '0' && c <= '9' }
func isToneDigit(c byte) bool { return c == '3' || c == '4' || c == '6' || c == '7' }

func isAllAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

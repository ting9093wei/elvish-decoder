package main

import (
	"strings"
)

// 注音分類
var initials = map[string]bool{
	"ㄅ": true, "ㄆ": true, "ㄇ": true, "ㄈ": true,
	"ㄉ": true, "ㄊ": true, "ㄋ": true, "ㄌ": true,
	"ㄍ": true, "ㄎ": true, "ㄏ": true,
	"ㄐ": true, "ㄑ": true, "ㄒ": true,
	"ㄓ": true, "ㄔ": true, "ㄕ": true, "ㄖ": true,
	"ㄗ": true, "ㄘ": true, "ㄙ": true,
}

var medials = map[string]bool{
	"ㄧ": true, "ㄨ": true, "ㄩ": true,
}

var finals = map[string]bool{
	"ㄚ": true, "ㄛ": true, "ㄜ": true, "ㄝ": true,
	"ㄞ": true, "ㄟ": true, "ㄠ": true, "ㄡ": true,
	"ㄢ": true, "ㄣ": true, "ㄤ": true, "ㄥ": true, "ㄦ": true,
}

var tones = map[string]bool{
	"ˊ": true, "ˇ": true, "ˋ": true, "˙": true,
}

func decode(text string) string {
	syllables := parseSyllables(text)
	return strings.Join(syllables, " ")
}

func parseSyllables(text string) []string {
	symbols := keysToZhuyinSymbols(text)
	var syllables []string
	var current []string

	for _, sym := range symbols {
		if tones[sym] {
			current = append(current, sym)
			syllables = append(syllables, strings.Join(current, ""))
			current = nil
		} else if sym == " " {
			if len(current) > 0 {
				syllables = append(syllables, strings.Join(current, ""))
				current = nil
			}
		} else if initials[sym] || medials[sym] || finals[sym] {
			if len(current) > 0 && shouldSplit(current, sym) {
				syllables = append(syllables, strings.Join(current, ""))
				current = []string{sym}
			} else {
				current = append(current, sym)
			}
		} else {
			if len(current) > 0 {
				syllables = append(syllables, strings.Join(current, ""))
				current = nil
			}
			syllables = append(syllables, sym)
		}
	}

	if len(current) > 0 {
		syllables = append(syllables, strings.Join(current, ""))
	}
	return syllables
}

func keysToZhuyinSymbols(text string) []string {
	var result []string
	lower := strings.ToLower(text)
	for i := 0; i < len(lower); i++ {
		c := lower[i]
		if z, ok := keyToZhuyin[c]; ok {
			result = append(result, z)
		} else if t, ok := toneKeys[c]; ok {
			result = append(result, t)
		} else {
			result = append(result, string(c))
		}
	}
	return result
}

func shouldSplit(current []string, nextSym string) bool {
	hasInitial := false
	hasMedial := false
	hasFinal := false
	for _, s := range current {
		if initials[s] {
			hasInitial = true
		}
		if medials[s] {
			hasMedial = true
		}
		if finals[s] {
			hasFinal = true
		}
	}

	if initials[nextSym] {
		if hasFinal || hasMedial || hasInitial {
			return true
		}
	}
	if medials[nextSym] {
		if hasFinal {
			return true
		}
		last := current[len(current)-1]
		if medials[last] {
			return true
		}
	}
	return false
}

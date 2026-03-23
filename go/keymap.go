package main

// 英文按鍵 → 注音符號
var keyToZhuyin = map[byte]string{
	// 數字列
	'1': "ㄅ", '2': "ㄉ", '5': "ㄓ",
	'8': "ㄚ", '9': "ㄞ", '0': "ㄢ", '-': "ㄦ",
	// QWERTY 列
	'q': "ㄆ", 'w': "ㄊ", 'e': "ㄍ", 'r': "ㄐ", 't': "ㄔ",
	'y': "ㄗ", 'u': "ㄧ", 'i': "ㄛ", 'o': "ㄟ", 'p': "ㄣ",
	// ASDF 列
	'a': "ㄇ", 's': "ㄋ", 'd': "ㄎ", 'f': "ㄑ", 'g': "ㄕ",
	'h': "ㄘ", 'j': "ㄨ", 'k': "ㄜ", 'l': "ㄠ", ';': "ㄤ",
	// ZXCV 列
	'z': "ㄈ", 'x': "ㄌ", 'c': "ㄏ", 'v': "ㄒ", 'b': "ㄖ",
	'n': "ㄙ", 'm': "ㄩ", ',': "ㄝ", '.': "ㄡ", '/': "ㄥ",
}

// 聲調按鍵
var toneKeys = map[byte]string{
	'3': "ˇ", '4': "ˋ", '6': "ˊ", '7': "˙",
}

// 所有注音相關按鍵（含聲調）
var allZhuyinKeys map[byte]bool

// 注音數字鍵（非聲調）: 1,2,5,8,9,0
var zhuyinDigitKeys = map[byte]bool{
	'1': true, '2': true, '5': true, '8': true, '9': true, '0': true,
}

func init() {
	allZhuyinKeys = make(map[byte]bool)
	for k := range keyToZhuyin {
		allZhuyinKeys[k] = true
	}
	for k := range toneKeys {
		allZhuyinKeys[k] = true
	}
}

"""偵測輸入是否為注音鍵盤誤打"""

import re
from .keymap import ALL_ZHUYIN_KEYS, INITIALS, MEDIALS, FINALS, KEY_TO_ZHUYIN, TONE_KEYS

# 預編譯 regex
_RE_CHINESE = re.compile(r'[\u4e00-\u9fff]')
_RE_TONE_MIXED = re.compile(r'[a-zA-Z][3467]|[3467][a-zA-Z]')
_RE_DIGIT_MIXED = re.compile(r'[a-z][0-9]|[0-9][a-z]')
_RE_ZHUYIN_DIGIT_MIXED = re.compile(r'[a-z][125890]|[125890][a-z]')
_RE_FILE_PATH = re.compile(r'[./\\][\w]+\.[a-zA-Z]{1,5}\b')
_RE_URL = re.compile(r'https?://')
_RE_CODE = re.compile(r'[(){}\[\]=<>]')


def is_likely_zhuyin(text: str) -> bool:
    """
    判斷輸入文字是否可能是注音鍵盤誤打。

    兩種模式：
    1. 有聲調鍵 (3467) → 強信號，寬鬆判定
    2. 無聲調鍵（全一聲）→ 看 token 結構：短、含注音數字鍵、非英文
    """
    text = text.strip()
    if not text:
        return False

    # 已包含中文
    if _RE_CHINESE.search(text):
        return False

    # 檔案路徑、URL、程式碼
    if _looks_like_code_or_path(text):
        return False

    tokens = text.split()
    if not tokens:
        return False

    has_tone = bool(_RE_TONE_MIXED.search(text))

    if has_tone:
        return _check_with_tone(tokens)
    else:
        return _check_no_tone(tokens)


def _check_with_tone(tokens: list[str]) -> bool:
    """有聲調鍵時的判定：只要足夠比例的 token 像注音就過"""
    zhuyin_count = 0
    for token in tokens:
        if _token_looks_like_zhuyin(token.lower(), require_tone=False):
            zhuyin_count += 1

    return zhuyin_count / len(tokens) >= 0.5


def _check_no_tone(tokens: list[str]) -> bool:
    """
    無聲調鍵（全一聲）的判定，更嚴格：
    - token 都很短（平均 ≤ 3）
    - 多數 token 含注音數字鍵 (1,2,5,8,9,0) 或不像英文
    - 幾乎沒有能被辨識為英文的 token
    """
    if len(tokens) < 2:
        return False

    total_len = sum(len(t) for t in tokens)
    avg_len = total_len / len(tokens)

    # 注音一聲 token 通常 1-4 字元
    if avg_len > 4:
        return False

    zhuyin_count = 0
    english_count = 0

    for token in tokens:
        t = token.lower()
        if _looks_like_english(t):
            english_count += 1
        elif _token_looks_like_first_tone(t):
            zhuyin_count += 1

    # 英文太多 → 不是注音
    if english_count > len(tokens) * 0.3:
        return False

    return zhuyin_count / len(tokens) >= 0.5


def _token_looks_like_first_tone(token: str) -> bool:
    """判斷 token 是否像一聲注音（無聲調鍵）"""
    if len(token) > 5 or len(token) < 1:
        return False

    # 含注音數字鍵 (1,2,5,8,9,0) 混在字母中 → 信號
    has_zhuyin_digit = bool(_RE_ZHUYIN_DIGIT_MIXED.search(token))

    # 覆蓋率
    zhuyin_chars = sum(1 for c in token if c in ALL_ZHUYIN_KEYS)
    coverage = zhuyin_chars / len(token)

    if has_zhuyin_digit and coverage >= 0.8:
        return True

    # 純字母但非常短(1-2) 且不像英文 → 可能是注音
    if token.isalpha() and len(token) <= 2 and not _looks_like_english(token):
        return True

    return False


def _looks_like_code_or_path(text: str) -> bool:
    """排除明顯是程式碼、檔案路徑、URL 的輸入"""
    if _RE_FILE_PATH.search(text):
        return True
    if _RE_URL.search(text):
        return True
    if _RE_CODE.search(text):
        return True
    return False


def _token_looks_like_zhuyin(token: str, require_tone: bool = False) -> bool:
    """判斷單一 token 是否像注音誤打"""
    if len(token) <= 1:
        return False

    if _looks_like_english(token):
        return False

    has_tone_digit = bool(_RE_TONE_MIXED.search(token))

    # 覆蓋率
    zhuyin_chars = sum(1 for c in token if c in ALL_ZHUYIN_KEYS)
    coverage = zhuyin_chars / len(token) if token else 0

    if has_tone_digit and coverage > 0.7:
        return True

    # 有數字混字母（含注音數字鍵）
    has_digit_mixed = bool(_RE_DIGIT_MIXED.search(token))
    if has_digit_mixed and coverage > 0.8:
        return True

    # 短 token + 高覆蓋率 + 非英文（可能是一聲 token 混在有聲調的句子中）
    if not require_tone and len(token) <= 4 and coverage >= 0.9 and not token.isalpha():
        return True

    return False


# 常見英文短字
_COMMON_WORDS = frozenset({
    'a', 'an', 'the', 'is', 'it', 'in', 'on', 'at', 'to', 'of',
    'i', 'we', 'he', 'me', 'my', 'no', 'do', 'go', 'so', 'or',
    'if', 'up', 'as', 'be', 'by', 'ok', 'hi', 'am', 'us',
    'and', 'for', 'but', 'not', 'you', 'all', 'can', 'had', 'her',
    'was', 'one', 'our', 'out', 'get', 'has', 'him', 'his', 'how',
    'its', 'let', 'may', 'new', 'now', 'old', 'see', 'way', 'who',
    'did', 'got', 'put', 'say', 'she', 'too', 'use', 'run', 'fix',
    'bug', 'add', 'set', 'try', 'why', 'yes', 'yet',
    'this', 'that', 'with', 'have', 'from', 'they', 'been', 'call',
    'code', 'file', 'test', 'help', 'make', 'will', 'what', 'when',
    'your', 'each', 'some', 'then', 'them', 'look', 'only', 'come',
    'could', 'would', 'should', 'about', 'after', 'their', 'which',
    'there', 'where', 'other', 'these', 'those', 'think', 'first',
    'class', 'function', 'return', 'import', 'print', 'while',
    'break', 'continue', 'except', 'raise', 'yield', 'async', 'await',
    'const', 'export', 'default', 'delete', 'typeof', 'switch',
    'git', 'npm', 'pip', 'ssh', 'api', 'url', 'cli', 'env',
    'readme', 'config', 'setup', 'index', 'main', 'init', 'info',
    'data', 'list', 'item', 'name', 'type', 'value', 'error',
    'true', 'false', 'null', 'none', 'self', 'void',
})


def _looks_like_english(token: str) -> bool:
    """粗略判斷是否像正常英文"""
    clean = token.replace('.', '').replace('_', '').replace('-', '').replace('/', '')

    if clean in _COMMON_WORDS or token in _COMMON_WORDS:
        return True

    # 純字母、有母音、長度 ≥ 3 → 可能是英文
    if clean.isalpha() and len(clean) >= 3:
        if any(c in 'aeiou' for c in clean):
            return True

    return False

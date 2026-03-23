"""將英文按鍵序列解析為注音音節"""

from .keymap import KEY_TO_ZHUYIN, TONE_KEYS, INITIALS, MEDIALS, FINALS


def keys_to_zhuyin_symbols(text: str) -> list[str]:
    """將每個英文字元轉換為對應的注音符號或聲調，無法對應的保持原樣"""
    result = []
    for ch in text.lower():
        if ch in KEY_TO_ZHUYIN:
            result.append(KEY_TO_ZHUYIN[ch])
        elif ch in TONE_KEYS:
            result.append(TONE_KEYS[ch])
        else:
            result.append(ch)
    return result


def parse_syllables(text: str) -> list[str]:
    """
    將英文按鍵序列解析成注音音節列表。

    注音輸入規則：
    - 每個音節由 [聲母] + [介音] + [韻母] + 聲調 組成
    - 聲調鍵 (3,4,6,7) 或空白鍵（一聲）結束一個音節
    - 空白鍵同時也可能是詞與詞之間的分隔

    回傳：音節字串列表，例如 ['ㄋㄧˇ', 'ㄏㄠˇ']
    """
    symbols = keys_to_zhuyin_symbols(text)
    syllables = []
    current = []

    for sym in symbols:
        if sym in ('ˊ', 'ˇ', 'ˋ', '˙'):
            # 聲調 = 音節結束
            current.append(sym)
            syllables.append(''.join(current))
            current = []
        elif sym == ' ':
            # 空白 = 一聲結束，或者單純分隔
            if current:
                syllables.append(''.join(current))
                current = []
        elif sym in INITIALS | MEDIALS | FINALS:
            # 如果當前音節已有韻母，且新符號是聲母，代表新音節開始（一聲省略）
            if current and _should_split(current, sym):
                syllables.append(''.join(current))
                current = [sym]
            else:
                current.append(sym)
        else:
            # 非注音字元，先結束當前音節
            if current:
                syllables.append(''.join(current))
                current = []
            syllables.append(sym)

    if current:
        syllables.append(''.join(current))

    return syllables


def _should_split(current: list[str], next_sym: str) -> bool:
    """判斷是否應該在此處切分音節（處理一聲省略的情況）"""
    if not current:
        return False

    last = current[-1]

    # 新的聲母出現，且前面已有韻母或介音 → 切分
    if next_sym in INITIALS:
        # 前面有韻母 → 一定切
        if any(s in FINALS for s in current):
            return True
        # 前面有介音且新的是聲母 → 切（例如 ㄧ 單獨成音後接聲母）
        if any(s in MEDIALS for s in current) and len(current) >= 1:
            return True
        # 前面已有聲母 → 切（不會有兩個聲母連續）
        if any(s in INITIALS for s in current):
            return True

    # 新的介音，但前面已有韻母 → 切
    if next_sym in MEDIALS and any(s in FINALS for s in current):
        return True

    # 新的介音，前面已有介音 → 看情況
    if next_sym in MEDIALS and any(s in MEDIALS for s in current):
        # ㄧㄨ, ㄧㄩ 等不合法，要切
        if last in MEDIALS:
            return True

    return False


def syllables_to_string(syllables: list[str]) -> str:
    """將音節列表組合為可讀的注音字串"""
    return ' '.join(syllables)

"""主要解碼邏輯：英文按鍵 → 注音 → 可讀輸出"""

from .parser import parse_syllables, syllables_to_string
from .detector import is_likely_zhuyin


def decode(text: str) -> str:
    """
    將注音鍵盤誤打的英文轉換為注音符號字串。

    例如: "su3cl3" → "ㄋㄧˇ ㄏㄠˇ"
    """
    syllables = parse_syllables(text)
    return syllables_to_string(syllables)


def detect_and_decode(text: str) -> dict:
    """
    偵測並解碼輸入。

    回傳:
        {
            "is_zhuyin": bool,       # 是否偵測為注音誤打
            "original": str,         # 原始輸入
            "zhuyin": str | None,    # 轉換後的注音（若偵測為誤打）
            "hint": str | None,      # 給 AI 的提示訊息
        }
    """
    if is_likely_zhuyin(text):
        zhuyin = decode(text)
        return {
            "is_zhuyin": True,
            "original": text,
            "zhuyin": zhuyin,
            "hint": f"使用者可能忘記切換輸入法。原始輸入「{text}」對應的注音是「{zhuyin}」，請以注音解讀回覆。",
        }
    return {
        "is_zhuyin": False,
        "original": text,
        "zhuyin": None,
        "hint": None,
    }

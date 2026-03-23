"""zhuyin-decoder: 注音鍵盤誤打偵測與轉換"""

from .decoder import decode, detect_and_decode
from .detector import is_likely_zhuyin
from .parser import parse_syllables, syllables_to_string

__all__ = [
    'decode',
    'detect_and_decode',
    'is_likely_zhuyin',
    'parse_syllables',
    'syllables_to_string',
]

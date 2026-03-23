"""測試注音解碼器"""

import sys
sys.path.insert(0, '.')

from zhuyin_decoder.keymap import KEY_TO_ZHUYIN, TONE_KEYS
from zhuyin_decoder.parser import keys_to_zhuyin_symbols, parse_syllables
from zhuyin_decoder.detector import is_likely_zhuyin
from zhuyin_decoder.decoder import decode, detect_and_decode


def test_key_mapping():
    """測試基本按鍵對應"""
    assert KEY_TO_ZHUYIN['s'] == 'ㄋ'
    assert KEY_TO_ZHUYIN['u'] == 'ㄧ'
    assert TONE_KEYS['3'] == 'ˇ'
    assert KEY_TO_ZHUYIN['c'] == 'ㄏ'
    assert KEY_TO_ZHUYIN['l'] == 'ㄠ'
    print("✓ key_mapping")


def test_keys_to_symbols():
    """測試字元轉注音符號"""
    symbols = keys_to_zhuyin_symbols('su3')
    assert symbols == ['ㄋ', 'ㄧ', 'ˇ']
    print("✓ keys_to_symbols")


def test_parse_syllables():
    """測試音節解析"""
    # 你好 = su3cl3
    syllables = parse_syllables('su3cl3')
    assert syllables == ['ㄋㄧˇ', 'ㄏㄠˇ']

    # 我 = ji3
    syllables = parse_syllables('ji3')
    assert syllables == ['ㄨㄛˇ']

    # 謝謝 = vu,4vu,4 (v=ㄒ u=ㄧ ,=ㄝ 4=ˋ)
    syllables = parse_syllables('vu,4vu,4')
    assert syllables == ['ㄒㄧㄝˋ', 'ㄒㄧㄝˋ'], f"got {syllables}"

    print("✓ parse_syllables")


def test_decode():
    """測試完整解碼"""
    # 你好
    result = decode('su3cl3')
    assert result == 'ㄋㄧˇ ㄏㄠˇ', f"got: {result}"

    # 我
    result = decode('ji3')
    assert result == 'ㄨㄛˇ', f"got: {result}"

    print("✓ decode")


def test_detect_not_zhuyin():
    """測試正常英文不會被誤判"""
    assert not is_likely_zhuyin('hello world')
    assert not is_likely_zhuyin('please fix this bug')
    assert not is_likely_zhuyin('git status')
    assert not is_likely_zhuyin('')
    assert not is_likely_zhuyin('你好')
    print("✓ detect_not_zhuyin")


def test_detect_zhuyin():
    """測試注音誤打能被偵測"""
    assert is_likely_zhuyin('su3cl3')        # 你好
    assert is_likely_zhuyin('ji3g4284')      # 我是誰 (rough)
    print("✓ detect_zhuyin")


def test_detect_and_decode():
    """測試偵測 + 解碼整合"""
    result = detect_and_decode('su3cl3')
    assert result['is_zhuyin'] is True
    assert result['zhuyin'] == 'ㄋㄧˇ ㄏㄠˇ'
    assert result['hint'] is not None

    result = detect_and_decode('hello')
    assert result['is_zhuyin'] is False
    print("✓ detect_and_decode")


if __name__ == '__main__':
    test_key_mapping()
    test_keys_to_symbols()
    test_parse_syllables()
    test_decode()
    test_detect_not_zhuyin()
    test_detect_zhuyin()
    test_detect_and_decode()
    print("\n全部測試通過！")

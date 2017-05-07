package strings

import (
    "unicode/utf16"
)

func utf16BeEncode(table *codePageTableInfo, str String) (bytes []byte) {
    for _, r := range utf16.Encode([]rune(string(str))) {
        bytes = append(bytes, byte(r >> 8))
        bytes = append(bytes, byte(r))
    }

    return
}

func utf16BeDecode(table *codePageTableInfo, bytes []byte) String {
    runes := []rune{}
    for _, r := range bytesToUInt16Array(bytes) {
        runes = append(runes, rune((r >> 8) | (r << 8)))
    }

    return String(string(runes))
}

func init() {
    cptable[CP_UTF16_BE] = codePageTableInfo{
                        CodePage    : uint16(CP_UTF16_BE),
                        initialized : true,
                        encode      : utf16BeEncode,
                        decode      : utf16BeDecode,
                    }
}

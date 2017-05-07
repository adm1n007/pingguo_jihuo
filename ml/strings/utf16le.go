package strings

import (
    "unicode/utf16"
)

func utf16LeEncode(table *codePageTableInfo, str String) (bytes []byte) {
    for _, r := range utf16.Encode([]rune(string(str))) {
        bytes = append(bytes, byte(r))
        bytes = append(bytes, byte(r >> 8))
    }

    return
}

func utf16LeDecode(table *codePageTableInfo, bytes []byte) String {
    runes := []rune{}
    for _, r := range bytesToUInt16Array(bytes) {
        runes = append(runes, rune(r))
    }

    return String(string(runes))
}

func init() {
    cptable[CP_UTF16_LE] = codePageTableInfo{
                        CodePage    : uint16(CP_UTF16_LE),
                        initialized : true,
                        encode      : utf16LeEncode,
                        decode      : utf16LeDecode,
                    }
}

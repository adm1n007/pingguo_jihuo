package strings

func utf8Encode(table *codePageTableInfo, str String) (bytes []byte) {
    for i := 0; i != len(str); i++ {
        bytes = append(bytes, str[i])
    }

    return
}

func utf8Decode(table *codePageTableInfo, bytes []byte) String {
    return String(string(bytes))
}

func init() {
    cptable[CP_UTF8] = codePageTableInfo{
                        CodePage    : uint16(CP_UTF8),
                        initialized : true,
                        encode      : utf8Encode,
                        decode      : utf8Decode,
                    }
}

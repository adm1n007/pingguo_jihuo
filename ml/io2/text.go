package io2

import (
    // "os"
    "bytes"
    // "io/ioutil"
    "ml/io2/filestream"
    "ml/strings"
    "unicode/utf8"
)

func isAnsiAsUtf8(buf []byte) bool {
    offset := 0
    for offset < len(buf) {
        r, size := utf8.DecodeRune(buf[offset:])
        if r == utf8.RuneError {
            return false
        }

        offset += size
    }

    return true
}

func ReadLines(filename strings.String) ([]strings.String) {
    // file, err := os.Open(string(filename))
    // if err != nil {
    //     return nil, err
    // }

    // defer file.Close()

    // buf, err := ioutil.ReadAll(file)
    // if err != nil {
    //     return nil, err
    // }

    file := filestream.Open(filename.String())
    defer file.Close()

    buf := file.ReadAll()

    var text strings.String

    switch {
        case bytes.Equal(buf[:3], strings.BOM_UTF8):
            text = strings.Decode(buf[3:], strings.CP_UTF8)

        case bytes.Equal(buf[:2], strings.BOM_UTF16_LE):
            text = strings.Decode(buf[2:], strings.CP_UTF16_LE)

        case bytes.Equal(buf[:2], strings.BOM_UTF16_BE):
            text = strings.Decode(buf[2:], strings.CP_UTF16_BE)

        default:
            if isAnsiAsUtf8(buf) {
                text = strings.Decode(buf, strings.CP_UTF8)
            } else {
                text = strings.Decode(buf, strings.CP_GBK)
            }
    }

    return text.SplitLines()
}

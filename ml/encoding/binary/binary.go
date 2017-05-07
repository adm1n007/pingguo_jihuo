package binary

import (
    . "active_apple/ml/trace"
    "encoding/binary"
    "bytes"
)

func IntToBytes(data interface{}, length int, byteOrder binary.ByteOrder) []byte {
    if length == 0 {
        Raisef("invalid length")
    }

    switch i := data.(type) {
        case int:
            data = int64(i)

        case uint:
            data = uint64(i)
    }

    w := bytes.NewBuffer(nil)
    err := binary.Write(w, byteOrder, data)
    RaiseIf(err)

    b := w.Bytes()

    switch byteOrder {
        case LittleEndian:
            index := 0
            for index != len(b) && b[index] != 0 {
                index++
            }

            b = b[:index]

        case BigEndian:
            index := 0
            for index != len(b) && b[index] == 0 {
                index++
            }

            b = b[index:]
    }

    if len(b) >= length {
        return b
    }

    buf := make([]byte, length)

    switch byteOrder {
        case LittleEndian:
            copy(buf, b)

        case BigEndian:
            copy(buf[length - len(b):], b)
    }

    return buf
}

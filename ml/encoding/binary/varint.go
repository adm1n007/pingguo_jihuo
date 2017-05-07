package binary

import (
    . "active_apple/ml/trace"
    "encoding/binary"
)

func ToVarint(data interface{}) []byte {
    buf := [binary.MaxVarintLen64]byte{}

    var length int

    switch v := data.(type) {
        case int8:
            length = binary.PutUvarint(buf[:], uint64(v))

        case int16:
            length = binary.PutUvarint(buf[:], uint64(v))

        case int32:
            length = binary.PutUvarint(buf[:], uint64(v))

        case int64:
            length = binary.PutUvarint(buf[:], uint64(v))

        case int:
            length = binary.PutUvarint(buf[:], uint64(v))

        case uint8:
            length = binary.PutUvarint(buf[:], uint64(v))

        case uint16:
            length = binary.PutUvarint(buf[:], uint64(v))

        case uint32:
            length = binary.PutUvarint(buf[:], uint64(v))

        case uint64:
            length = binary.PutUvarint(buf[:], uint64(v))

        case uint:
            length = binary.PutUvarint(buf[:], uint64(v))

        default:
            Raisef("UnsupportedType: %T", v)
    }

    return buf[:length]
}

func Varint(buf []byte) (uint64, int) {
    return binary.Uvarint(buf)
}

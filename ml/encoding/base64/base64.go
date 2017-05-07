package base64

import (
    . "ml/strings"
    "encoding/base64"
)

func EncodeToString(data []byte) String {
    return String(base64.StdEncoding.EncodeToString(data))
}

func DecodeString(s string) []byte {
    data, err := base64.StdEncoding.DecodeString(s)
    raiseBase64Error(err)
    return data
}

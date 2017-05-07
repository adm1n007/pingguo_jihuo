package ituneslib

import (
    . "ml/trace"
)

func EncryptJsSpToken(token []byte) []byte {
    t := make([]byte, len(token))
    copy(t, token)

    st, _, _ := itunes.EncryptJsSpToken.Call(2, bytesPtr(t))
    status := getStatus(st)
    if status != 0 {
        Raise(newiTunesHelperErrorf("EncryptJsSpToken return %X", uint32(status)))
    }

    return t
}

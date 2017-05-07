package packet

import (
    . "ml/trace"
)

type PacketError struct {
    *BaseException
}

func RaisePacketError(err error) {
    if err == nil {
        return
    }

    Raise(NewPacketError(err.Error()))
}

func NewPacketError(format string, args ...interface{}) *PacketError {
    return &PacketError{BaseException: NewBaseException(format, args...)}
}

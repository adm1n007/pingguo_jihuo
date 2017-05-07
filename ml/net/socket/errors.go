package socket

import (
    . "active_apple/ml/trace"
    "net"
)

type SocketError struct {
    *BaseException
}

func RaiseSocketError(err error) {
    if err == nil {
        return
    }

    if e, ok := err.(*net.OpError); ok {
        if e.Timeout() {
            Raise(NewTimeoutError(e.Error()))
        }
    }

    Raise(NewSocketError(err.Error()))
}

func NewSocketError(msg string) *SocketError {
    return &SocketError{
        BaseException: NewBaseException(msg),
    }
}


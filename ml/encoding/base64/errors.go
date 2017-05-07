package base64

import (
    . "ml/trace"
)

type Base64Error struct {
    *BaseException
}

func raiseBase64Error(err error) {
    if err == nil {
        return
    }

    Raise(NewBase64Error(err.Error()))
}

func NewBase64Error(format string, args ...interface{}) *Base64Error {
    return &Base64Error{
        BaseException: NewBaseException(format, args...),
    }
}

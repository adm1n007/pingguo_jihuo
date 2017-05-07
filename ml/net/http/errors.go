package http

import (
    . "fmt"
    . "ml/strings"
    . "ml/trace"
    "errors"
)

const (
    HTTP_ERROR_GENERIC              = 0
    HTTP_ERROR_TIMEOUT              = 1
    HTTP_ERROR_CONNECT_PROXY        = 2
    HTTP_ERROR_INVALID_RESPONSE     = 3
    HTTP_ERROR_TOO_MANY_REDIRECT    = 4
    HTTP_ERROR_CANNOT_CONNECT       = 5
    HTTP_ERROR_READ_ERROR           = 6
    HTTP_ERROR_RESPONSE_ERROR       = 7
    HTTP_ERROR_BAD_GATE_WAY         = 8
)

type HttpError struct {
    Op          string
    URL         string
    Err         error
    Type        int
}

func NewHttpError(t int, op, url, msg String) *HttpError {
    return &HttpError{string(op), string(url), errors.New(string(msg)), t}
}

func (self *HttpError) Error() string {
    return Sprintf("%v %v: (%v) %v", self.Op, self.URL, self.Type, self.Err.Error())
}

func RaiseHttpError(err error) {
    if err == nil {
        return
    }

    Raise(NewHttpError(HTTP_ERROR_GENERIC, "", "", String(err.Error())))
}

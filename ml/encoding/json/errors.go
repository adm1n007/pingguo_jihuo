package json

import (
    . "ml/trace"
)

type JSONDecodeError struct {
    *BaseException
}

func NewJSONDecodeError(msg string) *JSONDecodeError {
    return &JSONDecodeError{
        BaseException: NewBaseException(msg),
    }
}
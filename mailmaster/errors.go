package mailmaster

import (
    . "active_apple/ml/trace"
)

type LoginError struct {
    *BaseException
}

type PasswordError struct {
    *BaseException
}

type AccountError struct {
    *BaseException
}

func NewLoginError(format string, args ...interface{}) *LoginError {
    return &LoginError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewPasswordError(format string, args ...interface{}) *PasswordError {
    return &PasswordError{
        BaseException: NewBaseException(format, args...),
    }
}

func NewAccountError(format string, args ...interface{}) *AccountError {
    return &AccountError{
        BaseException: NewBaseException(format, args...),
    }
}

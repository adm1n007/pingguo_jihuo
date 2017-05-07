package rsa

import (
)

type BaseRSAError struct {
    Message string
}

type RSAInvalidKeyError struct {
    BaseRSAError
}

func (self *BaseRSAError) String() string {
    return self.Message
}

func (self *BaseRSAError) Error() string {
    return self.Message
}

func NewRSAInvalidKeyError(msg string) *RSAInvalidKeyError {
    return &RSAInvalidKeyError{
        BaseRSAError: BaseRSAError{
            msg,
        },
    }
}

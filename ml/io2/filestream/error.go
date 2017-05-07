package filestream

import (
    . "active_apple/ml/trace"
)

const (
    BUILTIN_FILE_ERROR   = 0
)

type FileErrorType int

func (self FileError) String() string {
    return "BUILTIN_FILE_ERROR"
}

type FileError struct {
    error
    Type FileErrorType
}

func raiseFileError(err error) {
    if err == nil {
        return
    }

    Raise(&FileError{err, BUILTIN_FILE_ERROR})
}

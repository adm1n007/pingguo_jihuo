package filestream

import (
    . "active_apple/ml/trace"
)

func raiseGenericError(err error) {
    if err == nil {
        return
    }

    Raise(NewFileGenericError(err.Error()))
}

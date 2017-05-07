package trace


type ErrorWrapper struct {
    *BaseException
    cause error
}

func NewWrapper(cause error, format string, args ...interface{}) *ErrorWrapper {
    return &ErrorWrapper{
        BaseException   : newBaseException(3, format, args...),
        cause           : cause,
    }
}

func (self *ErrorWrapper) Cause() error {
    return self.cause
}

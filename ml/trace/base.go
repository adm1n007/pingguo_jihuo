package trace

import (
    "fmt"
    "strings"
)

type baseExceptionInterface interface {
    GetTraceBack() []string
    GetTraceBackString() string
    getMessage() string
}

type callStackEntry struct {
    pc      uintptr
    file    string
    line    int
}

type BaseException struct {
    Message string
    callStack [callerDepth]callStackEntry
}

func NewBaseException(format string, args ...interface{}) *BaseException {
    return newBaseException(3, format, args...)
}

func NewBaseExceptionWithSkip(skip int, format string, args ...interface{}) *BaseException {
    return newBaseException(skip + 3, format, args...)
}

func newBaseException(skip int, format string, args ...interface{}) *BaseException {
    e := &BaseException{
        Message: fmt.Sprintf(format, args...),
    }

    getCallStack(skip, e.callStack[:])

    return e
}

func (self *BaseException) String() string {
    return "(traceback)\n" + self.GetTraceBackString() + "\n" + self.Message
}

func (self *BaseException) Error() string {
    return self.String()
}

func (self *BaseException) GetTraceBack() []string {
    return getTrackBack(self.callStack[:])
}

func (self *BaseException) GetTraceBackString() string {
    return strings.Join(self.GetTraceBack(), "\r\n")
}

func (self *BaseException) getMessage() string {
    return self.Message
}

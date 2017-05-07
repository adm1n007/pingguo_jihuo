package trace

type Exception struct {
    Message     string
    Traceback   string
    Value       interface{}
}

func (self *Exception) String() string {
    return "(traceback)\n" + self.Traceback + "\n" + self.Message
}

func (self *Exception) Error() string {
    return self.String()
}

/*

    exception types

*/

type IndexError struct {
    *BaseException
}

type AttributeError struct {
    *BaseException
}

type TimeoutError struct {
    *BaseException
}

type KeyError struct {
    *BaseException
}

type NotImplementedError struct {
    *BaseException
}

type FileGenericError struct {
    *BaseException
}

type FileNotFoundError struct {
    *BaseException
}

type PermissionError struct {
    *BaseException
}

func NewIndexError(format string, args ...interface{}) *IndexError {
    return &IndexError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewAttributeError(format string, args ...interface{}) *AttributeError {
    return &AttributeError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewTimeoutError(format string, args ...interface{}) *TimeoutError {
    return &TimeoutError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewKeyError(format string, args ...interface{}) *KeyError {
    return &KeyError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewNotImplementedError(format string, args ...interface{}) *NotImplementedError {
    return &NotImplementedError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewFileGenericError(format string, args ...interface{}) *FileGenericError {
    return &FileGenericError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewFileNotFoundError(format string, args ...interface{}) *FileNotFoundError {
    return &FileNotFoundError{
        BaseException: newBaseException(3, format, args...),
    }
}

func NewPermissionError(format string, args ...interface{}) *PermissionError {
    return &PermissionError{
        BaseException: newBaseException(3, format, args...),
    }
}
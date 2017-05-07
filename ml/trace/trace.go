package trace

import (
    "fmt"
    "runtime"
    "reflect"
    "path/filepath"
)

const (
    callerDepth = 10
)

func getCaller(skip int) (name, file string, line int) {
    pc, file, line, _ := runtime.Caller(skip)

    // file = filepath.Base(file)
    name = filepath.Base(runtime.FuncForPC(pc).Name())
    return
}

func raiseimpl(v interface{}) {
    var exp *Exception

    name, _, line := getCaller(3)

    switch e := v.(type) {
        case *Exception:
            exp = e

        default:
            be, isBaseException := e.(baseExceptionInterface)
            if isBaseException {
                exp = &Exception{
                    Traceback   : be.GetTraceBackString(),
                    Message     : fmt.Sprintf("[%s:%d] [%T] %v\n", name, line, v, be.getMessage()),
                    Value       : v,
                }
            } else {
                exp = &Exception{
                    Traceback   : string(stack(3, callerDepth)),
                    Message     : fmt.Sprintf("[%s:%d] [%T] %v\n", name, line, v, v),
                    Value       : v,
                }
            }
    }

    panic(exp)
}

func RaiseIf(err error) {
    if err != nil {
        raiseimpl(NewBaseException(err.Error()))
    }
}

func Raise(v interface{}) {
    raiseimpl(v)
}

func Raisef(v ...interface{}) {
    raiseimpl(NewBaseExceptionWithSkip(1, v[0].(string), v[1:]...))
}

func Catch(exp interface{}) *Exception {
    switch e := exp.(type) {
        case *Exception:
            return e

        case nil:
            return nil

        default:
            const skip = 2
            name, _, line := getCaller(skip)
            return &Exception{
                        Traceback   : string(stack(skip, callerDepth)),
                        Message     : fmt.Sprintf("[%s:%d] %v\n", name, line, e),
                        Value       : e,
                    }
    }
}

func IsError(exp, typ interface{}) bool {
    expType := reflect.TypeOf(exp)
    typType := reflect.TypeOf(typ)

    for expType.Kind() == reflect.Ptr {
        expType = expType.Elem()
    }

    for typType.Kind() == reflect.Ptr {
        typType = typType.Elem()
    }

    return expType == typType
}

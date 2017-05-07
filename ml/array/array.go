package array

import (
    . "active_apple/ml/trace"
    . "fmt"
)

type Array []interface{}

func NewArray(values ...interface{}) *Array {
    arr := &Array{}
    if len(values) != 0 {
        arr.Append(values...)
    }

    return arr
}

func (self *Array) Length() int {
    return len(*self)
}

func (self *Array) Clear() *Array {
    *self = (*self)[:0]
    return self
}

func (self *Array) Contains(value interface{}) bool {
    _, found := self.Index(value)
    return found
}

func (self *Array) Index(value interface{}) (index int, found bool) {
    index = 0
    found = false

    for i, v := range *self {
        if v == value {
            index = i
            found = true
            break
        }
    }

    return
}

func (self *Array) Append(values ...interface{}) *Array {
    *self = append(*self, values...)
    return self
}

func (self *Array) Remove(value interface{}) bool {
    index, ok := self.Index(value)

    if ok {
        self.Pop(index)
    }

    return ok
}

func (self *Array) Pop(index int) (interface{}, bool) {
    if index >= len(*self) {
        return nil, false
    }

    value := (*self)[index]
    *self = append((*self)[:index], (*self)[index + 1:]...)

    return value, true
}

func (self *Array) Peek(index int) interface{} {
    if index >= len(*self) {
        Raise(Sprintf("peek out of index: %d", index))
    }

    return (*self)[index]
}

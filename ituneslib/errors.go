package ituneslib

import (
    . "fmt"
)

type iTunesHelperError struct {
    msg string
}

func (self *iTunesHelperError) Error() string {
    return "iTunesHelper: " + self.msg
}

func newiTunesHelperErrorf(format string, args ...interface{}) *iTunesHelperError {
    return &iTunesHelperError{Sprintf(format, args...)}
}

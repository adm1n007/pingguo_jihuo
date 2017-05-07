package http

import (
    . "ml/array"
)

type RequestOptions struct {
    DontReadResponseBody    bool
    DontFollowRedirects     bool
    IgnoreEncodeKeys        Array
    OverwriteHeaders        bool
    AutoRetry               bool
    Ignore404               bool
    MaxTimeoutTimes         int
}

const (
    DefaultMaxTimeoutTimes = 3
)

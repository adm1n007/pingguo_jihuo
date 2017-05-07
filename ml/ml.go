package ml

import (
    "runtime"
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

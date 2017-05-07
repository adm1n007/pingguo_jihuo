package entry

import (
    "ituneslib"
    "runtime/pprof"
    "os"
)

func b2i(b bool) int {
    switch b {
        case true:
            return 1

        default:
            return 0
    }
}

func dumpGoroutines(filename string) {
    return
    callstack, _ := os.Create(filename)
    defer callstack.Close()

    pprof.Lookup("goroutine").WriteTo(callstack, 1)
}

func init() {
    ituneslib.Initialize()
}

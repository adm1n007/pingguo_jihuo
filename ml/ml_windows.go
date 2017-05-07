package ml

import (
    "syscall"
)

func fixTimePeriod() {
    syscall.MustLoadDLL("winmm.dll").MustFindProc("timeEndPeriod").Call(uintptr(1))
}

func init() {
    // fixTimePeriod()
}

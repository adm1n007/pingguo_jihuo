package os2

import (
    "syscall"
    "unsafe"
    "unicode/utf16"
)

func getExecutable() string {
    var GetModuleFileNameW = syscall.MustLoadDLL("kernel32.dll").MustFindProc("GetModuleFileNameW")
    b := make([]uint16, 0x7FFF)
    r, _, _ := GetModuleFileNameW.Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
    n := uint32(r)
    if n == 0 {
        return ""
    }
    return string(utf16.Decode(b[:n]))
}

package syscall

import (
    "syscall"
    "fmt"
)

func Call(proc uintptr, args ...uintptr) (r1, r2 uintptr, lastErr error) {
    addr := proc
    switch len(args) {
        case 0:
            return syscall.Syscall(addr, uintptr(len(args)), 0, 0, 0)
        case 1:
            return syscall.Syscall(addr, uintptr(len(args)), args[0], 0, 0)
        case 2:
            return syscall.Syscall(addr, uintptr(len(args)), args[0], args[1], 0)
        case 3:
            return syscall.Syscall(addr, uintptr(len(args)), args[0], args[1], args[2])
        case 4:
            return syscall.Syscall6(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], 0, 0)
        case 5:
            return syscall.Syscall6(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], 0)
        case 6:
            return syscall.Syscall6(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5])
        case 7:
            return syscall.Syscall9(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], 0, 0)
        case 8:
            return syscall.Syscall9(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], 0)
        case 9:
            return syscall.Syscall9(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8])
        case 10:
            return syscall.Syscall12(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], 0, 0)
        case 11:
            return syscall.Syscall12(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], 0)
        case 12:
            return syscall.Syscall12(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11])
        case 13:
            return syscall.Syscall15(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], 0, 0)
        case 14:
            return syscall.Syscall15(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], 0)
        case 15:
            return syscall.Syscall15(addr, uintptr(len(args)), args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12], args[13], args[14])
        default:
            panic(fmt.Sprintf("Call \"%d\" with too many arguments @ %d.", proc, len(args)))
    }

    return
}

package syscall

type Proc uintptr

func (p Proc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
    return Call(uintptr(p), a...)
}

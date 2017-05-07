package console

import (
    "syscall"
    "unsafe"
)

var (
    getch           *syscall.Proc
    SetConsoleTitle *syscall.Proc
)

func pause() {
    getch.Call()
}

func setTitle(text string) {
    SetConsoleTitle.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))))
}

func init() {
    getch = syscall.MustLoadDLL("msvcrt.dll").MustFindProc("_getch")
    SetConsoleTitle = syscall.MustLoadDLL("kernel32.dll").MustFindProc("SetConsoleTitleW")
}

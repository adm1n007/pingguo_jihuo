package os2

import (
    "path/filepath"
    "strconv"
)

func Executable() string {
    return getExecutable()
}

func ExecutableName() string {
    name := filepath.Base(Executable())
    return name[:len(name) - len(filepath.Ext(name))]
}

func ExecutablePath() string {
    return filepath.Dir(Executable())
}

const (
    PtrSize_32bits  = 32
    PtrSize_64bits  = 64
)

func PtrSize() int {
    return strconv.IntSize
}

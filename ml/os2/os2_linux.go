package os2

import (
    "os"
    "strings"
)

func getExecutable() string {
    const deletedTag = " (deleted)"
    execpath, err := os.Readlink("/proc/self/exe")
    if err != nil {
        return execpath
    }

    return strings.TrimPrefix(strings.TrimSuffix(execpath, deletedTag), deletedTag)
}

package filestream

import (
    "os"
)

type FileGeneric interface {
    Close() error
    Stat() (os.FileInfo, error)
    Truncate(size int64) error
    Seek(offset int64, whence int) (int64, error)
    Read(b []byte) (int, error)
    Write(b []byte) (int, error)
    Sync() error
}
package filestream

import (
    "os"
    "io"
    "time"
)

type bytesIOFileInfo struct {
    size    int64
}

func (fs *bytesIOFileInfo) Name() string       { return "<bytesIO>" }
func (fs *bytesIOFileInfo) Size() int64        { return fs.size }
func (fs *bytesIOFileInfo) Mode() os.FileMode  { return 0 }
func (fs *bytesIOFileInfo) ModTime() time.Time { return time.Time{} }
func (fs *bytesIOFileInfo) IsDir() bool        { return false }
func (fs *bytesIOFileInfo) Sys() interface{}   { return nil }

type bytesIO struct {
    buffer      []byte
    position    int64
}

func newBytesIO() *bytesIO {
    return &bytesIO{[]byte{}, 0}
}

func (self *bytesIO) length() int64 {
    return int64(len(self.buffer))
}

func (self *bytesIO) resize(size int64) {
    if self.length() >= size {
        return
    }

    buffer := make([]byte, size)
    copy(buffer, self.buffer)

    self.buffer = buffer
}

func (self *bytesIO) ensure(n int64) {
    remain := self.length() - self.position
    if remain >= n {
        return
    }

    self.resize(self.length() + n - remain)
}

func (self *bytesIO) Close() error {
    return nil
}

func (self *bytesIO) Stat() (os.FileInfo, error) {
    return &bytesIOFileInfo{size: self.length()}, nil
}

func (self *bytesIO) Truncate(size int64) error {
    if size >= self.length() {
        return nil
    }

    self.buffer = self.buffer[:size]
    return nil
}

func (self *bytesIO) Seek(offset int64, whence int) (int64, error) {
    length := self.length()
    var position int64

    switch whence {
        case os.SEEK_SET:
            position = offset

        case os.SEEK_CUR:
            position = self.position + offset

        case os.SEEK_END:
            position = length + offset
    }

    self.resize(position)
    self.position = position

    return position, nil
}

func (self *bytesIO) Read(b []byte) (int, error) {
    if self.position == self.length() {
        return 0, io.EOF
    }

    n := copy(b, self.buffer[self.position:])
    self.position += int64(n)
    return n, nil
}

func (self *bytesIO) Write(b []byte) (int, error) {
    self.ensure(int64(len(b)))
    n := copy(self.buffer[self.position:], b)
    self.position += int64(n)
    return n, nil
}

func (self *bytesIO) Sync() error {
    return nil
}

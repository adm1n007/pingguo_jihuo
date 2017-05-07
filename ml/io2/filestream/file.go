package filestream

import (
    . "active_apple/ml"
    . "active_apple/ml/strings"
    . "active_apple/ml/trace"

    "os"
    "io"
    "math"
    "unsafe"
    "reflect"
    "encoding/binary"
)

var (
    BigEndian       = &binary.BigEndian
    LittleEndian    = &binary.LittleEndian

    SEEK_SET        = os.SEEK_SET
    SEEK_CUR        = os.SEEK_CUR
    SEEK_END        = os.SEEK_END

    READ            = 1 << 0
    WRITE           = 1 << 1
    READWRITE       = READ | WRITE
    CREATE          = 1 << 2
)

const (
    END_OF_FILE     = -1
)

type File struct {
    file FileGeneric
    Endian binary.ByteOrder
}

func Open(name string) *File {
    return CreateFile(name, READ)
}

func Create(name string) *File {
    return CreateFile(name, READWRITE | CREATE)
}

func CreateMemory(buffer ...[]byte) *File {
    buf := &File{newBytesIO(), LittleEndian}

    if len(buffer) != 0 {
        buf.Write(buffer[0])
        buf.SetPosition(0)
    }

    return buf
}

func CreateFile(name string, mode int) *File {
    flag := 0

    switch {
        case (mode & READWRITE) == READWRITE:
            flag = os.O_RDWR

        case (mode & READ) == READ:
            flag = os.O_RDONLY

        case (mode & WRITE) == WRITE:
            flag = os.O_WRONLY
    }

    if (mode & CREATE) != 0 {
        flag |= os.O_TRUNC | os.O_CREATE
    }

    f, err := os.OpenFile(name, flag, 0666)

    if err == nil {
        return &File{f, LittleEndian}
    }

    switch {
        case os.IsNotExist(err):
            Raise(NewFileNotFoundError(err.Error()))

        case os.IsPermission(err):
            Raise(NewPermissionError(err.Error()))

        default:
            Raise(NewFileGenericError(err.Error()))
    }

    return nil
}

func (self *File) Close() {
    if self.file != nil {
        self.file.Close()
        self.file = nil
    }
}

func (self *File) Writer() io.Writer {
    return self.file
}

func (self *File) Length() int64 {
    fi, err := self.file.Stat()
    raiseGenericError(err)
    return fi.Size()
}

func (self *File) SetLength(length int64) {
    pos := self.Position()
    if pos > self.Length() {
        self.SetPosition(length)
        self.Write(byte(0))
    }

    err := self.file.Truncate(length)
    raiseGenericError(err)

    if length < pos {
        pos = length
    }

    self.SetPosition(pos)
}

func (self *File) Remaining() int64 {
    return self.Length() - self.Position()
}

func (self *File) Position() int64 {
    return self.Seek(0, SEEK_CUR)
}

func (self *File) SetPosition(offset int64) {
    if offset == END_OF_FILE {
        self.Seek(0, SEEK_END)
        return
    }

    self.Seek(offset, SEEK_SET)
}

func (self *File) Seek(offset int64, whence int) int64 {
    pos, err := self.file.Seek(offset, whence)
    raiseGenericError(err)
    return pos
}

func (self *File) Read(n int) []byte {
    buffer := [1024]byte{}
    data := []byte{}

    for n > 0 {
        bytesRead := If(n > len(buffer), len(buffer), n).(int)
        bytesRead, err := self.file.Read(buffer[:bytesRead])
        if err == io.EOF {
            data = append(data, buffer[:bytesRead]...)
            break
        }

        raiseGenericError(err)

        data = append(data, buffer[:bytesRead]...)
        n -= bytesRead
    }

    return data
}

func (self *File) Write(args ...interface{}) int {
    buffer := []byte{}

    data := args[0]

    switch b := data.(type) {
        case bool:
            buffer = append(buffer, If(b, 1, 0).(byte))

        case int8, byte:
            buffer = append(buffer, b.(byte))

        case int16, uint16:
            buf := [8]byte{}
            self.Endian.PutUint16(buf[:], b.(uint16))
            buffer = buf[:2]

        case int32, uint32:
            buf := [8]byte{}
            self.Endian.PutUint32(buf[:], b.(uint32))
            buffer = buf[:4]

        case int64, uint64:
            buf := [8]byte{}
            self.Endian.PutUint64(buf[:], b.(uint64))
            buffer = buf[:8]

        case float32:
            buf := [8]byte{}
            self.Endian.PutUint32(buf[:], math.Float32bits(b))
            buffer = buf[:4]

        case float64:
            buf := [8]byte{}
            self.Endian.PutUint64(buf[:], math.Float64bits(b))
            buffer = buf[:8]

        case []byte:
            buffer = b

        case String:
            codpage := CP_UTF8
            if len(args) > 1 {
                codpage = args[1].(Encoding)
            }

            buffer = b.Encode(codpage)

        case string:
            s := String(b)
            codpage := CP_UTF8
            if len(args) > 1 {
                codpage = args[1].(Encoding)
            }

            buffer = s.Encode(codpage)
    }

    n, err := self.file.Write(buffer)
    raiseGenericError(err)
    return n
}

func (self *File) ReadAll() (data []byte) {
    self.SetPosition(0)
    return self.ReadRemaining()
}

func (self *File) ReadRemaining() (data []byte) {
    length := self.Remaining()

    for length > 0 {
        read := self.Read(int(length))
        if len(read) == 0 {
            break
        }

        length -= int64(len(read))
        data = append(data, read...)
    }

    return
}

func (self *File) Flush() {
    raiseGenericError(self.file.Sync())
}

func (self *File) IsEndOfFile() bool {
    return self.Position() >= self.Length()
}

func (self *File) ReadBoolean() bool {
    return self.ReadByte() != 0
}

func (self *File) ReadChar() int {
    return int(int8(self.Read(1)[0]))
}

func (self *File) ReadByte() uint {
    return uint(self.Read(1)[0])
}

func (self *File) ReadShort() int {
    return int(int16(self.Endian.Uint16(self.Read(2))))
}

func (self *File) ReadUShort() uint {
    return uint(self.Endian.Uint16(self.Read(2)))
}

func (self *File) ReadLong() int {
    return int(self.Endian.Uint32(self.Read(4)))
}

func (self *File) ReadULong() uint {
    return uint(self.Endian.Uint32(self.Read(4)))
}

func (self *File) ReadLong64() int64 {
    return int64(self.Endian.Uint64(self.Read(8)))
}

func (self *File) ReadULong64() uint64 {
    return self.Endian.Uint64(self.Read(8))
}

func (self *File) ReadFloat() float32 {
    return math.Float32frombits(uint32(self.ReadULong()))
}

func (self *File) ReadDouble() float64 {
    return math.Float64frombits(self.ReadULong64())
}

func (self *File) ReadUVarint() (ret uint64) {
    buf := [binary.MaxVarintLen64]byte{}

    var b byte = 0x80

    for i := 0; b >= 0x80; i++ {
        b = byte(self.ReadByte())
        buf[i] = b
    }

    ret, _ = binary.Uvarint(buf[:])

    return
}

func (self *File) ReadMultiByte(encoding ...Encoding) String {
    bytes := []byte{}

    codepage := CP_UTF8

    switch len(encoding) {
        case 1:
            codepage = encoding[0]
    }

    for {
        ch := self.ReadByte()
        if ch == 0 {
            break
        }

        bytes = append(bytes, byte(ch))
    }

    return Decode(bytes, codepage)
}

func (self *File) ReadUTF16() String {
    bytes := []byte{}

    codepage := CP_UTF16_LE
    if self.Endian == BigEndian {
        codepage = CP_UTF16_BE
    }

    for {
        ch := self.ReadUShort()
        if ch == 0 {
            break
        }

        bytes = append(bytes, byte(ch), byte(ch >> 8))
    }

    return Decode(bytes, codepage)
}

func (self *File) ReadType(t interface{}) interface{} {
    typ := reflect.TypeOf(t)
    bytes:= self.Read(int(typ.Size()))
    return reflect.NewAt(typ, unsafe.Pointer(&bytes[0])).Elem().Interface()
}

func (self *File) WriteType(v interface{}) int {
    val := reflect.ValueOf(v)
    for val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    typ := val.Type()
    addr := val.Addr().Pointer()
    p := (*[^uint32(0) >> 1]byte)(unsafe.Pointer(addr))

    return self.Write(p[:typ.Size()])
}
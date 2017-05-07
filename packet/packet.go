package packet

import (
    . "fmt"
    . "active_apple/ml"
    . "active_apple/ml/dict"
    . "active_apple/ml/trace"
    . "active_apple/ml/strings"

    "crypto/aes"
    "crypto/cipher"
    "encoding/json"

    "active_apple/ml/encoding/binary"
    "active_apple/ml/io2/filestream"
    "active_apple/ml/net/socket"
)

type Packet struct {
    ServerId    int
    SerialId    int
    Aeskey      []byte
    AppId       int
    Headers     OrderedDict
    Body        OrderedDict
    Compressed  bool
    bodyData    []byte
}

func Recieve(sock socket.Socket, aesKey []byte) *Packet {
    magic := sock.ReadAll(1)
    if magic[0] != MAGIC {
        Raise(NewPacketError("incorrect packet magic"))
    }

    var headerLength    uint64 = 0
    var bodyLength      uint64 = 0

    flags           := int(binary.BigEndian.Uint16(sock.ReadAll(2)))
    appId           := int(binary.BigEndian.Uint16(sock.ReadAll(2)))
    serverId        := int(binary.BigEndian.Uint16(sock.ReadAll(2)))
    serialId        := int(binary.BigEndian.Uint16(sock.ReadAll(2)))

    if FlagOn(flags, UNKNOWN) {
        // status
        sock.ReadAll(2)
    }

    if FlagOn(flags, HAS_HEADERS) {
        buf := []byte{}
        for b := byte(0x80); b >= 0x80; {
            b = sock.ReadAll(1)[0]
            buf = append(buf, b)
        }

        headerLength, _ = binary.Varint(buf)
    }

    if FlagOn(flags, HAS_BODY) {
        buf := []byte{}
        for b := byte(0x80); b >= 0x80; {
            b = sock.ReadAll(1)[0]
            buf = append(buf, b)
        }

        bodyLength, _ = binary.Varint(buf)
    }

    packet := NewPacket(serverId)
    packet.SerialId = serialId
    packet.AppId = appId
    packet.Aeskey = aesKey

    if headerLength != 0 {
        headers := packet.decryptData(sock.ReadAll(int(headerLength)))
        buffer := filestream.CreateMemory(headers)
        buffer.Endian = filestream.BigEndian

        for buffer.IsEndOfFile() == false {
            itemType := buffer.ReadUShort()
            itemSize := buffer.ReadUVarint()
            packet.AddHeader(Headers(itemType), buffer.Read(int(itemSize)))
        }
    }

    if bodyLength != 0 {
        body := packet.decryptData(sock.ReadAll(int(bodyLength)))

        var b interface{}
        RaiseIf(json.Unmarshal(body, &b))

        switch t := b.(type) {
            case map[string]interface{}:
                for k, v := range t {
                    packet.AddBody(String(k), v)
                }

            case []interface{}:
                for i, v := range t {
                    packet.AddBody(String(Sprintf("%d", i)), v)
                }
        }
    }

    return packet
}

func NewPacket(serverId int) *Packet {
    return &Packet{
        ServerId    : serverId,
        SerialId    : 1,
        Aeskey      : nil,
        AppId       : 2,
        Headers     : NewOrderedDict(),
        Body        : NewOrderedDict(),
        Compressed  : false,
    }
}

func (self *Packet) AddHeader(key Headers, value []byte) {
    self.Headers.Set(key, value)
}

func (self *Packet) AddBody(key String, value interface{}) {
    self.Body.Set(key.String(), value)
    self.bodyData = nil
}

func (self *Packet) SetBodyData(body []byte) {
    self.Body.Clear()
    self.bodyData = body
}

func (self *Packet) ToBinary() (data []byte) {
    buffer := filestream.CreateMemory()
    buffer.Endian = filestream.BigEndian

    buffer.Write(byte(MAGIC))

    flags := 0

    flags |= If(self.Headers.Length()   != 0,   HAS_HEADERS,    0).(int)
    flags |= If(self.Body.Length()      != 0,   HAS_BODY,       0).(int)
    flags |= If(self.bodyData           != nil, HAS_BODY,       0).(int)
    flags |= If(self.Aeskey             != nil, ENCRYPTED,      0).(int)
    flags |= If(self.Compressed,                COMPRESSED,     0).(int)

    buffer.Write(uint16(flags))
    buffer.Write(uint16(self.AppId))
    buffer.Write(uint16(self.ServerId))
    buffer.Write(uint16(self.SerialId))

    headers := []byte{}
    body := []byte{}

    if (flags & HAS_HEADERS) != 0 {
        headers = self.encryptData(self.headersToBinary())
        buffer.Write(binary.ToVarint(len(headers)))
    }

    if (flags & HAS_BODY) != 0 {
        if self.bodyData == nil {
            body = self.encryptData(self.Body.Json())
        } else {
            body = self.encryptData(self.bodyData)
        }

        buffer.Write(binary.ToVarint(len(body)))
    }

    buffer.Write(headers)
    buffer.Write(body)

    return buffer.ReadAll()
}

func (self *Packet) headersToBinary() []byte {
    buffer := filestream.CreateMemory()
    buffer.Endian = filestream.BigEndian

    for _, item := range self.Headers.Items() {
        v := item.Value.([]byte)
        buffer.Write(uint16(item.Key.(Headers)))
        buffer.Write(binary.ToVarint(len(v)))
        buffer.Write(v)
    }

    return buffer.ReadAll()
}

func (self *Packet) encryptData(data []byte) []byte {
    data = data[:]

    if self.Aeskey == nil {
        return data
    }

    aes, err := aes.NewCipher(self.Aeskey)
    RaiseIf(err)

    padding := aes.BlockSize() - len(data) % aes.BlockSize()

    for i := 0; i != padding; i++ {
        data = append(data, byte(padding))
    }

    cipher.NewCBCEncrypter(aes, []byte("0102030405060708")).CryptBlocks(data, data)

    return data
}

func (self *Packet) decryptData(data []byte) []byte {
    data = data[:]

    if self.Aeskey == nil {
        return data
    }

    aes, err := aes.NewCipher(self.Aeskey)
    RaiseIf(err)

    cipher.NewCBCDecrypter(aes, []byte("0102030405060708")).CryptBlocks(data, data)

    padding := int(data[len(data) - 1])

    return data[:len(data) - padding]
}

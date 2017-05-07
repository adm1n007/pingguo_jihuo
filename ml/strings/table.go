package strings

import (
    "io/ioutil"
    "compress/zlib"
    "bytes"
    "encoding/binary"
    "sync"
)

const (
    MAXIMUM_LEADBYTES   = 12
    SIZE_OF_WCHAR       = 2
    UnicodeNull         = 0
)

const (
    MB_TBL_SIZE         = 256             /* size of MB tables */
    GLYPH_TBL_SIZE      = MB_TBL_SIZE     /* size of GLYPH tables */
    DBCS_TBL_SIZE       = 256             /* size of DBCS tables */
    GLYPH_HEADER        = 1               /* size of GLYPH table header */
    DBCS_HEADER         = 1               /* size of DBCS table header */
    LANG_HEADER         = 1               /* size of LANGUAGE file header */
    UP_HEADER           = 1               /* size of UPPERCASE table header */
    LO_HEADER           = 1               /* size of LOWERCASE table header */
)

type codePageTableInfo struct {
    CodePage                uint16                          // code page number
    MaximumCharacterSize    uint16                          // max length (bytes) of a char
    DefaultChar             uint16                          // default character (MB)
    UniDefaultChar          uint16                          // default character (Unicode)
    TransDefaultChar        uint16                          // translation of default char (Unicode)
    TransUniDefaultChar     uint16                          // translation of Unic default char (MB)
    DBCSCodePage            bool                            // Non 0 for DBCS code pages
    // LeadByte                [MAXIMUM_LEADBYTES]byte         // lead byte ranges
    MultiByteTable          []uint16                        // pointer to MB translation table
    WideCharTable           []uint16                        // pointer to WC translation table
    // DBCSRanges              uint                            // pointer to DBCS ranges
    TranslateTable          []uint16                        // pointer to DBCS offsets

    data []byte
    initialized bool

    encode func(table *codePageTableInfo, str String) []byte
    decode func(table *codePageTableInfo, bytes []byte) String
}

var lock = &sync.Mutex{}

func bytesToUInt16Array(bytes []byte) []uint16 {
    arr := make([]uint16, len(bytes) / 2)

    for i := range(arr) {
        arr[i] = uint16(bytes[i * 2]) | (uint16(bytes[i * 2 + 1]) << 8)
    }

    return arr
}

func decompressData(compressed []byte) []byte {
    reader, err := zlib.NewReader(bytes.NewReader(compressed))
    if err != nil {
        panic(err)
    }

    defer reader.Close()
    compressed, err = ioutil.ReadAll(reader)
    if err != nil {
        panic(err)
    }

    return compressed
}

func extractTable(bytes []byte) ([]uint16, int) {
    tableSize, offset := binary.Uvarint(bytes)
    table := bytes[offset:offset + int(tableSize)]
    return bytesToUInt16Array(table), offset + int(tableSize)
}

func (self *codePageTableInfo) initialize() {
    if self.initialized {
        return
    }

    lock.Lock()
    defer lock.Unlock()

    if self.initialized {
        return
    }

    uncompressed := decompressData(self.data)
    self.data = nil

    offset := 0
    size := 0

    self.MultiByteTable, size = extractTable(uncompressed)
    offset += size

    self.TranslateTable, size = extractTable(uncompressed[offset:])
    offset += size

    self.WideCharTable, size = extractTable(uncompressed[offset:])
    offset += size

    self.initialized = true
}

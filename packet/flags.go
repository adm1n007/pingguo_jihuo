package packet

const (
    MAGIC       = 0xD1
    COMPRESSED  = 0x0400
    HAS_BODY    = 0x1000
    HAS_HEADERS = 0x2000
    UNKNOWN     = 0x4000
    ENCRYPTED   = 0x8000
)

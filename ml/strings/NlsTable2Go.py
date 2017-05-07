from ml import *
import zlib

MAXIMUM_LEADBYTES   = 12
MB_TBL_SIZE         = 256
GLYPH_TBL_SIZE      = MB_TBL_SIZE
DBCS_TBL_SIZE       = 256
GLYPH_HEADER        = 1
DBCS_HEADER         = 1
LANG_HEADER         = 1
UP_HEADER           = 1
LO_HEADER           = 1

class CPTABLEINFO(Structure):
    _fields_ = [
        ('CodePage',                USHORT),
        ('MaximumCharacterSize',    USHORT),
        ('DefaultChar',             USHORT),
        ('UniDefaultChar',          USHORT),
        ('TransDefaultChar',        USHORT),
        ('TransUniDefaultChar',     USHORT),
        ('DBCSCodePage',            USHORT),
        ('LeadByte',                BYTE * MAXIMUM_LEADBYTES),
        ('MultiByteTable',          ULONG_PTR),
        ('WideCharTable',           ULONG_PTR),
        ('DBCSRanges',              ULONG_PTR),
        ('DBCSOffsets',             ULONG_PTR),
    ]

PCPTABLEINFO = ctypes.POINTER(CPTABLEINFO)

def RtlInitCodePageTable(table):
    info = CPTABLEINFO()
    buf = (BYTE * len(table)).from_buffer(table)

    windll.ntdll.RtlInitCodePageTable(byref(buf), PCPTABLEINFO(info))

    return info

def EncodeVarint(value):
    o = value
    varint = b''

    bits = value & 0x7f
    value >>= 7

    while value:
        varint += int.to_bytes(0x80 | bits, 1, 'little')
        bits = value & 0x7f
        value >>= 7

    varint += int.to_bytes(bits, 1, 'little')

    print('%X' % o, varint)

    return varint

def main():
    if len(sys.argv) == 1:
        files = EnumDirectoryFiles(os.path.dirname(sys.argv[0]), '*.nls')
        print(files)
        sys.argv.extend(files)

    for f in sys.argv[1:]:
        table = bytearray(open(f, 'rb').read())

        info = RtlInitCodePageTable(table)

        # print('CodePage             = %s' % info.CodePage)
        # print('MaximumCharacterSize = %X' % info.MaximumCharacterSize)
        # print('DefaultChar          = %X' % info.DefaultChar)
        # print('UniDefaultChar       = %X' % info.UniDefaultChar)
        # print('TransDefaultChar     = %X' % info.TransDefaultChar)
        # print('TransUniDefaultChar  = %X' % info.TransUniDefaultChar)
        # print('DBCSCodePage         = %s' % info.DBCSCodePage)
        # print('LeadByte             = %s' % bytes(info.LeadByte))
        # print('MultiByteTable       = %s' % (info.MultiByteTable))
        # print('WideCharTable        = %s' % (info.WideCharTable))
        # print('DBCSRanges           = %s' % (info.DBCSRanges))
        # print('DBCSOffsets          = %s' % (info.DBCSOffsets))

        cp = info.CodePage

        MultiByteTable = bytes((USHORT * MB_TBL_SIZE).from_address(info.MultiByteTable))

        TranslateTableSize = info.WideCharTable - info.DBCSOffsets - 2
        TranslateTable = bytes((USHORT * (TranslateTableSize // 2)).from_address(info.DBCSOffsets))
        WideCharTable = bytes((USHORT * 0x10000).from_address(info.WideCharTable))

        compressed = bytearray()

        # MultiByteTable = zlib.compress(MultiByteTable, 9)
        compressed.extend(EncodeVarint(len(MultiByteTable)))
        compressed.extend(MultiByteTable)

        # TranslateTable = zlib.compress(TranslateTable, 9)
        compressed.extend(EncodeVarint(len(TranslateTable)))
        compressed.extend(TranslateTable)

        # WideCharTable = zlib.compress(WideCharTable, 9)
        compressed.extend(EncodeVarint(len(WideCharTable)))
        compressed.extend(WideCharTable)

        compressed = zlib.compress(compressed, 9)

        data = []
        for i in range(0, len(compressed), 16):
            data.append(''.join(['0x%02X,' % ch for ch in compressed[i:i + 16]]).rstrip())

        src = [
            'package strings',
            '',
        ]

        src.extend([
            'func init() {',
            '    var compressedData = []byte{\n        %s' % ('\n        '.join(data)),
            '    }'
            '',
            '',
            '    cptable[%d] = codePageTableInfo{' % cp,
            '                        CodePage                : %d,' % cp,
            '                        MaximumCharacterSize    : 0x%X,' % info.MaximumCharacterSize,
            '                        DefaultChar             : 0x%X,' % info.DefaultChar,
            '                        UniDefaultChar          : 0x%X,' % info.UniDefaultChar,
            '                        TransDefaultChar        : 0x%X,' % info.TransDefaultChar,
            '                        TransUniDefaultChar     : 0x%X,' % info.TransUniDefaultChar,
            '                        DBCSCodePage            : %s,' % (info.DBCSCodePage and 'true' or 'false'),
            '                        data                    : compressedData,',
            '                        initialized             : false,',
        ])

        padding = '                        '

        # LeadByte = 'LeadByte                : []byte[]{%s},' % ','.join(['0x%X' % ch for ch in bytes(info.LeadByte)])
        # src.append(padding + LeadByte)

        src.extend([
            '                    }',
            '}',
        ])

        f = os.path.join(os.path.dirname(f), '%d.go' % cp)
        print(f)
        open(f, 'wb').write('\n'.join(src).encode('UTF8'))

if __name__ == '__main__':
    TryInvoke(main)

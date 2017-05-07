package strings

import (
    "unicode/utf8"
)

func CustomCPToUnicodeN(CustomCP *codePageTableInfo, CustomCPString []byte)(UnicodeString String) {

    BytesInCustomCPString := uint(len(CustomCPString))

    runes := []rune{}

    if CustomCP.DBCSCodePage == false {
        TranslateTable := CustomCP.MultiByteTable

        for i := uint(0); i != BytesInCustomCPString; i++ {
            runes = append(runes, rune(TranslateTable[CustomCPString[i]]))
        }

    } else {
        NlsCustomLeadByteInfo := CustomCP.TranslateTable
        TranslateTable := CustomCP.TranslateTable
        index := 0

        for BytesInCustomCPString != 0 {
            BytesInCustomCPString--

            Entry := uint(NlsCustomLeadByteInfo[CustomCPString[index]])

            if Entry != 0 {
                if BytesInCustomCPString == 0 {
                    break
                }

                TailByte := uint(CustomCPString[index + 1])
                runes = append(runes, rune(TranslateTable[Entry + TailByte]))

                index += 2
                BytesInCustomCPString--

            } else {
                runes = append(runes, rune(CustomCP.MultiByteTable[CustomCPString[index]]))
                index++
            }

        }
    }

    return String(string(runes))
}

func UnicodeToCustomCPN(CustomCP *codePageTableInfo, UnicodeString_ String) (CustomCPString []byte) {

    Ucs16String := []uint16{}
    UnicodeString := string(UnicodeString_)

    for len(UnicodeString) > 0 {
        r, size := utf8.DecodeRuneInString(UnicodeString)
        UnicodeString = UnicodeString[size:]
        Ucs16String = append(Ucs16String, uint16(r))
    }

    CharsInUnicodeString := len(Ucs16String)

    if CustomCP.DBCSCodePage == false {
        TranslateTable := CustomCP.WideCharTable

        for i := 0; i != CharsInUnicodeString; i++ {
            CustomCPString = append(CustomCPString, byte(TranslateTable[Ucs16String[i]]))
        }

    } else {
        WideTranslateTable := CustomCP.WideCharTable
        index := 0

        for ; CharsInUnicodeString != 0; CharsInUnicodeString-- {
            MbChar := WideTranslateTable[Ucs16String[index]]
            index++

            if (MbChar & 0xFF00 != 0) {
                CustomCPString = append(CustomCPString, byte(MbChar >> 8))    // lead byte
            }

            CustomCPString = append(CustomCPString, byte(MbChar & 0xFF))
        }
    }

    return
}

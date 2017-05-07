package utility

import (
	"encoding/json"
	"ml/random"
	. "ml/strings"

	"github.com/axgle/mahonia"
)

var pinyin = [][]string{}

func GenerateName(min, max int) String {
	var name String
	for n := random.IntRange(min, max+1); n > 0; {
		p := String(pinyin[random.ChoiceIndex(pinyin)][0])
		if p.Length() > 1 {
			continue
		}

		name += p
		n -= p.Length()
	}

	return String(name)
}

func GenerateChineseFamilyName() String {
	var name String
	name = String(familyNames[random.ChoiceIndex(familyNames)])
	return name
}

func GenerateChineseGivenName() String {
	var name String
	name = String(givenNames[random.ChoiceIndex(givenNames)])
	return name
}

func GenerateLevel1Hanzi() String {
	// high
	var prebyte = random.IntRange(0xB1, 0xD7+1)
	var sufmax = 0xFE
	if prebyte == 0xD7 {
		sufmax = 0xF9
	}
	// lo
	var sufbyte = random.IntRange(0xA1, sufmax+1)
	// [] slice, [2] array?
	gbkBytes := []byte{byte(prebyte), byte(sufbyte)}
	utf8 := mahonia.NewDecoder("gbk").ConvertString(string(gbkBytes))
	return String(utf8)
}

func GenerateLevel1Name(min, max int) String {
	var name String
	for n := random.IntRange(min, max+1); n > 0; {
		p := GenerateLevel1Hanzi()

		name += p
		n -= p.Length()
	}

	return String(name)
}

func GeneratePinyin() String {
	name := ""
	for n := random.IntRange(3, 5); n > 0; {
		p := pinyin[random.ChoiceIndex(pinyin)][1]
		if len(p) > 4 {
			continue
		}

		name += p
		n--
	}

	return String(name)
}

func init() {
	json.Unmarshal([]byte(pydata), &pinyin)
}

package utility

import (
	. "fmt"
	. "ml/strings"
	. "ml/trace"

	"math/rand"
	"regexp"

	"ml/console"
	"ml/random"
	"ml/uuid"

	"github.com/PuerkitoBio/goquery"
)

var _ = Println
var Pause = console.Pause
var SetTitle = console.SetTitle

var (
	DIGITS_PATTERN          = regexp.MustCompile(`\d`)
	CHECK_DUPLICATE_PATTERN = regexp.MustCompile(`(.)\\1{2}`)
)

func quickHash(msg String) String {
	var hi8 uint32 = 0
	var lo8 uint32 = 0
	if len(msg) == 0 {
		return "00000000"
	}
	ml := len(msg)
	var c uint8
	for i := 0; i < ml; i++ {
		c = msg[i]
		if i&1 != 0 {
			lo8 = (lo8 << 5) - lo8 + uint32(c)
		} else {
			hi8 = (hi8 << 5) - hi8 + uint32(c)
		}
	}

	return String(Sprintf("%x%x", hi8, lo8))
}

func GenerateDeviceInfo() String {
	guid := String("").Join([]String{
		GenerateRandomString(),
		GenerateRandomString(),
		GenerateRandomString(),
		GenerateRandomString(),
		GenerateRandomString(),
	})

	var b uint32 = 0
	var c uint32 = 0

	for index, ch := range guid {
		d := uint32(ch)

		switch index % 2 {
		case 0:
			c = (c << 5) - c + d

		default:
			b = (b << 5) - b + d
		}
	}

	return "d1-5a3a494bb3ef19f"

	return String(Sprintf("d1-%x%x", c, b))
}

func GenerateBrowserInfo() (width, height, colorDepth, innerWidth, innerHeight int, language String) {
	if false {
		return 1920, 1080, 24, 1920, 989, "zh-cn"

	} else {
		bpp := []int{16, 24}
		lang := []String{`zh-cn`, `en-us`}

		return 1280, 720, random.Choice(bpp).(int), 1280, 660, random.Choice(lang).(String)
	}

	// return 1920, 1080, random.Choice(bpp).(int), 1920, 1036, random.Choice(lang).(String)
}

func GenerateiOSBrowserInfo() (width, height, colorDepth, innerWidth, innerHeight int, language String) {
	return 320, 480, 32, 320, 460, "zh-cn"
	// return 1920, 1080, random.Choice(bpp).(int), 1920, 1036, random.Choice(lang).(String)
}

func GenerateRandomString() String {
	u, err := uuid.NewV4()
	RaiseIf(err)

	return String(u.String()).Replace("-", "")
}

func GenerateUUID() String {
	u, err := uuid.NewV4()
	RaiseIf(err)

	return String(u.String()).ToUpper()
}

func GenerateMachineGuid() String {
	guid := make([]String, 7)

	for i := 0; i != 7; i++ {
		guid[i] = GenerateRandomString()[:8]
	}

	guid[2] = "00000000"

	return String(".").Join(guid).ToUpper()
}

func GenerateOSXMachineGuid() String {
	return GenerateRandomString()[:6] + GenerateRandomString()[:6]
	return String(Sprintf("%d", random.IntRange(100000, 1000000))) + GenerateRandomString().ToUpper()[:6]
}

func GenerateiOSMachineGuid() String {
	return (GenerateRandomString() + GenerateRandomString())[:40]
}

func GeneratePassword() String {
	LOWER := []rune(`abcdefghijklmnopqrstuvwxyz`)
	UPPER := []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZ`)
	DIGIT := []rune(`1234567890`)

	for {
		password := []rune{}

		for _, set := range [][]rune{LOWER, UPPER, DIGIT} {
			for length := random.IntRange(5, 10); length != 0; length-- {
				password = append(password, set[random.ChoiceIndex(set)])
			}
		}

		p := make([]rune, len(password))
		for i, v := range rand.Perm(len(password)) {
			p[i] = password[v]
		}

		for i := 0; i < len(p); i++ {
			ch := p[i]
			if i+2 < len(p) && p[i+1] == ch && p[i+2] == ch {
				p = append(p[:i], p[i+1:]...)
				i--
			}
		}

		if len(p) < 8 {
			continue
		}

		pwd := String(p)

		return pwd
	}
}

func GenerateMachineName() String {
	return GeneratePinyin()
}

func GenerateAnswer() String {
	return GenerateName(3, 8)
}

func GenerateStreet1() String {
	// 街名和门牌号

	r := []string{"路", "街", "弄", "胡同", "大道"}
	s := []string{"号", "大院"}

	return String(Sprintf(
		"%s%s%d%s",
		GenerateLevel1Name(1, 3), //GenerateName(1, 3),
		random.Choice(r).(string),
		random.IntRange(1, 1000),
		random.Choice(s).(string),
	))
}

func GenerateStreet2() String {
	// 楼号、单元号、房间号

	b := []string{"栋", "区", "弄"}
	s := []string{"号", "单元", "街道", "校区", "宿舍", "教学楼", "号房"}
	i := []rune(`ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789`)

	return String(Sprintf(
		"%v%s%d%s",
		random.Choice(i).(rune),
		random.Choice(b).(string),
		random.IntRange(1, 1000),
		random.Choice(s).(string),
	))
}

func GenerateStreet3() String {
	// 街

	r := []string{"路", "街", "弄", "胡同", "大道"}

	return String(Sprintf(
		"%s%s",
		GenerateLevel1Name(2, 4), //GenerateName(2, 4),
		random.Choice(r).(string),
	))

	//return GenerateLevel1Name(4, 8) //GenerateName(4, 8)
}

func ParseHTML(html String) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(html.NewReader())
	RaiseIf(err)
	return doc
}

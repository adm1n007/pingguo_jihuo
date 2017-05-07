package register

import (
	. "fmt"
	. "ml"
	. "ml/array"
	. "ml/dict"
	. "ml/strings"
	. "ml/trace"

	"ml/random"
	urllib "net/url"
	"plistlib"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"globals"
)

func mergeDict(dict1, dict2 Dict) {
	for k, v := range dict2 {
		dict1[k] = v
	}
}

func mergeElement(e1, e2 elementDict) {
	for k, v := range e2 {
		e1[k] = v
	}
}

func peekProperty(text, property string) (value string) {
	pattern := regexp.MustCompile(Sprintf(`(?m:\.%s\s*=\s*(.*);?$)`, property))
	value = pattern.FindStringSubmatch(text)[1]
	return
}

func buildOrigin(url *urllib.URL) string {
	return Sprintf("%s://%s", url.Scheme, url.Host)
}

func buildNext(url *urllib.URL, next string) string {
	return Sprintf("%s://%s%s", url.Scheme, url.Host, next)
}

func selectOption(elem *goquery.Selection, argx ...Dict) string {
	filter := "option"
	attr := "value"
	maxIndex := 99999

	if len(argx) != 0 {
		args := argx[0]

		if v, ok := args["maxIndex"]; ok {
			maxIndex = v.(int)
		}

		if v, ok := args["filter"]; ok {
			filter = v.(string)
		}

		if v, ok := args["attr"]; ok {
			attr = v.(string)
		}
	}

	all := []*goquery.Selection{}
	elem.Find(filter).Each(func(i int, s *goquery.Selection) {
		all = append(all, s)
	})

	for len(all) != 0 {
		index := random.ChoiceIndex(all)
		if index >= maxIndex {
			all = append(all[:index], all[index+1:]...)
			continue
		}

		elem := all[index]
		value := elem.Attr2(attr)
		if _, err := strconv.Atoi(value); err != nil {
			all = append(all[:index], all[index+1:]...)
			continue
		}

		return value
	}

	return ""
}

func getAllInput(doc *goquery.Selection) *goquery.Selection {
	return doc.Find("input").Filter("[type=text],[type=password],[type=email],[type=url],[type=search],[type=tel]")
}

func getInputValue(s *goquery.Selection) string {
	return s.AttrOr("value", "")
}

func getInputName(s *goquery.Selection) string {
	id := s.AttrOr("id", "")
	return If(id == "", s.Attr2("name"), id).(string)
}

func getAllElement(doc *goquery.Selection, ids Array) map[string]*goquery.Selection {
	data := map[string]*goquery.Selection{}

	doc.Find("[id]").Each(func(i int, s *goquery.Selection) {
		i, exist := ids.Index(s.Attr2("id"))
		if exist {
			data[ids[i].(string)] = s
		}
	})

	return data
}

type POSITION struct {
	xmin, xmax int
	ymin, ymax int
}

var (
	tags = map[string]string{
		"timeSlice":           "ts",
		"start":               "st",
		"focusAndValueLength": "kk",
		"focusOnly":           "ff",
		"mouseMove":           "mm",
		"mouseClick":          "mc",
		"keyDown":             "kd",
		"focusBlur":           "fb",
		"submit":              "fs",
		"touchstart":          "te",
	}

	pcpositionChina = map[string]*POSITION{
		"codeRedemptionField": &POSITION{120, 420, 370, 390},
		"lastFirstName":       &POSITION{120, 260, 480, 500},
		"firstName":           &POSITION{280, 420, 480, 500},
		"street1":             &POSITION{120, 420, 515, 535},
		"street2":             &POSITION{120, 420, 550, 570},
		"street3":             &POSITION{120, 420, 585, 605},
		"city":                &POSITION{120, 240, 620, 640},
		"postalcode":          &POSITION{120, 215, 655, 675},
		"state":               &POSITION{230, 350, 655, 675},
		"phone1Number":        &POSITION{120, 240, 690, 710},
		"continue":            &POSITION{920, 1020, 850, 870},
	}

	pcpositionIndia = map[string]*POSITION{
		"codeRedemptionField": &POSITION{120, 420, 370, 390},
		"salutation":          &POSITION{280, 420, 480, 500},
		"firstName":           &POSITION{280, 420, 480, 500},
		"lastName":            &POSITION{120, 260, 480, 500},
		"street1":             &POSITION{120, 420, 515, 535},
		"street2":             &POSITION{120, 420, 550, 570},
		"street3":             &POSITION{120, 420, 585, 605},
		"city":                &POSITION{120, 240, 620, 640},
		"postalcode":          &POSITION{120, 215, 655, 675},
		"state":               &POSITION{230, 350, 655, 675},
		"phone1AreaCode":      &POSITION{120, 240, 690, 710},
		"phone1Number":        &POSITION{120, 240, 690, 710},
		"continue":            &POSITION{920, 1020, 850, 870},
	}

	/**
	  [Title]
	  [First Name] [Last Name]
	  [Cyty]
	  [Street1               ]
	  [Street2               ]
	  [postalcode] Taiwan
	  [AreaCode][Phone  ]
	*/

	pcpositionTaiwan = map[string]*POSITION{
		// "codeRedemptionField"   : &POSITION{120, 420, 370, 390},

		"salutation":    &POSITION{120, 170, 480, 500},
		"lastFirstName": &POSITION{120, 260, 515, 535}, "firstName": &POSITION{280, 420, 515, 535},
		"citypopup":      &POSITION{120, 170, 550, 570},
		"street1":        &POSITION{120, 420, 585, 605},
		"street2":        &POSITION{120, 420, 620, 640},
		"postalcode":     &POSITION{120, 200, 655, 675},
		"phone1AreaCode": &POSITION{120, 210, 690, 710}, "phone1Number": &POSITION{210, 300, 690, 710},

		"continue": &POSITION{920, 1020, 850, 870},
	}

	/*
	   [Title]
	   [First Name] [Last Name]
	   [Street                ]
	   [Street2               ]
	   [Suburb         ]
	   [postalcode] [City ] NewZealand
	   [AreaCode][Phone  ]
	*/

	pcpositionNewZealand = map[string]*POSITION{
		"codeRedemptionField": &POSITION{120, 420, 370, 390},
		"salutation":          &POSITION{120, 150, 480, 500},
		"firstName":           &POSITION{120, 260, 515, 535},
		"lastName":            &POSITION{280, 420, 515, 535},
		"street1":             &POSITION{120, 420, 550, 570},
		"street2":             &POSITION{120, 420, 585, 605},
		"suburb":              &POSITION{120, 240, 620, 640},
		"postalcode":          &POSITION{120, 200, 655, 675},
		"city":                &POSITION{240, 300, 655, 675},
		"phone1AreaCode":      &POSITION{120, 180, 690, 710},
		"phone1Number":        &POSITION{200, 250, 690, 710},
		"continue":            &POSITION{920, 1020, 850, 870},
	}

	pcpositionVietnam = map[string]*POSITION{
		"codeRedemptionField": &POSITION{120, 420, 370, 390},
		"salutation":          &POSITION{120, 150, 480, 500},
		"firstName":           &POSITION{120, 260, 515, 535},
		"lastName":            &POSITION{280, 420, 515, 535},
		"street1":             &POSITION{120, 420, 550, 570},
		"street2":             &POSITION{120, 420, 585, 605},
		"city":                &POSITION{120, 200, 620, 640},
		"postalcode":          &POSITION{240, 300, 620, 640},
		"phone1AreaCode":      &POSITION{120, 180, 655, 675},
		"phone1Number":        &POSITION{200, 250, 655, 675},
		"continue":            &POSITION{920, 1020, 850, 870},
	}

	iosposition = map[string]*POSITION{
		"":                 &POSITION{93, 291, 179, 193},
		"lastNameField":    &POSITION{146, 300, 555, 568},
		"firstNameField":   &POSITION{146, 300, 600, 610},
		"street1Field":     &POSITION{146, 300, 640, 660},
		"street2Field":     &POSITION{146, 300, 687, 707},
		"street3Field":     &POSITION{146, 300, 737, 750},
		"postalCodeField":  &POSITION{146, 300, 776, 795},
		"cityField":        &POSITION{146, 300, 822, 838},
		"stateField":       &POSITION{146, 300, 869, 884},
		"phoneNumberField": &POSITION{146, 300, 922, 936},
	}

	iospositionNone = map[string]*POSITION{
		"":                 &POSITION{0, 320, 377, 424}, // don't care?
		"lastNameField":    &POSITION{136, 297, 697, 721},
		"firstNameField":   &POSITION{136, 297, 742, 766},
		"street1Field":     &POSITION{136, 297, 788, 811},
		"street2Field":     &POSITION{136, 297, 833, 856},
		"street3Field":     &POSITION{136, 297, 877, 901},
		"postalCodeField":  &POSITION{136, 297, 923, 946},
		"cityField":        &POSITION{136, 297, 967, 991},
		"stateField":       &POSITION{136, 304, 1012, 1036},
		"phoneNumberField": &POSITION{136, 297, 1067, 1091},
	}
)

func init() {
	for _, v := range pcpositionChina {
		v.ymin -= 23
		v.ymax -= 23
	}

	for _, v := range pcpositionIndia {
		v.ymin -= 23
		v.ymax -= 23
	}

	for _, v := range pcpositionNewZealand {
		v.ymin -= 23
		v.ymax -= 23
	}
}

func findInputByPosition(x, y int, position map[string]*POSITION) String {
	for n, v := range position {
		if x > v.xmin && x < v.xmax && y > v.ymin && y < v.ymax {
			return String(n)
		}
	}

	return ""
}

func getRandomPosition(position map[string]*POSITION) (x, y int, name String) {
	name = ""
	for name == "continue" || name == "" {
		x = random.IntRange(50, 800)
		y = random.IntRange(50, 800)
		name = findInputByPosition(x, y, position)
	}

	return
}

func getRandomPositionForItem(pos *POSITION) (x, y int) {
	x = random.IntRange(pos.xmin, pos.xmax)
	y = random.IntRange(pos.ymin, pos.ymax)

	return
}

func generateInputRecordsChina(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) String {
	return generateInputRecordsWorker(doc, data, form, pcpositionChina)
}

func generateInputRecordsTaiwan(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) String {
	return generateInputRecordsWorker(doc, data, form, pcpositionTaiwan)
}

func generateInputRecordsIndia(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) String {
	return generateInputRecordsWorker(doc, data, form, pcpositionIndia)
}

func generateInputRecordsNewZealand(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) String {
	return generateInputRecordsWorker(doc, data, form, pcpositionNewZealand)
}

func generateInputRecordsVietnam(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) String {
	return generateInputRecordsWorker(doc, data, form, pcpositionVietnam)
}

func generateInputRecordsWorker(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict, position map[string]*POSITION) String {
	records := []String{}

	var timeFactor float64 = 1
	var lastMoveMouseTime int64 = 0

	var currentTime int64 = globals.GetCurrentTime()

	// currentTime -= 15 * 1000

	var lastTimeslice int64 = currentTime
	var initTimestamp int64 = currentTime
	var lastTimestamp int64 = currentTime

	var addIprItem2 func(tag string, timestamp int64, items *Array)
	var addIprItem func(tag string, items *Array, delta int, dontMoveMouse bool)

	addIprItem2 = func(tag string, timestamp int64, items *Array) {
		elapsed := timestamp - lastTimestamp
		lastTimestamp = timestamp

		if timeFactor > 1 {
			elapsed = int64(float64(elapsed)/timeFactor + 0.5)
		}

		v := []string{
			tag,
			Sprintf("%x", elapsed),
		}

		if items != nil {
			for _, e := range *items {
				switch t := e.(type) {
				case int:
					v = append(v, Sprintf("%x", int(t)))

				case int64:
					v = append(v, Sprintf("%x", int(t)))

				case float32:
					v = append(v, Sprintf("%x", int(t)))

				case float64:
					v = append(v, Sprintf("%x", int(t)))

				default:
					v = append(v, Sprint(t))
				}
			}
		}

		records = append(records, String(",").Join(v))

		if currentTime-lastTimeslice > 15000 {
			lastTimeslice = currentTime
			addIprItem2(tags["timeSlice"], currentTime, &Array{currentTime - initTimestamp})
		}
	}

	addIprItem = func(tag string, items *Array, delta int, dontMoveMouse bool) {
		var min int

		if tags["mouseClick"] == tag {
			min = 0x500
		} else {
			min = 50
		}

		max := 1500

		if delta < 0 {
			delta = random.IntRange(min, max)
		}

		if tag != tags["mouseMove"] {
			currentTime += int64(delta)
			addIprItem2(tag, currentTime, items)
		}

		if dontMoveMouse == false && currentTime-lastMoveMouseTime > 5000 {
			switch {
			case tag == tags["mouseMove"]:
				currentTime += int64(delta)
				lastMoveMouseTime = currentTime
				addIprItem2(tag, currentTime, items)

			case random.IntRange(0, 5) < 1:
				if random.IntRange(0, 2) == 1 {
					ps := []*POSITION{}
					pname := []string{}
					for n, v := range position {
						ps = append(ps, v)
						pname = append(pname, n)
					}
					index := random.ChoiceIndex(ps)
					pos := ps[index]
					addIprItem(
						tags["mouseMove"],
						&Array{
							random.IntRange(pos.xmin, pos.xmax),
							random.IntRange(pos.ymin, pos.ymax),
							pname[index],
						},
						-1,
						false,
					)
				} else {
					x, y, name := getRandomPosition(position)
					items := &Array{x, y}
					if name.IsEmpty() == false {
						items.Append(name)
					}
					addIprItem(tags["mouseMove"], items, -1, false)
				}
			}
		}
	}

	inputs := getAllInput(doc)
	inputArray := Array{}
	st := Array{}

	inputs.Each(func(i int, s *goquery.Selection) {
		st.Append(Sprintf("%s,%s", getInputName(s), "0"))
		inputArray.Append(s)
		// if getInputName(s) == "postalcode" {
		//     if state, exist := data["state"]; exist {
		//         inputArray.Append(state)
		//     }
		// }
	})

	{
		addIprItem2("ncip", currentTime, &Array{globals.GetCurrentTime() / 1000, 1, 1})
		addIprItem(tags["start"], &st, 0, true)

		x, y, name := getRandomPosition(position)
		items := &Array{x, y}
		if name.IsEmpty() == false {
			items.Append(name)
		}
		addIprItem(tags["mouseMove"], items, random.IntRange(0x10, 0x100), false)

		currentTime += int64(random.IntRange(800, 4000))
	}

	// inputArray = random.Shuffle(inputArray)
	for i, e := range inputArray {
		s := e.(*goquery.Selection)

		name := getInputName(s)

		Println("name = ", name)

		if _, exist := data[name]; exist == false {
			continue
		}

		pos := position[name]
		if pos == nil {
			continue
		}

		if s.Nodes[0].Data != "input" && name == "state" {
			addIprItem(
				tags["mouseMove"],
				&Array{
					random.IntRange(pos.xmin, pos.xmax),
					random.IntRange(pos.ymin, pos.ymax),
					name,
				},
				-1,
				false,
			)

			addIprItem(
				tags["mouseClick"],
				&Array{
					random.IntRange(pos.xmin, pos.xmax),
					random.IntRange(pos.ymin, pos.ymax),
					name,
				},
				-1,
				false,
			)

			addIprItem(
				tags["mouseMove"],
				&Array{
					random.IntRange(pos.xmin, pos.xmax),
					random.IntRange(pos.ymax+10, pos.ymax+210),
				},
				-1,
				false,
			)

			addIprItem(
				tags["mouseClick"],
				&Array{
					random.IntRange(pos.xmin, pos.xmax),
					random.IntRange(pos.ymax+10, pos.ymax+210),
				},
				-1,
				false,
			)

			continue
		}

		valueLength := len(Sprintf("%v", form[data[name].Attr2("name")]))

		if valueLength == 0 && random.IntRange(0, 10) < 5 {
			continue
		}

		useMouse := i == 0 || random.IntRange(0, 10) < 5

		if useMouse {
			addIprItem(
				tags["mouseMove"],
				&Array{
					random.IntRange(pos.xmin, pos.xmax),
					random.IntRange(pos.ymin, pos.ymax),
					name,
				},
				-1,
				false,
			)
		}

		addIprItem(tags["focusAndValueLength"], &Array{0, name}, random.IntRange(0, 2), true)
		addIprItem(tags["focusOnly"], &Array{name}, random.IntRange(0, 2), false)

		if useMouse {
			addIprItem(
				tags["mouseClick"],
				&Array{
					random.IntRange(pos.xmin, pos.xmax),
					random.IntRange(pos.ymin, pos.ymax),
					name,
				},
				random.IntRange(0x40, 0x70),
				false,
			)
		}

		currentTime += int64(random.IntRange(100, 4000))

		if valueLength != 0 {
			d := random.IntRange(0x300, 0x7400)
			for i := random.IntRange(valueLength/1, valueLength*3/2); i != 0; i-- {
				addIprItem(tags["keyDown"], nil, d, false)
				d = random.IntRange(0x30, 0xA0)
			}
		}

		currentTime += int64(random.IntRange(800, 1500))

		if i+1 != inputArray.Length() {
			addIprItem(tags["focusBlur"], &Array{name}, random.IntRange(1, 4), true)
			continue
		}

		// {
		//     x, y, name := getRandomPosition()
		//     items := &Array{x, y}
		//     if name.IsEmpty() == false {
		//         items.Append(name)
		//     }
		//     addIprItem(tags["mouseMove"], items, -1, false)
		// }

		ctnu := position["continue"]

		addIprItem(
			tags["mouseMove"],
			&Array{
				random.IntRange(ctnu.xmin, ctnu.xmax),
				random.IntRange(ctnu.ymin, ctnu.ymax),
				data["continue"].Attr2("name"),
			},
			random.IntRange(500, 1300),
			true,
		)

		addIprItem(tags["focusBlur"], &Array{name}, random.IntRange(1500, 2000), true)

		addIprItem(
			tags["mouseClick"],
			&Array{
				random.IntRange(ctnu.xmin, ctnu.xmax),
				random.IntRange(ctnu.ymin, ctnu.ymax),
				data["continue"].Attr2("name"),
			},
			random.IntRange(0x40, 0x100),
			true,
		)
		addIprItem(tags["submit"], &Array{0, 0}, random.IntRange(3, 7), false)

		break
	}

	return String(";").Join(records) + ",;"
}

func generateiOSInputRecords(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) (String, int64) {
	position := iospositionNone

	records := []String{}

	var timeFactor float64 = 1
	var lastMoveMouseTime int64 = 0
	var currentTime int64 = globals.GetCurrentTime()

	// currentTime -= 15 * 1000

	var lastTimeslice int64 = currentTime
	var initTimestamp int64 = currentTime
	var lastTimestamp int64 = currentTime

	var addIprItem2 func(tag string, timestamp int64, items *Array)
	var addIprItem func(tag string, items *Array, delta int, dontMoveMouse bool)

	addIprItem2 = func(tag string, timestamp int64, items *Array) {
		elapsed := timestamp - lastTimestamp
		lastTimestamp = timestamp

		if timeFactor > 1 {
			elapsed = int64(float64(elapsed)/timeFactor + 0.5)
		}

		if tag == "ncip" {
			elapsed = -1
		}

		v := []string{
			tag,
			Sprintf("%x", elapsed),
		}

		if items != nil {
			for _, e := range *items {
				switch t := e.(type) {
				case int:
					v = append(v, Sprintf("%x", int(t)))

				case int64:
					v = append(v, Sprintf("%x", int(t)))

				case float32:
					v = append(v, Sprintf("%x", int(t)))

				case float64:
					v = append(v, Sprintf("%x", int(t)))

				default:
					v = append(v, Sprint(t))
				}
			}
		}

		records = append(records, String(",").Join(v))

		if currentTime-lastTimeslice > 15000 {
			lastTimeslice = currentTime
			addIprItem2(tags["timeSlice"], currentTime, &Array{currentTime - initTimestamp})
		}
	}

	addIprItem = func(tag string, items *Array, delta int, dontMoveMouse bool) {
		var min int

		if tags["mouseClick"] == tag {
			min = 0x500
		} else {
			min = 50
		}

		max := 1500

		if delta < 0 {
			delta = random.IntRange(min, max)
		}

		if tag != tags["mouseMove"] {
			currentTime += int64(delta)
			addIprItem2(tag, currentTime, items)
		}

		if dontMoveMouse == false && currentTime-lastMoveMouseTime > 5000 {
			switch {
			case tag == tags["mouseMove"]:
				currentTime += int64(delta)
				lastMoveMouseTime = currentTime
				addIprItem2(tag, currentTime, items)

			case random.IntRange(0, 5) < 1:
				break
				if random.IntRange(0, 2) == 1 {
					ps := []*POSITION{}
					pname := []string{}
					for n, v := range position {
						ps = append(ps, v)
						pname = append(pname, n)
					}
					index := random.ChoiceIndex(ps)
					pos := ps[index]
					addIprItem(
						tags["mouseMove"],
						&Array{
							random.IntRange(pos.xmin, pos.xmax),
							random.IntRange(pos.ymin, pos.ymax),
							pname[index],
						},
						-1,
						false,
					)
				} else {
					x, y, name := getRandomPosition(position)
					items := &Array{x, y}
					if name.IsEmpty() == false {
						items.Append(name)
					}
					addIprItem(tags["mouseMove"], items, -1, false)
				}
			}
		}
	}

	inputs := getAllInput(doc)
	inputArray := Array{}
	st := Array{}

	inputs.Each(func(i int, s *goquery.Selection) {
		st.Append(Sprintf("%s,%s", getInputName(s), "0"))
		inputArray.Append(s)
		if getInputName(s) == "cityField" {
			inputArray.Append(data["stateField"])
		}
	})

	{
		addIprItem2("ncip", currentTime, &Array{globals.GetCurrentTime() / 1000, 1, 1})
		addIprItem(tags["start"], &st, 0, true)
		addIprItem(tags["touchstart"], &Array{-1, -1, "VISA"}, random.IntRange(16, 500), false)
		addIprItem(tags["touchstart"], &Array{-1, -1}, random.IntRange(16, 500), false)
		addIprItem(tags["touchstart"], &Array{-1, -1}, random.IntRange(16, 500), false)
		addIprItem(tags["touchstart"], &Array{-1, -1}, random.IntRange(16, 500), false)

		// bug?!
		//currentTime += int64(random.IntRange(512, 1024)) // delay to lastNameField
	}

	var prevName string

	for i, e := range inputArray {
		s := e.(*goquery.Selection)

		name := getInputName(s)

		if _, exist := data[name]; exist == false {
			continue
		}

		if s.Nodes[0].Data != "input" && name == "stateField" {
			pos := position[name]
			x, y := getRandomPositionForItem(pos)

			addIprItem(tags["touchstart"], &Array{-1, -1, name}, -1, false)
			addIprItem(tags["mouseMove"], &Array{x, y, name}, random.IntRange(0x10, 0x70), false)
			addIprItem(tags["focusBlur"], &Array{prevName}, random.IntRange(48, 174), false)
			addIprItem(tags["mouseClick"], &Array{x, y, name}, random.IntRange(0, 0x20), false)

			prevName = ""
			continue
		}

		valueLength := len(Sprintf("%v", form[data[name].Attr2("name")]))

		if valueLength == 0 {
			continue
		}

		pos := position[name]
		x, y := getRandomPositionForItem(pos)

		addIprItem(tags["touchstart"], &Array{-1, -1, name}, -1, false)
		addIprItem(tags["mouseMove"], &Array{x, y, name}, random.IntRange(0, 0x70), false)

		if prevName != "" {
			addIprItem(tags["focusBlur"], &Array{prevName}, random.IntRange(48, 174), true)
			prevName = ""
		}

		addIprItem(tags["focusAndValueLength"], &Array{0, name}, random.IntRange(2, 174), true)
		addIprItem(tags["focusOnly"], &Array{name}, random.IntRange(0, 1), false)
		addIprItem(tags["mouseClick"], &Array{x, y, name}, random.IntRange(0, 0x20), false)

		//currentTime += int64(random.IntRange(100, 4000))

		if valueLength != 0 {
			d := random.IntRange(0x400, 0x2200) // added?
			for i := valueLength + random.IntRange(0, 5); i != 0; i-- {
				addIprItem(tags["keyDown"], nil, d, false)
				d = random.IntRange(0x80, 0x520)
			}
		}

		//currentTime += int64(random.IntRange(800, 1500))

		if i+1 != inputArray.Length() {
			prevName = name
			continue
		}

		addIprItem(tags["focusBlur"], &Array{name}, random.IntRange(562, 3000), true)
		addIprItem(tags["mouseClick"], &Array{0, 0, "hiddenBottomRightButtonId"}, random.IntRange(0x40, 0x800), true)
		addIprItem(tags["submit"], &Array{0, 0}, random.IntRange(3, 15), false)

		break
	}

	return String(";").Join(records) + ",;", currentTime - initTimestamp
}

func generateiOSInputRecordsRetry(doc *goquery.Selection, data map[string]*goquery.Selection, form Dict) (String, int64) {
	position := iospositionNone

	records := []String{}

	var timeFactor float64 = 1
	var lastMoveMouseTime int64 = 0
	var currentTime int64 = globals.GetCurrentTime()

	// currentTime -= 15 * 1000

	var lastTimeslice int64 = currentTime
	var initTimestamp int64 = currentTime
	var lastTimestamp int64 = currentTime

	var addIprItem2 func(tag string, timestamp int64, items *Array)
	var addIprItem func(tag string, items *Array, delta int, dontMoveMouse bool)

	addIprItem2 = func(tag string, timestamp int64, items *Array) {
		elapsed := timestamp - lastTimestamp
		lastTimestamp = timestamp

		if timeFactor > 1 {
			elapsed = int64(float64(elapsed)/timeFactor + 0.5)
		}

		if tag == "ncip" {
			elapsed = -1
		}

		v := []string{
			tag,
			Sprintf("%x", elapsed),
		}

		if items != nil {
			for _, e := range *items {
				switch t := e.(type) {
				case int:
					v = append(v, Sprintf("%x", int(t)))

				case int64:
					v = append(v, Sprintf("%x", int(t)))

				case float32:
					v = append(v, Sprintf("%x", int(t)))

				case float64:
					v = append(v, Sprintf("%x", int(t)))

				default:
					v = append(v, Sprint(t))
				}
			}
		}

		records = append(records, String(",").Join(v))

		if currentTime-lastTimeslice > 15000 {
			lastTimeslice = currentTime
			addIprItem2(tags["timeSlice"], currentTime, &Array{currentTime - initTimestamp})
		}
	}

	addIprItem = func(tag string, items *Array, delta int, dontMoveMouse bool) {
		var min int

		if tags["mouseClick"] == tag {
			min = 0x500
		} else {
			min = 50
		}

		max := 1500

		if delta < 0 {
			delta = random.IntRange(min, max)
		}

		if tag != tags["mouseMove"] {
			currentTime += int64(delta)
			addIprItem2(tag, currentTime, items)
		}

		if dontMoveMouse == false && currentTime-lastMoveMouseTime > 5000 {
			switch {
			case tag == tags["mouseMove"]:
				currentTime += int64(delta)
				lastMoveMouseTime = currentTime
				addIprItem2(tag, currentTime, items)

			case random.IntRange(0, 5) < 1:
				break
				if random.IntRange(0, 2) == 1 {
					ps := []*POSITION{}
					pname := []string{}
					for n, v := range position {
						ps = append(ps, v)
						pname = append(pname, n)
					}
					index := random.ChoiceIndex(ps)
					pos := ps[index]
					addIprItem(
						tags["mouseMove"],
						&Array{
							random.IntRange(pos.xmin, pos.xmax),
							random.IntRange(pos.ymin, pos.ymax),
							pname[index],
						},
						-1,
						false,
					)
				} else {
					x, y, name := getRandomPosition(position)
					items := &Array{x, y}
					if name.IsEmpty() == false {
						items.Append(name)
					}
					addIprItem(tags["mouseMove"], items, -1, false)
				}
			}
		}
	}

	inputs := getAllInput(doc)
	//inputArray := Array{}
	st := Array{}

	inputs.Each(func(i int, s *goquery.Selection) {
		name := getInputName(s)
		valueLength := len([]rune(Sprintf("%v", form[data[name].Attr2("name")])))
		st.Append(Sprintf("%s,%d", name, valueLength))
	})

	{
		addIprItem2("ncip", currentTime, &Array{globals.GetCurrentTime() / 1000, 1, 1})
		addIprItem(tags["start"], &st, 0, true)
		//addIprItem(tags["touchstart"], &Array{-1, -1}, random.IntRange(0x10, 0x100), false)

		//currentTime += int64(random.IntRange(800, 4000))
	}

	mcdelay := random.IntRange(0x1500, 0x5000)
	addIprItem(tags["mouseClick"], &Array{0, 0, "hiddenBottomRightButtonId"}, mcdelay, true)
	addIprItem(tags["submit"], &Array{0, 0}, random.IntRange(3, 15), false)

	return String(";").Join(records) + ",;", currentTime - initTimestamp
}

func parsePlist(text string) Dict {
	plist := Dict{}

	RaiseIf(plistlib.Unmarshal(
		String("<plist version=\"1.0\">\r\n"+text+"</plist>").Encode(CP_UTF8),
		&plist,
	))

	return plist
}

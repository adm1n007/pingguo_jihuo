package inputRecord

import (
	"fmt"
	. "ml/dict"
	"ml/logging/logger"
	"ml/random"
	. "ml/strings"
	"time"

	"../../htmlutils"
	"github.com/PuerkitoBio/goquery"
)

type Options struct {
	TimeIncrement time.Duration
	DontMoveMouse bool
}

const (
	TimeAutoIncrement time.Duration = -1
)

type Generator struct {
	layout          ElementLayout
	Records         []*Record
	lookup          map[string]*ElementPosition
	randomMoveMouse bool

	initialTime       time.Time
	currentTime       time.Time
	lastTime          time.Time
	lastTimeslice     time.Time
	lastMoveMouseTime time.Time
}

func New(layout ElementLayout) *Generator {
	now := time.Now()

	lookup := map[string]*ElementPosition{}

	for _, e := range layout {
		lookup[e.Name.String()] = e
	}

	return &Generator{
		layout:          layout,
		lookup:          lookup,
		randomMoveMouse: true,
		initialTime:     now,
		currentTime:     now,
		lastTime:        now,
		lastTimeslice:   now,
	}
}

func (self *Generator) increseTimestamp(increment time.Duration) {
	if increment != 0 {
		self.currentTime = self.currentTime.Add(increment)
	}
}

func (self *Generator) addRecordWorder(record *Record, options *Options) {
	record.Elapsed = self.currentTime.Sub(self.lastTime)
	self.lastTime = self.currentTime

	self.Records = append(self.Records, record)

	if self.currentTime.Sub(self.lastTimeslice) > 15*time.Second {
		self.lastTimeslice = self.currentTime
		self.addRecordWorder(
			&Record{
				Type:      TimeSlice,
				Timestamp: int64(self.currentTime.Sub(self.initialTime) / time.Millisecond),
			},
			options,
		)
	}
}

func (self *Generator) AddRecord(record *Record, options *Options) {
	if options == nil {
		options = &Options{
			TimeIncrement: TimeAutoIncrement,
		}
	}

	if options.TimeIncrement == TimeAutoIncrement {
		max := 1500
		min := 50

		if record.Type == MouseClick {
			min = 0x500
		}

		options.TimeIncrement = randomTime(min, max)
	}

	if record.Type != MouseMove {
		self.increseTimestamp(options.TimeIncrement)
		self.addRecordWorder(record, options)
	}

	if options.DontMoveMouse {
		return
	}

	if self.currentTime.Sub(self.lastMoveMouseTime) <= 5*time.Second {
		return
	}

	if record.Type == MouseMove {
		self.increseTimestamp(options.TimeIncrement)
		self.lastMoveMouseTime = self.currentTime
		self.addRecordWorder(record, options)

		return
	}

	if self.randomMoveMouse == false || random.IntRange(0, 5) >= 1 {
		return
	}

	mustHover := random.IntRange(0, 2) == 0

	for {
		x, y, element := self.getRandomPosition()
		if mustHover && element.IsEmpty() {
			continue
		}

		self.moveMouse(x, y, element, nil)
		break
	}
}

func (self *Generator) findInputForPosition(x, y int) String {
	for _, e := range self.layout {
		if x > e.Left && x < e.Right && y > e.Top && y < e.Bottom {
			return e.Name
		}
	}

	return ""
}

func (self *Generator) getRandomPosition() (x, y int, element String) {
	element = "continue"
	for element == "continue" {
		x = random.IntRange(50, 800)
		y = random.IntRange(50, 800)

		element = self.findInputForPosition(x, y)
	}

	return
}

func (self *Generator) recordsToString(records []*Record) String {
	ipr := []String{}

	for _, r := range records {
		ipr = append(ipr, String(r.String()))
	}

	return String(";").Join(ipr) + ";"
}

func (self *Generator) GenerateRetry(elements map[string]*goquery.Selection, inputs []*goquery.Selection, form Dict) String {
	self.ncip()
	self.start(inputs)

	x, y, element := self.getRandomPosition()
	self.moveMouse(x, y, element,
		&Options{
			TimeIncrement: randomTime(0x10, 0x100),
		},
	)

	next := self.lookup["continue"]
	continuebtn := elements["continue"]

	self.clickMouse(
		random.IntRange(next.Left, next.Right),
		random.IntRange(next.Top, next.Bottom),
		String(continuebtn.Attr2("name")),
		&Options{
			TimeIncrement: randomTime(0x40, 0x100),
			DontMoveMouse: true,
		},
	)

	self.submit(randomTime(3, 7))

	return self.recordsToString(self.Records)
}

func (self *Generator) Generate(elements map[string]*goquery.Selection, inputs []*goquery.Selection, form Dict) String {
	self.ncip()
	self.start(inputs)

	x, y, element := self.getRandomPosition()
	self.moveMouse(x, y, element,
		&Options{
			TimeIncrement: randomTime(0x10, 0x100),
		},
	)

	self.increseTimestamp(randomTime(800, 4000))

	for index, elem := range inputs {
		name := htmlutils.GetInputName(elem)

		if _, exists := form[elem.Attr2("name")]; exists == false {
			continue
		}

		e := self.lookup[name]
		if e == nil {
			logger.Debug("can't find %v", name)
			continue
		}

		logger.Debug("search %v: %+v", name, e)

		if elem.Nodes[0].Data != "input" {
			self.selectListBox(self.lookup[name])
			continue
		}

		valueLength := len(fmt.Sprintf("%v", form[elem.Attr2("name")]))
		if valueLength == 0 && random.IntRange(0, 10) < 5 {
			continue
		}

		useMouse := index == 0 || random.IntRange(0, 10) < 5

		if useMouse {
			self.moveMouse(
				random.IntRange(e.Left, e.Right),
				random.IntRange(e.Top, e.Bottom),
				e.Name,
				nil,
			)
		}

		self.focusAndValueLength(e)
		self.focusOnly(e, randomTime(0, 2), false)

		if useMouse {
			self.clickMouse(
				random.IntRange(e.Left, e.Right),
				random.IntRange(e.Top, e.Bottom),
				e.Name,
				&Options{
					TimeIncrement: randomTime(0x40, 0x70),
				},
			)
		}

		self.increseTimestamp(randomTime(100, 4000))

		if valueLength != 0 {
			d := randomTime(0x300, 0x7400)
			for i := random.IntRange(valueLength/1, valueLength*2/2); i != 0; i-- {
				self.keyDown(d)
				d = randomTime(0x30, 0x300)
			}
		}

		self.increseTimestamp(randomTime(800, 1500))

		if index+1 != len(inputs) {
			self.focusBlur(e, randomTime(1, 4), true)
			continue
		}

		next := self.lookup["continue"]
		continuebtn := elements["continue"]

		self.moveMouse(
			random.IntRange(next.Left, next.Right),
			random.IntRange(next.Top, next.Bottom),
			String(continuebtn.Attr2("name")),
			&Options{
				TimeIncrement: randomTime(500, 1300),
				DontMoveMouse: true,
			},
		)

		self.focusBlur(e, randomTime(1500, 2000), true)

		self.clickMouse(
			random.IntRange(next.Left, next.Right),
			random.IntRange(next.Top, next.Bottom),
			String(continuebtn.Attr2("name")),
			&Options{
				TimeIncrement: randomTime(0x40, 0x100),
				DontMoveMouse: true,
			},
		)

		self.submit(randomTime(3, 7))

		break
	}

	ipr := []String{}

	for _, r := range self.Records {
		ipr = append(ipr, String(r.String()))
	}

	return String(";").Join(ipr) + ";"
}

func (self *Generator) GenerateiOS(elements map[string]*goquery.Selection, inputs []*goquery.Selection, form Dict) String {
	self.randomMoveMouse = false

	self.ncip()
	self.start(inputs)

	self.increseTimestamp(randomTime(800, 4000))

	for index, elem := range inputs {
		if _, exists := form[elem.Attr2("name")]; exists == false {
			continue
		}

		name := htmlutils.GetInputName(elem)
		e := self.lookup[name]
		logger.Debug("search %v: %+v", name, e)

		if elem.Nodes[0].Data != "input" {
			self.selectListBox(self.lookup[name])
			self.increseTimestamp(randomTime(0x500, 0xA00))
			continue
		}

		valueLength := len(fmt.Sprintf("%v", form[elem.Attr2("name")]))
		if valueLength == 0 {
			continue
		}

		self.moveMouse(
			random.IntRange(e.Left, e.Right),
			random.IntRange(e.Top, e.Bottom),
			e.Name,
			nil,
		)

		self.focusAndValueLength(e)
		self.focusOnly(e, randomTime(0, 2), false)

		self.clickMouse(
			random.IntRange(e.Left, e.Right),
			random.IntRange(e.Top, e.Bottom),
			e.Name,
			&Options{
				TimeIncrement: randomTime(0x5, 0x30),
			},
		)

		self.increseTimestamp(randomTime(100, 4000))

		d := randomTime(0x300, 0x7400)
		for i := random.IntRange(valueLength, valueLength*2); i != 0; i-- {
			self.keyDown(d)
			d = randomTime(0x30, 0x300)
		}

		self.increseTimestamp(randomTime(800, 1500))

		if index+1 != len(inputs) {
			nextElement := self.lookup[htmlutils.GetInputName(inputs[index+1])]

			self.moveMouse(
				random.IntRange(nextElement.Left, nextElement.Right),
				random.IntRange(nextElement.Top, nextElement.Bottom),
				nextElement.Name,
				&Options{
					TimeIncrement: randomTime(0x180, 0x250),
				},
			)

			self.focusBlur(e, randomTime(0x20, 0x200), true)
			continue
		}

		next := self.lookup["hiddenBottomRightButtonId"]
		continuebtn := elements["hiddenBottomRightButtonId"]

		self.moveMouse(
			random.IntRange(next.Left, next.Right),
			random.IntRange(next.Top, next.Bottom),
			String(continuebtn.Attr2("name")),
			&Options{
				TimeIncrement: randomTime(500, 1300),
				DontMoveMouse: true,
			},
		)

		self.focusBlur(e, randomTime(1500, 2000), true)

		logger.Debug("next btn %+v", next)

		self.clickMouse(
			random.IntRange(next.Left, next.Right),
			random.IntRange(next.Top, next.Bottom),
			String(continuebtn.Attr2("id")),
			&Options{
				TimeIncrement: randomTime(0x30, 0x800),
				DontMoveMouse: true,
			},
		)

		self.submit(randomTime(3, 7))

		break
	}

	ipr := []String{}

	for _, r := range self.Records {
		ipr = append(ipr, String(r.String()))
	}

	return String(";").Join(ipr) + ";"
}

func (self *Generator) GenerateiOSRetry(elements map[string]*goquery.Selection, inputs []*goquery.Selection, form Dict) String {
	// retry1
	//ncip,0,58b1e8b4,1,1;
	//st,0,sesame-id-input,0,national-id-input,0,creditCardNumberField,0,verificationNumberField,0,mobileNumberField,11,codeRedemptionField,0,lastNameField,7,firstNameField,5,street1Field,16,street2Field,6,street3Field,0,postalCodeField,6,cityField,4,phoneNumberField,11;
	//mc,2e20,0,0,hiddenBottomRightButtonId;
	//fs,4,0,0,;
	// retry2
	//ncip,-1,58b1e8df,1,1;
	//st,0,sesame-id-input,0,national-id-input,0,creditCardNumberField,0,verificationNumberField,0,mobileNumberField,11,codeRedemptionField,0,lastNameField,7,firstNameField,5,street1Field,16,street2Field,6,street3Field,0,postalCodeField,6,cityField,4,phoneNumberField,11;
	//mc,16d3,0,0,hiddenBottomRightButtonId;
	//fs,5,0,0,;
	self.ncip()
	self.start(inputs)

	next := self.lookup["hiddenBottomRightButtonId"]
	continuebtn := elements["hiddenBottomRightButtonId"]

	logger.Debug("next btn %+v", next)

	self.clickMouse(
		random.IntRange(next.Left, next.Right),
		random.IntRange(next.Top, next.Bottom),
		String(continuebtn.Attr2("id")),
		&Options{
			TimeIncrement: randomTime(0x30, 0x800),
			DontMoveMouse: true,
		},
	)

	self.submit(randomTime(3, 7))

	return self.recordsToString(self.Records)
}

func (self *Generator) Elapsed() time.Duration {
	return self.currentTime.Sub(self.initialTime)
}

func (self *Generator) selectListBox(e *ElementPosition) {
	self.moveMouse(
		random.IntRange(e.Left, e.Right),
		random.IntRange(e.Top, e.Bottom),
		e.Name,
		nil,
	)
	self.clickMouse(
		random.IntRange(e.Left, e.Right),
		random.IntRange(e.Top, e.Bottom),
		e.Name,
		nil,
	)
	self.moveMouse(
		random.IntRange(e.Left, e.Right),
		random.IntRange(e.Top+10, e.Bottom+100),
		e.Name,
		nil,
	)
}

func (self *Generator) ncip() {
	self.addRecordWorder(
		&Record{
			Type:      Ncip,
			Timestamp: self.currentTime.Unix(),
		},
		nil,
	)
}

func (self *Generator) start(inputs []*goquery.Selection) {
	self.AddRecord(
		&Record{
			Type:   Start,
			Inputs: inputs,
		},
		&Options{
			DontMoveMouse: true,
		},
	)
}

func (self *Generator) focusAndValueLength(element *ElementPosition) {
	self.AddRecord(
		&Record{
			Type:           FocusAndValueLength,
			ElementFocused: element.Name,
		},
		&Options{
			DontMoveMouse: true,
			TimeIncrement: randomTime(0, 2),
		},
	)
}

func (self *Generator) focusOnly(element *ElementPosition, increment time.Duration, dontMoveMouse bool) {
	self.AddRecord(
		&Record{
			Type:           FocusOnly,
			ElementFocused: element.Name,
		},
		&Options{
			DontMoveMouse: dontMoveMouse,
			TimeIncrement: increment,
		},
	)
}

func (self *Generator) focusBlur(element *ElementPosition, increment time.Duration, dontMoveMouse bool) {
	self.AddRecord(
		&Record{
			Type:           FocusBlur,
			ElementFocused: element.Name,
		},
		&Options{
			DontMoveMouse: dontMoveMouse,
			TimeIncrement: increment,
		},
	)
}

func (self *Generator) moveMouse(x, y int, element String, options *Options) {
	self.AddRecord(
		&Record{
			Type:           MouseMove,
			X:              x,
			Y:              y,
			ElementHovered: element,
		},
		options,
	)
}

func (self *Generator) clickMouse(x, y int, element String, options *Options) {
	self.AddRecord(
		&Record{
			Type:           MouseClick,
			X:              x,
			Y:              y,
			ElementHovered: element,
		},
		options,
	)
}

func (self *Generator) keyDown(increment time.Duration) {
	self.AddRecord(&Record{Type: KeyDown}, &Options{TimeIncrement: increment})
}

func (self *Generator) submit(increment time.Duration) {
	self.AddRecord(
		&Record{Type: Submit},
		&Options{
			TimeIncrement: increment,
			DontMoveMouse: true,
		},
	)
}

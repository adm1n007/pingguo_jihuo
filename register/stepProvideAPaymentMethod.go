package register

import (
	. "fmt"
	. "ml/array"
	. "ml/dict"
	. "ml/strings"
	. "ml/trace"

	"encoding/base64"
	"encoding/json"
	urllib "net/url"
//	"strconv"

	"ituneslib"
	httplib "ml/net/http"
	"ml/random"

	"github.com/PuerkitoBio/goquery"

	"globals"
//	"htmlutils"
	"utility"
	//"inputRecord"

	"ml/logging/logger"
)

func (self *AppleIdRegister) getProvidePaymentInfo(doc *goquery.Document, resp *httplib.Response, elementIds Array, callback fillFormCallback) (Dict, elementDict) {
	data := getAllElement(doc.Selection, elementIds)
	if len(data) != len(elementIds) {
		logger.Debug("getProvidePaymentInfo: doc: %v, doc.Selection: %v", doc, doc.Selection)
		logger.Debug("getProvidePaymentInfo: data: %#v, elementIds: %+v", data, elementIds)
		Raise(NewGenericError("len(data) != len(ids), resp incorrect"))
	}

	data["continue"] = doc.Find("input.continue[type=submit]")

	script := doc.Find("script[type*=text][type*=javascript]:contains('function ndpd_load')")

	formData := JsonDict{}
	widgetKeyData := JsonDict{}

	RaiseIf(json.Unmarshal(String(FORM_DATA_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &formData))
	RaiseIf(json.Unmarshal(String(WIDGET_KEY_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &widgetKeyData))

	widgetKey := "p"
	widgetKeyRandom := self.getWidgetKeyRandom()
	// widgetKey, widgetKeyRandom = self.getWidgetKey(widgetKeyData, resp.Request.URL)

	form := Dict{
		data[`country`].Attr2(`name`): self.account.Country.ShortName(),
		`sp`:  ``,
		`res`: ``,
		data[`codeRedemptionField`].Attr2(`name`): getInputValue(data[`codeRedemptionField`]),

		data[`ndpd-s`].Attr2(`name`):   formData[`s`],
		data[`ndpd-f`].Attr2(`name`):   formData[`f`],
		data[`ndpd-fm`].Attr2(`name`):  formData[`fm`],
		data[`ndpd-w`].Attr2(`name`):   formData[`w`],
		data[`ndpd-ipr`].Attr2(`name`): "",
		data[`ndpd-di`].Attr2(`name`):  self.deviceInfo,
		data[`ndpd-bi`].Attr2(`name`):  self.browserInfo,
		data[`ndpd-wk`].Attr2(`name`):  widgetKey,
		data[`ndpd-vk`].Attr2(`name`):  data[`ndpd-vk`].Attr2("value"),
		`ndpd-wkr`:                     widgetKeyRandom,

		data[`hidden-captcha-player-mode`].Attr2(`name`): data[`hidden-captcha-player-mode`].Attr2("value"),
		data[`continue`].Attr2(`name`):                   data[`continue`].Attr2("value"),

		data[`machineGUID`].Attr2(`name`): self.machineGUID,
		data[`pageUUID`].Attr2(`name`):    data[`pageUUID`].Attr2(`value`),
		data[`signature`].Attr2(`name`):   base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data[`pageUUID`].Attr2(`value`)) + self.machineGUID).Encode(CP_UTF8))),

		data[`longName`].Attr2(`name`): Sprintf(`%s,%s`, self.lastName, self.firstName),
	}

	if callback != nil {
		callback(form, data)
	}

	return form, data
}

func (self *AppleIdRegister) stepProvideAPaymentMethod(refer *urllib.URL, next string, form Dict) bool {
	self.info("stepProvideAPaymentMethod")

	// form missed UPCC
	resp := self.post(
		buildNext(refer, next),
		Dict{
			"headers": Dict{
				"Referer":         refer.String(),
				"Origin":          buildOrigin(refer),
				"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"Accept-Encoding": "gzip, deflate",
				"Content-Type":    "application/x-www-form-urlencoded",
				`Connection`:      `keep-alive`,
			},
			"body": form,
		},
	)

	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	switch self.account.Country {
	case ituneslib.CountryID_China:
		form = self.stepProvideAPaymentMethod_China(doc, resp)

	case ituneslib.CountryID_India:
		form = self.stepProvideAPaymentMethod_India(doc, resp)

	case ituneslib.CountryID_NewZealand:
		form = self.stepProvideAPaymentMethod_NewZealand(doc, resp)

	case ituneslib.CountryID_Vietnam:
		form = self.stepProvideAPaymentMethod_Vietnam(doc, resp)

	case ituneslib.CountryID_Taiwan:
		form = self.stepProvideAPaymentMethod_Taiwan(doc, resp)

	default:
		Raisef("NotImplemented: %v", self.account.Country)
	}

	return self.stepVerifyApppleId(refer, doc.Find("form[method=post]").Attr2("action"), form)
}

func (self *AppleIdRegister) stepProvideAPaymentMethod_China(doc *goquery.Document, resp *httplib.Response) Dict {
	/*
	self.info("stepProvideAPaymentMethod_China")

	elements := Array{
		`country`,
		"mobile-phone",
		`codeRedemptionField`,

		`lastFirstName`,
		`firstName`,
		`street1`,
		`street2`,
		`street3`,
		`city`,
		`postalcode`,
		`state`,
		`phone1Number`,

		`ndpd-s`,
		`ndpd-f`,
		`ndpd-fm`,
		`ndpd-w`,
		`ndpd-ipr`,
		`ndpd-di`,
		`ndpd-bi`,
		`ndpd-wk`,
		`ndpd-vk`,

		`hidden-captcha-player-mode`,
		"card_type_id",

		`pageUUID`,
		`machineGUID`,
		`signature`,

		`longName`,
	}

	form, _ := self.getProvidePaymentInfo(doc, resp, elements, func(form Dict, elements elementDict) {
		state := selectOption(elements["state"])
		stateIndex, _ := strconv.Atoi(state)
		stateName := elements["state"].Find("option").Slice(stateIndex+1, stateIndex+2).Text()
		cityName, postalcode := utility.GeneratePostalCode(stateName)

		e := Dict{
			elements["mobile-phone"].Attr2("name"):  getInputValue(elements["mobile-phone"]),
			`credit-card-type`:                      ``,
			elements["card_type_id"].Attr2("name"):  elements["card_type_id"].Attr2("value"),
			elements[`lastFirstName`].Attr2(`name`): self.lastName,
			elements[`firstName`].Attr2(`name`):     self.firstName,
			elements[`street1`].Attr2(`name`):       utility.GenerateStreet1(),
			elements[`street2`].Attr2(`name`):       "", // utility.GenerateStreet2(),
			elements[`street3`].Attr2(`name`):       "", // utility.GenerateStreet3(),
			elements[`city`].Attr2(`name`):          cityName,
			elements[`postalcode`].Attr2(`name`):    postalcode,
			elements[`state`].Attr2(`name`):         state,
			elements[`phone1Number`].Attr2(`name`):  utility.GenerateMobileNumber(),
		}

		mergeDict(form, e)

		form[elements[`street2`].Attr2(`name`)] = utility.GenerateStreet2()
		form[elements[`street3`].Attr2(`name`)] = utility.GenerateStreet3()

		inputs := []*goquery.Selection{}
		htmlutils.FindInputs(doc.Selection).Each(func(i int, s *goquery.Selection) {
			inputs = append(inputs, s)
		})

		self.info("Generate")
		generator := inputRecord.New(inputRecord.LayoutForCountry(self.account.Country))
		form[elements["ndpd-ipr"].Attr2("name")] = generator.Generate(elements, inputs, form)

		// delete(form, elements["ndpd-s"].Attr2("name"))
		// delete(form, elements["ndpd-f"].Attr2("name"))
		// delete(form, elements["ndpd-fm"].Attr2("name"))
		// delete(form, elements["ndpd-w"].Attr2("name"))
		// delete(form, elements["ndpd-ipr"].Attr2("name"))
		// delete(form, elements["ndpd-di"].Attr2("name"))
		// delete(form, elements["ndpd-bi"].Attr2("name"))
		// delete(form, elements["ndpd-wk"].Attr2("name"))
		// delete(form, elements["ndpd-vk"].Attr2("name"))
		// delete(form, elements["ndpd-ipr"].Attr2("name"))
		// delete(form, "ndpd-wkr")
	})
*/
	return nil
}

func (self *AppleIdRegister) stepProvideAPaymentMethod_India(doc *goquery.Document, resp *httplib.Response) Dict {
	/*
	self.info("stepProvideAPaymentMethod_India")

	elements := Array{
		`country`,
		`codeRedemptionField`,

		"salutation",
		"firstName",
		"lastName",
		"street1",
		"street2",
		"street3",
		"city",
		"postalcode",
		"state",
		"phone1AreaCode",
		"phone1Number",

		`ndpd-s`,
		`ndpd-f`,
		`ndpd-fm`,
		`ndpd-w`,
		`ndpd-ipr`,
		`ndpd-di`,
		`ndpd-bi`,
		`ndpd-wk`,
		`ndpd-vk`,

		`hidden-captcha-player-mode`,

		`pageUUID`,
		`machineGUID`,
		`signature`,

		`longName`,
	}

	form, _ := self.getProvidePaymentInfo(doc, resp, elements, func(form Dict, elements elementDict) {
		state := selectOption(elements["state"])

		e := Dict{
			elements[`salutation`].Attr2(`name`):     selectOption(elements[`salutation`]),
			elements[`firstName`].Attr2(`name`):      self.firstName,
			elements[`lastName`].Attr2(`name`):       self.lastName,
			elements[`street1`].Attr2(`name`):        utility.GenerateStreet1(), // not valid for india.
			elements[`street2`].Attr2(`name`):        "",
			elements[`street3`].Attr2(`name`):        "",
			elements[`city`].Attr2(`name`):           utility.GenerateRandomString()[:5],
			elements[`state`].Attr2(`name`):          state,
			elements[`postalcode`].Attr2(`name`):     Sprintf("%06d", random.IntRange(10000, 999999)),
			elements[`phone1AreaCode`].Attr2(`name`): random.IntRange(10, 99),
			elements[`phone1Number`].Attr2(`name`):   utility.GenerateMobileNumber(),
		}

		mergeDict(form, e)
		Raise(NewNotImplementedError("no inputRecord impl for %v", self.account.Country))
		// form[elements["ndpd-ipr"].Attr2("name")] = generateInputRecordsIndia(doc.Selection, elements, form)
	})
*/
	return nil
}

func (self *AppleIdRegister) stepProvideAPaymentMethod_NewZealand(doc *goquery.Document, resp *httplib.Response) Dict {
	self.info("stepProvideAPaymentMethod_NewZealand")

	elements := Array{
		`country`,
		`codeRedemptionField`,

		"salutation",
		"firstName",
		"lastName",
		"street1",
		"street2",
		"suburb",
		"postalcode",
		"city",
		"phone1AreaCode",
		"phone1Number",

		`ndpd-s`,
		`ndpd-f`,
		`ndpd-fm`,
		`ndpd-w`,
		`ndpd-ipr`,
		`ndpd-di`,
		`ndpd-bi`,
		`ndpd-wk`,
		`ndpd-vk`,

		`hidden-captcha-player-mode`,

		`pageUUID`,
		`machineGUID`,
		`signature`,

		`longName`,
	}

	form, _ := self.getProvidePaymentInfo(doc, resp, elements, func(form Dict, elements elementDict) {

		e := Dict{
			elements[`salutation`].Attr2(`name`):     selectOption(elements[`salutation`]),
			elements[`firstName`].Attr2(`name`):      self.firstName,
			elements[`lastName`].Attr2(`name`):       self.lastName,
			elements[`street1`].Attr2(`name`):        utility.GenerateRandomString()[:8],
			elements[`street2`].Attr2(`name`):        utility.GenerateRandomString()[:8],
			elements[`suburb`].Attr2(`name`):         utility.GenerateRandomString()[:6],
			elements[`postalcode`].Attr2(`name`):     Sprintf("%04d", random.IntRange(100, 9999)),
			elements[`city`].Attr2(`name`):           utility.GenerateRandomString()[:5],
			elements[`phone1AreaCode`].Attr2(`name`): random.IntRange(10, 99),
			elements[`phone1Number`].Attr2(`name`):   utility.GenerateMobileNumber(),
		}

		mergeDict(form, e)
		Raise(NewNotImplementedError("no inputRecord impl for %v", self.account.Country))
		// form[elements["ndpd-ipr"].Attr2("name")] = generateInputRecordsNewZealand(doc.Selection, elements, form)
	})

	return form
}

func (self *AppleIdRegister) stepProvideAPaymentMethod_Vietnam(doc *goquery.Document, resp *httplib.Response) Dict {
	self.info("stepProvideAPaymentMethod_Vietnam")

	elements := Array{
		`country`,
		`codeRedemptionField`,

		"salutation",
		"firstName",
		"lastName",
		"street1",
		"street2",
		"city",
		"postalcode",
		"phone1AreaCode",
		"phone1Number",

		`ndpd-s`,
		`ndpd-f`,
		`ndpd-fm`,
		`ndpd-w`,
		`ndpd-ipr`,
		`ndpd-di`,
		`ndpd-bi`,
		`ndpd-wk`,
		`ndpd-vk`,

		`hidden-captcha-player-mode`,

		`pageUUID`,
		`machineGUID`,
		`signature`,

		`longName`,
	}

	form, _ := self.getProvidePaymentInfo(doc, resp, elements, func(form Dict, elements elementDict) {
		e := Dict{
			elements[`salutation`].Attr2(`name`):     selectOption(elements[`salutation`]),
			elements[`firstName`].Attr2(`name`):      self.firstName,
			elements[`lastName`].Attr2(`name`):       self.lastName,
			elements[`street1`].Attr2(`name`):        utility.GenerateRandomString()[:8],
			elements[`street2`].Attr2(`name`):        utility.GenerateRandomString()[:8],
			elements[`city`].Attr2(`name`):           utility.GenerateRandomString()[:5],
			elements[`postalcode`].Attr2(`name`):     Sprintf("%06d", random.IntRange(10000, 999999)),
			elements[`phone1AreaCode`].Attr2(`name`): random.IntRange(10, 99),
			elements[`phone1Number`].Attr2(`name`):   utility.GenerateMobileNumber(),
		}

		mergeDict(form, e)
		Raise(NewNotImplementedError("no inputRecord impl for %v", self.account.Country))
		// form[elements["ndpd-ipr"].Attr2("name")] = generateInputRecordsVietnam(doc.Selection, elements, form)
	})

	return form
}

func (self *AppleIdRegister) stepProvideAPaymentMethod_Taiwan(doc *goquery.Document, resp *httplib.Response) Dict {
	self.info("stepProvideAPaymentMethod_Taiwan")

	elements := Array{
		`country`,
		`codeRedemptionField`,

		"salutation",
		"lastFirstName",
		"firstName",
		"citypopup",
		"street1",
		"street2",
		"postalcode",
		"phone1AreaCode",
		"phone1Number",

		`ndpd-s`,
		`ndpd-f`,
		`ndpd-fm`,
		`ndpd-w`,
		`ndpd-ipr`,
		`ndpd-di`,
		`ndpd-bi`,
		`ndpd-wk`,
		`ndpd-vk`,

		`hidden-captcha-player-mode`,

		`pageUUID`,
		`machineGUID`,
		`signature`,

		`longName`,
	}

	form, _ := self.getProvidePaymentInfo(doc, resp, elements, func(form Dict, elements elementDict) {
		e := Dict{
			`credit-card-type`:                       ``,
			elements[`salutation`].Attr2(`name`):     selectOption(elements[`salutation`]),
			elements[`lastFirstName`].Attr2(`name`):  self.firstName,
			elements[`firstName`].Attr2(`name`):      self.lastName,
			elements[`citypopup`].Attr2(`name`):      selectOption(elements["citypopup"]),
			elements[`street1`].Attr2(`name`):        utility.GenerateStreet1(),
			elements[`street2`].Attr2(`name`):        utility.GenerateStreet2(),
			elements[`postalcode`].Attr2(`name`):     Sprintf("%05d", random.IntRange(10000, 100000)),
			elements[`phone1AreaCode`].Attr2(`name`): random.IntRange(886, 1000),
			elements[`phone1Number`].Attr2(`name`):   utility.GenerateMobileNumber(),
		}

		mergeDict(form, e)
		Raise(NewNotImplementedError("no inputRecord impl for %v", self.account.Country))
		// form[elements["ndpd-ipr"].Attr2("name")] = generateInputRecordsTaiwan(doc.Selection, elements, form)

		form[elements[`street1`].Attr2(`name`)] = utility.GeneratePinyin()
		form[elements[`street2`].Attr2(`name`)] = utility.GeneratePinyin()
		form[elements["ndpd-ipr"].Attr2("name")] = Sprintf("ncip,0,%x,1,1;st,0,codeRedemptionField,0,lastFirstName,0,firstName,0,street1,0,street2,0,postalcode,0,phone1AreaCode,0,phone1Number,0;mm,2d,223,194;mc,1a78,8a,1d8,salutation;mm,0,9d,203,lastFirstName;kk,223,0,lastFirstName;ff,0,lastFirstName;mc,5f,a1,1fb,lastFirstName;kd,6fc;kd,19f;kd,16;kd,64;fb,190,lastFirstName;kk,1,0,firstName;ff,0,firstName;mc,a,142,1fe,firstName;kd,1b;kd,64;kd,63;kd,26;kd,71;kd,7d;kd,3e;kd,d;mm,90,8c,227,citypopup;fb,c2,firstName;mc,436,83,226,citypopup;kk,138,0,street1;ff,0,street1;mc,47,a4,24a,street1;kd,71;kd,3e;kd,19;kd,8a;kd,19;kd,7d;kd,57;kd,7d;fb,10d,street1;kk,0,0,street2;ff,0,street2;mc,4e,dd,25d,street2;kd,10a;kd,19;kd,8d;kd,ac;kd,4e;kd,48;kd,25;kd,64;kd,af;kd,d;kd,4b;kd,bb;kd,c;kd,19;ts,0,3a98;kd,4f;fb,14e,street2;kk,1,0,postalcode;ff,0,postalcode;mc,31,b0,282,postalcode;kd,4fb;mm,74,b0,281,postalcode;kd,48;kd,21c;kd,14e;kd,10a;fb,365,postalcode;kk,1,0,phone1AreaCode;ff,0,phone1AreaCode;mc,46,8c,2a8,phone1AreaCode;kd,3d6;kd,64;kd,109;fb,14e,phone1AreaCode;kk,1,0,phone1Number;ff,0,phone1Number;mc,46,d5,2a2,phone1Number;mm,a1,d7,2a2,phone1Number;kd,368;kd,3e;kd,3f;kd,d1;kd,c8;kd,b2;kd,93;mc,52c,d3,2a9,phone1Number;mc,a0,d3,2a9,phone1Number;kd,133;kd,70;kd,14f;kd,12c;kd,bc;kd,fd;kd,89;kd,104;kd,190;kd,ae;kd,15e;mm,c0,dd,2a9,phone1Number;fb,459,phone1Number;mc,40,3b3,330,2.0.1.1.3.0.7.11.3.9.1;fs,3,0,0,;", globals.GetCurrentTime()/1000)
	})

	return form
}

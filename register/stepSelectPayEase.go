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

	"ituneslib"
	"ml/random"

	"github.com/PuerkitoBio/goquery"

	"globals"
	"utility"

	"ml/logging/logger"
)

func (self *AppleIdRegister) getSelectPaymentInfo(doc *goquery.Document, elementIds Array, callback fillFormCallback) (Dict, elementDict) {

	data := getAllElement(doc.Selection, elementIds)
	if len(data) != len(elementIds) {
		logger.Debug("getSelectPaymentInfo: data: %+v, elementIds: %+v", data, elementIds)
		Raise(NewGenericError("len(data) != len(elementIds), resp incorrect"))
	}

	data["continue"] = doc.Find("input.continue[type=submit]")

	script := doc.Find("script[type*=text][type*=javascript]:contains('function ndpd_load')")

	formData := JsonDict{}
	widgetKeyData := JsonDict{}

	RaiseIf(json.Unmarshal(String(FORM_DATA_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &formData))
	RaiseIf(json.Unmarshal(String(WIDGET_KEY_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &widgetKeyData))

	//logger.Debug("formData: %+v, widgetKeyData: %+v", formData, widgetKeyData)

	widgetKey := "p"
	widgetKeyRandom := self.getWidgetKeyRandom()
	// widgetKey, widgetKeyRandom = self.getWidgetKey(widgetKeyData, resp.Request.URL)

	inputs := []string{}
	getAllInput(doc.Selection).Each(func(i int, s *goquery.Selection) {
		inputs = append(inputs, Sprintf("%s,%d", getInputName(s), String(s.AttrOr("value", "")).Length()))
	})

	// prefer none instead payease
	payment := doc.Find("input[type=submit][name=None]") // 无
	if payment == nil || payment.Length() == 0 {
		payment = doc.Find("input[type=submit][name=PayEase]")
		self.debug("select PayEase") // 银行卡
	}

	inputRecordPattern := []string{
		Sprintf(`ncip,0,%x,1,1`, globals.GetCurrentTime()/1000),
		string(`st,0,` + String(",").Join(inputs)),
		Sprintf(`mm,%x,%x,%x`, random.IntRange(0x10, 0x20), random.IntRange(0x100, 0x200), random.IntRange(0x100, 0x200)),
		Sprintf(`mc,%x,%x,%x`, random.IntRange(0x250, 1000), random.IntRange(0x97, 0xCB), random.IntRange(0x97, 0xCB)),
		Sprintf(`mc,%x,%x,%x`, random.IntRange(0x250, 1000), random.IntRange(0x97, 0xCB), random.IntRange(0x97, 0xCB)),
		Sprintf(`mc,%x,%x,%x,%s`, random.IntRange(0x250, 1000), random.IntRange(210, 340), random.IntRange(275, 300), getInputName(payment)),
		Sprintf(`fs,%x,0,0,;`, random.IntRange(0, 0x10)),
	}

	// some key became "None" here! 'None': '无',
	// maybe apple check the key order?
	form := Dict{
		data[`country`].Attr2(`name`): self.account.Country.ShortName(), // CN
		`credit-card-type`:            ``,
		`sp`:                          ``,
		`res`:                         ``,
		payment.Attr2("name"): getInputValue(payment), // "None": "无"
		//`UPCC`:                        `UnionPay`,
		data[`email-regex`].Attr2(`name`):         data[`email-regex`].Attr2(`value`),         // `^[A-Za-z0-9._%+-]+@`,
		data[`phone-regex`].Attr2(`name`):         data[`phone-regex`].Attr2(`value`),         // `^[0-9-]{10,30}$`,
		data[`codeRedemptionField`].Attr2(`name`): getInputValue(data[`codeRedemptionField`]), // 2.0.1.1.3.0.7.11.3.1.0.5.21.1.3.1.3.2.3, 兑换码

		data[`sesame-id`].Attr2(`name`):    ``,
		data[`national-id`].Attr2(`name`):  ``,
		data[`country-code`].Attr2(`name`): `0`,

		data[`ndpd-s`].Attr2(`name`):   formData[`s`],
		data[`ndpd-f`].Attr2(`name`):   formData[`f`],
		data[`ndpd-fm`].Attr2(`name`):  formData[`fm`],
		data[`ndpd-w`].Attr2(`name`):   formData[`w`],
		data[`ndpd-ipr`].Attr2(`name`): String(`;`).Join(inputRecordPattern),
		data[`ndpd-di`].Attr2(`name`):  self.deviceInfo,
		data[`ndpd-bi`].Attr2(`name`):  self.browserInfo,

		data[`ndpd-probe`].Attr2(`name`): ``,            //data[`ndpd-probe`].Attr2("value"),
		data[`ndpd-af`].Attr2(`name`):    ``,            //data[`ndpd-af`].Attr2("value"),
		data[`ndpd-fv`].Attr2(`name`):    `fv,mp4`,      //data[`ndpd-fv`].Attr2("value"), // fv,mp4
		data[`ndpd-fa`].Attr2(`name`):    `fa,mpeg,wav`, //data[`ndpd-fa`].Attr2("value"), // fa,mpeg,wav
		// TODO: separate browserInfo with plugin for Mac/PC
		data[`ndpd-bp`].Attr2(`name`): `p,Java Applet Plug-in,Quartz Composer Plug-In,QuickTime Plug-in+7`, // p,Java+Applet+Plug-in,Quartz+Composer+Plug-In,QuickTime+Plug-in+7

		data[`ndpd-wk`].Attr2(`name`): widgetKey,
		data[`ndpd-vk`].Attr2(`name`): data[`ndpd-vk`].Attr2("value"),
		`ndpd-wkr`:                    widgetKeyRandom,

		data[`hidden-captcha-player-mode`].Attr2(`name`): data[`hidden-captcha-player-mode`].Attr2("value"), // VIDEO

		data[`machineGUID`].Attr2(`name`): self.machineGUID,
		data[`pageUUID`].Attr2(`name`):    data[`pageUUID`].Attr2(`value`),
		data[`signature`].Attr2(`name`):   base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data[`pageUUID`].Attr2(`value`)) + self.machineGUID).Encode(CP_UTF8))),
	}

	//logger.Debug("data[`codeRedemptionField`]: %+v, name: %#v", data[`codeRedemptionField`], data[`codeRedemptionField`].Attr2(`name`))
	//logger.Debug("paymentname %#v: %#v", payment.Attr2("name"), getInputValue(payment))
	//logger.Debug("getSelectPaymentInfo: data: %+v,  form: %+v", data, form)

	callback(form, data)

	return form, data
}

func (self *AppleIdRegister) stepSelectPayEase(refer *urllib.URL, next string, form Dict) bool {
	self.info("stepSelectPayEase")

	switch self.account.Country {
	case ituneslib.CountryID_China:
		break

	case ituneslib.CountryID_India,
		ituneslib.CountryID_Taiwan,
		ituneslib.CountryID_NewZealand,
		ituneslib.CountryID_Vietnam:
		return self.stepProvideAPaymentMethod(refer, next, form)
	}

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

	self.pingUrl(doc, resp, true)

	switch self.account.Country {
	case ituneslib.CountryID_China:
		form = self.stepSelectPayEase_China(doc)

	default:
		Raise(NewNotImplementedError("NotImplemented: %v", self.account.Country))
	}

	return self.stepProvideAPaymentMethod(
		resp.Request.URL,
		doc.Find("form[method=post]").Attr2("action"),
		form,
	)
}

func (self *AppleIdRegister) stepSelectPayEase_China(doc *goquery.Document) Dict {
	self.info("stepSelectPayEase_China")

	elements := Array{
		"country",
		`codeRedemptionField`,
		//"cc_number",
		//"cc_month",
		//"cc_year",
		//"cc_ccv",
		"mobile-phone",
		"card_type_id",

		"lastFirstName",
		"firstName",
		"street1",
		"street2",
		"street3",
		"city",
		"postalcode",
		"state",
		"phone1Number",

		"phone-regex",
		"email-regex",

		"ndpd-s",
		"ndpd-f",
		"ndpd-fm",
		"ndpd-w",
		"ndpd-ipr",
		"ndpd-di",
		"ndpd-bi",
		"ndpd-wk",
		"ndpd-vk",

		"ndpd-probe",
		"ndpd-af",
		"ndpd-fv",
		"ndpd-fa",
		"ndpd-bp",

		"sesame-id",
		"national-id",
		"country-code",

		"hidden-captcha-player-mode",

		"pageUUID",
		"machineGUID",
		"signature",

		"longName",
	}

	form, _ := self.getSelectPaymentInfo(doc, elements, func(form Dict, elements elementDict) {
		// anonymouse func?
		e := Dict{
			//elements[`cc_number`].Attr2(`name`):    getInputValue(elements[`cc_number`]),
			//elements[`cc_month`].Attr2(`name`):     `0`,
			//elements[`cc_year`].Attr2(`name`):      `0`,
			//elements[`cc_ccv`].Attr2(`name`):       getInputValue(elements[`cc_ccv`]),
			elements[`mobile-phone`].Attr2(`name`): getInputValue(elements[`mobile-phone`]),
			elements[`card_type_id`].Attr2(`name`): getInputValue(elements[`card_type_id`]),

			elements[`lastFirstName`].Attr2(`name`): getInputValue(elements[`lastFirstName`]),
			elements[`firstName`].Attr2(`name`):     getInputValue(elements[`firstName`]),
			elements[`street1`].Attr2(`name`):       getInputValue(elements[`street1`]),
			elements[`street2`].Attr2(`name`):       getInputValue(elements[`street2`]),
			elements[`street3`].Attr2(`name`):       getInputValue(elements[`street3`]),
			elements[`city`].Attr2(`name`):          getInputValue(elements[`city`]),
			elements[`postalcode`].Attr2(`name`):    getInputValue(elements[`postalcode`]),
			elements[`state`].Attr2(`name`):         elements[`state`].Find("option").First().Attr2(`value`),
			elements[`phone1Number`].Attr2(`name`):  getInputValue(elements[`phone1Number`]),
			elements[`longName`].Attr2(`name`):      Sprintf(`%s,%s`, getInputValue(elements[`lastFirstName`]), getInputValue(elements[`firstName`])),
		}

		for k, v := range e {
			form[k] = v
		}

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

	return form
}

func (self *AppleIdRegister) stepSelectPayEase_India(doc *goquery.Document) Dict {
	self.info("stepSelectPayEase_India")

	elements := Array{
		"country",
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
		"phone1Number",

		"ndpd-s",
		"ndpd-f",
		"ndpd-fm",
		"ndpd-w",
		"ndpd-ipr",
		"ndpd-di",
		"ndpd-bi",
		"ndpd-wk",
		"ndpd-vk",

		"hidden-captcha-player-mode",

		"pageUUID",
		"machineGUID",
		"signature",

		"longName",
	}

	form, _ := self.getSelectPaymentInfo(doc, elements, func(form Dict, elements elementDict) {
		e := Dict{
			elements[`salutation`].Attr2(`name`):   elements[`salutation`].Find("option").First().Attr2(`value`),
			elements[`firstName`].Attr2(`name`):    getInputValue(elements[`firstName`]),
			elements[`lastName`].Attr2(`name`):     getInputValue(elements[`lastName`]),
			elements[`street1`].Attr2(`name`):      getInputValue(elements[`street1`]),
			elements[`street2`].Attr2(`name`):      getInputValue(elements[`street2`]),
			elements[`street3`].Attr2(`name`):      getInputValue(elements[`street3`]),
			elements[`city`].Attr2(`name`):         getInputValue(elements[`city`]),
			elements[`postalcode`].Attr2(`name`):   getInputValue(elements[`postalcode`]),
			elements[`state`].Attr2(`name`):        elements[`state`].Find("option").First().Attr2(`value`),
			elements[`phone1Number`].Attr2(`name`): getInputValue(elements[`phone1Number`]),
			elements[`longName`].Attr2(`name`):     Sprintf(`%s,%s`, getInputValue(elements[`lastName`]), getInputValue(elements[`firstName`])),
		}

		for k, v := range e {
			form[k] = v
		}
	})

	return form
}

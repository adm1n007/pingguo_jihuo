package register

import (
	. "fmt"
	. "ml/array"
	. "ml/dict"
	. "ml/strings"
	. "ml/trace"

	"encoding/base64"
	"encoding/json"
	"math"
	"math/rand"
	urllib "net/url"
	"plistlib"
	"regexp"
	"strconv"
	"time"

	"ml/logging/logger"
	"ml/net/http"
	"ml/random"

	"github.com/PuerkitoBio/goquery"

	"ituneslib"

	"account"
	"globals"
	"proxy"
	"utility"
	//"./inputRecord"
)

// var (
//     FORM_DATA_PATTERN       = regexp.MustCompile(`"fd"\s*:\s*({.*?})`)
//     WIDGET_KEY_PATTERN      = regexp.MustCompile(`"wk"\s*:\s*({.*?})`)
//     STRIP_CALLBACK_PATTERN  = regexp.MustCompile(`.*\(({.*?})\).*`)
// )

type AppleIdRegisteriOS struct {
	session     *http.Session
	bag         Dict
	sapSession  *ituneslib.SapSession
	account     *account.AppleAccount
	machineGUID String
	buyVID      String
	machineName String
	deviceInfo  String
	browserInfo String
	width       int
	height      int
	colorDepth  int
	innerWidth  int
	innerHeight int
	lang        String
	firstName   String
	lastName    String
	retry0      int64
	proxy       *proxy.Proxy
}

func NewRegisteriOS(account *account.AppleAccount, proxy *proxy.Proxy) *AppleIdRegisteriOS {
	session := http.NewSession()

	RaiseIf(session.SetProxy(proxy.Host, proxy.Port, proxy.User, proxy.Password))
	// session.SetProxy("localhost", 6789)

	session.SetHeaders(Dict{
		"Accept":                  "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding":         "gzip, deflate",
		"Accept-Language":         "zh-cn",
		"User-Agent":              "AppStore/2.0 iOS/7.1.2 model/iPhone3,3 build/11D257 (5; dt:70)",
		"X-Apple-Client-Versions": "GameCenter/2.0",
		"X-Apple-Connection-Type": "WiFi",
		"X-Apple-Store-Front":     "143465-19,24 t:native",
		"Connection":              "keep-alive",
	})

	register := &AppleIdRegisteriOS{
		session:    session,
		bag:        Dict{},
		sapSession: ituneslib.NewSapSession(),
		account:    account,
	}

	register.proxy = proxy
	register.init()

	return register
}

func (self *AppleIdRegisteriOS) init() {
	if len(self.account.ApplePassword) == 0 {
		self.account.ApplePassword = string(utility.GeneratePassword())
	}

	width, height, colorDepth, innerWidth, innerHeight, lang := utility.GenerateiOSBrowserInfo()

	self.machineName = utility.GeneratePinyin()
	self.machineGUID = utility.GenerateiOSMachineGuid()
	self.machineGUID = utility.GenerateiOSMachineGuid()
	self.machineGUID = "0fb84ae483ca4a19e9e84547c36473ec29d95c2e"
	self.buyVID = "C0214840-0331-4DEF-96F6-61AB2707E4B4" //utility.GenerateUUID() // C0214840-0331-4DEF-96F6-61AB2707E4B4
	self.deviceInfo = utility.GenerateDeviceInfo()
	self.browserInfo = String(Sprintf("b1.%dx%d %dx%d %v %v.-480.%v", width, height, innerWidth, innerHeight, colorDepth, colorDepth, lang))
	self.width = width
	self.height = height
	self.colorDepth = colorDepth
	self.innerWidth = innerWidth
	self.innerHeight = innerHeight
	self.lang = lang

	switch self.account.Country {
	case ituneslib.CountryID_China:
		self.firstName = utility.GenerateChineseGivenName() //utility.GenerateName(1, 2)
		self.lastName = utility.GenerateChineseFamilyName() //utility.GenerateName(1, 2)

	default:
		self.firstName = utility.GenerateRandomString()[:6]
		self.lastName = utility.GenerateRandomString()[:6]
	}
}

func (self *AppleIdRegisteriOS) Close() {
	if self.sapSession != nil {
		self.sapSession.Close()
		self.sapSession = nil
	}

	self.session.Close()
}

func (self *AppleIdRegisteriOS) info(format interface{}, args ...interface{}) {
	logger.Info("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *AppleIdRegisteriOS) debug(format interface{}, args ...interface{}) {
	logger.Debug("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *AppleIdRegisteriOS) request(method, url interface{}, params ...Dict) (resp *http.Response) {
	self.debug("request: %v, params: %v", url, params)

	for {
		exp := Try(func() { resp = self.session.Request(method, url, params...) })

		if exp != nil {
			self.debug("%v", exp)
			e := exp.Value.(*http.HttpError)

			switch e.Type {
			case http.HTTP_ERROR_INVALID_RESPONSE,
				http.HTTP_ERROR_BAD_GATE_WAY,
				http.HTTP_ERROR_TIMEOUT,
				http.HTTP_ERROR_GENERIC:
				time.Sleep(time.Second)
				continue

			case http.HTTP_ERROR_CONNECT_PROXY:
				fallthrough
			default:
				Raise(exp)
			}
		}

		switch resp.StatusCode {
		case http.StatusOK:
		case http.StatusFound:
		case http.StatusNotModified:
			break

		default:
			time.Sleep(time.Second)
			continue
		}

		break
	}

	// time.Sleep(time.Second * 15)
	return
}

func (self *AppleIdRegisteriOS) get(url interface{}, params ...Dict) *http.Response {
	return self.request("GET", url, params...)
}

func (self *AppleIdRegisteriOS) post(url interface{}, params ...Dict) *http.Response {
	return self.request("POST", url, params...)
}

func (self *AppleIdRegisteriOS) checkPageError(doc *goquery.Document) {
	pageText := String(doc.Text())

	if pageText.Contains("您的会话已超时。请再试一次。") {
		Raise(NewRegisterError("您的会话已超时。请再试一次。"))
	}

	if pageText.Contains("如需帮助，请联系 iTunes 支持") {
		Raise(NewRegisterError("如需帮助，请联系 iTunes 支持"))
	}

	pageError := doc.Find(".page-error")
	if pageError.Length() == 0 || pageError.Children().Length() == 0 {
		return
	}

	genericError := pageError.Find(".generic-error")
	if genericError.Length() != 0 {
		text := String(genericError.Text())
		if text.Contains("您输入的电子邮件地址已用于") {
			Raise(NewAlreadyExistsError(text.String()))
		}

		Raise(NewGenericError(text.String()))
	}

	err := pageError.Find(".error")
	if err.Length() != 0 {
		Raise(NewRegisterError(err.Text()))
	}

	self.debug("%v", pageError.Text())
	self.debug("unknown error")
	utility.SetTitle("unknown error")
	level := logger.Level()
	logger.SetLevel(9999)
	utility.Pause()
	logger.SetLevel(level)
}

func (self *AppleIdRegisteriOS) getWidgetKeyRandom() int {
	return int(math.Floor(rand.Float64()*1e6) + 1e3)
}

func (self *AppleIdRegisteriOS) getWidgetKey(widgetKeyData map[string]interface{}, refer *urllib.URL) (widgetKey, widgetKeyRandom int) {
	widgetKeyRandom = self.getWidgetKeyRandom()

	resp := self.get(
		widgetKeyData["r"].(string),
		Dict{
			//"Referer": refer.String(),
			"headers": Dict{
				"Accept":          "*/*",
				"Accept-Language": "zh-Hans",
				"Referer":         refer.String(),
			},
			"params": Dict{
				"r":  widgetKeyRandom,
				"wt": widgetKeyData["w"].(string),
			},
		},
	)

	data := map[string]interface{}{}
	RaiseIf(json.Unmarshal(String(STRIP_CALLBACK_PATTERN.FindStringSubmatch(string(resp.Text()))[1]).Encode(CP_UTF8), &data))

	widgetKey = int(data["wk"].(float64))

	return
}

func (self *AppleIdRegisteriOS) initbag() {
	var resp *http.Response

	resp = self.get("https://init.itunes.apple.com/bag.xml?ix=5&os=7&locale=en_US")

	plist := Dict{}
	resp.Plist(&plist)
	plistlib.Unmarshal(plist["bag"].([]byte), &self.bag)
}

func (self *AppleIdRegisteriOS) initsap() {
	signSapSetupCert := Dict{}
	self.get(
		self.bag["sign-sap-setup-cert"],
		Dict{
			"headers": Dict{
				"Accept":          "*/*",
				"Accept-Language": "zh-Hans",
			},
		},
	).Plist(&signSapSetupCert)

	cert := self.sapSession.ExchangeData(ituneslib.SAP_TYPE_REGISTER, signSapSetupCert["sign-sap-setup-cert"].([]byte))
	body, err := plistlib.MarshalIndent(Dict{"sign-sap-setup-buffer": cert}, plistlib.XMLFormat, "    ")
	RaiseIf(err)

	signSapSetupBuffer := Dict{}
	self.post(
		self.bag["sign-sap-setup"],
		Dict{
			"headers": Dict{
				"Accept":          "*/*",
				"Accept-Language": "zh-Hans",
				"Content-Type":    "application/x-www-form-urlencoded",
			},
			"body": body,
		},
	).Plist(&signSapSetupBuffer)

	self.sapSession.ExchangeData(ituneslib.SAP_TYPE_REGISTER, signSapSetupBuffer["sign-sap-setup-buffer"].([]byte))
	self.sapSession.Initialized = true
}

func (self *AppleIdRegisteriOS) Initialize() {
	self.initbag()
	self.initsap()
}

func (self *AppleIdRegisteriOS) pingUrl(doc *goquery.Document, resp *http.Response, switches ...bool) {
	return

	// defer func () {
	//     Println(Catch(recover()))
	//     time.Sleep(time.Hour)
	// }()

	// logger.SetLevel(99999999)

	skipItsMetricsR := false

	switch len(switches) {
	case 1:
		skipItsMetricsR = switches[0]
		fallthrough

	default:
	}

	peekProperty := func(text, property string) (value String) {
		pattern := regexp.MustCompile(Sprintf(`(?m:\.%s\s*=\s*(.*);?$)`, property))
		value = String(pattern.FindStringSubmatch(text)[1])
		value = value.Trim(";")
		if value[:1] == `"` {
			value = value[1 : value.Length()-1]
		}

		return
	}

	script := doc.Find("script[type*=text][type*=javascript]:contains('var iTSMetricsCallbackFunction')")
	protocol := doc.Find("script[type*=text][type*=x-apple-plist][id=protocol]")

	text := script.Text()

	reportingSuite := peekProperty(text, "reportingSuite")
	pageName := peekProperty(text, "pageName")
	channel := peekProperty(text, "channel")
	prop22 := peekProperty(text, "prop22")
	eVar22 := peekProperty(text, "eVar22")

	// itsMetricsR=Signup%20View%20Terms-CN@@Signup@@@@
	// itsMetricsR := "pageName@@channelName@@pageNameExtrasForBuyMetrics@@location"

	if skipItsMetricsR == false {
		referrerPageInfoKey := String("@@").Join([]string{
			pageName.Replace(" ", "%20").String(),
			channel.Replace(" ", "%20").String(),
			"",
			"",
		})

		self.debug("referrerPageInfoKey = %v @ %v", referrerPageInfoKey, resp.URL.String())

		self.session.SetCookies(
			String(resp.URL.String()),
			Dict{
				"itsMetricsR": referrerPageInfoKey,
			},
		)
	}

	plist := parsePlist(protocol.Text())

	url := plist["pings"].([]interface{})[0].(string)
	r := self.get(url)

	now := time.Unix(0, globals.GetCurrentTime()*int64(time.Millisecond)).AddDate(0, -1, 0)

	query, _ := urllib.ParseQuery(r.Request.URL.RawQuery)

	params := Dict{
		"AQB":      "1",
		"ndh":      "1",
		"t":        Sprintf("%d/%d/%d %d:%d:%d %d -480", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second(), now.Weekday()),
		"ce":       "UTF-8",
		"cl":       query.Get("cl"),
		"pageName": pageName,
		"g":        resp.URL.String(),
		"ch":       channel,
		"h5":       reportingSuite,
		"c12":      resp.Request.Header.Get("User-Agent"),
		"v12":      resp.Request.Header.Get("User-Agent"),
		"c22":      prop22,
		"v22":      eVar22,
		"s":        Sprintf("%dx%d", self.width, self.height),
		"c":        Sprintf("%d", self.colorDepth),
		"v":        "Y",
		"k":        "Y",
		"bw":       Sprintf("%d", self.innerWidth),
		"bh":       Sprintf("%d", self.innerHeight),
		"AQE":      "1",
		"sfcustom": "1",
	}

	referer := resp.Request.Header.Get("Referer")
	if len(referer) != 0 {
		params["r"] = referer
	}

	salt := Sprintf("s%014d", int64(math.Floor(float64(now.UnixNano())/float64(time.Microsecond)/10800000.0))%10+int64(math.Floor(rand.Float64()*1e13)))

	r = self.get(
		"https://securemetrics.apple.com/b/ss/applesuperglobal/1/H.20.3/"+salt,
		Dict{
			"params": params,
			"headers": Dict{
				"Accept":          "*/*",
				"Accept-Encoding": "gzip",
				"Referer":         resp.Request.URL.String(),
			},
		},
	)
}

func (self *AppleIdRegisteriOS) Signup() int {
	var exp *Exception
	var result int
	var reason string

	start := time.Now()
	success := false

	exp = Try(func() {
		//success = self.stepSignupWizard()
		success = self.stepBuyProduct()
	})

	result = SIGNUP_FAILED
	self.account.CreationTime = globals.GetCurrentTime()

	self.debug("%v", exp)

	switch {
	case success:
		result = SIGNUP_SUCCESS

	case exp == nil:
		result = SIGNUP_FAILED
		reason = "no exception but failed"

	default:
		switch e := exp.Value.(type) {
		case *AlreadyExistsError:
			reason = e.Message
			result = SIGNUP_EXISTS

		case *RegisterError:
			reason = e.Message
			result = SIGNUP_BANNED

		case *GenericError:
			reason = e.Message
			result = SIGNUP_FAILED

		case *http.HttpError:
			Raise(exp)

		default:
			msg := String(exp.Message)
			switch {
			case msg.Contains("can't find attribute"),
				msg.Contains("index out of range"):

				reason = exp.Message
				result = SIGNUP_FAILED

			default:
				self.debug("what ? debug me\n%v", exp)
				utility.SetTitle("debug me")
				level := logger.Level()
				logger.SetLevel(9999)
				utility.Pause("debug me")
				logger.SetLevel(level)
			}
		}
	}

	elapsed := time.Now().Sub(start)
	self.info("signup_result:%s proxy:%s reason:%s elapsed:%v", signupResultText[result], self.proxy, reason, elapsed)

	return result
}

func (self *AppleIdRegisteriOS) stepBuyProduct() bool {
	self.info("stepBuyProduct")

	bodyfmt := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>appExtVrsId</key>
	<string>820595060</string>
	<key>guid</key>
	<string>0fb84ae483ca4a19e9e84547c36473ec29d95c2e</string>
	<key>pg</key>
	<string>default</string>
	<key>price</key>
	<string>0</string>
	<key>pricingParameters</key>
	<string>STDQ</string>
	<key>productType</key>
	<string>C</string>
	<key>salableAdamId</key>
	<string>444934666</string>
	<key>vid</key>
	<string>%s</string>
</dict>
</plist>
`
	body := Sprintf(bodyfmt, self.buyVID)

	resp := self.post(
		self.bag[`buyProduct`],
		Dict{
			"headers": Dict{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			"body": body,
		},
	)

	self.debug("resp text: \n%v", resp.Text())

	// TODO: bag
	resp = self.get(
		"https://xp.apple.com/register",
		Dict{
			"headers": Dict{
				"Accept": "*/*",
			},
		},
	)

	self.debug("resp text: \n%v", resp.Text())

	return self.stepSignupWizard()
}

func (self *AppleIdRegisteriOS) stepSignupWizard() bool {
	self.info("stepSignupWizard")

	signature := base64.StdEncoding.EncodeToString(self.sapSession.CreatePrimeSignature())

	resp := self.get(
		self.bag[`signup`],
		Dict{
			"headers": Dict{
				`X-Apple-ActionSignature`: signature,
				"X-Apple-Partner":         "origin.0",
				"Accept-Language":         "zh-Hans",
			},
			"params": Dict{
				"guid":    self.machineGUID,
				"product": "productType=C&price=0&salableAdamId=444934666&pricingParameters=STDQ&pg=default&appExtVrsId=820595060&vid=" + self.buyVID,
			},
		},
	)

	sig, err := base64.StdEncoding.DecodeString(resp.Header.Get("X-Apple-ActionSignature"))
	RaiseIf(err)

	self.sapSession.VerifyPrimeSignature(sig)

	logger.Debug("resp body: \n%v", resp.Text())

	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	ids := Array{
		"storeFrontField",
		"hiddenBottomRightButtonId",
	}

	data := getAllElement(doc.Selection, ids)
	if len(data) != len(ids) {
		Raise(NewGenericError("len(data) != len(ids), resp incorrect"))
	}

	form := Dict{
		data["storeFrontField"].Attr2("name"):           getInputValue(data["storeFrontField"]),
		data["hiddenBottomRightButtonId"].Attr2("name"): getInputValue(data["hiddenBottomRightButtonId"]),
	}

	return self.stepTermsAndConditions(resp.Request.URL, doc.Find("form[method=post]").Attr2("action"), form)
}

func (self *AppleIdRegisteriOS) stepTermsAndConditions(refer *urllib.URL, next string, form Dict) bool {
	self.info(`stepTermsAndConditions`)

	resp := self.post(
		buildNext(refer, next),
		Dict{
			"headers": Dict{
				"Referer":         refer.String(),
				"Origin":          buildOrigin(refer),
				"X-Apple-Partner": "origin.0",
				"Content-Type":    "application/x-www-form-urlencoded",
			},
			"body": form,
		},
	)

	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	ids := Array{
		"pageUUID",
		"machineGUID",
		"signature",
		"hiddenBottomRightButtonId",
	}

	data := getAllElement(doc.Selection, ids)
	if len(data) != len(ids) {
		Raise(NewGenericError("len(data) != len(ids), resp incorrect"))
	}

	form = Dict{
		data["pageUUID"].Attr2("name"):                  data["pageUUID"].Attr2("value"),
		data["machineGUID"].Attr2("name"):               self.machineGUID,
		data["signature"].Attr2("name"):                 base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data["pageUUID"].Attr2("value")) + self.machineGUID).Encode(CP_UTF8))),
		data["hiddenBottomRightButtonId"].Attr2("name"): data["hiddenBottomRightButtonId"].Attr2("value"),
	}

	return self.stepProvideAppleIdDetails(resp.Request.URL, doc.Find("form[method=post]").Attr2("action"), form)
}

func (self *AppleIdRegisteriOS) stepProvideAppleIdDetails(refer *urllib.URL, next string, form Dict) bool {
	self.info(`stepProvideAppleIdDetails`)

	resp := self.post(
		buildNext(refer, next),
		Dict{
			"headers": Dict{
				"Referer":         refer.String(),
				"Origin":          buildOrigin(refer),
				"X-Apple-Partner": "origin.0",
				"Content-Type":    "application/x-www-form-urlencoded",
			},
			"body": form,
		},
	)

	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	ids := Array{
		"pageUUID",
		"machineGUID",
		"signature",

		"accountNameField",
		"passwordField",
		"passwordVerificationField",

		"questionField1Input",
		"answerField1",
		"questionField2Input",
		"answerField2",
		"questionField3Input",
		"answerField3",

		"recoveryEmailField",

		"birthYear",
		"birthMonthPopup",
		"birthDayPopup",

		"hiddenBottomRightButtonId",
	}

	data := getAllElement(doc.Selection, ids)
	if len(data) != len(ids) {
		Raise(NewGenericError("len(data) != len(ids), resp incorrect"))
	}

	questionsDictionary := String(doc.Find("script[type*=text][type*=javascript]:contains('questionsDictionary')").Text())
	questionsDictionary = "{" + questionsDictionary.Split("{", 1)[1].RSplit(";", 1)[0]

	type Question struct {
		Id       int
		Question string
	}

	type Questions struct {
		SecurityQuestions        []Question
		SecurityQuestionsPerPage int
	}

	questions := Questions{}
	RaiseIf(json.Unmarshal(questionsDictionary.Encode(CP_UTF8), &questions))

	selectQuestion := func(group int) Question {
		return questions.SecurityQuestions[group*questions.SecurityQuestionsPerPage+random.IntRange(0, questions.SecurityQuestionsPerPage)]
	}

	form = Dict{
		//
		// signature
		//

		data["pageUUID"].Attr2("name"):    data["pageUUID"].Attr2("value"),
		data["machineGUID"].Attr2("name"): self.machineGUID,
		data["signature"].Attr2("name"):   base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data["pageUUID"].Attr2("value")) + self.machineGUID + String(self.account.UserName)).Encode(CP_UTF8))),

		//
		// apple id and password
		//

		data["accountNameField"].Attr2("name"):          self.account.UserName,
		data["passwordField"].Attr2("name"):             self.account.ApplePassword,
		data["passwordVerificationField"].Attr2("name"): self.account.ApplePassword,

		//
		// 3 questions and answers
		//

		data["questionField1Input"].Attr2("name"): selectQuestion(0).Id,
		data["answerField1"].Attr2("name"):        utility.GenerateAnswer(),
		data["questionField2Input"].Attr2("name"): selectQuestion(1).Id,
		data["answerField2"].Attr2("name"):        utility.GenerateAnswer(),
		data["questionField3Input"].Attr2("name"): selectQuestion(2).Id,
		data["answerField3"].Attr2("name"):        utility.GenerateAnswer(),

		//
		// recoveryEmail
		//

		data["recoveryEmailField"].Attr2("name"): "",

		//
		// apple id and password
		//

		data["birthYear"].Attr2("name"):       random.IntRange(1980, 2000),
		data["birthMonthPopup"].Attr2("name"): selectOption(data["birthMonthPopup"]),
		data["birthDayPopup"].Attr2("name"):   selectOption(data["birthDayPopup"], Dict{"maxIndex": 25}),

		//
		// continue
		//

		data["hiddenBottomRightButtonId"].Attr2("name"): data["hiddenBottomRightButtonId"].Attr2("value"),
	}

	self.account.Question1 = Sprintf("%v", form[data["questionField1Input"].Attr2("name")])
	self.account.Answer1 = Sprintf("%v", form[data["answerField1"].Attr2("name")])
	self.account.Question2 = Sprintf("%v", form[data["questionField2Input"].Attr2("name")])
	self.account.Answer2 = Sprintf("%v", form[data["answerField2"].Attr2("name")])
	self.account.Question3 = Sprintf("%v", form[data["questionField3Input"].Attr2("name")])
	self.account.Answer3 = Sprintf("%v", form[data["answerField3"].Attr2("name")])
	self.account.RecoveryEmail = Sprintf("%v", form[data["recoveryEmailField"].Attr2("name")])
	self.account.Birth = Sprintf("%v-%v-%v", form[data["birthYear"].Attr2("name")], form[data["birthMonthPopup"].Attr2("name")], form[data["birthDayPopup"].Attr2("name")])

	return self.stepProvideAPaymentMethod(resp.Request.URL, doc.Find("form[method=post]").Attr2("action"), form)
}

func (self *AppleIdRegisteriOS) stepProvideAPaymentMethod(refer *urllib.URL, next string, form Dict) bool {
	self.info("stepProvideAPaymentMethod")

	// post quesions and id/password, return empty billing address form
	resp := self.post(
		buildNext(refer, next),
		Dict{
			"headers": Dict{
				"Referer":         refer.String(),
				"Origin":          buildOrigin(refer),
				"X-Apple-Partner": "origin.0",
				"Content-Type":    "application/x-www-form-urlencoded",
			},
			"body": form,
		},
	)

	// TODO: check registered
	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	ids := Array{
		"pageUUID",
		"machineGUID",
		"signature",

		"longName",
		"paymentTypeField",
		// sp
		// res

		"phone_regex",
		"email_regex",

		"creditCardNumberField",
		"verificationNumberField",
		"expirationMonthField",
		"expirationYearField",
		"mobileNumberField",
		"card_type_id",

		"sesame-id-input",
		"national-id-input",
		"country-code",

		"kddr-std",
		"codeRedemptionField",

		"lastNameField",
		"firstNameField",

		"street1Field",
		"street2Field",
		"street3Field",

		"postalCodeField",
		"cityField",
		"stateField",

		"phoneNumberField",

		"ndpd-s",
		"ndpd-f",
		"ndpd-fm",
		"ndpd-w",
		"ndpd-ipr",
		"ndpd-di",
		"ndpd-bi",

		"ndpd-probe",
		"ndpd-af",
		"ndpd-fv",
		"ndpd-fa",
		"ndpd-bp",

		"ndpd-wk",
		"ndpd-vk",
		"hidden-captcha-player-mode",

		"hiddenBottomRightButtonId",
	}

	data := getAllElement(doc.Selection, ids)
	if _, exists := data["kddr-std"]; exists == false {
		getCardTypeUrl := getAllElement(doc.Selection, Array{"getCardTypeUrl"})
		if len(getCardTypeUrl) != 0 {
			data["getCardTypeUrl"] = getCardTypeUrl["getCardTypeUrl"]
		}
	}

	if len(data) != len(ids) {
		Raise(NewGenericError("len(data) != len(ids), resp incorrect"))
	}

	script := doc.Find("script[type*=text][type*=javascript]:contains('function ndpd_load')")

	formData := map[string]interface{}{}
	widgetKeyData := map[string]interface{}{}

	RaiseIf(json.Unmarshal(String(FORM_DATA_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &formData))
	RaiseIf(json.Unmarshal(String(WIDGET_KEY_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &widgetKeyData))

	widgetKey, widgetKeyRandom := self.getWidgetKey(widgetKeyData, resp.Request.URL)

	state := selectOption(data["stateField"]) // random select?!
	stateIndex, _ := strconv.Atoi(state)
	stateName := data["stateField"].Find("option").Slice(stateIndex+1, stateIndex+2).Text()
	cityName, postalcode := utility.GeneratePostalCode(stateName)

	// prefer none instead payease
	payment := doc.Find("div[id=NONE]") // 无
	if payment == nil || payment.Length() == 0 {
		payment = doc.Find("div[id=PEAS]")
		self.debug("select PEAS") // 银行卡
	}

	form = Dict{
		//
		// signature
		//

		data["pageUUID"].Attr2("name"):    data["pageUUID"].Attr2("value"),
		data["machineGUID"].Attr2("name"): self.machineGUID,
		data["signature"].Attr2("name"):   base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data["pageUUID"].Attr2("value")) + self.machineGUID).Encode(CP_UTF8))),

		//
		// credit card info
		//

		data["longName"].Attr2("name"):         Sprintf(`%s,%s`, self.lastName, self.firstName),
		data["paymentTypeField"].Attr2("name"): payment.Attr2("input-value"),
		"sp":  "",
		"res": "",

		data[`email_regex`].Attr2(`name`): data[`email_regex`].Attr2(`value`), // `^[A-Za-z0-9._%+-]+@`,
		data[`phone_regex`].Attr2(`name`): data[`phone_regex`].Attr2(`value`), // `^[0-9-]{10,30}$`,

		data["creditCardNumberField"].Attr2("name"):   getInputValue(data["creditCardNumberField"]),
		data["verificationNumberField"].Attr2("name"): getInputValue(data["verificationNumberField"]),
		data["expirationMonthField"].Attr2("name"):    data["expirationMonthField"].Find("option").First().Attr2("value"),
		data["expirationYearField"].Attr2("name"):     data["expirationYearField"].Find("option").First().Attr2("value"),
		data["mobileNumberField"].Attr2("name"):       getInputValue(data["mobileNumberField"]),
		data["card_type_id"].Attr2("name"):            getInputValue(data["card_type_id"]),

		data["sesame-id-input"].Attr2("name"):   "",
		data["national-id-input"].Attr2("name"): "",

		//
		// code
		//

		data["codeRedemptionField"].Attr2("name"): getInputValue(data["codeRedemptionField"]),

		//
		// name
		//

		data["lastNameField"].Attr2("name"):  self.lastName,
		data["firstNameField"].Attr2("name"): self.firstName,

		//
		// address
		//

		data["street1Field"].Attr2("name"): utility.GenerateStreet1(),
		data["street2Field"].Attr2("name"): utility.GenerateStreet2(),
		data["street3Field"].Attr2("name"): utility.GenerateStreet3(),

		data["postalCodeField"].Attr2("name"): postalcode,
		data["cityField"].Attr2("name"):       cityName,
		data["stateField"].Attr2("name"):      state,

		data["phoneNumberField"].Attr2("name"): utility.GenerateMobileNumber(),

		data["country-code"].Attr2("name"): "86", // TODO: GetInput

		//
		// input
		//

		// data["ndpd-s"].Attr2("name")                        : getInputValue(data["ndpd-s"]),
		// data["ndpd-f"].Attr2("name")                        : getInputValue(data["ndpd-f"]),
		// data["ndpd-fm"].Attr2("name")                       : getInputValue(data["ndpd-fm"]),
		// data["ndpd-w"].Attr2("name")                        : getInputValue(data["ndpd-w"]),
		// data["ndpd-ipr"].Attr2("name")                      : getInputValue(data["ndpd-ipr"]),
		// data["ndpd-di"].Attr2("name")                       : getInputValue(data["ndpd-di"]),
		// data["ndpd-bi"].Attr2("name")                       : getInputValue(data["ndpd-bi"]),
		// data["ndpd-wk"].Attr2("name")                       : getInputValue(data["ndpd-wk"]),
		// data["ndpd-vk"].Attr2("name")                       : getInputValue(data["ndpd-vk"]),

		data["ndpd-s"].Attr2("name"):   formData["s"],
		data["ndpd-f"].Attr2("name"):   formData["f"],
		data["ndpd-fm"].Attr2("name"):  formData["fm"],
		data["ndpd-w"].Attr2("name"):   formData["w"],
		data["ndpd-ipr"].Attr2("name"): "", // Sprintf("ncip,0,%x,1,1;st,0,creditCardNumberField,0,verificationNumberField,0,mobileNumberField,0,codeRedemptionField,0,lastNameField,0,firstNameField,0,street1Field,0,street2Field,0,street3Field,0,postalCodeField,0,cityField,0,phoneNumberField,0;te,6fa,-1,-1;te,3bb,-1,-1;mm,4a,17e,1cc,codeRedemptionField;kk,5d,0,codeRedemptionField;ff,1,codeRedemptionField;mc,14,17e,1cc,codeRedemptionField;te,5a6,-1,-1;te,373,-1,-1,lastNameField;fb,67,codeRedemptionField;kk,2,0,lastNameField;ff,1,lastNameField;mc,4,16b,22c,lastNameField;kd,21e;kd,40;kd,8d;kd,41;kd,65;te,c2,-1,-1,firstNameField;fb,89,lastNameField;kk,3,0,firstNameField;ff,0,firstNameField;mc,7,164,25e,firstNameField;kd,1b5;kd,5a;kd,62;kd,50;kd,63;kd,c8;te,52,-1,-1,street1Field;mm,5b,159,287,street1Field;fb,20,firstNameField;kk,3,0,street1Field;ff,0,street1Field;mc,c,159,287,street1Field;kd,1e0;kd,5d;kd,3a;kd,b1;kd,57;te,da,-1,-1,street2Field;fb,7c,street1Field;kk,8,0,street2Field;ff,0,street2Field;mc,8,16c,2b4,street2Field;kd,13f;kd,4c;kd,2e;kd,a1;kd,21;kd,9d;te,86,-1,-1,street3Field;fb,84,street2Field;kk,3,0,street3Field;ff,0,street3Field;mc,7,17e,2e3,street3Field;kd,14c;kd,4e;kd,3e;kd,84;kd,2d;kd,c0;te,85,-1,-1,postalCodeField;fb,6d,street3Field;kk,a,0,postalCodeField;ff,1,postalCodeField;mc,1e,16d,305,postalCodeField;kd,67d;kd,41c;kd,19f;kd,184;kd,1c1;ts,0,3b98;kd,c1;kd,bb;kd,cc;te,10b,-1,-1,cityField;mm,49,161,336,cityField;fb,21,postalCodeField;kk,3,0,cityField;ff,0,cityField;mc,5,161,336,cityField;kd,13b;kd,70;kd,1a;kd,70;kd,20;kd,9e;te,51,-1,-1;fb,c7,cityField;mc,8,14e,36c,stateField;te,5c2,-1,-1,cityField;kk,91,6,cityField;ff,0,cityField;mc,18,135,33a,cityField;kd,55c;kd,ac;kd,8e;kd,95;kd,af;kd,96;kd,89;kd,101;kd,80;kd,1c;kd,60;kd,57;kd,58;kd,40;kd,6c;te,f6,-1,-1;te,2bc,-1,-1,phoneNumberField;mm,4a,148,3a5,phoneNumberField;fb,c7,cityField;kk,3,0,phoneNumberField;ff,0,phoneNumberField;mc,6,148,3a5,phoneNumberField;kd,373;kd,267;kd,16e;kd,29c;kd,e8;kd,d8;kd,28c;kd,2fd;kd,12c;kd,a7;kd,d8;fb,2f5,phoneNumberField;te,30a,-1,-1,navNextButton;mm,b6,158,4ab,navNextButton;mc,c,158,4ab,navNextButton;fs,2,0,0,;", globals.GetCurrentTime() / 1000),
		data["ndpd-di"].Attr2("name"):  self.deviceInfo,
		data["ndpd-bi"].Attr2("name"):  Sprintf("b2|320x480 320x460 32 32|-480|zh-cn|bp1-%s|false|%s|AppStore/2.0 iOS/7.1.2 model/iPhone3,3 build/11D257 (5; dt:70)|Not Supported", "3a4102651cc481", refer), //self.browserInfo,

		data[`ndpd-probe`].Attr2(`name`): ``, //data[`ndpd-probe`].Attr2("value"),
		data[`ndpd-af`].Attr2(`name`):    ``, //data[`ndpd-af`].Attr2("value"),
		data[`ndpd-fv`].Attr2(`name`):    `fv,mp4`,
		data[`ndpd-fa`].Attr2(`name`):    `fa,mpeg,wav`,
		data[`ndpd-bp`].Attr2(`name`):    `p,QuickTime Plug-in`,

		data["ndpd-wk"].Attr2("name"): widgetKey,
		data["ndpd-vk"].Attr2("name"): getInputValue(data["ndpd-vk"]),
		"ndpd-wkr":                    widgetKeyRandom,

		data["hidden-captcha-player-mode"].Attr2("name"): getInputValue(data["hidden-captcha-player-mode"]),
		data["hiddenBottomRightButtonId"].Attr2("name"):  getInputValue(data["hiddenBottomRightButtonId"]),
	}

	if kddr_std, exists := data["kddr-std"]; exists != false {
		form[kddr_std.Attr2("name")] = getInputValue(kddr_std)
	} else {
		getCardTypeUrl := data["getCardTypeUrl"]
		form[getCardTypeUrl.Attr2("name")] = getInputValue(getCardTypeUrl)
	}

	// fuck!!! not generator.generateiOS!!!
	form[data["ndpd-ipr"].Attr2("name")], self.retry0 = generateiOSInputRecords(doc.Selection, data, form)

	return self.stepVerifyApppleId(resp.Request.URL, doc.Find("form[method=post]").Attr2("action"), form)
}

func (self *AppleIdRegisteriOS) stepVerifyApppleId(refer *urllib.URL, next string, form Dict) bool {
	self.info("stepVerifyApppleId")

	self.debug("sleep %v ms", self.retry0)
	time.Sleep(time.Duration(self.retry0) * time.Millisecond)

	self.debug("form\n%v", form)

	// TODO: merge form and retry
	var resp *http.Response
	var doc *goquery.Document
	var elapsed int64

	// extend next to full url path
	next = buildNext(refer, next)

	for i := 0; ; i++ {

		resp = self.post(
			next,
			Dict{
				//"X-Apple-Partner": "origin.0",
				"headers": Dict{
					"Referer":      refer.String(),
					"Origin":       buildOrigin(refer),
					"Content-Type": "application/x-www-form-urlencoded",
				},
				"body": form,
			},
		)

		logger.Debug("resp body: \n%v", resp.Text())

		doc = utility.ParseHTML(resp.Text())
		exp := Try(func() {
			self.checkPageError(doc)
		})

		if exp == nil {
			break
		}

		switch {
		case i == 5,
			IsError(exp.Value, RegisterError{}) == false:
			Raise(exp)
		}

		self.debug("verifyApppleId retry %d", i+1)
		// TODO: re-create form from resp?!
		ids := Array{
			"pageUUID",
			"machineGUID",
			"signature",

			"longName",
			"paymentTypeField",
			// sp
			// res

			"phone_regex",
			"email_regex",

			"creditCardNumberField",
			"verificationNumberField",
			"expirationMonthField",
			"expirationYearField",
			"mobileNumberField",
			"card_type_id",

			"sesame-id-input",
			"national-id-input",
			"country-code",

			"kddr-std",
			"codeRedemptionField",

			"lastNameField",
			"firstNameField",

			"street1Field",
			"street2Field",
			"street3Field",

			"postalCodeField",
			"cityField",
			"stateField",

			"phoneNumberField",

			"ndpd-s",
			"ndpd-f",
			"ndpd-fm",
			"ndpd-w",
			"ndpd-ipr",
			"ndpd-di",
			"ndpd-bi",

			"ndpd-probe",
			"ndpd-af",
			"ndpd-fv",
			"ndpd-fa",
			"ndpd-bp",

			"ndpd-wk",
			"ndpd-vk",
			"hidden-captcha-player-mode",

			"hiddenBottomRightButtonId",
		}

		form2, elements := self.getProvidePaymentInfo(doc, resp, ids, nil)

		form.MergeFrom(form2)

		//generator := inputRecord.New(inputRecord.LayoutForCountry(self.account.Country))
		//form[elements["ndpd-ipr"].Attr2("name")] = generator.GenerateiOSRetry(elements, inputs, form)
		form[elements["ndpd-bi"].Attr2("name")] = Sprintf("b2|320x480 320x460 32 32|-480|zh-cn|bp1-%s|false|%s|AppStore/2.0 iOS/7.1.2 model/iPhone3,3 build/11D257 (5; dt:70)|Not Supported", "3a4102651cc481", refer)
		form[elements["ndpd-ipr"].Attr2("name")], elapsed = generateiOSInputRecordsRetry(doc.Selection, elements, form)

		// need update refer? next is full url now
		refer.Parse(next)
		// next = current root path + current post target, used as post url in next loop
		next = buildNext(refer, doc.Find("form[method=post]").Attr2("action"))

		self.debug("wait retry, sleep %v ms", elapsed)
		time.Sleep(time.Duration(elapsed) * time.Millisecond)

		//return false
	}

	self.debug("find email-verification")

	success := doc.Find("body[id=MHSignupSuccess]").Length() != 0
	if success == false {
		self.debug("can't find body MHSignupSuccess")
		return success
	}

	script := doc.Find("script[id=protocol][type*=text][type*=plist]")
	if script.Length() == 0 {
		self.debug("can't find dsid")
		return true
	}

	plist := Dict{}
	err := plistlib.Unmarshal(
		String("<plist version=\"1.0\">\r\n"+script.Text()+"</plist>").Encode(CP_UTF8),
		&plist,
	)

	if err != nil {
		self.debug("unmarshal plist failed: %v", err)
		return true
	}

	dsid, err := strconv.ParseInt(plist["set-current-user"].(map[string]interface{})["dsPersonId"].(string), 10, 64)
	RaiseIf(err)

	self.account.Dsid = dsid

	self.debug("found dsid: %d", dsid)

	return true
}

func (self *AppleIdRegisteriOS) getProvidePaymentInfo(doc *goquery.Document, resp *http.Response, elementIds Array, callback fillFormCallback) (Dict, elementDict) {
	data := getAllElement(doc.Selection, elementIds)
	if _, exists := data["kddr-std"]; exists == false {
		getCardTypeUrl := getAllElement(doc.Selection, Array{"getCardTypeUrl"})
		if len(getCardTypeUrl) != 0 {
			data["getCardTypeUrl"] = getCardTypeUrl["getCardTypeUrl"]
		}
	}

	if len(data) != len(elementIds) {
		Raise(NewGenericError("len(data) != len(elementIds), retry doc incorrect"))
	}

	script := doc.Find("script[type*=text][type*=javascript]:contains('function ndpd_load')")

	formData := map[string]interface{}{}
	widgetKeyData := map[string]interface{}{}

	RaiseIf(json.Unmarshal(String(FORM_DATA_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &formData))
	RaiseIf(json.Unmarshal(String(WIDGET_KEY_PATTERN.FindStringSubmatch(script.Text())[1]).Encode(CP_UTF8), &widgetKeyData))

	widgetKey, widgetKeyRandom := self.getWidgetKey(widgetKeyData, resp.Request.URL)

	form := Dict{
		//
		// signature
		//

		data["pageUUID"].Attr2("name"):    data["pageUUID"].Attr2("value"),
		data["machineGUID"].Attr2("name"): self.machineGUID, // hidden empty?
		data["signature"].Attr2("name"):   base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(data["pageUUID"].Attr2("value")) + self.machineGUID).Encode(CP_UTF8))),

		//
		// credit card info
		//

		// is empty hidden
		data["longName"].Attr2("name"):         Sprintf(`%s,%s`, self.lastName, self.firstName), //getInputValue(data["longName"]),
		data["paymentTypeField"].Attr2("name"): data["paymentTypeField"].Attr2(`value`),         // use value instead assume
		"sp":  "",
		"res": "",

		data[`email_regex`].Attr2(`name`): data[`email_regex`].Attr2(`value`), // `^[A-Za-z0-9._%+-]+@`,
		data[`phone_regex`].Attr2(`name`): data[`phone_regex`].Attr2(`value`), // `^[0-9-]{10,30}$`,

		data["creditCardNumberField"].Attr2("name"):   getInputValue(data["creditCardNumberField"]),
		data["verificationNumberField"].Attr2("name"): getInputValue(data["verificationNumberField"]),
		data["expirationMonthField"].Attr2("name"):    data["expirationMonthField"].Find("option").First().Attr2("value"),
		data["expirationYearField"].Attr2("name"):     data["expirationYearField"].Find("option").First().Attr2("value"),
		data["mobileNumberField"].Attr2("name"):       getInputValue(data["mobileNumberField"]),
		data["card_type_id"].Attr2("name"):            getInputValue(data["card_type_id"]),

		data["sesame-id-input"].Attr2("name"):   getInputValue(data["sesame-id-input"]),   // empty
		data["national-id-input"].Attr2("name"): getInputValue(data["national-id-input"]), // empty

		//
		// code
		//

		data["codeRedemptionField"].Attr2("name"): getInputValue(data["codeRedemptionField"]),

		//
		// name
		//

		data["lastNameField"].Attr2("name"):  getInputValue(data["lastNameField"]),
		data["firstNameField"].Attr2("name"): getInputValue(data["firstNameField"]),

		//
		// address
		//

		data["street1Field"].Attr2("name"): getInputValue(data["street1Field"]),
		data["street2Field"].Attr2("name"): getInputValue(data["street2Field"]),
		data["street3Field"].Attr2("name"): getInputValue(data["street3Field"]),

		data["postalCodeField"].Attr2("name"): getInputValue(data["postalCodeField"]),
		data["cityField"].Attr2("name"):       getInputValue(data["cityField"]),
		//data["stateField"].Attr2("name"):      getInputValue(data["stateField"]),

		data["phoneNumberField"].Attr2("name"): getInputValue(data["phoneNumberField"]),

		data["country-code"].Attr2("name"): "86",

		//
		// input
		//

		// data["ndpd-s"].Attr2("name")                        : getInputValue(data["ndpd-s"]),
		// data["ndpd-f"].Attr2("name")                        : getInputValue(data["ndpd-f"]),
		// data["ndpd-fm"].Attr2("name")                       : getInputValue(data["ndpd-fm"]),
		// data["ndpd-w"].Attr2("name")                        : getInputValue(data["ndpd-w"]),
		// data["ndpd-ipr"].Attr2("name")                      : getInputValue(data["ndpd-ipr"]),
		// data["ndpd-di"].Attr2("name")                       : getInputValue(data["ndpd-di"]),
		// data["ndpd-bi"].Attr2("name")                       : getInputValue(data["ndpd-bi"]),
		// data["ndpd-wk"].Attr2("name")                       : getInputValue(data["ndpd-wk"]),
		// data["ndpd-vk"].Attr2("name")                       : getInputValue(data["ndpd-vk"]),

		data["ndpd-s"].Attr2("name"):   formData["s"],
		data["ndpd-f"].Attr2("name"):   formData["f"],
		data["ndpd-fm"].Attr2("name"):  formData["fm"],
		data["ndpd-w"].Attr2("name"):   formData["w"],
		data["ndpd-ipr"].Attr2("name"): "", // will filled with retryInputRecord.
		data["ndpd-di"].Attr2("name"):  self.deviceInfo,
		data["ndpd-bi"].Attr2("name"):  self.browserInfo, // step url may change, current value is p

		data[`ndpd-probe`].Attr2(`name`): ``, //data[`ndpd-probe`].Attr2("value"),
		data[`ndpd-af`].Attr2(`name`):    ``, //data[`ndpd-af`].Attr2("value"),
		data[`ndpd-fv`].Attr2(`name`):    `fv,mp4`,
		data[`ndpd-fa`].Attr2(`name`):    `fa,mpeg,wav`,
		data[`ndpd-bp`].Attr2(`name`):    `p,QuickTime Plug-in`,

		data["ndpd-wk"].Attr2("name"): widgetKey, // captcha
		data["ndpd-vk"].Attr2("name"): getInputValue(data["ndpd-vk"]),
		"ndpd-wkr":                    widgetKeyRandom,

		data["hidden-captcha-player-mode"].Attr2("name"): getInputValue(data["hidden-captcha-player-mode"]),
		data["hiddenBottomRightButtonId"].Attr2("name"):  getInputValue(data["hiddenBottomRightButtonId"]),
	}

	if kddr_std, exists := data["kddr-std"]; exists != false {
		form[kddr_std.Attr2("name")] = getInputValue(kddr_std)
	} else {
		getCardTypeUrl := data["getCardTypeUrl"]
		form[getCardTypeUrl.Attr2("name")] = getInputValue(getCardTypeUrl)
	}

	//form[data["ndpd-ipr"].Attr2("name")] = generateiOSInputRecords(doc.Selection, data, form)

	if callback != nil {
		callback(form, data)
	}

	return form, data
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

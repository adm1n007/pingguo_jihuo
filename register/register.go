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
//	"strconv"
	"time"

	"ml/logging/logger"
	"ml/net/http"
	"ml/strings"
	"ml/uuid"

	"github.com/PuerkitoBio/goquery"

	"ituneslib"

	"account"
	"globals"
	"proxy"
	"utility"
	// "inputRecord"
)

var (
	FORM_DATA_PATTERN      = regexp.MustCompile(`"fd"\s*:\s*({.*?})`)
	WIDGET_KEY_PATTERN     = regexp.MustCompile(`"wk"\s*:\s*({.*?})`)
	STRIP_CALLBACK_PATTERN = regexp.MustCompile(`.*\(({.*?})\).*`)
)

type Question struct {
	Id       int
	Question string
}

type Questions struct {
	SecurityQuestions        []Question
	SecurityQuestionsPerPage int
	QuestionsIdTable         map[string]int
	QuestionsTextTable       map[int]string
}

type AppleIdRegister struct {
	session     *http.Session
	bag         Dict
	sapSession  *ituneslib.SapSession
	account     *account.AppleAccount
	machineGUID String
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
	proxy       *proxy.Proxy
}

func createHtppSession(userAgent string, country ituneslib.CountryID, proxy *proxy.Proxy) *http.Session {
	session := http.NewSession()
	session.DefaultOptions.AutoRetry = true
	session.DefaultOptions.Ignore404 = false

	RaiseIf(session.SetProxy(proxy.Host, proxy.Port, proxy.User, proxy.Password))

	if globals.Debugging && globals.UseFiddler {
		session.SetProxy("localhost", 80)
	}

	session.SetHeaders(Dict{
		"User-Agent":          userAgent,
		"Accept-Encoding":     "gzip",
		"Accept-Language":     "zh-cn, zh;q=0.75, en-us;q=0.50, en;q=0.25",
		"X-Apple-Store-Front": country.StoreFront(),
		"X-Apple-Tz":          "28800",
	})

	return session
}

func NewRegisterWin(account *account.AppleAccount, proxy *proxy.Proxy) *AppleIdRegister {
	session := createHtppSession("iTunes/12.3 (Windows; Microsoft Windows 8.1 x64 Business Edition (Build 9200); x64) AppleWebKit/7601.1056.1.1", account.Country, proxy)

	register := &AppleIdRegister{
		session:    session,
		bag:        Dict{},
		sapSession: ituneslib.NewSapSession(),
		account:    account,
		proxy:      proxy,
	}
	register.init()

	setProxy := func(sap *ituneslib.SapSession) {
		RaiseIf(sap.HttpSession.SetProxy(proxy.Host, proxy.Port, proxy.User, proxy.Password))
	}

	setProxy(register.sapSession)

	return register
}

func NewRegisterOSX(account *account.AppleAccount, proxy *proxy.Proxy) *AppleIdRegister {
	ua := "iTunes/12.3 (Macintosh; OS X 10.10.5) AppleWebKit/600.8.9"
	ua = "iTunes/12.3.2 (Macintosh; OS X 10.10.5) AppleWebKit/600.8.9"
	session := createHtppSession(ua, account.Country, proxy)

	register := &AppleIdRegister{
		session:    session,
		bag:        Dict{},
		sapSession: ituneslib.NewSapSession(),
		account:    account,
		proxy:      proxy,
	}

	register.init()

	register.machineGUID = utility.GenerateOSXMachineGuid().ToUpper()
	// register.machineGUID = "A820664EF2C8"
	// register.machineName += "’s MacBook Pro"

	setProxy := func(sap *ituneslib.SapSession) {
		RaiseIf(sap.HttpSession.SetProxy(proxy.Host, proxy.Port, proxy.User, proxy.Password))
	}

	setProxy(register.sapSession)
	// confilict with initsap?!
	//register.sapSession.Initialize(ituneslib.DefaultOSXUserAgent, ituneslib.CountryID_China, ituneslib.SAP_TYPE_REGISTER)

	return register
}

func (self *AppleIdRegister) init() {
	if len(self.account.ApplePassword) == 0 {
		self.account.CreateRandomPassword()
	}

	width, height, colorDepth, innerWidth, innerHeight, lang := utility.GenerateBrowserInfo()

	self.machineName = utility.GeneratePinyin()
	self.machineGUID = utility.GenerateMachineGuid()
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

func (self *AppleIdRegister) Close() {
	self.session.Close()
	self.sapSession.Close()
}

func (self *AppleIdRegister) info(format interface{}, args ...interface{}) {
	logger.Info("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *AppleIdRegister) debug(format interface{}, args ...interface{}) {
	logger.Debug("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *AppleIdRegister) request(method, url interface{}, params ...Dict) (resp *http.Response) {
	var headers Dict

	switch len(params) {
	case 1:
		p := params[0]
		if h, ok := p["headers"]; ok {
			headers = h.(Dict)
		} else {
			headers = Dict{}
		}

	case 0:
		params = append(params, Dict{})
		headers = Dict{}
	}

	now := time.Now().UTC()

	headers["X-Apple-I-Client-Time"] = now.Format("2006-01-02T15:04:05Z")
	headers["Date"] = now.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	params[0]["headers"] = headers

	return self.session.Request(method, url, params...)

	timeoutCount := 0

	for {
		self.debug("request: %v", url)
		exp := Try(func() { resp = self.session.Request(method, url, params...) })

		if exp != nil {
			self.debug("%v", exp)
			e := exp.Value.(*http.HttpError)

			switch e.Type {
			case http.HTTP_ERROR_TIMEOUT:
				timeoutCount += 1
				if timeoutCount > 5 {
					e.Type = http.HTTP_ERROR_CONNECT_PROXY
					Raise(exp)
				}
				fallthrough

			case http.HTTP_ERROR_INVALID_RESPONSE,
				http.HTTP_ERROR_BAD_GATE_WAY,
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

func (self *AppleIdRegister) get(url interface{}, params ...Dict) *http.Response {
	return self.request("GET", url, params...)
}

func (self *AppleIdRegister) post(url interface{}, params ...Dict) *http.Response {
	return self.request("POST", url, params...)
}

func (self *AppleIdRegister) checkPageError(doc *goquery.Document) {
	pageError := doc.Find(".page-error")
	if pageError.Length() == 0 || pageError.Children().Length() == 0 {
		return
	}

	genericError := pageError.Find(".generic-error")
	if genericError.Length() != 0 {
		text := String(genericError.Text())
		lower := text.ToLower()

		switch {
		case text.Contains("您输入的电子邮件地址已被用于 Apple 帐户"),
			lower.Contains("the email address you entered is already being used for an apple account"):
			Raise(NewAlreadyExistsError(text.String()))

		// 可能有问题, 不是被ban
		case text.Contains("您输入的电子邮件地址无效"):
			Raise(NewAccountBannedError(text.String()))

		case text.Contains("您需要设定更复杂的密码"),
			text.Contains("连续相同的字符"),
			lower.Contains("identical characters"),
			lower.Contains("a more complex password is required"):
			Raise(NewPasswordTooSimpleError(text.String()))
		}

		Raise(NewGenericError(text.String()))
	}

	err := pageError.Find(".error")
	if err.Length() != 0 {
		text := String(err.Text())
		switch {
		case text.Contains("此电子邮件地址不可用。请使用其他地址重试。"),
			text.Contains("此电子邮件地址不可用作"):
			Raise(NewAccountBannedError(text.String()))
		}

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

func (self *AppleIdRegister) getWidgetKeyRandom() int {
	return int(math.Floor(rand.Float64()*1000000) + 1000)
}

func (self *AppleIdRegister) getWidgetKey(widgetKeyData map[string]interface{}, refer *urllib.URL) (widgetKey string, widgetKeyRandom int) {
	widgetKeyRandom = self.getWidgetKeyRandom()

	resp := self.get(
		widgetKeyData["r"].(string),
		Dict{
			"Referer": refer.String(),
			"params": Dict{
				"r":  widgetKeyRandom,
				"wt": widgetKeyData["w"].(string),
			},
		},
	)

	data := map[string]interface{}{}
	RaiseIf(json.Unmarshal(String(STRIP_CALLBACK_PATTERN.FindStringSubmatch(string(resp.Text()))[1]).Encode(CP_UTF8), &data))

	widgetKey = Sprintf("%d", int(data["wk"].(float64)))

	return
}

func (self *AppleIdRegister) AgreeNewTermsAndConditions() {
	self.initbag()
	self.initsap(ituneslib.SAP_TYPE_LOGIN)

	self.session.SetHeaders(Dict{})
	self.session.SetHeaders(Dict{
		"User-Agent":          "iTunes/12.3 (Windows; Microsoft Windows 10 x64 Pro (Build 10547); x64) AppleWebKit/7601.1056.1.1",
		"Accept-Encoding":     "gzip",
		"Accept-Language":     "zh-cn, zh;q=0.75, en-us;q=0.50, en;q=0.25",
		"X-Apple-Store-Front": "143465-19,32",
		"X-Apple-Tz":          "28800",
	})

	p := self.Login()

	self.debug("dsid = %v", p["dsPersonId"])

	self.session.SetHeaders(Dict{
		"X-Token": p["passwordToken"].(string),
		"X-Dsid":  p["dsPersonId"].(string),
	})

	return

	url := Format(
		`https://buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/com.apple.jingle.app.finance.DirectAction/termsPage?userInfo=shouldCancelPurchaseBatch%3Dtrue&product=salableAdamId%3D342115564%26origPage2%3DSoftware-CN-Hipstamatic%2C%20LLC-HIPSTAMATIC%20Camera-342115564%26origPageLocation%3DBuy%26needDiv%3D1%26productType%3DC%26machineName%3D{{.machineName}}%26price%3D0%26guid%3D{{.machineGuid}}%26origPage%3DSoftware-CN-Hipstamatic%2C%20LLC-HIPSTAMATIC%20Camera-342115564%26origPageCh%3DSoftware%20Pages%26pricingParameters%3DSTDQ%26appExtVrsId%3D813518897%26origPageCh2%3DSoftware%20Pages&isBuyExcepn=true`,
		Dict{
			"machineGuid": utility.GenerateMachineGuid(),
			"machineName": utility.GenerateMachineName(),
		},
	)

	resp := self.get(url.String())

	doc := utility.ParseHTML(resp.Text())

	data := getAllElement(doc.Selection, Array{"iagree", "continue"})
	form := Dict{
		data["iagree"].Attr2("name"):   getInputValue(data["iagree"]),
		data["continue"].Attr2("name"): getInputValue(data["continue"]),
	}

	self.post(
		buildNext(resp.URL, doc.Find("form[method=post]").Attr2("action")),
		Dict{
			"headers": Dict{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			"body": form,
		},
	)
}

func (self *AppleIdRegister) initbag() {
	var resp *http.Response

	resp = self.get("https://init.itunes.apple.com/bag.xml?ix=5&ign-bsn=1")

	plist := Dict{}
	resp.Plist(&plist)
	plistlib.Unmarshal(plist["bag"].([]byte), &self.bag)
}

func (self *AppleIdRegister) initsap(sapType ituneslib.SapCertType) {
	signSapSetupCert := Dict{}
	self.get(self.bag["sign-sap-setup-cert"]).Plist(&signSapSetupCert)

	// SapExchangeData return FFFF5B9B
	cert := self.sapSession.ExchangeData(sapType, signSapSetupCert["sign-sap-setup-cert"].([]byte))
	self.debug("ExchangeData phase 1")

	body, err := plistlib.MarshalIndent(Dict{"sign-sap-setup-buffer": cert}, plistlib.XMLFormat, "    ")
	RaiseIf(err)

	signSapSetupBuffer := Dict{}
	self.post(
		self.bag["sign-sap-setup"],
		Dict{
			"headers": Dict{
				"Content-Type": "application/x-apple-plist",
			},

			"body": body,
		},
	).Plist(&signSapSetupBuffer)

	self.sapSession.ExchangeData(sapType, signSapSetupBuffer["sign-sap-setup-buffer"].([]byte))
	self.sapSession.Initialized = true
}

func (self *AppleIdRegister) Initialize() {
	self.initbag()
	self.initsap(ituneslib.SAP_TYPE_REGISTER)
}

func (self *AppleIdRegister) Login() Dict {
	const loginFormStub = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
    <key>appleId</key>
    <string>{{.userName}}</string>
    <key>attempt</key>
    <integer>1</integer>
    <key>createSession</key>
    <string>true</string>
    <key>guid</key>
    <string>{{.machineGuid}}</string>
    <key>machineName</key>
    <string>{{.machineName}}</string>
    <key>password</key>
    <string>{{.password}}</string>
    <key>why</key>
    <string>signIn</string>
</dict>
</plist>`

	// var err error

	data := strings.Format(
		loginFormStub,
		Dict{
			"userName":    self.account.UserName,
			"machineGuid": utility.GenerateMachineGuid(),
			"machineName": utility.GenerateMachineName(),
			"password":    self.account.ApplePassword,
		},
	).Encode(strings.CP_UTF8)

	resp := self.post(
		self.bag["authenticateAccount"],
		Dict{
			"headers": Dict{
				"X-Apple-ActionSignature": base64.StdEncoding.EncodeToString(self.sapSession.SignData(data)),
				"Content-Type":            "application/x-apple-plist",
			},
			"body": data,
		},
	)

	p := Dict{}
	resp.Plist(&p)

	// dsid, err = strconv.ParseInt(p["dsPersonId"].(string), 10, 64)
	// RaiseIf(err)

	return p
}

func (self *AppleIdRegister) pingUrl(doc *goquery.Document, resp *http.Response, switches ...bool) {
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

func (self *AppleIdRegister) Signup() int {
	var exp *Exception
	var result int
	var reason string

	start := time.Now()
	self.proxy.PreSignupLock(String(self.account.UserName))
	success := false

	exp = Try(func() {
		success = self.stepSignupWizard()
	})

	result = SIGNUP_FAILED
	self.account.CreationTime = globals.GetCurrentTime()

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

		case *PasswordTooSimpleError:
			self.account.CreateRandomPassword()
			reason = e.Message
			result = SIGNUP_PASSWORD_TOO_SIMPLE

		case *RegisterError:
			reason = e.Message
			result = SIGNUP_BANNED

		case *AccountBannedError:
			reason = e.Message
			result = SIGNUP_EMAIL_INVALID

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
	//if result == SIGNUP_SUCCESS || result == SIGNUP_BANNED {
	self.proxy.LeftCount--
	//}
	self.proxy.PostSignupUnlock(String(self.account.UserName))
	// TODO: force a remove call?
	self.info("signup_result:%s proxy:%s reason:%s elapsed:%v", signupResultText[result], self.proxy, reason, elapsed)

	if exp != nil && result != SIGNUP_BANNED {
		self.debug("exception: %v", exp)
	}

	return result
}

func (self *AppleIdRegister) stepSignupWizard() bool {
	self.info("stepSignupWizard")

	signature := base64.StdEncoding.EncodeToString(self.sapSession.CreatePrimeSignature())
	boundary := utility.GenerateRandomString().ToUpper()

	body := String("\r\n").Join([]string{
		`--{{.boundary}}`,
		`Content-Disposition: form-data; name="uuid"`,
		``,
		`{{.uuid}}`,
		`--{{.boundary}}`,
		`Content-Disposition: form-data; name="product"`,
		``,
		`{{.product}}`,
		`--{{.boundary}}`,
		`Content-Disposition: form-data; name="machineName"`,
		``,
		`{{.machineName}}`,
		`--{{.boundary}}`,
		`Content-Disposition: form-data; name="guid"`,
		``,
		`{{.guid}}`,
		`--{{.boundary}}--`,
		``,
	}).Format(Dict{
		"boundary":    boundary,
		"guid":        self.machineGUID,
		"machineName": self.machineName,
		"product":     `productType=C&price=0&salableAdamId=903247390&pricingParameters=STDQ&pg=default&appExtVrsId=812773075`,
		"uuid":        uuid.UUIDV4().ToUpper(),
	})

	// HTTP 307 is done by http package?! because boundary is same in two post
	resp := self.post(
		self.bag[`signup`],
		Dict{
			"headers": Dict{
				`X-Apple-ActionSignature`: signature,
				`Content-Type`:            Sprintf(`multipart/form-data; boundary=%s`, boundary),
				`Connection`:              `close`,
			},
			"body": body,
		},
	)

	self.debug("headers: %v", resp.Header)

	sig, err := base64.StdEncoding.DecodeString(resp.Header.Get("X-Apple-ActionSignature"))
	RaiseIf(err)

	self.sapSession.VerifyPrimeSignature(sig)

	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	// <key>failureType</key><string>5002</string>
	//logger.Debug("resp: %v\n, resp.Text: %v", resp, resp.Text())
	continueButton := doc.Find("button.emphasized[type=button]")

	logger.Debug("continueButton: %v", continueButton)
	// (*Selection).Attr2: Raise("can't find attribute: '" + attrName + "'")
	href := continueButton.Parent().Attr2("href")

	return self.stepTermsAndConditions(resp.Request.URL, href)
}

func (self *AppleIdRegister) stepTermsAndConditions(refer *urllib.URL, next string) bool {
	self.info(`stepTermsAndConditions`)

	resp := self.get(
		buildNext(refer, next),
		Dict{
			"headers": Dict{
				"Referer":         refer.String(),
				"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"Accept-Encoding": "gzip, deflate",
				`Connection`:      `keep-alive`,
			},
		},
	)

	doc := utility.ParseHTML(resp.Text())
	self.checkPageError(doc)

	self.pingUrl(doc, resp)

	acceptTerms := doc.Find("div.accept-terms")
	action := acceptTerms.Find("form").Attr2("action")

	inputs := map[string]*goquery.Selection{}

	acceptTerms.Find("input").Each(func(i int, s *goquery.Selection) {
		inputs[s.Attr2("id")] = s
	})

	form := Dict{
		inputs["iagree"].Attr2("name"):      inputs["iagree"].Attr2("value"),
		inputs["pageUUID"].Attr2("name"):    inputs["pageUUID"].Attr2("value"),
		inputs["machineGUID"].Attr2("name"): self.machineGUID,
		inputs["signature"].Attr2("name"):   base64.StdEncoding.EncodeToString(self.sapSession.SignData((String(inputs["pageUUID"].Attr2("value")) + self.machineGUID).Encode(CP_UTF8))),
	}

	return self.stepProvideAppleIdDetails(resp.Request.URL, action, form)
}

func (self *AppleIdRegister) parseQuestionsDictionary(doc *goquery.Document) Questions {
	questionsDictionary := String(doc.Find("script[type*=text][type*=javascript]:contains('questionsDictionary')").Text())
	questionsDictionary = "{" + questionsDictionary.Split("{", 1)[1].RSplit("};", 1)[0] + "}"

	questions := Questions{
		QuestionsIdTable:   map[string]int{},
		QuestionsTextTable: map[int]string{},
	}
	RaiseIf(json.Unmarshal(questionsDictionary.Encode(CP_UTF8), &questions))

	for _, question := range questions.SecurityQuestions {
		questions.QuestionsIdTable[question.Question] = question.Id
		questions.QuestionsTextTable[question.Id] = question.Question
	}

	return questions
}

func (self *AppleIdRegister) stepVerifyApppleId(refer *urllib.URL, next string, form Dict) bool {
	/*
	self.info("stepVerifyApppleId")

	time.Sleep(globals.Preferences.Register.VerifyApppleIdDelay)

	self.debug("form\n%v", form)

	var resp *http.Response
	var doc *goquery.Document

	next = buildNext(refer, next)

	for i := 0; ; i++ {
		resp = self.post(
			next,
			Dict{
				"headers": Dict{
					"Referer":         refer.String(),
					"Origin":          buildOrigin(refer),
					"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,q=0.8",
					"Accept-Encoding": "gzip, deflate",
					"Content-Type":    "application/x-www-form-urlencoded",
					`Connection`:      `keep-alive`,
				},
				"body": form,
			},
		)

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

		form2, elements := self.getProvidePaymentInfo(doc, resp, Array{
			"country",
			"codeRedemptionField",

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
		},
			nil,
		)

		form.MergeFrom(form2)

		inputs := []*goquery.Selection{}
		getAllInput(doc.Selection).Each(func(i int, s *goquery.Selection) {
			inputs = append(inputs, s)
		})

		generator := inputRecord.New(inputRecord.LayoutForCountry(self.account.Country))
		form[elements["ndpd-ipr"].Attr2("name")] = generator.GenerateRetry(elements, inputs, form)

		// refer need re-take from current url?!
		next = buildNext(refer, doc.Find("form[method=post]").Attr2("action"))

		self.debug("wait retry %v", generator.Elapsed())
		time.Sleep(generator.Elapsed())
	}

	self.debug("find email-verification")

	success := doc.Find(".email-verification").Length() != 0
	if success == false {
		self.debug("can't find email-verification")
		return success
	}

	script := doc.Find("script[id=protocol][type*=text][type*=x-apple-plist]")
	if script.Length() == 0 {
		self.debug("can't find dsid")
		return success
	}

	plist := Dict{}
	err := plistlib.Unmarshal(
		String("<plist version=\"1.0\">\r\n"+script.Text()+"</plist>").Encode(CP_UTF8),
		&plist,
	)

	if err != nil {
		self.debug("unmarshal plist failed: %v", err)
		return success
	}

	dsid, err := strconv.ParseInt(plist["set-current-user"].(map[string]interface{})["dsPersonId"].(string), 10, 64)
	RaiseIf(err)

	self.account.Dsid = dsid

	self.debug("found dsid: %d", dsid)
*/

	return true
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

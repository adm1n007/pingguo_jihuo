package activator

import (
	. "fmt"
	. "active_apple/ml/dict"
	. "active_apple/ml/strings"
	. "active_apple/ml/trace"

	"path"
	"strings"
	"time"

	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"

	"active_apple/ituneslib"
	"active_apple/ml/encoding/base64"
	"active_apple/ml/logging/logger"
	"active_apple/ml/net/http"
	"active_apple/ml/net/pop3"

	"github.com/PuerkitoBio/goquery"

	"active_apple/account"
	"active_apple/proxy"
	"active_apple/utility"
)

type AppleIdActivator struct {
	session      *http.Session
	account      *account.AppleAccount
	proxyManager proxy.Manager
	sapSession   *ituneslib.SapSession
}

func NewActivator(account *account.AppleAccount, proxyManager proxy.Manager) *AppleIdActivator {
	session := http.NewSession()
	session.SetHeaders(Dict{
		`User-Agent`:      `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36`,
		`Accept`:          `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`,
		`Accept-Encoding`: `gzip, deflate`,
		`Accept-Language`: `zh-CN,en-US;q=0.8,en;q=0.5,zh-HK;q=0.3`,
		`Content-Type`:    `application/x-www-form-urlencoded`,
	})

	session.DefaultOptions.AutoRetry = true
	session.DefaultOptions.Ignore404 = false

	activator := &AppleIdActivator{
		session:      session,
		account:      account,
		proxyManager: proxyManager,
	}

	return activator
}

func (self *AppleIdActivator) Close() {
	self.session.Close()
	if self.sapSession != nil {
		self.sapSession.Close()
	}
}

func (self *AppleIdActivator) info(format interface{}, args ...interface{}) {
	logger.Info("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *AppleIdActivator) debug(format interface{}, args ...interface{}) {
	logger.Debug("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *AppleIdActivator) request(method, url interface{}, params ...Dict) (resp *http.Response) {
	self.debug("request %v: %v", method, url)
	for {
		resp = self.session.Request(method, url, params...)
		if resp.StatusCode == http.StatusNotFound {
			continue
		}

		return
	}
}

func (self *AppleIdActivator) get(url interface{}, params ...Dict) *http.Response {
	return self.request("GET", url, params...)
}

func (self *AppleIdActivator) post(url interface{}, params ...Dict) *http.Response {
	return self.request("POST", url, params...)
}

func (self *AppleIdActivator) createSapSession(proxyManager proxy.Manager) *ituneslib.SapSession {
	setProxy := func(sap *ituneslib.SapSession) {
		if proxyManager != nil {
			proxy := proxyManager.GetProxy(30 * time.Second)
			RaiseIf(sap.HttpSession.SetProxy(proxy.Host, proxy.Port, proxy.User, proxy.Password))
		}
	}

	if true {
		sapSession := ituneslib.NewSapSession()
		sap := sapSession

		defer func() {
			if sap != nil {
				sap.Close()
			}
		}()

		setProxy(sapSession)
		sapSession.Initialize(ituneslib.DefaultWindowsUserAgent, ituneslib.CountryID_China, ituneslib.SAP_TYPE_LOGIN)

		sap = nil
		return sapSession
	}

	if self.sapSession != nil {
		setProxy(self.sapSession)
		return self.sapSession
	}

	self.sapSession = ituneslib.NewSapSession()

	setProxy(self.sapSession)
	self.sapSession.Initialize(ituneslib.DefaultWindowsUserAgent, ituneslib.CountryID_China, ituneslib.SAP_TYPE_LOGIN)
	return self.sapSession
}

func (self *AppleIdActivator) FindVerifyUrls() (result ActivateResult, urls []String) {
	if String(self.account.UserName).ToLower().EndsWith("@163.com") {
        //帐号为163邮箱，执行登录163邮箱查找激活邮件
		return self.FindVerifyUrls163()
	}

	result = ACTIVATE_MAIL_LOGIN_FAILED

	self.debug("connect pop3")

	var client *pop3.Client
	var err error

	for i := 1; ; i++ {
		client, err = pop3.DialTLS(Sprintf("%s:995", self.account.PopAddress))
		if err == nil {
			break
		}

		self.debug("connect %v failed %d: %v", self.account.PopAddress, i, err)

		if String(err.Error()).Contains("DialWithDialer timed out") == false {
			break
		}
	}

	if err != nil {
		result = ACTIVATE_MAIL_CONNECT_FAILED
		self.debug("connect %v failed: %v", self.account.PopAddress, err)
		return
	}

	defer client.Close()

	self.debug("auth user = %s pass = %s", self.account.UserName, self.account.MailPassword)

	err = client.Auth(self.account.UserName, self.account.MailPassword)
	if err != nil {
		self.debug("auth failed: %v", err)
		return
	}

	result = ACTIVATE_EMAIL_NOT_FOUND

	_, sizes, err := client.ListAll()
	if err != nil {
		self.debug("ListAll failed: %v", err)
		return
	}

	for mailIndex := 1; mailIndex != len(sizes)+1; mailIndex++ {
		text, err := client.Retr(mailIndex)
		if err != nil {
			self.debug("Retr %d failed: %v", mailIndex, err)
			continue
		}

		message, err := mail.ReadMessage(strings.NewReader(text))
		if err != nil {
			self.debug("ReadMessage %d failed: %v", mailIndex, err)
			continue
		}

		var parts []String

		exp := Try(func() {
			parts = self.parseMessage(message)
		})

		if exp != nil {
			self.debug("parseMessage error: %v", exp)
			continue
		}

		for _, p := range parts {
			for _, line := range p.SplitLines() {
				if line.StartsWith("https://id.apple.com/cgi-bin/verify.cgi") {
					urls = append(urls, line)
					break
				}
			}
		}
	}

	if len(urls) != 0 {
		result = ACTIVATE_SUCCESS
	}

	return
}

func (self *AppleIdActivator) parseMessage(message *mail.Message) (parts []String) {
	mediaType_, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
	RaiseIf(err)

	mediaType := String(mediaType_)

	var payload String
	var slurp []byte

	switch {
	case mediaType.StartsWith("multipart/"):
		mr := multipart.NewReader(message.Body, params["boundary"])
		for {
			var part *multipart.Part

			part, err = mr.NextPart()
			if err != nil {
				return
			}

			slurp, err = ioutil.ReadAll(part)
			RaiseIf(err)

			payload, err = self.decodePayload(
				slurp,
				String(part.Header.Get("Content-Transfer-Encoding")),
				self.getCharset(part.Header),
			)
			RaiseIf(err)
			parts = append(parts, payload)
		}

	case mediaType.StartsWith("text/"):
		slurp, err = ioutil.ReadAll(message.Body)
		RaiseIf(err)

		payload, err = self.decodePayload(
			slurp,
			String(message.Header.Get("Content-Transfer-Encoding")),
			self.getCharset(message.Header),
		)
		RaiseIf(err)
		parts = append(parts, payload)
	}

	return
}

type MapGetter interface {
	Get(key string) string
}

func (self *AppleIdActivator) getCharset(header MapGetter) Encoding {
	charset := CP_UTF8
	switch String(header.Get("charset")).ToLower() {
	case "gb2312", "gbk":
		charset = CP_GBK

	case "big5":
		charset = CP_BIG5

	case "shiftjis":
		charset = CP_SHIFT_JIS
	}

	return charset
}

func (self *AppleIdActivator) decodePayload(payload []byte, encoding String, charset Encoding) (String, error) {
	switch encoding.ToLower() {
	case "base64":
		bytes := base64.DecodeString(string(payload))
		return Decode(bytes, charset), nil

	case "quoted-printable":
		r := quotedprintable.NewReader(strings.NewReader(string(payload)))
		bytes, err := ioutil.ReadAll(r)
		return Decode(bytes, charset), err
	}

	return Decode(payload, CP_UTF8), nil
}

func (self *AppleIdActivator) VerifyByUrl(url String, proxy *proxy.Proxy) (result ActivateResult) {

	if self.session.SetProxy(proxy.Host, proxy.Port, proxy.User, proxy.Password) != nil {
		return ACTIVATE_PROXY_INVALID
	}

	// self.session.SetProxy("127.0.0.1", 6789)

	self.session.SetTimeout(15 * time.Second)

	resp := self.get(url)
	text := resp.Text()

	// text = debug_html

	result = self.checkStatusFromText(text)
	// self.debug("checkStatusFromText: %d", result)
	switch result {
        case ACTIVATE_FAILED:
            break

        case ACTIVATE_LINK_EXPIRED:
            self.ResendVerificationEmail()
            return ACTIVATE_DELAY_PROCESS

        default:
            return
	}

	doc := utility.ParseHTML(text)
	form := doc.Find("form[id=command]")
	if form.Length() == 0 {
		form = doc.Find("form[id=form1]")
	}

	inputs := form.Find("input").Filter("[name]")

	data := Dict{}

	inputs.Each(func(i int, s *goquery.Selection) {
		data[s.Attr2("name")] = s.AttrOr("value", "")
	})

	if result != ACTIVATE_LINK_EXPIRED {
		delete(data, "theAccountName")
		delete(data, "accountPassword")

		data["fdDetails"] = "-480|648x1152x24x614x1152"
		data["appleId"] = self.account.UserName
		data["accountPassword"] = self.account.ApplePassword
	}

	// self.debug("%v", data)

	u := resp.Request.URL
	link := Sprintf("%s://%s%s", u.Scheme, u.Host, path.Join(path.Dir(u.Path), form.Attr2("action")))

	switch result {
	case ACTIVATE_LINK_EXPIRED:
		self.debug("open resend email link: %v", link)

	default:
		self.debug("login and open verify link: %v", link)
	}

	resp = self.post(link, Dict{"body": data})
	text = resp.Text()

	result = self.checkStatusFromText(text)
	switch result {
	case ACTIVATE_EMAIL_RESENT:
		result = ACTIVATE_DELAY_PROCESS

	case ACTIVATE_SUCCESS,
		ACTIVATE_SESSION_TIMEOUT,
		ACTIVATE_ACCOUNT_INVALID,
		ACTIVATE_CANT_ACTIVATE,
		ACTIVATE_DELAY_PROCESS:
		break

	default:
		self.debug("verify failed: %v", text)
	}

	return
}

// var (
//     SERVICE_KEY_PATTERN = regexp.MustCompile(`serviceKey\s*:\s*'(.*?)'`)
//     SESSION_ID_PATTERN  = regexp.MustCompile(`sessionId\s*:\s*'(.*?)'`)
//     SCNT_PATTERN        = regexp.MustCompile(`scnt\s*:\s*'(.*?)'`)
// )

// func (self *AppleIdActivator) VerifyByVerifyCode(proxyManager proxy.Manager) ActivateResult {
//     url := "https://appleid.apple.com/"

//     resp := self.get(url)
//     text := resp.Text()

//     serviceKey := SERVICE_KEY_PATTERN.FindStringSubmatch(text.String())[1]
//     domainId := 1

//     _ = `
//     TF1;020;;;;;;;;;;;;;;;;;;;;;;Mozilla;Netscape;5.0%20%28Windows%20NT%2010.0%3B%20WOW64%29%20AppleWebKit/537.36%20%28KHTML%2C%20like%20Gecko%29%20Chrome/48.0.2564.116%20Safari/537.36;20030107;undefined;true;;true;Win32;undefined;Mozilla/5.0%20%28Windows%20NT%2010.0%3B%20WOW64%29%20AppleWebKit/537.36%20%28KHTML%2C%20like%20Gecko%29%20Chrome/48.0.2564.116%20Safari/537.36;en-US;GBK;appleid.apple.com;undefined;undefined;undefined;undefined;false;false;1455857604897;8;6/7/2005%2C%209%3A33%3A44%20PM;1333;888;;20.0;;;;;503999;-480;-480;2/19/2016%2C%2012%3A45%3A09%20PM;24;1333;798;0;0;;;;;;Shockwave%20Flash%7CShockwave%20Flash%2020.0%20r0;;;;;;;;;;;;;23;;;;;;;;;;;;;;;5.6.1-0;;

//     TF1;020;;;;;;;;;;;;;;;;;;;;;;Mozilla;Netscape;5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36;20030107;undefined;true;;true;Win32;undefined;Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36;en-US;GBK;appleid.apple.com;undefined;undefined;undefined;undefined;false;false;1455857604897;8;6/7/2005, 9:33:44 PM;1333;888;;20.0;;;;;503999;-480;-480;2/19/2016, 12:45:09 PM;24;1333;798;0;0;;;;;;Shockwave Flash|Shockwave Flash 20.0 r0;;;;;;;;;;;;;23;;;;;;;;;;;;;;;5.6.1-0;;
//     `

//     resp = self.post("https://idmsa.apple.com/appleauth/auth/signin",
//                 Dict{
//                     "headers": Dict{
//                         "X-Requested-With"          : "XMLHttpRequest",
//                         "X-Apple-Domain-Id"         : domainId,
//                         "X-Apple-Widget-Key"        : serviceKey,
//                         "X-Apple-I-FD-Client-Info"  : `{"U":"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36","L":"en-US","Z":"GMT+08:00","V":"1.1","F":"Nla44j1e3NlY5BSo9z4ofjb75PaK4Vpjt.gEngMQEjZr_WhXTA6FL.26y8GGEDd5ihORoVyFGh8cmvSuCKzIlnY6xljQlpRD2phteB.CaJ6hO3f9p_nH1u_eH3BhxUC550ialT0iakA2zGUMnGWFfwMHDCQyFA2wv4qnvtCtNQNbXlo4V.lzXJJIneGffLMC7EZ3QHPBirTYKUowRslzRQqwSM2VxQTPY.2O_0vLG9mhORoVjnjk3nKy_.7rqbPaRgwe98vDdYejftckuyPBDjaY2ftckZZLQ084akJ8YSI9vFW0q0UfR0odm_dhrxbuJjkWxv5iK9JruxThveWhYjTLy.EKY.6ekcsjNtG.5vm_UdBzW2wHCSFQ_H_CtTrLv37lhQwMAj9htsfHOrf8M2Lz4mvmfTT9oaSzeWhyr1BNlrK1BNlYCa1nkBMfs.Bl4"}`,
//                     },
//                     "body": Dict{
//                         "accountName"   : self.account.UserName,
//                         "password"      : self.account.ApplePassword,
//                         "rememberMe"    : false,
//                     },
//                 },
//             )

//     resp = self.get("https://appleid.apple.com/widget/account/?rv=1&language=en_US_USA&widgetKey=" + serviceKey)

//     text = resp.Text()
//     sessionId := SESSION_ID_PATTERN.FindStringSubmatch(text.String())[1]
//     scnt := SCNT_PATTERN.FindStringSubmatch(text.String())[1]

//     for retry := 10; retry >= 0; retry-- {
//         result, codes := self.FindVerifyCode163()
//         time.Sleep(5 * time.Second)
//     }

//     return ACTIVATE_SUCCESS
// }

func (self *AppleIdActivator) VerifyPassword(proxyManager proxy.Manager) ActivateResult {
	// sap is safe?
	sap := self.createSapSession(proxyManager)
	defer sap.Close()
	// re-use sap httpsession (so is dangerous?)
	session := sap.HttpSession

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
</plist>
`
	data := Format(
		loginFormStub,
		Dict{
			"userName":    self.account.UserName,
			"machineGuid": utility.GenerateMachineGuid(),
			"machineName": utility.GeneratePinyin(),
			"password":    self.account.ApplePassword,
		},
	).Encode(CP_UTF8)

	// Request(method, url interface{}, params ...Dict)
	resp := session.Post(
		sap.UrlBag["authenticateAccount"],
		Dict{
			"headers": Dict{
				"X-Apple-ActionSignature": base64.EncodeToString(sap.SignData(data)),
				"Content-Type":            "application/x-apple-plist",
			},
			"body": data,
			// timeout here?
		},
	)

	p := JsonDict{}
	resp.Plist(&p)

	self.debug("login appleId result: %v", p)

	if failureType := p["failureType"]; failureType != nil {
		if failureType.(string) == "-5000" {
			return ACTIVATE_PASSWORD_ERROR
		}
	}

	if customerMessage := p["customerMessage"]; customerMessage != nil {
		msg := String(customerMessage.(string))

		switch {
		case msg.Contains("密码输入有误"):
			return ACTIVATE_PASSWORD_ERROR

		case msg.Contains("被禁用"):
			return ACTIVATE_ACCOUNT_BANNED
		}
	}

	if passwordToken := p["passwordToken"]; passwordToken != nil {
		token, ok := passwordToken.(string)

		if ok && token != "" {
			return ACTIVATE_SUCCESS
		}
	}

    //liangxu新增逻辑
    //已激活
    if p["m-allowed"].(bool) {
        self.debug("appleId login m-allowed is true")
        return ACTIVATE_SUCCESS
    }
    dialog := p.Map("dialog")
    okButtonAction := dialog.Map("okButtonAction")
    url := okButtonAction["url"].(string)
    resp = session.Get(url)
    doc := utility.ParseHTML(resp.Text())
    resendLink := doc.Find("a[class=resend-link]")
    resp = session.Get(Sprintf("%s://%s%s", resp.Request.URL.Scheme, resp.Request.URL.Host, resendLink.Attr2("href")))
    text := resp.Text()
    if self.checkStatusFromText(text) != ACTIVATE_EMAIL_RESENT {
        self.debug("login appleid resent activeemail error")
        return ACTIVATE_CANT_ACTIVATE
    } else {
        self.debug("resent active email success")
    }
    //liangxu新增逻辑

	return ACTIVATE_INVALID_PASSWORD_TOKEN
}

func (self *AppleIdActivator) ResendVerificationEmail() ActivateResult {
	proxyManager := self.proxyManager

	sap := self.createSapSession(proxyManager)
	defer sap.Close()
	session := sap.HttpSession

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
</plist>
`
	data := Format(
		loginFormStub,
		Dict{
			"userName":    self.account.UserName,
			"machineGuid": utility.GenerateMachineGuid(),
			"machineName": utility.GeneratePinyin(),
			"password":    self.account.ApplePassword,
		},
	).Encode(CP_UTF8)

	resp := session.Post(
		sap.UrlBag["authenticateAccount"],
		Dict{
			"headers": Dict{
				"X-Apple-ActionSignature": base64.EncodeToString(sap.SignData(data)),
				"Content-Type":            "application/x-apple-plist",
			},
			"body": data,
		},
	)

	p := JsonDict{}
	resp.Plist(&p)

	self.debug("ResendVerificationEmail login apple result Json: %v", p)

	if failureType := p["failureType"]; failureType != nil {
		if failureType.(string) == "-5000" {
			return ACTIVATE_PASSWORD_ERROR
		}
	}

	if p["m-allowed"].(bool) {
		self.debug("ResendVerificationEmail login result: %v", p)
		return ACTIVATE_SUCCESS
	}

	dialog := p.Map("dialog")


	okButtonAction := dialog.Map("okButtonAction")
	url := okButtonAction["url"].(string)

	resp = session.Get(url)

	doc := utility.ParseHTML(resp.Text())

	resendLink := doc.Find("a[class=resend-link]")
	resp = session.Get(Sprintf("%s://%s%s", resp.Request.URL.Scheme, resp.Request.URL.Host, resendLink.Attr2("href")))
	text := resp.Text()

	switch self.checkStatusFromText(text) {
	case ACTIVATE_EMAIL_RESENT:
		return ACTIVATE_DELAY_PROCESS

	default:
		self.debug("ResendVerificationEmail page\n%v", text)
		return ACTIVATE_CANT_ACTIVATE
	}
}

func (self *AppleIdActivator) checkStatusFromText(text String) ActivateResult {
	lower := text.ToLower()

	switch {
	case text.Contains("您的验证链接已过期"):
		return ACTIVATE_LINK_EXPIRED

	case text.Contains("验证电子邮件已发送"):
		return ACTIVATE_EMAIL_RESENT

	case text.Contains("无法验证电子邮件地址"):
		return ACTIVATE_CANT_ACTIVATE
		self.account.CantVerify = true
		return ACTIVATE_DELAY_PROCESS

	case text.Contains("您的会话已超时"),
		text.Contains("您的阶段作业已超时"):
		return ACTIVATE_SESSION_TIMEOUT

	case text.Contains("电子邮件地址之前已经验证"):
		self.debug("already verified")
		return ACTIVATE_SUCCESS

	case text.Contains("电子邮件地址已验证"),
		text.Contains("已验证电子邮件地址"),
		lower.Contains("email address verified"),
		lower.Contains("has already been verified"):
		return ACTIVATE_SUCCESS

	case text.Contains("E-Mail-Adresse bestätigt"):
		return ACTIVATE_PASSWORD_ERROR

	case text.Contains("无法完成您的请求"),
		text.Contains("无法验证此电子邮件地址，因为验证链接已过期或失效"):
		return ACTIVATE_CANT_FINISH_REQUEST

	case text.Contains("您的 Apple ID 或密码输入有误"),
		lower.Contains("your apple id or password was entered incorrectly"):
		return ACTIVATE_ACCOUNT_INVALID

	case text.Contains("此 Apple ID 已由于安全原因被禁用"):
		return ACTIVATE_CANT_ACTIVATE
	}

	return ACTIVATE_FAILED
}

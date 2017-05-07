package register

import (
    . "fmt"
    . "ml/trace"
    . "ml/dict"
    . "ml/strings"

    "time"
    "encoding/base64"
    "encoding/hex"

    "ml/logging/logger"
    "ml/net/http"
    "ituneslib"

    "account"
    "proxy"
    "utility"
)

const (
    BASE_URL    = "https://setup.icloud.com/"
)

type ICloudRegister struct {
    session     *http.Session
    sapSession  *ituneslib.SapSession
    account     *account.AppleAccount
    firstName   String
    lastName    String
    serverInfo  Dict
    proxy       *proxy.Proxy
}

func NewiCloudRegister(account *account.AppleAccount, proxy *proxy.Proxy) *ICloudRegister {
    session := http.NewSession()

    session.SetHeaders(Dict{
        "User-Agent"        : "iOS iPhone 12H143 iPhone Setup Assistant",
        "Accept-Encoding"   : "gzip, deflate",
        "Accept-Language"   : "zh-cn",
        "Accept"            : "*/*",
        "Connection"        : "keep-alive",
        "X-MMe-Client-Info" : "<iPhone5,4> <iPhone OS;8.4;12H143> <com.apple.AppleAccount/1.0 (com.apple.Preferences/1.0)>",
        "X-MMe-Country"     : "CN",
    })

    register := &ICloudRegister{
                        session     : session,
                        sapSession  : ituneslib.NewSapSession(),
                        account     : account,
                        proxy       : proxy,
                }

    register.setProxy(session)
    register.init()

    return register
}

func (self *ICloudRegister) init() {
    if len(self.account.ApplePassword) == 0 {
        self.account.ApplePassword  = string(utility.GeneratePassword())
    }

    self.firstName = utility.GenerateRandomString()[:8]
    self.lastName = utility.GenerateRandomString()[:8]
}

func (self *ICloudRegister) setProxy(session *http.Session) {
    RaiseIf(session.SetProxy(self.proxy.Host, self.proxy.Port, self.proxy.User, self.proxy.Password))
    session.SetProxy("localhost", 6789)
}

func (self *ICloudRegister) Close() {
    if self.sapSession != nil {
        self.sapSession.Close()
        self.sapSession = nil
    }

    self.session.Close()
}

func (self *ICloudRegister) info(format interface{}, args ...interface{}) {
    logger.Info("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *ICloudRegister) debug(format interface{}, args ...interface{}) {
    logger.Debug("%s", Sprintf(Sprintf("[%s] %v", self.account.UserName, format), args...))
}

func (self *ICloudRegister) request(session *http.Session, method, url interface{}, params ...Dict) (resp *http.Response) {
    self.debug("request: %v", url)

    for {
        exp := Try(func() { resp = session.Request(method, url, params...) })

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

    return
}

func (self *ICloudRegister) get(url interface{}, params ...Dict) *http.Response {
    return self.request(self.session, "GET", url, params...)
}

func (self *ICloudRegister) post(url interface{}, params ...Dict) *http.Response {
    return self.request(self.session, "POST", url, params...)
}

func (self *ICloudRegister) Initialize() {
}

func (self *ICloudRegister) Signup() int {
    start := time.Now()

    result := SIGNUP_SUCCESS
    reason := ""

    exp := Try(func() {
        self.getQualifySession()
    })

    Println(exp)

    defer time.Sleep(time.Hour)

    self.debug(exp)

    elapsed := time.Now().Sub(start)
    self.info("signup_result:%s proxy:%s reason:%s elapsed:%v", signupResultText[result], self.proxy, reason, elapsed)

    return result
}

func (self *ICloudRegister) getQualifySession() []byte {
    session := http.NewSession()
    defer session.Close()

    self.setProxy(session)
    session.SetHeaders(Dict{
        "User-Agent"            : "%E8%AE%BE%E7%BD%AE/1.0 CFNetwork/711.4.6 Darwin/14.0.0",
        "Accept-Encoding"       : "gzip, deflate",
        "Accept-Language"       : "zh-cn",
        "Accept"                : "*/*",
        "Connection"            : "keep-alive",
    })

    configurations := Dict{}
    self.request(session, "GET",
        "https://setup.icloud.com/configurations/init?context=buddy",
        Dict{
            "headers": Dict{
                "X-MMe-Client-Info" : "<iPhone5,4> <iPhone OS;8.4;12H143> <com.apple.AppleAccount/1.0 (com.apple.Preferences/1.0)>",
                "X-MMe-Country"     : "CN",
            },
        },
    ).Plist(&configurations)

    cert := self.request(session, "GET", configurations["urls"].(map[string]interface{})["qualifyCert"].(string)).Content

    cert = self.sapSession.ExchangeData(ituneslib.SAP_TYPE_LOGIN, cert)

    ret := map[string]interface{}{}

    self.get(
        "https://setup.icloud.com/setup/qualify/session",
        Dict{
            "headers": Dict{
                "X-MMe-Client-Info" : "<iPhone5,4> <iPhone OS;8.4;12H143> <com.apple.AppleAccount/1.0 (com.apple.Preferences/1.0)>",
                "X-MMe-Country"     : "CN",
                "X-MMe-Nas-Session" : base64.StdEncoding.EncodeToString(cert),
            },
        },
    ).Json(&ret)

    success := ret["success"].(bool)
    if success == false {
        Raise("get session info failed")
    }

    Println("session-info =", ret["session-info"].(string))

    sessionInfo, err := base64.StdEncoding.DecodeString(ret["session-info"].(string))
    RaiseIf(err)

    Println(hex.EncodeToString(sessionInfo))
    self.sapSession.ExchangeData(ituneslib.SAP_TYPE_LOGIN, sessionInfo)

    logger.SetLevel(9999)

    return sessionInfo
}

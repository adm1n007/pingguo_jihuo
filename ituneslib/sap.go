package ituneslib

import (
    . "ml/trace"
    . "ml/dict"
    "unsafe"
    "plistlib"
    "ml/net/http"
)

const (
    OSXUserAgent12_3         = "iTunes/12.3 (Macintosh; OS X 10.10.5) AppleWebKit/600.8.9"

    DefaultOSXUserAgent      = "iTunes/12.4.1 (Macintosh; OS X 10.11.5) AppleWebKit/601.6.17"
    DefaultWindowsUserAgent  = "iTunes/12.3 (Windows; Microsoft Windows 8.1 x64 Business Edition (Build 9200); x64) AppleWebKit/7601.1056.1.1"
    DefaultAppStoreUserAgent = "MacAppStore/2.0 (Macintosh; OS X 10.10.2; 14C109) AppleWebKit/0600.3.18"
    // DefaultWindowsUserAgent = "iTunes/12.4 (Windows; Microsoft Windows 10.0 x64 Business Edition (Build 9200); x64) AppleWebKit/7601.6016.1000.1"
)

type SapSession struct {
    session         uintptr
    primeSignature  []byte
    deviceId        *FairPlayHWInfo
    HttpSession     http.HttpSesstion
    UrlBag          Dict
    Initialized     bool
}

const useSapPool = false
var sapSessionPool = make(chan *SapSession, 10000)

func sapInitialize() {
    if useSapPool == false {
        return
    }

    for i := cap(sapSessionPool); i != 0; i-- {
        sapSessionPool <- createSapSession(nil)
    }
}

func NewSapSession() (session *SapSession) {
    if useSapPool {
        return <- sapSessionPool
    }

    return createSapSession(nil)
}

func NewSapSessionWithDeviceId(deviceId *FairPlayHWInfo) (session *SapSession) {
    return createSapSession(deviceId)
}

func createSapSession(deviceId *FairPlayHWInfo) (session *SapSession) {
    var sapSession uintptr

    if deviceId == nil {
        deviceId = NewRandomFairPlayHWInfo()
    }

    st, _, _ := itunes.SapCreateSession.Call(uintptr(unsafe.Pointer(&sapSession)), uintptr(unsafe.Pointer(deviceId)))

    if int32(st) != 0 {
        Raise(newiTunesHelperErrorf("SapCreateSession failed: %X", uint32(st)))
    }

    h := http.NewSession()

    session = &SapSession{
                    session         : sapSession,
                    primeSignature  : []byte{},
                    deviceId        : deviceId,
                    HttpSession     : h,
                    UrlBag          : Dict{},
                }
    return
}

func (self *SapSession) Close() {
    self.HttpSession.Close()

    if self.session == 0 {
        return
    }

    itunes.SapCloseSession.Call(self.session)
    self.session = 0

    if useSapPool {
        sapSessionPool <- createSapSession(nil)
    }
}

func (self *SapSession) Initialize(userAgent string, country CountryID, sapType SapCertType) {
    self.HttpSession.SetHeaders(Dict{
        "User-Agent"            : userAgent,
        "Accept-Encoding"       : "gzip",
        "Accept-Language"       : "zh-cn, zh;q=0.75, en-us;q=0.50, en;q=0.25",
        "X-Apple-Store-Front"   : country.StoreFront(),
        "X-Apple-Tz"            : country.TimeZone(),
    })

    self.HttpSession.GetDefaultOptions().AutoRetry = true
    self.HttpSession.GetDefaultOptions().Ignore404 = false

    self.initUrlbag()
    self.establishContext(sapType)

    self.Initialized = true
}

var cachedUrlBag Dict

func (self *SapSession) initUrlbag() {
    if len(cachedUrlBag) != 0 {
        self.UrlBag = cachedUrlBag
        return
    }

    resp := self.HttpSession.Get("https://init.itunes.apple.com/bag.xml?ix=5&ign-bsn=1")

    plist := Dict{}
    resp.Plist(&plist)
    plistlib.Unmarshal(plist["bag"].([]byte), &self.UrlBag)

    cachedUrlBag = self.UrlBag
}

func (self *SapSession) establishContext(sapType SapCertType) {
    certificateData := self.loadCertificateData()
    data := self.ExchangeData(sapType, certificateData)
    data = self.postExchangeData(data)
    self.ExchangeData(sapType, data)
}

var cachedSignSapSetupCert []byte

func (self *SapSession) loadCertificateData() []byte {
    if cachedSignSapSetupCert == nil {
        signSapSetupCert := Dict{}
        self.HttpSession.Get(self.UrlBag["sign-sap-setup-cert"]).Plist(&signSapSetupCert)
        cachedSignSapSetupCert = signSapSetupCert["sign-sap-setup-cert"].([]byte)
    }

    return cachedSignSapSetupCert
}

func (self *SapSession) postExchangeData(data []byte) []byte {
    body, err := plistlib.MarshalIndent(
                    Dict{
                        "sign-sap-setup-buffer": data,
                    },
                    plistlib.XMLFormat,
                    "    ",
                )

    RaiseIf(err)

    signSapSetupBuffer := Dict{}
    self.HttpSession.Post(
        self.UrlBag["sign-sap-setup"],
        Dict{
            "headers": Dict{
                "Content-Type": "application/x-apple-plist",
            },

            "body": body,
        },
    ).Plist(&signSapSetupBuffer)

    return signSapSetupBuffer["sign-sap-setup-buffer"].([]byte)
}

func (self *SapSession) CreatePrimeSignature() []byte {
    if self.Initialized == false {
        return []byte{}
    }

    if len(self.primeSignature) == 0 {
        var buf *byte
        var size int

        st, _, _ := itunes.SapCreatePrimeSignature.Call(
                            self.session,
                            uintptr(unsafe.Pointer(&buf)),
                            uintptr(unsafe.Pointer(&size)),
                        )
        status := getStatus(st)
        if status != 0 {
            return self.primeSignature
        }

        defer FreeSessionData(buf)
        self.primeSignature = toBytes(buf, size)
    }

    return self.primeSignature
}

func (self *SapSession) ExchangeData(sapType SapCertType, data []byte) (cert []byte) {
    var buf *byte
    var size int

    st, _, _ := itunes.SapExchangeData.Call(
                        self.session,
                        uintptr(sapType),
                        uintptr(unsafe.Pointer(self.deviceId)),
                        bytesPtr(data),
                        bytesLen(data),
                        uintptr(unsafe.Pointer(&buf)),
                        uintptr(unsafe.Pointer(&size)),
                    )

    status := getStatus(st)
    if status != 0 {
        Raise(newiTunesHelperErrorf("SapExchangeData return %X", uint32(status)))
    }

    if buf == nil {
        return
    }

    defer FreeSessionData(buf)
    cert = toBytes(buf, size)

    return
}

func (self *SapSession) VerifyPrimeSignature(signature []byte) {
    if self.Initialized == false {
        return
    }

    st, _, _ := itunes.SapVerifyPrimeSignature.Call(
                        self.session,
                        uintptr(unsafe.Pointer(&signature[0])),
                        uintptr(len(signature)),
                    )

    status := getStatus(st)
    if status != 0 {
        Raise(newiTunesHelperErrorf("SapVerifyPrimeSignature return %X", uint32(status)))
    }
}

func (self *SapSession) SignData(data []byte) (signature []byte) {
    var buf *byte
    var size int

    if self.Initialized == false {
        return []byte{}
    }

    st, _, _ := itunes.SapSignData.Call(
                        self.session,
                        uintptr(unsafe.Pointer(&data[0])),
                        uintptr(len(data)),
                        uintptr(unsafe.Pointer(&buf)),
                        uintptr(unsafe.Pointer(&size)),
                    )

    status := getStatus(st)
    if status != 0 {
        Raise(newiTunesHelperErrorf("SapSignData return %X", uint32(status)))
    }

    defer FreeSessionData(buf)
    signature = toBytes(buf, size)

    return
}

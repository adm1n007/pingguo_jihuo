package http

import (
    . "ml/strings"
    . "ml/dict"
    . "ml/trace"
    . "ml/array"

    urllib "net/url"
    httplib "net/http"
    netlib "net"

    "net/http/cookiejar"
    "fmt"
    "bytes"
    "io"
    "time"
    "io/ioutil"
    "ml/logging/logger"
    "ml/net/socket"
)

var cancelFollowRedirects = fmt.Errorf("")

const (
    NoneProxy   = iota
    Socks5Proxy
    HttpProxy
)

type HttpSesstion interface {
    Close()
    GetDefaultOptions() *RequestOptions
    Request(method, url interface{}, params ...Dict) (resp *Response)
    Get(url interface{}, params ...Dict) (resp *Response)
    Post(url interface{}, params ...Dict) (resp *Response)
    ClearHeaders()
    SetCookies(url String, cookies Dict)
    Headers() httplib.Header
    SetHeaders(headers Dict)
    AddHeaders(headers Dict)
    ProxyType() int
    SetSocks5Proxy(host String, port int, auth *socket.Auth)
    SetProxy(host String, port int, userAndPassword ...String) (err error)
    SetTimeout(timeout time.Duration)
}

type Session struct {
    cookie             *cookiejar.Jar
    client             *httplib.Client
    headers             httplib.Header
    defaultTransport   *Transport
    DefaultOptions     *RequestOptions

    proxyType           int
}

func getDefaultDialer() *netlib.Dialer {
    return &netlib.Dialer{
                Timeout     : 30 * time.Second,
                KeepAlive   : 30 * time.Second,
            }
}

func NewSession() *Session {
    jar, err := cookiejar.New(nil)
    if err != nil {
        RaiseHttpError(err)
    }

    defaultTransport := newTransport(&httplib.Transport{
        Proxy               : nil,
        DisableKeepAlives   : false,
        Dial                : getDefaultDialer().Dial,

        TLSHandshakeTimeout : 10 * time.Second,
    })

    client := &httplib.Client{
        CheckRedirect   : nil,
        Jar             : jar,
        Timeout         : 30 * time.Second,
        Transport       : defaultTransport,
    }

    return &Session{
                cookie              : jar,
                client              : client,
                headers             : make(httplib.Header),
                defaultTransport    : defaultTransport,
                DefaultOptions      : &RequestOptions{
                                            MaxTimeoutTimes : DefaultMaxTimeoutTimes,
                                            Ignore404       : true,
                                            AutoRetry       : true,
                                        },
            }
}

func (self *Session) Close() {
    self.defaultTransport.CloseIdleConnections()
}

func (self *Session) GetDefaultOptions() *RequestOptions {
    return self.DefaultOptions
}

func toString(value interface{}) String {
    switch v := value.(type) {
        case string:
            return String(v)

        case String:
            return v

        case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
            return String(fmt.Sprintf("%v", v))

        default:
            fmt.Printf("unknown value type %v\n", value)
            return String(v.(string))
    }
}

func (self *Session) getRequestOptions(params ...Dict) *RequestOptions {
    switch len(params) {
        case 0:
            return self.DefaultOptions

        case 1:
            switch opt := params[0]["options"].(type) {
                case RequestOptions:
                    return &opt

                case *RequestOptions:
                    return opt
            }
    }

    return self.DefaultOptions
}

func dictToValues(d Dict, encoding Encoding) urllib.Values {
    values := urllib.Values{}
    for k, v := range d {

        key := toString(k)
        value := toString(v)

        values.Set(string(key.Encode(encoding)), string(value.Encode(encoding)))
    }

    return values
}

func orderedDictToBody(d OrderedDict, encoding Encoding) []byte {
    keys := d.Keys()
    values := make([]String, keys.Length())

    for index, key := range keys {
        k := toString(key)
        v := toString(d.Get(key))
        values[index] = String(fmt.Sprintf("%s=%s", urllib.QueryEscape(k.String()), urllib.QueryEscape(v.String())))
    }

    return String("&").Join(values).Encode(encoding)
}

func applyHeadersToRequest(request *httplib.Request, defaultHeaders httplib.Header, extraHeaders Dict) {
    for k, vs := range defaultHeaders {
        for _, v := range vs {
            request.Header.Add(k, v)
        }
    }

    for k, v := range extraHeaders {
        request.Header.Set(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v))
    }
}

func (self *Session) requestImpl(methodi, urli interface{}, params_ ...Dict) (*Response) {
    var bodyReader      io.Reader
    var bodyData        []byte
    var params          Dict
    var encoding        Encoding
    var err             error

    method := toString(methodi)
    url := toString(urli)
    requestParams := urllib.Values{}

    options := self.getRequestOptions(params_...)

    {
        u, err := urllib.Parse(url.String())
        RaiseHttpError(err)
        query, err := urllib.ParseQuery(u.RawQuery)
        RaiseHttpError(err)

        for k, vs := range query {
            for _, v := range vs {
                requestParams.Add(k, v)
            }
        }
    }

    switch (len(params_)) {
        case 1:
            params = params_[0]

        case 0:
            params = Dict{}

        default:
            Raise(NewHttpError(HTTP_ERROR_GENERIC, method, url, String(fmt.Sprintf("invalid params: %d", len(params_)))))
    }

    switch v := params["encoding"].(type) {
        case int, Encoding:
            encoding = v.(Encoding)

        default:
            encoding = CP_UTF8
    }

    switch body := params["body"].(type) {
        case string:
            b := String(body)
            bodyData = b.Encode(encoding)

        case String:
            bodyData = body.Encode(encoding)

        case []byte:
            bodyData = body

        case Dict:
            bodyData = String(dictToValues(body, encoding).Encode()).Encode(encoding)

        case OrderedDict:
            bodyData = orderedDictToBody(body, encoding)

        default:
            bodyReader = nil
    }

    if bodyData != nil {
        bodyReader = bytes.NewBuffer(bodyData)
    }

    request, err := httplib.NewRequest(method.String(), url.String(), bodyReader)
    RaiseHttpError(err)

    switch query := params["params"].(type) {
        case Dict:
            for k, v := range query {
                requestParams.Add(string(toString(k).Encode(encoding)), string(toString(v).Encode(encoding)))
            }
    }

    extraHeaders := Dict{}

    switch headers := params["headers"].(type) {
        case Dict:
            extraHeaders = headers
    }

    switch options.OverwriteHeaders {
        case true:
            applyHeadersToRequest(request, nil, extraHeaders)

        case false:
            applyHeadersToRequest(request, self.headers, extraHeaders)
    }

    if len(requestParams) != 0 {
        queryString := ""

        ignoreEncodeKeys := options.IgnoreEncodeKeys
        if ignoreEncodeKeys == nil {
            ignoreEncodeKeys = Array{}
        }

        for k, vs := range requestParams {
            for _, v := range vs {
                if len(queryString) != 0 {
                    queryString += "&"
                }

                if ignoreEncodeKeys.Contains(k) == false {
                    k = urllib.QueryEscape(k)
                    v = urllib.QueryEscape(v)
                }

                queryString += fmt.Sprintf("%s=%s", k, v)
            }
        }

        request.URL.RawQuery = queryString
    }

    self.client.CheckRedirect = func(request *httplib.Request, via []*httplib.Request) error {
        if options.DontFollowRedirects {
            return cancelFollowRedirects
        }

        if len(via) >= 10 {
            Raise(NewHttpError(HTTP_ERROR_TOO_MANY_REDIRECT, method, url, "stopped after 10 redirects"))
        }

        request.Method = method.String()

        if bodyData != nil {
            request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyData))
            request.ContentLength = int64(len(bodyData))
        }

        applyHeadersToRequest(request, self.headers, extraHeaders)
        return nil
    }


    request.Close = false
    resp, err := self.client.Do(request)

    timeout := self.defaultTransport.RemoveCancelledRequest(request)

    self.client.CheckRedirect = nil

    if resp != nil && resp.Body != nil {
        defer func() {
            if resp.Body != nil {
                resp.Body.Close()
            }
        }()
    }

    switch {
        case err != nil:
            if timeout == false {
                self.defaultTransport.CancelRequest(request)
                self.defaultTransport.RemoveCancelledRequest(request)
            }

            uerr := err.(*urllib.Error)
            if uerr.Err == cancelFollowRedirects {
                break
            }

            herr := &HttpError{
                        Op      : uerr.Op,
                        URL     : uerr.URL,
                        Err     : uerr.Err,
                    }

            msg := String(herr.Err.Error())

            switch {
                case timeout,
                     msg.Contains("TLS handshake timeout"),
                     msg.Contains("Client.Timeout exceeded"):
                    herr.Type = HTTP_ERROR_TIMEOUT

                case msg.Contains("error connecting to proxy"):
                    herr.Type = HTTP_ERROR_CONNECT_PROXY

                case msg.Contains("unexpected EOF"):
                    herr.Type = HTTP_ERROR_INVALID_RESPONSE

                case msg.Contains("wsarecv"):
                    herr.Type = HTTP_ERROR_READ_ERROR

                case msg.Contains("Bad Gateway"):
                    herr.Type = HTTP_ERROR_BAD_GATE_WAY

                case msg.Contains("connectex:"):
                    herr.Type = HTTP_ERROR_CANNOT_CONNECT

                default:
                    herr.Type = HTTP_ERROR_GENERIC
            }

            Raise(herr)

            return nil
    }

    switch HttpStatusCode(resp.StatusCode) {
        case StatusMovedPermanently, StatusFound, StatusSeeOther, StatusTemporaryRedirect:
            resp.Body.Close()
            resp.Body = nil
    }

    // fix net/http/client/shouldRedirectPost does not follow StatusTemporaryRedirect

    switch HttpStatusCode(resp.StatusCode) {
        case StatusTemporaryRedirect:
            if options.DontFollowRedirects {
                break
            }

            if location := resp.Header.Get("Location"); location != "" {
                if resp != nil && resp.Body != nil {
                    resp.Body.Close()
                }

                return self.requestImpl(method, location, params_...)
            }
    }

    r := NewResponse(resp, options)

    return r
}

func (self *Session) Request(method, url interface{}, params ...Dict) (resp *Response) {
    options := self.getRequestOptions(params...)

    if options.AutoRetry == false {
        return self.requestImpl(method, url, params...)
    }

    maxTimeoutTimes := options.MaxTimeoutTimes
    timeoutTimes := 0

    for {
        exp := Try(func() { resp = self.requestImpl(method, url, params...) })

        if exp != nil {
            e := exp.Value.(*HttpError)

            switch e.Type {
                case HTTP_ERROR_TIMEOUT:
                    timeoutTimes += 1
                    if timeoutTimes > maxTimeoutTimes {
                        e.Type = HTTP_ERROR_CONNECT_PROXY
                        Raise(exp)
                    }
                    fallthrough

                case HTTP_ERROR_INVALID_RESPONSE,
                     HTTP_ERROR_BAD_GATE_WAY:
                     // HTTP_ERROR_GENERIC:
                    time.Sleep(time.Second)
                    continue

                case HTTP_ERROR_CONNECT_PROXY:
                    fallthrough
                default:
                    Raise(exp)
            }
        }

        switch resp.StatusCode {
            case StatusOK,
                 StatusCreated,
                 StatusNoContent,
                 StatusFound,
                 StatusPreconditionFailed,
                 StatusConferenceNotFound,
                 StatusNotModified:
                break

            case StatusNotFound:
                if options.Ignore404 {
                    break
                }

                fallthrough

            case StatusBadGateway,
                 StatusServiceUnavailable,
                 StatusGatewayTimeout:
                time.Sleep(time.Second)
                continue

            default:
                logger.Debug("unknown StatusCode: %d %v", resp.StatusCode, resp.StatusCode)
        }

        break
    }

    // time.Sleep(time.Second * 15)
    return
}

func (self *Session) Get(url interface{}, params ...Dict) (resp *Response) {
    return self.Request(MethodGet, url, params...)
}

func (self *Session) Post(url interface{}, params ...Dict) (resp *Response) {
    return self.Request(MethodPost, url, params...)
}

func (self *Session) ClearHeaders() {
    self.headers = httplib.Header{}
}

func (self *Session) SetCookies(url String, cookies Dict) {
    u, err := urllib.Parse(url.String())
    RaiseHttpError(err)

    c := []*httplib.Cookie{}

    for k, v := range cookies {
        c = append(c, &httplib.Cookie{
                Name    : fmt.Sprintf("%v", k),
                Value   : fmt.Sprintf("%v", v),
                Domain  : u.Host,
            },
        )
    }

    self.cookie.SetCookies(u, c)
}

func (self *Session) Headers() httplib.Header {
    return self.headers
}

func (self *Session) SetHeaders(headers Dict) {
    if len(headers) == 0 {
        self.ClearHeaders()
        return
    }

    for k, v := range headers {
        self.headers.Set(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v))
    }
}

func (self *Session) AddHeaders(headers Dict) {
    if len(headers) == 0 {
        self.ClearHeaders()
        return
    }

    for k, v := range headers {
        self.headers.Add(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v))
    }
}

func (self *Session) ProxyType() int {
    return self.proxyType
}

func (self *Session) SetSocks5Proxy(host String, port int, auth *socket.Auth) {
    if host.IsEmpty() {
        self.defaultTransport.Dial = getDefaultDialer().Dial
        self.proxyType = NoneProxy
        return
    }

    self.defaultTransport.Dial = func (network, address string) (c netlib.Conn, err error) {
        socks5 := socket.NewSocks5Dialer("tcp", host.String(), port, auth)
        a := String(address).Split(":", 1)
        return socks5.Connect(network, a[0].String(), a[1].ToInt(), 30 * time.Second)
    }

    self.SetProxy("", 0)
    self.proxyType = Socks5Proxy
}

func (self *Session) SetProxy(host String, port int, userAndPassword ...String) (err error) {
    if host.IsEmpty() {
        self.defaultTransport.Proxy = nil
        self.proxyType = NoneProxy

    } else {
        var proxyUrl *urllib.URL
        var user, pass String

        switch len(userAndPassword) {
            case 2:
                user, pass = userAndPassword[0], userAndPassword[1]
                if user.IsEmpty() == false && pass.IsEmpty() == false {
                    proxyUrl, err = urllib.Parse(fmt.Sprintf("http://%s:%s@%s:%d", user, pass, host, port))
                    break
                }

                fallthrough

            case 1:
                user = userAndPassword[0]
                if user.IsEmpty() == false {
                    proxyUrl, err = urllib.Parse(fmt.Sprintf("http://%s@%s:%d", user, host, port))
                    break
                }

                fallthrough

            default:
                proxyUrl, err = urllib.Parse(fmt.Sprintf("http://%s:%d", host, port))
        }

        if err != nil {
            return
        }

        self.defaultTransport.Proxy = httplib.ProxyURL(proxyUrl)
        self.SetSocks5Proxy("", 0, nil)
        self.proxyType = HttpProxy
    }

    return
}

func (self *Session) SetTimeout(timeout time.Duration) {
    self.client.Timeout = timeout
}

package mailmaster

import (
    "active_apple/packet"
    . "fmt"
    . "active_apple/ml/strings"
    . "active_apple/ml/dict"

    "time"

    "crypto/md5"

    "active_apple/ml/net/socket"
    "active_apple/ml/random"
    "active_apple/ml/uuid"
    "active_apple/ml/crypto/rsa"
    "active_apple/ml/encoding/binary"

	"log"
)

type MailMaster struct {
    userName    String
    password    String
    sock        socket.Socket
    serverId    int
    serialId    int
    timeout     time.Duration
    cookies     []String
    sid         String
    host        String
    aeskey      []byte

    proxyHost   String
    proxyPort   int
    proxyAuth   *socket.Auth
}

func NewMailMaster(user, password String) *MailMaster {
    mm := &MailMaster{
        userName    : user,
        password    : String(Sprintf("%x", md5.Sum(password.Encode(CP_UTF8)))),
        sock        : socket.NewTcpSocket(),
        serverId    : 4,
        serialId    : 1,
        timeout     : time.Second * 8,
        aeskey      : []byte("3090111110251509"),
    }

    return mm
}

func (self *MailMaster) Close() {
    self.sock.Close()
}

func (self *MailMaster) SetProxy(host String, port int, auth *socket.Auth) {
    self.proxyHost = host
    self.proxyPort = port
    self.proxyAuth = auth

    self.setProxyForSocket(self.sock)
}

func (self *MailMaster) getServer() (host String, port int) {
    //获取163邮件socket服务器ip和端口,返回格式类似这样 %!(EXTRA []strings.String=[220.181.13.217:443 220.181.13.218:443 220.181.13.219:443 220.181.13.220:443 220.181.13.229:443 220.181.13.232:443] )
    svrList := self.getServerListForUser(self.userName)

    //随机获取一个
    address := random.Choice(svrList).(String).Split(":")
    host, port = address[0], address[1].ToInt()

    return
}

func (self *MailMaster) setProxyForSocket(sock socket.Socket) {
    if self.proxyHost.IsEmpty() == false {
        sock.SetSocks5Proxy(self.proxyHost.String(), self.proxyPort, self.proxyAuth)
    }
}

func (self *MailMaster) connect(host String, port int) {
    self.setProxyForSocket(self.sock)

    self.sock.SetTimeout(self.timeout)
    self.sock.Connect(host.String(), int(port), self.timeout)
    self.handShake()

    return
}

func (self *MailMaster) Login() {
    //获取163可用的socket服务器
    host, port := self.getServer()

    //与163建立socket连接,timeout为8s
    self.connect(host, port)

    //尝试登录163邮箱
    self.login(host, port)
}

func (self *MailMaster) RegisterEntrance() {
/*
body = "cmd=register.entrance&flow=d_mail";

headers = {
    '256': b'/regall/unireg/call.do',
     '257': b'https://ssl.mail.163.com',
     '258': {'cmd': 'register.entrance', 'flow': 'd_mail'},
     '259': {
                'Accept-Language': 'zh-Hans;q=1',
                'User-Agent': 'mail/4.7.1 (iPod touch; iOS 8.3; Scale/2.00)'
            }
 }
*/
}

func (self *MailMaster) CheckName(userName String) {
/*
{
 '256': '/regall/unireg/call.do',
 '257': 'https://ssl.mail.163.com',
 '258': {
            'cmd': 'urs.checkName',
             'name': 'qwergdabnml',
             'product': 'mail_master',
             's_rev': '1',
             'sig': 'u6kTcyvlSFdYapHmjph9rcFkcIa-6UGIhcZEZw1bah8',
             'ts': '1457599292'
         },
 '259':{
            'Accept-Language': 'zh-Hans;q=1',
            'User-Agent': 'mail/4.7.1 (iPod touch; iOS 8.3; Scale/2.00)'
        }
}

    body = '&'.join([
        'cmd=urs.checkName',
        'name=qwergdabnml',
        'product=mail_master',
        's_rev=1',
        'ts=1457599292',
        'sig=u6kTcyvlSFdYapHmjph9rcFkcIa-6UGIhcZEZw1bah8',
    ])

*/
}

func (self *MailMaster) VerifyCode() {
/*
{
 '256': b'/regall/unireg/call.do',
 '257': b'https://ssl.mail.163.com',
 '258': {
            'cmd': 'register.verifyCode',
            'env': '625738110216',
            't': 1457601067616,
            'vt': 'd_mail'
        },
 '259': {
            'Accept-Language': 'zh-Hans;q=1',
            'Cookie': 'mailsync=f609bd1fcdee90b086950dabb0b0f99359b12750a7dd72d68f8227f6ec618b65f14cc9b17aa2b51aa8dc42bce97830b5,JSESSIONID=3072A9B27C566CD6A876CDBDF6804161,ser_adapter=INTERNAL134',
            'User-Agent': 'mail/4.7.1 (iPod touch; iOS 8.3; Scale/2.00)'
        },
}

    body = '&'.join([
                'cmd=register.verifyCode',
                'vt=d_mail',
                'env=625738110216',
                't=1457601067616',
            ])
*/
}

func (self *MailMaster) ListMessages(sentDate String) (mailIDs []String) {
    resp := self.sendMessage("wzp:listMessages", nil, Dict{
                "sentDate"  : sentDate,
                "desc"      : true,
                "fids"      : []int{1},
                "order"     : "date",
                "limit"     : 50,
                "start"     : 0,
            })

    for _, mid := range resp.Items() {
        mailIDs = append(mailIDs, String(mid.Value.(string)))
    }

    return
}

func (self *MailMaster) GetMessageInfos(mailIDs []String) (mails JsonArray) {
    resp := self.sendMessage("wzp:getMessageInfos", nil, Dict{
                "ids"               : mailIDs,
                "returnTag"         : true,
                "returnThreadsInfo" : true,
            })

    for _, mail := range resp.Values() {
        mails.Append(JsonDict(mail.(map[string]interface{})))
    }

    return
}

func (self *MailMaster) AsyncReadMessage(mailId, mode String, autoName bool) (message JsonDict) {
    resp := self.sendMessage("wzpm:asyncReadMessage", nil, Dict{
                "id"        : mailId,
                "mode"      : mode,
                "autoName"  : autoName,
            })

    message = JsonDict{}
    for _, item := range resp.Items() {
        message[item.Key.(string)] = item.Value
        // message.Set(item.Key, item.Value)
    }

    return
}

func (self *MailMaster) getServerListForUser(user String) []String {
    pkt := packet.NewPacket(1)
    pkt.SerialId = 1

    pkt.AddBody("srvids", []int{self.serverId})
    pkt.AddBody("appid", 2)
    pkt.AddBody("user", user)

    sock := socket.NewTcpSocket()
    defer sock.Close()

    self.setProxyForSocket(sock)
    sock.SetTimeout(self.timeout)
    sock.Connect(locateDomain["locate_domain"].(string), locateDomain["locate_port"].(int), self.timeout)
    sock.Write(pkt.ToBinary())

    //获取sock的数据
    pkt = packet.Recieve(sock, nil)

    //可能是邮件列表
    list := pkt.Body.Get("srvs").(map[string]interface{})[Sprintf("%d", self.serverId)].([]interface{})

    svrList := []String{}

    //循环邮件列表，放到一个svrList切片并返回
    for _, svr := range list {
        svrList = append(svrList, String(svr.(string)))
    }

    return svrList
}

func (self *MailMaster) login(host String, port int) {
    pkt := packet.NewPacket(self.serverId)

    pkt.SerialId = 2

    //注释接手前的随机uuid,更换为以邮箱名字作为参数的uuid
    //deviceId, _ := uuid.NewV4()
    deviceId, _ := uuid.NewV3(&uuid.UUID{}, []byte(self.userName))

    pkt.AddBody("device-id",    String(deviceId.String()).ToUpper())
    pkt.AddBody("type",         "1")
    pkt.AddBody("product",      "mail_client")
    pkt.AddBody("password",     self.password)
    pkt.AddBody("passtype",     "0")
    pkt.AddBody("ignore",       "false")
    pkt.AddBody("app-version",  "4.1.1.508")
    pkt.AddBody("funcid",       "loginone")
    pkt.AddBody("df",           "mailmaster_iphone")
    pkt.AddBody("username",     self.userName)

    sock := socket.NewTcpSocket()

    self.setProxyForSocket(sock)
    sock.Connect(host.String(), port, self.timeout)
    sock.Write(pkt.ToBinary())

    pkt = packet.Recieve(sock, nil)


    /*
        {
          "cookie": [Coremail:1456641018166%uGpLpgvHZnGwBLafryHHWdcDPkVAbMAW%g1a162.mail.163.com MAIL_SESS:gbI60JERnG2Qqs3q_L78EoKVxVwTv2zD7pNH8Oi9abpllwr82lQrsyxvDFVikoOxiROCsunEnj9b9k6hMQzqsNGiD7tOwao8vzN9d7_C3hqPXuBCdlHhe0YTErCfRc2BJnRDVXHJfhUOT0Fx64F6UN.N5faOAXmUvWYmyeRLqoHb8 NTES_SESS:gbI60JERnG2Qqs3q_L78EoKVxVwTv2zD7pNH8Oi9abpllwr82lQrsyxvDFVikoOxiROCsunEnj9b9k6hMQzqsNGiD7tOwao8vzN9d7_C3hqPXuBCdlHhe0YTErCfRc2BJnRDVXHJfhUOT0Fx64F6UN.N5faOAXmUvWYmyeRLqoHb8 P_INFO:jiuhao7221409@163.com],
          "host": "192.168.193.162",
          "result": true,
          "sid": "uGpLpgvHZnGwBLafryHHWdcDPkVAbMAW",
        }

        if pkt.Body.Get("result").(bool) == false {
            pkt.Body.Get("cause")
        }

    */



    if pkt.Body.Get("result").(bool) == false {
        switch pkt.Body.Get("cause").(string) {
            case "NULL":
                //Raise(NewLoginError("body = %v", pkt.Body))
                log.Panic("body = %v", pkt.Body)

            case "LOGINONE_URS_AUTHFAIL":
                //Raise(NewPasswordError(""))
                log.Panic("email password error")

            case "LOGIN_RETURN":
                log.Panic("body = %v", pkt.Body)
        }
    }

    self.host   = String(pkt.Body.Get("host").(string))
    self.sid    = String(pkt.Body.Get("sid").(string))
    self.cookies = []String{}

    for _, cookie := range pkt.Body.Get("cookie").([]interface{}) {
        self.cookies = append(self.cookies, String(cookie.(string)))
    }
}

func (self *MailMaster) handShake() {
    cipher := rsa.NewRSACipher(locateDomain["pubkey"])
    pk := cipher.PublicKey()
    pkt := packet.NewPacket(0x80)
    pubkey := md5.Sum([]byte(Sprintf("%X,%06X", pk.N, pk.E)))

    pkt.AddHeader(packet.PubkeyVersion, binary.IntToBytes(locateDomain["pubkey_version"].(int), 2, binary.BigEndian))
    pkt.AddHeader(packet.Pubkey, pubkey[:])
    pkt.SetBodyData(cipher.EncryptPKCS1v15(self.aeskey))

    self.sendPacket(pkt)
}

func (self *MailMaster) sendPacket(data *packet.Packet) *packet.Packet {
    self.serialId++

    data.SerialId = self.serialId
    self.sock.Write(data.ToBinary())

    return packet.Recieve(self.sock, data.Aeskey)
}

func (self *MailMaster) sendMessage(api String, headers Dict, body Dict) OrderedDict {
    message := packet.NewPacket(0x80)
    message.Aeskey = self.aeskey

    switch {
        case headers == nil:
            message.AddHeader(packet.Host,      []byte(self.host))
            message.AddHeader(packet.Sid,       []byte(self.sid))
            message.AddHeader(packet.Params,    []byte(self.cookies[0]))
            message.AddHeader(packet.Api,       []byte(api))

        case headers != nil:
            for k, v := range headers {
                message.AddHeader(k.(packet.Headers), v.([]byte))
            }
    }

    if body != nil {
        for k, v := range body {
            message.AddBody(String(Sprintf("%v", k)), v)
        }
    }

    return self.sendPacket(message).Body
}

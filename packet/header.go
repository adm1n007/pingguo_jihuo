package packet

import (
    . "fmt"
)


/*

    257
    https://ssl.mail.163.com

    258
    {"cmd":"register.entrance","flow":"d_mail"}

    256
    /regall/unireg/call.do

    259
    {"User-Agent":"mail\/4.7.1 (iPod touch; iOS 8.3; Scale\/2.00)","Accept-Language":"zh-Hans;q=1"}

*/

type Headers int

const (
    PubkeyVersion   = Headers(1)
    Pubkey          = Headers(2)
    Api             = Headers(256)
    Sid             = Headers(257)
    Params          = Headers(258)
    HttpHeaders     = Headers(259)
    Status          = Headers(338)
    Host            = Headers(512)
)

var headerText = map[Headers]string{
    PubkeyVersion   : "PubkeyVersion",
    Pubkey          : "Pubkey",
    Sid             : "Sid",
    Host            : "Host",
    Params          : "Params",
    Api             : "Api",
    Status          : "Status",
}

func (self Headers) String() string {
    s, ok := headerText[self]
    if ok == false {
        return Sprintf("%d", int(self))
    }

    return s
}

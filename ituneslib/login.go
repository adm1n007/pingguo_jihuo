package ituneslib

import (
    . "ml/strings"
    . "ml/dict"

    "plistlib"
    "ml/encoding/base64"
    "ml/net/http"
    // "ml/logging/logger"
)

type LoginReason int

const (
    LoginReason_SignIn              LoginReason = iota
    LoginReason_MachineAuthorize
    LoginReason_Purchase
    LoginReason_MachineDeauthorize
    LoginReason_ServerDialog
    LoginReason_ViewAccount
)

func (self LoginReason) String() string {
    return map[LoginReason]string{
        LoginReason_SignIn              : "signIn",
        LoginReason_MachineAuthorize    : "machineAuthorize",
        LoginReason_Purchase            : "purchase",
        LoginReason_MachineDeauthorize  : "machineDeauthorize",
        LoginReason_ServerDialog        : "serverDialog",
        LoginReason_ViewAccount         : "viewAccount",
    }[self]
}

type Login struct {
    http.HttpSesstion

    Response        JsonDict
    PasswordToken   String
    Dsid            int64
}

func NewLogin(httpSession http.HttpSesstion) *Login {
    return &Login{
        HttpSesstion: httpSession,
    }
}

func (self *Login) Close() {
    self.HttpSesstion.Close()
}

func (self *Login) Login(sap *SapSession, params Dict) bool {
    form := JsonDict{
        "appleId"       : params["userName"],
        "attempt"       : 1,
        "createSession" : "true",
        "guid"          : params["machineGuid"],
        "machineName"   : params["machineName"],
        "password"      : params["password"],
        "why"           : params["why"].(LoginReason).String(),
    }

    body, _ := plistlib.MarshalIndent(form, plistlib.XMLFormat, "    ")

    resp := self.Post(
                sap.UrlBag["authenticateAccount"],
                Dict{
                    "headers": Dict{
                        "X-Apple-ActionSignature"   : base64.EncodeToString(sap.SignData(body)),
                        "Content-Type"              : "application/x-apple-plist",
                    },
                    "body": body,
                },
            )

    // logger.Debug("[%v] login ret:\n%s", params["userName"], resp.Content)

    var p JsonDict
    resp.Plist(&p)

    self.Response = p

    // x, _ := plistlib.MarshalIndent(p, plistlib.XMLFormat, "  ")

    token := p["passwordToken"]
    if token == nil {
        return false
    }

    self.PasswordToken = String(token.(string))
    self.Dsid = String(p["dsPersonId"].(string)).ToInt64()

    self.SetHeaders(Dict{
        "X-Dsid"    : self.Dsid,
        "X-Token"   : token,
    })

    return true
}

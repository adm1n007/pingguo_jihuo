package activator

import (
    . "fmt"
    . "active_apple/ml/strings"
)

func (self *AppleIdActivator) FindVerifyUrls2() (result ActivateResult, urls []String) {
    if String(self.account.UserName).ToLower().EndsWith("@163.com") {
        return self.FindVerifyUrls163()
    }

    link := Sprintf("http://itunes.ccjxqj.com:8000/query?mails=%s", self.account.UserName)

    resp := self.get(link)

    ret := map[string]string{}
    resp.Json(&ret)

    activationLink, exists := ret[self.account.UserName]
    if exists == false {
        result = ACTIVATE_EMAIL_NOT_FOUND
        return
    }

    urls = append(urls, String(activationLink))

    return
}

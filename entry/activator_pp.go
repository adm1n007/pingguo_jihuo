package entry

import (
    . "ml/strings"

    "time"
    "account"
    "proxy"
    "activator"
    "globals"
)

func findVerifyUrlsPP(act *activator.AppleIdActivator) (result activator.ActivateResult, links []String) {
    failureCount := 0

    for ; failureCount < 30 && globals.Exiting == false; failureCount++ {
        //循环30次，直到激活成功 或 激活email找不到,failureCount恢复0
        result, links = act.FindVerifyUrls()
        if result == activator.ACTIVATE_SUCCESS ||
           result == activator.ACTIVATE_EMAIL_NOT_FOUND {

            failureCount = 0
            break
        }

        time.Sleep(time.Second * time.Duration(failureCount))
    }

    if failureCount != 0 {
        // result = activator.ACTIVATE_MAIL_LOGIN_FAILED
        return
    }

    //激活失败，以为着links长度为0，设置result为找不到激活邮件
    if len(links) == 0 {
        result = activator.ACTIVATE_EMAIL_NOT_FOUND
        return
    }

    return
}

func RunActivatorPP() {
    //163获取激活邮件
    callbacks.findVerifyUrls = findVerifyUrlsPP

    //帐号管理模块
    accountManager := account.NewManager(account.UnactivatedManager)

    //代理模块
    proxyManager := proxy.NewManagerPP("appleidactivate")

    //跑程序
    activatorRun(proxyManager, accountManager)
}

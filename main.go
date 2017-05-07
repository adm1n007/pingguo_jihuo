package main

import (
	"active_apple/mailmaster"
	. "active_apple/ml/strings"
	"log"

)

func main() {
	//返回一个带用户名密码,socket连接等指针结构体
	mm := mailmaster.NewMailMaster(String("bcd302388l2@163.com"), String("72575144"))

	defer mm.Close()

	//尝试登录163邮箱
	mm.Login()

	//不知道原因
	mids := mm.ListMessages("00011230080543:")
	//获取邮件信息
	mailInfos := mm.GetMessageInfos(mids)

	//循环邮件信息
	for index := range mailInfos {
	    log.Println("ffffffffffffff")
	    info := mailInfos.Map(index)
	    // msg := mm.AsyncReadMessage(String(info.Get("mid").(string)), "both", true)
	    msg := mm.AsyncReadMessage(String(info["mid"].(string)), "both", true)

	    text := msg.Map("text")
	    if text == nil {
	        continue
	    }

	    content := text["content"]
	    if content == nil {
	        continue
	    }
	}
}

package activator

import (
    . "active_apple/ml/strings"
    . "active_apple/ml/trace"

    "time"
    "active_apple/ml/net/socket"

    "active_apple/mailmaster"
    "active_apple/mailmaster/packet"
    "active_apple/globals"
)

func (self *AppleIdActivator) read163Emails(callback func (content String) bool) {
    var exp *Exception

    for i := 0; globals.Exiting == false && i != 100; {
        exp = Try(func () {
            //返回一个带用户名密码,socket连接等指针结构体
            mm := mailmaster.NewMailMaster(String(self.account.UserName), String(self.account.MailPassword))
            defer mm.Close()

            //是否设置为socket代理登录163邮箱
            if globals.Preferences.UseSocks5Proxy {
                p := self.proxyManager.GetProxy(-1)
                //port := p.Port
                //注释掉原来的port，写定port为socketport。因为163使用socket
                port := 8880
                self.debug("read163Email use socketProxy: %v", p)
                var auth *socket.Auth
                //由于是使用NewProxyWithJobID创建出来的p,所以User属性为空
                if p.User.IsEmpty() == false {
                    auth = &socket.Auth{
                        User: p.User.String(),
                        Password: p.Password.String(),
                    }
                    port = 61080
                }
                //设置代理连接属性
                mm.SetProxy(p.Host, port, auth)
            }

            //尝试登录163邮箱
            mm.Login()
            //不知道原因
            mids := mm.ListMessages("00011230080543:")
            //获取邮件信息
            mailInfos := mm.GetMessageInfos(mids)

            //循环邮件信息
            for index := range mailInfos {
                self.debug("async read")
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

                if readNext := callback(String(content.(string))); readNext == false {
                    break
                }
            }
        })

        if exp == nil {
            break
        }

        switch exp.Value.(type) {
            case *packet.PacketError:
                self.debug("packet corrupt: %v", exp)

            case *TimeoutError:
                self.debug("connect timeout: %v", exp)

            case *mailmaster.LoginError:
                self.debug("163 login error, ip banned: %v", exp)

            case *mailmaster.PasswordError:
                self.debug("163 password error")
                Raise(exp)

            default:
                i++
                self.debug("connect 163 failed: %v", exp)
                time.Sleep(time.Second)
        }

        self.debug("163 error: %v", exp)
    }

    if exp != nil {
        Raise(NewFileNotFoundError("%v", exp))
    }
}

func (self *AppleIdActivator) FindVerifyUrls163() (result ActivateResult, urls []String) {
    //默认为找不到激活邮件
    result = ACTIVATE_EMAIL_NOT_FOUND

    self.read163Emails(func (content String) (readNext bool) {
        for _, line := range content.SplitLines() {
            if line.Contains("https://id.apple.com") == false {
                continue
            }

            urls = append(urls, line)
            result = ACTIVATE_SUCCESS

            return false
        }

        return true
    })

    return
}

func (self *AppleIdActivator) FindVerifyCode163() (result ActivateResult, codes []String) {
    result = ACTIVATE_EMAIL_NOT_FOUND

    self.read163Emails(func (content String) bool {
        lines := content.SplitLines()
        for i, line := range lines {
            if line.Contains("输入下方的验证码：") == false {
                continue
            }

            if i + 3 >= len(lines) {
                break
            }

            for x := i + 1; x <= i + 3; x++ {
                if lines[x].IsEmpty() == false && lines[x].ToInt() != 0 {
                    codes = append(codes, lines[x])
                    result = ACTIVATE_SUCCESS
                    return false
                }
            }
        }

        return true
    })

    return
}

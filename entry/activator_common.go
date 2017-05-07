package entry

import (
	. "fmt"
	. "ml"
	. "ml/strings"
	. "ml/trace"

	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"ml/console"
	"ml/logging/logger"
	"ml/net/http"
	"ml/timer"

	"account"
	"activator"
	"globals"
	"mailmaster"
	"proxy"
)

func verifyAccount(act *activator.AppleIdActivator, proxyManager proxy.Manager) (result activator.ActivateResult) {
	for globals.Exiting == false {

		result = activator.ACTIVATE_NETWORK_ERROR

		// dangerous usage of proxy?
		exp := Try(func() {
			result = act.VerifyPassword(proxyManager)
		})

		if exp == nil {
			break
		}

		if _, ok := exp.Value.(*http.HttpError); ok {
			continue
		}

		logger.Debug("VerifyPassword exception:\n%v", exp)
	}

	return
}

func activate(acc *account.AppleAccount, proxyManager proxy.Manager) (result activator.ActivateResult, exp *Exception) {
    //默认结果为激活失败
	result = activator.ACTIVATE_FAILED

	logger.Debug("[%v] activate", acc.UserName)

	exp = Try(func() {
		act := activator.NewActivator(acc, proxyManager)
		defer act.Close()

        //登录苹果id，确认密码是否正确
		if acc.PasswordVerified == false {
			result = verifyAccount(act, proxyManager)

			switch result {
			case activator.ACTIVATE_PASSWORD_ERROR, //appleId 密码错误
				activator.ACTIVATE_ACCOUNT_BANNED,  //appleId 被封
                activator.ACTIVATE_CANT_ACTIVATE:   //不能重发激活邮件
				return
			}

			acc.PasswordVerified = true
		}

        time.Sleep(time.Second * 15)
        //激活时间从这里开始算,从163注册开始
		activateStartTime := time.Now()
        var links []String

        //重新定义result:激活结果,link:xxx .用于163
		switch {
            case acc.Status != account.ACCOUNT_EMAIL_NOT_FOUND,
                 acc.VerificationEmailResent:
                    //不管163是否已经有激活邮件，都重复一封激活邮件
                    //e := Try(func() {
                    //    result = act.ResendVerificationEmail()
                    //})
                    ////检查163重复邮件是否有http错误
                    //if e != nil {
                    //    if _, isHttpError := e.Value.(*http.HttpError); isHttpError {
                    //        Raise(e)
                    //    }
                    //
                    //    logger.Debug("[%v] ResendVerificationEmail error: %v", acc.UserName, e)
                    //}
                    //
                    ////如果重复邮件结果为延迟处理，则暂时30秒再执行程序
                    //if activator.ACTIVATE_DELAY_PROCESS == result {
                    //    time.Sleep(time.Second * 15)
                    //} else {
                    //    return
                    //}

                    //帐号状态不在6 和 帐号重发邮件 bool，就获取link
                    result, links = callbacks.findVerifyUrls(act)
            default:
                //找不到激活邮件
                result = activator.ACTIVATE_EMAIL_NOT_FOUND
		}

		if globals.Exiting {
			result = activator.ACTIVATE_FAILED
			return
		}

		switch result {
            case activator.ACTIVATE_EMAIL_NOT_FOUND:
                //if acc.VerificationEmailResent == false {
                //    result = activator.ACTIVATE_CANT_ACTIVATE
                //
                //    e := Try(func() {
                //        result = act.ResendVerificationEmail()
                //    })
                //
                //    if e != nil {
                //        if _, isHttpError := e.Value.(*http.HttpError); isHttpError {
                //            Raise(e)
                //        }
                //
                //        logger.Debug("[%v] ResendVerificationEmail error: %v", acc.UserName, e)
                //    }
                //
                //    acc.VerificationEmailResent = true
                //}
                return

            case activator.ACTIVATE_SUCCESS:
                break

            default:
                return
		}

		result = activator.ACTIVATE_FAILED

        //新建一个代理，用于打开激活的apple url
		proxy := proxyManager.GetProxy(30 * time.Second)
		if proxy == nil {
			result = activator.ACTIVATE_PROXY_DISCONNECTED
			return
		}
		proxy.PreSignupLock(String(acc.UserName))

		verifyStartTime := time.Now()
		result = act.VerifyByUrl(links[len(links)-1], proxy)
		logger.Debug("[%v] open verify link took %v", acc.UserName, time.Now().Sub(verifyStartTime))
		logger.Debug("[%v] activate took %v", acc.UserName, time.Now().Sub(activateStartTime))

		if result == activator.ACTIVATE_SUCCESS {
			result = verifyAccount(act, proxyManager)
			logger.Debug("[%v] verifyAccount: %v", acc.UserName, result)
		}
		proxy.PostSignupUnlock(String(acc.UserName))
	})

	return
}

func activatorLoop(accountManager *account.Manager, proxyManager proxy.Manager) {
    //activatorLoop2里面有个死循环，用于不断检测需要激活的帐号
	if exp := Try(func() { activatorLoop2(accountManager, proxyManager) }); exp != nil {
		logger.Debug("activatorLoop exception:\n%v", exp)
	}

	logger.Debug("activatorLoop2 exit")
}

func activatorLoop2(accountManager *account.Manager, proxyManager proxy.Manager) {
	var acc *account.AppleAccount

	statusMap := map[activator.ActivateResult]int64{
		activator.ACTIVATE_SUCCESS:                account.ACCOUNT_ACTIVATED,
		activator.ACTIVATE_EMAIL_NOT_FOUND:        account.ACCOUNT_EMAIL_NOT_FOUND,
		activator.ACTIVATE_MAIL_LOGIN_FAILED:      account.ACCOUNT_LOGIN_FAILED,
		activator.ACTIVATE_ACCOUNT_INVALID:        account.ACCOUNT_EXISTS,
		activator.ACTIVATE_CANT_FINISH_REQUEST:    account.ACCOUNT_ACTIVATE_FAILED,
		activator.ACTIVATE_CANT_ACTIVATE:          account.ACCOUNT_CANT_ACTIVATE,
		activator.ACTIVATE_PASSWORD_ERROR:         account.ACCOUNT_PASSWORD_ERROR,
		activator.ACTIVATE_INVALID_PASSWORD_TOKEN: account.ACCOUNT_CANT_ACTIVATE,
		activator.ACTIVATE_ACCOUNT_BANNED:         account.ACCOUNT_SECURITY_BANNED,
	}

	retry := false

	logger.Debug("activatorLoop2++")
	for globals.Exiting == false {
		if retry == false {
			acc = accountManager.GetAccount(true)
			if acc == nil {
				logger.Debug("got nil account")
				break
			}
		}

        retry = false

        //elapsed := time.Now().Sub(time.Unix(0, acc.CreationTime*int64(time.Millisecond)))
        if time.Now().Unix() < acc.CreationTime + 120 {
            //如果帐号是刚注册2分钟内的话，将帐号放回到通道等待重新取出来，不再log。日志太多
			//logger.Debug("ReleaseAccount", acc.UserName)
            accountManager.ReleaseAccount(acc)
            continue
		}

        //跑激活流程，包括登录苹果和163邮箱
		result, exp := activate(acc, proxyManager)

		if exp != nil {
			switch err := exp.Value.(type) {
			case *http.HttpError:
				switch err.Type {
				case http.HTTP_ERROR_CANNOT_CONNECT,
					http.HTTP_ERROR_CONNECT_PROXY:
					logger.Debug("[%s] connect proxy timeout:\n%v", acc.UserName, err)
					result = activator.ACTIVATE_PROXY_DISCONNECTED

				case http.HTTP_ERROR_READ_ERROR,
					http.HTTP_ERROR_RESPONSE_ERROR:
					logger.Debug("[%s] network error:\n%v", acc.UserName, err)
					result = activator.ACTIVATE_NETWORK_ERROR

				default:
					logger.Debug("[%s] unknown http exception:\n%v", acc.UserName, exp)
					result = activator.ACTIVATE_NETWORK_ERROR
				}

			case *AttributeError:
				logger.Debug("[%s] AttributeError\n%v", acc.UserName, exp)
				result = activator.ACTIVATE_NETWORK_ERROR

			case *FileNotFoundError:
				logger.Debug("[%s] FileNotFoundError\n%v", acc.UserName, exp)
				result = activator.ACTIVATE_MAIL_LOGIN_FAILED

			case *mailmaster.PasswordError:
				logger.Debug("[%s] PasswordError\n%v", acc.UserName, exp)
				result = activator.ACTIVATE_MAIL_LOGIN_FAILED

			default:
				logger.Debug("[%s] unknown exception:\n%v", acc.UserName, exp)
				result = activator.ACTIVATE_FAILED
			}
		}

		logger.Debug("[%v] activate result: %v", acc.UserName, result)

		status, ok := statusMap[result]

		switch {
		case ok:
			switch status {
				case account.ACCOUNT_ACTIVATED:
					atomic.AddInt32(&successTotal, 1)

				case account.ACCOUNT_PASSWORD_ERROR,
					account.ACCOUNT_SECURITY_BANNED,
					account.ACCOUNT_CANT_ACTIVATE:
					atomic.AddInt32(&failureTotal, 1)
				}
			acc.Status = status
			accountManager.Commit(acc)

		case result == activator.ACTIVATE_DELAY_PROCESS:
			acc.CreationTime = globals.GetCurrentTime()
			accountManager.ReleaseAccount(acc)

		case result == activator.ACTIVATE_PROXY_DISCONNECTED,
			result == activator.ACTIVATE_NETWORK_ERROR,
			result == activator.ACTIVATE_SESSION_TIMEOUT,
			result == activator.ACTIVATE_MAIL_CONNECT_FAILED:
			atomic.AddInt32(&networkError, 1)
			retry = true
			logger.Debug("[%v] retry", acc.UserName)

		default:
			logger.Debug("[%v] unknown result: %v", acc.UserName, result)
			Raisef("[%v] unknown result: %v", acc.UserName, result)
		}
	}
	logger.Debug("activatorLoop2--")
}

func activatorRun(proxyManager proxy.Manager, accountManager *account.Manager) {
	defer dumpGoroutines("activator_callstack.txt")
	defer console.Pause("done")

	logger.LogToFile(true)
	sigc := make(chan os.Signal, 10)

	proxyManager.DisableCounter(true)

	signal.Notify(sigc, os.Interrupt)

	go func() {
		dumper := timer.NewTicker(time.Minute * 30)
		updater := timer.NewTicker(time.Second)

		second := 0
		speed := 0
		speeds := [60]int{}

		for {
			select {
			case <-dumper.C:
				dumpGoroutines("activator_callstack.txt")

			case <-updater.C:
				speed = int(successTotal)
				speeds[second] = speed
				second = (second + 1) % len(speeds)
				speed = (speed - speeds[second])

				// logger.Info(Sprintf("Success:%d Failure:%d Speed:%d", successTotal, failureTotal, speed))
				console.SetTitle(
					Sprintf("%v%d @ S:%d F:%d N:%d",
						If(globals.Exiting, "exiting ", ""),
						speed,
						successTotal,
						failureTotal,
						networkError,
					),
				)
			}
		}

		dumper.Stop()
		updater.Stop()
	}()

	go func() {
		for {
			_, ok := <-sigc
			if ok == false {
				break
			}

			dumpGoroutines("activator_callstack.txt")

			logger.Debug("exiting")
			globals.Exiting = true
		}
	}()

	wg := sync.WaitGroup{}

	// emailNotFoundManager := account.NewManager(account.EmailNotFoundManager)

	// emailNotFoundManager.Start()
	accountManager.Start()
	proxyManager.Start()

	startTime = time.Now()

	for i := 0; i < globals.Preferences.MaxActivatorWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			// logger.Debug("goroutine %d running", id)
			activatorLoop(accountManager, proxyManager)
			wg.Done()
		}(i)

		// wg.Add(1)
		// go func (id int) {
		//     activatorLoop(emailNotFoundManager, proxyManager)
		//     wg.Done()
		// }(i)
	}

    //直到上面的激活worker都关闭了才能执行
	wg.Wait()

	logger.Debug("all goroutine exit")

	proxyManager.Close()
	accountManager.Close()
	// emailNotFoundManager.Stop()

	signal.Stop(sigc)
}

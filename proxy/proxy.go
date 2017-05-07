package proxy

import (
	. "fmt"
	. "ml/strings"

	"time"

	"globals"
)

type Proxy struct {
	Host       String
	Port       int
	Port2      int
	User       String
	Password   String
	ExpireTime time.Time
    CallbackTime time.Time
	LeftCount  int
	// TODO: pending accounts map
	RegCount int
	JobID    String
}

func NewProxy(host String, alive int) *Proxy {
	p := &Proxy{
		Host:       host,
		Port:       globals.Preferences.Proxy.Port,
		User:       String(globals.Preferences.Proxy.User),
		Password:   String(globals.Preferences.Proxy.Password),
		ExpireTime: time.Now().Add(time.Duration(alive) * time.Second),
	}

	return p
}

func NewProxyWithPort(host String, port int, alive int) *Proxy {
	p := &Proxy{
		Host:       host,
		Port:       port,
		ExpireTime: time.Now().Add(time.Duration(alive) * time.Second),
	}

	return p
}

func NewProxyWithJobID(host String, leftCount int, jobID String) *Proxy {
	p := &Proxy{
		Host:       host,
		Port:       globals.Preferences.Proxy.Port,
		LeftCount:  leftCount,
		JobID:      jobID,
		ExpireTime: time.Now().Add(time.Duration(3*60) * time.Second), // 允许获取的时间
		CallbackTime: time.Now().Add(time.Duration(4*60) * time.Second), // 回收时间，比获取时间长，防止最后一个进行中被回收
	}

	return p
}

func (self *Proxy) String() string {
	//return Sprintf("host: %s alive: %d expire: %s", self.Host, self.ExpireTime.Sub(time.Now())/time.Second, self.ExpireTime.Format("2006-01-02 15:04:05"))
	return Sprintf("host: %s, jobid: %s, left: %d", self.Host, self.JobID, self.LeftCount)
}

func (self *Proxy) Expired() bool {
    //1. 用于start代理，检查是否超时，没有超时返回false
	return self.ExpireTime.Before(time.Now())
	//return self.LeftCount <= 0 && self.RegCount <= 0
}

func (self *Proxy) Callbacked() bool {
    //用于检查代理是否到回收时间,没有超时返回false
    return self.CallbackTime.Before(time.Now())
}

func (self *Proxy) PreSignupLock(account String) {
	self.RegCount++
}

func (self *Proxy) PostSignupUnlock(account String) {
	self.RegCount--
	if self.RegCount < 0 {
		self.RegCount = 0
	}
	if self.LeftCount < 0 {
		self.LeftCount = 0
	}
}

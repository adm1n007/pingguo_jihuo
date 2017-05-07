package proxy

import (
	// . "fmt"
	. "ml/array"
	. "ml/strings"
	. "ml/trace"

	"sync"
	"time"

	"ml/logging/logger"
	"ml/random"
	"ml/sync2"

	. "fmt"
	"ml/encoding/json"
	"ml/net/http"

	"globals"
)

type ManagerPP struct {
	Manager

	stop    bool
	running bool
	lock    *sync.Mutex
	event   *sync2.Event

	//getProxyPacket []byte

	pool    map[String]*Proxy
	counter map[String]int

	counterDisabled bool
}

type ProxyResp struct {
	Status int `json:"status,omitempty"`
	Data   struct {
		Job_id String `json:"job_id,omitempty"`
		IP     String `json:"ip,omitempty"`
	}
	Msg  String `json:"msg,omitempty"`
	Time int    `json:"time,omitempty"`
}

func NewManagerPP(proxyTag string) Manager {
	mgr := &ManagerPP{
		lock:  &sync.Mutex{},
		event: sync2.NewEvent(),

		//getProxyPacket: append([]byte("\x12\x00\x07\x00\x00\x00\x00\x00"), []byte(proxyTag)...),

		pool:    make(map[String]*Proxy),
		counter: make(map[String]int),

		counterDisabled: false,
	}

	//mgr.getProxyPacket[0] = byte(len(mgr.getProxyPacket))

	return mgr
}

func (self *ManagerPP) Close() {
	self.Stop()
}

func (self *ManagerPP) collect() {
	// logger.Debug("begin collect")
	if len(self.pool) > globals.Preferences.MaxProxys {
		return
	}


	//proxy := NewProxyWithJobID("localhost", 1, "localhost")
	//self.lock.Lock()
	//self.pool[proxy.Host] = proxy
	//self.lock.Unlock()
	//logger.Debug("got proxy %v", proxy)

	//return

	session := http.NewSession()
	link := Sprintf("%s/%s?action=%s&ex=%d&lo=%s&t=%d", globals.Preferences.ProxyServer.Url, "ip/get", "getjob", 300, "sh", int(time.Now().Unix()))
	resp := session.Request("GET", link)
	session.Close()

	if resp.StatusCode == http.StatusNotFound {
		logger.Debug("get ip failed!")
		//session.Close()
		return
	}
	text := resp.Text()
	logger.Debug("Getproxy api result: %v", text)

	var proxyjson ProxyResp
	json.LoadString(string(text), &proxyjson)
	//logger.Debug("parsed json: %v", proxyjson)

	if proxyjson.Status == 0 {
		//job_id 是内网ip
		proxy := NewProxyWithJobID(proxyjson.Data.IP, 1, proxyjson.Data.Job_id)

		self.lock.Lock()
		self.pool[proxy.Host] = proxy
		self.lock.Unlock()
	}
}

func (self *ManagerPP) remove() {
	// logger.Debug("remove expired proxy")

	self.lock.Lock()
	defer self.lock.Unlock()

	keys := []String{}
	jobs := []String{}
	for ip, proxy := range self.pool {
		if proxy.Callbacked() == false {
			continue
		}

		keys = append(keys, ip)
		jobs = append(jobs, proxy.JobID)
	}

	session := http.NewSession()
	for _, jobid := range jobs {
		if jobid == "localhost" {
			continue
		}
		// need unregister
		link := Sprintf("%s/%s?job_id=%s", globals.Preferences.ProxyServer.Url, "ip/discard", jobid)
		resp := session.Request("GET", link)
		if resp.StatusCode == http.StatusNotFound {
			continue
		}
		text := resp.Text()
		logger.Debug("Discard Ip api result :%v, jobId: %s", text, jobid)

		var proxyjson ProxyResp
		json.LoadString(string(text), &proxyjson)
		//logger.Debug("parsed json: %v", proxyjson)
	}
	session.Close()

	for _, ip := range keys {
		delete(self.pool, ip)
		delete(self.counter, ip)
	}
}

func (self *ManagerPP) mainLoop2() {
	// TODO: need a immediately call after each signup to check removal of expired proxy?
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for self.stop == false {
		if exp := Try(self.collect); exp != nil && exp.Value.(error).Error() != "EOF" {
			logger.Debug("collect error: %s", exp)
		}

		select {
		case <-ticker.C:
			self.remove()

		default:
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (self *ManagerPP) mainLoop() {
	self.mainLoop2()

	logger.Debug("mainLoop2 exit")
	self.event.Signal()
}

func (self *ManagerPP) Start() {
	if self.running {
		return
	}

	self.running = true
	self.stop = false
	go self.mainLoop()
}

func (self *ManagerPP) Stop() {
	if self.running == false {
		return
	}

	self.stop = true
	self.event.Wait()

	self.stop = false
	self.running = false
}

func (self *ManagerPP) GetProxy(minAlive time.Duration) *Proxy {
	//return &Proxy{Host: "localhost", Port: 8888, LeftCount: 4}
	if minAlive == -1 {
		minAlive = 30 * time.Second
	}

	// FIXME: force a ticker.Reset() instead invoke remove here

	self.remove()

    //不断循环,直到获取到代理为止
	for globals.Exiting == false {
		//now := time.Now()

		proxy := func() *Proxy {
			self.lock.Lock()
			defer self.lock.Unlock()

			proxyPool := Array{}

			// make(map[String]*Proxy) to array
			for _, proxy := range self.pool {
				proxyPool.Append(proxy)
			}

			for _, obj := range random.Shuffle(proxyPool) {
				// obj is array element
				//logger.Debug("type of obj: %T", obj)
				proxy := obj.(*Proxy)

				if proxy.Expired() {
					continue
				}

				//if proxy.ExpireTime.Sub(now) < minAlive {
				//	continue
				//}

				if self.counterDisabled {
					return proxy
				}

				if self.counter[proxy.Host] >= 1500 {
					continue
				}

				self.counter[proxy.Host] += 1
				return proxy
			}

			return nil
		}()

        //代理不为空，立即返回
		if proxy != nil {
			//proxy.LeftCount--
			return proxy
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (self *ManagerPP) DisableCounter(disable bool) {
	self.counterDisabled = disable
}

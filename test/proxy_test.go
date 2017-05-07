package test

import (
    . "ml/trace"
    "fmt"
    "testing"
    "time"
    "sync"
    "ml/net/http"
    "../proxy"
)

func TestProxyManager(t *testing.T) {
    proxymgr := proxy.NewManagerLocal()
    proxymgr.Start()
    defer proxymgr.Stop()

    proxymgr.DisableCounter(true)

    // wg := sync.WaitGroup{}

    // for i := 0; i != 256; i++ {
    //     wg.Add(1)
    //     go func() {
    //         h := http.NewSession()
    //         h.DefaultOptions.AutoRetry = true
    //         h.SetTimeout(15 * time.Second)
    //         for {
    //             fmt.Println("get proxy")
    //             // p := proxymgr.GetProxy(30 * time.Second)
    //             fmt.Println("got proxy")
    //             // h.SetProxy(p.Host, p.Port, p.User, p.Password)
    //             Try(func() { h.Get("http://www.apple.com/") })
    //             fmt.Println("got")
    //         }
    //         h.Close()
    //     }()
    // }

    // fmt.Println("fuck")
    // wg.Wait()

    // for {
    //     fmt.Printf("%+v\n", *proxymgr.GetProxy(-1))
    //     time.Sleep(time.Second)
    // }

    time.Sleep(time.Minute * 120)
}

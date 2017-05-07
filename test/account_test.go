package test

import (
    // . "fmt"
    "testing"
    "./account"
)

func TestProxyManager(t *testing.T) {
    mgr := account.NewManager(account.UnregisteredManager)
    defer mgr.Stop()
    mgr.Start()

    for i := 0; i != 100; i++ {
        acc := mgr.GetAccount(true)
        t.Log(acc)
    }
}

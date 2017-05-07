package proxy

import (
    "time"
)

type Manager interface {
    Close()
    Start()
    Stop()
    GetProxy(minAlive time.Duration) *Proxy
    DisableCounter(enable bool)
}

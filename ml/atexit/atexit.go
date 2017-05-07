package atexit

import (
    "sync"
)

var lock = &sync.Mutex{}
var cb = map[func()]bool{}

func Register(f func()) {
    lock.Lock()
    defer lock.Unlock()

    cb[f] = true
}

func UnRegister(f func()) {
    lock.Lock()
    defer lock.Unlock()

    delete(cb, f)
}

func Fire() {
    lock.Lock()
    defer lock.Unlock()

    for f, _ := range cb {
        f()
    }

    cb = map[func()]bool{}
}

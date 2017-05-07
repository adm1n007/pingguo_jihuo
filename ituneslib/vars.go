package ituneslib

import (
    "sync"
)

var lock = sync.Mutex{}
var deviceMap = map[uintptr]*IosDevice{}
var deviceNotifications = map[uintptr]func(device *IosDevice, remove bool){}

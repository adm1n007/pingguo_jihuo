package ituneslib

import (
    "reflect"
    "sync/atomic"
)

var itunesInitialized uintptr = 0

func Initialize() {
    if atomic.CompareAndSwapUintptr(&itunesInitialized, 0, 1) == false {
        return
    }

    itunesDllInitialize()
    itunes.Initialize.Call()
    sapInitialize()
}

func FreeSessionData(ptr interface{}) {
    itunes.iTunesFreeMemory.Call(reflect.ValueOf(ptr).Pointer())
}

func freeMachineData(ptr interface{}) {
    itunes.MachineDataFree.Call(reflect.ValueOf(ptr).Pointer())
}

func FreeMemory(ptr interface{}) {
    itunes.FreeMemory.Call(reflect.ValueOf(ptr).Pointer())
}

func init() {
    // base := syscall.MustLoadDLL(filepath.Join(os2.ExecutablePath(), "iTunesHelper.dll"))

    // t := reflect.TypeOf(itunes)
    // v := reflect.ValueOf(&itunes).Elem()
    // for i, n := 0, t.NumField(); i != n; i++ {
    //     name := t.Field(i).Name
    //     v.Field(i).Set(reflect.ValueOf(base.MustFindProc(name)))
    // }
}

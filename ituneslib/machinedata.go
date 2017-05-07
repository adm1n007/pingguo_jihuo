package ituneslib

import (
    . "ml/trace"
    "unsafe"
    "ml/os2"
)

type MachineDataSession struct {
    session     uintptr
}

func NewMachineDataSession() *MachineDataSession {
    md := &MachineDataSession{}
    md.initialize()
    return md
}

func (self *MachineDataSession) initialize() {

}

func (self *MachineDataSession) Close() {
    if self.session != 0 {
        itunes.MachineDataClose.Call(self.session)
        self.session = 0
    }
}

func (self *MachineDataSession) StartProvisioning(dsid int64, actionData []byte) (clientData []byte) {
    var buf *byte
    var size int
    var status int

    switch os2.PtrSize() {
        case os2.PtrSize_32bits:
            st, _, _ := itunes.MachineDataStartProvisioning.Call(
                            uintptr(dsid & 0xFFFFFFFF),
                            uintptr((dsid >> 32) & 0xFFFFFFFF),
                            bytesPtr(actionData),
                            bytesLen(actionData),
                            uintptr(unsafe.Pointer(&buf)),
                            uintptr(unsafe.Pointer(&size)),
                            uintptr(unsafe.Pointer(&self.session)),
                        )
            status = getStatus(st)

        case os2.PtrSize_64bits:
            st, _, _ := itunes.MachineDataStartProvisioning.Call(
                            uintptr(dsid),
                            bytesPtr(actionData),
                            bytesLen(actionData),
                            uintptr(unsafe.Pointer(&buf)),
                            uintptr(unsafe.Pointer(&size)),
                            uintptr(unsafe.Pointer(&self.session)),
                        )
            status = getStatus(st)
    }

    if status != 0 {
        Raise(newiTunesHelperErrorf("MD.StartProvisioning failed: %d", status))
    }

    defer freeMachineData(buf)
    clientData = toBytes(buf, size)
    return
}

func (self *MachineDataSession) FinishProvisioning(settingInfo, transportKey []byte) {
    st, _, _ := itunes.MachineDataFinishProvisioning.Call(
                    self.session,
                    bytesPtr(settingInfo),
                    bytesLen(settingInfo),
                    bytesPtr(transportKey),
                    bytesLen(transportKey),
                )
    status := getStatus(st)

    if status != 0 {
        Raise(newiTunesHelperErrorf("MD.FinishProvisioning failed: %d", status))
    }
}

func (self *MachineDataSession) GetData(dsid int64) (machineData, signature []byte) {
    var status int
    var data *byte
    var dataSize int
    var sig *byte
    var sigSize int

    switch os2.PtrSize() {
        case os2.PtrSize_32bits:
            st, _, _ := itunes.MachineDataGetData.Call(
                            uintptr(dsid & 0xFFFFFFFF),
                            uintptr((dsid >> 32) & 0xFFFFFFFF),
                            uintptr(unsafe.Pointer(&data)),
                            uintptr(unsafe.Pointer(&dataSize)),
                            uintptr(unsafe.Pointer(&sig)),
                            uintptr(unsafe.Pointer(&sigSize)),
                        )
            status = getStatus(st)

        case os2.PtrSize_64bits:
            st, _, _ := itunes.MachineDataGetData.Call(
                            uintptr(dsid),
                            uintptr(unsafe.Pointer(&data)),
                            uintptr(unsafe.Pointer(&dataSize)),
                            uintptr(unsafe.Pointer(&sig)),
                            uintptr(unsafe.Pointer(&sigSize)),
                        )
            status = getStatus(st)
    }

    if status != 0 {
        return
        // Raise(newiTunesHelperErrorf("MD.GetData failed: %d", status))
    }

    defer freeMachineData(data)
    defer freeMachineData(sig)

    machineData = toBytes(data, dataSize)
    signature = toBytes(sig, sigSize)

    return
}

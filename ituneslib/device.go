package ituneslib

import (
    . "fmt"
    "unsafe"
    "ml/syscall"
)

type IosDevice struct {
    device  uintptr
    ptr     uintptr

    ProductType     string
    DeviceType      string
    DeviceName      string
    DeviceClass     string
    ProductVersion  string
    CpuArchitecture string
    UniqueDeviceId  string
    Activated       bool
    JailBroken      bool
}

// com.apple.mobile.installation_proxy
// {'ClientOptions': {'ApplicationType': 'Any'}, 'Command': 'Browse'}

func newIosDevice(device uintptr) *IosDevice {
    dev := &IosDevice{
        device: device,
    }

    itunes.iOSDeviceCreate.Call(device, uintptr(unsafe.Pointer(&dev.ptr)))
    return dev
}

func (self *IosDevice) Close() {
    itunes.iOSDeviceClose.Call(self.ptr)
    self.ptr = 0
}

func (self *IosDevice) Refresh() {
    self.Connect()
    defer self.Disconnect()

    self.GetDeviceType()
    self.GetDeviceName()
    self.GetDeviceClass()
    self.GetProductVersion()
    self.GetCpuArchitecture()
    self.GetUniqueDeviceId()
    self.IsActivated()
    self.IsJailBroken()
}

func (self *IosDevice) Connect() bool {
    r, _, _ := itunes.iOSDeviceConnect.Call(self.ptr)
    return r == 0
}

func (self *IosDevice) Disconnect() {
    itunes.iOSDeviceDisconnect.Call(self.ptr)
}

func (self *IosDevice) getStringValue(proc syscall.Proc) string {
    var buffer *byte
    var status int32

    st, _, _ := proc.Call(self.ptr, uintptr(unsafe.Pointer(&buffer)))
    status = int32(st)
    if status != 0 {
        return Sprintf("<call '%q' failed: %08X>", proc, status)
    }

    defer FreeMemory(buffer)

    return toString(buffer)
}

func (self *IosDevice) getBytesValue(proc syscall.Proc) []byte {
    var buffer *byte
    var size int
    var status int32

    st, _, _ := proc.Call(self.ptr, uintptr(unsafe.Pointer(&buffer)), uintptr(unsafe.Pointer(&size)))
    status = int32(st)
    if status != 0 {
        return []byte(Sprintf("<call '%q' failed: %08X>", proc, status))
    }

    defer FreeMemory(buffer)

    return toBytes(buffer, size)
}

func (self *IosDevice) GetProductType() string {
    self.ProductType = self.getStringValue(itunes.iOSDeviceGetProductType)
    return self.ProductType
}

func (self *IosDevice) GetDeviceType() string {
    self.DeviceType = productTypeToName[self.GetProductType()]
    return self.DeviceType
}

func (self *IosDevice) GetDeviceName() string {
    self.DeviceName = self.getStringValue(itunes.iOSDeviceGetDeviceName)
    return self.DeviceName
}

func (self *IosDevice) GetDeviceClass() string {
    self.DeviceClass = self.getStringValue(itunes.iOSDeviceGetDeviceClass)
    return self.DeviceClass
}

func (self *IosDevice) GetProductVersion() string {
    self.ProductVersion = self.getStringValue(itunes.iOSDeviceGetProductVersion)
    return self.ProductVersion
}

func (self *IosDevice) IsJailBroken() bool {
    ret, _, _ := itunes.iOSDeviceIsJailBroken.Call(self.ptr)
    self.JailBroken = int32(ret) != 0
    return self.JailBroken
}

func (self *IosDevice) IsActivated() bool {
    self.Activated = self.getStringValue(itunes.iOSDeviceGetActivationState) == "Activated"
    return self.Activated
}

func (self *IosDevice) GetCpuArchitecture() string {
    self.CpuArchitecture = self.getStringValue(itunes.iOSDeviceGetCpuArchitecture)
    return self.CpuArchitecture
}

func (self *IosDevice) GetUniqueDeviceId() string {
    self.UniqueDeviceId = self.getStringValue(itunes.iOSDeviceGetUniqueDeviceID)
    return self.UniqueDeviceId
}

func (self *IosDevice) GetUniqueDeviceIDData() string {
    self.UniqueDeviceId = self.getStringValue(itunes.iOSDeviceGetUniqueDeviceIDData)
    return self.UniqueDeviceId
}

package ituneslib

import (
    _ "unsafe"
    "ml/syscall"
)

//go:linkname proc_Initialize proc_Initialize
//go:linkname proc_FreeMemory proc_FreeMemory
//go:linkname proc_iTunesFreeMemory proc_iTunesFreeMemory
//go:linkname proc_DeviceNotificationSubscribe proc_DeviceNotificationSubscribe
//go:linkname proc_DeviceWaitForDeviceConnectionChanged proc_DeviceWaitForDeviceConnectionChanged
//go:linkname proc_iOSDeviceCreate proc_iOSDeviceCreate
//go:linkname proc_iOSDeviceClose proc_iOSDeviceClose
//go:linkname proc_iOSDeviceConnect proc_iOSDeviceConnect
//go:linkname proc_iOSDeviceDisconnect proc_iOSDeviceDisconnect
//go:linkname proc_iOSDeviceGetProductType proc_iOSDeviceGetProductType
//go:linkname proc_iOSDeviceGetDeviceName proc_iOSDeviceGetDeviceName
//go:linkname proc_iOSDeviceGetDeviceClass proc_iOSDeviceGetDeviceClass
//go:linkname proc_iOSDeviceGetProductVersion proc_iOSDeviceGetProductVersion
//go:linkname proc_iOSDeviceGetCpuArchitecture proc_iOSDeviceGetCpuArchitecture
//go:linkname proc_iOSDeviceGetActivationState proc_iOSDeviceGetActivationState
//go:linkname proc_iOSDeviceGetUniqueDeviceID proc_iOSDeviceGetUniqueDeviceID
//go:linkname proc_iOSDeviceGetUniqueDeviceIDData proc_iOSDeviceGetUniqueDeviceIDData
//go:linkname proc_iOSDeviceIsJailBroken proc_iOSDeviceIsJailBroken
//go:linkname proc_iOSDeviceAuthorizeDsids proc_iOSDeviceAuthorizeDsids
//go:linkname proc_SapCreateSession proc_SapCreateSession
//go:linkname proc_SapCloseSession proc_SapCloseSession
//go:linkname proc_SapCreatePrimeSignature proc_SapCreatePrimeSignature
//go:linkname proc_SapVerifyPrimeSignature proc_SapVerifyPrimeSignature
//go:linkname proc_SapExchangeData proc_SapExchangeData
//go:linkname proc_SapSignData proc_SapSignData
//go:linkname proc_KbsyncCreateSession proc_KbsyncCreateSession
//go:linkname proc_KbsyncValidate proc_KbsyncValidate
//go:linkname proc_KbsyncGetData proc_KbsyncGetData
//go:linkname proc_KbsyncImport proc_KbsyncImport
//go:linkname proc_KbsyncCloseSession proc_KbsyncCloseSession
//go:linkname proc_KbsyncSaveDsid proc_KbsyncSaveDsid
//go:linkname proc_MachineDataStartProvisioning proc_MachineDataStartProvisioning
//go:linkname proc_MachineDataFinishProvisioning proc_MachineDataFinishProvisioning
//go:linkname proc_MachineDataFree proc_MachineDataFree
//go:linkname proc_MachineDataClose proc_MachineDataClose
//go:linkname proc_MachineDataGetData proc_MachineDataGetData
//go:linkname proc_EncryptJsSpToken proc_EncryptJsSpToken

type _itunesAPI struct {
    Initialize,
    FreeMemory,
    iTunesFreeMemory,

    DeviceNotificationSubscribe,
    DeviceWaitForDeviceConnectionChanged,

    iOSDeviceCreate,
    iOSDeviceClose,
    iOSDeviceConnect,
    iOSDeviceDisconnect,
    iOSDeviceGetProductType,
    iOSDeviceGetDeviceName,
    iOSDeviceGetDeviceClass,
    iOSDeviceGetProductVersion,
    iOSDeviceGetCpuArchitecture,
    iOSDeviceGetActivationState,
    iOSDeviceGetUniqueDeviceID,
    iOSDeviceGetUniqueDeviceIDData,
    iOSDeviceIsJailBroken,
    iOSDeviceAuthorizeDsids,

    SapCreateSession,
    SapCloseSession,
    SapCreatePrimeSignature,
    SapVerifyPrimeSignature,
    SapExchangeData,
    SapSignData,

    KbsyncCreateSession,
    KbsyncValidate,
    KbsyncGetData,
    KbsyncImport,
    KbsyncCloseSession,
    KbsyncSaveDsid,

    MachineDataStartProvisioning,
    MachineDataFinishProvisioning,
    MachineDataFree,
    MachineDataClose,
    MachineDataGetData syscall.Proc

    EncryptJsSpToken syscall.Proc
}

var itunes _itunesAPI

func itunesDllInitialize()  {
    itunes = _itunesAPI{
        proc_Initialize,
        proc_FreeMemory,
        proc_iTunesFreeMemory,
        proc_DeviceNotificationSubscribe,
        proc_DeviceWaitForDeviceConnectionChanged,
        proc_iOSDeviceCreate,
        proc_iOSDeviceClose,
        proc_iOSDeviceConnect,
        proc_iOSDeviceDisconnect,
        proc_iOSDeviceGetProductType,
        proc_iOSDeviceGetDeviceName,
        proc_iOSDeviceGetDeviceClass,
        proc_iOSDeviceGetProductVersion,
        proc_iOSDeviceGetCpuArchitecture,
        proc_iOSDeviceGetActivationState,
        proc_iOSDeviceGetUniqueDeviceID,
        proc_iOSDeviceGetUniqueDeviceIDData,
        proc_iOSDeviceIsJailBroken,
        proc_iOSDeviceAuthorizeDsids,
        proc_SapCreateSession,
        proc_SapCloseSession,
        proc_SapCreatePrimeSignature,
        proc_SapVerifyPrimeSignature,
        proc_SapExchangeData,
        proc_SapSignData,
        proc_KbsyncCreateSession,
        proc_KbsyncValidate,
        proc_KbsyncGetData,
        proc_KbsyncImport,
        proc_KbsyncCloseSession,
        proc_KbsyncSaveDsid,
        proc_MachineDataStartProvisioning,
        proc_MachineDataFinishProvisioning,
        proc_MachineDataFree,
        proc_MachineDataClose,
        proc_MachineDataGetData,

        proc_EncryptJsSpToken,
    }
}

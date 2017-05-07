package ituneslib

import (
    "ml/syscall"
)

//go:cgo_import_dynamic proc_Initialize Initialize "iTunesHelper.dll"
//go:cgo_import_dynamic proc_FreeMemory FreeMemory "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iTunesFreeMemory iTunesFreeMemory "iTunesHelper.dll"

//go:cgo_import_dynamic proc_DeviceNotificationSubscribe DeviceNotificationSubscribe "iTunesHelper.dll"
//go:cgo_import_dynamic proc_DeviceWaitForDeviceConnectionChanged DeviceWaitForDeviceConnectionChanged "iTunesHelper.dll"

//go:cgo_import_dynamic proc_iOSDeviceCreate iOSDeviceCreate "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceClose iOSDeviceClose "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceConnect iOSDeviceConnect "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceDisconnect iOSDeviceDisconnect "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetProductType iOSDeviceGetProductType "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetDeviceName iOSDeviceGetDeviceName "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetDeviceClass iOSDeviceGetDeviceClass "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetProductVersion iOSDeviceGetProductVersion "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetCpuArchitecture iOSDeviceGetCpuArchitecture "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetActivationState iOSDeviceGetActivationState "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetUniqueDeviceID iOSDeviceGetUniqueDeviceID "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceGetUniqueDeviceIDData iOSDeviceGetUniqueDeviceIDData "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceIsJailBroken iOSDeviceIsJailBroken "iTunesHelper.dll"
//go:cgo_import_dynamic proc_iOSDeviceAuthorizeDsids iOSDeviceAuthorizeDsids "iTunesHelper.dll"

//go:cgo_import_dynamic proc_SapCreateSession SapCreateSession "iTunesHelper.dll"
//go:cgo_import_dynamic proc_SapCloseSession SapCloseSession "iTunesHelper.dll"
//go:cgo_import_dynamic proc_SapCreatePrimeSignature SapCreatePrimeSignature "iTunesHelper.dll"
//go:cgo_import_dynamic proc_SapVerifyPrimeSignature SapVerifyPrimeSignature "iTunesHelper.dll"
//go:cgo_import_dynamic proc_SapExchangeData SapExchangeData "iTunesHelper.dll"
//go:cgo_import_dynamic proc_SapSignData SapSignData "iTunesHelper.dll"

//go:cgo_import_dynamic proc_KbsyncCreateSession KbsyncCreateSession "iTunesHelper.dll"
//go:cgo_import_dynamic proc_KbsyncValidate KbsyncValidate "iTunesHelper.dll"
//go:cgo_import_dynamic proc_KbsyncGetData KbsyncGetData "iTunesHelper.dll"
//go:cgo_import_dynamic proc_KbsyncImport KbsyncImport "iTunesHelper.dll"
//go:cgo_import_dynamic proc_KbsyncCloseSession KbsyncCloseSession "iTunesHelper.dll"
//go:cgo_import_dynamic proc_KbsyncSaveDsid KbsyncSaveDsid "iTunesHelper.dll"

//go:cgo_import_dynamic proc_MachineDataStartProvisioning MachineDataStartProvisioning "iTunesHelper.dll"
//go:cgo_import_dynamic proc_MachineDataFinishProvisioning MachineDataFinishProvisioning "iTunesHelper.dll"
//go:cgo_import_dynamic proc_MachineDataFree MachineDataFree "iTunesHelper.dll"
//go:cgo_import_dynamic proc_MachineDataClose MachineDataClose "iTunesHelper.dll"
//go:cgo_import_dynamic proc_MachineDataGetData MachineDataGetData "iTunesHelper.dll"

//go:cgo_import_dynamic proc_EncryptJsSpToken EncryptJsSpToken "iTunesHelper.dll"

var (
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
    proc_MachineDataGetData syscall.Proc

    proc_EncryptJsSpToken syscall.Proc
)

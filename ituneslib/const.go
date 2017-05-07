package ituneslib

const (
    STATE_CONNECT       = 1
    STATE_DISCONNECT    = 2
    STATE_UNSUBSCRIBE   = 3
)

const (
    AnonymousDsid   = int64(-1)
)

var productTypeToName = map[string]string{
    "iPad1,1"     : "iPad 1",
    "iPad2,1"     : "iPad 2 WiFi",
    "iPad2,2"     : "iPad 2 WiFi GSM",
    "iPad2,3"     : "iPad 2 WiFi CDMA",
    "iPad2,4"     : "iPad 2 WiFi(32nm)",
    "iPad2,5"     : "iPad mini WiFi",
    "iPad2,6"     : "iPad mini GSM",
    "iPad2,7"     : "iPad mini CDMA",
    "iPad3,1"     : "iPad 3 WiFi",
    "iPad3,2"     : "iPad 3 WiFi 4G CDMA",
    "iPad3,3"     : "iPad 3 WiFi 4G",
    "iPad3,4"     : "iPad 4 WiFi",
    "iPad3,5"     : "iPad 4 GSM",
    "iPad3,6"     : "iPad 4 CDMA",
    "iPad4,1"     : "iPad Air WiFi",
    "iPad4,2"     : "iPad Air WiFi+Cellular",
    "iPad4,3"     : "iPad Air",
    "iPad4,4"     : "iPad mini 2 WiFi",
    "iPad4,5"     : "iPad mini 2 WiFi+Cellular",
    "iPad4,6"     : "iPad mini 2",
    "iPad5,3"     : "iPad Air 2",
    "iPad5,4"     : "iPad Air 2",
    "iPad4,7"     : "iPad mini 3",
    "iPad4,8"     : "iPad mini 3",
    "iPad4,9"     : "iPad mini 3",
    "iPhone1,1"   : "iPhone 1",
    "iPhone1,2"   : "iPhone 3G",
    "iPhone2,1"   : "iPhone 3GS",
    "iPhone3,1"   : "iPhone 4",
    "iPhone3,2"   : "iPhone 4",
    "iPhone3,3"   : "iPhone 4 CDMA",
    "iPhone4,1"   : "iPhone 4S",
    "iPhone5,1"   : "iPhone 5",
    "iPhone5,2"   : "iPhone 5",
    "iPhone5,3"   : "iPhone 5c",
    "iPhone5,4"   : "iPhone 5c",
    "iPhone6,1"   : "iPhone 5s",
    "iPhone6,2"   : "iPhone 5s",
    "iPhone6,3"   : "iPhone 5s",
    "iPhone7,1"   : "iPhone 6 Plus",
    "iPhone7,2"   : "iPhone 6",
    "iPod1,1"     : "iPod touch 1G",
    "iPod2,1"     : "iPod touch 2G",
    "iPod3,1"     : "iPod touch 3G",
    "iPod4,1"     : "iPod touch 4G",
    "iPod5,1"     : "iPod touch 5G",
}

package ituneslib

import (
    "ml/uuid"
)

type FairPlayHWInfo struct {
    Length  int32
    Id      [20]byte
}

func NewRandomFairPlayHWInfo() *FairPlayHWInfo {
    deviceId := &FairPlayHWInfo{}
    deviceId.Length = 6

    u, _ := uuid.NewV4()

    copy(deviceId.Id[:], u[:6])

    return deviceId
}

type SapCertType uint

const (
    SAP_TYPE_REGISTER  = SapCertType(0xD2)
    SAP_TYPE_LOGIN     = SapCertType(0xC8)
)

func (self SapCertType) String() string {
    return map[SapCertType]string{
        SAP_TYPE_LOGIN      : "SAP_TYPE_LOGIN",
        SAP_TYPE_REGISTER   : "SAP_TYPE_REGISTER",
    }[self]
}

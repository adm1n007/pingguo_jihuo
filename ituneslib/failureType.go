package ituneslib

import (
    "fmt"
)

type FailureType int

const (
    Failure_Success          = FailureType(0)
    Failure_NoAccount        = FailureType(1001)
    Failure_NeedLogin        = FailureType(2072)
    Failure_TermsNotAccepted = FailureType(3038)
    Failure_CheckBillingInfo = FailureType(5001)
    Failure_ServerDialog     = FailureType(-1)
)

func (self FailureType) String() string {
    r, ok := map[FailureType]string{
                Failure_Success          : "Success",
                Failure_NoAccount        : "NoAccount",
                Failure_NeedLogin        : "NeedLogin",
                Failure_TermsNotAccepted : "TermsNotAccepted",
                Failure_CheckBillingInfo : "CheckBillingInfo",
                Failure_ServerDialog     : "ServerDialog",
            }[self]

    if ok {
        return r
    }

    return fmt.Sprintf("unknown FailureType: %d", self)
}

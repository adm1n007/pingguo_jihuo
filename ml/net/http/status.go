package http

import (
    httplib "net/http"
)

type HttpStatusCode int

const (

    StatusContinue                      = HttpStatusCode(httplib.StatusContinue)                        // 100
    StatusSwitchingProtocols            = HttpStatusCode(httplib.StatusSwitchingProtocols)              // 101

    StatusOK                            = HttpStatusCode(httplib.StatusOK)                              // 200
    StatusCreated                       = HttpStatusCode(httplib.StatusCreated)                         // 201
    StatusAccepted                      = HttpStatusCode(httplib.StatusAccepted)                        // 202
    StatusNonAuthoritativeInfo          = HttpStatusCode(httplib.StatusNonAuthoritativeInfo)            // 203
    StatusNoContent                     = HttpStatusCode(httplib.StatusNoContent)                       // 204
    StatusResetContent                  = HttpStatusCode(httplib.StatusResetContent)                    // 205
    StatusPartialContent                = HttpStatusCode(httplib.StatusPartialContent)                  // 206

    StatusMultipleChoices               = HttpStatusCode(httplib.StatusMultipleChoices)                 // 300
    StatusMovedPermanently              = HttpStatusCode(httplib.StatusMovedPermanently)                // 301
    StatusFound                         = HttpStatusCode(httplib.StatusFound)                           // 302
    StatusSeeOther                      = HttpStatusCode(httplib.StatusSeeOther)                        // 303
    StatusNotModified                   = HttpStatusCode(httplib.StatusNotModified)                     // 304
    StatusUseProxy                      = HttpStatusCode(httplib.StatusUseProxy)                        // 305
    StatusTemporaryRedirect             = HttpStatusCode(httplib.StatusTemporaryRedirect)               // 307

    StatusBadRequest                    = HttpStatusCode(httplib.StatusBadRequest)                      // 400
    StatusUnauthorized                  = HttpStatusCode(httplib.StatusUnauthorized)                    // 401
    StatusPaymentRequired               = HttpStatusCode(httplib.StatusPaymentRequired)                 // 402
    StatusForbidden                     = HttpStatusCode(httplib.StatusForbidden)                       // 403
    StatusNotFound                      = HttpStatusCode(httplib.StatusNotFound)                        // 404
    StatusMethodNotAllowed              = HttpStatusCode(httplib.StatusMethodNotAllowed)                // 405
    StatusNotAcceptable                 = HttpStatusCode(httplib.StatusNotAcceptable)                   // 406
    StatusProxyAuthRequired             = HttpStatusCode(httplib.StatusProxyAuthRequired)               // 407
    StatusRequestTimeout                = HttpStatusCode(httplib.StatusRequestTimeout)                  // 408
    StatusConflict                      = HttpStatusCode(httplib.StatusConflict)                        // 409
    StatusGone                          = HttpStatusCode(httplib.StatusGone)                            // 410
    StatusLengthRequired                = HttpStatusCode(httplib.StatusLengthRequired)                  // 411
    StatusPreconditionFailed            = HttpStatusCode(httplib.StatusPreconditionFailed)              // 412
    StatusRequestEntityTooLarge         = HttpStatusCode(httplib.StatusRequestEntityTooLarge)           // 413
    StatusRequestURITooLong             = HttpStatusCode(httplib.StatusRequestURITooLong)               // 414
    StatusUnsupportedMediaType          = HttpStatusCode(httplib.StatusUnsupportedMediaType)            // 415
    StatusRequestedRangeNotSatisfiable  = HttpStatusCode(httplib.StatusRequestedRangeNotSatisfiable)    // 416
    StatusExpectationFailed             = HttpStatusCode(httplib.StatusExpectationFailed)               // 417
    StatusTeapot                        = HttpStatusCode(httplib.StatusTeapot)                          // 418
    StatusPreconditionRequired          = HttpStatusCode(httplib.StatusPreconditionRequired)            // 428
    StatusTooManyRequests               = HttpStatusCode(httplib.StatusTooManyRequests)                 // 429
    StatusRequestHeaderFieldsTooLarge   = HttpStatusCode(httplib.StatusRequestHeaderFieldsTooLarge)     // 431
    StatusUnavailableForLegalReasons    = HttpStatusCode(httplib.StatusUnavailableForLegalReasons)      // 451
    StatusConferenceNotFound            = HttpStatusCode(452)

    StatusInternalServerError           = HttpStatusCode(httplib.StatusInternalServerError)             // 500
    StatusNotImplemented                = HttpStatusCode(httplib.StatusNotImplemented)                  // 501
    StatusBadGateway                    = HttpStatusCode(httplib.StatusBadGateway)                      // 502
    StatusServiceUnavailable            = HttpStatusCode(httplib.StatusServiceUnavailable)              // 503
    StatusGatewayTimeout                = HttpStatusCode(httplib.StatusGatewayTimeout)                  // 504
    StatusHTTPVersionNotSupported       = HttpStatusCode(httplib.StatusHTTPVersionNotSupported)         // 505
    StatusNetworkAuthenticationRequired = HttpStatusCode(httplib.StatusNetworkAuthenticationRequired)   // 511
)

func (self HttpStatusCode) String() string {
    return httplib.StatusText(int(self))
}

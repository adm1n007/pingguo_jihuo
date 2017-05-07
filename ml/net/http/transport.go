package http

import (
    httplib "net/http"
    "sync"
)

type Transport struct {
    *httplib.Transport
    cancelledRequests    map[*httplib.Request]bool
    requestsLock        sync.Mutex
}

func newTransport(transport *httplib.Transport) *Transport {
    return &Transport{
        Transport           : transport,
        cancelledRequests    : map[*httplib.Request]bool{},
        requestsLock        : sync.Mutex{},
    }
}

func (self *Transport) CancelRequest(request *httplib.Request) {
    self.Transport.CancelRequest(request)

    self.requestsLock.Lock()

    self.cancelledRequests[request] = true

    self.requestsLock.Unlock()
}

func (self *Transport) RemoveCancelledRequest(request *httplib.Request) bool {
    self.requestsLock.Lock()

    cancelled := self.cancelledRequests[request]
    delete(self.cancelledRequests, request)

    self.requestsLock.Unlock()

    return cancelled
}

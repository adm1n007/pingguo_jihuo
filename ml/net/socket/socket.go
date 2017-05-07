package socket

import (
    "net"
    "time"
)

type Auth struct {
    User string
    Password string
}

type Socket interface {
    SetSocks5Proxy(host string, port int, auth *Auth)
    Connect(host string, port int, timeout time.Duration)
    Read(n int) (buf []byte)
    ReadAll(n int) (buf []byte)
    Write(buf []byte) (n int)
    Close()

    LocalAddr() net.Addr
    RemoteAddr() net.Addr

    SetTimeout(t time.Duration)
    SetReadTimeout(t time.Duration)
    SetWriteTimeout(t time.Duration)

    SetDeadline(t time.Time)
    SetReadDeadline(t time.Time)
    SetWriteDeadline(t time.Time)
}

type Proxy interface {
    Connect(network, host string, port int, timeout time.Duration) (net.Conn, error)
}

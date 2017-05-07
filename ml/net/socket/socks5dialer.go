package socket

import (
    "fmt"
    "net"
    "time"
    "golang.org/x/net/proxy"
)

type Socks5Dialer struct {
    socks5 proxy.Dialer
    timeout time.Duration
}

func NewSocks5Dialer(network string, proxyAddress string, proxyPort int, auth *Auth) *Socks5Dialer {
    s := &Socks5Dialer{
        timeout: 0,
    }

    var err error
    var a *proxy.Auth

    if auth != nil {
        a = &proxy.Auth{
            User: auth.User,
            Password: auth.Password,
        }
    }

    s.socks5, err = proxy.SOCKS5(network, fmt.Sprintf("%s:%d", proxyAddress, proxyPort), a, s)
    RaiseSocketError(err)

    return s
}

func (self *Socks5Dialer) SetTimeout(timeout time.Duration) {
    if timeout >= 0 {
        self.timeout = timeout
    }
}

func (self *Socks5Dialer) Connect(network, host string, port int, timeout time.Duration) (net.Conn, error) {
    self.SetTimeout(timeout)
    return self.socks5.Dial(network, fmt.Sprintf("%s:%d", host, port))
}

func (self *Socks5Dialer) Dial(network, address string) (c net.Conn, err error) {
    if self.timeout <= 0 {
        c, err = net.Dial(network, address)
    } else {
        c, err = net.DialTimeout(network, address, self.timeout)
        if err == nil {
            c.SetDeadline(time.Now().Add(self.timeout * 2))
        }
    }

    return
}

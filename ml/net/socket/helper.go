package socket

import (
    "strings"
)

func mapHost(host string) string {
    switch strings.ToLower(host) {
        case "localhost":
            return "127.0.0.1"

        default:
            return host
    }
}

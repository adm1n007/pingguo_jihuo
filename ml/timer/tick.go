package timer

import (
    "time"
)

type Ticker struct {
    *time.Ticker
    Duration time.Duration
}

func NewTicker(d time.Duration) *Ticker {
    return &Ticker{
        Ticker  : time.NewTicker(d),
        Duration: d,
    }
}

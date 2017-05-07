package sync2

import (
    "sync"
)

type Semaphore struct {
    wg *sync.WaitGroup
}

func NewSemaphore() *Semaphore {
    return &Semaphore{
        &sync.WaitGroup{},
    }
}

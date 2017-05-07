package sync2

import (
    "sync"
)

type Event struct {
    cond *sync.Cond
}

func NewEvent() *Event {
    return &Event{
                cond : sync.NewCond(&sync.Mutex{}),
            }
}

func (self *Event) Wait() {
    self.cond.L.Lock()
    defer self.cond.L.Unlock()
    self.cond.Wait()
}

func (self *Event) Signal() {
    self.cond.Signal()
}

func (self *Event) Broadcast() {
    self.cond.Broadcast()
}

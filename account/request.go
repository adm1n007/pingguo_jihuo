package account

import (
    "active_apple/ml/sync2"
)



type CommitRequest struct {
    event       *sync2.Event
    account     *AppleAccount
}

func NewRequest(account *AppleAccount) *CommitRequest {
    return &CommitRequest{
                event   : sync2.NewEvent(),
                account : account,
            }
}

func (self *CommitRequest) Wait() {
    self.event.Wait()
}

func (self *CommitRequest) Done() {
    self.event.Signal()
}

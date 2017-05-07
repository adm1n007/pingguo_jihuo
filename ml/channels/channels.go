package channels

type Channel interface {
    In() chan <- interface{}
    Out() <-chan interface{}
    Close()
    Length() int
}

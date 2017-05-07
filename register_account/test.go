package main
import(
    "sync"
    "fmt"
    "time"
)

var l sync.Mutex
type Exception struct {
    Message     string
    Traceback   string
    Value       interface{}
}


func main() {
    var exp *Exception
    panic(exp)
}

func f() {
    l.Lock()
    fmt.Println("i am going")
    time.Sleep(300 * time.Millisecond)
    fmt.Println("hello world")
    l.Unlock()
}

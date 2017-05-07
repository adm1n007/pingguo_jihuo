package console

import (
    "fmt"
)

func SetTitle(text interface{}) {
    setTitle(fmt.Sprintf("%v", text))
}

func Pause(text ...interface{}) {
    if len(text) != 0 {
        fmt.Printf("%v\n", text[0])
    }

    pause()
}

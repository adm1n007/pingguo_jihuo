package console

import (
    "bufio"
    "os"
)

func pause() {
    bufio.NewReader(os.Stdin).ReadByte()
}

func setTitle(text string) {
}

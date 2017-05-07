package test

import (
    . "ml/dict"
    . "ml/strings"
    . "fmt"
    "../mailmaster"
    "testing"
    "encoding/hex"
)

func TestFuckYou(t *testing.T) {
    mm := mailmaster.NewMailMaster("oji537@163.com", "pk0364")
    defer mm.Close()
    mm.Initialize()
    mids := mm.ListMessages("00011230080543:")
    mailInfos := mm.GetMessageInfos(mids)

    for _, info := range mailInfos {
        Println(info)

        msg := mm.AsyncReadMessage(String(info.Get("mid").(string)), "both", true)
        text := String(msg.Get("text").(map[string]interface{})["content"].(string))

        for _, line := range text.SplitLines() {
            if line.Contains("https://id.apple.com") {
                Println(line)
            }
        }
    }
}

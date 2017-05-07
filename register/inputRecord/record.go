package inputRecord

import (
    . "ml/strings"
    . "ml/trace"
    "fmt"
    "time"
    "github.com/PuerkitoBio/goquery"
    "../../htmlutils"
)

type Record struct {
    Type        RecordType
    Elapsed     time.Duration

    Timestamp   int64

    X int
    Y int

    ElementFocused  String
    ElementHovered  String
    Inputs          []*goquery.Selection
}

func (self *Record) String() string {
    prefix := fmt.Sprintf("%s,%x,", self.Type.Tag(), int64(self.Elapsed / time.Millisecond))

    var suffix string

    switch self.Type {
        case Ncip:
            suffix = fmt.Sprintf("%x,1,1", self.Timestamp)

        case Start:
            st := []String{}
            for _, e := range self.Inputs {
                if e.Nodes[0].Data != "input" {
                    continue
                }

                name := htmlutils.GetInputName(e)
                length := String(htmlutils.GetInputValue(e)).Length()
                st = append(st, String(fmt.Sprintf("%s,%d", name, length)))
            }

            suffix = String(",").Join(st).String()

        case TimeSlice:
            suffix = fmt.Sprintf("%x", self.Timestamp)

        case FocusAndValueLength:
            suffix = fmt.Sprintf("0,%s", self.ElementFocused)

        case FocusOnly, FocusBlur:
            suffix = self.ElementFocused.String()

        case MouseMove, MouseClick:
            suffix = fmt.Sprintf("%x,%x", self.X, self.Y)
            if self.ElementHovered.IsEmpty() == false {
                suffix += "," + self.ElementHovered.String()
            }

        case KeyDown:
            break

        case Submit:
            suffix = "0,0,"

        default:
            Raise(NewNotImplementedError("%v", self.Type))
    }

    if len(suffix) == 0 {
        return prefix[:len(prefix) - 1]
    }

    return prefix + suffix
}

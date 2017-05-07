package inputRecord

type RecordType int

const (
    Ncip RecordType = iota
    TimeSlice
    Start
    FocusAndValueLength
    FocusOnly
    FocusBlur
    MouseMove
    MouseClick
    KeyDown
    Submit
    TouchStart
)

func (self RecordType) Tag() string {
    return map[RecordType]string{
        Ncip                : "ncip",
        TimeSlice           : "ts",
        Start               : "st",
        FocusAndValueLength : "kk",
        FocusOnly           : "ff",
        FocusBlur           : "fb",
        MouseMove           : "mm",
        MouseClick          : "mc",
        KeyDown             : "kd",
        Submit              : "fs",
        TouchStart          : "te",
    }[self]
}

func (self RecordType) String() string {
    return map[RecordType]string{
        Ncip                : "Ncip",
        TimeSlice           : "TimeSlice",
        Start               : "Start",
        FocusAndValueLength : "FocusAndValueLength",
        FocusOnly           : "FocusOnly",
        MouseMove           : "MouseMove",
        MouseClick          : "MouseClick",
        KeyDown             : "KeyDown",
        FocusBlur           : "FocusBlur",
        Submit              : "Submit",
        TouchStart          : "TouchStart",
    }[self]
}

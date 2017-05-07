package ml

func If(cond bool, True interface{}, False interface{}) interface{} {
    if cond {
        return True
    }

    return False
}

func FlagOn(flags, bit int) bool {
    return (flags & bit) != 0
}

func FlagOff(flags, bit int) bool {
    return !FlagOn(flags, bit)
}

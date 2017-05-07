package trace

func Try(f func()) (exp *Exception) {
    defer func() {
        exp = Catch(recover())
    }()
    f()

    return
}

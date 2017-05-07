package trace

type config struct {
    ReadSource      bool
}

var Config = config{
    ReadSource: true,
}

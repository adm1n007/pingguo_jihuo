package logging

import (
    "fmt"
    "io"
    "os"
    "sync"
    "time"
    "runtime"
    . "ml/strings"
    "ml/os2"
    "path/filepath"
)

const (
    CRITICAL = 50
    FATAL    = CRITICAL
    ERROR    = 40
    WARNING  = 30
    WARN     = WARNING
    INFO     = 20
    DEBUG    = 10
    NOTSET   = 0

    skipFuncs   = 3
)

var levelToName = map[int]string {
    CRITICAL    : "CRITICAL",
    ERROR       : "ERROR",
    WARNING     : "WARNING",
    INFO        : "INFO",
    DEBUG       : "DEBUG",
    NOTSET      : "NOTSET",
}

var nameToLevel = map[string]int {
    "CRITICAL"  : CRITICAL,
    "ERROR"     : ERROR,
    "WARN"      : WARNING,
    "WARNING"   : WARNING,
    "INFO"      : INFO,
    "DEBUG"     : DEBUG,
    "NOTSET"    : NOTSET,
}

var lock = sync.Mutex{}
var defaultFormatter = String("[%date %time][%tag][%file][%func:%line]")

type Logger struct {
    formatter   String
    level       int
    callDepth   int
    skip        int
    tag         String
    lock        sync.Mutex
    out         []io.WriteCloser
}

func NewLogger(tag interface{}) *Logger {
    return &Logger{
               formatter    : defaultFormatter,
               level        : DEBUG,
               tag          : String(fmt.Sprint(tag)),
               callDepth    : 10,
               skip         : skipFuncs,
               out          : []io.WriteCloser{os.Stdout},
            }
}

func (self *Logger) String() string {
    return fmt.Sprintf("%s logger", self.tag)
}

func (self *Logger) getCaller(skip int) (pc uintptr, name, file string, line int, ok bool) {
    for {

        pc, file, line, ok = runtime.Caller(skip)
        if ok == false {
            break
        }

        names := String(runtime.FuncForPC(pc).Name()).RSplit(".", 1)
        n := names[len(names) - 1]

        if String("debug.info.warning.error.fatal").Contains(n.ToLower()) == false {
            name = string(n)
            break
        }

        skip++
    }

    return
}

func (self *Logger) output(level int, format interface{}, args ...interface{}) {
    if level < self.level {
        return
    }

    lock.Lock()
    defer lock.Unlock()

    t := time.Now()

    // pc, file, line, ok := runtime.Caller(self.skip)
    _, name, file, line, ok := self.getCaller(self.skip)
    if ok == false {
        file = "???"
        name = "???"
        line = -1
        // pc = 0
    }

    text := fmt.Sprintf(fmt.Sprintf("%v", format), args...)

    formatter := self.formatter

    if formatter.Find("%tag") != -1 {
        formatter = formatter.Replace("%tag", self.tag)
    }

    if formatter.Find("%file") != -1 {
        formatter = formatter.Replace("%file", String(filepath.Base(file)))
    }

    if formatter.Find("%func") != -1 {
        // names := String(runtime.FuncForPC(pc).Name()).RSplit(".", 1)
        // formatter = formatter.Replace("%func", names[len(names) - 1])
        formatter = formatter.Replace("%func", String(name))
    }

    if formatter.Find("%line") != -1 {
        formatter = formatter.Replace("%line", String(fmt.Sprintf("%d", line)))
    }

    if formatter.Find("%date") != -1 {
        year, month, day := t.Date()
        date := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

        formatter = formatter.Replace("%date", String(date))
    }

    if formatter.Find("%time") != -1 {
        hour, min, sec := t.Clock()
        time := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)

        formatter = formatter.Replace("%time", String(time))
    }

    formatter += "[" + String(levelToName[level]) + "] "

    buf := []byte(string(formatter) + text + "\n")

    for _, out := range self.out {
        out.Write(buf)
    }
}

func (self *Logger) Debug(format interface{}, args ...interface{}) {
    self.output(DEBUG, format, args...)
}

func (self *Logger) Info(format interface{}, args ...interface{}) {
    self.output(INFO, format, args...)
}

func (self *Logger) Warning(format interface{}, args ...interface{}) {
    self.output(WARNING, format, args...)
}

func (self *Logger) Error(format interface{}, args ...interface{}) {
    self.output(ERROR, format, args...)
}

func (self *Logger) Fatal(format interface{}, args ...interface{}) {
    self.output(FATAL, format, args...)
}

func (self *Logger) SetLevel(level int) {
    self.level = level
}

func (self *Logger) Level() int {
    return self.level
}

func (self *Logger) SetSkip(skip int) {
    self.skip = skip + skipFuncs
}

func (self *Logger) SetFormater(formatter ...String) {
    switch len(formatter) {
        case 0:
            self.formatter = defaultFormatter
        default:
            self.formatter = formatter[0]
    }
}

func (self *Logger) LogToFile(enable bool, path ...String) error {
    var filename string
    var output io.WriteCloser = nil

    self.lock.Lock()
    defer self.lock.Unlock()

    if enable == false {
        if len(self.out) == 2 {
            self.removeOutput(self.out[1])
        }

        return nil
    }

    switch len(path) {
        case 0:
            filename = self.getDefaultFileName()
        default:
            filename = string(path[0])
    }

    if len(self.out) == 2 {
        self.removeOutput(self.out[1])
    }

    output, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
    if err != nil {
        return err
    }

    output.Write([]byte{0xEF, 0xBB, 0xBF})
    self.out = append(self.out, output)

    return nil
}

func (self *Logger) getDefaultFileName() string {
    now := time.Now()
    year, month, day := now.Date()
    hour, minute, second := now.Clock()
    exe := os2.Executable()
    exeName := filepath.Base(exe)

    baseName := fmt.Sprintf("[%s][%s][%04d-%02d-%02d %02d.%02d.%02d][%d].txt", exeName[:len(exeName) - len(filepath.Ext(exeName))], self.tag, year, month, day, hour, minute, second, os.Getpid())

    logPath := filepath.Join(filepath.Dir(exe), "logs")
    err := os.MkdirAll(logPath, 0666)
    if err != nil {
        return ""
    }

    return filepath.Join(logPath, baseName)
}

func (self *Logger) removeOutput(output io.WriteCloser) {
    for i := range self.out {
        if self.out[i] != output {
            continue
        }

        switch output {
            case os.Stdin:
            case os.Stdout:
            case os.Stderr:

            default:
                output.Close()
        }

        self.out = append(self.out[:i], self.out[i + 1:]...)
        break
    }
}
